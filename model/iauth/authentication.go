// Copyright (c) 2021 The BFE Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/itxn"
	"github.com/bfenetworks/api-server/stateful"
)

type key string

var keyUser key = "user"

var (
	UserTypeNormal int8 = 0
	UserTypeToken  int8 = 1
)

const (
	AuthTypePassword   = "Password"
	AuthTypeSessionKey = "Session"
	AuthTypeToken      = "Token"
	AuthTypeSkip       = "Skip"
)

type Loginer interface {
	GetName() string
	GetScopes() []string
	GetType() int8
	IsAdmin() bool
}

var (
	_ Loginer = &Visitor{}
	_ Loginer = &User{}
	_ Loginer = &Token{}
)

type Visitor struct {
	User  *User
	Token *Token
}

func (v *Visitor) GetName() string {
	if v.User != nil {
		return v.User.GetName()
	}

	return v.Token.GetName()
}

func (v *Visitor) GetScopes() []string {
	if v.User != nil {
		return v.User.GetScopes()
	}

	return v.Token.GetScopes()
}

func (v *Visitor) GetType() int8 {
	if v.User != nil {
		return v.User.GetType()
	}

	return v.Token.GetType()
}

func (v *Visitor) IsAdmin() bool {
	if v.User != nil {
		return v.User.IsAdmin()
	}

	return v.Token.IsAdmin()
}

type User struct {
	ID                 int64
	Name               string
	Type               int8
	Admin              bool
	Password           string
	SessionKey         string
	SessionKeyCreateAt time.Time
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetScopes() []string {
	if u.Admin {
		return []string{ScopeSystem}
	}

	return []string{ScopeProduct}
}

func (u *User) GetType() int8 {
	return u.Type
}

func (u *User) IsAdmin() bool {
	return u.Admin
}

type Token struct {
	ID    int64
	Name  string
	Token string
	Scope string

	Product *ibasic.Product
}

func (t *Token) GetName() string {
	return t.Name
}

func (t *Token) GetScopes() []string {
	return []string{t.Scope}
}

func (t *Token) GetType() int8 {
	return UserTypeToken
}

func (t *Token) IsAdmin() bool {
	return t.Scope == ScopeSystem
}

func NewVisitorContext(ctx context.Context, visitor *Visitor) context.Context {
	return context.WithValue(ctx, keyUser, visitor)
}

func MustGetVisitor(ctx context.Context) (*Visitor, error) {
	obj := ctx.Value(keyUser)
	if obj != nil {
		return obj.(*Visitor), nil
	}

	if stateful.DefaultConfig.RunTime.SkipTokenValidate {
		return newFakeVisitor(ScopeSystem), nil
	}

	return nil, xerror.WrapAuthenticateFailErrorWithMsg("User Not Login")
}

func newFakeVisitor(scope string) *Visitor {
	return &Visitor{
		Token: &Token{
			ID:    1,
			Name:  "SkipUser",
			Scope: scope,
		},
	}
}

type UserParam struct {
	Name               *string
	Password           *string
	Scopes             []string
	SessionKey         *string
	SessionKeyCreateAt *time.Time
}

type UserFilter struct {
	IDs        []int64
	Name       *string
	SessionKey *string
	Type       *int8
	Types      []int8
}

type AuthenticateParam struct {
	Type     string
	Identify string
	Extend   string
}

type AuthenticateStorager interface {
	FetchUserList(ctx context.Context, param *UserFilter) ([]*User, error)
	FetchUser(ctx context.Context, param *UserFilter) (*User, error)
	UpdateUser(ctx context.Context, user *User, param *UserParam) error
	CreateUser(ctx context.Context, param *UserParam) error
	DeleteUser(ctx context.Context, user *User) error

	FetchTokens(ctx context.Context, param *TokenFilter) ([]*Token, error)
	CreateToken(ctx context.Context, token *TokenParam) error
	DeleteToken(ctx context.Context, param *Token) error
}

type AuthenticateManager struct {
	txn               itxn.TxnStorager
	storager          AuthenticateStorager // the storager for authentication
	authorizeStorager AuthorizeStorager    // the storager for authorization
}

func NewAuthenticateManager(txn itxn.TxnStorager, storager AuthenticateStorager,
	authorizeStorage AuthorizeStorager) *AuthenticateManager {
	return &AuthenticateManager{
		txn:               txn,
		storager:          storager,
		authorizeStorager: authorizeStorage,
	}
}

func sessionKeyFactory(len int) (string, error) {
	bs := make([]byte, len)
	_, err := rand.Read(bs)
	if err != nil {
		return "", xerror.WrapModelErrorWithMsg("rand.Read data error: %s", err.Error())
	}

	return base64.URLEncoding.EncodeToString(bs), nil
}

func authTypePassword(ctx context.Context, param *AuthenticateParam, manager *AuthenticateManager) (v *Visitor, err error) {

	err = manager.txn.AtomExecute(ctx, func(ctx context.Context) error {
		userName := param.Identify
		user, err := manager.storager.FetchUser(ctx, &UserFilter{
			Name: &userName,
		})
		if err != nil {
			return err
		}

		if user == nil {
			return xerror.WrapAuthenticateFailErrorWithMsg("User %s Not Exist", userName)
		}

		if param.Extend != "SKIP" && user.Password != param.Extend {
			return xerror.WrapAuthenticateFailErrorWithMsg("Password Wrong")
		}

		// update session key
		for {
			sessionKey, err := sessionKeyFactory(15)
			if err != nil {
				return err
			}
			user1, err := manager.storager.FetchUser(ctx, &UserFilter{
				SessionKey: &sessionKey,
			})
			if err != nil {
				return err
			}

			if user1 != nil {
				continue
			}

			user.SessionKey = sessionKey

			v = &Visitor{
				User: user,
			}

			return manager.storager.UpdateUser(ctx, user, &UserParam{
				SessionKey:         &sessionKey,
				SessionKeyCreateAt: lib.PTimeNow(),
			})
		}
	})

	return
}

var Authenticators = map[string]func(ctx context.Context, param *AuthenticateParam, manager *AuthenticateManager) (*Visitor, error){
	AuthTypePassword: authTypePassword,

	AuthTypeSessionKey: func(ctx context.Context, param *AuthenticateParam, manager *AuthenticateManager) (v *Visitor, err error) {
		user, err := manager.storager.FetchUser(ctx, &UserFilter{
			SessionKey: &param.Identify,
		})
		if err != nil {
			return nil, err
		}

		if user == nil {
			return nil, xerror.WrapAuthenticateFailErrorWithMsg("Session Key Wrong")
		}

		if user.SessionKeyCreateAt.AddDate(0, 0, stateful.DefaultConfig.RunTime.SessionExpireDay).Before(time.Now()) {
			return nil, xerror.WrapAuthenticateFailErrorWithMsg("Session Key Expired")
		}

		return &Visitor{
			User: user,
		}, nil
	},

	AuthTypeToken: func(ctx context.Context, param *AuthenticateParam, manager *AuthenticateManager) (v *Visitor, err error) {
		tokens, err := manager.storager.FetchTokens(ctx, &TokenFilter{
			Token: &param.Identify,
		})
		if err != nil {
			return nil, err
		}

		if len(tokens) == 0 {
			return nil, xerror.WrapAuthenticateFailErrorWithMsg("Token Wrong")
		}

		return &Visitor{
			Token: tokens[0],
		}, nil
	},

	AuthTypeSkip: func(ctx context.Context, param *AuthenticateParam, manager *AuthenticateManager) (v *Visitor, err error) {
		if !stateful.DefaultConfig.RunTime.SkipTokenValidate {
			return nil, xerror.WrapAuthenticateFailErrorWithMsg("Bad Authorization Flag")
		}

		return newFakeVisitor(param.Identify), nil
	},
}

func (m *AuthenticateManager) Authenticate(ctx context.Context, param *AuthenticateParam) (v *Visitor, err error) {
	handler := Authenticators[param.Type]
	if handler == nil {
		return nil, xerror.WrapParamErrorWithMsg("Illegal Authenticate Type %s", param.Type)
	}

	return handler(ctx, param, m)
}

func (m *AuthenticateManager) DestroySessionKey(ctx context.Context, sessionKey string) (err error) {
	return m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		user, err := m.storager.FetchUser(ctx, &UserFilter{
			SessionKey: &sessionKey,
		})
		if err != nil {
			return err
		}

		if user == nil {
			return xerror.WrapAuthenticateFailErrorWithMsg("Session Key Not Exist")
		}

		return m.storager.UpdateUser(ctx, user, &UserParam{
			SessionKey:         lib.PString(""),
			SessionKeyCreateAt: lib.PTime(time.Time{}.AddDate(0, 1, 1)),
		})
	})
}

type TokenFilter struct {
	IDs   []int64
	Name  *string
	Token *string
}

type TokenParam struct {
	Name  *string
	Token *string
	Scope *string
}

func (m *AuthenticateManager) FetchTokens(ctx context.Context, filter *TokenFilter) (list []*Token, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err = m.storager.FetchTokens(ctx, filter)
		if err != nil {
			return err
		}

		products, err := m.authorizeStorager.BatchFetchTokenProduct(ctx, list)
		if err != nil {
			return err
		}

		for _, one := range list {
			one.Product = products[one.ID]
		}
		return err
	})

	return
}

func (m *AuthenticateManager) DeleteToken(ctx context.Context, token *Token) (err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		err = m.storager.DeleteToken(ctx, token)
		if err != nil {
			return err
		}

		return m.authorizeStorager.UnbindTokenAllProduct(ctx, token)
	})

	return err
}

func (m *AuthenticateManager) CreateToken(ctx context.Context, param *TokenParam, product *ibasic.Product) (token *Token, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		tokens, err := m.storager.FetchTokens(ctx, &TokenFilter{
			Name: param.Name,
		})
		if len(tokens) > 0 {
			return xerror.WrapModelErrorWithMsg("Token Existed")
		}

		var tokenVal string

		for i := 0; i < 5; i++ {
			tokenVal, err = sessionKeyFactory(15)
			if err != nil {
				continue
			}

			tokens, err = m.storager.FetchTokens(ctx, &TokenFilter{
				Token: &tokenVal,
			})
			if err != nil {
				continue
			}
			if len(tokens) > 0 {
				continue
			}

			token = &Token{
				Name:  *param.Name,
				Token: tokenVal,
				Scope: *param.Scope,
			}

			err = m.storager.CreateToken(ctx, &TokenParam{
				Name:  param.Name,
				Token: &tokenVal,
				Scope: param.Scope,
			})
			if err != nil {
				return err
			}

			tokens, err = m.storager.FetchTokens(ctx, &TokenFilter{
				Name: param.Name,
			})
			if err != nil {
				return err
			}

			token = tokens[0]

			if product != nil {
				err = m.authorizeStorager.BindTokenProduct(ctx, token, product)
			}
			return err
		}

		return err
	})

	return
}

func (m *AuthenticateManager) CreateUser(ctx context.Context, param *UserParam) (err error) {
	if param.Password != nil {
		if err = passwordCheck(*param.Password); err != nil {
			return err
		}
	}

	return m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		user, err := m.storager.FetchUser(ctx, &UserFilter{
			Name: param.Name,
		})
		if err != nil {
			return err
		}

		if user != nil {
			return xerror.WrapModelErrorWithMsg("User Existed")
		}

		return m.storager.CreateUser(ctx, param)
	})
}

func (m *AuthenticateManager) DeleteUser(ctx context.Context, userName string) (err error) {
	return m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		user, err := m.storager.FetchUser(ctx, &UserFilter{
			Name: &userName,
		})
		if err != nil {
			return err
		}

		if user == nil {
			return xerror.WrapRecordNotExist("User")
		}

		if err = m.storager.DeleteUser(ctx, user); err != nil {
			return err
		}

		return m.authorizeStorager.UnbindUserAllProduct(ctx, user)
	})
}

func passwordCheck(password string) error {
	if len(password) < 6 {
		return xerror.WrapParamErrorWithMsg("Password Lenght Must Bigger Than 6")
	}

	return nil
}

type PasswordChangeData struct {
	UserName    string
	OldPassword string
	Password    string
}

func (m *AuthenticateManager) UpdateUserPassword(ctx context.Context, pcd *PasswordChangeData) (err error) {
	if err = passwordCheck(pcd.Password); err != nil {
		return err
	}

	return m.updateUser(ctx, &UserFilter{
		Name: &pcd.UserName,
	}, func(user *User) error {
		if pcd.OldPassword != "" {
			if user.Password != pcd.OldPassword {
				return xerror.WrapModelErrorWithMsg("Invalid Password")
			}
		}
		return nil
	}, &UserParam{
		Password:           &pcd.Password,
		SessionKey:         lib.PString(""),
		SessionKeyCreateAt: lib.PTime(time.Time{}.AddDate(0, 1, 1)),
	})
}

func (m *AuthenticateManager) updateUser(ctx context.Context, filter *UserFilter,
	userChecker func(*User) error, newData *UserParam) (err error) {

	return m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		user, err := m.storager.FetchUser(ctx, filter)
		if err != nil {
			return err
		}

		if user == nil {
			return xerror.WrapModelErrorWithMsg("User Not Exist")
		}

		if userChecker != nil {
			if err = userChecker(user); err != nil {
				return err
			}
		}

		return m.storager.UpdateUser(ctx, user, newData)
	})
}

func (m *AuthenticateManager) FetchUser(ctx context.Context, filter *UserFilter) (user *User, err error) {
	list, err := m.FetchUserList(ctx, filter)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		user = list[0]
	}

	return
}

func (m *AuthenticateManager) FetchUserList(ctx context.Context, param *UserFilter) (users []*User, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		users, err = m.storager.FetchUserList(ctx, param)

		return err
	})

	return
}
