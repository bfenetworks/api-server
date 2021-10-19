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

const tClusterTableName = "clusters"

// TCluster Query Result
type TCluster struct {
	ID                     int64     `db:"id"`
	Name                   string    `db:"name"`
	Description            string    `db:"description"`
	ProductID              int64     `db:"product_id"`
	MaxIdleConnPerHost     int16     `db:"max_idle_conn_per_host"`
	TimeoutConnServ        int32     `db:"timeout_conn_serv"`
	TimeoutResponseHeader  int32     `db:"timeout_response_header"`
	TimeoutReadbodyClient  int32     `db:"timeout_readbody_client"`
	TimeoutReadClientAgain int32     `db:"timeout_read_client_again"`
	TimeoutWriteClient     int32     `db:"timeout_write_client"`
	HealthcheckSchem       string    `db:"healthcheck_schem"`
	HealthcheckInterval    int32     `db:"healthcheck_interval"`
	HealthcheckFailnum     int32     `db:"healthcheck_failnum"`
	HealthcheckHost        string    `db:"healthcheck_host"`
	HealthcheckUri         string    `db:"healthcheck_uri"`
	HealthcheckStatuscode  int32     `db:"healthcheck_statuscode"`
	ClientipCarry          bool      `db:"clientip_carry"`
	PortCarry              bool      `db:"port_carry"`
	MaxRetryInCluster      int8      `db:"max_retry_in_cluster"`
	MaxRetryCrossCluster   int8      `db:"max_retry_cross_cluster"`
	Ready                  bool      `db:"ready"`
	HashStrategy           int32     `db:"hash_strategy"`
	CookieKey              string    `db:"cookie_key"`
	HashHeader             string    `db:"hash_header"`
	SessionSticky          bool      `db:"session_sticky"`
	ReqWriteBufferSize     int32     `db:"req_write_buffer_size"`
	ReqFlushInterval       int32     `db:"req_flush_interval"`
	ResFlushInterval       int32     `db:"res_flush_interval"`
	CancelOnClientClose    bool      `db:"cancel_on_client_close"`
	FailureStatus          bool      `db:"failure_status"`
	CreatedAt              time.Time `db:"created_at"`
	UpdatedAt              time.Time `db:"updated_at"`
}

// TClusterOne Query One
// return (nil, nil) if record not existed
func TClusterOne(dbCtx lib.DBContexter, where *TClusterParam) (*TCluster, error) {
	t := &TCluster{}
	err := internal.QueryOne(dbCtx, tClusterTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TClusterList Query Multiple
func TClusterList(dbCtx lib.DBContexter, where *TClusterParam) ([]*TCluster, error) {
	t := []*TCluster{}
	err := internal.QueryList(dbCtx, tClusterTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TClusterParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TClusterParam struct {
	IDs []int64 `db:"id,in"`

	ID                     *int64     `db:"id"`
	Name                   *string    `db:"name"`
	Names                  []string   `db:"name,in"`
	Description            *string    `db:"description"`
	ProductID              *int64     `db:"product_id"`
	MaxIdleConnPerHost     *int16     `db:"max_idle_conn_per_host"`
	TimeoutConnServ        *int32     `db:"timeout_conn_serv"`
	TimeoutResponseHeader  *int32     `db:"timeout_response_header"`
	TimeoutReadbodyClient  *int32     `db:"timeout_readbody_client"`
	TimeoutReadClientAgain *int32     `db:"timeout_read_client_again"`
	TimeoutWriteClient     *int32     `db:"timeout_write_client"`
	HealthcheckSchem       *string    `db:"healthcheck_schem"`
	HealthcheckInterval    *int32     `db:"healthcheck_interval"`
	HealthcheckFailnum     *int32     `db:"healthcheck_failnum"`
	HealthcheckHost        *string    `db:"healthcheck_host"`
	HealthcheckUri         *string    `db:"healthcheck_uri"`
	HealthcheckStatuscode  *int32     `db:"healthcheck_statuscode"`
	ClientipCarry          *bool      `db:"clientip_carry"`
	PortCarry              *bool      `db:"port_carry"`
	MaxRetryInCluster      *int8      `db:"max_retry_in_cluster"`
	MaxRetryCrossCluster   *int8      `db:"max_retry_cross_cluster"`
	Ready                  *bool      `db:"ready"`
	HashStrategy           *int32     `db:"hash_strategy"`
	CookieKey              *string    `db:"cookie_key"`
	HashHeader             *string    `db:"hash_header"`
	SessionSticky          *bool      `db:"session_sticky"`
	ReqWriteBufferSize     *int32     `db:"req_write_buffer_size"`
	ReqFlushInterval       *int32     `db:"req_flush_interval"`
	ResFlushInterval       *int32     `db:"res_flush_interval"`
	CancelOnClientClose    *bool      `db:"cancel_on_client_close"`
	FailureStatus          *bool      `db:"failure_status"`
	CreatedAt              *time.Time `db:"created_at"`
	UpdatedAt              *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TClusterCreate One/Multiple
func TClusterCreate(dbCtx lib.DBContexter, data ...*TClusterParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tClusterTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tClusterTableName, list...)
}

// TClusterUpdate Update One
func TClusterUpdate(dbCtx lib.DBContexter, val, where *TClusterParam) (int64, error) {
	return internal.Update(dbCtx, tClusterTableName, where, val)
}

// TClusterDelete Delete One/Multiple
func TClusterDelete(dbCtx lib.DBContexter, where *TClusterParam) (int64, error) {
	return internal.Delete(dbCtx, tClusterTableName, where)
}
