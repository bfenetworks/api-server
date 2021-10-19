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
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bfenetworks/api-server/lib/xreq"
)

func convert(handler func(*http.Request) (*http.Request, error)) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			newReq, err := handler(r)
			if err != nil {
				xreq.ErrorRender(err, rw, r)
				return
			}

			next.ServeHTTP(rw, newReq)
		})
	}
}

var (
	McProductProbe = convert(ProductProbeAction)
	McUserProbe    = convert(UserProbeAction)
)
