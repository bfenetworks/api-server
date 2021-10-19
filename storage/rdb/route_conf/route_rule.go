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

package route_conf

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/model/iroute_conf"
	"github.com/bfenetworks/api-server/model/iversion_control"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao"
)

var _ iroute_conf.RouteRuleStorager = &RouteRuleStorager{}

func NewRouteRuleStorager(dbCtxFactory lib.DBContextFactory,
	versionControlStorager iversion_control.VersionControlStorager) *RouteRuleStorager {
	return &RouteRuleStorager{
		dbCtxFactory:           dbCtxFactory,
		versionControlStorager: versionControlStorager,
	}
}

type RouteRuleStorager struct {
	dbCtxFactory lib.DBContextFactory

	versionControlStorager iversion_control.VersionControlStorager
}

func newDaoRouteBasicRuleParam(product *ibasic.Product, rule *iroute_conf.BasicRouteRule) (*dao.TRouteBasicRuleParam, error) {
	bsHostNames, err := json.Marshal(rule.HostNames)
	if err != nil {
		return nil, xerror.WrapDirtyDataError(err)
	}
	bsPaths, err := json.Marshal(rule.Paths)
	if err != nil {
		return nil, xerror.WrapDirtyDataError(err)
	}

	return &dao.TRouteBasicRuleParam{
		Description: lib.PString(rule.Description),
		ProductID:   lib.PInt64(product.ID),
		HostNames:   bsHostNames,
		ClusterID:   &rule.ClusterID,
		Paths:       bsPaths,
	}, nil
}

func (rs *RouteRuleStorager) FetchProductRule(ctx context.Context, product *ibasic.Product,
	clusterList []*icluster_conf.Cluster) (*iroute_conf.ProductRouteRule, error) {
	m, err := rs.FetchRoutRules(ctx, []*ibasic.Product{product}, clusterList)
	if err != nil {
		return nil, err
	}

	return m[product.ID], nil
}

func (rs *RouteRuleStorager) FetchRoutRules(ctx context.Context, products []*ibasic.Product,
	clusterList []*icluster_conf.Cluster) (map[int64]*iroute_conf.ProductRouteRule, error) {

	var productIDs []int64
	productID2Name := map[int64]string{}
	for _, one := range products {
		productIDs = append(productIDs, one.ID)
		productID2Name[one.ID] = one.Name
	}

	dbCtx, err := rs.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	advanceRules, err := dao.TRouteAdvanceRuleList(dbCtx, &dao.TRouteAdvanceRuleParam{
		ProductIDs: productIDs,
	})
	if err != nil {
		return nil, err
	}

	basicRules, err := dao.TRouteBasicRuleList(dbCtx, &dao.TRouteBasicRuleParam{
		ProductIDs: productIDs,
	})
	if err != nil {
		return nil, err
	}

	clusterMap := icluster_conf.ClusterList2MapByID(clusterList)

	product2ProductRouteRule := map[int64]*iroute_conf.ProductRouteRule{}
	for _, one := range advanceRules {
		rule := product2ProductRouteRule[one.ProductID]
		if rule == nil {
			rule = &iroute_conf.ProductRouteRule{}
			product2ProductRouteRule[one.ProductID] = rule
		}

		cluster := clusterMap[one.ClusterID]
		if cluster == nil {
			continue
		}
		rule.AdvanceRouteRules = append(rule.AdvanceRouteRules, &iroute_conf.AdvanceRouteRule{
			Name:        one.Name,
			Description: one.Description,
			Expression:  one.Expression,
			ClusterName: cluster.Name,
			ClusterID:   one.ClusterID,
		})
	}

	for _, one := range basicRules {
		rule := product2ProductRouteRule[one.ProductID]
		if rule == nil {
			rule = &iroute_conf.ProductRouteRule{}
			product2ProductRouteRule[one.ProductID] = rule
		}

		cluster := clusterMap[one.ClusterID]
		if cluster == nil {
			return nil, xerror.WrapDirtyDataErrorWithMsg("Cluster %d Not Existed", one.ClusterID)
		}
		basicRule, err := newRouteBasicRule(one, cluster)
		if err != nil {
			return nil, err
		}
		rule.BasicRouteRules = append(rule.BasicRouteRules, basicRule)
	}

	return product2ProductRouteRule, nil
}

func newRouteBasicRule(param *dao.TRouteBasicRule, cluster *icluster_conf.Cluster) (*iroute_conf.BasicRouteRule, error) {
	rule := &iroute_conf.BasicRouteRule{
		Description: param.Description,
		ClusterID:   param.ClusterID,
		ClusterName: cluster.Name,
	}
	if err := json.Unmarshal(param.HostNames, &rule.HostNames); err != nil {
		return nil, xerror.WrapDirtyDataErrorWithMsg("HostNames: %s, err: %v", string(param.HostNames), err)
	}
	if err := json.Unmarshal(param.Paths, &rule.Paths); err != nil {
		return nil, xerror.WrapDirtyDataErrorWithMsg("Paths: %s, err: %v", string(param.Paths), err)
	}

	return rule, nil
}

func (rs *RouteRuleStorager) UpsertProductRule(ctx context.Context, product *ibasic.Product,
	rule *iroute_conf.ProductRouteRule) error {

	now := lib.PTime(time.Now())

	// prepare db data
	daoBasicRules := []*dao.TRouteBasicRuleParam{}
	for _, one := range rule.BasicRouteRules {
		data, err := newDaoRouteBasicRuleParam(product, one)
		if err != nil {
			return err
		}
		data.CreatedAt = now
		data.UpdatedAt = now
		daoBasicRules = append(daoBasicRules, data)
	}

	daoAdvanceRules := []*dao.TRouteAdvanceRuleParam{}
	for _, one := range rule.AdvanceRouteRules {
		daoAdvanceRules = append(daoAdvanceRules, &dao.TRouteAdvanceRuleParam{
			ProductID:   &product.ID,
			Name:        lib.PString(one.Name),
			ClusterID:   lib.PInt64(one.ClusterID),
			Expression:  lib.PString(one.Expression),
			Description: lib.PString(one.Description),
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	dbCtx, err := rs.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	if _, err := dao.TRouteAdvanceRuleList(dbCtx, &dao.TRouteAdvanceRuleParam{
		ProductID: &product.ID,
		LockMode:  &dao.ModeForUpdate,
	}); err != nil {
		return err
	}

	if _, err := dao.TRouteAdvanceRuleDelete(dbCtx, &dao.TRouteAdvanceRuleParam{
		ProductID: &product.ID,
	}); err != nil {
		return err
	}

	if _, err := dao.TRouteBasicRuleDelete(dbCtx, &dao.TRouteBasicRuleParam{
		ProductID: &product.ID,
	}); err != nil {
		return err
	}

	if len(daoAdvanceRules) > 0 {
		if _, err := dao.TRouteAdvanceRuleCreate(dbCtx, daoAdvanceRules...); err != nil {
			return err
		}
	}

	if len(daoBasicRules) > 0 {
		if _, err := dao.TRouteBasicRuleCreate(dbCtx, daoBasicRules...); err != nil {
			return err
		}
	}

	return nil
}
