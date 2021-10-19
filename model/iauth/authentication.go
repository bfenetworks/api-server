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
	"github.com/bfenetworks/api-server/model/itxn"
	"github.com/bfenetworks/api-server/stateful"
)

type key string

var keyUser key = "user"

func NewUserContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, keyUser, user)
}

func MustGetUser(ctx context.Context) (*User, error) {
	obj := ctx.Value(keyUser)
	if obj != nil {
		return obj.(*User), nil
	}

	if stateful.DefaultConfig.RunTime.SkipTokenValidate {
		return newFakeUser(&Role{
			Name: RoleNameAdmin,
		}), nil
	}

	return nil, xerror.WrapAuthenticateFailErrorWithMsg("User Not Login")
}

func newFakeUser(role *Role) *User {
	return &User{
		ID:    1,
		Name:  "SkipUser",
		Roles: []*Role{role},
	}
}

type User struct {
	ID              int64
	Name            string
	Password        string
	SessionKey      string
	SessionCreateAt time.Time

	Roles []*Role
}

type UserFilter struct {
	Name       *string
	SessionKey *string
}

type UserParam struct {
	Name            *string
	SessionKey      *string
	Password        *string
	Roles           []*Role
	SessionCreateAt *time.Time
}

func (u *User) IsAdmin() bool {
	for _, one := range u.Roles {
		if one.Name == RoleNameAdmin {
			return true
		}
	}
	return false
}

func (u *User) IsInner() bool {
	for _, one := range u.Roles {
		if one.Name == RoleNameInner {
			return true
		}
	}
	return false
}

const (
	RoleNameAdmin   = "admin"
	RoleNameProduct = "product"
	RoleNameInner   = "inner"
)

var AllowRoles = map[string]bool{
	RoleNameAdmin:   true,
	RoleNameProduct: true,
	RoleNameInner:   true,
}

var RoleInner = &Role{
	Name: RoleNameInner,
}

type Role struct {
	Name string
}

const (
	AuthTypePassword   = "Password"
	AuthTypeSessionKey = "Session"
	AuthTypeSkip       = "Skip"
)

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
}

func RoleList(rs []string) ([]*Role, error) {
	var roles []*Role
	for _, r := range rs {
		if !AllowRoles[r] {
			return nil, xerror.WrapParamErrorWithMsg("Role %s illegal", r)
		}

		roles = append(roles, &Role{
			Name: r,
		})
	}

	return roles, nil
}

type AuthenticateManager struct {
	storager AuthenticateStorager
	txn      itxn.TxnStorager
}

func NewAuthenticateManager(txn itxn.TxnStorager, storager AuthenticateStorager) *AuthenticateManager {
	return &AuthenticateManager{
		txn:      txn,
		storager: storager,
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

func authTypePassword(ctx context.Context, param *AuthenticateParam, manager *AuthenticateManager) (user *User, err error) {

	err = manager.txn.AtomExecute(ctx, func(ctx context.Context) error {
		userName := param.Identify
		user, err = manager.storager.FetchUser(ctx, &UserFilter{
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

			return manager.storager.UpdateUser(ctx, user, &UserParam{
				SessionKey:      &sessionKey,
				SessionCreateAt: lib.PTimeNow(),
			})
		}
	})

	return
}

var Authenticators = map[string]func(ctx context.Context, param *AuthenticateParam, manager *AuthenticateManager) (*User, error){
	AuthTypePassword: authTypePassword,

	AuthTypeSessionKey: func(ctx context.Context, param *AuthenticateParam, manager *AuthenticateManager) (user *User, err error) {
		user, err = manager.storager.FetchUser(ctx, &UserFilter{
			SessionKey: &param.Identify,
		})
		if err != nil {
			return nil, err
		}

		if user == nil {
			return nil, xerror.WrapAuthenticateFailErrorWithMsg("Session Key Wrong")
		}

		if user.SessionCreateAt.AddDate(0, 0, stateful.DefaultConfig.RunTime.SessionExpireDay).Before(time.Now()) {
			return nil, xerror.WrapAuthenticateFailErrorWithMsg("Session Key Expired")
		}

		return user, nil
	},

	AuthTypeSkip: func(ctx context.Context, param *AuthenticateParam, manager *AuthenticateManager) (user *User, err error) {
		if !stateful.DefaultConfig.RunTime.SkipTokenValidate {
			return nil, xerror.WrapAuthenticateFailErrorWithMsg("Bad Authorization Flag")
		}

		return newFakeUser(&Role{
			Name: param.Identify,
		}), nil
	},
}

func (m *AuthenticateManager) Authenticate(ctx context.Context, param *AuthenticateParam) (user *User, err error) {
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

		if user.IsInner() {
			return m.storager.DeleteUser(ctx, user)
		}

		return m.storager.UpdateUser(ctx, user, &UserParam{
			SessionKey:      lib.PString(""),
			SessionCreateAt: lib.PTime(time.Time{}),
		})
	})
}

func (m *AuthenticateManager) CreateInnerUser(ctx context.Context) (user *User, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		var sessionKey string

		for i := 0; i < 5; i++ {
			sessionKey, err = sessionKeyFactory(15)
			if err != nil {
				continue
			}

			user, err = m.storager.FetchUser(ctx, &UserFilter{
				Name: &sessionKey,
			})
			if err != nil {
				continue
			}
			if user != nil {
				continue
			}

			user = &User{
				Name:       sessionKey,
				SessionKey: sessionKey,
			}

			return m.storager.CreateUser(ctx, &UserParam{
				Name:            &sessionKey,
				SessionKey:      &sessionKey,
				Roles:           []*Role{RoleInner},
				SessionCreateAt: lib.PTime(time.Now().AddDate(100, 0, 0)), // never expire
			})
		}

		return err
	})

	return
}

func (m *AuthenticateManager) CreateUser(ctx context.Context, param *UserParam) (err error) {
	if err = passwordCheck(*param.Password); err != nil {
		return err
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

		return m.storager.CreateUser(ctx, &UserParam{
			Name:     param.Name,
			Password: param.Password,
			Roles:    param.Roles,
		})
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
			return xerror.WrapModelErrorWithMsg("User Not Exist")
		}

		return m.storager.DeleteUser(ctx, user)
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

	return m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		user, err := m.storager.FetchUser(ctx, &UserFilter{
			Name: &pcd.UserName,
		})
		if err != nil {
			return err
		}

		if user == nil {
			return xerror.WrapModelErrorWithMsg("User Not Exist")
		}

		if pcd.OldPassword != "" {
			if user.Password != pcd.OldPassword {
				return xerror.WrapModelErrorWithMsg("Invalid Password")
			}
		}

		return m.storager.UpdateUser(ctx, user, &UserParam{
			Password:        &pcd.Password,
			SessionKey:      lib.PString(""),
			SessionCreateAt: lib.PTime(time.Time{}),
		})
	})
}

func (m *AuthenticateManager) FetchUser(ctx context.Context, userName string) (user *User, err error) {
	list, err := m.fetchUserList(ctx, &UserFilter{
		Name: &userName,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		user = list[0]
	}

	return
}

func (m *AuthenticateManager) fetchUserList(ctx context.Context, param *UserFilter) (users []*User, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		users, err = m.storager.FetchUserList(ctx, param)

		return err
	})

	return
}

func (m *AuthenticateManager) FetchInnerUserList(ctx context.Context) (users []*User, err error) {
	users, err = m.fetchUserList(ctx, nil)

	var tmpUserList []*User
	for _, one := range users {
		if !one.IsInner() {
			continue
		}

		tmpUserList = append(tmpUserList, one)
	}
	users = tmpUserList

	return
}

func (m *AuthenticateManager) FetchNormalUserList(ctx context.Context, param *UserParam) (users []*User, err error) {
	users, err = m.fetchUserList(ctx, nil)

	var tmpUserList []*User
	for _, one := range users {
		if one.IsInner() {
			continue
		}

		tmpUserList = append(tmpUserList, one)
	}
	users = tmpUserList

	return
}
