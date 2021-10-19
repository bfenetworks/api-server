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

package protocol

import (
	"context"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/model/iprotocol"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao"
)

type RDBCertificateStorager struct {
	dbCtxFactory lib.DBContextFactory
}

var _ iprotocol.CertificateStorager = &RDBCertificateStorager{}

func NewCertificateStorager(dbCtxFactory lib.DBContextFactory) *RDBCertificateStorager {
	return &RDBCertificateStorager{
		dbCtxFactory: dbCtxFactory,
	}
}

func (ps *RDBCertificateStorager) DeleteCertificate(ctx context.Context, certificate *iprotocol.Certificate) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TCertificateDelete(dbCtx, &dao.TCertificateParam{
		CertName: &certificate.CertName,
	})
	return err
}

func (ps *RDBCertificateStorager) CreateCertificate(ctx context.Context, pp *iprotocol.CertificateParam) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TCertificateCreate(dbCtx, certificateParamI2D(pp))
	return err
}

func (ps *RDBCertificateStorager) UpdateCertificate(ctx context.Context, Certificate *iprotocol.Certificate,
	pp *iprotocol.CertificateParam) error {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TCertificateUpdate(dbCtx, certificateParamI2D(pp), &dao.TCertificateParam{
		CertName: &Certificate.CertName,
	})
	return err
}

func (ps *RDBCertificateStorager) FetchCertificates(ctx context.Context, filter *iprotocol.CertificateFilter) ([]*iprotocol.Certificate, error) {
	dbCtx, err := ps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	list, err := dao.TCertificateList(dbCtx, filte2param(filter))
	if err != nil {
		return nil, err
	}

	rst := make([]*iprotocol.Certificate, len(list))
	for i, one := range list {
		rst[i] = certificated2i(one)
	}
	return rst, nil
}

func filte2param(filter *iprotocol.CertificateFilter) *dao.TCertificateParam {
	if filter == nil {
		return nil
	}

	return &dao.TCertificateParam{
		CertName:  filter.CertName,
		IsDefault: filter.IsDefault,
	}
}

func certificateParamI2D(pp *iprotocol.CertificateParam) *dao.TCertificateParam {
	if pp == nil {
		return nil
	}

	return &dao.TCertificateParam{
		CertName:     pp.CertName,
		Description:  pp.Description,
		IsDefault:    pp.IsDefault,
		CertFileName: pp.CertFileName,
		CertFilePath: pp.CertFilePath,
		KeyFileName:  pp.KeyFileName,
		KeyFilePath:  pp.KeyFilePath,
		ExpiredDate:  pp.ExpiredDate,
	}
}

func certificated2i(pp *dao.TCertificate) *iprotocol.Certificate {
	return &iprotocol.Certificate{
		CertName:     pp.CertName,
		Description:  pp.Description,
		IsDefault:    pp.IsDefault,
		CertFileName: pp.CertFileName,
		CertFilePath: pp.CertFilePath,
		KeyFileName:  pp.KeyFileName,
		KeyFilePath:  pp.KeyFilePath,
		ExpiredDate:  pp.ExpiredDate,
	}
}
