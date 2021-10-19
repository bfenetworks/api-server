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

const tBfeClusterTableName = "bfe_clusters"

// TBfeCluster Query Result
type TBfeCluster struct {
	ID                 int64     `db:"id"`
	Name               string    `db:"name"`
	PoolName           string    `db:"pool_name"`
	Capacity           int64     `db:"capacity"`
	Enabled            bool      `db:"enabled"`
	GtcEnabled         bool      `db:"gtc_enabled"`
	GtcManualEnabled   bool      `db:"gtc_manual_enabled"`
	ExemptTrafficCheck bool      `db:"exempt_traffic_check"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// TBfeClusterOne Query One
// return (nil, nil) if record not existed
func TBfeClusterOne(dbCtx lib.DBContexter, where *TBfeClusterParam) (*TBfeCluster, error) {
	t := &TBfeCluster{}
	err := internal.QueryOne(dbCtx, tBfeClusterTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TBfeClusterList Query Multiple
func TBfeClusterList(dbCtx lib.DBContexter, where *TBfeClusterParam) ([]*TBfeCluster, error) {
	t := []*TBfeCluster{}
	err := internal.QueryList(dbCtx, tBfeClusterTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TBfeClusterParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TBfeClusterParam struct {
	// IDs              []int64     `db:"id,in"`

	ID                 *int64     `db:"id"`
	Name               *string    `db:"name"`
	PoolName           *string    `db:"pool_name"`
	Capacity           *int64     `db:"capacity"`
	Enabled            *bool      `db:"enabled"`
	GtcEnabled         *bool      `db:"gtc_enabled"`
	GtcManualEnabled   *bool      `db:"gtc_manual_enabled"`
	ExemptTrafficCheck *bool      `db:"exempt_traffic_check"`
	CreatedAt          *time.Time `db:"created_at"`
	UpdatedAt          *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TBfeClusterCreate One/Multiple
func TBfeClusterCreate(dbCtx lib.DBContexter, data ...*TBfeClusterParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tBfeClusterTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tBfeClusterTableName, list...)
}

// TBfeClusterUpdate Update One
func TBfeClusterUpdate(dbCtx lib.DBContexter, val, where *TBfeClusterParam) (int64, error) {
	return internal.Update(dbCtx, tBfeClusterTableName, where, val)
}

// TBfeClusterDelete Delete One/Multiple
func TBfeClusterDelete(dbCtx lib.DBContexter, where *TBfeClusterParam) (int64, error) {
	return internal.Delete(dbCtx, tBfeClusterTableName, where)
}
