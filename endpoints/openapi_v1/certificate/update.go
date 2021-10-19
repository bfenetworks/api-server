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

package certificate

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/iprotocol"
	"github.com/bfenetworks/api-server/stateful/container"
)

// UpdateParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type UpdateParam struct {
	CertName *string `uri:"cert_name" validate:"required,min=2"`
}

// UpdateRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UpdateEndpoint = &xreq.Endpoint{
	Path:       "/certificates/{cert_name}/default",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(UpdateAction),
	Authorizer: iauth.FA(iauth.FeatureCert, iauth.ActionUpdate),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newUpdateParam4Update(req *http.Request) (*UpdateParam, error) {
	param := &UpdateParam{}
	err := xreq.BindURI(req, param)

	return param, err
}

func updateActionProcess(req *http.Request, param *UpdateParam) (*OneData, error) {
	list, err := container.CertificateManager.FetchCertificates(req.Context(), &iprotocol.CertificateFilter{
		CertName: param.CertName,
	})
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, xerror.WrapRecordNotExist()
	}

	if err := container.CertificateManager.UpdateAsDefaultCertificate(req.Context(), list[0]); err != nil {
		return nil, err
	}

	list, err = container.CertificateManager.FetchCertificates(req.Context(), &iprotocol.CertificateFilter{
		CertName: param.CertName,
	})
	if err != nil {
		return nil, err
	}

	return newOneData(list[0]), nil
}

var _ xreq.Handler = UpdateAction

// UpdateAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func UpdateAction(req *http.Request) (interface{}, error) {
	param, err := newUpdateParam4Update(req)
	if err != nil {
		return nil, err
	}

	return updateActionProcess(req, param)
}
