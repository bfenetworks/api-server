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
	"errors"
	"fmt"
	"os"

	"github.com/baidu/go-lib/log"
	"github.com/baidu/go-lib/log/log4go"
)

type LoggerConfig struct {
	LogName     string `validate:"required,min=1"`
	LogLevel    string `validate:"required,oneof=DEBUG TRACE INFO WARNING ERROR CRITICAL"`
	RotateWhen  string `validate:"required,oneof=M H D MIDNIGHT"` // rotate time
	BackupCount int    `validate:"required,min=1"`                // backup files
	Format      string
	StdOut      bool
}

func (config *Config) InitLog() error {
	log4go.SetLogBufferLength(10000)
	log4go.SetLogWithBlocking(false)

	var err error
	loggers := map[string]log4go.Logger{}
	for _, c := range config.Loggers {
		defaultFormat := log4go.LogFormat
		if c.Format != "" {
			log4go.SetLogFormat(c.Format)
		}
		loggers[c.LogName], err = log.Create(c.LogName, c.LogLevel, config.LogDir, c.StdOut, c.RotateWhen, c.BackupCount)
		if err != nil {
			return err
		}
		log4go.SetLogFormat(defaultFormat)
	}

	AccessLogger = loggers["access"]
	if AccessLogger == nil {
		return errors.New("logger access must be set")
	}

	SQLLogger = loggers["sql"]
	if SQLLogger == nil && config.RunTime.RecordSQL {
		return errors.New("logger access must be set when RunTime.RecordSQL")
	}

	log.Logger = AccessLogger

	return nil
}

var (
	SQLLogger    log4go.Logger
	AccessLogger log4go.Logger
)

func CloseLog() {
	log.Logger.Close()
}

func Exit(topic string, err error, code int) {
	msg := fmt.Sprintf("%s - %s", topic, err)
	fmt.Println(msg)

	if AccessLogger != nil {
		AccessLogger.Error(msg)
	}

	os.Exit(code)
}
