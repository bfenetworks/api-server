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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/stateful"
)

type Handler func(req *http.Request) (interface{}, error)

func Convert(h Handler) func(req *http.Request) *Result {
	return func(req *http.Request) *Result {
		data, err := h(req)
		return &Result{
			OriginErr: err,
			Data:      data,
		}
	}
}

func RawConvert(h Handler) func(req *http.Request) *Result {
	return func(req *http.Request) *Result {
		data, err := h(req)
		return &Result{
			OriginErr: err,
			Data:      data,
			Render: func(w http.ResponseWriter, req *http.Request, res *Result) {
				if err == nil {
					w.Header().Add("Content-Type", "application/octet-stream")
					w.Write(data.([]byte))
					return
				}

				// json response
				Render(w, req, res)
			},
		}
	}
}

type Endpoint struct {
	Path    string
	Method  string
	Handler func(*http.Request) *Result

	RegisterHandler func(*mux.Router) *mux.Route

	Authorizer iauth.Authorizer
}

func (ep *Endpoint) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rst := ep.Handler(req)

	render := Render
	if rst != nil && rst.Render != nil {
		render = rst.Render
	}

	render(rw, req, rst)
}

func (re *Endpoint) String() string {
	return fmt.Sprintf("[%-7s --> %-20s ", re.Method+"]", re.Path)
}

func (ep *Endpoint) Register(router *mux.Router) *mux.Router {
	if authorizer := ep.Authorizer; authorizer != nil {
		router = router.NewRoute().Subrouter()
		router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				ctx := req.Context()
				GetRequestInfo(ctx).URLPattern = ep.Path

				user, err := iauth.MustGetUser(ctx)
				if err != nil {
					ErrorRender(err, rw, req)
					return
				}

				ok, err := iauth.Authorizate(ctx, authorizer, user)
				if err != nil {
					ErrorRender(err, rw, req)
					return
				}

				if !ok {
					ErrorRender(xerror.WrapAuthorizateFailErrorWithMsg("Auth Deny"), rw, req)
					return
				}

				next.ServeHTTP(rw, req)
			})
		})
	}

	if ep.RegisterHandler == nil {
		router.Handle(ep.Path, ep).Methods(ep.Method)
	} else {
		ep.RegisterHandler(router).Handler(ep)
	}

	return router
}

// Result used for result process
type Result struct {
	OriginErr error `json:"-"`

	Code      int         `json:"ErrNum"`         // http status code
	ErrMsg    string      `json:"ErrMsg"`         // return message if failed
	Data      interface{} `json:"Data,omitempty"` // return data if success
	ManualURL string      `json:"Manual,omitempty"`

	Render func(w http.ResponseWriter, req *http.Request, res *Result) `json:"-"`

	ErrorDetail string `json:"-"` // detail error msg, presents in access log not response
}

// Error error interface
func (res *Result) Error() string {
	if res.OriginErr == nil {
		return "<nil>"
	}

	return res.OriginErr.Error()
}

func (res *Result) parseError() {
	err := res.OriginErr
	if err == nil {
		return
	}

	rr := xerror.Resolve(err)
	res.Code = rr.ErrNo
	res.ErrMsg = rr.Type
	if rr.Msg != "" {
		res.ErrMsg += (": " + rr.Msg)
	}
	res.ErrorDetail = rr.FullMsg()
}

func (res *Result) IsSucc() bool {
	if res.OriginErr != nil {
		return false
	}
	return res.Code == 200 || res.Code == 0
}

func ErrorRender(err error, w http.ResponseWriter, req *http.Request) {
	rst := &Result{
		OriginErr: err,
	}

	Render(w, req, rst)
}

func Render(w http.ResponseWriter, req *http.Request, res *Result) {
	if res.Code == 0 {
		res.Code = 200
	}
	if res.ErrMsg == "" {
		res.ErrMsg = "success"
	}

	isSucc := res.IsSucc()
	if !isSucc {
		res.parseError()
	}

	_, requestInfo := InitRequestInfo(req.Context(), req)
	requestInfo.RetMsg = res.ErrMsg

	if stateful.DefaultConfig.RunTime.Debug {
		requestInfo.ErrDetail = res.ErrorDetail
	}

	if !isSucc {
		res.ErrMsg = stateful.TryMappingErrMsg(req, res.ErrMsg)
	}

	// process according to http code
	switch {
	case 200 <= res.Code && res.Code < 300:
		res.Code = http.StatusOK
	case 400 <= res.Code && res.Code < 500:
	case res.Code == 555 || res.Code == 556:
	// case 300 <= res.Code && res.Code < 400:
	// case 500 <= res.Code && res.Code < 600:
	// dont change code value
	default:
		res.Code = http.StatusInternalServerError
	}

	requestInfo.StatusCode = res.Code

	w.Header().Add("Req-ID", requestInfo.LogID)

	var rspContent []byte
	w.Header().Add("Content-Type", "application/json")
	var data interface{} = res
	var err error
	rspContent, err = json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(res.Code)
	w.Write(rspContent)
}
