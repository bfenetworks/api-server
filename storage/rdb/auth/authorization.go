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

type RDBAuthorizeStorager struct {
	dbCtxFactory lib.DBContextFactory

	productStorager      ibasic.ProductStorager
	authenticateStorager iauth.AuthenticateStorager
}

var _ iauth.AuthorizeStorager = &RDBAuthorizeStorager{}

func NewAuthorizeStorager(dbCtxFactory lib.DBContextFactory,
	productStorager ibasic.ProductStorager,
	authenticateStorager iauth.AuthenticateStorager) *RDBAuthorizeStorager {
	return &RDBAuthorizeStorager{
		dbCtxFactory:         dbCtxFactory,
		productStorager:      productStorager,
		authenticateStorager: authenticateStorager,
	}
}

func (ps *RDBAuthorizeStorager) UnbindUserProduct(ctx context.Context, user *iauth.User, product *ibasic.Product) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TUserProductDelete(dbCtx, &dao.TUserProductParam{
		UserID:    &user.ID,
		ProductID: &product.ID,
	})

	return err
}

func (ps *RDBAuthorizeStorager) UnbindUserAllProduct(ctx context.Context, user *iauth.User) error {
	return ps.unbindTokenAllProduct(ctx, user.ID)
}

func (ps *RDBAuthorizeStorager) BindUserProduct(ctx context.Context, user *iauth.User, product *ibasic.Product) error {
	return ps.bindUserProduct(ctx, user.ID, product)
}

func (ps *RDBAuthorizeStorager) FetchUserProducts(ctx context.Context, user *iauth.User) ([]*ibasic.Product, error) {
	return ps.fetchUserProducts(ctx, user.ID)
}

func (ps *RDBAuthorizeStorager) FetchProductUsers(ctx context.Context, product *ibasic.Product) ([]*iauth.User, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	list, err := dao.TUserProductList(dbCtx, &dao.TUserProductParam{
		ProductID: &product.ID,
	})
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, err
	}

	userMap := map[int64]bool{}
	for _, one := range list {
		userMap[one.UserID] = true
	}

	return ps.authenticateStorager.FetchUserList(dbCtx, &iauth.UserFilter{
		IDs: lib.Int64BoolMap2Slice(userMap),
	})
}

func (ps *RDBAuthorizeStorager) IsUserProductGranted(ctx context.Context, user *iauth.User, product *ibasic.Product) (bool, error) {
	return ps.isTokenProductGranted(ctx, user.ID, product)
}

func (ps *RDBAuthorizeStorager) FetchTokenProduct(ctx context.Context, token *iauth.Token) (*ibasic.Product, error) {
	list, err := ps.fetchUserProducts(ctx, token.ID)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return list[0], nil
}

func (ps *RDBAuthorizeStorager) BatchFetchTokenProduct(ctx context.Context, tokens []*iauth.Token) (map[int64]*ibasic.Product, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	tokenIDs := []int64{}
	for _, one := range tokens {
		tokenIDs = append(tokenIDs, one.ID)
	}

	bindList, err := dao.TUserProductList(dbCtx, &dao.TUserProductParam{
		UserIDs: tokenIDs,
	})
	if err != nil {
		return nil, err
	}

	if len(bindList) == 0 {
		return nil, err
	}

	productMap := map[int64]bool{}
	for _, one := range bindList {
		productMap[one.ProductID] = true
	}

	productList, err := ps.productStorager.FetchProducts(dbCtx, &ibasic.ProductFilter{
		IDs: lib.Int64BoolMap2Slice(productMap),
	})
	if err != nil {
		return nil, err
	}

	productObjMap := ibasic.ProductIDMap(productList)
	rst := map[int64]*ibasic.Product{}
	for _, one := range bindList {
		rst[one.UserID] = productObjMap[one.ProductID]
	}

	return rst, nil
}

func (ps *RDBAuthorizeStorager) fetchUserProducts(ctx context.Context, userID int64) ([]*ibasic.Product, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	list, err := dao.TUserProductList(dbCtx, &dao.TUserProductParam{
		UserID: &userID,
	})
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, err
	}

	productMap := map[int64]bool{}
	for _, one := range list {
		productMap[one.ProductID] = true
	}

	return ps.productStorager.FetchProducts(dbCtx, &ibasic.ProductFilter{
		IDs: lib.Int64BoolMap2Slice(productMap),
	})
}

func (ps *RDBAuthorizeStorager) FetchProductTokens(ctx context.Context, product *ibasic.Product) ([]*iauth.Token, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	list, err := dao.TUserProductList(dbCtx, &dao.TUserProductParam{
		ProductID: &product.ID,
	})
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, err
	}

	tokeMap := map[int64]bool{}
	for _, one := range list {
		tokeMap[one.UserID] = true
	}

	return ps.authenticateStorager.FetchTokens(dbCtx, &iauth.TokenFilter{
		IDs: lib.Int64BoolMap2Slice(tokeMap),
	})
}

func (ps *RDBAuthorizeStorager) UpdateUserScopes(ctx context.Context, user *iauth.User, scopes []string) error {
	return ps.authenticateStorager.UpdateUser(ctx, user, &iauth.UserParam{
		Scopes: scopes,
	})
}

func (ps *RDBAuthorizeStorager) IsTokenProductGranted(ctx context.Context, token *iauth.Token, product *ibasic.Product) (bool, error) {
	return ps.isTokenProductGranted(ctx, token.ID, product)
}

func (ps *RDBAuthorizeStorager) isTokenProductGranted(ctx context.Context, userID int64, product *ibasic.Product) (bool, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return false, err
	}

	one, err := dao.TUserProductOne(dbCtx, &dao.TUserProductParam{
		UserID:    &userID,
		ProductID: &product.ID,
	})

	return one != nil, err
}

func (ps *RDBAuthorizeStorager) BindTokenProduct(ctx context.Context, token *iauth.Token, product *ibasic.Product) error {
	return ps.bindUserProduct(ctx, token.ID, product)
}

func (ps *RDBAuthorizeStorager) bindUserProduct(ctx context.Context, userID int64, product *ibasic.Product) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	one, err := dao.TUserProductOne(dbCtx, &dao.TUserProductParam{
		UserID:    &userID,
		ProductID: &product.ID,
	})
	if err != nil {
		return err
	}
	if one != nil {
		return nil
	}

	_, err = dao.TUserProductCreate(dbCtx, &dao.TUserProductParam{
		UserID:    &userID,
		ProductID: &product.ID,
	})

	return err
}

func (ps *RDBAuthorizeStorager) UnbindTokenAllProduct(ctx context.Context, token *iauth.Token) error {
	return ps.unbindTokenAllProduct(ctx, token.ID)
}

func (ps *RDBAuthorizeStorager) unbindTokenAllProduct(ctx context.Context, userID int64) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TUserProductDelete(dbCtx, &dao.TUserProductParam{
		UserID: &userID,
	})

	return err
}
