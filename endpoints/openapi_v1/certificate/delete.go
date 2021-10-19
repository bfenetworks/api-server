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

// OneParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type OneParam struct {
	CertName *string `uri:"cert_name" validate:"required,min=2"`
}

// OneData Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type OneData struct {
	CertName    string `json:"cert_name" uri:"cert_name"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`

	CertFileName string `json:"cert_file_name"`
	KeyFileName  string `json:"key_file_name"`
	ExpiredDate  string `json:"expired_date"`
}

func newOneData(param *iprotocol.Certificate) *OneData {
	if param == nil {
		return nil
	}

	return &OneData{
		CertName:     param.CertName,
		Description:  param.Description,
		IsDefault:    param.IsDefault,
		CertFileName: param.CertFileName,
		KeyFileName:  param.KeyFileName,
		ExpiredDate:  param.ExpiredDate,
	}
}

// DeleteRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var DeleteEndpoint = &xreq.Endpoint{
	Path:       "/certificates/{cert_name}",
	Method:     http.MethodDelete,
	Handler:    xreq.Convert(DeleteAction),
	Authorizer: iauth.FA(iauth.FeatureCert, iauth.ActionDelete),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newOneParamFromReq(req *http.Request) (*OneParam, error) {
	param := &OneParam{}
	err := xreq.BindURI(req, param)
	return param, err
}

func deleteActionProcess(req *http.Request, param *OneParam) (*OneData, error) {
	list, err := container.CertificateManager.FetchCertificates(req.Context(), &iprotocol.CertificateFilter{
		CertName: param.CertName,
	})
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, xerror.WrapRecordNotExist()
	}

	if err = container.CertificateManager.DeleteCertificate(req.Context(), list[0]); err != nil {
		return nil, err
	}

	return newOneData(list[0]), nil
}

var _ xreq.Handler = DeleteAction

// DeleteAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func DeleteAction(req *http.Request) (interface{}, error) {
	param, err := newOneParamFromReq(req)
	if err != nil {
		return nil, err
	}

	return deleteActionProcess(req, param)
}
