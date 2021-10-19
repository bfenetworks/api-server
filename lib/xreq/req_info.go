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

package xreq

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type requestInfoKey string

var (
	_requestInfoKey requestInfoKey = "request_info"
)

type RequestInfo struct {
	StartTime time.Time
	Duration  time.Duration

	URLPath  string
	Method   string
	ClientIP string
	LogID    string

	URLPattern string

	StatusCode int
	RetMsg     string
	ErrDetail  string
}

func InitRequestInfo(ctx context.Context, req *http.Request) (context.Context, *RequestInfo) {
	requestInfo := GetRequestInfo(ctx)
	if requestInfo != nil {
		return ctx, requestInfo
	}
	requestInfo = &RequestInfo{
		StartTime:  time.Now(),
		URLPath:    req.URL.Path,
		ClientIP:   req.Header.Get("ClientIp"),
		Method:     req.Method,
		StatusCode: 200,
	}
	return context.WithValue(ctx, _requestInfoKey, requestInfo), requestInfo
}

func GetRequestInfo(ctx context.Context) *RequestInfo {
	if v := ctx.Value(_requestInfoKey); v != nil {
		return v.(*RequestInfo)
	}

	return nil
}

func (requestInfo *RequestInfo) String() string {
	return fmt.Sprintf("[%s]cost_ms[%d] method[%s] pattern[%s] path[%s] client_ip[%s] status_code[%d] ret_msg[%s] err_detail[%s]]",
		requestInfo.LogID, requestInfo.Duration.Milliseconds(), requestInfo.Method, requestInfo.URLPattern, requestInfo.URLPath, requestInfo.ClientIP,
		requestInfo.StatusCode, requestInfo.RetMsg, requestInfo.ErrDetail)
}
