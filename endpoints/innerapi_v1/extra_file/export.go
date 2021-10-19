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

package extra_file

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/stateful/container"
)

// ExportRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ExportExtraFileEndpoint = &xreq.Endpoint{
	RegisterHandler: func(router *mux.Router) *mux.Route {
		return router.PathPrefix("/configs/extra_files/").Methods(http.MethodGet)
	},
	Handler:    xreq.RawConvert(ExportExtraFileAction),
	Authorizer: iauth.FA(iauth.FeatureExtraFile, iauth.ActionExport),
}

func ExportExtraFileActionProcess(req *http.Request, fileName string) ([]byte, error) {
	extraFile, err := container.ExtraFileManager.FetchExtraFile(req.Context(), fileName)
	if err != nil {
		return nil, err
	}

	if extraFile == nil {
		return nil, xerror.WrapRecordNotExist()
	}

	return extraFile.Content, nil
}

var _ xreq.Handler = ExportExtraFileAction

// ExportExtraFileAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ExportExtraFileAction(req *http.Request) (interface{}, error) {
	return ExportExtraFileActionProcess(req, strings.SplitN(req.URL.Path, "extra_files/", 2)[1])
}
