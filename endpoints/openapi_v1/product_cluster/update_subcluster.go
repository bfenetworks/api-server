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

// BindSubCluster Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type BindSubCluster struct {
	ClusterName *string  `uri:"cluster_name" validate:"required,min=1"`
	SubClusters []string `json:"sub_clusters" uri:"sub_clusters" validate:"required,min=1"`
}

// BindSubClusterRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var BindSubClusterEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/clusters/{cluster_name}/sub-clusters",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(BindSubClusterAction),
	Authorizer: iauth.FAP(iauth.FeatureProductCluster, iauth.ActionUpdate),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newBindSubCluster4BindSubCluster(req *http.Request) (*BindSubCluster, error) {
	bindSubCluster := &BindSubCluster{}
	err := xreq.Bind(req, bindSubCluster)
	return bindSubCluster, err
}

func bindSubClusterActionProcess(req *http.Request, param *BindSubCluster) (*ClusterData, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	cluster, err := container.ClusterManager.FetchCluster(req.Context(), &icluster_conf.ClusterFilter{
		Name:    param.ClusterName,
		Product: product,
	})
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, xerror.WrapRecordNotExist("Cluster")
	}

	if err = container.ClusterManager.RebindSubCluster(req.Context(), product, cluster, param.SubClusters); err != nil {
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

var _ xreq.Handler = BindSubClusterAction

// BindSubClusterAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func BindSubClusterAction(req *http.Request) (interface{}, error) {
	bindSubCluster, err := newBindSubCluster4BindSubCluster(req)
	if err != nil {
		return nil, err
	}

	return bindSubClusterActionProcess(req, bindSubCluster)
}
