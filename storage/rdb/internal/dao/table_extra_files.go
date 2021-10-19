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

const tExtraFileTableName = "extra_files"

// TExtraFile Query Result
type TExtraFile struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	ProductID   int64     `db:"product_id"`
	Description string    `db:"description"`
	Md5         []byte    `db:"md5"`
	Content     []byte    `db:"content"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// TExtraFileOne Query One
// return (nil, nil) if record not existed
func TExtraFileOne(dbCtx lib.DBContexter, where *TExtraFileParam) (*TExtraFile, error) {
	t := &TExtraFile{}
	err := internal.QueryOne(dbCtx, tExtraFileTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TExtraFileList Query Multiple
func TExtraFileList(dbCtx lib.DBContexter, where *TExtraFileParam) ([]*TExtraFile, error) {
	t := []*TExtraFile{}
	err := internal.QueryList(dbCtx, tExtraFileTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TExtraFileParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TExtraFileParam struct {
	// IDs              []int64     `db:"id,in"`

	ID          *int64     `db:"id"`
	Name        *string    `db:"name"`
	Names       []string   `db:"name,in"`
	ProductID   *int64     `db:"product_id"`
	Description *string    `db:"description"`
	Md5         []byte     `db:"md5"`
	Content     []byte     `db:"content"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TExtraFileCreate One/Multiple
func TExtraFileCreate(dbCtx lib.DBContexter, data ...*TExtraFileParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tExtraFileTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tExtraFileTableName, list...)
}

// TExtraFileUpdate Update One
func TExtraFileUpdate(dbCtx lib.DBContexter, val, where *TExtraFileParam) (int64, error) {
	return internal.Update(dbCtx, tExtraFileTableName, where, val)
}

// TExtraFileDelete Delete One/Multiple
func TExtraFileDelete(dbCtx lib.DBContexter, where *TExtraFileParam) (int64, error) {
	return internal.Delete(dbCtx, tExtraFileTableName, where)
}
