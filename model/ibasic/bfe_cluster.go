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

package ibasic

import (
	"context"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/itxn"
)

type BFECluster struct {
	ID                 int64
	Name               string
	Pool               string
	Enabled            bool
	ExemptTrafficCheck bool
	Capacity           int64
}

type BFEClusterParam struct {
	Name     *string
	Pool     *string
	Capacity *int64
}

type BFEClusterFilter struct {
	Name *string
	Pool *string
}

type BFEClusterStorager interface {
	DeleteBFECluster(context.Context, *BFECluster) error
	CreateBFECluster(context.Context, *BFEClusterParam) error
	FetchBFEClusters(context.Context, *BFEClusterFilter) ([]*BFECluster, error)
}

type BFEClusterManager struct {
	storager BFEClusterStorager
	txn      itxn.TxnStorager
}

func NewBFEClusterManager(txn itxn.TxnStorager, storager BFEClusterStorager) *BFEClusterManager {
	return &BFEClusterManager{
		txn:      txn,
		storager: storager,
	}
}

func (pm *BFEClusterManager) FetchBFEClusters(ctx context.Context, param *BFEClusterFilter) (list []*BFECluster, err error) {
	err = pm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err = pm.storager.FetchBFEClusters(ctx, param)
		return err
	})

	return
}

func (pm *BFEClusterManager) CreateBFECluster(ctx context.Context, param *BFEClusterParam) (err error) {
	return pm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err := pm.storager.FetchBFEClusters(ctx, &BFEClusterFilter{
			Name: param.Name,
		})
		if err != nil {
			return err
		}
		if len(list) != 0 {
			return xerror.WrapRecordExisted("BFE Cluster")
		}

		return pm.storager.CreateBFECluster(ctx, param)
	})
}

func (pm *BFEClusterManager) DeleteBFECluster(ctx context.Context, param *BFEClusterParam) (err error) {
	return pm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		_param := &BFEClusterFilter{
			Name: param.Name,
		}
		list, err := pm.storager.FetchBFEClusters(ctx, _param)
		if err != nil {
			return err
		}
		if len(list) == 0 {
			return xerror.WrapRecordNotExist("BFE Cluster")
		}

		return pm.storager.DeleteBFECluster(ctx, list[0])
	})
}

func BFEClusterID2NameMap(list []*BFECluster) map[int64]string {
	m := make(map[int64]string, len(list))
	for _, one := range list {
		m[one.ID] = one.Name
	}

	return m
}

func BFEClusterIDMap(list []*BFECluster) map[int64]*BFECluster {
	m := make(map[int64]*BFECluster, len(list))
	for _, one := range list {
		m[one.ID] = one
	}

	return m
}

func BFEClusterNameMap(list []*BFECluster) map[string]*BFECluster {
	m := make(map[string]*BFECluster, len(list))
	for _, one := range list {
		m[one.Name] = one
	}

	return m
}
