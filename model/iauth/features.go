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

package iauth

type (
	Feature string
	Action  int64
)

type FeatureAuthorition struct {
	Feature Feature
	Action  Action
}

type Authorization struct {
	FeatureAuthorizer *FeatureAuthorition
	ValidateProduct   bool
}

func NewFeatureAuthorization(f Feature, a Action) *Authorization {
	return &Authorization{
		FeatureAuthorizer: &FeatureAuthorition{
			Feature: f,
			Action:  a,
		},
	}
}

func NewFeatureAuthorizerWithFactoryWithProduct(f Feature, a Action) *Authorization {
	tmp := NewFeatureAuthorization(f, a)
	tmp.ValidateProduct = true

	return tmp
}

var (
	FA  = NewFeatureAuthorization
	FAP = NewFeatureAuthorizerWithFactoryWithProduct
)

const (
	ActionDeny    Action = 1 << iota // 000001
	ActionRead                       // 000010
	ActionReadAll                    // 000010
	ActionUpdate                     // 000100
	ActionCreate                     // 001000
	ActionDelete                     // 010000
	ActionExport                     // 100000
)

func (a Action) Revoke(b Action) Action {
	return a & (^b)
}

func (a Action) Grant(b Action) Action {
	return a | b
}

func (a Action) IsAllowed(b Action) bool {
	return a&b != 0
}

const (
	// global resource
	FeatureProxyPool  Feature = "ProxyPool"
	FeatureBFECluster Feature = "BFECluster"
	FeatureBFEPool    Feature = "BFEPool"
	FeatureArea       Feature = "Area"
	FeatureDomain     Feature = "Domain"
	FeatureProduct    Feature = "Product"
	FeatureExtraFile  Feature = "ExtraFile"

	// product resource
	FeatureProductPool       Feature = "ProductPool"
	FeatureRoute             Feature = "Route"
	FeatureSubCluster        Feature = "SubCluster"
	FeatureProductCluster    Feature = "ProductCluster"
	FeatureTraffic           Feature = "Traffic"
	FeatureCert              Feature = "Cert"
	FeatureActiveHealthCheck Feature = "ActiveHealthCheck"

	// auth
	FeatureProductUser Feature = "AuthProductUser"
	FeatureUser        Feature = "User"
	FeatureToken       Feature = "Token"

	// nlb resource
	FeatureNLBPool    Feature = "NLBPool"
	FeatureNLBCluster Feature = "NLBCluster"
)

var (
	actionProductNormal = ActionDeny.
				Grant(ActionRead).
				Grant(ActionUpdate).
				Grant(ActionCreate).
				Grant(ActionReadAll).
				Grant(ActionDelete)

	actionAll = ActionDeny.
			Grant(ActionRead).
			Grant(ActionUpdate).
			Grant(ActionCreate).
			Grant(ActionDelete).
			Grant(ActionReadAll).
			Grant(ActionExport)
)

var scope2permission = map[string]map[Feature]Action{
	ScopeSystem: {
		FeatureProxyPool:  actionAll,
		FeatureBFECluster: actionAll,
		FeatureBFEPool:    actionAll,
		FeatureArea:       actionAll,
		FeatureDomain:     actionAll,
		FeatureProduct:    actionAll,

		FeatureProductPool:       actionAll,
		FeatureRoute:             actionAll,
		FeatureSubCluster:        actionAll,
		FeatureProductCluster:    actionAll,
		FeatureTraffic:           actionAll,
		FeatureCert:              actionAll,
		FeatureActiveHealthCheck: actionAll,

		FeatureProductUser: actionAll,
		FeatureUser:        actionAll,
		FeatureToken:       actionAll,

		FeatureNLBPool:    actionAll,
		FeatureNLBCluster: actionAll,
	},
	ScopeProduct: {
		FeatureUser:       ActionReadAll,
		FeatureToken:      ActionReadAll,
		FeatureProxyPool:  actionProductNormal.Grant(ActionReadAll),
		FeatureBFECluster: actionProductNormal,
		FeatureBFEPool:    actionProductNormal,
		FeatureArea:       actionProductNormal.Grant(ActionReadAll),
		FeatureDomain:     actionProductNormal.Grant(ActionReadAll),
		FeatureProduct:    actionProductNormal.Grant(ActionReadAll),

		FeatureProductPool:       actionProductNormal,
		FeatureRoute:             actionProductNormal,
		FeatureSubCluster:        actionProductNormal,
		FeatureProductCluster:    actionProductNormal,
		FeatureTraffic:           actionProductNormal,
		FeatureCert:              actionProductNormal,
		FeatureActiveHealthCheck: actionProductNormal,

		FeatureProductUser: actionProductNormal,

		FeatureNLBPool:    actionProductNormal,
		FeatureNLBCluster: actionProductNormal,
	},
	ScopeSupport: {
		FeatureProxyPool:         ActionExport,
		FeatureRoute:             ActionExport,
		FeatureCert:              ActionExport,
		FeatureActiveHealthCheck: ActionExport,
		FeatureExtraFile:         ActionExport,
	},
}
