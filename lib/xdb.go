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

package lib

import (
	"context"
	"database/sql"
	"strings"
)

type DBContexter interface {
	context.Context
	Conn() *sql.DB
}

type DBContext struct {
	conn *sql.DB
	tx   *sql.Tx

	context.Context
}

func NewDBContext(ctx context.Context, conn *sql.DB) *DBContext {
	return &DBContext{
		conn:    conn,
		Context: ctx,
	}
}

func (ctx *DBContext) Conn() *sql.DB {
	return ctx.conn
}

func (ctx *DBContext) BeginTrans() error {
	if ctx.tx != nil {
		return nil
	}

	var err error
	ctx.tx, err = ctx.conn.Begin()
	return err
}

type Op struct {
	blockWrite bool
	openTxn    bool
}

func BlockWrite() *Op {
	return &Op{
		blockWrite: true,
	}
}

func WantBlockWrite(ops ...*Op) bool {
	for _, op := range ops {
		if op.blockWrite {
			return true
		}
	}

	return false
}

func OpenTxn() *Op {
	return &Op{
		openTxn: true,
	}
}

func WantOpenTxn(ops ...*Op) bool {
	for _, op := range ops {
		if op.openTxn {
			return true
		}
	}

	return false
}

type DBContextFactory func(ctx context.Context, ops ...*Op) (*DBContext, error)

func RDBTxnExecute(dc *DBContext, handler func(context.Context) error) error {
	var err error
	if dc.tx == nil {
		dc.tx, err = dc.conn.Begin()
		if err != nil {
			return err
		}
	}

	err = handler(dc)
	return commitOrRollback(dc.tx, err)
}

func commitOrRollback(tx *sql.Tx, err error) error {
	if err != nil {
		if strings.Contains(err.Error(), "invalid connection") {
			return err
		}

		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func Transaction(conn *sql.DB, do func(*sql.Tx) error) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	if err := do(tx); err != nil {

		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func DuplicateEntryError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Error 1062: Duplicate entry")
}
