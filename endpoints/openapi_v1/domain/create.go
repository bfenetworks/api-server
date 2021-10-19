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

package domain

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/iroute_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

// CreateParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type CreateParam struct {
	Name *string `json:"name" uri:"name" validate:"required,min=2"`
}

// CreateRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var CreateEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/domains",
	Method:     http.MethodPost,
	Handler:    xreq.Convert(CreateAction),
	Authorizer: iauth.FAP(iauth.FeatureDomain, iauth.ActionCreate),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newCreateParam4Create(req *http.Request) (*CreateParam, error) {
	param := &CreateParam{}
	err := xreq.BindJSON(req, param)
	return param, err
}

func CreateActionProcess(req *http.Request, param *CreateParam) (*OneData, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	err = container.DomainManager.CreateDomain(req.Context(), product, &iroute_conf.DomainParam{
		ProductID: &product.ID,
		Name:      param.Name,
	})
	if err != nil {
		return nil, err
	}

	return &OneData{
		Name: *param.Name,
	}, err
}

var _ xreq.Handler = CreateAction

// CreateAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func CreateAction(req *http.Request) (interface{}, error) {
	param, err := newCreateParam4Create(req)
	if err != nil {
		return nil, err
	}

	return CreateActionProcess(req, param)
}
