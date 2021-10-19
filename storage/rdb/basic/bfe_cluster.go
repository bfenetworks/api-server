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

package basic

import (
	"context"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao"
)

var (
	DeleteStatusNo  int8 = 0
	DeleteStatusYes int8 = 1

	TagBFE     int8 = 1
	TagProduct int8 = 3
)

var PoolsTagMap = map[int8]bool{
	TagBFE:     true,
	TagProduct: true,
}

type RDBBFEClusterStorager struct {
	dbCtxFactory lib.DBContextFactory
}

var _ ibasic.BFEClusterStorager = &RDBBFEClusterStorager{}

func NewRDBBFEClusterStorager(dbCtxFactory lib.DBContextFactory) *RDBBFEClusterStorager {
	return &RDBBFEClusterStorager{
		dbCtxFactory: dbCtxFactory,
	}
}

func (ps *RDBBFEClusterStorager) FetchBFEClusters(ctx context.Context, filter *ibasic.BFEClusterFilter) ([]*ibasic.BFECluster, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	var where *dao.TBfeClusterParam
	if filter != nil {
		where = &dao.TBfeClusterParam{
			Name:     filter.Name,
			PoolName: filter.Pool,
		}
	}

	list, err := dao.TBfeClusterList(dbCtx, where)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}

	rst := make([]*ibasic.BFECluster, len(list))
	for i, one := range list {
		rst[i] = &ibasic.BFECluster{
			Name:               one.Name,
			Pool:               one.PoolName,
			Enabled:            one.Enabled,
			ExemptTrafficCheck: one.ExemptTrafficCheck,
			Capacity:           one.Capacity,
		}
	}
	return rst, nil
}

func (ps *RDBBFEClusterStorager) DeleteBFECluster(ctx context.Context, pp *ibasic.BFECluster) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TBfeClusterDelete(dbCtx, &dao.TBfeClusterParam{
		Name: &pp.Name,
	})

	return err
}

func (ps *RDBBFEClusterStorager) CreateBFECluster(ctx context.Context, pp *ibasic.BFEClusterParam) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	pool, err := dao.TPoolsOne(dbCtx, &dao.TPoolsParam{
		Name: pp.Pool,
		Tag:  &TagBFE,
	})
	if err != nil {
		return err
	}
	if pool == nil {
		return xerror.WrapParamErrorWithMsg("Pool %s Not Existed", *pp.Pool)
	}

	_, err = dao.TBfeClusterCreate(dbCtx, &dao.TBfeClusterParam{
		Name:     pp.Name,
		PoolName: &pool.Name,
		Capacity: pp.Capacity,
	})

	return err
}
