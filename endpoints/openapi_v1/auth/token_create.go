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

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/stateful/container"
)

// SessionKeyByInnerRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var TokenCreateEndpoint = &xreq.Endpoint{
	Path:       "/auth/tokens",
	Method:     http.MethodPost,
	Handler:    xreq.Convert(TokenCreateAction),
	Authorizer: iauth.FA(iauth.FeatureToken, iauth.ActionCreate),
}

type TokenCreateParam struct {
	Name        *string `json:"name" uri:"name" validate:"required,min=1"`
	Scope       *string `json:"scope" validate:"oneof=System Product Support"`
	ProductName *string `json:"product_name" validate:""`
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newTokenCreateParam(req *http.Request) (*TokenCreateParam, error) {
	param := &TokenCreateParam{}
	err := xreq.BindJSON(req, param)
	if err != nil {
		return param, err
	}

	if *param.Scope == iauth.ScopeProduct && param.ProductName == nil {
		return nil, xerror.WrapParamErrorWithMsg("ProductName Required When Scope Is Product")
	}

	return param, err
}

var _ xreq.Handler = TokenCreateAction

// TokenCreateAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func TokenCreateAction(req *http.Request) (interface{}, error) {
	param, err := newTokenCreateParam(req)
	if err != nil {
		return nil, err
	}

	var product *ibasic.Product
	if param.ProductName != nil {
		products, err := container.ProductManager.FetchProducts(req.Context(), &ibasic.ProductFilter{
			Name: param.ProductName,
		})
		if err != nil {
			return nil, err
		}

		if len(products) == 0 {
			return nil, xerror.WrapParamErrorWithMsg("Product Not Exist")
		}
		product = products[0]
	}

	token, err := container.AuthenticateManager.CreateToken(req.Context(), &iauth.TokenParam{
		Name:  param.Name,
		Scope: param.Scope,
	}, product)
	if err != nil {
		return nil, err
	}

	return &TokenCreateData{
		Token: token.Token,
	}, nil
}

type TokenCreateData struct {
	Token string `json:"token"`
}
