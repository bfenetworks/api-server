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

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

// ListRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ListEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/sub_clusters",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ListAction),
	Authorizer: iauth.FAP(iauth.FeatureSubCluster, iauth.ActionRead),
}

var _ xreq.Handler = ListAction

// ListAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ListAction(req *http.Request) (interface{}, error) {
	return listActionProcess(req)
}

func listActionProcess(req *http.Request) ([]*OneData, error) {
	// get product info
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	subClusterList, err := container.SubClusterManager.SubClusterList(req.Context(),
		&icluster_conf.SubClusterFilter{
			Product: product,
		})

	if err != nil {
		return nil, err
	}

	if len(subClusterList) == 0 {
		return nil, nil
	}

	list := make([]*OneData, len(subClusterList))
	for i, one := range subClusterList {
		list[i] = newOneData(one)

	}
	return list, nil
}
