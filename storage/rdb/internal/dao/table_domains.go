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

const tDomainTableName = "domains"

var (
	DomainTypeGeneral int32 = 7
)

// TDomain Query Result
type TDomain struct {
	ID                    int64     `db:"id"`
	Name                  string    `db:"name"`
	ProductID             int64     `db:"product_id"`
	Type                  int32     `db:"type"`
	UsingAdvancedRedirect int8      `db:"using_advanced_redirect"`
	UsingAdvancedHsts     int8      `db:"using_advanced_hsts"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`
}

// TDomainOne Query One
// return (nil, nil) if record not existed
func TDomainOne(dbCtx lib.DBContexter, where *TDomainParam) (*TDomain, error) {
	t := &TDomain{}
	err := internal.QueryOne(dbCtx, tDomainTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TDomainList Query Multiple
func TDomainList(dbCtx lib.DBContexter, where *TDomainParam) ([]*TDomain, error) {
	t := []*TDomain{}
	err := internal.QueryList(dbCtx, tDomainTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TDomainParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TDomainParam struct {
	// IDs              []int64     `db:"id,in"`

	ID                    *int64     `db:"id"`
	Name                  *string    `db:"name"`
	ProductID             *int64     `db:"product_id"`
	Type                  *int32     `db:"type"`
	UsingAdvancedRedirect *int8      `db:"using_advanced_redirect"`
	UsingAdvancedHsts     *int8      `db:"using_advanced_hsts"`
	CreatedAt             *time.Time `db:"created_at"`
	UpdatedAt             *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TDomainCreate One/Multiple
func TDomainCreate(dbCtx lib.DBContexter, data ...*TDomainParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tDomainTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tDomainTableName, list...)
}

// TDomainUpdate Update One
func TDomainUpdate(dbCtx lib.DBContexter, val, where *TDomainParam) (int64, error) {
	return internal.Update(dbCtx, tDomainTableName, where, val)
}

// TDomainDelete Delete One/Multiple
func TDomainDelete(dbCtx lib.DBContexter, where *TDomainParam) (int64, error) {
	return internal.Delete(dbCtx, tDomainTableName, where)
}
