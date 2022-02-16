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
	"encoding/json"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao"
)

type RDBPoolStorager struct {
	dbCtxFactory lib.DBContextFactory

	productStorager ibasic.ProductStorager
}

func NewRDBPoolStorager(dbCtxFactory lib.DBContextFactory,
	productStorager ibasic.ProductStorager) *RDBPoolStorager {

	return &RDBPoolStorager{
		dbCtxFactory:    dbCtxFactory,
		productStorager: productStorager,
	}
}

var _ icluster_conf.PoolStorager = &RDBPoolStorager{}

func poolFilter2Param(filter *icluster_conf.PoolFilter) *dao.TPoolsParam {
	if filter == nil {
		return nil
	}

	return &dao.TPoolsParam{
		Id:        filter.ID,
		Ids:       filter.IDs,
		Name:      filter.Name,
		ProductID: filter.ProductID,
	}
}

func poolParami2d(data *icluster_conf.PoolParam) (*dao.TPoolsParam, error) {
	if data == nil {
		return nil, nil
	}

	return &dao.TPoolsParam{
		Id:        data.ID,
		Name:      data.Name,
		ProductID: data.ProductID,
		Type:      data.Type,
		Tag:       data.Tag,
	}, nil
}

func (rpps *RDBPoolStorager) CreatePool(ctx context.Context, product *ibasic.Product,
	data *icluster_conf.PoolParam) (*icluster_conf.Pool, error) {

	data.ProductID = &product.ID
	param, err := poolParami2d(data)
	if err != nil {
		return nil, err
	}

	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	_, err = dao.TPoolsCreate(dbCtx, param)
	if err != nil {
		return nil, err
	}

	return rpps.FetchPool(ctx, *param.Name)

}

func (rpps *RDBPoolStorager) FetchPool(ctx context.Context, name string) (*icluster_conf.Pool, error) {
	list, err := rpps.FetchPools(ctx, &icluster_conf.PoolFilter{
		Name: &name,
	})
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func newPool(pp *dao.TPools, product *ibasic.Product) (*icluster_conf.Pool, error) {
	data := &icluster_conf.Pool{
		ID:      pp.Id,
		Name:    pp.Name,
		Ready:   pp.Ready,
		Type:    pp.Type,
		Product: product,

		Tag: pp.Tag,
	}

	// get default instance list
	if pp.InstanceDetail == "" || pp.InstanceDetail == "NULL" {
		pp.InstanceDetail = "[]"
	}
	is := []icluster_conf.Instance{}
	if err := json.Unmarshal([]byte(pp.InstanceDetail), &is); err != nil {
		return nil, xerror.WrapDirtyDataErrorWithMsg("pool %s, raw: %s, err: %v", pp.Name, pp.InstanceDetail, err)
	}
	data.SetDefaultInstances(is)

	return data, nil
}

func (rpps *RDBPoolStorager) FetchPools(ctx context.Context, filter *icluster_conf.PoolFilter) ([]*icluster_conf.Pool, error) {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	ppList, err := dao.TPoolsList(dbCtx, poolFilter2Param(filter))
	if err != nil {
		return nil, err
	}

	if len(ppList) == 0 {
		return nil, nil
	}

	productIDMap := map[int64]bool{}
	for _, one := range ppList {
		productIDMap[one.ProductId] = true
	}

	productList, err := rpps.productStorager.FetchProducts(dbCtx, &ibasic.ProductFilter{
		IDs: lib.Int64BoolMap2Slice(productIDMap),
	})
	if err != nil {
		return nil, err
	}

	productMap := ibasic.ProductIDMap(productList)

	rst := []*icluster_conf.Pool{}
	for _, one := range ppList {
		p, err := newPool(one, productMap[one.ProductId])
		if err != nil {
			return nil, err
		}
		rst = append(rst, p)
	}
	// rpps.registerServier.GetRegisteredInstance(rst)

	return rst, nil
}

func (rpps *RDBPoolStorager) UpdatePool(ctx context.Context, oldData *icluster_conf.Pool,
	diff *icluster_conf.PoolParam) error {

	p, err := poolParami2d(diff)
	if err != nil {
		return err
	}

	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return nil
	}

	_, err = dao.TPoolsUpdate(dbCtx, p, &dao.TPoolsParam{
		Id: &oldData.ID,
	})

	return err
}

func (rpps *RDBPoolStorager) DeletePool(ctx context.Context, pool *icluster_conf.Pool) error {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return nil
	}

	_, err = dao.TPoolsDelete(dbCtx, &dao.TPoolsParam{
		Name: &pool.Name,
	})

	return err
}
