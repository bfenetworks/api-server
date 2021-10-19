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

package protocol

import (
	"net/http"

	"github.com/bfenetworks/api-server/endpoints/innerapi_v1/export_util"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/iprotocol"
	"github.com/bfenetworks/api-server/stateful/container"
)

// ExportRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ServertCertExportEndpoint = &xreq.Endpoint{
	Path:       "/configs/protocol/server_cert_conf",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ServerCertExportAction),
	Authorizer: iauth.FA(iauth.FeatureCert, iauth.ActionExport),
}

func exportActionProcess(req *http.Request, param *export_util.ExportParam) (*iprotocol.ServerCertConf, error) {
	return container.CertificateManager.ExportServerCert(req.Context(), param.Version)
}

var _ xreq.Handler = ServerCertExportAction

// ServerCertExportAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ServerCertExportAction(req *http.Request) (interface{}, error) {
	param, err := export_util.NewExportFromReq(req)
	if err != nil {
		return nil, err
	}

	return exportActionProcess(req, param)
}
