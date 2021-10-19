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

package gslb_data

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

type GSLBExportParam struct {
	Version    string `form:"version"`
	BFECluster string `form:"bfe_cluster" validate:"required,min=1"`
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func NewGSLBExportFromReq(req *http.Request) (*GSLBExportParam, error) {
	param := &GSLBExportParam{}
	err := xreq.BindForm(req, param)
	if err != nil {
		return nil, err
	}

	return param, err
}

// ExportRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ExportGSLBEndpoint = &xreq.Endpoint{
	Path:       "/configs/gslb_data/gslb",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ExportGSLBAction),
	Authorizer: iauth.FA(iauth.FeatureRoute, iauth.ActionExport),
}

func ExportGSLBActionProcess(req *http.Request, param *GSLBExportParam) (*icluster_conf.GSLBConf, error) {
	return container.ClusterManager.ExportGSLB(req.Context(), param.Version, param.BFECluster)
}

var _ xreq.Handler = ExportGSLBAction

// ExportGSLBAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ExportGSLBAction(req *http.Request) (interface{}, error) {
	param, err := NewGSLBExportFromReq(req)
	if err != nil {
		return nil, err
	}

	return ExportGSLBActionProcess(req, param)
}
