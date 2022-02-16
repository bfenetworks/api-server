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

package stateful

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosRegisterConfig struct {
	ServerConfig []constant.ServerConfig
	ClientConfig constant.ClientConfig
}

var NacosClient naming_client.INamingClient

func (d *Config) InitNacos() error {
	if d.NacosRegsiter.ClientConfig.NamespaceId == "" || len(d.NacosRegsiter.ServerConfig) == 0 {
		fmt.Println("The configuration is not sound and naocs will not be started")
		return nil
	}
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &d.NacosRegsiter.ClientConfig,
			ServerConfigs: d.NacosRegsiter.ServerConfig,
		},
	)
	if err != nil {
		msg := fmt.Sprintf("nacos start err - %s", err)
		fmt.Println("nacos start err :", msg)
		return err
	}
	NacosClient = client
	return nil
}

func GetNacosClient() (naming_client.INamingClient, error) {
	return NacosClient, nil
}
