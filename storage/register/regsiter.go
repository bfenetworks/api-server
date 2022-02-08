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
	register "github.com/bfenetworks/api-server/storage/register/nacos"
)

type Register interface {
	SetRegisterInfo(registerInfo stateful.RegisterInfo)
	GetInstance(name string) ([]icluster_conf.Instance, error)
	//GetInstances(name []string) []icluster_conf.Instance
	Init() error
}

type RegisterServier struct {
	RegisterConfig  *stateful.RegisterMainConfig
	RegisterExample map[string]Register
	TypeMapper      map[int]string
}

func (registerServier *RegisterServier) Init() {
	registerServier.RegisterExample = make(map[string]Register)
	registerServier.TypeMapper = make(map[int]string)
	registerServier.TypeMapper[1] = "nacos"
	for _, registerInfo := range registerServier.RegisterConfig.Registers {
		if registerInfo.Type == "nacos" {
			registerObject := register.RegsiterNacos{RegisterInfo: registerInfo}
			registerObject.Init()
			registerServier.RegisterExample[registerInfo.Type] = &registerObject
		}
	}
}

func (registerServier *RegisterServier) GetRegisteredInstance(pools []*icluster_conf.Pool) {
	panic("")
	// for _, pool := range pools {
	// if registerType, ok := registerServier.TypeMapper[int(pool.Type)]; ok {
	// 	registerObject, _ok := registerServier.RegisterExample[registerType]
	// 	if !_ok {
	// 		continue
	// 	}
	// 	name := pool.Name[strings.Index(pool.Name, ".")+1:]
	// instances, _ := registerObject.GetInstance(name)
	// pool.Instances = instances
	// }
	// }
}
