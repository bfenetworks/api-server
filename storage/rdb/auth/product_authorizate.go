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

package auth

import (
	"context"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao"
)

type RDBProductAuthorizateStorager struct {
	dbCtxFactory lib.DBContextFactory
}

var _ iauth.ProductAuthorizateStorager = &RDBProductAuthorizateStorager{}

func NewProductAuthorizateStorager(dbCtxFactory lib.DBContextFactory) *RDBProductAuthorizateStorager {
	return &RDBProductAuthorizateStorager{
		dbCtxFactory: dbCtxFactory,
	}
}

func (ps *RDBProductAuthorizateStorager) FetchGrantedUsers(ctx context.Context, product *ibasic.Product) ([]*iauth.User, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	ups, err := dao.TUserProductList(dbCtx, &dao.TUserProductParam{
		ProductID: &product.ID,
	})
	if err != nil {
		return nil, err
	}
	if len(ups) == 0 {
		return nil, nil
	}

	userIDs := []int64{}
	for _, one := range ups {
		userIDs = append(userIDs, one.UserID)
	}

	us, err := dao.TUserList(dbCtx, &dao.TUserParam{
		IDs: userIDs,
	})
	if err != nil {
		return nil, err
	}

	users := []*iauth.User{}
	for _, one := range us {
		users = append(users, userD2M(one))
	}

	return users, nil
}

func (ps *RDBProductAuthorizateStorager) FetchUser(ctx context.Context, product *ibasic.Product, user *iauth.User) (*iauth.User, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	up, err := dao.TUserProductOne(dbCtx, &dao.TUserProductParam{
		ProductID: &product.ID,
		UserID:    &user.ID,
	})
	if err != nil {
		return nil, err
	}
	if up == nil {
		return nil, nil
	}

	return user, nil
}

func (ps *RDBProductAuthorizateStorager) Grant(ctx context.Context, product *ibasic.Product, user *iauth.User) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TUserProductCreate(dbCtx, &dao.TUserProductParam{
		ProductID: &product.ID,
		UserID:    &user.ID,
	})

	return err
}

func (ps *RDBProductAuthorizateStorager) Revoke(ctx context.Context, product *ibasic.Product, user *iauth.User) error {
	u, err := ps.FetchUser(ctx, product, user)
	if err != nil {
		return err
	}
	if u == nil {
		return nil
	}

	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TUserProductDelete(dbCtx, &dao.TUserProductParam{
		ProductID: &product.ID,
		UserID:    &user.ID,
	})

	return err
}

func (ps *RDBProductAuthorizateStorager) FetchProducts(ctx context.Context, user *iauth.User) ([]int64, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	ups, err := dao.TUserProductList(dbCtx, &dao.TUserProductParam{
		UserID: &user.ID,
	})
	if err != nil {
		return nil, err
	}

	if len(ups) == 0 {
		return nil, nil
	}

	ids := []int64{}
	for _, one := range ups {
		ids = append(ids, one.ProductID)
	}

	return ids, nil
}
