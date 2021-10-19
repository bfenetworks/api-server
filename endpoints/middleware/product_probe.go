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

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/stateful/container"
)

type ProductProbeParam struct {
	ProductID   *int64  `uri:"product_id"`
	ProductName *string `uri:"product_name"`
}

func newProductProbeParam(req *http.Request) (*ProductProbeParam, error) {
	param := &ProductProbeParam{}
	err := xreq.BindURI(req, param)

	return param, err
}

func ProductProbeAction(req *http.Request) (*http.Request, error) {
	param, err := newProductProbeParam(req)
	if err != nil {
		return nil, err
	}

	if param.ProductID == nil && param.ProductName == nil {
		return req, nil
	}

	products, err := container.ProductManager.FetchProducts(req.Context(), &ibasic.ProductFilter{
		ID:   param.ProductID,
		Name: param.ProductName,
	})
	if err != nil {
		return nil, err
	}

	if len(products) != 1 {
		return nil, xerror.WrapParamErrorWithMsg("Product Not Exist")
	}

	return req.WithContext(ibasic.NewProductContext(req.Context(), products[0])), nil
}

func NewMiddleWareFunc(handler func(*http.Request) (*http.Request, error)) func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		newReq, err := handler(r)
		if err != nil {
			xreq.ErrorRender(err, rw, r)
			return
		}

		next(rw, newReq)
	}
}
