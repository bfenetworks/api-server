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

package bfe_pool

import (
	"net/http"

	"github.com/bfenetworks/api-server/endpoints/openapi_v1/product_pool"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

// OneRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var OneEndpoint = &xreq.Endpoint{
	Path:       "/bfe-pools/{instance_pool_name}",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(OneAction),
	Authorizer: iauth.FA(iauth.FeatureBFEPool, iauth.ActionReadAll),
}

var _ xreq.Handler = OneAction

// OneAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func OneAction(req *http.Request) (interface{}, error) {
	param, err := product_pool.NewOneParam(req)
	if err != nil {
		return nil, err
	}

	one, err := container.PoolManager.FetchBFEPool(req.Context(), param.InstancePoolName)
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, xerror.WrapRecordNotExist("Instance Pool")
	}

	pism, err := container.PoolInstancesManager.BatchFetchInstances(req.Context(), []*icluster_conf.Pool{one})
	if err != nil {
		return nil, err
	}
	return product_pool.NewOneData(one, pism[one.Name]), nil

}
