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

const tSubClusterTableName = "sub_clusters"

// TSubCluster Query Result
type TSubCluster struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	ClusterID   int64     `db:"cluster_id"`
	ProductID   int64     `db:"product_id"`
	Description string    `db:"description"`
	PoolsID     int64     `db:"bns_name_id"`
	Enabled     bool      `db:"enabled"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// TSubClusterOne Query One
func TSubClusterOne(dbCtx lib.DBContexter, where *TSubClusterParam) (*TSubCluster, error) {
	t := &TSubCluster{}
	err := internal.QueryOne(dbCtx, tSubClusterTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TSubClusterList Query Multiple
func TSubClusterList(dbCtx lib.DBContexter, where *TSubClusterParam) ([]*TSubCluster, error) {
	t := []*TSubCluster{}
	err := internal.QueryList(dbCtx, tSubClusterTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TSubClusterParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TSubClusterParam struct {
	IDs []int64 `db:"id,in"`

	ID          *int64     `db:"id"`
	Name        *string    `db:"name"`
	Names       []string   `db:"name,in"`
	ClusterID   *int64     `db:"cluster_id"`
	ClusterIDs  []int64    `db:"cluster_id,in"`
	ProductID   *int64     `db:"product_id"`
	Description *string    `db:"description"`
	PoolsID     *int64     `db:"bns_name_id"`
	PoolsIDs    []int64    `db:"bns_name_id,in"`
	Enabled     *bool      `db:"enabled"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TSubClusterCreate One/Multiple
func TSubClusterCreate(dbCtx lib.DBContexter, data ...*TSubClusterParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tSubClusterTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tSubClusterTableName, list...)
}

// TSubClusterUpdate Update One
func TSubClusterUpdate(dbCtx lib.DBContexter, val, where *TSubClusterParam) (int64, error) {
	return internal.Update(dbCtx, tSubClusterTableName, where, val)
}

// TSubClusterDelete Delete One/Multiple
func TSubClusterDelete(dbCtx lib.DBContexter, where *TSubClusterParam) (int64, error) {
	return internal.Delete(dbCtx, tSubClusterTableName, where)
}
