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

	"github.com/bfenetworks/api-server/model/ibasic"
)

type ProductAuthorizateStorager interface {
	FetchGrantedUsers(context.Context, *ibasic.Product) ([]*User, error)
	FetchUser(context.Context, *ibasic.Product, *User) (*User, error)
	FetchProducts(context.Context, *User) ([]int64, error)
	Grant(context.Context, *ibasic.Product, *User) error
	Revoke(context.Context, *ibasic.Product, *User) error
}

func NewProductAuthorizateManager(storager ProductAuthorizateStorager) *ProductAuthorizateManager {
	return &ProductAuthorizateManager{
		storager: storager,
	}
}

type ProductAuthorizateManager struct {
	storager ProductAuthorizateStorager
}

func (pa *ProductAuthorizateManager) Authorizate(ctx context.Context, user *User) (bool, error) {
	if user.IsAdmin() {
		return true, nil
	}

	product, err := ibasic.MustGetProduct(ctx)
	if err != nil {
		return false, err
	}

	user, err = pa.storager.FetchUser(ctx, product, user)
	return user != nil, err
}

func (pa *ProductAuthorizateManager) Grant(ctx context.Context, user *User) error {
	product, err := ibasic.MustGetProduct(ctx)
	if err != nil {
		return err
	}

	return pa.storager.Grant(ctx, product, user)
}

func (pa *ProductAuthorizateManager) Revoke(ctx context.Context, user *User) error {
	product, err := ibasic.MustGetProduct(ctx)
	if err != nil {
		return err
	}
	return pa.storager.Revoke(ctx, product, user)
}

func (pa *ProductAuthorizateManager) GrantedUsers(ctx context.Context) ([]*User, error) {
	product, err := ibasic.MustGetProduct(ctx)
	if err != nil {
		return nil, err
	}

	return pa.storager.FetchGrantedUsers(ctx, product)
}

func (pa *ProductAuthorizateManager) FetchProducts(ctx context.Context, user *User) ([]int64, error) {
	return pa.storager.FetchProducts(ctx, user)
}
