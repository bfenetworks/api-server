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

package middleware

import (
	"net/http"
	"strings"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/stateful/container"
)

type UserProbeParam struct {
	UserID   *int64  `bind:"product_id"`
	UserName *string `bind:"product_name"`
}

func newUserProbeParam(req *http.Request) (*UserProbeParam, error) {
	param := &UserProbeParam{}
	err := xreq.BindURI(req, param)

	return param, err
}

func UserProbeAction(req *http.Request) (*http.Request, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return req, nil
	}

	ss := strings.SplitN(authHeader, " ", 2)
	if len(ss) != 2 {
		return req, xerror.WrapParamErrorWithMsg("Bad Format Header Authorization")
	}

	param := &iauth.AuthenticateParam{
		Type:     ss[0],
		Identify: ss[1],
	}

	user, err := container.AuthenticateManager.Authenticate(req.Context(), param)
	if err != nil {
		return nil, err
	}

	return req.WithContext(iauth.NewUserContext(req.Context(), user)), nil
}
