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

package basic

import (
	"context"
	"strings"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao"
)

type RDBProductStorager struct {
	dbCtxFactory lib.DBContextFactory
}

var _ ibasic.ProductStorager = &RDBProductStorager{}

func NewProductManager(dbCtxFactory lib.DBContextFactory) *RDBProductStorager {
	return &RDBProductStorager{
		dbCtxFactory: dbCtxFactory,
	}
}

func (ps *RDBProductStorager) DeleteProduct(ctx context.Context, product *ibasic.Product) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	return dao.TProductDeleteByProductID(dbCtx, product.ID)
}

func (ps *RDBProductStorager) CreateProduct(ctx context.Context, pp *ibasic.ProductParam) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TProductCreate(dbCtx, productParami2d(pp))
	return err
}

func (ps *RDBProductStorager) UpdateProduct(ctx context.Context, product *ibasic.Product,
	pp *ibasic.ProductParam) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TProductUpdate(dbCtx, productParami2d(pp), &dao.TProductParam{
		Name: &product.Name,
	})
	return err
}

func (ps *RDBProductStorager) FetchProducts(ctx context.Context, filter *ibasic.ProductFilter) ([]*ibasic.Product, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	list, err := dao.TProductList(dbCtx, productFilter2Param(filter))
	if err != nil {
		return nil, err
	}

	rst := make([]*ibasic.Product, len(list))
	for i, one := range list {
		rst[i] = productd2i(one)
	}
	return rst, nil
}

func productFilter2Param(filter *ibasic.ProductFilter) *dao.TProductParam {
	if filter == nil {
		return nil
	}

	return &dao.TProductParam{
		IDs:  filter.IDs,
		Id:   filter.ID,
		NeId: filter.NeID,
		Name: filter.Name,
	}
}

func productParami2d(pp *ibasic.ProductParam) *dao.TProductParam {
	if pp == nil {
		return nil
	}

	marshal := func(s []string) *string {
		if s == nil {
			return nil
		}

		ss := strings.Join(s, ";")
		return &ss
	}

	return &dao.TProductParam{
		Id:          pp.ID,
		Name:        pp.Name,
		Description: pp.Description,

		MailList:      marshal(pp.MailList),
		ContactPerson: marshal(pp.ContactPersonList),
		SmsList:       marshal(pp.PhoneList),
	}
}

func productd2i(p *dao.TProduct) *ibasic.Product {
	unmarshal := func(s string) []string {
		if s == "" {
			return []string{}
		}

		return strings.Split(s, ";")
	}

	return &ibasic.Product{
		ID:                p.Id,
		Name:              p.Name,
		Description:       p.Description,
		MailList:          unmarshal(p.MailList),
		ContactPersonList: unmarshal(p.ContactPerson),
		PhoneList:         unmarshal(p.SmsList),

		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
