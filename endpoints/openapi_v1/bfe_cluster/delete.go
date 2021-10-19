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

package bfe_cluster

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/stateful/container"
)

// BFEClusterDeleteParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type BFEClusterDeleteParam struct {
	Name *string `json:"name" uri:"name" validate:"required,min=1"`
}

// DeleteRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var DeleteEndpoint = &xreq.Endpoint{
	Path:       "/bfe-clusters/{name}",
	Method:     http.MethodDelete,
	Handler:    xreq.Convert(DeleteAction),
	Authorizer: iauth.FA(iauth.FeatureBFECluster, iauth.ActionDelete),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newBFEClusterDeleteParam4Delete(req *http.Request) (*BFEClusterDeleteParam, error) {
	bfeClusterDeleteParam := &BFEClusterDeleteParam{}
	err := xreq.BindURI(req, bfeClusterDeleteParam)
	return bfeClusterDeleteParam, err
}

func deleteActionProcess(req *http.Request, param *BFEClusterDeleteParam) error {
	err := container.BFEClusterManager.DeleteBFECluster(req.Context(), &ibasic.BFEClusterParam{
		Name: param.Name,
	})

	return err
}

var _ xreq.Handler = DeleteAction

// DeleteAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func DeleteAction(req *http.Request) (interface{}, error) {
	param, err := newBFEClusterDeleteParam4Delete(req)
	if err != nil {
		return nil, err
	}

	return nil, deleteActionProcess(req, param)
}
