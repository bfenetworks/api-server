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
	"fmt"
	"strings"

	"github.com/bfenetworks/bfe/bfe_basic"
	"github.com/bfenetworks/bfe/bfe_http"
)

func ReqFactory(url string, reqHeader, rspHeader map[string]string, method string) (*bfe_basic.Request, error) {
	httpReq, err := bfe_http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("bfe_http.NewRequest(): %s", err.Error())
	}

	// add case's header to http request
	for key, val := range reqHeader {
		if key == "Host" {
			httpReq.Host = val
		} else {
			httpReq.Header.Add(key, val)
		}
	}

	req := &bfe_basic.Request{
		HttpRequest: httpReq,
		HttpResponse: &bfe_http.Response{
			Header: bfe_http.Header{},
		},
		Session: &bfe_basic.Session{
			Proto:    httpReq.Proto,
			IsSecure: strings.HasPrefix(strings.ToLower(url), "https"),
		},
	}

	for key, val := range rspHeader {
		req.HttpResponse.Header.Add(key, val)
	}

	return req, nil
}
