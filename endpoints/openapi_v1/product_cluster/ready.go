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

package product_cluster

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

// ReadyRspParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type ReadyRspParam struct {
	Name  string `json:"name"`
	Ready bool   `json:"ready"`
}

// ReadyRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ReadyEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/clusters/{cluster_name}/ready",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ReadyAction),
	Authorizer: iauth.FAP(iauth.FeatureProductCluster, iauth.ActionRead),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newOneParam4Ready(req *http.Request) (*OneParam, error) {
	param := &OneParam{}
	err := xreq.BindURI(req, param)
	return param, err
}

func readyActionProcess(req *http.Request, param *OneParam) (*ReadyRspParam, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	cluster, err := container.ClusterManager.FetchCluster(req.Context(), &icluster_conf.ClusterFilter{
		Name:    param.Name,
		Product: product,
	})
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, xerror.WrapRecordNotExist("Cluster")
	}

	return &ReadyRspParam{
		Name:  cluster.Name,
		Ready: cluster.Ready,
	}, nil
}

var _ xreq.Handler = ReadyAction

// ReadyAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ReadyAction(req *http.Request) (interface{}, error) {
	param, err := newOneParam4Ready(req)
	if err != nil {
		return nil, err
	}

	return readyActionProcess(req, param)
}
