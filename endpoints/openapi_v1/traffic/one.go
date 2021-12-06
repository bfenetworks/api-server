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

package traffic

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

// OneParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type OneParam struct {
	ClusterName string `json:"cluster_name" uri:"cluster_name" validate:"required,min=2"`
}

// AutoScheduler Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type AutoScheduler struct {
	MaxRegionLoad    float64          `json:"max_region_load" uri:"max_region_load"`
	MaxBlackholeLoad float64          `json:"max_blackhole_load" uri:"max_blackhole_load"`
	BlackholeEnabled bool             `json:"blackhole_enabled" uri:"blackhole_enabled"`
	Capacity         map[string]int64 `json:"capacity"`
}

// SubCluster Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type SubCluster struct {
	Name     string `json:"name" uri:"name" validate:"required,min=1"`
	Capacity int64  `json:"capacity" uri:"capacity" validate:"required,min=1"`
}

// OneData Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type OneData struct {
	Cluster       string                    `json:"cluster" uri:"cluster"`
	Scheduler     map[string]map[string]int `json:"scheduler,omitempty" uri:"scheduler"`
	AutoScheduler *AutoScheduler            `json:"auto_scheduler,omitempty" uri:"auto_scheduler"`
}

// OneRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var OneEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/clusters/{cluster_name}/scheduler",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(OneAction),
	Authorizer: iauth.FAP(iauth.FeatureTraffic, iauth.ActionRead),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newOneParam4One(req *http.Request) (*OneParam, error) {
	param := &OneParam{}
	err := xreq.BindURI(req, param)
	return param, err
}

func oneActionProcess(req *http.Request, param *OneParam) (*OneData, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	cluster, err := container.ClusterManager.FetchCluster(req.Context(), &icluster_conf.ClusterFilter{
		Product: product,
		Name:    &param.ClusterName,
	})
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, xerror.WrapRecordNotExist("Cluster")
	}

	if cluster.Scheduler == nil {
		return nil, xerror.WrapDirtyDataErrorWithMsg("Manual Mode, But Without LbMatrix Setting")
	}

	return &OneData{
		Cluster:   cluster.Name,
		Scheduler: cluster.Scheduler,
	}, nil
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
