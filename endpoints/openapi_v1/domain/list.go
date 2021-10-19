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

// ListRsp Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type ListRsp struct {
	List []*string `json:"list" uri:"list"`
}

// OneData Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type OneData struct {
	Name string `json:"name" uri:"name"`
}

// ListRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ListEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/domains",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ListAction),
	Authorizer: iauth.FAP(iauth.FeatureDomain, iauth.ActionRead),
}
var _ xreq.Handler = ListAction

// ListAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ListAction(req *http.Request) (interface{}, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	list, err := container.DomainManager.DomainList(req.Context(), &iroute_conf.DomainFilter{
		Product: product,
	})
	if err != nil {
		return nil, err
	}

	rst := make([]string, len(list))
	for i, one := range list {
		rst[i] = one.Name
	}
	return rst, nil
}
