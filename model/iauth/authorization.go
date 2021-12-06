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

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/itxn"
)

const (
	ScopeAlwaysAllowed = "Allowed"
	ScopeSystem        = "System"
	ScopeProduct       = "Product"
	ScopeSupport       = "Support"
)

type AuthorizeStorager interface {
	UnbindUserProduct(ctx context.Context, user *User, product *ibasic.Product) error
	UnbindUserAllProduct(ctx context.Context, user *User) error
	BindUserProduct(ctx context.Context, user *User, product *ibasic.Product) error
	FetchUserProducts(ctx context.Context, user *User) ([]*ibasic.Product, error)
	FetchProductUsers(ctx context.Context, product *ibasic.Product) ([]*User, error)
	UpdateUserScopes(ctx context.Context, user *User, scopes []string) error
	IsUserProductGranted(ctx context.Context, user *User, product *ibasic.Product) (bool, error)

	UnbindTokenAllProduct(ctx context.Context, token *Token) error
	BindTokenProduct(ctx context.Context, token *Token, product *ibasic.Product) error
	FetchProductTokens(ctx context.Context, product *ibasic.Product) ([]*Token, error)
	IsTokenProductGranted(ctx context.Context, token *Token, product *ibasic.Product) (bool, error)
	FetchTokenProduct(ctx context.Context, token *Token) (*ibasic.Product, error)
	BatchFetchTokenProduct(ctx context.Context, token []*Token) (map[int64]*ibasic.Product, error)
}

type AuthorizeManager struct {
	storager AuthorizeStorager
	txn      itxn.TxnStorager
}

func NewAuthorizeManager(txn itxn.TxnStorager, storager AuthorizeStorager) *AuthorizeManager {
	return &AuthorizeManager{
		txn:      txn,
		storager: storager,
	}
}

func (m *AuthorizeManager) UpdateUserIsAdmin(ctx context.Context, user *User, isAdmin bool) (err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		mapping := map[bool][]string{
			false: {ScopeProduct},
			true:  {ScopeSystem},
		}

		return m.storager.UpdateUserScopes(ctx, user, mapping[isAdmin])
	})

	return
}

func (m *AuthorizeManager) Authorizate(ctx context.Context, authrizer *Authorization) (err error) {
	vistor, err := MustGetVistor(ctx)
	if err != nil {
		return
	}

	if vistor.IsAdmin() {
		return nil
	}

	featureGranted := false
	for _, scope := range vistor.GetScopes() {
		permissions, ok := scope2permission[scope]
		if !ok {
			continue
		}

		action, ok := permissions[authrizer.FeatureAuthorizer.Feature]
		if !ok {
			continue
		}

		featureGranted = action.IsAllowed(authrizer.FeatureAuthorizer.Action)
		if featureGranted {
			break
		}
	}

	if !featureGranted {
		return xerror.WrapAuthorizateFailErrorWithMsg("Feature Access Deny")
	}

	if authrizer.ValidateProduct {
		product, err := ibasic.MustGetProduct(ctx)
		if err != nil {
			return err
		}
		ok, err := m.IsVistorProductGranted(ctx, vistor, product)
		if err != nil {
			return err
		}

		if !ok {
			return xerror.WrapAuthorizateFailErrorWithMsg("Product Access Deny")
		}
	}

	return nil
}

func (m *AuthorizeManager) IsVistorProductGranted(ctx context.Context, v *Visitor, product *ibasic.Product) (bound bool, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		if user := v.User; user != nil {
			bound, err = m.storager.IsUserProductGranted(ctx, user, product)
			return err
		}

		bound, err = m.storager.IsTokenProductGranted(ctx, v.Token, product)
		return err
	})

	return
}

func (m *AuthorizeManager) FetchVistorProductList(ctx context.Context, v *Visitor) (userProducts []*ibasic.Product, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		if user := v.User; user != nil {
			userProducts, err = m.storager.FetchUserProducts(ctx, user)
			return err
		}

		tokenProduct, err := m.storager.FetchTokenProduct(ctx, v.Token)
		if err != nil {
			return err
		}

		userProducts = []*ibasic.Product{tokenProduct}
		return nil
	})

	return
}

func (m *AuthorizeManager) FetchProductUsers(ctx context.Context, product *ibasic.Product) (users []*User, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		users, err = m.storager.FetchProductUsers(ctx, product)

		return err
	})

	return
}

func (m *AuthorizeManager) BindUserProduct(ctx context.Context, user *User, product *ibasic.Product) (err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		err = m.storager.BindUserProduct(ctx, user, product)

		return err
	})

	return
}

func (m *AuthorizeManager) UnBindUserProduct(ctx context.Context, user *User, product *ibasic.Product) (err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		err = m.storager.UnbindUserProduct(ctx, user, product)

		return err
	})

	return
}

func (m *AuthorizeManager) FetchProductTokens(ctx context.Context, product *ibasic.Product) (tokens []*Token, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		tokens, err = m.storager.FetchProductTokens(ctx, product)

		return err
	})

	return
}
