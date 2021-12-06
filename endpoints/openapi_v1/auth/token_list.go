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

// InnerUserListRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var TokenListEndpoint = &xreq.Endpoint{
	Path:       "/auth/tokens",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(TokenListAction),
	Authorizer: iauth.FA(iauth.FeatureToken, iauth.ActionReadAll),
}

func tokenListActionProcess(req *http.Request) ([]*TokenData, error) {
	list, err := container.AuthenticateManager.FetchTokens(req.Context(), nil)
	if err != nil {
		return nil, err
	}

	var tokens []*TokenData
	for _, one := range list {
		tokens = append(tokens, newTokenData(one, true))
	}

	return tokens, nil
}

var _ xreq.Handler = TokenListAction

// TokenListAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func TokenListAction(req *http.Request) (interface{}, error) {
	return tokenListActionProcess(req)
}
