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
	register "github.com/bfenetworks/api-server/model/register/nacos"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

type NacosRegisterConfig struct {
	ServerConfig []constant.ServerConfig
	ClientConfig constant.ClientConfig
}

func (d *Config) InitRegister() error {

	d.initNacosRegister()
	return nil
}

func (d *Config) initNacosRegister() error {
	registerObject := register.RegsiterNacos{ClientConfig: d.NacosRegsiter.ClientConfig, ServerConfig: d.NacosRegsiter.ServerConfig}
	registerObject.Init()
	return nil
}
