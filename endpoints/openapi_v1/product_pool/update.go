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

package product_pool

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

// UpdateRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UpdateEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/instance-pools/{instance_pool_name}",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(UpdateAction),
	Authorizer: iauth.FAP(iauth.FeatureProductPool, iauth.ActionUpdate),
}

var _ xreq.Handler = UpdateAction

// UpdateAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func UpdateAction(req *http.Request) (interface{}, error) {
	param, err := NewUpsertParam(req)
	if err != nil {
		return nil, err
	}

	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}
	one, err := container.PoolManager.FetchProductPool(req.Context(), product, *param.Name)
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, xerror.WrapRecordNotExist("Instance Pool")
	}

	err = container.PoolManager.UpdateProductPool(req.Context(), product, one, &icluster_conf.PoolParam{
		Instances: Instancesc2i(param.Instances),
	})

	one.Instances = Instancesc2i(param.Instances)

	return NewOneData(one), err
}
