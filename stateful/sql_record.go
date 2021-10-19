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
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bfenetworks/api-server/lib"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
)

type SQLRecord struct {
	Cost time.Duration
	SQL  string
	Args []interface{}
	Err  error
}

var inReg = regexp.MustCompile(`(?i)(?U)in.?\(.*\).?`)

func (sr *SQLRecord) format() {
	// INSERT INTO mod_header_rules (actions,cond,created_at,is_last,product_id) VALUES ('[{"Cmd":"REQ_HEADER_SET",
	if strings.HasPrefix(sr.SQL, "INSERT") {
		if i := strings.Index(sr.SQL, " VALUES "); i != -1 {
			sr.SQL = sr.SQL[0:i] + "intercept"
		}
	}

	// XXX WHERE X in (?, ?, ...)
	sr.SQL = inReg.ReplaceAllString(sr.SQL, "IN(intercept)")
}

func (sr *SQLRecord) UpdateMonitor() {
	MetricSQLAccessCounter.With(prometheus.Labels{"sql": sr.SQL}).Inc()
	MetricSQLCostCounter.With(prometheus.Labels{"sql": sr.SQL}).Add(float64(sr.Cost.Milliseconds()))
}

func (sr *SQLRecord) String(ctx context.Context) string {
	logID := lib.GainLogID(ctx)
	return fmt.Sprintf("[%s]type[SQL] cost_ms[%d] sql[%s] params[%v] err[%v]",
		logID, sr.Cost.Milliseconds(), sr.SQL, sr.Args, sr.Err)
}

func (sr *SQLRecord) Print(ctx context.Context) {
	sr.format()
	sr.UpdateMonitor()

	if !DefaultConfig.RunTime.RecordSQL {
		return
	}
	SQLLogger.Info(sr.String(ctx))
}
