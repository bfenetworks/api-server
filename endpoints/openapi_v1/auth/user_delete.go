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

// UserNameParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type UserNameParam struct {
	UserName *string `uri:"user_name" validate:"required,min=1"`
}

// UserDeleteRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UserDeleteEndpoint = &xreq.Endpoint{
	Path:       "/auth/users/{user_name}",
	Method:     http.MethodDelete,
	Handler:    xreq.Convert(UserDeleteAction),
	Authorizer: iauth.FA(iauth.FeatureUser, iauth.ActionDelete),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newUserNameParam(req *http.Request) (*UserNameParam, error) {
	codeLoginParam := &UserNameParam{}
	err := xreq.BindURI(req, codeLoginParam)
	return codeLoginParam, err
}

func userDeleteActionProcess(req *http.Request, param *UserNameParam) error {
	return container.AuthenticateManager.DeleteUser(req.Context(), *param.UserName)
}

var _ xreq.Handler = UserDeleteAction

// UserDeleteAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func UserDeleteAction(req *http.Request) (interface{}, error) {
	param, err := newUserNameParam(req)
	if err != nil {
		return nil, err
	}

	return nil, userDeleteActionProcess(req, param)
}
