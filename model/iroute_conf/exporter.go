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

	"github.com/bfenetworks/bfe/bfe_config/bfe_cluster_conf/cluster_conf"
	"github.com/bfenetworks/bfe/bfe_config/bfe_route_conf/host_rule_conf"
	"github.com/bfenetworks/bfe/bfe_config/bfe_route_conf/route_rule_conf"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/model/iversion_control"
)

type RouteRuleExportData struct {
	Version     string
	HostTable   *host_rule_conf.HostTableConf
	RouteTable  *route_rule_conf.RouteTableFile
	ClusterConf *cluster_conf.BfeClusterConf
}

func (rred *RouteRuleExportData) UpdateVersion(version string) error {
	rred.Version = version
	rred.RouteTable.Version = &version
	rred.HostTable.Version = &version
	rred.ClusterConf.Version = &version

	return nil
}

const (
	ConfigTopicRouteRule = "route_rule"
)

func (rm *RouteRuleManager) ExportRouteRule(ctx context.Context, lastVersion string) (*RouteRuleExportData, error) {
	ed, err := rm.versionControlManager.ExportConfig(ctx, ConfigTopicRouteRule, rm.exportRouteRule)
	if err != nil {
		return nil, err
	}

	conf := ed.DataWithoutVersion.(*RouteRuleExportData)
	if conf.Version == lastVersion {
		return nil, nil
	}

	return conf, nil
}

func (rm *RouteRuleManager) exportRouteRule(ctx context.Context) (*iversion_control.ExportData, error) {
	domains, err := rm.domainStorager.FetchDomains(ctx, nil)
	if err != nil {
		return nil, err
	}

	clusters, err := rm.clusterStorager.FetchClusterList(ctx, nil)
	if err != nil {
		return nil, err
	}
	clusters = icluster_conf.AppendAdvancedRuleCluster(clusters)

	routeRules, err := rm.storager.FetchRoutRules(ctx, nil, clusters)
	if err != nil {
		return nil, err
	}

	products, err := rm.productStorager.FetchProducts(ctx, nil)
	if err != nil {
		return nil, err
	}

	productMapID2Name := map[int64]string{}
	tmp := ibasic.ProductIDMap(products)
	for _, domain := range domains {
		product, ok := tmp[domain.ProductID]
		if !ok {
			return nil, xerror.WrapDirtyDataErrorWithMsg("Domain refer Not Existed Product %d", domain.ProductID)
		}
		productMapID2Name[domain.ProductID] = product.Name
	}

	emptyVersion := iversion_control.ZeroVersion
	rred := &RouteRuleExportData{
		Version:     emptyVersion,
		RouteTable:  newRouteTableFile(emptyVersion, productMapID2Name, routeRules),
		HostTable:   newHostTableConf(emptyVersion, productMapID2Name, domains),
		ClusterConf: icluster_conf.NewBfeClusterConf(emptyVersion, clusters),
	}

	return &iversion_control.ExportData{
		Topic:              ConfigTopicRouteRule,
		DataWithoutVersion: rred,
	}, nil
}

var defaultProduct = "bfe"

func newRouteTableFile(version string, productMapID2Name map[int64]string,
	routeRules map[int64]*ProductRouteRule) *route_rule_conf.RouteTableFile {

	basicRule := route_rule_conf.ProductBasicRouteRuleFile{}
	advanceRule := route_rule_conf.ProductAdvancedRouteRuleFile{}

	clusterName := func(cn string) *string {
		if cn == icluster_conf.RouteAdvancedModeClusterName4DP {
			return &cn
		}

		tmp := "cluster_" + cn
		return &tmp
	}

	newBasicRouteRuleFiles := func(brrs []*BasicRouteRule) (bs []route_rule_conf.BasicRouteRuleFile) {
		for _, brr := range brrs {
			bs = append(bs, route_rule_conf.BasicRouteRuleFile{
				Hostname:    brr.HostNames,
				Path:        brr.Paths,
				ClusterName: clusterName(brr.ClusterName),
			})
		}

		return bs
	}

	newAdvanceRouteRuleFiles := func(arrs []*AdvanceRouteRule) (as []route_rule_conf.AdvancedRouteRuleFile) {
		for _, arr := range arrs {
			as = append(as, route_rule_conf.AdvancedRouteRuleFile{
				Cond:        &arr.Expression,
				ClusterName: clusterName(arr.ClusterName),
			})
		}

		return as
	}

	// sort by product_id
	if len(routeRules) > 0 {
		maxProductID := len(routeRules)
		for pid := range routeRules {
			maxProductID += int(pid)
			break
		}

		for _, pid := range lib.SortMapInt642String(productMapID2Name) {
			rule, ok := routeRules[pid]
			if !ok {
				continue
			}

			pName, ok := productMapID2Name[pid]
			if !ok {
				continue
			}

			basicRule[pName] = newBasicRouteRuleFiles(rule.BasicRouteRules)
			advanceRule[pName] = newAdvanceRouteRuleFiles(rule.AdvanceRouteRules)
		}
	}

	return &route_rule_conf.RouteTableFile{
		Version:     &version,
		BasicRule:   &basicRule,
		ProductRule: &advanceRule,
	}
}
