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

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/iprotocol"
	"github.com/bfenetworks/api-server/stateful/container"
)

// CreateRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var CreateEndpoint = &xreq.Endpoint{
	Path:       "/certificates",
	Method:     http.MethodPost,
	Handler:    xreq.Convert(CreateAction),
	Authorizer: iauth.FA(iauth.FeatureCert, iauth.ActionCreate),
}

// CreateParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type CreateParam struct {
	Name *string `json:"cert_name" validate:"required,min=2"`

	Description *string `json:"description" validate:"required,min=2"`
	IsDefault   *bool   `json:"is_default" validate:"required"`

	CertFileName    *string `json:"cert_file_name" validate:"required,min=2"`
	CertFileContent *string `json:"cert_file_content" validate:"required,min=2"`
	KeyFileName     *string `json:"key_file_name" validate:"required,min=2"`
	KeyFileContent  *string `json:"key_file_content" validate:"required,min=2"`
	ExpiredDate     *string `json:"expired_date" validate:"required,min=2"`
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newCreateParam4Create(req *http.Request) (*CreateParam, error) {
	param := &CreateParam{}
	err := xreq.BindJSON(req, param)
	return param, err
}

func createActionProcess(req *http.Request, param *CreateParam) (*OneData, error) {
	if err := container.CertificateManager.CreateCertificate(req.Context(), &iprotocol.CertificateParam{
		CertName:    param.Name,
		Description: param.Description,
		IsDefault:   param.IsDefault,

		CertFileName:    param.CertFileName,
		CertFileContent: param.CertFileContent,
		KeyFileName:     param.KeyFileName,
		KeyFileContent:  param.KeyFileContent,
		ExpiredDate:     param.ExpiredDate,
	}); err != nil {
		return nil, err
	}

	list, err := container.CertificateManager.FetchCertificates(req.Context(), &iprotocol.CertificateFilter{
		CertName: param.Name,
	})
	if err != nil {
		return nil, err
	}

	return newOneData(list[0]), nil
}

var _ xreq.Handler = CreateAction

// CreateAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func CreateAction(req *http.Request) (interface{}, error) {
	param, err := newCreateParam4Create(req)
	if err != nil {
		return nil, err
	}

	return createActionProcess(req, param)
}
