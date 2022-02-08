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

// OneParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type OneParam struct {
	InstancePoolName string `json:"instance_pool_name" uri:"instance_pool_name" validate:"required,min=2"`
}

// Instance Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type Instance struct {
	Hostname string            `json:"hostname" uri:"hostname" validate:"required,min=2"`
	IP       string            `json:"ip" uri:"ip" validate:"required,ip"`
	Weight   int64             `json:"weight" uri:"weight" validate:"min=0,max=100"`
	Ports    map[string]int    `json:"ports" uri:"ports" validate:"required,min=1"`
	Tags     map[string]string `json:"tags" uri:"tags" validate:"required,min=1"`
}

// OneData Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type OneData struct {
	Name      string      `json:"name" uri:"name"`
	Instances []*Instance `json:"instances" uri:"instances"`
}

func NewOneData(pool *icluster_conf.Pool, pis *icluster_conf.PoolInstances) *OneData {
	is := []*Instance{}
	if pis != nil {
		for _, one := range pis.Instances {
			is = append(is, &Instance{
				Hostname: one.HostName,
				IP:       one.IP,
				Weight:   one.Weight,
				Ports:    one.Ports,
				Tags:     one.Tags,
			})
		}
	}

	return &OneData{
		Name:      pool.Name,
		Instances: is,
	}
}

// OneRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var OneEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/instance-pools/{instance_pool_name}",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(OneAction),
	Authorizer: iauth.FAP(iauth.FeatureProductPool, iauth.ActionRead),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func NewOneParam(req *http.Request) (*OneParam, error) {
	param := &OneParam{}
	err := xreq.BindURI(req, param)
	return param, err
}

var _ xreq.Handler = OneAction

// OneAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func OneAction(req *http.Request) (interface{}, error) {
	param, err := NewOneParam(req)
	if err != nil {
		return nil, err
	}

	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	one, err := container.PoolManager.FetchProductPool(req.Context(), product, param.InstancePoolName)
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
	return NewOneData(one, pism[one.Name]), nil
}
