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

// UserUpdatePasswordRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UserUpdatePasswordEndpoint = &xreq.Endpoint{
	Path:       "/auth/users/{user_name}/passwd",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(UserUpdatePasswordAction),
	Authorizer: iauth.FA(iauth.FeatureUser, iauth.ActionUpdate),
}

// UserUpdatePasswordParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type UserUpdatePasswordParam struct {
	UserName    *string `uri:"user_name" validate:"required,min=1"`
	OldPassword string  `json:"old_password" validate:"omitempty"`
	Password    *string `json:"password" validate:"required,min=6"`
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newUserUpdatePasswordParam(req *http.Request) (*UserUpdatePasswordParam, error) {
	param := &UserUpdatePasswordParam{}
	err := xreq.Bind(req, param)
	return param, err
}

func userUpdatePasswordActionProcess(req *http.Request, param *UserUpdatePasswordParam) error {
	user, err := iauth.MustGetVistor(req.Context())
	if err != nil {
		return err
	}

	if user.GetName() == *param.UserName {
		if param.OldPassword == "" {
			return xerror.WrapParamErrorWithMsg("Want Old Password")
		}
	}

	return container.AuthenticateManager.UpdateUserPassword(req.Context(), &iauth.PasswordChangeData{
		UserName:    *param.UserName,
		OldPassword: param.OldPassword,
		Password:    *param.Password,
	})
}

var _ xreq.Handler = UserUpdatePasswordAction

// UserUpdatePasswordAction action
// Admin update other user's password
func UserUpdatePasswordAction(req *http.Request) (interface{}, error) {
	param, err := newUserUpdatePasswordParam(req)
	if err != nil {
		return nil, err
	}

	return nil, userUpdatePasswordActionProcess(req, param)
}
