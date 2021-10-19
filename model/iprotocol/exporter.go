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
	"strings"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/iversion_control"
	"github.com/bfenetworks/bfe/bfe_config/bfe_tls_conf/server_cert_conf"
)

const (
	ConfigTopicServerCert = "certificate"

	tlsConfDir = "tls_conf"
)

type ServerCertConf struct {
	server_cert_conf.BfeServerCertConf
}

func (scc *ServerCertConf) UpdateVersion(version string) error {
	scc.Version = version

	for certFileName, certConfig := range scc.BfeServerCertConf.Config.CertConf {
		i := strings.Index(certConfig.ServerCertFile, "/")
		if i == -1 { // impossible
			return xerror.WrapDirtyDataErrorWithMsg("ServerCertFile must has /, path: %s", certConfig.ServerCertFile)
		}
		certConfig.ServerCertFile = tlsConfDir + "_" + version + certConfig.ServerCertFile[i:]

		i = strings.Index(certConfig.ServerKeyFile, "/")
		if i == -1 { // impossible
			return xerror.WrapDirtyDataErrorWithMsg("ServerKeyFile must has /, path: %s", certConfig.ServerKeyFile)
		}
		certConfig.ServerKeyFile = tlsConfDir + "_" + version + certConfig.ServerKeyFile[i:]

		scc.Config.CertConf[certFileName] = certConfig
	}

	return nil
}

func (pm *CertificateManager) certificateGenerator(ctx context.Context) (*iversion_control.ExportData, error) {
	certificates, err := pm.storager.FetchCertificates(ctx, nil)
	if err != nil {
		return nil, err
	}

	defaultCertName := ""
	certConf := map[string]server_cert_conf.ServerCertConf{}

	for _, cert := range certificates {
		if cert.IsDefault {
			defaultCertName = cert.CertName
		}

		certConf[cert.CertName] = server_cert_conf.ServerCertConf{
			ServerCertFile: cert.CertFilePath,
			ServerKeyFile:  cert.KeyFilePath,
		}
	}

	scc := &ServerCertConf{
		BfeServerCertConf: server_cert_conf.BfeServerCertConf{
			Config: server_cert_conf.ServerCertConfMap{
				Default:  defaultCertName,
				CertConf: certConf,
			},
		},
	}
	scc.UpdateVersion(iversion_control.ZeroVersion)

	return &iversion_control.ExportData{
		Topic:              ConfigTopicServerCert,
		DataWithoutVersion: scc,
	}, nil
}

func (pm *CertificateManager) ExportServerCert(ctx context.Context, lastVersion string) (*ServerCertConf, error) {
	ed, err := pm.versionControlManager.ExportConfig(ctx, ConfigTopicServerCert, pm.certificateGenerator)
	if err != nil {
		return nil, err
	}

	conf := ed.DataWithoutVersion.(*ServerCertConf)
	if conf.Version == lastVersion {
		return nil, nil
	}

	return conf, nil
}
