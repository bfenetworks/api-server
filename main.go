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

package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"path/filepath"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"gopkg.in/tylerb/graceful.v1"

	"github.com/bfenetworks/api-server/endpoints"
	"github.com/bfenetworks/api-server/stateful"
	"github.com/bfenetworks/api-server/stateful/container/rdb"
	"github.com/bfenetworks/api-server/storage/register"
	"github.com/bfenetworks/api-server/version"
)

var (
	help       *bool   = flag.Bool("h", false, "to show help")
	showVer    *bool   = flag.Bool("v", false, "to show version")
	confDir    *string = flag.String("c", "./conf/", "API configure dir")
	serverConf *string = flag.String("sc", "api_server.toml", "server conf file")

	logDir *string = flag.String("l", "./log", "dir path of log")
)

func main() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}
	if *showVer {
		fmt.Printf("version %s\n", version.Version)
		return
	}

	if err := stateful.LoadConfig(filepath.Join(*confDir, *serverConf)); err != nil {
		stateful.Exit("LoadConfig", err, -1)
	}

	config := stateful.DefaultConfig
	config.LogDir = *logDir
	config.ConfigDir = *confDir

	config.Vars["conf_dir"] = *confDir
	config.Vars["log_dir"] = *logDir

	if err := config.Init(); err != nil {
		stateful.Exit("config.Init", err, -1)
	}

	if err := config.InitLog(); err != nil {
		stateful.Exit("config.InitLog", err, -1)
	}

	defer func() {
		time.Sleep(time.Second)
		stateful.CloseLog()
	}()

	if err := config.Depends.Init(); err != nil {
		stateful.Exit("config.Depends.Init", err, -1)
	}

	if err := config.InitDB(); err != nil {
		stateful.Exit("config.InitDB", err, -1)
	}

	registerConfig, _ := stateful.GetRegisterConfig(confDir)

	registerServier := register.RegisterServier{RegisterConfig: registerConfig}
	registerServier.Init()

	rdb.Init(&registerServier)
	serverStartUp()
}

func serverStartUp() {
	serverConfig := stateful.DefaultConfig.Server

	if serverConfig.MonitorPort > 0 {
		stateful.NewMonitorServerWithRun(version.Version, serverConfig.MonitorPort)
	}

	n := negroni.New()
	router := mux.NewRouter()
	endpoints.RegisterRouters(router)
	n.UseHandler(router)

	timeout := time.Duration(serverConfig.GracefulTimeOutMs) * time.Millisecond
	address := fmt.Sprintf("0.0.0.0:%d", serverConfig.ServerPort)
	fmt.Println("Run Server At:", address)
	graceful.Run(address, timeout, n)
}
