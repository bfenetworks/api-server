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
	"runtime"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/stateful"
)

type Recovery struct {
	StackAll  bool
	StackSize int
}

// NewRecovery returns a new instance of Recovery
func NewRecovery() *Recovery {
	rec := &Recovery{StackAll: false, StackSize: 1024 * 8}
	return rec
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	ctx := lib.NewLogContext(req.Context())
	ctx, requestInfo := xreq.InitRequestInfo(ctx, req)
	requestInfo.LogID = lib.GainLogID(ctx)

	req = req.WithContext(ctx)

	defer func() {
		if err := recover(); err != nil {
			stateful.MetricPaincCounter.Inc()

			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]
			stackString := string(stack)

			requestInfo.StatusCode = 500
			requestInfo.RetMsg = "system error"
			requestInfo.ErrDetail = fmt.Sprintf("PANIC: ERR:%s STACK:%s", err, strings.ReplaceAll(stackString, "\n", "\\n"))

			stateful.AccessLogger.Warn(requestInfo.String())

			r := &xreq.Result{
				Code:   requestInfo.StatusCode,
				ErrMsg: requestInfo.RetMsg,
			}

			xreq.Render(rw, req, r)
		}
	}()

	next(rw, req)
}

func McConvert(handler negroni.Handler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(rw, r, next.ServeHTTP)
		})
	}
}

var (
	MCRecovery = McConvert(NewRecovery())
	MCLogger   = McConvert(NewLoggerMiddleWare())
	MCCors     = McConvert(NewCors())
)
