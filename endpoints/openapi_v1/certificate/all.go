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
	"github.com/bfenetworks/api-server/stateful/container"
)

// AllRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var AllEndpoint = &xreq.Endpoint{
	Path:       "/certificates",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(AllAction),
	Authorizer: iauth.FA(iauth.FeatureCert, iauth.ActionReadAll),
}

func allActionProcess(req *http.Request) ([]*OneData, error) {
	list, err := container.CertificateManager.FetchCertificates(req.Context(), nil)
	if err != nil {
		return nil, err
	}

	result := []*OneData{}
	for _, one := range list {
		result = append(result, newOneData(one))
	}
	return result, nil
}

var _ xreq.Handler = AllAction

// AllAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func AllAction(req *http.Request) (interface{}, error) {
	return allActionProcess(req)
}
