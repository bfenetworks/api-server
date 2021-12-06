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

// ProductUserUnbindRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ProductUserUnbindEndpoint = &xreq.Endpoint{
	Path:       "/auth/users/{user_name}/products/{product_name}",
	Method:     http.MethodDelete,
	Handler:    xreq.Convert(ProductUserUnbindAction),
	Authorizer: iauth.FA(iauth.FeatureProductUser, iauth.ActionDelete),
}

func productUserUnbindActionProcess(req *http.Request, param *UserNameParam) error {
	user, err := container.AuthenticateManager.FetchUser(req.Context(), &iauth.UserFilter{
		Name: param.UserName,
		Type: &iauth.UserTypeNormal,
	})
	if err != nil {
		return err
	}

	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return err
	}

	if user == nil {
		return xerror.WrapRecordNotExist("User")
	}

	return container.AuthorizeManager.UnBindUserProduct(req.Context(), user, product)
}

var _ xreq.Handler = ProductUserUnbindAction

// ProductUserUnbindAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ProductUserUnbindAction(req *http.Request) (interface{}, error) {
	param, err := newUserNameParam(req)
	if err != nil {
		return nil, err
	}

	return nil, productUserUnbindActionProcess(req, param)
}
