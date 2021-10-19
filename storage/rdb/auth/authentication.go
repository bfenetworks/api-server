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

func (ps *RDBAuthenticateStorager) FetchUserList(ctx context.Context, param *iauth.UserFilter) ([]*iauth.User, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	list, err := dao.TUserList(dbCtx, userFilter2Param(param))
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
		Name: &user.Name,
	})

	return err
}

func userD2M(param *dao.TUser) *iauth.User {
	if param == nil {
		return nil
	}
	user := &iauth.User{
		ID:              param.ID,
		Name:            param.Name,
		Password:        param.Password,
		SessionKey:      param.SessionKey,
		SessionCreateAt: param.SessionKeyCreatedAt,
	}
	for _, one := range strings.Split(param.Roles, ",") {
		user.Roles = append(user.Roles, &iauth.Role{
			Name: one,
		})
	}

	return user
}

func userFilter2Param(filter *iauth.UserFilter) *dao.TUserParam {
	if filter == nil {
		return nil
	}

	return &dao.TUserParam{
		Name:       filter.Name,
		SessionKey: filter.SessionKey,
	}
}

func userParamM2D(param *iauth.UserParam) *dao.TUserParam {
	if param == nil {
		return nil
	}

	var (
		roles []string
		rs    *string
	)
	if param.Roles != nil {
		for _, one := range param.Roles {
			roles = append(roles, one.Name)
		}
		rs = lib.PString(strings.Join(roles, ","))
	}

	return &dao.TUserParam{
		Name:                param.Name,
		SessionKey:          param.SessionKey,
		SessionKeyCreatedAt: param.SessionCreateAt,
		Roles:               rs,
		Password:            param.Password,
	}
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
		Name: &user.Name,
	})

	return err
}
