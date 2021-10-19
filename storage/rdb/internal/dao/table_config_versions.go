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

const tConfigVersionTableName = "config_versions"

// TConfigVersion Query Result
type TConfigVersion struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	DataSign  string    `db:"data_sign"`
	Version   string    `db:"version"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// TConfigVersionOne Query One
// return (nil, nil) if record not existed
func TConfigVersionOne(dbCtx lib.DBContexter, where *TConfigVersionParam) (*TConfigVersion, error) {
	t := &TConfigVersion{}
	err := internal.QueryOne(dbCtx, tConfigVersionTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TConfigVersionList Query Multiple
func TConfigVersionList(dbCtx lib.DBContexter, where *TConfigVersionParam) ([]*TConfigVersion, error) {
	t := []*TConfigVersion{}
	err := internal.QueryList(dbCtx, tConfigVersionTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TConfigVersionParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TConfigVersionParam struct {
	// IDs              []int64     `db:"id,in"`

	ID        *int64     `db:"id"`
	Name      *string    `db:"name"`
	DataSign  *string    `db:"data_sign"`
	Version   *string    `db:"version"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TConfigVersionCreate One/Multiple
func TConfigVersionCreate(dbCtx lib.DBContexter, data ...*TConfigVersionParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tConfigVersionTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tConfigVersionTableName, list...)
}

// TConfigVersionUpdate Update One
func TConfigVersionUpdate(dbCtx lib.DBContexter, val, where *TConfigVersionParam) (int64, error) {
	return internal.Update(dbCtx, tConfigVersionTableName, where, val)
}

// TConfigVersionDelete Delete One/Multiple
func TConfigVersionDelete(dbCtx lib.DBContexter, where *TConfigVersionParam) (int64, error) {
	return internal.Delete(dbCtx, tConfigVersionTableName, where)
}
