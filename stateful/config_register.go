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
	"io/ioutil"
	"path"

	"github.com/baidu/go-lib/log"
	"gopkg.in/yaml.v3"
)

type Address struct {
	IpAddr string `yaml:"ipAddr"`
	Port   uint64 `yaml:"port"`
}

type RegisterInfo struct {
	Name      string    `yaml:"name"`
	Type      string    `yaml:"type"`
	Address   []Address `yaml:"address"`
	NameSpace string    `yaml:"nameSpace"`
	Config map[string]string `yaml:"config"`
}

type RegisterMainConfig struct {
	Registers []RegisterInfo `yaml:"register"`
}

func GetRegisterConfig(confDir *string) (*RegisterMainConfig, error) {
	var config *RegisterMainConfig
	confPath := path.Join(*confDir, "bfe_register.yaml")
	buffer, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Logger.Error("confg_register.getRegisterConfig(): in BfeRegisterConfigLoad():%s", err.Error())
		return nil, err
	}
	err = yaml.Unmarshal(buffer, &config)
	if err != nil {
		log.Logger.Error("confg_register.getRegisterConfig(): in BfeRegisterConfigLoad():%s", err.Error())
		return nil, err
	}
	return config, nil

}
