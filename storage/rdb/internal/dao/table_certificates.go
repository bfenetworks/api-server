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

const tCertificateTableName = "certificates"

// TCertificate Query Result
type TCertificate struct {
	ID           int64     `db:"id"`
	CertName     string    `db:"cert_name"`
	Description  string    `db:"description"`
	IsDefault    bool      `db:"is_default"`
	ExpiredDate  string    `db:"expired_date"`
	CertFileName string    `db:"cert_file_name"`
	CertFilePath string    `db:"cert_file_path"`
	KeyFileName  string    `db:"key_file_name"`
	KeyFilePath  string    `db:"key_file_path"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// TCertificateOne Query One
// return (nil, nil) if record not existed
func TCertificateOne(dbCtx lib.DBContexter, where *TCertificateParam) (*TCertificate, error) {
	t := &TCertificate{}
	err := internal.QueryOne(dbCtx, tCertificateTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TCertificateList Query Multiple
func TCertificateList(dbCtx lib.DBContexter, where *TCertificateParam) ([]*TCertificate, error) {
	t := []*TCertificate{}
	err := internal.QueryList(dbCtx, tCertificateTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TCertificateParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TCertificateParam struct {
	// IDs              []int64     `db:"id,in"`

	ID           *int64     `db:"id"`
	CertName     *string    `db:"cert_name"`
	Description  *string    `db:"description"`
	IsDefault    *bool      `db:"is_default"`
	ExpiredDate  *string    `db:"expired_date"`
	CertFileName *string    `db:"cert_file_name"`
	CertFilePath *string    `db:"cert_file_path"`
	KeyFileName  *string    `db:"key_file_name"`
	KeyFilePath  *string    `db:"key_file_path"`
	CreatedAt    *time.Time `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TCertificateCreate One/Multiple
func TCertificateCreate(dbCtx lib.DBContexter, data ...*TCertificateParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tCertificateTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tCertificateTableName, list...)
}

// TCertificateUpdate Update One
func TCertificateUpdate(dbCtx lib.DBContexter, val, where *TCertificateParam) (int64, error) {
	return internal.Update(dbCtx, tCertificateTableName, where, val)
}

// TCertificateDelete Delete One/Multiple
func TCertificateDelete(dbCtx lib.DBContexter, where *TCertificateParam) (int64, error) {
	return internal.Delete(dbCtx, tCertificateTableName, where)
}
