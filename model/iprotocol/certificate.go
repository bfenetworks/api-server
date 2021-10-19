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

package iprotocol

import (
	"context"
	"encoding/pem"
	"strings"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/itxn"
	"github.com/bfenetworks/api-server/model/iversion_control"
	"github.com/bfenetworks/bfe/bfe_tls"
)

type Certificate struct {
	CertName    string
	Description string
	IsDefault   bool

	CertFileName string
	CertFilePath string
	KeyFileName  string
	KeyFilePath  string
	ExpiredDate  string

	Products []*ibasic.Product
}

type CertificateFilter struct {
	CertName  *string
	IsDefault *bool
}

type CertificateParam struct {
	CertName    *string
	Description *string
	IsDefault   *bool

	CertFileName    *string
	CertFilePath    *string
	CertFileContent *string
	KeyFileName     *string
	KeyFileContent  *string
	KeyFilePath     *string
	ExpiredDate     *string
}

type CertificateStorager interface {
	FetchCertificates(context.Context, *CertificateFilter) ([]*Certificate, error)
	DeleteCertificate(context.Context, *Certificate) error
	CreateCertificate(context.Context, *CertificateParam) error
	UpdateCertificate(context.Context, *Certificate, *CertificateParam) error
}

type CertificateManager struct {
	storager          CertificateStorager
	txn               itxn.TxnStorager
	extraFileStorager ibasic.ExtraFileStorager

	versionControlManager *iversion_control.VersionControlManager
}

func NewCertificateManager(txn itxn.TxnStorager, storager CertificateStorager,
	versionControlManager *iversion_control.VersionControlManager,
	extraFileStorager ibasic.ExtraFileStorager) *CertificateManager {
	return &CertificateManager{
		txn:               txn,
		storager:          storager,
		extraFileStorager: extraFileStorager,

		versionControlManager: versionControlManager,
	}
}

func (pm *CertificateManager) FetchCertificates(ctx context.Context, param *CertificateFilter) (list []*Certificate, err error) {
	err = pm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err = pm.storager.FetchCertificates(ctx, param)
		return err
	})

	return
}

func (pm *CertificateManager) DeleteCertificate(ctx context.Context, certificate *Certificate) (err error) {
	if certificate.IsDefault {
		return xerror.WrapModelErrorWithMsg("Cant Delete Default Certificate")
	}

	if len(certificate.Products) > 0 {
		return xerror.WrapModelErrorWithMsg("Cant Delete Certificate Be Refer By Product")
	}

	return pm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		if err := pm.extraFileStorager.DeleteExtraFile(ctx, &ibasic.ExtraFileFilter{
			Names: []string{certificate.CertFilePath, certificate.KeyFilePath},
		}); err != nil {
			return err
		}

		return pm.storager.DeleteCertificate(ctx, certificate)
	})
}

func validateCertPair(certFileName string, certFileContent string, keyFileName string, keyFileContent string) error {
	checkCertFileInfo := func(certPEMBlock []byte) error {
		var certDERBlock *pem.Block
		for {
			certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
			if certDERBlock == nil {
				break
			}

			if certDERBlock.Type != "CERTIFICATE" {
				return xerror.WrapParamErrorWithMsg("Certificate File Format Must Be PEM")
			}
		}
		return nil
	}

	if err := checkCertFileInfo([]byte(certFileContent)); err != nil {
		return err
	}

	checkKeyFileInfo := func(keyPEMBlock []byte) error {
		var keyDERBlock *pem.Block
		keyDERBlock, _ = pem.Decode(keyPEMBlock)
		if keyDERBlock == nil {
			return xerror.WrapParamErrorWithMsg("Certificate Private Key File Format Must Be PEM")
		}
		if keyDERBlock.Type != "PRIVATE KEY" && !strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY") {
			return xerror.WrapParamErrorWithMsg("Certificate Private Key File Format Must Be PEM")
		}
		return nil
	}

	if err := checkKeyFileInfo([]byte(keyFileContent)); err != nil {
		return err
	}

	checkCertKeyPair := func(certPEMBlock, keyPEMBlock []byte) error {
		_, err := bfe_tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		return err
	}

	if err := checkCertKeyPair([]byte(certFileContent), []byte(keyFileContent)); err != nil {
		return err
	}

	return nil
}

func (pm *CertificateManager) CreateCertificate(ctx context.Context, param *CertificateParam) (err error) {
	if err = validateCertPair(*param.CertFileName, *param.CertFileContent, *param.KeyFileName, *param.KeyFileContent); err != nil {
		return err
	}

	names := []string{
		*param.CertFileName,
		*param.KeyFileName,
	}

	existedNames := map[string]bool{}
	for _, one := range names {
		if existedNames[one] {
			return xerror.WrapParamErrorWithMsg("Certificate File Name %s Existed", one)
		}

		existedNames[one] = true
	}

	err = pm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err := pm.storager.FetchCertificates(ctx, nil)
		if err != nil {
			return err
		}

		var defaultCertificate *Certificate
		// validate name unique
		for _, one := range list {
			if one.CertName == *param.CertName {
				return xerror.WrapRecordExisted("Certification")
			}

			for _, name := range names {
				if one.CertFileName == name || one.KeyFileName == name {
					return xerror.WrapModelErrorWithMsg("Certificate File Name %s By Used By %s", name, one.CertName)
				}
			}

			if one.IsDefault {
				defaultCertificate = one
			}
		}

		if defaultCertificate == nil && (param.IsDefault == nil || !*param.IsDefault) {
			return xerror.WrapModelErrorWithMsg("Must Has Default Certificate")
		}
		// check/update default
		if isDefault := param.IsDefault; isDefault != nil && *isDefault && defaultCertificate != nil {
			if err = pm.storager.UpdateCertificate(ctx, defaultCertificate, &CertificateParam{
				IsDefault: lib.PBool(false),
			}); err != nil {
				return err
			}
		}

		param.CertFilePath = lib.PString(ibasic.ExtraFilePath(tlsConfDir, ibasic.BuildinProduct, *param.CertFileName))
		param.KeyFilePath = lib.PString(ibasic.ExtraFilePath(tlsConfDir, ibasic.BuildinProduct, *param.KeyFileName))

		if err := pm.extraFileStorager.CreateExtraFile(ctx, ibasic.BuildinProduct, &ibasic.ExtraFileParam{
			Name:    param.CertFilePath,
			Content: []byte(*param.CertFileContent),
		}, &ibasic.ExtraFileParam{
			Name:    param.KeyFilePath,
			Content: []byte(*param.KeyFileContent),
		}); err != nil {
			return err
		}

		return pm.storager.CreateCertificate(ctx, param)
	})

	return
}

func (pm *CertificateManager) UpdateAsDefaultCertificate(ctx context.Context, cert *Certificate) (err error) {
	if cert.IsDefault {
		return nil
	}

	return pm.txn.AtomExecute(ctx, func(ctx context.Context) error {

		list, err := pm.storager.FetchCertificates(ctx, &CertificateFilter{
			IsDefault: lib.PBool(true),
		})
		if err != nil {
			return err
		}
		if len(list) != 0 {
			if err = pm.storager.UpdateCertificate(ctx, list[0], &CertificateParam{
				IsDefault: lib.PBool(false),
			}); err != nil {
				return err
			}
		}

		return pm.storager.UpdateCertificate(ctx, cert, &CertificateParam{
			IsDefault: lib.PBool(true),
		})
	})
}
