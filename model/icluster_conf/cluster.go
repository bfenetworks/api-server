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

	"github.com/bfenetworks/bfe/bfe_config/bfe_cluster_conf/cluster_conf"
	"github.com/bfenetworks/bfe/bfe_config/bfe_route_conf/route_rule_conf"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/itxn"
	"github.com/bfenetworks/api-server/model/iversion_control"
	"github.com/bfenetworks/api-server/stateful"
)

var (
	ClusterHashStrategyClientIDOnlyI     int32 = 0
	ClusterHashStrategyClientIPOnlyI     int32 = 1
	ClusterHashStrategyClientIDPreferedI int32 = 2

	ClusterHealthCheckHTTP      = "http"
	ClusterHealthCheckTCP       = "tcp"
	ClusterHealthCheckSchemaMap = map[string]bool{
		ClusterHealthCheckHTTP: true,
		ClusterHealthCheckTCP:  true,
	}

	ClusterStickTypeSubCluster = "SUB_CLUSTER"
	ClusterStickTypeInstance   = "INSTANCE"

	ClusterDefaultReqFlushInterval int32 = 0
	ClusterDefaultResFlushInterval int32 = -1 // -1: write response directly without using timing refresh
)

type ClusterBasicConnectionParam struct {
	MaxIdleConnPerRs    *int16
	CancelOnClientClose *bool
}

type ClusterBasicBuffersParam struct {
	ReqWriteBufferSize *int32
	ReqFlushInterval   *int32
	ResFlushInterval   *int32
}

type ClusterBasicRetriesParam struct {
	MaxRetryInSubcluster    *int8
	MaxRetryCrossSubcluster *int8
}

type ClusterBasicTimeoutsParam struct {
	TimeoutConnServ        *int32
	TimeoutResponseHeader  *int32
	TimeoutReadbodyClient  *int32
	TimeoutReadClientAgain *int32
	TimeoutWriteClient     *int32
}

type ClusterBasicParam struct {
	Connection *ClusterBasicConnectionParam
	Retries    *ClusterBasicRetriesParam
	Buffers    *ClusterBasicBuffersParam
	Timeouts   *ClusterBasicTimeoutsParam
}

type ClusterStickySessionsParam struct {
	SessionSticky *bool
	HashStrategy  *int32
	HashHeader    *string
}

type ClusterPassiveHealthCheckParam struct {
	Schema     *string
	Interval   *int32
	Failnum    *int32
	Statuscode *int32
	Host       *string
	Uri        *string
}

type ClusterFilter struct {
	ID  *int64
	IDs []int64

	Names []string
	Name  *string

	Product *ibasic.Product
}

type ClusterParam struct {
	ID *int64

	Name *string

	ProductID *int64

	Description *string

	Basic *ClusterBasicParam

	StickySessions *ClusterStickySessionsParam

	SubClusters []string

	Scheduler map[string]map[string]int

	PassiveHealthCheck *ClusterPassiveHealthCheckParam
}

type ClusterBasicConnection struct {
	MaxIdleConnPerRs    int16
	CancelOnClientClose bool
}

type ClusterBasicBuffers struct {
	ReqWriteBufferSize int32
	ReqFlushInterval   int32
	ResFlushInterval   int32
}

type ClusterBasicRetries struct {
	MaxRetryInSubcluster    int8
	MaxRetryCrossSubcluster int8
}

type ClusterBasicTimeouts struct {
	TimeoutConnServ        int32
	TimeoutResponseHeader  int32
	TimeoutReadbodyClient  int32
	TimeoutReadClientAgain int32
	TimeoutWriteClient     int32
}

type ClusterBasic struct {
	Connection *ClusterBasicConnection
	Retries    *ClusterBasicRetries
	Buffers    *ClusterBasicBuffers
	Timeouts   *ClusterBasicTimeouts
}

type ClusterStickySessions struct {
	SessionSticky bool
	HashStrategy  int32
	HashHeader    string
}

type ClusterPassiveHealthCheck struct {
	Schema     string
	Interval   int32
	Failnum    int32
	Statuscode int32
	Host       string
	Uri        string
}

func (cphc *ClusterPassiveHealthCheck) toBackendCheck() *cluster_conf.BackendCheck {
	if cphc == nil {
		return nil
	}

	int322intp := func(i int32) *int {
		tmp := int(i)
		return &tmp
	}

	return &cluster_conf.BackendCheck{
		Schem:         &cphc.Schema,
		Uri:           &cphc.Uri,
		Host:          &cphc.Host,
		FailNum:       int322intp(cphc.Failnum),
		CheckInterval: int322intp(cphc.Interval),
		StatusCode:    int322intp(cphc.Statuscode),
	}
}

type Cluster struct {
	ID          int64
	Name        string
	Description string
	Ready       bool
	ProductID   int64

	Basic              *ClusterBasic
	StickySessions     *ClusterStickySessions
	SubClusters        []*SubCluster
	Scheduler          map[string]map[string]int
	PassiveHealthCheck *ClusterPassiveHealthCheck
}

func (cluster *Cluster) SubClusterNames() []string {
	var names []string
	for _, one := range cluster.SubClusters {
		names = append(names, one.Name)
	}

	return names
}

func ClusterList2MapByName(list []*Cluster) map[string]*Cluster {
	m := map[string]*Cluster{}
	for _, one := range list {
		m[one.Name] = one
	}

	return m
}

const (
	ResourceClusterRule = "cluster_rule"
)

func ClusterList2MapByID(list []*Cluster) map[int64]*Cluster {
	m := map[int64]*Cluster{}
	for _, one := range list {
		m[one.ID] = one
	}

	return m
}

func NewClusterManager(txn itxn.TxnStorager, storager ClusterStorager,
	subClusterStorager SubClusterStorager, bfeClusterStorager ibasic.BFEClusterStorager,
	versionControlManager *iversion_control.VersionControlManager,
	deleteCheckers map[string]func(context.Context, *ibasic.Product, *Cluster) error) *ClusterManager {

	return &ClusterManager{
		txn:                   txn,
		storager:              storager,
		subClusterStorager:    subClusterStorager,
		bfeClusterStorager:    bfeClusterStorager,
		versionControlManager: versionControlManager,

		deleteCheckers: deleteCheckers,
	}
}

type ClusterStorager interface {
	FetchCluster(ctx context.Context, param *ClusterFilter) (*Cluster, error)
	FetchClusterList(ctx context.Context, param *ClusterFilter) ([]*Cluster, error)
	ClusterUpdate(ctx context.Context, product *ibasic.Product, old *Cluster, param *ClusterParam) error
	ClusterCreate(ctx context.Context, product *ibasic.Product, param *ClusterParam, subClusters []*SubCluster) (int64, error)
	ClusterDelete(ctx context.Context, product *ibasic.Product, cluster *Cluster) error
	BindSubCluster(ctx context.Context, cluster *Cluster, appendSubClusters, unbindSubClusters []*SubCluster) error
}

type ClusterManager struct {
	txn itxn.TxnStorager

	storager           ClusterStorager
	subClusterStorager SubClusterStorager
	bfeClusterStorager ibasic.BFEClusterStorager

	versionControlManager *iversion_control.VersionControlManager

	deleteCheckers map[string]func(context.Context, *ibasic.Product, *Cluster) error
}

func (rm *ClusterManager) FetchClusterList(ctx context.Context, param *ClusterFilter) (list []*Cluster, err error) {
	err = rm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err = rm.storager.FetchClusterList(ctx, param)
		return err
	})

	return
}

func (cm *ClusterManager) FetchCluster(ctx context.Context, param *ClusterFilter) (*Cluster, error) {
	list, err := cm.FetchClusterList(ctx, param)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (cm *ClusterManager) CreateCluster(ctx context.Context, product *ibasic.Product, param *ClusterParam) (err error) {
	param.ProductID = &product.ID

	err = cm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		old, err := cm.storager.FetchClusterList(ctx, &ClusterFilter{
			Name: param.Name,
		})
		if err != nil {
			return err
		}
		if len(old) != 0 {
			return xerror.WrapRecordExisted("cluster")
		}

		bindingSubClusters, err := cm.subClusterStorager.FetchSubClusterList(ctx, &SubClusterFilter{
			Names:   param.SubClusters,
			Product: product,
		})
		if err != nil {
			return err
		}
		err = cm.checkBindingSubClusters(ctx, nil, param.SubClusters, bindingSubClusters)
		if err != nil {
			return err
		}
		if err := cm.checkManualLB(ctx, nil, param); err != nil {
			return err
		}

		if param.Scheduler == nil {
			if param.Scheduler, err = cm.constructDefaultScheduler(ctx, bindingSubClusters); err != nil {
				return err
			}
		}

		clusterID, err := cm.storager.ClusterCreate(ctx, product, param, bindingSubClusters)
		if err != nil {
			return err
		}

		return cm.storager.BindSubCluster(ctx, &Cluster{
			ID: clusterID,
		}, bindingSubClusters, nil)
	})

	return
}

const BlackHole = "GSLB_BLACKHOLE"

func (cm *ClusterManager) constructDefaultScheduler(ctx context.Context, subClusters []*SubCluster) (map[string]map[string]int, error) {
	bfeClusters, err := cm.bfeClusterStorager.FetchBFEClusters(ctx, nil)
	if err != nil {
		return nil, err
	}

	lbMatrix := map[string]map[string]int{}

	rate := 100 / len(subClusters)
	mod := 100 - rate*len(subClusters)
	for _, bfeCluster := range bfeClusters {
		tmp := map[string]int{
			BlackHole: mod,
		}
		for _, subCluster := range subClusters {
			tmp[subCluster.Name] = rate
		}

		lbMatrix[bfeCluster.Name] = tmp
	}

	return lbMatrix, nil
}

func (cm *ClusterManager) checkManualLB(ctx context.Context, old *Cluster, param *ClusterParam) error {
	if param.Scheduler == nil {
		return nil
	}

	bfeClusters, err := cm.bfeClusterStorager.FetchBFEClusters(ctx, nil)
	if err != nil {
		return err
	}

	lbMatrix := param.Scheduler
	if len(bfeClusters) != len(lbMatrix) {
		return xerror.WrapParamErrorWithMsg("LbMatrix Config Illegal, Want All BFE Cluster Exist")
	}

	subClusters := param.SubClusters
	if subClusters == nil && old != nil {
		subClusters = old.SubClusterNames()
	}

	bfeClusterMap := ibasic.BFEClusterNameMap(bfeClusters)
	for bfeClusterName, subClusterLbMatrix := range lbMatrix {
		if _, ok := bfeClusterMap[bfeClusterName]; !ok {
			return xerror.WrapParamErrorWithMsg("LbMatrix Config Illegal, BFE Cluster %s Not Exist", bfeClusterName)
		}

		total := 0
		for subClusterName, rate := range subClusterLbMatrix {
			if subClusterName != BlackHole && !lib.StringSliceHasElement(subClusters, subClusterName) {
				return xerror.WrapParamErrorWithMsg("LbMatrix Config Illegal, SubCluster %s Not In BFE Cluster %s Config", bfeClusterName, subClusterName)
			}
			if rate < 0 {
				return xerror.WrapParamErrorWithMsg("LbMatrix Config Illegal, BFE Cluster %s Rate Must Bigger Than 0, Got %d", bfeClusterName, rate)
			}
			total += rate
		}
		if total != 100 {
			return xerror.WrapParamErrorWithMsg("LbMatrix Config Illegal, BFE Cluster %s Total Rate Is %d, Want 100", bfeClusterName, total)
		}

		for _, subCluster := range subClusters {
			if _, ok := subClusterLbMatrix[subCluster]; !ok {
				return xerror.WrapParamErrorWithMsg("LbMatrix Config Illegal, SubCluster %s Not In BFE Cluster %s Config", bfeClusterName, subCluster)
			}
		}
	}

	return nil
}

func (cm *ClusterManager) checkBindingSubClusters(ctx context.Context, cluster *Cluster,
	bindingSubClusterNames []string, bindSubClusters []*SubCluster) error {

	if len(bindSubClusters) == 0 {
		return xerror.WrapModelErrorWithMsg("Cluster Want At Least On SubCluster")
	}

	var oldClusterID int64
	if cluster != nil {
		oldClusterID = cluster.ID
	}
	subClusterMap := SubClusterList2MapByName(bindSubClusters)
	for _, scName := range bindingSubClusterNames {
		subCluster := subClusterMap[scName]
		if subCluster == nil {
			return xerror.WrapModelErrorWithMsg("SubCluster %s Not Exist", scName)
		}

		if subCluster.ClusterID != oldClusterID && subCluster.ClusterID > 0 {
			return xerror.WrapModelErrorWithMsg("SubCluster %s be Mounted With Cluster %d", scName, subCluster.ClusterID)
		}

		if !stateful.IgnoreBNSStatusCheck && !subCluster.Ready {
			return xerror.WrapDependentUnReadyErrorWithMsg("SubCluster %s Not Ready", scName)
		}
	}

	return nil
}

func (cm *ClusterManager) UpdateCluster(ctx context.Context, product *ibasic.Product, oldData *Cluster,
	param *ClusterParam) (err error) {

	err = cm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		if err = cm.checkManualLB(ctx, oldData, param); err != nil {
			return err
		}

		return cm.storager.ClusterUpdate(ctx, product, oldData, param)
	})

	return
}

func (cm *ClusterManager) checkLbMatrix(cluster *Cluster, unbindSubClusters, appendSubClusters []string) (map[string]map[string]int, error) {
	unbindSubClusterMap := lib.StringSlice2Map(unbindSubClusters)
	newManualLbMatrix := map[string]map[string]int{}

	for bfeCluster, subClusterRate := range cluster.Scheduler {
		newManualLbMatrix[bfeCluster] = map[string]int{}
		for subClusterName, rate := range subClusterRate {
			if unbindSubClusterMap[subClusterName] {
				if rate := subClusterRate[subClusterName]; rate != 0 {
					return nil, xerror.WrapModelErrorWithMsg(
						"BFE Cluster %s, SubCluster: %s Rate is %d, Set to 0 Before Unbind", bfeCluster, subClusterName, rate)
				}
			} else {
				newManualLbMatrix[bfeCluster][subClusterName] = rate
			}

			for _, subClusterName := range appendSubClusters {
				newManualLbMatrix[bfeCluster][subClusterName] = 0
			}
		}
	}

	return newManualLbMatrix, nil
}

func (cm *ClusterManager) RebindSubCluster(ctx context.Context, product *ibasic.Product, cluster *Cluster,
	bindingSubClusterNames []string) error {

	unbindSubClusterNames := lib.StringSliceSubtract(cluster.SubClusterNames(), bindingSubClusterNames)
	appendSubClusterNames := lib.StringSliceSubtract(bindingSubClusterNames, cluster.SubClusterNames())
	if len(unbindSubClusterNames) == 0 && len(appendSubClusterNames) == 0 {
		return nil
	}

	manualLbMatrix, err := cm.checkLbMatrix(cluster, unbindSubClusterNames, appendSubClusterNames)
	if err != nil {
		return err
	}

	return cm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		bindingSubClusters, err := cm.subClusterStorager.FetchSubClusterList(ctx, &SubClusterFilter{
			Names:   bindingSubClusterNames,
			Product: product,
		})
		if err != nil {
			return err
		}

		if err = cm.checkBindingSubClusters(ctx, cluster, bindingSubClusterNames, bindingSubClusters); err != nil {
			return err
		}

		if manualLbMatrix != nil {
			if err := cm.storager.ClusterUpdate(ctx, product, cluster, &ClusterParam{
				Scheduler: manualLbMatrix,
			}); err != nil {
				return err
			}
		}

		var appendSubClusters, unbindSubClusters []*SubCluster
		if len(unbindSubClusterNames) > 0 {
			scMap := SubClusterList2MapByName(cluster.SubClusters)
			for _, one := range unbindSubClusterNames {
				unbindSubClusters = append(unbindSubClusters, scMap[one])
			}
		}

		scMap := SubClusterList2MapByName(bindingSubClusters)
		for _, one := range appendSubClusterNames {
			appendSubClusters = append(appendSubClusters, scMap[one])
		}

		// U should check param by yourself
		return cm.storager.BindSubCluster(ctx, cluster, appendSubClusters, unbindSubClusters)
	})
}

func (cm *ClusterManager) DeleteCluster(ctx context.Context, product *ibasic.Product, cluster *Cluster) (err error) {
	err = cm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		for _, checker := range cm.deleteCheckers {
			err = checker(ctx, product, cluster)
			if err != nil {
				return err
			}
		}

		if err = cm.storager.BindSubCluster(ctx, cluster, nil, cluster.SubClusters); err != nil {
			return err
		}

		if err = cm.storager.ClusterDelete(ctx, product, cluster); err != nil {
			return err
		}

		return nil
	})

	return
}

var (
	UnMountedClusterID int64 = -1

	RouteAdvancedModeClusterName4DP       = route_rule_conf.AdvancedMode
	RouteAdvancedModeClusterName          = "GO_TO_ADVANCED_RULES"
	RouteAdvancedModeClusterID      int64 = -1

	SystemKeepRouteNames = map[string]bool{
		RouteAdvancedModeClusterName:    true,
		RouteAdvancedModeClusterName4DP: true,
	}
)

func AppendAdvancedRuleCluster(list []*Cluster) []*Cluster {
	return append(list, &Cluster{
		ID:   RouteAdvancedModeClusterID,
		Name: RouteAdvancedModeClusterName4DP,
	})
}

func NewBfeClusterConf(version string, clusters []*Cluster) *cluster_conf.BfeClusterConf {
	clusterConfMap := cluster_conf.ClusterToConf{}

	int322intp := func(i int32) *int {
		tmp := int(i)
		return &tmp
	}
	int162intp := func(i int16) *int {
		tmp := int(i)
		return &tmp
	}
	int82intp := func(i int8) *int {
		tmp := int(i)
		return &tmp
	}

	for _, cluster := range clusters {
		if SystemKeepRouteNames[cluster.Name] {
			continue
		}

		clusterConfMap[cluster.Name] = cluster_conf.ClusterConf{
			BackendConf: &cluster_conf.BackendBasic{
				Protocol:              lib.PString("http"),
				TimeoutConnSrv:        int322intp(cluster.Basic.Timeouts.TimeoutConnServ),
				TimeoutResponseHeader: int322intp(cluster.Basic.Timeouts.TimeoutResponseHeader),
				MaxIdleConnsPerHost:   int162intp(cluster.Basic.Connection.MaxIdleConnPerRs),
			},
			CheckConf: cluster.PassiveHealthCheck.toBackendCheck(),
			GslbBasic: &cluster_conf.GslbBasicConf{
				CrossRetry: int82intp(cluster.Basic.Retries.MaxRetryCrossSubcluster),
				RetryMax:   int82intp(cluster.Basic.Retries.MaxRetryInSubcluster),
				HashConf: &cluster_conf.HashConf{
					HashStrategy:  int322intp(cluster.StickySessions.HashStrategy),
					HashHeader:    &cluster.StickySessions.HashHeader,
					SessionSticky: &cluster.StickySessions.SessionSticky,
				},
				BalanceMode: lib.PString("WRR"),
			},
			ClusterBasic: &cluster_conf.ClusterBasicConf{
				TimeoutReadClient:      int322intp(cluster.Basic.Timeouts.TimeoutReadbodyClient),
				TimeoutWriteClient:     int322intp(cluster.Basic.Timeouts.TimeoutWriteClient),
				TimeoutReadClientAgain: int322intp(cluster.Basic.Timeouts.TimeoutReadClientAgain),
				ReqWriteBufferSize:     int322intp(cluster.Basic.Buffers.ReqWriteBufferSize),
				ReqFlushInterval:       int322intp(cluster.Basic.Buffers.ReqFlushInterval),
				ResFlushInterval:       int322intp(cluster.Basic.Buffers.ResFlushInterval),
				CancelOnClientClose:    &cluster.Basic.Connection.CancelOnClientClose,
			},
		}

	}
	return &cluster_conf.BfeClusterConf{
		Version: &version,
		Config:  &clusterConfMap,
	}
}
