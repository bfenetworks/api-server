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

package route

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/model/iroute_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

// ListRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ListEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/routes",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ListAction),
	Authorizer: iauth.FAP(iauth.FeatureRoute, iauth.ActionRead),
}

func routeRule2routeRuleParam(p *iroute_conf.ProductRouteRule) *ProductRouteRuleParam {
	afrs := []*AdvanceRouteRule{}
	for _, one := range p.AdvanceRouteRules {
		afrs = append(afrs, &AdvanceRouteRule{
			Description: one.Description,
			ClusterName: one.ClusterName,
			Expression:  one.Expression,
			Name:        one.Name,
		})
	}

	bfrs := []*BasicRouteRule{}
	for _, one := range p.BasicRouteRules {
		clusterName := one.ClusterName
		if clusterName == icluster_conf.RouteAdvancedModeClusterName4DP {
			clusterName = icluster_conf.RouteAdvancedModeClusterName
		}
		bfrs = append(bfrs, &BasicRouteRule{
			HostNames:   one.HostNames,
			Paths:       one.Paths,
			Description: one.Description,
			ClusterName: clusterName,
		})
	}

	return &ProductRouteRuleParam{
		BasicRouteRules:   bfrs,
		AdvanceRouteRules: afrs,
	}
}

func listActionProcess(req *http.Request) (*ProductRouteRuleData, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	rule, err := container.RouteRuleManager.FetchProductRule(req.Context(), product)
	if err != nil {
		return nil, err
	}

	if rule == nil {
		return nullRule, nil
	}

	return newProductRouteRuleData(routeRule2routeRuleParam(rule)), nil
}

var nullRule = &ProductRouteRuleData{
	BasicRouteRules:   []*BasicRouteRule{},
	AdvanceRouteRules: []*AdvanceRouteRule{},
}

var _ xreq.Handler = ListAction

// ListAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ListAction(req *http.Request) (interface{}, error) {
	return listActionProcess(req)
}
