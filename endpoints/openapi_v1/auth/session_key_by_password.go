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

// UserIdentifyData Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type UserIdentifyData struct {
	UserName   string   `json:"user_name,omitempty" uri:"user_name"`
	SessionKey string   `json:"session_key,omitempty" uri:"session_key"`
	Roles      []string `json:"roles,omitempty" uri:"roles"`
}

// UserNamePasswordParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type UserNamePasswordParam struct {
	UserName *string `json:"user_name" uri:"user_name" validate:"required,min=1"`
	Password *string `json:"password" uri:"password" validate:"required,min=1"`
}

// SessionKeyByPasswordRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var SessionKeyByPasswordEndpoint = &xreq.Endpoint{
	Path:       "/auth/session-keys/password",
	Method:     http.MethodPost,
	Handler:    xreq.Convert(SessionKeyByPasswordAction),
	Authorizer: nil,
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newUserNamePasswordParam4SessionKeyByPassword(req *http.Request) (*UserNamePasswordParam, error) {
	userNamePasswordParam := &UserNamePasswordParam{}
	err := xreq.BindJSON(req, userNamePasswordParam)
	return userNamePasswordParam, err
}

func sessionKeyByPasswordActionProcess(req *http.Request, param *UserNamePasswordParam) (*UserIdentifyData, error) {
	user, err := container.AuthenticateManager.Authenticate(req.Context(), &iauth.AuthenticateParam{
		Type:     iauth.AuthTypePassword,
		Identify: *param.UserName,
		Extend:   *param.Password,
	})
	if err != nil {
		return nil, err
	}

	return newUserIdentifyData(user, true), nil
}

var _ xreq.Handler = SessionKeyByPasswordAction

// SessionKeyByPasswordAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func SessionKeyByPasswordAction(req *http.Request) (interface{}, error) {
	userNamePasswordParam, err := newUserNamePasswordParam4SessionKeyByPassword(req)
	if err != nil {
		return nil, err
	}

	return sessionKeyByPasswordActionProcess(req, userNamePasswordParam)
}

func newUserIdentifyData(user *iauth.User, withSessionKey bool) *UserIdentifyData {
	roles := []string{}
	for _, one := range user.Roles {
		roles = append(roles, one.Name)
	}

	data := &UserIdentifyData{
		UserName: user.Name,
		Roles:    roles,
	}
	if withSessionKey {
		data.SessionKey = user.SessionKey
	}

	return data
}
