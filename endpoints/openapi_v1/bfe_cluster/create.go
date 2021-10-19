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

package bfe_cluster

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/stateful/container"
)

// BFEClusterCreateParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type BFEClusterCreateParam struct {
	Name *string `json:"name" uri:"name" validate:"required,min=1"`
	Pool *string `json:"pool" uri:"pool" validate:"required,min=1"`
}

// CreateRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var CreateEndpoint = &xreq.Endpoint{
	Path:       "/bfe-clusters",
	Method:     http.MethodPost,
	Handler:    xreq.Convert(CreateAction),
	Authorizer: iauth.FA(iauth.FeatureBFECluster, iauth.ActionCreate),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newBFEClusterCreateParam4Create(req *http.Request) (*BFEClusterCreateParam, error) {
	bfeClusterCreateParam := &BFEClusterCreateParam{}
	err := xreq.BindJSON(req, bfeClusterCreateParam)
	return bfeClusterCreateParam, err
}

func createActionProcess(req *http.Request, param *BFEClusterCreateParam) error {
	return container.BFEClusterManager.CreateBFECluster(req.Context(), &ibasic.BFEClusterParam{
		Name:     param.Name,
		Pool:     param.Pool,
		Capacity: lib.PInt64(0),
	})
}

var _ xreq.Handler = CreateAction

// CreateAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func CreateAction(req *http.Request) (interface{}, error) {
	param, err := newBFEClusterCreateParam4Create(req)
	if err != nil {
		return nil, err
	}

	return nil, createActionProcess(req, param)
}
