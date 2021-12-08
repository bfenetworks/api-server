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

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xreq"
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/model/iroute_conf"
	"github.com/bfenetworks/api-server/stateful/container"
)

type ProductListParam struct {
	Domain  *string `form:"domain"`
	Cluster *string `form:"cluster"`
}

func newProductData(pp *ibasic.Product) *ProductData {
	return &ProductData{
		Name:              pp.Name,
		Description:       pp.Description,
		MailList:          pp.MailList,
		PhoneList:         pp.PhoneList,
		ContactPersonList: pp.ContactPersonList,
	}
}

var _ xreq.Handler = ProductListAction

// ProductListRoute route
// Auto Gen By ctrl, Modify as U Need
var ProductListEndpoint = &xreq.Endpoint{
	Path:       "/products",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ProductListAction),
	Authorizer: iauth.FA(iauth.FeatureProduct, iauth.ActionRead),
}

// ProductListAction get list
// Auto Gen By ctrl, Modify as U Need
func ProductListAction(req *http.Request) (interface{}, error) {
	param := &ProductListParam{}
	if err := xreq.BindForm(req, param); err != nil {
		return nil, err
	}

	var productID *int64
	rst := []*ProductData{}
	if domain := param.Domain; domain != nil && *domain != "" {
		domains, err := container.DomainManager.DomainList(req.Context(), &iroute_conf.DomainFilter{
			Name: domain,
		})
		if err != nil {
			return nil, err
		}
		if len(domains) == 0 {
			return rst, nil
		}

		productID = &(domains[0].ProductID)
	} else if cluster := param.Cluster; cluster != nil && *cluster != "" {
		clusterObj, err := container.ClusterManager.FetchCluster(req.Context(), &icluster_conf.ClusterFilter{
			Name: cluster,
		})
		if err != nil {
			return nil, err
		}
		if clusterObj == nil {
			return rst, nil
		}

		productID = &(clusterObj.ProductID)
	}

	visitor, err := iauth.MustGetVisitor(req.Context())
	if err != nil {
		return nil, err
	}

	var grantedProducts []*ibasic.Product
	if !visitor.IsAdmin() {
		grantedProducts, err = container.AuthorizeManager.FetchVisitorProductList(req.Context(), visitor)
		if err != nil {
			return nil, err
		}
		if len(grantedProducts) == 0 {
			return []*ProductData{}, nil
		}

		for _, pp := range grantedProducts {
			if productID == nil || *productID == pp.ID {
				rst = append(rst, newProductData(pp))
			}
		}

		return rst, nil
	}

	list, err := container.ProductManager.FetchProducts(req.Context(), &ibasic.ProductFilter{
		ID:   productID,
		NeID: lib.PInt64(1),
	})
	if err != nil {
		return nil, err
	}

	for _, pp := range list {
		rst = append(rst, newProductData(pp))
	}

	return rst, nil
}
