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

package iroute_conf

import (
	"context"
	"fmt"
	"strings"

	"github.com/bfenetworks/bfe/bfe_basic/condition"
	"github.com/bfenetworks/bfe/bfe_config/bfe_route_conf/route_rule_conf"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/model/itxn"
	"github.com/bfenetworks/api-server/model/iversion_control"
	"github.com/bfenetworks/api-server/stateful"
)

type BasicRouteRule struct {
	HostNames   []string
	Paths       []string
	ClusterName string
	ClusterID   int64
	Description string
}

type AdvanceRouteRule struct {
	Name        string
	Description string
	Expression  string
	ClusterName string
	ClusterID   int64
}

type RouteRuleCase struct {
	Description   string
	URL           string
	Method        string
	Header        map[string]string
	ExpectCluster string
}

type ProductRouteRule struct {
	BasicRouteRules   []*BasicRouteRule
	AdvanceRouteRules []*AdvanceRouteRule

	RouteCases []*RouteRuleCase
}

type HostUsedInfo struct {
	Type   string
	Detail string
}

func (prr *ProductRouteRule) HostBeUsed(host string) *HostUsedInfo {
	for _, brr := range prr.BasicRouteRules {
		for _, h := range brr.HostNames {
			if host == h {
				return &HostUsedInfo{
					Type:   "BasicConditionExpression",
					Detail: host,
				}
			}
		}
	}

	keyword := fmt.Sprintf(`req_host_in("%s")`, host)

	for _, arr := range prr.AdvanceRouteRules {
		if strings.Contains(arr.Expression, keyword) {
			return &HostUsedInfo{
				Type:   "AdvanceConditionExpression",
				Detail: arr.Expression,
			}
		}
	}

	return nil
}

type ProductRouteRuleConvertResult struct {
	BasicRouteRuleFiles    []route_rule_conf.BasicRouteRuleFile
	AdvancedRouteRuleFiles []route_rule_conf.AdvancedRouteRuleFile

	ReferClusterNames []string
}

type RouteRuleRunCaseResult struct {
	RouteRuleCase
	ActualCluster string
	Pass          bool
}

func (pfr *ProductRouteRule) Convert() (*ProductRouteRuleConvertResult, error) {
	if len(pfr.AdvanceRouteRules) == 0 ||
		pfr.AdvanceRouteRules[len(pfr.AdvanceRouteRules)-1].Expression != "default_t()" {
		return nil, xerror.WrapParamErrorWithMsg("Last ForwardRule Expression Must Be default_t()")
	}

	clusterNameMap := map[string]bool{}
	basicRules := route_rule_conf.BasicRouteRuleFiles{}
	for _, one := range pfr.BasicRouteRules {
		clusterNameMap[one.ClusterName] = true
		basicRules = append(basicRules, route_rule_conf.BasicRouteRuleFile{
			Hostname:    one.HostNames,
			Path:        one.Paths,
			ClusterName: &one.ClusterName,
		})
	}

	advanceRules := route_rule_conf.AdvancedRouteRuleFiles{}
	for _, one := range pfr.AdvanceRouteRules {
		clusterNameMap[one.ClusterName] = true
		advanceRules = append(advanceRules, route_rule_conf.AdvancedRouteRuleFile{
			Cond:        lib.PString(one.Expression),
			ClusterName: &one.ClusterName,
		})
	}

	return &ProductRouteRuleConvertResult{
		BasicRouteRuleFiles:    basicRules,
		AdvancedRouteRuleFiles: advanceRules,

		ReferClusterNames: lib.StringMap2Slice(clusterNameMap),
	}, nil
}

type RouteRuleStorager interface {
	UpsertProductRule(ctx context.Context, product *ibasic.Product, rule *ProductRouteRule) error
	FetchProductRule(ctx context.Context, product *ibasic.Product,
		clusterList []*icluster_conf.Cluster) (*ProductRouteRule, error)
	FetchRoutRules(ctx context.Context, products []*ibasic.Product,
		clusterList []*icluster_conf.Cluster) (map[int64]*ProductRouteRule, error)
}

func NewRouteRuleManager(txn itxn.TxnStorager, storager RouteRuleStorager, clusterStorager icluster_conf.ClusterStorager,
	productStorager ibasic.ProductStorager, versionControlManager *iversion_control.VersionControlManager,
	domainStorager DomainStorager) *RouteRuleManager {
	return &RouteRuleManager{
		txn:             txn,
		storager:        storager,
		clusterStorager: clusterStorager,
		productStorager: productStorager,

		versionControlManager: versionControlManager,
		domainStorager:        domainStorager,
	}
}

type RouteRuleManager struct {
	versionControlManager *iversion_control.VersionControlManager

	txn             itxn.TxnStorager
	storager        RouteRuleStorager
	clusterStorager icluster_conf.ClusterStorager
	productStorager ibasic.ProductStorager
	domainStorager  DomainStorager
}

func (rm *RouteRuleManager) ExpressionVerify(ctx context.Context, expression string) (err error) {
	_, err = condition.Build(expression)
	return err
}

func (rm *RouteRuleManager) FetchProductRule(ctx context.Context, product *ibasic.Product) (prr *ProductRouteRule, err error) {
	err = rm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		clusters, err := rm.clusterStorager.FetchClusterList(ctx, &icluster_conf.ClusterFilter{
			Product: product,
		})
		if err != nil {
			return err
		}

		clusters = icluster_conf.AppendAdvancedRuleCluster(clusters)

		m, err := rm.storager.FetchRoutRules(ctx, []*ibasic.Product{product}, clusters)
		if err != nil {
			return err
		}
		prr = m[product.ID]

		return nil
	})

	return
}

func (rm *RouteRuleManager) UpsertProductRule(ctx context.Context, product *ibasic.Product, rule *ProductRouteRule) error {
	cr, err := rule.Convert()
	if err != nil {
		return err
	}

	var clusterList []*icluster_conf.Cluster
	var clusterMap map[string]*icluster_conf.Cluster

	err = rm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		// verify cluster
		if referClusters := cr.ReferClusterNames; len(referClusters) > 0 {
			clusterList, err = rm.clusterStorager.FetchClusterList(ctx, &icluster_conf.ClusterFilter{
				Names:   referClusters,
				Product: product,
			})
			if err != nil {
				return err
			}

			clusterList = icluster_conf.AppendAdvancedRuleCluster(clusterList)
			clusterMap = icluster_conf.ClusterList2MapByName(clusterList)
			for _, clusterName := range referClusters {
				if icluster_conf.SystemKeepRouteNames[clusterName] {
					continue
				}

				cluster, ok := clusterMap[clusterName]
				if !ok {
					return xerror.WrapModelErrorWithMsg("Cluster %s Not Exist", clusterName)
				}

				if !stateful.IgnoreBNSStatusCheck {
					if !cluster.Ready {
						return xerror.WrapModelErrorWithMsg("Cluster %s Not Ready", clusterName)
					}
				}
			}
		}

		for _, one := range rule.AdvanceRouteRules {
			one.ClusterID = clusterMap[one.ClusterName].ID
		}
		for _, one := range rule.BasicRouteRules {
			one.ClusterID = clusterMap[one.ClusterName].ID
		}

		if err := rm.storager.UpsertProductRule(ctx, product, rule); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (rm *RouteRuleManager) ClusterDeleteChecker(ctx context.Context, product *ibasic.Product, cluster *icluster_conf.Cluster) error {
	m, err := rm.storager.FetchRoutRules(ctx, []*ibasic.Product{product}, []*icluster_conf.Cluster{cluster})
	if err != nil {
		return err
	}

	if m == nil || m[product.ID] == nil {
		return nil
	}

	rule := m[product.ID]
	if len(rule.AdvanceRouteRules) > 0 {
		return xerror.WrapModelErrorWithMsg("Rule %s Refer To This Cluster", rule.AdvanceRouteRules[0].Name)
	}

	if len(rule.BasicRouteRules) > 0 {
		return xerror.WrapModelErrorWithMsg("Rule %s Refer To This Cluster", rule.BasicRouteRules[0].Description)
	}

	return nil
}
