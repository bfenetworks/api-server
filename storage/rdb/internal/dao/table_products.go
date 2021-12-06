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
	"fmt"
	"strings"
	"time"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao/internal"
)

const tProductTableName = "products"

// TProduct Query Result
type TProduct struct {
	Id            int64     `db:"id"`
	Name          string    `db:"name"`
	MailList      string    `db:"mail_list"`
	ContactPerson string    `db:"contact_person"`
	SmsList       string    `db:"sms_list"`
	Description   string    `db:"description"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// TProductOne Query One
func TProductOne(dbCtx lib.DBContexter, where *TProductParam) (*TProduct, error) {
	t := &TProduct{}
	err := internal.QueryOne(dbCtx, tProductTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TProductList Query Multiple
func TProductList(dbCtx lib.DBContexter, where *TProductParam) ([]*TProduct, error) {
	t := []*TProduct{}
	err := internal.QueryList(dbCtx, tProductTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TProductParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TProductParam struct {
	// Ids              []int64     `db:"id,in"`

	Id            *int64     `db:"id"`
	Name          *string    `db:"name"`
	MailList      *string    `db:"mail_list"`
	ContactPerson *string    `db:"contact_person"`
	SmsList       *string    `db:"sms_list"`
	Description   *string    `db:"description"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`

	IDs  []int64 `db:"id,in"`
	NeId *int64  `db:"id,!="`

	OrderBy *string `db:"_orderby"`
}

// TProductCreate One/Multiple
func TProductCreate(dbCtx lib.DBContexter, data ...*TProductParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tProductTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tProductTableName, list...)
}

// TProductUpdate Update One
func TProductUpdate(dbCtx lib.DBContexter, val, where *TProductParam) (int64, error) {
	return internal.Update(dbCtx, tProductTableName, where, val)
}

// TProductDelete Delete One/Multiple
func TProductDelete(dbCtx lib.DBContexter, where *TProductParam) (int64, error) {
	return internal.Delete(dbCtx, tProductTableName, where)
}

var deleteSQL = `
DELETE FROM products  				WHERE id = xxx;
DELETE FROM domains  				WHERE product_id = xxx;
DELETE FROM clusters  				WHERE product_id = xxx;
DELETE FROM lb_matrices  			WHERE product_id = xxx;
DELETE FROM sub_clusters  			WHERE product_id = xxx;
DELETE FROM pools  					WHERE product_id = xxx;
DELETE FROM route_basic_rules  		WHERE product_id = xxx;
DELETE FROM route_advance_rules 	WHERE product_id = xxx;
DELETE FROM user_products  			WHERE product_id = xxx;
DELETE FROM extra_files  			WHERE product_id = xxx; `

func TProductDeleteByProductID(dbCtx lib.DBContexter, productID int64) error {
	sql := strings.Replace(deleteSQL, "xxx", fmt.Sprintf("%d", productID), -1)

	_, err := dbCtx.Conn().ExecContext(dbCtx, sql)
	return err
}
