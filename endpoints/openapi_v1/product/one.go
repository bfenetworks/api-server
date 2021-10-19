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
	"strings"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
)

var _ xreq.Handler = ProductOneAction

// ProductOneRoute route
// Auto Gen By ctrl, Modify as U Need
var ProductOneEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ProductOneAction),
	Authorizer: iauth.FAP(iauth.FeatureProduct, iauth.ActionRead),
}

// ProductData one response
// Auto Gen By ctrl, Modify as U Need
type ProductData struct {
	Id                int64    `json:"id,omitempty"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	MailList          []string `json:"mail_list"`
	PhoneList         []string `json:"phone_list"`
	ContactPersonList []string `json:"contact_person_list"`
}

// ProductOneAction get one
// Auto Gen By ctrl, Modify as U Need
func ProductOneAction(req *http.Request) (interface{}, error) {
	_p, err := getProduct(req)
	if err != nil {
		return nil, err
	}

	return newProductData(_p), nil
}

func getProduct(req *http.Request) (*ibasic.Product, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	if strings.ToUpper(product.Name) == ibasic.BuildinProduct.Name {
		return nil, xerror.WrapModelErrorWithMsg("Dont Modify Buid-in Product")
	}

	return product, nil
}
