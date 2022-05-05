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

package cluster_conf

import (
	"context"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao"
)

type RDBSubClusterStorager struct {
	dbCtxFactory lib.DBContextFactory

	poolStorage     icluster_conf.PoolStorage
	productStorager ibasic.ProductStorager
}

func NewRDBSubClusterStorager(dbCtxFactory lib.DBContextFactory, poolStorage icluster_conf.PoolStorage,
	productStorager ibasic.ProductStorager) *RDBSubClusterStorager {

	return &RDBSubClusterStorager{
		dbCtxFactory:    dbCtxFactory,
		poolStorage:     poolStorage,
		productStorager: productStorager,
	}
}

var _ icluster_conf.SubClusterStorager = &RDBSubClusterStorager{}

func subClusterFilter2Param(filter *icluster_conf.SubClusterFilter) *dao.TSubClusterParam {
	if filter == nil {
		return nil
	}

	tmp := &dao.TSubClusterParam{
		PoolsIDs: filter.PoolIDs,

		Name:  filter.Name,
		Names: filter.Names,

		ClusterIDs: filter.ClusterIDs,
	}

	if filter.Product != nil {
		tmp.ProductID = &filter.Product.ID
	}

	if filter.InstancePool != nil {
		tmp.PoolsID = &filter.InstancePool.ID
	}

	return tmp
}

func subClusterParami2d(data *icluster_conf.SubClusterParam) *dao.TSubClusterParam {
	if data == nil {
		return nil
	}

	tmp := &dao.TSubClusterParam{
		PoolsID: data.PoolID,

		Name: data.Name,

		ClusterIDs:  data.ClusterIDs,
		Description: data.Description,
	}

	if data.Product != nil {
		tmp.ProductID = &data.Product.ID
	}

	if data.Cluster != nil {
		tmp.ClusterID = &data.Cluster.ID
	}

	if data.InstancePool != nil {
		tmp.PoolsID = &data.InstancePool.ID
	}

	return tmp
}

func newSubCluster(pp *dao.TSubCluster, pool *icluster_conf.Pool, product *ibasic.Product) *icluster_conf.SubCluster {
	data := &icluster_conf.SubCluster{
		ID:          pp.ID,
		Name:        pp.Name,
		Enabled:     pp.Enabled,
		Description: pp.Description,
		ClusterID:   pp.ClusterID,

		InstancePool: pool,
	}

	if pool != nil {
		data.Ready = pool.Ready
	}

	return data
}

func (rpps *RDBSubClusterStorager) FetchSubClusterList(ctx context.Context,
	filter *icluster_conf.SubClusterFilter) ([]*icluster_conf.SubCluster, error) {

	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	subClusterList, err := dao.TSubClusterList(dbCtx, subClusterFilter2Param(filter))
	if err != nil {
		return nil, err
	}

	if len(subClusterList) == 0 {
		return nil, nil
	}

	poolIDs, productIDs := map[int64]bool{}, map[int64]bool{}
	for _, one := range subClusterList {
		poolIDs[one.PoolsID] = true
		productIDs[one.ProductID] = true
	}
	poolList, err := rpps.poolStorage.FetchPools(dbCtx, &icluster_conf.PoolFilter{
		IDs: lib.Int64BoolMap2Slice(poolIDs),
	})
	if err != nil {
		return nil, err
	}

	productList, err := rpps.productStorager.FetchProducts(dbCtx, &ibasic.ProductFilter{
		IDs: lib.Int64BoolMap2Slice(productIDs),
	})
	if err != nil {
		return nil, err
	}

	poolMap := icluster_conf.PoolList2Map(poolList)
	productMap := ibasic.ProductIDMap(productList)

	rst := []*icluster_conf.SubCluster{}
	for _, one := range subClusterList {
		rst = append(rst, newSubCluster(one, poolMap[one.PoolsID], productMap[one.ProductID]))
	}

	return rst, nil
}

func (rpps *RDBSubClusterStorager) CreateSubCluster(ctx context.Context, param *icluster_conf.SubClusterParam) error {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TSubClusterCreate(dbCtx, subClusterParami2d(param))
	if err != nil {
		return err
	}

	return nil
}

func (rpps *RDBSubClusterStorager) DeleteSubCluster(ctx context.Context, oldOne *icluster_conf.SubCluster) error {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TSubClusterDelete(dbCtx, &dao.TSubClusterParam{ID: &oldOne.ID})

	return err
}

func (rpps *RDBSubClusterStorager) UpdateSubCluster(ctx context.Context, oldOne *icluster_conf.SubCluster, param *icluster_conf.SubClusterParam) error {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TSubClusterUpdate(dbCtx, subClusterParami2d(param), &dao.TSubClusterParam{ID: &oldOne.ID})

	return err
}
