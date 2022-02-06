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

package register

import (
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/stateful"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type RegsiterNacos struct {
	RegisterInfo stateful.RegisterInfo
	client       naming_client.INamingClient
}

func (register *RegsiterNacos) SetRegisterInfo(registerInfo stateful.RegisterInfo) {
	register.RegisterInfo = registerInfo
}

func (register *RegsiterNacos) Init() error {
	registerInfo := register.RegisterInfo
	sc := make([]constant.ServerConfig, len(registerInfo.Address))
	for addressIndex, address := range registerInfo.Address {
		sc[addressIndex] = constant.ServerConfig{
			IpAddr: address.IpAddr,
			Port:   address.Port,
		}
	}
	cc := constant.ClientConfig{
		NamespaceId:         register.RegisterInfo.NameSpace, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: false,
		LogDir:              "./log",
		CacheDir:            "./cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return err
	}
	register.client = client
	return nil
}

func (regsiter *RegsiterNacos) GetInstance(name string) ([]icluster_conf.Instance, error) {
	selectInstancesParam := vo.SelectInstancesParam{
		ServiceName: name,
		HealthyOnly: true,
	}
	instances, err := regsiter.client.SelectInstances(selectInstancesParam)
	if err != nil {
		return nil, err
	}
	bfeInstances := make([]icluster_conf.Instance, len(instances))
	for index, instance := range instances {
		bfeInstances[index] = CreateBfeInstance(instance)
	}
	return bfeInstances, nil
}

func CreateBfeInstance(instance model.Instance) icluster_conf.Instance {
	bfeInstance := icluster_conf.Instance{
		IP:       instance.Ip,
		Ports:    map[string]int{"Default": int(instance.Port)},
		Weight:   int64(instance.Weight),
		HostName: instance.ServiceName,
		Tags:     instance.Metadata,
	}
	return bfeInstance
}
