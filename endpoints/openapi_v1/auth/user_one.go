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

// UserOneRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UserOneEndpoint = &xreq.Endpoint{
	Path:       "/auth/users/{user_name}",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(UserOneAction),
	Authorizer: iauth.FA(iauth.FeatureUser, iauth.ActionReadAll),
}

func userOneActionProcess(req *http.Request, param *UserNameParam) (*UserData, error) {
	user, err := container.AuthenticateManager.FetchUser(req.Context(), &iauth.UserFilter{
		Name: param.UserName,
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, xerror.WrapRecordNotExist("User")
	}

	ups, err := container.AuthorizeManager.FetchVisitorProductList(req.Context(), &iauth.Visitor{
		User: user,
	})
	if err != nil {
		return nil, err
	}

	userData := newUserData(user)
	for _, one := range ups {
		userData.Products = append(userData.Products, one.Name)
	}

	return userData, nil
}

var _ xreq.Handler = UserOneAction

// UserOneAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func UserOneAction(req *http.Request) (interface{}, error) {
	param, err := newUserNameParam(req)
	if err != nil {
		return nil, err
	}

	return userOneActionProcess(req, param)
}

// UserData Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type UserData struct {
	UserName string `json:"user_name,omitempty"`
	IsAdmin  bool   `json:"is_admin"`

	SessionKey string   `json:"session_key,omitempty"`
	Products   []string `json:"products,omitempty"`
}

func newUserData(user *iauth.User) *UserData {
	data := &UserData{
		UserName: user.Name,
		IsAdmin:  user.Admin,
	}

	return data
}
