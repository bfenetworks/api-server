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

package dao

import (
	"time"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao/internal"
)

const tUserTableName = "users"

// TUser Query Result
type TUser struct {
	ID                  int64     `db:"id"`
	Name                string    `db:"name"`
	SessionKey          string    `db:"session_key"`
	SessionKeyCreatedAt time.Time `db:"session_key_created_at"`
	Roles               string    `db:"roles"`
	Password            string    `db:"password"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
}

// TUserOne Query One
// return (nil, nil) if record not existed
func TUserOne(dbCtx lib.DBContexter, where *TUserParam) (*TUser, error) {
	t := &TUser{}
	err := internal.QueryOne(dbCtx, tUserTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TUserList Query Multiple
func TUserList(dbCtx lib.DBContexter, where *TUserParam) ([]*TUser, error) {
	t := []*TUser{}
	err := internal.QueryList(dbCtx, tUserTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TUserParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TUserParam struct {
	IDs []int64 `db:"id,in"`

	ID                  *int64     `db:"id"`
	Name                *string    `db:"name"`
	SessionKey          *string    `db:"session_key"`
	SessionKeyCreatedAt *time.Time `db:"session_key_created_at"`
	Roles               *string    `db:"roles"`
	Password            *string    `db:"password"`
	CreatedAt           *time.Time `db:"created_at"`
	UpdatedAt           *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TUserCreate One/Multiple
func TUserCreate(dbCtx lib.DBContexter, data ...*TUserParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tUserTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tUserTableName, list...)
}

// TUserUpdate Update One
func TUserUpdate(dbCtx lib.DBContexter, val, where *TUserParam) (int64, error) {
	return internal.Update(dbCtx, tUserTableName, where, val)
}

// TUserDelete Delete One/Multiple
func TUserDelete(dbCtx lib.DBContexter, where *TUserParam) (int64, error) {
	return internal.Delete(dbCtx, tUserTableName, where)
}
