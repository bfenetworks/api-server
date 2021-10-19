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

import (
	"context"
)

type Feature string

type Action int64

type FeatureAuthorizer struct {
	Feature Feature
	Action  Action
}

func NewFeatureAuthorizer(f Feature, a Action) *FeatureAuthorizer {
	return &FeatureAuthorizer{
		Feature: f,
		Action:  a,
	}
}

var DefaultProductAuthorizateManager *ProductAuthorizateManager

func NewFeatureAuthorizerWithFactoryWithProduct(f Feature, a Action) *MultiAuthorizer {
	return NewMultiAuthorizerWithFactory([]func() Authorizer{
		func() Authorizer { return NewFeatureAuthorizer(f, a) },
		func() Authorizer { return DefaultProductAuthorizateManager },
	}, OpAnd)
}

var (
	FA  = NewFeatureAuthorizer
	FAP = NewFeatureAuthorizerWithFactoryWithProduct
)

func (pa *FeatureAuthorizer) Authorizate(ctx context.Context, user *User) (bool, error) {
	for _, role := range user.Roles {
		feature2action := role2permission[role.Name]
		if feature2action == nil {
			continue
		}

		action := feature2action[pa.Feature]
		if action.IsAllowed(pa.Action) {
			return true, nil
		}
	}

	return false, nil
}

func (pa *FeatureAuthorizer) Grant(ctx context.Context, user *User) error {
	panic("implemented")
}

func (pa *FeatureAuthorizer) Revoke(ctx context.Context, user *User) error {
	panic("implemented")
}

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

	// product resource, module
	FeatureModHeader  Feature = "mod.header"
	FeatureModRewrite Feature = "mod.rewrite"

	// auth
	FeatureProductUser Feature = "AuthProductUser"
	FeatureUser        Feature = "User"

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

var role2permission = map[string]map[Feature]Action{
	"admin": {
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

		FeatureModHeader:  actionAll,
		FeatureModRewrite: actionAll,

		FeatureProductUser: actionAll,
		FeatureUser:        actionAll,

		FeatureNLBPool:    actionAll,
		FeatureNLBCluster: actionAll,
	},
	"product": {
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

		FeatureModHeader:  actionProductNormal,
		FeatureModRewrite: actionProductNormal,

		FeatureProductUser: actionProductNormal,
		FeatureUser:        actionProductNormal,

		FeatureNLBPool:    actionProductNormal,
		FeatureNLBCluster: actionProductNormal,
	},
	"inner": {
		FeatureProxyPool:         ActionExport,
		FeatureRoute:             ActionExport,
		FeatureCert:              ActionExport,
		FeatureActiveHealthCheck: ActionExport,
		FeatureModHeader:         ActionExport,
		FeatureModRewrite:        ActionExport,
		FeatureExtraFile:         ActionExport,
	},
}
