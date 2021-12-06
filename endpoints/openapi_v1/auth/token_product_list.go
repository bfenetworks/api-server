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

package auth

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/stateful/container"
)

// ProductTokenListEndpoint route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ProductTokenListEndpoint = &xreq.Endpoint{
	Path:       "/auth/tokens/actions/search-by-product/{product_name}",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ProductTokenListAction),
	Authorizer: iauth.FA(iauth.FeatureToken, iauth.ActionReadAll),
}

var _ xreq.Handler = ProductTokenListAction

// ProductTokenListAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ProductTokenListAction(req *http.Request) (interface{}, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}
	return productTokenListActionProcess(req, &ProductTokenListParam{
		Product: product,
	})
}

func productTokenListActionProcess(req *http.Request, param *ProductTokenListParam) ([]*TokenData, error) {
	list, err := container.AuthorizeManager.FetchProductTokens(req.Context(), param.Product)
	if err != nil {
		return nil, err
	}

	var tokens []*TokenData
	for _, one := range list {
		tokens = append(tokens, newTokenData(one, true))
	}

	return tokens, nil
}

type ProductTokenListParam struct {
	Product *ibasic.Product
}
