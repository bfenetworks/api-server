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
	"strings"

	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosPoolInstanceStorage struct {
	client naming_client.INamingClient
}

func NewNacosInstancePoolStorage(client naming_client.INamingClient) *NacosPoolInstanceStorage {
	return &NacosPoolInstanceStorage{
		client: client,
	}
}

func (rpps *NacosPoolInstanceStorage) UpdateInstances(ctx context.Context, pool *icluster_conf.Pool, pis *icluster_conf.InstancePool) error {

	return nil
}

func (rpps *NacosPoolInstanceStorage) BatchFetchInstances(ctx context.Context, poolList []*icluster_conf.Pool) (map[string]*icluster_conf.InstancePool, error) {
	m := map[string]*icluster_conf.InstancePool{}
	for _, one := range poolList {
		pi, err := rpps.GetInstance(one.Name[strings.Index(one.Name, ".")+1:])
		pi.Name = one.Name
		if err != nil {
			return nil, err
		}
		m[pi.Name] = pi
	}

	return m, nil
}

func (rpps *NacosPoolInstanceStorage) GetInstance(name string) (*icluster_conf.InstancePool, error) {
	selectInstancesParam := vo.SelectInstancesParam{
		ServiceName: name,
		HealthyOnly: true,
	}
	instances, err := rpps.client.SelectInstances(selectInstancesParam)
	if err != nil {
		return nil, err
	}
	bfeInstances := make([]icluster_conf.Instance, len(instances))
	for index, instance := range instances {
		bfeInstances[index] = newBFEInstance(instance)
	}

	return &icluster_conf.InstancePool{Name: name, Instances: bfeInstances}, nil
}

func newBFEInstance(instance model.Instance) icluster_conf.Instance {
	bfeInstance := icluster_conf.Instance{
		IP:       instance.Ip,
		Ports:    map[string]int{"Default": int(instance.Port)},
		Weight:   int64(instance.Weight),
		HostName: instance.ServiceName,
		Tags:     instance.Metadata,
	}
	return bfeInstance
}
