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

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/iroute_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

// DeleteParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type DeleteParam struct {
	DomainName string `json:"domain_name" uri:"domain_name" validate:"required,min=2"`
}

// DeleteRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var DeleteEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/domains/{domain_name}",
	Method:     http.MethodDelete,
	Handler:    xreq.Convert(DeleteAction),
	Authorizer: iauth.FAP(iauth.FeatureDomain, iauth.ActionDelete),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newDeleteParam4Delete(req *http.Request) (*DeleteParam, error) {
	param := &DeleteParam{}
	err := xreq.BindURI(req, param)
	return param, err
}

var _ xreq.Handler = DeleteAction

// DeleteAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func DeleteAction(req *http.Request) (interface{}, error) {
	param, err := newDeleteParam4Delete(req)
	if err != nil {
		return nil, err
	}
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	list, err := container.DomainManager.DomainList(req.Context(), &iroute_conf.DomainFilter{
		Name: &param.DomainName,
	})
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, xerror.WrapRecordNotExist()
	}

	if err := container.DomainManager.DeleteDomain(req.Context(), product, list[0]); err != nil {
		return nil, err
	}

	return &OneData{
		Name: list[0].Name,
	}, nil
}
