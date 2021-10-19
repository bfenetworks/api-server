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

package endpoints

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"

	"github.com/bfenetworks/api-server/endpoints/innerapi_v1"
	"github.com/bfenetworks/api-server/endpoints/middleware"
	"github.com/bfenetworks/api-server/endpoints/openapi_v1"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/stateful"
)

func fileHandler(root http.Dir, fs http.Handler) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/") {
			r.URL.Path = "/" + r.URL.Path
		}
		r.URL.Path = path.Clean(r.URL.Path)

		requestInfo := xreq.GetRequestInfo(r.Context())
		defer middleware.Record(requestInfo)

		f, err := root.Open(r.URL.Path)
		if err != nil {
			if os.IsNotExist(err) {
				r.URL.Path = "/"
			}
		} else {
			f.Close()
		}

		fs.ServeHTTP(rw, r)
	}
}

func RegisterRouters(router *mux.Router) {
	fileServerRoot, err := filepath.Abs(stateful.DefaultConfig.RunTime.StaticFilePath)
	if err != nil {
		panic(err)
	}
	root := http.Dir(fileServerRoot)
	fs := http.FileServer(root)
	fh := fileHandler(root, fs)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.MCRecovery.Middleware(http.HandlerFunc(fh)).ServeHTTP(w, r)
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := &xreq.Result{Code: 405, ErrMsg: "Method Not Allowed"}
		xreq.Render(w, r, res)
	})

	router.Use(middleware.MCRecovery)
	router.Use(middleware.MCLogger)
	router.Use(middleware.MCCors)

	openapi_v1.RegisterEndpoints(router)
	innerapi_v1.RegisterRouter(router)
}
