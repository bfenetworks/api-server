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
	"strings"

	"github.com/bfenetworks/api-server/model/icluster_conf"
)

type Register interface {
	GetInstance(name string) ([]icluster_conf.Instance, error)
	Init() error
}

type RegisterServier struct {
	RegisterExample map[string]Register
	TypeMapper      map[int]string
}

func (registerServier *RegisterServier) GetRegisteredInstance(pools []*icluster_conf.Pool) {
	for _, pool := range pools {
		if registerType, ok := registerServier.TypeMapper[int(pool.Type)]; ok {
			registerObject, _ok := registerServier.RegisterExample[registerType]
			if !_ok {
				continue
			}
			name := pool.Name[strings.Index(pool.Name, ".")+1:]
			instances, _ := registerObject.GetInstance(name)
			pool.Instances = instances
		}
	}
}
