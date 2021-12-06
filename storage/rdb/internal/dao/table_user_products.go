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

const tUserProductTableName = "user_products"

// TUserProduct Query Result
type TUserProduct struct {
	UserID    int64     `db:"user_id"`
	ProductID int64     `db:"product_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// TUserProductOne Query One
// return (nil, nil) if record not existed
func TUserProductOne(dbCtx lib.DBContexter, where *TUserProductParam) (*TUserProduct, error) {
	t := &TUserProduct{}
	err := internal.QueryOne(dbCtx, tUserProductTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TUserProductList Query Multiple
func TUserProductList(dbCtx lib.DBContexter, where *TUserProductParam) ([]*TUserProduct, error) {
	t := []*TUserProduct{}
	err := internal.QueryList(dbCtx, tUserProductTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TUserProductParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TUserProductParam struct {
	// IDs              []int64     `db:"id,in"`

	UserID    *int64     `db:"user_id"`
	UserIDs   []int64    `db:"user_id,in"`
	ProductID *int64     `db:"product_id"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TUserProductCreate One/Multiple
func TUserProductCreate(dbCtx lib.DBContexter, data ...*TUserProductParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tUserProductTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tUserProductTableName, list...)
}

// TUserProductUpdate Update One
func TUserProductUpdate(dbCtx lib.DBContexter, val, where *TUserProductParam) (int64, error) {
	return internal.Update(dbCtx, tUserProductTableName, where, val)
}

// TUserProductDelete Delete One/Multiple
func TUserProductDelete(dbCtx lib.DBContexter, where *TUserProductParam) (int64, error) {
	return internal.Delete(dbCtx, tUserProductTableName, where)
}
