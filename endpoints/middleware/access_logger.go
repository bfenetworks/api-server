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

package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/stateful"
)

func UpdateMonitor(req *http.Request, requestInfo *xreq.RequestInfo) {
	method := strings.ToLower(requestInfo.Method)
	code := fmt.Sprint(requestInfo.StatusCode)

	stateful.MetricAPIAccessCounter.With(prometheus.Labels{
		"pattern":     requestInfo.URLPattern,
		"method":      method,
		"status_code": code,
	}).Inc()
	stateful.MetricAPICostHisCounter.With(prometheus.Labels{
		"pattern":     requestInfo.URLPattern,
		"method":      method,
		"status_code": code,
	}).Add(float64(requestInfo.Duration.Milliseconds()))
}

// API access logger
type LoggerMiddleWare struct{}

func NewLoggerMiddleWare() *LoggerMiddleWare {
	return &LoggerMiddleWare{}
}

func GetClientIp(r *http.Request) string {
	return r.Header.Get("Clientip")
}

func (l *LoggerMiddleWare) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)

	ctx := r.Context()
	requestInfo := xreq.GetRequestInfo(ctx)
	requestInfo.Duration = time.Since(requestInfo.StartTime)

	nrw := rw.(negroni.ResponseWriter)
	requestInfo.StatusCode = nrw.Status()

	Record(requestInfo)
	UpdateMonitor(r, requestInfo)
}

func Record(requestInfo *xreq.RequestInfo) {
	if requestInfo.StatusCode == 200 {
		stateful.AccessLogger.Info(requestInfo.String())
	} else {
		stateful.AccessLogger.Warn(requestInfo.String())
	}
}
