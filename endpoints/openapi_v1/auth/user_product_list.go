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

package auth

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/stateful/container"
)

// ProductUserListRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var ProductUserBindListEndpoint = &xreq.Endpoint{
	Path:       "/auth/users/actions/search-by-product/{product_name}",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ProductUserListAction),
	Authorizer: iauth.FA(iauth.FeatureProductUser, iauth.ActionReadAll),
}

// ProducctUserListParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type ProducctUserListParam struct {
	Type *string `json:"type" form:"type" validate:"omitempty,oneof=jwt normal"`
}

func productUserListActionProcess(req *http.Request, param *ProducctUserListParam) ([]*UserData, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	list, err := container.AuthorizeManager.FetchProductUsers(req.Context(), product)
	if err != nil {
		return nil, err
	}

	var userDataList []*UserData
	for _, one := range list {
		userDataList = append(userDataList, newUserData(one))
	}

	return userDataList, nil
}

var _ xreq.Handler = ProductUserListAction

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newProducctUserListParam(req *http.Request) (*ProducctUserListParam, error) {
	param := &ProducctUserListParam{}
	err := xreq.BindForm(req, param)
	return param, err
}

// ProductUserListAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func ProductUserListAction(req *http.Request) (interface{}, error) {
	param, err := newProducctUserListParam(req)
	if err != nil {
		return nil, err
	}

	return productUserListActionProcess(req, param)
}
