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

// ProdcutUserListRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ProdcutUserListEndpoint = &xreq.Endpoint{
	Path:       "/auth/products/{product_name}/users",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ProdcutUserListAction),
	Authorizer: iauth.FA(iauth.FeatureProductUser, iauth.ActionReadAll),
}

func productUserListActionProcess(req *http.Request) ([]*UserIdentifyData, error) {
	list, err := container.ProductAuthorizateManager.GrantedUsers(req.Context())
	if err != nil {
		return nil, err
	}

	var rst []*UserIdentifyData
	for _, one := range list {
		rst = append(rst, newUserIdentifyData(one, false))
	}

	return rst, nil
}

var _ xreq.Handler = ProdcutUserListAction

// ProdcutUserListAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ProdcutUserListAction(req *http.Request) (interface{}, error) {
	return productUserListActionProcess(req)
}
