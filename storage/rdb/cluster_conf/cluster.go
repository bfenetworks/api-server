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

var _ icluster_conf.ClusterStorager = &RDBClusterStorager{}

func NewRDBClusterStorager(dbCtxFactory lib.DBContextFactory,
	subClusterStorager icluster_conf.SubClusterStorager) *RDBClusterStorager {

	return &RDBClusterStorager{
		dbCtxFactory:       dbCtxFactory,
		subClusterStorager: subClusterStorager,
	}
}

type RDBClusterStorager struct {
	dbCtxFactory       lib.DBContextFactory
	subClusterStorager icluster_conf.SubClusterStorager
}

func (rm *RDBClusterStorager) ClusterUpdate(ctx context.Context, product *ibasic.Product, old *icluster_conf.Cluster,
	param *icluster_conf.ClusterParam) error {

	dbCtx, err := rm.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	clusterID := old.ID

	if mlb := param.Scheduler; mlb != nil {
		if err = rm.upsertLBMatrix(dbCtx, old, product, mlb); err != nil {
			return err
		}
	}

	daoClusterParam := newDaoClusterParam(param)
	daoClusterParam.UpdatedAt = lib.PTimeNow()
	_, err = dao.TClusterUpdate(dbCtx, daoClusterParam, &dao.TClusterParam{
		ID: &clusterID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (rm *RDBClusterStorager) ClusterCreate(ctx context.Context, product *ibasic.Product,
	param *icluster_conf.ClusterParam, subClusters []*icluster_conf.SubCluster) (int64, error) {

	dbCtx, err := rm.dbCtxFactory(ctx)
	if err != nil {
		return 0, err
	}

	daoClusterParam := newDaoClusterParam(param)
	clusterID, err := dao.TClusterCreate(dbCtx, daoClusterParam)
	if err != nil {
		return 0, err
	}

	cluster := &icluster_conf.Cluster{
		ID: clusterID,
	}

	if mlb := param.Scheduler; mlb != nil {
		if err = rm.upsertLBMatrix(dbCtx, cluster, product, mlb); err != nil {
			return 0, err
		}
	}

	return clusterID, nil
}

func (rm *RDBClusterStorager) ClusterDelete(ctx context.Context, product *ibasic.Product, cluster *icluster_conf.Cluster) error {
	dbCtx, err := rm.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	clusterID := cluster.ID
	if _, err = dao.TLbMatrixDelete(dbCtx, &dao.TLbMatrixParam{
		ClusterID: &clusterID,
	}); err != nil {
		return err
	}

	_, err = dao.TClusterDelete(dbCtx, &dao.TClusterParam{
		ID: &clusterID,
	})

	return err
}

func (rm *RDBClusterStorager) upsertLBMatrix(dbCtx *lib.DBContext, cluster *icluster_conf.Cluster,
	product *ibasic.Product, lbMatrix map[string]map[string]int) error {

	clusterID := cluster.ID
	_, err := dao.TLbMatrixDelete(dbCtx, &dao.TLbMatrixParam{
		ClusterID: &clusterID,
	})
	if err != nil {
		return err
	}

	lbMatrixs, err := json.Marshal(lbMatrix)
	if err != nil {
		return xerror.WrapParamErrorWithMsg("LBMatrix Marshal fail, err: %v", err)
	}

	_, err = dao.TLbMatrixCreate(dbCtx, &dao.TLbMatrixParam{
		ClusterID: &clusterID,
		ProductID: &product.ID,
		LbMatrix:  lib.PString(string(lbMatrixs)),
		UpdatedAt: lib.PTimeNow(),
		CreatedAt: lib.PTimeNow(),
	})

	return err
}

func (rm *RDBClusterStorager) FetchClusterList(ctx context.Context, filter *icluster_conf.ClusterFilter) ([]*icluster_conf.Cluster, error) {
	dbCtx, err := rm.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	list, err := dao.TClusterList(dbCtx, clusterFilter2Param(filter))
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}

	clusterIDs := []int64{}
	for _, one := range list {
		clusterIDs = append(clusterIDs, one.ID)
	}

	_subClusterList, err := rm.subClusterStorager.FetchSubClusterList(dbCtx, &icluster_conf.SubClusterFilter{
		ClusterIDs: clusterIDs,
	})
	if err != nil {
		return nil, err
	}

	lbMatrixMap, err := rm.fetchLBMatrixs(dbCtx, clusterIDs)
	if err != nil {
		return nil, err
	}

	rst := []*icluster_conf.Cluster{}
	for _, one := range list {
		subClusterList, capacites := []*icluster_conf.SubCluster{}, map[string]int64{}
		for _, subCluster := range _subClusterList {
			if subCluster.ClusterID == one.ID {
				capacites[subCluster.Name] = subCluster.Capacity
				subClusterList = append(subClusterList, subCluster)
			}
		}

		rst = append(rst, newCluster(one, subClusterList, capacites, lbMatrixMap[one.ID]))
	}

	return rst, nil
}

func (rm *RDBClusterStorager) fetchLBMatrixs(dbCtx *lib.DBContext, clusterIDs []int64) (map[int64]map[string]map[string]int, error) {
	manualLbList, err := dao.TLbMatrixList(dbCtx, &dao.TLbMatrixParam{
		ClusterIDs: clusterIDs,
	})
	if err != nil {
		return nil, err
	}

	rst := map[int64]map[string]map[string]int{}
	for _, one := range manualLbList {
		data := map[string]map[string]int{}
		if err := json.Unmarshal([]byte(one.LbMatrix), &data); err != nil {
			return nil, xerror.WrapDirtyDataErrorWithMsg("LbMatrix, err: %v, raw datat: %s", err, one.LbMatrix)
		}
		rst[one.ClusterID] = data
	}

	return rst, nil
}

func newCluster(dc *dao.TCluster, subClusters []*icluster_conf.SubCluster, capacities map[string]int64,
	scheduler map[string]map[string]int) *icluster_conf.Cluster {

	return &icluster_conf.Cluster{
		ID:          dc.ID,
		Name:        dc.Name,
		ProductID:   dc.ProductID,
		Ready:       dc.Ready,
		Description: dc.Description,

		Basic: &icluster_conf.ClusterBasic{
			Connection: &icluster_conf.ClusterBasicConnection{
				MaxIdleConnPerRs:    dc.MaxIdleConnPerHost,
				CancelOnClientClose: dc.CancelOnClientClose,
			},
			Retries: &icluster_conf.ClusterBasicRetries{
				MaxRetryInSubcluster:    dc.MaxRetryInCluster,
				MaxRetryCrossSubcluster: dc.MaxRetryCrossCluster,
			},
			Buffers: &icluster_conf.ClusterBasicBuffers{
				ReqWriteBufferSize: dc.ReqWriteBufferSize,
				ReqFlushInterval:   dc.ReqFlushInterval,
				ResFlushInterval:   dc.ResFlushInterval,
			},
			Timeouts: &icluster_conf.ClusterBasicTimeouts{
				TimeoutConnServ:        dc.TimeoutConnServ,
				TimeoutResponseHeader:  dc.TimeoutResponseHeader,
				TimeoutReadbodyClient:  dc.TimeoutReadbodyClient,
				TimeoutReadClientAgain: dc.TimeoutReadClientAgain,
				TimeoutWriteClient:     dc.TimeoutWriteClient,
			},
		},

		StickySessions: &icluster_conf.ClusterStickySessions{
			HashStrategy:  dc.HashStrategy,
			HashHeader:    dc.HashHeader,
			SessionSticky: dc.SessionSticky,
		},

		Scheduler: scheduler,
		PassiveHealthCheck: &icluster_conf.ClusterPassiveHealthCheck{
			Schema:     dc.HealthcheckSchem,
			Interval:   dc.HealthcheckInterval,
			Failnum:    dc.HealthcheckFailnum,
			Host:       dc.HealthcheckHost,
			Uri:        dc.HealthcheckUri,
			Statuscode: dc.HealthcheckStatuscode,
		},

		SubClusters: subClusters,
	}
}

func clusterFilter2Param(filter *icluster_conf.ClusterFilter) *dao.TClusterParam {
	if filter == nil {
		return nil
	}

	var productID *int64
	if filter.Product != nil {
		productID = &filter.Product.ID
	}
	return &dao.TClusterParam{
		ID:  filter.ID,
		IDs: filter.IDs,

		Name:  filter.Name,
		Names: filter.Names,

		ProductID: productID,
	}
}

func newDaoClusterParam(param *icluster_conf.ClusterParam) *dao.TClusterParam {
	if param == nil {
		return nil
	}

	dc := &dao.TClusterParam{
		ID:          param.ID,
		Name:        param.Name,
		ProductID:   param.ProductID,
		Description: param.Description,
	}

	if basic := param.Basic; basic != nil {
		if conn := basic.Connection; conn != nil {
			dc.MaxIdleConnPerHost = conn.MaxIdleConnPerRs
			dc.CancelOnClientClose = conn.CancelOnClientClose
		}

		if retries := basic.Retries; retries != nil {
			dc.MaxRetryInCluster = retries.MaxRetryInSubcluster
			dc.MaxRetryCrossCluster = retries.MaxRetryCrossSubcluster
		}

		if buffer := basic.Buffers; buffer != nil {
			dc.ReqWriteBufferSize = buffer.ReqWriteBufferSize
			dc.ReqFlushInterval = buffer.ReqFlushInterval
			dc.ResFlushInterval = buffer.ResFlushInterval
		}

		if timeouts := basic.Timeouts; timeouts != nil {
			dc.TimeoutConnServ = timeouts.TimeoutConnServ
			dc.TimeoutResponseHeader = timeouts.TimeoutResponseHeader
			dc.TimeoutReadbodyClient = timeouts.TimeoutReadbodyClient
			dc.TimeoutReadClientAgain = timeouts.TimeoutReadClientAgain
			dc.TimeoutWriteClient = timeouts.TimeoutWriteClient
		}
	}

	if stickySessions := param.StickySessions; stickySessions != nil {
		dc.HashStrategy = stickySessions.HashStrategy
		dc.HashHeader = stickySessions.HashHeader
		dc.SessionSticky = stickySessions.SessionSticky
	}

	if passive := param.PassiveHealthCheck; passive != nil {
		dc.HealthcheckSchem = passive.Schema
		dc.HealthcheckInterval = passive.Interval
		dc.HealthcheckFailnum = passive.Failnum
		dc.HealthcheckHost = passive.Host
		dc.HealthcheckUri = passive.Uri
		dc.HealthcheckStatuscode = passive.Statuscode
	}

	return dc
}

func (rm *RDBClusterStorager) BindSubCluster(ctx context.Context, cluster *icluster_conf.Cluster,
	appendSubClusters, unbindSubClusters []*icluster_conf.SubCluster) error {

	dbCtx, err := rm.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	var addingPoolNames, removePoolNames []string
	if len(unbindSubClusters) > 0 {
		if _, err = dao.TSubClusterUpdate(dbCtx, &dao.TSubClusterParam{
			ClusterID: lib.PInt64(-1),
		}, &dao.TSubClusterParam{
			ClusterID: &cluster.ID,
			IDs:       icluster_conf.SubClusterList2IDSlice(unbindSubClusters),
		}); err != nil {
			return err
		}

		for _, one := range unbindSubClusters {
			if one.InstancePool.Tag == icluster_conf.PoolTagProduct {
				removePoolNames = append(removePoolNames, one.InstancePool.Name)
			}
		}
	}

	if len(appendSubClusters) > 0 {
		if _, err = dao.TSubClusterUpdate(dbCtx, &dao.TSubClusterParam{
			ClusterID: lib.PInt64(cluster.ID),
		}, &dao.TSubClusterParam{
			IDs: icluster_conf.SubClusterList2IDSlice(appendSubClusters),
		}); err != nil {
			return err
		}

		for _, one := range appendSubClusters {
			if one.InstancePool.Tag == icluster_conf.PoolTagProduct {
				addingPoolNames = append(addingPoolNames, one.InstancePool.Name)
			}
		}
	}

	return nil
}

func (rm *RDBClusterStorager) FetchCluster(ctx context.Context, filter *icluster_conf.ClusterFilter) (*icluster_conf.Cluster, error) {
	list, err := rm.FetchClusterList(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}
