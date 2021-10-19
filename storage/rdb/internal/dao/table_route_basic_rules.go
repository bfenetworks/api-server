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

const tRouteBasicRuleTableName = "route_basic_rules"

// TRouteBasicRule Query Result
type TRouteBasicRule struct {
	ID          int64     `db:"id"`
	Description string    `db:"description"`
	ProductID   int64     `db:"product_id"`
	HostNames   []byte    `db:"host_names"`
	Paths       []byte    `db:"paths"`
	ClusterID   int64     `db:"cluster_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// TRouteBasicRuleOne Query One
// return (nil, nil) if record not existed
func TRouteBasicRuleOne(dbCtx lib.DBContexter, where *TRouteBasicRuleParam) (*TRouteBasicRule, error) {
	t := &TRouteBasicRule{}
	err := internal.QueryOne(dbCtx, tRouteBasicRuleTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TRouteBasicRuleList Query Multiple
func TRouteBasicRuleList(dbCtx lib.DBContexter, where *TRouteBasicRuleParam) ([]*TRouteBasicRule, error) {
	t := []*TRouteBasicRule{}
	err := internal.QueryList(dbCtx, tRouteBasicRuleTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TRouteBasicRuleParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TRouteBasicRuleParam struct {
	// IDs              []int64     `db:"id,in"`

	ID          *int64     `db:"id"`
	Description *string    `db:"description"`
	ProductID   *int64     `db:"product_id"`
	ProductIDs  []int64    `db:"product_id,in"`
	HostNames   []byte     `db:"host_names"`
	Paths       []byte     `db:"paths"`
	ClusterID   *int64     `db:"cluster_id"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TRouteBasicRuleCreate One/Multiple
func TRouteBasicRuleCreate(dbCtx lib.DBContexter, data ...*TRouteBasicRuleParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tRouteBasicRuleTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tRouteBasicRuleTableName, list...)
}

// TRouteBasicRuleUpdate Update One
func TRouteBasicRuleUpdate(dbCtx lib.DBContexter, val, where *TRouteBasicRuleParam) (int64, error) {
	return internal.Update(dbCtx, tRouteBasicRuleTableName, where, val)
}

// TRouteBasicRuleDelete Delete One/Multiple
func TRouteBasicRuleDelete(dbCtx lib.DBContexter, where *TRouteBasicRuleParam) (int64, error) {
	return internal.Delete(dbCtx, tRouteBasicRuleTableName, where)
}
