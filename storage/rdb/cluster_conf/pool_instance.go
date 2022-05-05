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
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao"
)

type RDBInstancePoolStorage struct {
	dbCtxFactory lib.DBContextFactory
}

func NewRDBInstancePoolStorage(dbCtxFactory lib.DBContextFactory) *RDBInstancePoolStorage {
	return &RDBInstancePoolStorage{
		dbCtxFactory: dbCtxFactory,
	}
}

var _ icluster_conf.InstancePoolStorage = &RDBInstancePoolStorage{}

func (rpps *RDBInstancePoolStorage) UpdateInstances(ctx context.Context, pool *icluster_conf.Pool,
	pis *icluster_conf.InstancePool) error {

	var detail *string
	if pis.Instances != nil {
		bs, err := json.Marshal(pis.Instances)
		if err != nil {
			return xerror.WrapParamErrorWithMsg("Instances Marshal, err: %s", err)
		}

		detail = lib.PString(string(bs))
	}

	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}
	_, err = dao.TPoolsUpdate(dbCtx, &dao.TPoolsParam{
		InstanceDetail: detail,
	}, &dao.TPoolsParam{
		Id: &pool.ID,
	})

	return err
}

func (rpps *RDBInstancePoolStorage) BatchFetchInstances(ctx context.Context,
	poolList []*icluster_conf.Pool) (map[string]*icluster_conf.InstancePool, error) {

	m := map[string]*icluster_conf.InstancePool{}
	for _, one := range poolList {
		// because of RDBPoolStorage.FetchPools will get pool list
		// it's trick
		pi := one.GetDefaultPool()
		m[pi.Name] = pi
	}

	return m, nil
}
