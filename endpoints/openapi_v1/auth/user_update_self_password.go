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

// UserUpdateSelfPasswordRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UserUpdateSelfPasswordEndpoint = &xreq.Endpoint{
	Path:       "/auth/passwd",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(UserUpdateSelfPasswordAction),
	Authorizer: nil,
}

// UserUpdateSelfPasswordParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type UserUpdateSelfPasswordParam struct {
	OldPassword *string `json:"old_password" validate:"required,min=6"`
	Password    *string `json:"password" validate:"required,min=6"`
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newUserUpdateSelfPasswordParam(req *http.Request) (*UserUpdateSelfPasswordParam, error) {
	param := &UserUpdateSelfPasswordParam{}
	err := xreq.BindJSON(req, param)
	return param, err
}

func userUpdateSelfPasswordActionProcess(req *http.Request, param *UserUpdateSelfPasswordParam) error {
	user, err := iauth.MustGetUser(req.Context())
	if err != nil {
		return err
	}

	return container.AuthenticateManager.UpdateUserPassword(req.Context(), &iauth.PasswordChangeData{
		UserName:    user.Name,
		Password:    *param.Password,
		OldPassword: *param.OldPassword,
	})
}

var _ xreq.Handler = UserUpdateSelfPasswordAction

// UserUpdateSelfPasswordAction action
// anyone can modify self password
func UserUpdateSelfPasswordAction(req *http.Request) (interface{}, error) {
	param, err := newUserUpdateSelfPasswordParam(req)
	if err != nil {
		return nil, err
	}

	return nil, userUpdateSelfPasswordActionProcess(req, param)
}
