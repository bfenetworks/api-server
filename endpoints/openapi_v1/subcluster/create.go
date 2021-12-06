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

// CreateParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type CreateParam struct {
	Name         *string `json:"name" uri:"name" validate:"required,min=2"`
	InstancePool *string `json:"instance_pool" uri:"instance_pool" validate:"required,min=2"`
	Description  *string `json:"description" uri:"description" validate:"required,min=2"`
}

// CreateRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var CreateEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/sub-clusters",
	Method:     http.MethodPost,
	Handler:    xreq.Convert(CreateAction),
	Authorizer: iauth.FAP(iauth.FeatureSubCluster, iauth.ActionCreate),
}

var _ xreq.Handler = CreateAction

// CreateAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func CreateAction(req *http.Request) (interface{}, error) {
	param, err := newCreateParam4Create(req)
	if err != nil {
		return nil, err
	}

	return CreateProcess(req, param)
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newCreateParam4Create(req *http.Request) (*CreateParam, error) {
	param := &CreateParam{}
	err := xreq.BindJSON(req, param)
	return param, err
}

func CreateProcess(req *http.Request, param *CreateParam) (*OneData, error) {
	// get product info
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	err = container.SubClusterManager.CreateSubCluster(req.Context(), product,
		&icluster_conf.SubClusterParam{
			Name:        param.Name,
			Product:     product,
			PoolName:    param.InstancePool,
			Description: param.Description,
		})
	if err != nil {
		return nil, err
	}

	newCluster, err := one(req, &icluster_conf.SubClusterFilter{
		Name:    param.Name,
		Product: product,
	})
	if err != nil {
		return nil, err
	}

	return newOneData(newCluster), nil
}
