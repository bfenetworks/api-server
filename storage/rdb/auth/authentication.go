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
	"strings"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao"
)

type RDBAuthenticateStorager struct {
	dbCtxFactory lib.DBContextFactory
}

var _ iauth.AuthenticateStorager = &RDBAuthenticateStorager{}

func NewAuthenticateStorager(dbCtxFactory lib.DBContextFactory) *RDBAuthenticateStorager {
	return &RDBAuthenticateStorager{
		dbCtxFactory: dbCtxFactory,
	}
}

func (ps *RDBAuthenticateStorager) FetchUserList(ctx context.Context, filter *iauth.UserFilter) ([]*iauth.User, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		filter = &iauth.UserFilter{}
	}
	filter.Type = &iauth.UserTypeNormal
	list, err := dao.TUserList(dbCtx, userFilter2Param(filter))
	if err != nil {
		return nil, err
	}

	rst := []*iauth.User{}
	for _, one := range list {
		rst = append(rst, userD2M(one))
	}

	return rst, nil
}

func (ps *RDBAuthenticateStorager) FetchUser(ctx context.Context, filter *iauth.UserFilter) (*iauth.User, error) {
	list, err := ps.FetchUserList(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, err
	}

	return list[0], nil
}

func (ps *RDBAuthenticateStorager) UpdateUser(ctx context.Context, user *iauth.User, param *iauth.UserParam) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TUserUpdate(dbCtx, userParamM2D(param), &dao.TUserParam{
		ID: &user.ID,
	})

	return err
}

func (ps *RDBAuthenticateStorager) CreateUser(ctx context.Context, param *iauth.UserParam) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TUserCreate(dbCtx, userParamM2D(param))

	return err
}

func (ps *RDBAuthenticateStorager) DeleteUser(ctx context.Context, user *iauth.User) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TUserDelete(dbCtx, &dao.TUserParam{
		ID: &user.ID,
	})

	return err
}

func userD2M(param *dao.TUser) *iauth.User {
	if param == nil {
		return nil
	}

	user := &iauth.User{
		ID:                 param.ID,
		Name:               param.Name,
		Type:               param.Type,
		Admin:              strings.Contains(param.Scopes, iauth.ScopeSystem),
		Password:           param.Password,
		SessionKey:         param.Ticket,
		SessionKeyCreateAt: param.TicketCreatedAt,
	}

	return user
}

func userFilter2Param(filter *iauth.UserFilter) *dao.TUserParam {
	if filter == nil {
		return nil
	}

	return &dao.TUserParam{
		IDs:    filter.IDs,
		Name:   filter.Name,
		Ticket: filter.SessionKey,
		Type:   filter.Type,
		Types:  filter.Types,
	}
}

func userParamM2D(param *iauth.UserParam) *dao.TUserParam {
	if param == nil {
		return nil
	}

	var scopes *string
	if param.Scopes != nil {
		scopes = lib.PString(strings.Join(param.Scopes, ","))
	}

	return &dao.TUserParam{
		Name:            param.Name,
		Type:            lib.PInt8(iauth.UserTypeNormal),
		Password:        param.Password,
		Scopes:          scopes,
		Ticket:          param.SessionKey,
		TicketCreatedAt: param.SessionKeyCreateAt,
	}
}

func (ps *RDBAuthenticateStorager) FetchTokens(ctx context.Context, filter *iauth.TokenFilter) ([]*iauth.Token, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	list, err := dao.TUserList(dbCtx, tokenFilter2Param(filter))
	if err != nil {
		return nil, err
	}

	rst := []*iauth.Token{}
	for _, one := range list {
		rst = append(rst, tokenD2M(one))
	}

	return rst, nil
}

func (ps *RDBAuthenticateStorager) CreateToken(ctx context.Context, param *iauth.TokenParam) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TUserCreate(dbCtx, tokenParamM2D(param))

	return err
}
func (ps *RDBAuthenticateStorager) DeleteToken(ctx context.Context, token *iauth.Token) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TUserDelete(dbCtx, &dao.TUserParam{
		Name: &token.Name,
		Type: lib.PInt8(iauth.UserTypeToken),
	})

	return err
}

func tokenD2M(param *dao.TUser) *iauth.Token {
	if param == nil {
		return nil
	}

	return &iauth.Token{
		ID:    param.ID,
		Name:  param.Name,
		Token: param.Ticket,
		Scope: param.Scopes,
	}
}

func tokenFilter2Param(filter *iauth.TokenFilter) *dao.TUserParam {
	if filter == nil {
		filter = &iauth.TokenFilter{}
	}

	return &dao.TUserParam{
		IDs:    filter.IDs,
		Name:   filter.Name,
		Type:   lib.PInt8(iauth.UserTypeToken),
		Ticket: filter.Token,
	}
}

func tokenParamM2D(param *iauth.TokenParam) *dao.TUserParam {
	if param == nil {
		return nil
	}

	return &dao.TUserParam{
		Name:   param.Name,
		Type:   lib.PInt8(iauth.UserTypeToken),
		Ticket: param.Token,
		Scopes: param.Scope,
	}
}
