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

// UpdateParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type UpdateParam struct {
	Name           *string `json:"name" uri:"name" validate:"min=2"`
	Description    *string `json:"description" uri:"description" validate:"omitempty,min=2"`
	SubClusterName *string `uri:"sub_cluster_name" validate:"required,min=2"`
}

// UpdateRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UpdateEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/sub-clusters/{sub_cluster_name}",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(UpdateAction),
	Authorizer: iauth.FAP(iauth.FeatureSubCluster, iauth.ActionUpdate),
}

var _ xreq.Handler = UpdateAction

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newUpdateParam4Update(req *http.Request) (*UpdateParam, error) {
	param := &UpdateParam{}
	err := xreq.Bind(req, param)
	return param, err
}

// UpdateAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func UpdateAction(req *http.Request) (interface{}, error) {
	param, err := newUpdateParam4Update(req)
	if err != nil {
		return nil, err
	}

	return updateActionProcess(req, param)
}

func updateActionProcess(req *http.Request, param *UpdateParam) (*OneData, error) {
	// get product info
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	// check if exist one
	oldOne, err := one(req, &icluster_conf.SubClusterFilter{
		Name:    param.SubClusterName,
		Product: product,
	})
	if err != nil {
		return nil, err
	}
	if oldOne == nil {
		return nil, xerror.WrapRecordNotExist("SubCluster")
	}

	err = container.SubClusterManager.UpdateSubCluster(req.Context(), oldOne, &icluster_conf.SubClusterParam{
		Description: param.Description,
	})
	if err != nil {
		return nil, err
	}

	newSubCluster, err := one(req, &icluster_conf.SubClusterFilter{
		ID: &oldOne.ID,
	})
	if err != nil {
		return nil, err
	}

	return newOneData(newSubCluster), nil
}
