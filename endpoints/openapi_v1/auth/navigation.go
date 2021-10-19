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
	"github.com/bfenetworks/api-server/stateful"
)

var NavigationEndpoint = &xreq.Endpoint{
	Path:       "/meta",
	Method:     "GET",
	Handler:    xreq.Convert(NavigationProcess),
	Authorizer: nil,
}

var _ xreq.Handler = NavigationProcess

func NavigationProcess(req *http.Request) (interface{}, error) {
	return map[string]interface{}{
		"nav":  stateful.DefaultConfig.Depends.Role2Nav(),
		"icon": stateful.DefaultConfig.Depends.UIIcon,
		"logo": stateful.DefaultConfig.Depends.UILogo,
	}, nil
}
