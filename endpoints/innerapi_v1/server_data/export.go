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

package server_data

import (
	"net/http"

	"github.com/bfenetworks/api-server/endpoints/innerapi_v1/export_util"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/iroute_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

// ExportRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ExportEndpoint = &xreq.Endpoint{
	Path:       "/configs/tls_conf/server_data_conf",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ExportAction),
	Authorizer: iauth.FA(iauth.FeatureRoute, iauth.ActionExport),
}

func ExportActionProcess(req *http.Request, param *export_util.ExportParam) (*iroute_conf.RouteRuleExportData, error) {
	return container.RouteRuleManager.ExportRouteRule(req.Context(), param.Version)
}

var _ xreq.Handler = ExportAction

// ExportAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ExportAction(req *http.Request) (interface{}, error) {
	param, err := export_util.NewExportFromReq(req)
	if err != nil {
		return nil, err
	}

	return ExportActionProcess(req, param)
}
