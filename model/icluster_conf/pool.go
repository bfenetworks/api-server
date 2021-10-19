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
	"fmt"
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
	Instances []Instance

	Tag *int8
}

type Pool struct {
	ID        int64
	Name      string
	Ready     bool
	Product   *ibasic.Product
	Instances []Instance
	Tag       int8
}

type Instance struct {
	HostName string            `json:"Name"`
	IP       string            `json:"Addr"`
	Port     int               `json:"Port"`
	Ports    map[string]int    `json:"Ports,omitempty"`
	Tags     map[string]string `json:"tags,omitempty"`
	Weight   int64             `json:"Weight"`
	Disable  bool              `json:"Disable"`
}

func (i *Instance) IPWithPort() string {
	if i.Port == 0 {
		i.Port = i.Ports["Default"]
	}

	return fmt.Sprintf("%s:%d", i.IP, i.Port)
}

type PoolStorager interface {
	FetchPool(ctx context.Context, name string) (*Pool, error)
	FetchPools(ctx context.Context, param *PoolFilter) ([]*Pool, error)

	CreatePool(ctx context.Context, product *ibasic.Product, data *PoolParam) (*Pool, error)
	UpdatePool(ctx context.Context, oldData *Pool, diff *PoolParam) error
	DeletePool(ctx context.Context, pool *Pool) error
}

type PoolManager struct {
	storager           PoolStorager
	bfeClusterStorager ibasic.BFEClusterStorager
	subClusterStorager SubClusterStorager
	txn                itxn.TxnStorager
}

func NewPoolManager(txn itxn.TxnStorager, storager PoolStorager,
	bfeClusterStorager ibasic.BFEClusterStorager, subClusterStorager SubClusterStorager) *PoolManager {

	return &PoolManager{
		txn:                txn,
		storager:           storager,
		bfeClusterStorager: bfeClusterStorager,
		subClusterStorager: subClusterStorager,
	}
}

func (rppm *PoolManager) FetchPoolByName(ctx context.Context, name string) (one *Pool, err error) {
	err = rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		one, err = rppm.storager.FetchPool(ctx, name)
		return err
	})

	return
}

func (rppm *PoolManager) FetchBFEPool(ctx context.Context, name string) (one *Pool, err error) {
	return rppm.FetchProductPool(ctx, ibasic.BuildinProduct, name)
}

func (rppm *PoolManager) FetchProductPool(ctx context.Context, product *ibasic.Product, name string) (one *Pool, err error) {
	name, err = poolNameJudger(product.Name, name)
	if err != nil {
		return
	}

	err = rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		one, err = rppm.storager.FetchPool(ctx, name)
		return err
	})

	return
}

func (rppm *PoolManager) FetchBFEPools(ctx context.Context) (list []*Pool, err error) {
	return rppm.FetchProductPools(ctx, ibasic.BuildinProduct)
}

func (rppm *PoolManager) FetchProductPools(ctx context.Context, product *ibasic.Product) (list []*Pool, err error) {
	err = rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err = rppm.storager.FetchPools(ctx, &PoolFilter{
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
func (rppm *PoolManager) CanDelete(ctx context.Context, pool *Pool) error {
	bfeClusters, err := rppm.bfeClusterStorager.FetchBFEClusters(ctx, &ibasic.BFEClusterFilter{
		Pool: &pool.Name,
	})
	if err != nil {
		return err
	}
	if len(bfeClusters) != 0 {
		return xerror.WrapModelErrorWithMsg("BFECluster %s Refer To This Pool", bfeClusters[0].Name)
	}

	subClusters, err := rppm.subClusterStorager.FetchSubClusterList(ctx, &SubClusterFilter{
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

func (rppm *PoolManager) DeleteBFEPool(ctx context.Context, name string) (one *Pool, err error) {
	return rppm.DeleteProductPool(ctx, ibasic.BuildinProduct, name)
}

func (rppm *PoolManager) DeleteProductPool(ctx context.Context, product *ibasic.Product, name string) (one *Pool, err error) {
	name, err = poolNameJudger(product.Name, name)
	if err != nil {
		return
	}

	err = rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		one, err = rppm.storager.FetchPool(ctx, name)
		if err != nil {
			return err
		}

		if one == nil {
			return xerror.WrapRecordNotExist("Pool")
		}

		if err = rppm.CanDelete(ctx, one); err != nil {
			return err
		}

		return rppm.storager.DeletePool(ctx, one)
	})

	return
}

func (rppm *PoolManager) CreateBFEPool(ctx context.Context, pool *PoolParam) (one *Pool, err error) {
	pool.Tag = &PoolTagBFE
	return rppm.CreateProductPool(ctx, ibasic.BuildinProduct, pool)
}

func (rppm *PoolManager) CreateProductPool(ctx context.Context, product *ibasic.Product, pool *PoolParam) (one *Pool, err error) {
	var pN string
	pN, err = poolNameJudger(product.Name, *pool.Name)
	if err != nil {
		return
	}
	pool.Name = &pN
	if pool.Tag == nil {
		pool.Tag = &PoolTagProduct
	}

	err = rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		old, err := rppm.storager.FetchPool(ctx, *pool.Name)
		if err != nil {
			return err
		}
		if old != nil {
			return xerror.WrapRecordExisted()
		}

		one, err = rppm.storager.CreatePool(ctx, product, pool)
		return err
	})

	return
}

func (rppm *PoolManager) UpdateBFEPool(ctx context.Context, pool *Pool, diff *PoolParam) (err error) {
	return rppm.UpdateProductPool(ctx, ibasic.BuildinProduct, pool, diff)
}

func (rppm *PoolManager) UpdateProductPool(ctx context.Context, product *ibasic.Product, pool *Pool, diff *PoolParam) (err error) {
	err = rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		return rppm.storager.UpdatePool(ctx, pool, diff)
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

func (rppm *PoolManager) GetPoolByName(ctx context.Context, poolName *string) (pool *Pool, err error) {
	err = rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		if poolName == nil || *poolName == "" {
			return xerror.WrapParamErrorWithMsg("Pool Name Illegal")
		}

		pool, err = rppm.storager.FetchPool(ctx, *poolName)
		return err
	})

	return
}
