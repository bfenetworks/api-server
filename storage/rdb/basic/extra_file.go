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
	"crypto/md5"
	"fmt"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao"
)

type RDBExtraFileStorager struct {
	dbCtxFactory lib.DBContextFactory
}

var _ ibasic.ExtraFileStorager = &RDBExtraFileStorager{}

func NewRDBExtraFileStorager(dbCtxFactory lib.DBContextFactory) *RDBExtraFileStorager {
	return &RDBExtraFileStorager{
		dbCtxFactory: dbCtxFactory,
	}
}

func (ps *RDBExtraFileStorager) FetchExtraFiles(ctx context.Context, where *ibasic.ExtraFileFilter) ([]*ibasic.ExtraFile, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	list, err := dao.TExtraFileList(dbCtx, filter2Param(where))
	if err != nil {
		return nil, err
	}

	var rst []*ibasic.ExtraFile
	for _, one := range list {
		rst = append(rst, &ibasic.ExtraFile{
			ID:          one.ID,
			ProductID:   one.ProductID,
			Name:        one.Name,
			Description: one.Description,
			Md5:         one.Md5,
			Content:     one.Content,
		})
	}

	return rst, nil
}

func filter2Param(filter *ibasic.ExtraFileFilter) *dao.TExtraFileParam {
	if filter == nil {
		return nil
	}

	return &dao.TExtraFileParam{
		Name:  filter.Name,
		Names: filter.Names,
	}
}

func (ps *RDBExtraFileStorager) DeleteExtraFile(ctx context.Context, pp *ibasic.ExtraFileFilter) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TExtraFileDelete(dbCtx, filter2Param(pp))

	return err
}

func (ps *RDBExtraFileStorager) CreateExtraFile(ctx context.Context, product *ibasic.Product, pps ...*ibasic.ExtraFileParam) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	tps := []*dao.TExtraFileParam{}
	for _, pp := range pps {
		if pp.Md5 == nil {
			pp.Md5 = []byte(fmt.Sprintf("%x", md5.Sum(pp.Content)))
		}

		tps = append(tps, &dao.TExtraFileParam{
			ProductID:   &product.ID,
			Name:        pp.Name,
			Description: pp.Description,
			Md5:         pp.Md5,
			Content:     pp.Content,
		})
	}

	_, err = dao.TExtraFileCreate(dbCtx, tps...)

	return err
}
