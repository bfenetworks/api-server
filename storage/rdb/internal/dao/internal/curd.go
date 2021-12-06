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

package internal

import (
	"time"

	"github.com/didi/gendry/scanner"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/stateful"
)

func init() {
	scanner.SetTagName(tagName)
}

func QueryOne(dbCtx lib.DBContexter, table string, where interface{}, rst interface{}) error {
	return queryList(dbCtx, table, where, rst, true)
}

func QueryList(dbCtx lib.DBContexter, table string, where interface{}, rst interface{}) error {
	return queryList(dbCtx, table, where, rst, false)
}

func queryList(dbCtx lib.DBContexter, table string, where interface{}, rst interface{}, queryOne bool) error {
	tmp := Struct2Where(where)
	if tmp != nil && queryOne {
		tmp["_limit"] = []uint{0, 1}
	}

	build := NewSelectBuilder(table, tmp, nil)
	sql, args, err := build.Compile()
	if err != nil {
		return xerror.WrapDaoError(err)
	}

	now := time.Now()
	rows, err := dbCtx.Conn().QueryContext(dbCtx, sql, args...) // .ScanContext(dbCtx, rst).Error
	record := &stateful.SQLRecord{
		SQL:  sql,
		Args: args,
		Err:  err,
		Cost: time.Since(now),
	}
	defer record.Print(dbCtx)

	if err != nil {
		return xerror.WrapDaoError(err)
	}

	record.Err = scanner.Scan(rows, rst)
	if record.Err == scanner.ErrEmptyResult {
		return ErrRecordNotFound
	}
	return xerror.WrapDaoError(record.Err)
}

func Create(dbCtx lib.DBContexter, table string, data ...interface{}) (int64, error) {
	build := NewInsertBuilder(table, Struct2AssignList(data...))
	sql, args, err := build.Compile()
	if err != nil {
		return 0, xerror.WrapDaoError(err)
	}

	now := time.Now()

	rst, err := dbCtx.Conn().ExecContext(dbCtx, sql, args...)
	sr := &stateful.SQLRecord{
		SQL:  sql,
		Args: args,
		Err:  err,
		Cost: time.Since(now),
	}
	sr.Print(dbCtx)
	if err != nil {
		return 0, xerror.WrapDaoError(err)
	}
	id, err := rst.LastInsertId()
	if err != nil {
		return 0, xerror.WrapDaoError(err)
	}
	return id, nil
}

func Update(dbCtx lib.DBContexter, table string, where interface{}, data interface{}) (int64, error) {
	build := NewUpdateBuilder(table, Struct2Where(where), Struct2Assign(data))
	sql, args, err := build.Compile()
	if err != nil {
		return 0, xerror.WrapDaoError(err)
	}

	now := time.Now()
	rst, err := dbCtx.Conn().ExecContext(dbCtx, sql, args...)
	sr := &stateful.SQLRecord{
		SQL:  sql,
		Args: args,
		Err:  err,
		Cost: time.Since(now),
	}
	defer sr.Print(dbCtx)
	if err != nil {
		return 0, xerror.WrapDaoError(err)
	}
	var rows int64
	rows, sr.Err = rst.RowsAffected()
	return rows, err
}

func Delete(dbCtx lib.DBContexter, table string, where interface{}) (int64, error) {
	build := NewDeleteBuilder(table, Struct2Where(where))
	sql, args, err := build.Compile()
	if err != nil {
		return 0, xerror.WrapDaoError(err)
	}

	now := time.Now()

	rst, err := dbCtx.Conn().ExecContext(dbCtx, sql, args...)
	sr := &stateful.SQLRecord{
		SQL:  sql,
		Args: args,
		Err:  err,
		Cost: time.Since(now),
	}
	sr.Print(dbCtx)
	if err != nil {
		return 0, xerror.WrapDaoError(err)
	}
	rows, err := rst.RowsAffected()
	if err != nil {
		return 0, xerror.WrapDaoError(err)
	}
	return rows, nil
}
