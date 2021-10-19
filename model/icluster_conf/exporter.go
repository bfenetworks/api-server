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

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/iversion_control"
	"github.com/bfenetworks/bfe/bfe_config/bfe_cluster_conf/cluster_table_conf"
	"github.com/bfenetworks/bfe/bfe_config/bfe_cluster_conf/gslb_conf"
)

const (
	ConfigTopicClusterTable = "cluster_table"
	ConfigTopicGSLB         = "gslb"
)

type ClusterTableConf struct {
	cluster_table_conf.ClusterTableConf
}

func (ctc *ClusterTableConf) UpdateVersion(version string) error {
	ctc.Version = &version

	return nil
}

func (rm *ClusterManager) clusterTableConfGenerator(ctx context.Context) (*iversion_control.ExportData, error) {
	clusters, err := rm.storager.FetchClusterList(ctx, nil)
	if err != nil {
		return nil, err
	}

	allClusters := cluster_table_conf.AllClusterBackend{}
	for _, cluster := range clusters {
		clusterBackend := map[string]cluster_table_conf.SubClusterBackend{}

		for _, subCluster := range cluster.SubClusters {
			if subCluster.InstancePool == nil || len(subCluster.InstancePool.Instances) == 0 {
				continue
			}

			subClusterBackend := make(cluster_table_conf.SubClusterBackend, 0, len(subCluster.InstancePool.Instances))
			for _, instance := range subCluster.InstancePool.Instances {
				subClusterBackend = append(subClusterBackend, &cluster_table_conf.BackendConf{
					Name:   lib.PString(instance.HostName),
					Addr:   lib.PString(instance.IP),
					Port:   lib.PInt(instance.Port),
					Weight: lib.PInt(int(instance.Weight)),
				})
			}

			clusterBackend[subCluster.Name] = subClusterBackend
		}

		allClusters[cluster.Name] = clusterBackend
	}

	clusterTableConf := &ClusterTableConf{
		ClusterTableConf: cluster_table_conf.ClusterTableConf{
			Config: &allClusters,
		},
	}
	clusterTableConf.UpdateVersion(iversion_control.ZeroVersion)

	if err := clusterTableConf.ClusterTableConf.Config.Check(); err != nil {
		return nil, xerror.WrapModelError(err)
	}

	return &iversion_control.ExportData{
		Topic:              ConfigTopicClusterTable,
		DataWithoutVersion: clusterTableConf,
	}, nil
}

func (rm *ClusterManager) ExportClusterTable(ctx context.Context, lastVersion string) (*ClusterTableConf, error) {
	ed, err := rm.versionControlManager.ExportConfig(ctx, ConfigTopicClusterTable, rm.clusterTableConfGenerator)
	if err != nil {
		return nil, err
	}

	conf := ed.DataWithoutVersion.(*ClusterTableConf)
	if *conf.Version == lastVersion {
		return nil, nil
	}

	return conf, nil
}

type GSLBConf struct {
	Version string
	gslb_conf.GslbConf
}

func (gc *GSLBConf) UpdateVersion(version string) error {
	gc.Version = version
	gc.Ts = &version

	return nil
}

func (rm *ClusterManager) gslbConfGenerator(bfeClusterName string) func(ctx context.Context) (*iversion_control.ExportData, error) {
	return func(ctx context.Context) (*iversion_control.ExportData, error) {
		topic := ConfigTopicGSLB + "." + bfeClusterName
		clusters, err := rm.storager.FetchClusterList(ctx, nil)
		if err != nil {
			return nil, err
		}

		gslbClustersConf := gslb_conf.GslbClustersConf{}
		for _, cluster := range clusters {
			manualScheduler := cluster.ManualScheduler
			if len(manualScheduler) == 0 {
				continue
			}

			lbMatraix := manualScheduler[bfeClusterName]
			if lbMatraix == nil {
				return nil, xerror.WrapParamErrorWithMsg("BFECluster %s Not Exist", bfeClusterName)
			}

			gslbClustersConf[cluster.Name] = lbMatraix
		}

		gslbConf := &GSLBConf{
			GslbConf: gslb_conf.GslbConf{
				Clusters: &gslbClustersConf,
				Hostname: lib.PString("gslb.manual.com"),
			},
		}
		gslbConf.UpdateVersion(iversion_control.ZeroVersion)

		if err := gslbConf.GslbConf.Check(); err != nil {
			return nil, xerror.WrapModelError(err)
		}

		return &iversion_control.ExportData{
			Topic:              topic,
			DataWithoutVersion: gslbConf,
		}, nil
	}
}

func (rm *ClusterManager) ExportGSLB(ctx context.Context, lastVersion, bfeClusterName string) (*GSLBConf, error) {
	topic := ConfigTopicGSLB + "." + bfeClusterName
	ed, err := rm.versionControlManager.ExportConfig(ctx, topic, rm.gslbConfGenerator(bfeClusterName))
	if err != nil {
		return nil, err
	}

	conf := ed.DataWithoutVersion.(*GSLBConf)
	if conf.Version == lastVersion {
		return nil, nil
	}

	return conf, nil
}
