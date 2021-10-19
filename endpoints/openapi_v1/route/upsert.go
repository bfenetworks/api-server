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

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"

	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/model/iroute_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

/*
{
	"Route_rules": [{
			"name": "rule1",
			"description": "message",
			"expression": "req_host_in(\"b.com\") || $TestVar",
			"cluster_name": "cluster-demo"
		}
	],
	"basic_Route_rules": [{
		"host_names": ["*.baidu.com", "t.com"],
		"paths": ["/aaa", "/abc"],
		"cluster_name": "cluster-demo",
		"description": ""
	}],
	"default_cluster_name": "cluster-demo"
}
*/

type AdvanceRouteRule struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Expression  string `json:"expression" validate:"required,min=1"`
	ClusterName string `json:"cluster_name" validate:"required,min=1"`
}

type BasicRouteRule struct {
	HostNames   []string `json:"host_names"`
	Paths       []string `json:"paths"`
	ClusterName string   `json:"cluster_name" validate:"required,min=1"`
	Description string   `json:"description"`
}

type ProductRouteRuleParam struct {
	BasicRouteRules   []*BasicRouteRule   `json:"basic_forward_rules" validate:"dive"`
	AdvanceRouteRules []*AdvanceRouteRule `json:"forward_rules" validate:"dive"`
}

type ProductRouteRuleData struct {
	BasicRouteRules   []*BasicRouteRule   `json:"basic_forward_rules"`
	AdvanceRouteRules []*AdvanceRouteRule `json:"forward_rules"`

	RouteCasesCode int `json:"forward_cases_code,omitempty"`
}

func newProductRouteRuleData(pfr *ProductRouteRuleParam) *ProductRouteRuleData {
	return &ProductRouteRuleData{
		BasicRouteRules:   pfr.BasicRouteRules,
		AdvanceRouteRules: pfr.AdvanceRouteRules,
	}
}

func newProductRouteRule(req *http.Request) (*ProductRouteRuleParam, error) {
	pfr := &ProductRouteRuleParam{}
	if err := xreq.BindJSON(req, pfr); err != nil {
		return nil, err
	}

	return pfr, nil
}

func routeRuleParam2routeRule(p *ProductRouteRuleParam) *iroute_conf.ProductRouteRule {
	afrs := []*iroute_conf.AdvanceRouteRule{}
	for _, one := range p.AdvanceRouteRules {
		afrs = append(afrs, &iroute_conf.AdvanceRouteRule{
			Description: one.Description,
			ClusterName: one.ClusterName,
			Expression:  one.Expression,
			Name:        one.Name,
		})
	}

	bfrs := []*iroute_conf.BasicRouteRule{}
	for _, one := range p.BasicRouteRules {
		clusterName := one.ClusterName
		if clusterName == icluster_conf.RouteAdvancedModeClusterName {
			clusterName = icluster_conf.RouteAdvancedModeClusterName4DP
		}
		bfrs = append(bfrs, &iroute_conf.BasicRouteRule{
			HostNames:   one.HostNames,
			Paths:       one.Paths,
			Description: one.Description,
			ClusterName: clusterName,
		})
	}

	return &iroute_conf.ProductRouteRule{
		BasicRouteRules:   bfrs,
		AdvanceRouteRules: afrs,
	}
}

// UpsertRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UpsertEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/routes",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(UpsertAction),
	Authorizer: iauth.FAP(iauth.FeatureRoute, iauth.ActionUpdate),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newRuleInfoFromReq(req *http.Request) (*ProductRouteRuleParam, error) {
	rule := &ProductRouteRuleParam{}
	err := xreq.BindJSON(req, rule)
	if err != nil {
		return nil, err
	}

	for _, one := range rule.AdvanceRouteRules {
		if one == nil {
			return nil, xerror.WrapParamErrorWithMsg("AdvanceRouteRules element cant be nil")
		}
	}
	for _, one := range rule.BasicRouteRules {
		if one == nil {
			return nil, xerror.WrapParamErrorWithMsg("BasicRouteRules element cant be nil")
		}
	}

	return rule, err
}

func UpsertActionProcess(req *http.Request, rule *ProductRouteRuleParam) (*ProductRouteRuleData, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	ipfr := routeRuleParam2routeRule(rule)

	err = container.RouteRuleManager.UpsertProductRule(req.Context(), product, ipfr)
	if err != nil {
		return nil, err
	}

	return newProductRouteRuleData(rule), nil
}

var _ xreq.Handler = UpsertAction

// UpsertAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func UpsertAction(req *http.Request) (interface{}, error) {
	rule, err := newRuleInfoFromReq(req)
	if err != nil {
		return nil, err
	}

	return UpsertActionProcess(req, rule)
}
