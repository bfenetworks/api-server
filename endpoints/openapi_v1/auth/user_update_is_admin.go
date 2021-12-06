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

// UserUpdateIsAdminRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UserUpdateIsAdminEndpoint = &xreq.Endpoint{
	Path:       "/auth/users/{user_name}/is_admin",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(UserUpdateIsAdminAction),
	Authorizer: iauth.FA(iauth.FeatureUser, iauth.ActionUpdate),
}

// UserUpdateIsAdminParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type UserUpdateIsAdminParam struct {
	UserName *string `uri:"user_name" validate:"required,min=1"`
	IsAdmin  bool    `json:"is_admin"`
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newUserUpdateIsAdminParam(req *http.Request) (*UserUpdateIsAdminParam, error) {
	param := &UserUpdateIsAdminParam{}
	err := xreq.Bind(req, param)
	return param, err
}

func userUpdateIsAdminActionProcess(req *http.Request, param *UserUpdateIsAdminParam) error {
	user, err := container.AuthenticateManager.FetchUser(req.Context(), &iauth.UserFilter{
		Name: param.UserName,
	})
	if err != nil {
		return err
	}
	if user == nil {
		return xerror.WrapRecordNotExist("User")
	}
	return container.AuthorizeManager.UpdateUserIsAdmin(req.Context(), user, param.IsAdmin)
}

var _ xreq.Handler = UserUpdateIsAdminAction

// UserUpdateIsAdminAction action
// Admin update other user's password
func UserUpdateIsAdminAction(req *http.Request) (interface{}, error) {
	param, err := newUserUpdateIsAdminParam(req)
	if err != nil {
		return nil, err
	}

	return nil, userUpdateIsAdminActionProcess(req, param)
}
