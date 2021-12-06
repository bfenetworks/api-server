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
	"strings"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

// UpsertParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type UpsertParam struct {
	Name      *string     `json:"name" uri:"instance_pool_name" validate:"required,min=2"`
	Instances []*Instance `json:"instances" uri:"instances" validate:"min=1,dive"`
}

// CreateRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var CreateEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/instance-pools",
	Method:     http.MethodPost,
	Handler:    xreq.Convert(CreateAction),
	Authorizer: iauth.FAP(iauth.FeatureProductPool, iauth.ActionCreate),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func NewUpsertParam(req *http.Request) (*UpsertParam, error) {
	param := &UpsertParam{}
	err := xreq.Bind(req, param)
	if err != nil {
		return nil, err
	}

	return param, err
}

var _ xreq.Handler = CreateAction

// CreateAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func CreateAction(req *http.Request) (interface{}, error) {
	param, err := NewUpsertParam(req)
	if err != nil {
		return nil, err
	}

	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(*param.Name, product.Name+".") {
		return nil, xerror.WrapParamErrorWithMsg("Want Prefix %s.", product.Name)
	}
	if len(*param.Name) == len(product.Name)+1 {
		return nil, xerror.WrapParamErrorWithMsg("Want Pool Name")
	}

	oneData, err := CreateProcess(req, product, param)
	if err != nil {
		return nil, err
	}

	return NewOneData(oneData), nil
}

func Instancesc2i(is []*Instance) []icluster_conf.Instance {
	rst := []icluster_conf.Instance{}
	for _, instance := range is {
		port := 0
		if instance.Ports != nil {
			port = instance.Ports["Default"]
		}
		rst = append(rst, icluster_conf.Instance{
			HostName: instance.Hostname,
			IP:       instance.IP,
			Weight:   instance.Weight,
			Ports:    instance.Ports,
			Port:     port,
			Tags:     instance.Tags,
		})
	}

	return rst
}

func CreateProcess(req *http.Request, product *ibasic.Product, param *UpsertParam) (*icluster_conf.Pool, error) {
	return container.PoolManager.CreateProductPool(req.Context(), product, &icluster_conf.PoolParam{
		Name:      param.Name,
		Instances: Instancesc2i(param.Instances),
	})
}
