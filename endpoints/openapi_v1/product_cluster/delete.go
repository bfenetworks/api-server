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

// DeleteRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var DeleteEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/clusters/{cluster_name}",
	Method:     http.MethodDelete,
	Handler:    xreq.Convert(DeleteAction),
	Authorizer: iauth.FAP(iauth.FeatureProductCluster, iauth.ActionDelete),
}

func deleteActionProcess(req *http.Request, param *OneParam) (*ClusterData, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	one, err := container.ClusterManager.FetchCluster(req.Context(), &icluster_conf.ClusterFilter{
		Name:    param.Name,
		Product: product,
	})
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, xerror.WrapRecordNotExist("Cluster")
	}

	if err := container.ClusterManager.DeleteCluster(req.Context(), product, one); err != nil {
		return nil, err
	}

	return clusterModel2Control(one), nil
}

var _ xreq.Handler = DeleteAction

// DeleteAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func DeleteAction(req *http.Request) (interface{}, error) {
	param, err := newOneParam4One(req)
	if err != nil {
		return nil, err
	}

	return deleteActionProcess(req, param)
}

// CreateCluster Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type CreateCluster struct {
	Name *string `json:"name" validate:"required,min=1"`
}
