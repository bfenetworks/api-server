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

const tLbMatrixTableName = "lb_matrices"

// TLbMatrix Query Result
type TLbMatrix struct {
	ClusterID int64     `db:"cluster_id"`
	LbMatrix  string    `db:"lb_matrix"`
	ProductID int64     `db:"product_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// TLbMatrixOne Query One
// return (nil, nil) if record not existed
func TLbMatrixOne(dbCtx lib.DBContexter, where *TLbMatrixParam) (*TLbMatrix, error) {
	t := &TLbMatrix{}
	err := internal.QueryOne(dbCtx, tLbMatrixTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TLbMatrixList Query Multiple
func TLbMatrixList(dbCtx lib.DBContexter, where *TLbMatrixParam) ([]*TLbMatrix, error) {
	t := []*TLbMatrix{}
	err := internal.QueryList(dbCtx, tLbMatrixTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TLbMatrixParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TLbMatrixParam struct {
	ClusterIDs []int64 `db:"cluster_id,in"`

	ClusterID *int64     `db:"cluster_id"`
	LbMatrix  *string    `db:"lb_matrix"`
	ProductID *int64     `db:"product_id"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TLbMatrixCreate One/Multiple
func TLbMatrixCreate(dbCtx lib.DBContexter, data ...*TLbMatrixParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tLbMatrixTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tLbMatrixTableName, list...)
}

// TLbMatrixUpdate Update One
func TLbMatrixUpdate(dbCtx lib.DBContexter, val, where *TLbMatrixParam) (int64, error) {
	return internal.Update(dbCtx, tLbMatrixTableName, where, val)
}

// TLbMatrixDelete Delete One/Multiple
func TLbMatrixDelete(dbCtx lib.DBContexter, where *TLbMatrixParam) (int64, error) {
	return internal.Delete(dbCtx, tLbMatrixTableName, where)
}
