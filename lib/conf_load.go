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

package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/gcfg.v1"
)

// LoadConf load config from file
func LoadConf(fileName string, data interface{}, fileSuffix string) error {
	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}

	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	switch strings.ToLower(fileSuffix) {
	case ".json":
		return json.Unmarshal(bs, data)
	case ".toml":
		return toml.Unmarshal(bs, data)
	case ".conf":
		return gcfg.ReadStringInto(data, string(bs))
	default:
		return fmt.Errorf("LoadConf, unsupported fileSuffix: %s, file: %s", fileSuffix, fileName)
	}
}

// LoadConfAuto load config from file. auto calculate file type by suffix
func LoadConfAuto(fileName string, data interface{}) error {
	return LoadConf(fileName, data, filepath.Ext(fileName))
}
