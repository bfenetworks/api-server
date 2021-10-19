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

// UserCreateParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type UserCreateParam struct {
	UserName *string  `json:"user_name" uri:"user_name" validate:"required,min=1"`
	Password *string  `json:"password" uri:"password" validate:"required,min=6"`
	Roles    []string `json:"roles" uri:"roles" validate:"min=1"`
}

// UserCreateRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UserCreateEndpoint = &xreq.Endpoint{
	Path:       "/auth/users",
	Method:     http.MethodPost,
	Handler:    xreq.Convert(UserCreateAction),
	Authorizer: iauth.FA(iauth.FeatureUser, iauth.ActionCreate),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newUserCreateParam(req *http.Request) (*UserCreateParam, error) {
	codeLoginParam := &UserCreateParam{}
	err := xreq.BindJSON(req, codeLoginParam)
	return codeLoginParam, err
}

func userCreateActionProcess(req *http.Request, param *UserCreateParam) error {
	roles, err := iauth.RoleList(param.Roles)
	if err != nil {
		return err
	}

	return container.AuthenticateManager.CreateUser(req.Context(), &iauth.UserParam{
		Name:     param.UserName,
		Password: param.Password,
		Roles:    roles,
	})
}

var _ xreq.Handler = UserCreateAction

// UserCreateAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func UserCreateAction(req *http.Request) (interface{}, error) {
	param, err := newUserCreateParam(req)
	if err != nil {
		return nil, err
	}

	return nil, userCreateActionProcess(req, param)
}
