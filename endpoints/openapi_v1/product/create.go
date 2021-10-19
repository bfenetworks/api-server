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
	"context"
	"net/http"

	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/stateful/container"
)

var _ xreq.Handler = ProductCreateAction

// ProductCreateRoute route
// Auto Gen By ctrl, Modify as U Need
var ProductCreateEndpoint = &xreq.Endpoint{
	Path:       "/products",
	Method:     http.MethodPost,
	Handler:    xreq.Convert(ProductCreateAction),
	Authorizer: iauth.FA(iauth.FeatureProduct, iauth.ActionCreate),
}

// ProductCreateParam Param
// Auto Gen By ctrl, Modify as U Need
type ProductCreateParam struct {
	Name              string   `json:"name" validate:"required,min=2"`
	Description       string   `json:"description"`
	MailList          []string `json:"mail_list"`
	PhoneList         []string `json:"phone_list"`
	ContactPersonList []string `json:"contact_person_list"`
}

// ProductCreateAction add one
// Auto Gen By ctrl, Modify as U Need
func ProductCreateAction(req *http.Request) (interface{}, error) {
	param := &ProductCreateParam{}
	if err := xreq.BindJSON(req, param); err != nil {
		return nil, err
	}

	newOne, err := ProductCreateProcess(req.Context(), param)
	if err != nil {
		return nil, err
	}

	return newProductData(newOne), nil
}

func ProductCreateProcess(ctx context.Context, param *ProductCreateParam) (*ibasic.Product, error) {
	if err := xreq.ValidateData(param, nil); err != nil {
		return nil, err
	}

	err := container.ProductManager.CreateProduct(ctx, &ibasic.ProductParam{
		Name:              &param.Name,
		Description:       &param.Description,
		MailList:          param.MailList,
		ContactPersonList: param.ContactPersonList,
		PhoneList:         param.PhoneList,
	})

	if err != nil {
		return nil, err
	}

	list, err := container.ProductManager.FetchProducts(ctx, &ibasic.ProductFilter{
		Name: &param.Name,
	})
	if err != nil {
		return nil, err
	}

	return list[0], nil
}
