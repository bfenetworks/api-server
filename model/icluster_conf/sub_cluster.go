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

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/itxn"
)

type SubCluster struct {
	ID   int64
	Name string

	ClusterID int64

	ProductID   int64
	ProductName string

	InstancePool *Pool

	Capacity    int64
	Enabled     bool
	Ready       bool
	Description string
}

type SubClusterFilter struct {
	ID    *int64
	Name  *string
	Names []string

	InstancePool *Pool
	PoolIDs      []int64

	Product *ibasic.Product

	ClusterIDs []int64
}

type SubClusterParam struct {
	ID   *int64
	Name *string

	PoolID       *int64
	PoolName     *string
	InstancePool *Pool

	Product *ibasic.Product

	ClusterName *string
	Cluster     *Cluster
	ClusterIDs  []int64

	Description *string
}

type SubClusterStorager interface {
	FetchSubClusterList(ctx context.Context, param *SubClusterFilter) ([]*SubCluster, error)
	CreateSubCluster(ctx context.Context, param *SubClusterParam) error
	DeleteSubCluster(ctx context.Context, param *SubCluster) error
	UpdateSubCluster(ctx context.Context, one *SubCluster, param *SubClusterParam) error
}

type SubClusterManager struct {
	txn      itxn.TxnStorager
	storager SubClusterStorager

	productStorager ibasic.ProductStorager
	poolStorage     PoolStorage
	clusterStorager ClusterStorager
}

func NewSubClusterManager(txn itxn.TxnStorager, storager SubClusterStorager,
	productStorager ibasic.ProductStorager, poolStorage PoolStorage,
	clusterStorager ClusterStorager) *SubClusterManager {
	return &SubClusterManager{
		txn:             txn,
		storager:        storager,
		productStorager: productStorager,
		poolStorage:     poolStorage,
		clusterStorager: clusterStorager,
	}
}

func SubClusterList2MapByName(list []*SubCluster) map[string]*SubCluster {
	m := map[string]*SubCluster{}
	for _, one := range list {
		m[one.Name] = one
	}

	return m
}

func SubClusterList2MapByID(list []*SubCluster) map[int64]*SubCluster {
	m := map[int64]*SubCluster{}
	for _, one := range list {
		m[one.ID] = one
	}

	return m
}

func SubClusterList2IDSlice(list []*SubCluster) []int64 {
	var s []int64
	for _, one := range list {
		s = append(s, one.ID)
	}

	return s
}

func SubClusterList2NameSlice(list []*SubCluster) []string {
	var s []string
	for _, one := range list {
		s = append(s, one.Name)
	}

	return s
}

func (scm *SubClusterManager) SubClusterList(ctx context.Context, param *SubClusterFilter) (list []*SubCluster, err error) {

	err = scm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err = scm.storager.FetchSubClusterList(ctx, param)
		if err != nil {
			return err
		}
		if len(list) == 0 {
			return nil
		}

		clusterIDs := map[int64]bool{}
		for _, one := range list {
			clusterIDs[one.ClusterID] = true
		}

		return err
	})

	return
}

func (scm *SubClusterManager) FetchSubCluster(ctx context.Context, param *SubClusterFilter) (subCluster *SubCluster, err error) {
	list, err := scm.storager.FetchSubClusterList(ctx, param)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}

	return list[0], nil
}

func (scm *SubClusterManager) CreateSubCluster(ctx context.Context, product *ibasic.Product, param *SubClusterParam) (err error) {
	err = scm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err := scm.storager.FetchSubClusterList(ctx, &SubClusterFilter{
			Name:    param.Name,
			Product: product,
		})
		if err != nil {
			return err
		}
		if len(list) != 0 {
			return xerror.WrapRecordExisted("SubCluster")
		}

		pool, err := scm.poolStorage.FetchPool(ctx, *param.PoolName)
		if err != nil {
			return err
		}
		if pool == nil {
			return xerror.WrapParamErrorWithMsg("Pool Not Exist")
		}

		if pool.Product != nil && pool.Product.ID != product.ID && pool.Product.ID != 1 {
			return xerror.WrapParamErrorWithMsg("Pool Not Valid")
		}

		param.InstancePool = pool
		param.Cluster = &Cluster{
			ID:   -1,
			Name: "unbinding",
		}
		param.Product = product
		return scm.storager.CreateSubCluster(ctx, param)
	})

	return
}

func (scm *SubClusterManager) DeleteSubCluster(ctx context.Context, subCluster *SubCluster) (err error) {
	if subCluster.ClusterID > 0 { // be mounted to cluster
		return xerror.WrapModelErrorWithMsg("SubCluster %s be Mounted With Cluster %d", subCluster.Name, subCluster.ClusterID)
	}
	err = scm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		return scm.storager.DeleteSubCluster(ctx, subCluster)
	})

	return
}

func (scm *SubClusterManager) UpdateSubCluster(ctx context.Context, subCluster *SubCluster, param *SubClusterParam) (err error) {
	err = scm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		return scm.storager.UpdateSubCluster(ctx, subCluster, param)
	})

	return
}
