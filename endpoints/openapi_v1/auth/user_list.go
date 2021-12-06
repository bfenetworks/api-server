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

// UserListRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UserListEndpoint = &xreq.Endpoint{
	Path:       "/auth/users",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(UserListAction),
	Authorizer: iauth.FA(iauth.FeatureUser, iauth.ActionReadAll),
}

func userListActionProcess(req *http.Request) ([]*UserData, error) {
	list, err := container.AuthenticateManager.FetchUserList(req.Context(), nil)
	if err != nil {
		return nil, err
	}

	users := []*UserData{}
	for _, one := range list {
		users = append(users, newUserData(one))
	}

	return users, nil
}

var _ xreq.Handler = UserListAction

// UserListAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func UserListAction(req *http.Request) (interface{}, error) {
	return userListActionProcess(req)
}
