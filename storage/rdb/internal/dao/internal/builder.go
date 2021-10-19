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
	"fmt"
	"reflect"
	"strings"

	"github.com/didi/gendry/builder"
)

var tagName = "db"

func Struct2Where(raw interface{}) map[string]interface{} {
	return struct2map(raw, false)
}

func Struct2AssignList(raws ...interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(raws))
	for _, raw := range raws {
		rst = append(rst, struct2map(raw, true))
	}
	return rst
}

func Struct2Assign(raw interface{}) map[string]interface{} {
	return struct2map(raw, true)
}

func struct2map(raw interface{}, ignoreOpt bool) map[string]interface{} {
	rst := map[string]interface{}{}
	if raw == nil {
		return rst
	}
	structType := reflect.TypeOf(raw)
	if kind := structType.Kind(); kind == reflect.Ptr || kind == reflect.Interface {
		structType = structType.Elem()
	}
	structVal := reflect.ValueOf(raw)
	if structVal.IsZero() {
		return rst
	}
	if structVal.Kind() == reflect.Ptr {
		structVal = structVal.Elem()
	}
	if structType.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < structVal.NumField(); i++ {
		valField := structVal.Field(i)
		valFieldKind := valField.Kind()
		if valFieldKind == reflect.Ptr {
			if valField.IsNil() {
				continue
			}
			valField = valField.Elem()
		}
		if valFieldKind == reflect.Slice {
			if valField.IsZero() {
				continue
			}
		}
		typeField := structType.Field(i)
		dbTag := typeField.Tag.Get(tagName)
		if dbTag == "-" {
			continue
		}
		key, opt := tagSplitter(dbTag)
		if key == "" {
			key = typeField.Name
		}
		if ignoreOpt {
			rst[key] = valField.Interface()
		} else {
			if opt == "" || opt == "=" {
				rst[key] = valField.Interface()
			} else {
				rst[key+" "+opt] = valField.Interface()
			}
		}
	}
	return rst
}

func tagSplitter(dbTag string) (key, opt string) {
	if dbTag == "" {
		return "", ""
	}
	i := strings.Index(dbTag, ",")
	if i == -1 {
		return dbTag, ""
	}
	return strings.TrimSpace(dbTag[:i]), strings.TrimSpace(dbTag[i+1:])
}

type SQLBuilder interface {
	Compile() (sql string, args []interface{}, err error)
}

func NewSelectBuilder(table string, where map[string]interface{}, fields []string) SQLBuilder {
	return &SelectBuilder{
		table:  table,
		where:  where,
		fields: fields,
	}
}

type SelectBuilder struct {
	table  string
	where  map[string]interface{}
	fields []string
}

func (s *SelectBuilder) From(table string) *SelectBuilder {
	s.table = table
	return s
}

func (s *SelectBuilder) Where(wheres map[string]interface{}) *SelectBuilder {
	s.where = wheres
	return s
}

func (s *SelectBuilder) Select(selectFields []string) *SelectBuilder {
	s.fields = selectFields
	return s
}

func (s *SelectBuilder) Compile() (sql string, args []interface{}, err error) {
	return builder.BuildSelect(s.table, s.where, s.fields)
}

var _ SQLBuilder = (*SelectBuilder)(nil)

func NewDeleteBuilder(table string, wheres map[string]interface{}) SQLBuilder {
	return &DeleteBuilder{
		table: table,
		where: wheres,
	}
}

type DeleteBuilder struct {
	table string
	where map[string]interface{}
}

func (d *DeleteBuilder) From(table string) *DeleteBuilder {
	d.table = table
	return d
}

func (d *DeleteBuilder) Where(wheres map[string]interface{}) *DeleteBuilder {
	d.where = wheres
	return d
}

func (d *DeleteBuilder) Compile() (sql string, args []interface{}, err error) {
	return builder.BuildDelete(d.table, d.where)
}

var _ SQLBuilder = (*DeleteBuilder)(nil)

func NewUpdateBuilder(table string, wheres map[string]interface{}, assigns map[string]interface{}) SQLBuilder {
	return &UpdateBuilder{
		table:   table,
		where:   wheres,
		assigns: assigns,
	}
}

type UpdateBuilder struct {
	table   string
	where   map[string]interface{}
	assigns map[string]interface{}
}

func (u *UpdateBuilder) From(table string) *UpdateBuilder {
	u.table = table
	return u
}

func (u *UpdateBuilder) Where(wheres map[string]interface{}) *UpdateBuilder {
	u.where = wheres
	return u
}

func (u *UpdateBuilder) Assign(assign map[string]interface{}) *UpdateBuilder {
	u.assigns = assign
	return u
}

func (u *UpdateBuilder) Compile() (sql string, args []interface{}, err error) {
	return builder.BuildUpdate(u.table, u.where, u.assigns)
}

var _ SQLBuilder = (*UpdateBuilder)(nil)

func NewInsertBuilder(table string, assigns []map[string]interface{}) SQLBuilder {
	return &InsertBuilder{
		table:   table,
		assigns: assigns,
		typ:     insertCommon,
	}
}

func NewIgnoreInsertBuilder(table string, assigns []map[string]interface{}) SQLBuilder {
	return &InsertBuilder{
		table:   table,
		assigns: assigns,
		typ:     insertIgnore,
	}
}

func NewReplaceInsertBuilder(table string, assigns []map[string]interface{}) SQLBuilder {
	return &InsertBuilder{
		table:   table,
		assigns: assigns,
		typ:     insertReplace,
	}
}

type InsertBuilder struct {
	table   string
	assigns []map[string]interface{}
	typ     int
}

func (i *InsertBuilder) Values(assigns []map[string]interface{}) *InsertBuilder {
	i.assigns = assigns
	return i
}

func (i *InsertBuilder) InsertInto(table string) *InsertBuilder {
	i.table = table
	i.typ = insertCommon
	return i
}

func (i *InsertBuilder) InsertIgnoreInto(table string) *InsertBuilder {
	i.table = table
	i.typ = insertIgnore
	return i
}

func (i *InsertBuilder) ReplaceInto(table string) *InsertBuilder {
	i.table = table
	i.typ = insertReplace
	return i
}

const (
	insertCommon = iota
	insertIgnore
	insertReplace
)

func (i *InsertBuilder) Compile() (sql string, args []interface{}, err error) {
	switch i.typ {
	case insertCommon:
		return builder.BuildInsert(i.table, i.assigns)
	case insertIgnore:
		return builder.BuildInsertIgnore(i.table, i.assigns)
	case insertReplace:
		return builder.BuildReplaceInsert(i.table, i.assigns)
	default:
		return "", nil, fmt.Errorf("unknown type=%d", i.typ)
	}
}

var _ SQLBuilder = (*InsertBuilder)(nil)
