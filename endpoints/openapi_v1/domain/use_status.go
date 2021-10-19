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

// OneParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type OneParam struct {
	DomainName string `json:"domain_name" uri:"domain_name"`
}

// UseStatusRsp Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type UseStatusRsp struct {
	DepType string `json:"dep_type"`
	DepName string `json:"dep_name"`
	BeUsed  bool   `json:"be_used"`
}

// UseStatusRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UseStatusEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/domains/{domain_name}/use-status",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(UseStatusAction),
	Authorizer: iauth.FAP(iauth.FeatureDomain, iauth.ActionRead),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newOneParam4UseStatus(req *http.Request) (*OneParam, error) {
	param := &OneParam{}
	err := xreq.BindURI(req, param)
	return param, err
}

func useStatusActionProcess(req *http.Request, param *OneParam) (*UseStatusRsp, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	list, err := container.DomainManager.DomainList(req.Context(), &iroute_conf.DomainFilter{
		Product: product,
		Name:    &param.DomainName,
	})
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, xerror.WrapRecordNotExist("Domain")
	}

	domain := list[0]
	info, err := container.DomainManager.BeUsed(req.Context(), product, domain)
	if err != nil {
		return nil, err
	}

	if info == nil {
		return &UseStatusRsp{
			BeUsed: false,
		}, nil
	}

	typ, name := info.Dependent()
	return &UseStatusRsp{
		BeUsed:  true,
		DepType: typ,
		DepName: name,
	}, nil
}

var _ xreq.Handler = UseStatusAction

// UseStatusAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func UseStatusAction(req *http.Request) (interface{}, error) {
	param, err := newOneParam4UseStatus(req)
	if err != nil {
		return nil, err
	}

	return useStatusActionProcess(req, param)
}
