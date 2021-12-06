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

// UpdateBasicRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UpdateBasicEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/clusters/{cluster_name}",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(UpdateAction),
	Authorizer: iauth.FAP(iauth.FeatureProductCluster, iauth.ActionUpdate),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newUpdateParam4Update(req *http.Request) (*UpsertParam, error) {
	param := &UpsertParam{}
	if err := xreq.Bind(req, param); err != nil {
		return nil, err
	}

	if param.SubClusters != nil {
		return nil, xerror.WrapParamErrorWithMsg("Invoke Bind API To Modify SubCluster")
	}

	if param.Scheduler != nil {
		return nil, xerror.WrapParamErrorWithMsg("Invoke Scheduler API To Modify Scheduler Setting")
	}

	if ss := param.StickySessions; ss != nil {
		if *ss.HashStrategy != clusterHashStrategyClientIPOnly && ss.HashHeader == nil {
			return nil, xerror.WrapParamErrorWithMsg("StickySessions.HashHeader Want Be Set")
		}
	}

	return param, nil
}

func updateActionProcess(req *http.Request, param *UpsertParam) (*ClusterData, error) {
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

	if err := container.ClusterManager.UpdateCluster(req.Context(), product, cluster, clusterParamControlModel(param)); err != nil {
		return nil, err
	}

	cluster, err = container.ClusterManager.FetchCluster(req.Context(), &icluster_conf.ClusterFilter{
		ID: &cluster.ID,
	})
	if err != nil {
		return nil, err
	}
	return clusterModel2Control(cluster), nil
}

var _ xreq.Handler = UpdateAction

// UpdateAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func UpdateAction(req *http.Request) (interface{}, error) {
	param, err := newUpdateParam4Update(req)
	if err != nil {
		return nil, err
	}

	return updateActionProcess(req, param)
}
