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
	"os"
	"time"

	"github.com/bfenetworks/api-server/lib"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

type ServerConfig struct {
	ServerPort        int `validate:"required,min=1"` // service port
	MonitorPort       int // monitor port
	GracefulTimeOutMs int `validate:"required,min=1"` // time out setting for graceful shutdown
}

type RunTimeConfig struct {
	SessionExpireDay  int  `validate:"required,min=1"`
	SkipTokenValidate bool // skip user identify, you can open it when debug
	RecordSQL         bool
	StaticFilePath    string
	Debug             bool
}

type Config struct {
	Server    ServerConfig
	Loggers   map[string]*LoggerConfig `validate:"dive"`
	Databases map[string]*DbConfig     `validate:"dive"`
	Depends   DependsConfig
	RunTime   RunTimeConfig

	Vars      map[string]string
	LogDir    string
	ConfigDir string
}

var DefaultConfig *Config

func LoadConfig(file string) error {
	config := &Config{
		Depends: DependsConfig{
			NavTreeFile: "${conf_dir}/nav_tree.toml",
			I18nDir:     "${conf_dir}/i18n",
		},
		RunTime: RunTimeConfig{
			StaticFilePath: "./static",
		},
		Vars: map[string]string{},
		Databases: map[string]*DbConfig{
			"bfe_db": {
				Config: mysql.Config{
					Loc: time.Local,
				},
			},
		},
	}

	if err := lib.LoadConfAuto(file, config); err != nil {
		return err
	}

	if err := validator.New().Struct(config); err != nil {
		return err
	}

	DefaultConfig = config
	return nil
}

func (config *Config) Init() error {
	mapping := func(k string) string {
		return config.Vars[k]
	}

	config.Depends.NavTreeFile = os.Expand(config.Depends.NavTreeFile, mapping)
	config.Depends.I18nDir = os.Expand(config.Depends.I18nDir, mapping)

	return nil
}

var (
	BFEProductID int64 = 1

	IgnoreBNSStatusCheck bool
)
