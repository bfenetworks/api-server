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
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/baidu/go-lib/web-monitor/web_monitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	MetricAPICostHisCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_cost",
	}, []string{"pattern", "status_code", "method"})
	MetricAPIAccessCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_access",
	}, []string{"pattern", "status_code", "method"})
	MetricSQLCostCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "sql_cost",
	}, []string{"sql"})
	MetricSQLAccessCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "sql_access",
	}, []string{"sql"})

	MetricPaincCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "panic",
	})
)

func init() {
	prometheus.MustRegister(
		MetricAPIAccessCounter,
		MetricAPICostHisCounter,
		MetricSQLAccessCounter,
		MetricSQLCostCounter,
		MetricPaincCounter)
}

func NewMonitorServerWithRun(version string, port int) *web_monitor.MonitorServer {
	monitorServer := web_monitor.NewMonitorServer("BFE_API_SERVER", version, port)

	monitorServer.RegisterHandler(web_monitor.WebHandleMonitor, "metrics", func(p url.Values) ([]byte, error) {
		rsp, req := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
		promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{},
		).ServeHTTP(rsp, req)

		return []byte(rsp.Body.String()), nil
	})

	go monitorServer.Start()

	return monitorServer
}
