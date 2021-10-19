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

package product

// Auto Gen By ctrl, Modify as U Need

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/stateful/container"
)

var _ xreq.Handler = ProductDeleteAction

// ProductDeleteEndpoint
// Auto Gen By ctrl, Modify as U Need
var ProductDeleteEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}",
	Method:     http.MethodDelete,
	Handler:    xreq.Convert(ProductDeleteAction),
	Authorizer: iauth.FAP(iauth.FeatureProduct, iauth.ActionDelete),
}

// ProductDeleteAction delete product
// Auto Gen By ctrl, Modify as U Need
func ProductDeleteAction(req *http.Request) (interface{}, error) {
	p, err := getProduct(req)
	if err != nil {
		return nil, err
	}

	if err := container.ProductManager.DeleteProduct(req.Context(), p); err != nil {
		return nil, err
	}

	return newProductData(p), nil
}
