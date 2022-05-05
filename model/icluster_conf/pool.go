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

package icluster_conf

import (
	"context"
	"strings"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/itxn"
)

var (
	PoolTagBFE     int8 = 1
	PoolTagProduct int8 = 2
)

type PoolFilter struct {
	Name      *string
	IDs       []int64
	ID        *int64
	ProductID *int64
}

type PoolParam struct {
	ID        *int64
	Name      *string
	ProductID *int64
	Type      *int8
	Tag       *int8
}

type Pool struct {
	ID      int64
	Name    string
	Type    int8
	Ready   bool
	Tag     int8
	Product *ibasic.Product

	instances []Instance
}

func (p *Pool) SetDefaultInstances(is []Instance) {
	p.instances = is
}

func (p *Pool) GetDefaultPool() *InstancePool {
	return &InstancePool{
		Name:      p.Name,
		Instances: p.instances,
	}
}

type PoolStorage interface {
	FetchPool(ctx context.Context, name string) (*Pool, error)
	FetchPools(ctx context.Context, param *PoolFilter) ([]*Pool, error)

	CreatePool(ctx context.Context, product *ibasic.Product, data *PoolParam) (*Pool, error)
	UpdatePool(ctx context.Context, oldData *Pool, diff *PoolParam) error
	DeletePool(ctx context.Context, pool *Pool) error
}

type PoolManager struct {
	storage            PoolStorage
	bfeClusterStorager ibasic.BFEClusterStorager
	subClusterStorager SubClusterStorager
	txn                itxn.TxnStorager

	instancePoolManager *InstancePoolManager
}

func NewPoolManager(txn itxn.TxnStorager, storage PoolStorage,
	bfeClusterStorager ibasic.BFEClusterStorager, subClusterStorager SubClusterStorager,
	instancePoolManager *InstancePoolManager) *PoolManager {

	return &PoolManager{
		txn:                txn,
		storage:            storage,
		bfeClusterStorager: bfeClusterStorager,
		subClusterStorager: subClusterStorager,

		instancePoolManager: instancePoolManager,
	}
}

func (m *PoolManager) FetchPoolByName(ctx context.Context, name string) (one *Pool, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		one, err = m.storage.FetchPool(ctx, name)
		return err
	})

	return
}

func (m *PoolManager) FetchBFEPool(ctx context.Context, name string) (one *Pool, err error) {
	return m.FetchProductPool(ctx, ibasic.BuildinProduct, name)
}

func (m *PoolManager) FetchProductPool(ctx context.Context, product *ibasic.Product, name string) (one *Pool, err error) {
	name, err = poolNameJudger(product.Name, name)
	if err != nil {
		return
	}

	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		one, err = m.storage.FetchPool(ctx, name)
		return err
	})

	return
}

func (m *PoolManager) FetchBFEPools(ctx context.Context) (list []*Pool, err error) {
	return m.FetchProductPools(ctx, ibasic.BuildinProduct)
}

func (m *PoolManager) FetchProductPools(ctx context.Context, product *ibasic.Product) (list []*Pool, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err = m.storage.FetchPools(ctx, &PoolFilter{
			ProductID: &product.ID,
		})
		return err
	})

	return
}

func poolNameJudger(productName string, poolName string) (realName string, err error) {
	ss := strings.SplitN(poolName, ".", 2)
	if len(ss) == 2 {
		if ss[0] != productName {
			return "", xerror.WrapParamErrorWithMsg("Pool Name Must Use Product Name as Prefix")
		}

		return poolName, nil
	}

	return productName + "." + poolName, nil
}

// CanDelete check whether pool can be deleted, Check Logic:
// 1. Not BFE Cluster Refer To
// 2. Not SubCluster Refer To
func (m *PoolManager) CanDelete(ctx context.Context, pool *Pool) error {
	bfeClusters, err := m.bfeClusterStorager.FetchBFEClusters(ctx, &ibasic.BFEClusterFilter{
		Pool: &pool.Name,
	})
	if err != nil {
		return err
	}
	if len(bfeClusters) != 0 {
		return xerror.WrapModelErrorWithMsg("BFECluster %s Refer To This Pool", bfeClusters[0].Name)
	}

	subClusters, err := m.subClusterStorager.FetchSubClusterList(ctx, &SubClusterFilter{
		InstancePool: pool,
	})
	if err != nil {
		return err
	}
	if len(subClusters) != 0 {
		return xerror.WrapModelErrorWithMsg("SubCluster %s Refer To This Pool", subClusters[0].Name)
	}

	return nil
}

func (m *PoolManager) DeleteBFEPool(ctx context.Context, name string) (one *Pool, err error) {
	return m.DeleteProductPool(ctx, ibasic.BuildinProduct, name)
}

func (m *PoolManager) DeleteProductPool(ctx context.Context, product *ibasic.Product, name string) (one *Pool, err error) {
	name, err = poolNameJudger(product.Name, name)
	if err != nil {
		return
	}

	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		one, err = m.storage.FetchPool(ctx, name)
		if err != nil {
			return err
		}

		if one == nil {
			return xerror.WrapRecordNotExist("Pool")
		}

		if err = m.CanDelete(ctx, one); err != nil {
			return err
		}

		return m.storage.DeletePool(ctx, one)
	})

	return
}

func (m *PoolManager) CreateBFEPool(ctx context.Context, pool *PoolParam, pis *InstancePool) (one *Pool, err error) {
	pool.Tag = &PoolTagBFE
	return m.CreateProductPool(ctx, ibasic.BuildinProduct, pool, pis)
}

func (m *PoolManager) CreateProductPool(ctx context.Context, product *ibasic.Product, param *PoolParam, pool *InstancePool) (one *Pool, err error) {
	var pN string
	pN, err = poolNameJudger(product.Name, *param.Name)
	if err != nil {
		return
	}
	param.Name = &pN
	if param.Tag == nil {
		param.Tag = &PoolTagProduct
	}

	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		old, err := m.storage.FetchPool(ctx, *param.Name)
		if err != nil {
			return err
		}
		if old != nil {
			return xerror.WrapRecordExisted()
		}

		one, err = m.storage.CreatePool(ctx, product, param)
		if err != nil {
			return err
		}

		if pool != nil {
			err = m.instancePoolManager.UpdateInstances(ctx, one, pool)
		}
		return err
	})

	return
}

func (m *PoolManager) UpdateBFEPool(ctx context.Context, pool *Pool, diff *PoolParam) (err error) {
	return m.UpdateProductPool(ctx, ibasic.BuildinProduct, pool, diff)
}

func (m *PoolManager) UpdateProductPool(ctx context.Context, product *ibasic.Product, pool *Pool, diff *PoolParam) (err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		return m.storage.UpdatePool(ctx, pool, diff)
	})

	return
}

func PoolList2Map(list []*Pool) map[int64]*Pool {
	m := map[int64]*Pool{}
	for _, one := range list {
		m[one.ID] = one
	}

	return m
}

func (m *PoolManager) GetPoolByName(ctx context.Context, poolName *string) (pool *Pool, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		if poolName == nil || *poolName == "" {
			return xerror.WrapParamErrorWithMsg("Pool Name Illegal")
		}

		pool, err = m.storage.FetchPool(ctx, *poolName)
		return err
	})

	return
}
