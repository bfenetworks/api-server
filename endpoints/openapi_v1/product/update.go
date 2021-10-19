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

package product

import (
	"net/http"

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/stateful/container"
)

var _ xreq.Handler = ProductUpdateAction

// ProductUpdateRoute
// Auto Gen By ctrl, Modify as U Need
var ProductUpdateEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(ProductUpdateAction),
	Authorizer: iauth.FAP(iauth.FeatureProduct, iauth.ActionUpdate),
}

// ProductUpdateParam
// Auto Gen By ctrl, Modify as U Need
type ProductUpdateParam struct {
	// Name              *string  `uri:"name" validate:"required,min=2"`
	Description       *string  `json:"description"`
	MailList          []string `json:"mail_list"`
	PhoneList         []string `json:"phone_list"`
	ContactPersonList []string `json:"contact_person_list"`
}

// ProductUpdateAction 更新
// Auto Gen By ctrl, Modify as U Need
func ProductUpdateAction(req *http.Request) (interface{}, error) {
	param := &ProductUpdateParam{}
	if err := xreq.BindJSON(req, param); err != nil {
		return nil, err
	}

	oldOne, err := getProduct(req)
	if err != nil {
		return nil, err
	}

	err = container.ProductManager.UpdateProduct(req.Context(), oldOne, &ibasic.ProductParam{
		Description:       param.Description,
		MailList:          param.MailList,
		PhoneList:         param.PhoneList,
		ContactPersonList: param.ContactPersonList,
	})

	if err != nil {
		return nil, err
	}

	list, err := container.ProductManager.FetchProducts(req.Context(), &ibasic.ProductFilter{
		Name: &oldOne.Name,
	})
	if err != nil {
		return nil, err
	}

	return newProductData(list[0]), nil
}
