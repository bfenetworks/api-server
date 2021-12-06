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
	"github.com/bfenetworks/api-server/stateful/container"
)

// InnerUserOneRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var TokenOneEndpoint = &xreq.Endpoint{
	Path:       "/auth/tokens/{token_name}",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(TokenOneAction),
	Authorizer: iauth.FA(iauth.FeatureToken, iauth.ActionReadAll),
}

func tokenOneActionProcess(req *http.Request) (*TokenData, error) {
	param, err := newTokenNameParam(req)
	if err != nil {
		return nil, err
	}

	list, err := container.AuthenticateManager.FetchTokens(req.Context(), &iauth.TokenFilter{
		Name: param.TokenName,
	})
	if err != nil {
		return nil, err
	}

	return newTokenData(list[0], true), nil
}

var _ xreq.Handler = TokenOneAction

// TokenOneAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func TokenOneAction(req *http.Request) (interface{}, error) {
	return tokenOneActionProcess(req)
}

type TokenData struct {
	Name        string `json:"name"`
	ProductName string `json:"product_name,omitempty"`
	Token       string `json:"token,omitempty"`
	Scope       string `json:"scope"`
}

func newTokenData(token *iauth.Token, withToken bool) *TokenData {
	t := token.Token
	if !withToken {
		t = ""
	}

	productName := ""
	if token.Product != nil {
		productName = token.Product.Name
	}

	return &TokenData{
		Name:  token.Name,
		Token: t,
		Scope: token.Scope,

		ProductName: productName,
	}
}
