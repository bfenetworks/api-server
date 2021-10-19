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

package subcluster

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
	Path:       "/products/{product_name}/sub_clusters/{sub_cluster_name}",
	Method:     http.MethodDelete,
	Handler:    xreq.Convert(DeleteAction),
	Authorizer: iauth.FAP(iauth.FeatureSubCluster, iauth.ActionDelete),
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

func deleteActionProcess(req *http.Request, param *OneParam) (*OneData, error) {
	// get product info
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	subClusterName := param.SubClusterName
	subCluster, err := one(req, &icluster_conf.SubClusterFilter{
		Name:    &subClusterName,
		Product: product,
	})
	if err != nil {
		return nil, err
	}
	if subCluster == nil {
		return nil, xerror.WrapRecordNotExist()
	}

	err = container.SubClusterManager.DeleteSubCluster(req.Context(), subCluster)

	if err != nil {
		return nil, err
	}

	return newOneData(subCluster), nil
}
