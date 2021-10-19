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

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

// ListRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ListEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/clusters",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ListAction),
	Authorizer: iauth.FAP(iauth.FeatureProductCluster, iauth.ActionRead),
}

func listActionProcess(req *http.Request) ([]*ClusterData, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	list, err := container.ClusterManager.FetchClusterList(req.Context(), &icluster_conf.ClusterFilter{
		Product: product,
	})
	if err != nil {
		return nil, err
	}

	rsp := []*ClusterData{}
	for _, one := range list {
		rsp = append(rsp, clusterModel2Control(one))
	}
	return rsp, nil
}

var _ xreq.Handler = ListAction

// ListAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ListAction(req *http.Request) (interface{}, error) {
	return listActionProcess(req)
}
