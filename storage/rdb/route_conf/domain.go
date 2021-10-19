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

package route_conf

import (
	"context"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/iroute_conf"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao"
)

var _ iroute_conf.DomainStorager = &DomainStorager{}

var RDBDomainStoragerSingleton iroute_conf.DomainStorager

func NewDomainStorager(dbCtxFactory lib.DBContextFactory) *DomainStorager {
	return &DomainStorager{
		dbCtxFactory: dbCtxFactory,
	}
}

type DomainStorager struct {
	dbCtxFactory lib.DBContextFactory
}

func (rs *DomainStorager) FetchDomains(ctx context.Context, filter *iroute_conf.DomainFilter) ([]*iroute_conf.Domain, error) {
	dbCtx, err := rs.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	list, err := dao.TDomainList(dbCtx, domainFilter2Param(filter))
	if err != nil {
		return nil, err
	}

	rst := make([]*iroute_conf.Domain, len(list))
	for i, one := range list {
		rst[i] = domaind2i(one)
	}

	return rst, nil
}

func (rs *DomainStorager) CreateDomain(ctx context.Context, product *ibasic.Product, param *iroute_conf.DomainParam) error {
	param.ProductID = &product.ID

	dbCtx, err := rs.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_p := domainParami2d(param)
	_p.Type = &dao.DomainTypeGeneral
	_, err = dao.TDomainCreate(dbCtx, _p)
	return err
}

func (rs *DomainStorager) DeleteDomain(ctx context.Context, product *ibasic.Product, domain *iroute_conf.Domain) error {
	dbCtx, err := rs.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TDomainDelete(dbCtx, &dao.TDomainParam{
		ID: &domain.ID,
	})
	return err
}

func domainFilter2Param(filter *iroute_conf.DomainFilter) *dao.TDomainParam {
	if filter == nil {
		return nil
	}

	var pid *int64
	if filter.Product != nil {
		pid = &filter.Product.ID
	}
	return &dao.TDomainParam{
		ProductID: pid,
		Name:      filter.Name,
	}
}

func domainParami2d(p *iroute_conf.DomainParam) *dao.TDomainParam {
	if p == nil {
		return nil
	}

	return &dao.TDomainParam{
		ProductID:             p.ProductID,
		Name:                  p.Name,
		UsingAdvancedRedirect: p.UsingAdvancedRedirect,
		UsingAdvancedHsts:     p.UsingAdvancedHsts,
	}
}

func domaind2i(p *dao.TDomain) *iroute_conf.Domain {
	if p == nil {
		return nil
	}

	return &iroute_conf.Domain{
		ProductID:             p.ProductID,
		Name:                  p.Name,
		UsingAdvancedRedirect: p.UsingAdvancedRedirect,
		UsingAdvancedHsts:     p.UsingAdvancedHsts,
		ID:                    p.ID,
	}
}
