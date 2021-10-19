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

package product_cluster

import (
	"net/http"

	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/stateful/container"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
)

// StickySessions Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type StickySessions struct {
	SessionStickyType string `json:"session_sticky_type" uri:"session_sticky_type"`
	HashStrategy      string `json:"hash_strategy" uri:"hash_strategy"`
	HashHeader        string `json:"hash_header" uri:"hash_header"`
}

// Basic Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type Basic struct {
	Connection *Connection `json:"connection" uri:"connection"`
	Retries    *Retries    `json:"retries" uri:"retries"`
	Buffers    *Buffers    `json:"buffers" uri:"buffers"`
	Timeouts   *Timeouts   `json:"timeouts" uri:"timeouts"`
}

// Connection Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type Connection struct {
	MaxIdleConnPerRs    int16 `json:"max_idle_conn_per_rs" uri:"max_idle_conn_per_rs"`
	CancelOnClientClose bool  `json:"cancel_on_client_close" uri:"cancel_on_client_close"`
}

// Retries Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type Retries struct {
	MaxRetryInSubcluster    int8 `json:"max_retry_in_subcluster" uri:"max_retry_in_subcluster"`
	MaxRetryCrossSubcluster int8 `json:"max_retry_cross_subcluster" uri:"max_retry_cross_subcluster"`
}

// Buffers Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type Buffers struct {
	ReqWriteBufferSize int32 `json:"req_write_buffer_size" uri:"req_write_buffer_size"`
}

// Timeouts Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type Timeouts struct {
	TimeoutConnServ        int32 `json:"timeout_conn_serv" uri:"timeout_conn_serv"`
	TimeoutResponseHeader  int32 `json:"timeout_response_header" uri:"timeout_response_header"`
	TimeoutReadbodyClient  int32 `json:"timeout_readbody_client" uri:"timeout_readbody_client"`
	TimeoutReadClientAgain int32 `json:"timeout_read_client_again" uri:"timeout_read_client_again"`
	TimeoutWriteClient     int32 `json:"timeout_write_client" uri:"timeout_write_client"`
}

// OneParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type OneParam struct {
	Name *string `uri:"cluster_name" validate:"required,min=2"`
}

type PassiveHealthCheck struct {
	Schema     string `json:"schema"`
	Interval   int32  `json:"interval" validate:"required,min=1"`
	Failnum    int32  `json:"failnum" validate:"required,min=1"`
	Statuscode int32  `json:"statuscode" validate:"required,min=0"`
	Host       string `json:"host" validate:"required,min=1"`
	Uri        string `json:"uri" validate:"required,min=1,startswith=/"`
}

// ClusterData Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type ClusterData struct {
	Name           string          `json:"name" uri:"name"`
	Description    string          `json:"description" uri:"description"`
	Basic          *Basic          `json:"basic" uri:"basic"`
	StickySessions *StickySessions `json:"sticky_sessions" uri:"sticky_sessions"`
	Ready          bool            `json:"ready"`

	PassiveHealthCheck *PassiveHealthCheck `json:"passive_health_check"`

	SubClusters []string `json:"sub_clusters"`

	ManualScheduler map[string]map[string]int `json:"manual_scheduler,omitempty"`
}

type AutoLbMatrix struct {
	MaxRegionLoad    float64          `json:"max_region_load"`
	MaxBlackholeLoad float64          `json:"max_blackhole_load"`
	BlackholeEnabled bool             `json:"blackhole_enabled"`
	Capacity         map[string]int64 `json:"capacity"`
}

func clusterModel2Control(cluster *icluster_conf.Cluster) *ClusterData {
	rsp := &ClusterData{
		Name:        cluster.Name,
		Description: cluster.Description,
		Ready:       cluster.Ready,
		Basic: &Basic{
			Connection: &Connection{
				MaxIdleConnPerRs:    cluster.Basic.Connection.MaxIdleConnPerRs,
				CancelOnClientClose: cluster.Basic.Connection.CancelOnClientClose,
			},
			Retries: &Retries{
				MaxRetryCrossSubcluster: cluster.Basic.Retries.MaxRetryCrossSubcluster,
				MaxRetryInSubcluster:    cluster.Basic.Retries.MaxRetryInSubcluster,
			},
			Buffers: &Buffers{
				ReqWriteBufferSize: cluster.Basic.Buffers.ReqWriteBufferSize,
			},
			Timeouts: &Timeouts{
				TimeoutConnServ:        cluster.Basic.Timeouts.TimeoutConnServ,
				TimeoutResponseHeader:  cluster.Basic.Timeouts.TimeoutResponseHeader,
				TimeoutReadbodyClient:  cluster.Basic.Timeouts.TimeoutReadbodyClient,
				TimeoutReadClientAgain: cluster.Basic.Timeouts.TimeoutReadClientAgain,
				TimeoutWriteClient:     cluster.Basic.Timeouts.TimeoutWriteClient,
			},
		},
		StickySessions: &StickySessions{
			SessionStickyType: map[bool]string{
				true:  icluster_conf.ClusterStickTypeInstance,
				false: icluster_conf.ClusterStickTypeSubCluster,
			}[cluster.StickySessions.SessionSticky],
			HashStrategy: map[int32]string{
				icluster_conf.ClusterHashStrategyClientIDOnlyI:     clusterHashStrategyClientIDOnly,
				icluster_conf.ClusterHashStrategyClientIPOnlyI:     clusterHashStrategyClientIPOnly,
				icluster_conf.ClusterHashStrategyClientIDPreferedI: clusterHashStrategyClientIDPrefered,
			}[cluster.StickySessions.HashStrategy],
			HashHeader: cluster.StickySessions.HashHeader,
		},

		SubClusters: cluster.SubClusterNames(),

		ManualScheduler: cluster.ManualScheduler,

		PassiveHealthCheck: PassiveHealthCheckM2C(cluster.PassiveHealthCheck),
	}

	return rsp
}

func PassiveHealthCheckM2C(phc *icluster_conf.ClusterPassiveHealthCheck) *PassiveHealthCheck {
	if phc == nil {
		return nil
	}

	return &PassiveHealthCheck{
		Schema:     phc.Schema,
		Interval:   phc.Interval,
		Failnum:    phc.Failnum,
		Statuscode: phc.Statuscode,
		Host:       phc.Host,
		Uri:        phc.Uri,
	}

}

// OneRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var OneEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/clusters/{cluster_name}",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(OneAction),
	Authorizer: iauth.FAP(iauth.FeatureProductCluster, iauth.ActionRead),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newOneParam4One(req *http.Request) (*OneParam, error) {
	param := &OneParam{}
	err := xreq.BindURI(req, param)
	return param, err
}

var _ xreq.Handler = OneAction

// OneAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func OneAction(req *http.Request) (interface{}, error) {
	param, err := newOneParam4One(req)
	if err != nil {
		return nil, err
	}

	return oneActionProcess(req, param)
}

func oneActionProcess(req *http.Request, param *OneParam) (*ClusterData, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	one, err := container.ClusterManager.FetchCluster(req.Context(), &icluster_conf.ClusterFilter{
		Name:    param.Name,
		Product: product,
	})
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, xerror.WrapRecordNotExist("Cluster")
	}

	return clusterModel2Control(one), nil
}
