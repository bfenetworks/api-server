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

package route

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/stateful/container"
)

// Deprecated
var ExpressionVerifyEndpoint = &xreq.Endpoint{
	Path:    "/expression/verify",
	Method:  http.MethodPatch,
	Handler: xreq.Convert(ExpressionVerifyAction),

	Authorizer: nil,
}

type ExpressionVerifyParam struct {
	Expression string `json:"expression" validate:"required,min=1"`
}

type VerifyResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func newExpressionVerifyParam(req *http.Request) (*ExpressionVerifyParam, error) {
	pfr := &ExpressionVerifyParam{}
	if err := xreq.BindJSON(req, pfr); err != nil {
		return nil, err
	}

	return pfr, nil
}

func ExpressionVerifyActionProcess(req *http.Request, param *ExpressionVerifyParam) (*VerifyResult, error) {
	err := container.RouteRuleManager.ExpressionVerify(req.Context(), param.Expression)
	if err != nil {
		return &VerifyResult{
			Code:    500,
			Message: err.Error(),
		}, err
	}

	return nil, nil
}

var _ xreq.Handler = ExpressionVerifyAction

// ExpressionVerifyAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ExpressionVerifyAction(req *http.Request) (interface{}, error) {
	rule, err := newExpressionVerifyParam(req)
	if err != nil {
		return nil, err
	}

	return ExpressionVerifyActionProcess(req, rule)
}
