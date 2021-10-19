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

package container

import (
	"github.com/bfenetworks/api-server/model/iauth"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/icluster_conf"
	"github.com/bfenetworks/api-server/model/iprotocol"
	"github.com/bfenetworks/api-server/model/iroute_conf"
	"github.com/bfenetworks/api-server/model/itxn"
	"github.com/bfenetworks/api-server/model/iversion_control"
)

var (
	TxnStoragerSingleton                itxn.TxnStorager
	VersionControlStoragerSingleton     iversion_control.VersionControlStorager
	RouteRuleStoragerSingleton          iroute_conf.RouteRuleStorager
	ProductStoragerSingleton            ibasic.ProductStorager
	BFEClusterStoragerSingleton         ibasic.BFEClusterStorager
	DomainStoragerSingleton             iroute_conf.DomainStorager
	ClusterStoragerSingleton            icluster_conf.ClusterStorager
	PoolStoragerSingleton               icluster_conf.PoolStorager
	SubClusterStoragerSingleton         icluster_conf.SubClusterStorager
	CertificateStoragerSingleton        iprotocol.CertificateStorager
	AuthenticateStoragerSingleton       iauth.AuthenticateStorager
	ProductAuthorizateStoragerSingleton iauth.ProductAuthorizateStorager
	ExtraFileStoragerSingleton          ibasic.ExtraFileStorager

	ExtraFileManager            *ibasic.ExtraFileManager
	ProductManager              *ibasic.ProductManager
	DomainManager               *iroute_conf.DomainManager
	BFEClusterManager           *ibasic.BFEClusterManager
	VersionControlManager       *iversion_control.VersionControlManager
	RouteRuleManager            *iroute_conf.RouteRuleManager
	ClusterManager              *icluster_conf.ClusterManager
	SubClusterManager           *icluster_conf.SubClusterManager
	CertificateManager          *iprotocol.CertificateManager
	AuthenticateManager         *iauth.AuthenticateManager
	ProductAuthorizateManager   *iauth.ProductAuthorizateManager
	FeatureAuthorizerrSingleton *iauth.FeatureAuthorizer
	PoolManager                 *icluster_conf.PoolManager
)
