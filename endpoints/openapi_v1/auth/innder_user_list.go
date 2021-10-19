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

// InnerUserListRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var InnerUserListEndpoint = &xreq.Endpoint{
	Path:       "/auth/inner-users",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(InnerUserListAction),
	Authorizer: iauth.FA(iauth.FeatureUser, iauth.ActionReadAll),
}

func innerUserListActionProcess(req *http.Request) ([]string, error) {
	list, err := container.AuthenticateManager.FetchInnerUserList(req.Context())
	if err != nil {
		return nil, err
	}

	var ss []string
	for _, one := range list {
		ss = append(ss, one.SessionKey)
	}

	return ss, nil
}

var _ xreq.Handler = InnerUserListAction

// InnerUserListAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func InnerUserListAction(req *http.Request) (interface{}, error) {
	return innerUserListActionProcess(req)
}
