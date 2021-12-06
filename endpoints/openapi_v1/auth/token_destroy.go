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
	"github.com/bfenetworks/api-server/stateful/container"
)

// TokenNameParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type TokenNameParam struct {
	TokenName *string `json:"token_name" uri:"token_name" validate:"required,min=1"`
}

// TokenDestroyRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var TokenDestroyEndpoint = &xreq.Endpoint{
	Path:       "/auth/tokens/{token_name}",
	Method:     http.MethodDelete,
	Handler:    xreq.Convert(TokenDestroyAction),
	Authorizer: iauth.FA(iauth.FeatureToken, iauth.ActionDelete),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newTokenNameParam(req *http.Request) (*TokenNameParam, error) {
	param := &TokenNameParam{}
	err := xreq.BindURI(req, param)
	return param, err
}

func TokenDestroyActionProcess(req *http.Request, param *TokenNameParam) error {
	tokens, err := container.AuthenticateManager.FetchTokens(req.Context(), &iauth.TokenFilter{
		Name: param.TokenName,
	})
	if err != nil {
		return err
	}

	if len(tokens) != 1 {
		return xerror.WrapRecordNotExist("Token")
	}

	return container.AuthenticateManager.DeleteToken(req.Context(), tokens[0])
}

var _ xreq.Handler = TokenDestroyAction

// TokenDestroyAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func TokenDestroyAction(req *http.Request) (interface{}, error) {
	param, err := newTokenNameParam(req)
	if err != nil {
		return nil, err
	}

	return nil, TokenDestroyActionProcess(req, param)
}
