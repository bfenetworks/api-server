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

package traffic

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

// ManualUpdateRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ManualUpdateEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/clusters/{cluster_name}/scheduler/manual",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(ManualUpdateAction),
	Authorizer: iauth.FAP(iauth.FeatureTraffic, iauth.ActionUpdate),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newManualUpdateParam4ManualUpdate(req *http.Request) (map[string]map[string]int, error) {
	param := map[string]map[string]int{}
	err := xreq.JSONDeserializer(req, &param)
	return param, err
}

func manualUpdateActionProcess(req *http.Request, updateParam map[string]map[string]int) (*OneData, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	param, err := newOneParam4One(req)
	if err != nil {
		return nil, err
	}

	cluster, err := container.ClusterManager.FetchCluster(req.Context(), &icluster_conf.ClusterFilter{
		Product: product,
		Name:    &param.ClusterName,
	})
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, xerror.WrapRecordNotExist("Cluster")
	}

	err = container.ClusterManager.UpdateCluster(req.Context(), product, cluster, &icluster_conf.ClusterParam{
		ManualScheduler: updateParam,
	})
	if err != nil {
		return nil, err
	}

	return oneActionProcess(req, &OneParam{
		ClusterName: cluster.Name,
	})
}

var _ xreq.Handler = ManualUpdateAction

// ManualUpdateAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ManualUpdateAction(req *http.Request) (interface{}, error) {
	manualUpdateParam, err := newManualUpdateParam4ManualUpdate(req)
	if err != nil {
		return nil, err
	}

	return manualUpdateActionProcess(req, manualUpdateParam)
}
