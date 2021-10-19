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

const tRouteAdvanceRuleTableName = "route_advance_rules"

// TRouteAdvanceRule Query Result
type TRouteAdvanceRule struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	ProductID   int64     `db:"product_id"`
	Expression  string    `db:"expression"`
	ClusterID   int64     `db:"cluster_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// TRouteAdvanceRuleOne Query One
// return (nil, nil) if record not existed
func TRouteAdvanceRuleOne(dbCtx lib.DBContexter, where *TRouteAdvanceRuleParam) (*TRouteAdvanceRule, error) {
	t := &TRouteAdvanceRule{}
	err := internal.QueryOne(dbCtx, tRouteAdvanceRuleTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TRouteAdvanceRuleList Query Multiple
func TRouteAdvanceRuleList(dbCtx lib.DBContexter, where *TRouteAdvanceRuleParam) ([]*TRouteAdvanceRule, error) {
	t := []*TRouteAdvanceRule{}
	err := internal.QueryList(dbCtx, tRouteAdvanceRuleTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TRouteAdvanceRuleParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TRouteAdvanceRuleParam struct {
	// IDs              []int64     `db:"id,in"`

	ID          *int64     `db:"id"`
	Name        *string    `db:"name"`
	Description *string    `db:"description"`
	ProductID   *int64     `db:"product_id"`
	ProductIDs  []int64    `db:"product_id,in"`
	ClusterID   *int64     `db:"cluster_id"`
	Expression  *string    `db:"expression"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`

	LockMode *string `db:"_lockMode"`
}

var (
	ModeForUpdate = "exclusive"
)

// TRouteAdvanceRuleCreate One/Multiple
func TRouteAdvanceRuleCreate(dbCtx lib.DBContexter, data ...*TRouteAdvanceRuleParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tRouteAdvanceRuleTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tRouteAdvanceRuleTableName, list...)
}

// TRouteAdvanceRuleUpdate Update One
func TRouteAdvanceRuleUpdate(dbCtx lib.DBContexter, val, where *TRouteAdvanceRuleParam) (int64, error) {
	return internal.Update(dbCtx, tRouteAdvanceRuleTableName, where, val)
}

// TRouteAdvanceRuleDelete Delete One/Multiple
func TRouteAdvanceRuleDelete(dbCtx lib.DBContexter, where *TRouteAdvanceRuleParam) (int64, error) {
	return internal.Delete(dbCtx, tRouteAdvanceRuleTableName, where)
}
