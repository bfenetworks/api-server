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

const tPoolsTableName = "pools"

// TPools Query Result
type TPools struct {
	Id             int64     `db:"id"`
	Name           string    `db:"name"`
	Ready          bool      `db:"ready"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	ProductId      int64     `db:"product_id"`
	Type           int8      `db:"type"`
	InstanceDetail string    `db:"instance_detail"`
	Tag            int8      `db:"tag"`
}

// TPoolsOne Query One
func TPoolsOne(dbCtx lib.DBContexter, where *TPoolsParam) (*TPools, error) {
	t := &TPools{}
	err := internal.QueryOne(dbCtx, tPoolsTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TPoolsList Query Multiple
func TPoolsList(dbCtx lib.DBContexter, where *TPoolsParam) ([]*TPools, error) {
	t := []*TPools{}
	err := internal.QueryList(dbCtx, tPoolsTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TPoolsParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TPoolsParam struct {
	Ids []int64 `db:"id,in"`

	Id   *int64  `db:"id"`
	Name *string `db:"name"`

	Ready          *bool      `db:"ready"`
	ProductID      *int64     `db:"product_id"`
	Type           *int8      `db:"type"`
	InstanceDetail *string    `db:"instance_detail"`
	Tag            *int8      `db:"tag"`
	CreatedAt      *time.Time `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TPoolsCreate One/Multiple
func TPoolsCreate(dbCtx lib.DBContexter, data ...*TPoolsParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tPoolsTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tPoolsTableName, list...)
}

// TPoolsUpdate Update One
func TPoolsUpdate(dbCtx lib.DBContexter, val, where *TPoolsParam) (int64, error) {
	return internal.Update(dbCtx, tPoolsTableName, where, val)
}

// TPoolsDelete Delete One/Multiple
func TPoolsDelete(dbCtx lib.DBContexter, where *TPoolsParam) (int64, error) {
	return internal.Delete(dbCtx, tPoolsTableName, where)
}

func PoolsList2Map(list []*TPools) map[int64]*TPools {
	m := map[int64]*TPools{}
	for _, one := range list {
		m[one.Id] = one
	}

	return m
}
