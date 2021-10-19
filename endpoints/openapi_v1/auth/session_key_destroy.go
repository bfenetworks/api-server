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
	"github.com/bfenetworks/api-server/stateful/container"
)

// SessionKeyParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type SessionKeyParam struct {
	SessionKey *string `json:"session_key" uri:"session_key" validate:"required,min=1"`
}

// SessionKeyDestroyRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var SessionKeyDestroyEndpoint = &xreq.Endpoint{
	Path:       "/auth/session-keys/{session_key}",
	Method:     http.MethodDelete,
	Handler:    xreq.Convert(SessionKeyDestroyAction),
	Authorizer: nil,
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newSessionKeyParam4SessionKeyDestroy(req *http.Request) (*SessionKeyParam, error) {
	sessionKeyParam := &SessionKeyParam{}
	err := xreq.BindURI(req, sessionKeyParam)
	return sessionKeyParam, err
}

func sessionKeyDestroyActionProcess(req *http.Request, param *SessionKeyParam) error {
	return container.AuthenticateManager.DestroySessionKey(req.Context(), *param.SessionKey)
}

var _ xreq.Handler = SessionKeyDestroyAction

// SessionKeyDestroyAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func SessionKeyDestroyAction(req *http.Request) (interface{}, error) {
	param, err := newSessionKeyParam4SessionKeyDestroy(req)
	if err != nil {
		return nil, err
	}

	return nil, sessionKeyDestroyActionProcess(req, param)
}
