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

package version_control

import (
	"context"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/model/iversion_control"
	"github.com/bfenetworks/api-server/storage/rdb/internal/dao"
)

var _ iversion_control.VersionControlStorager = &VersionControlStorager{}

type VersionControlStorager struct {
	dbCtxFactory lib.DBContextFactory
}

func NewVersionControllerStorage(dbCtxFactory lib.DBContextFactory) *VersionControlStorager {
	return &VersionControlStorager{
		dbCtxFactory: dbCtxFactory,
	}
}

func (vcs *VersionControlStorager) UpsertConfigLastExportedVersion(ctx context.Context, css *iversion_control.ExportData) (string, error) {
	dbCtx, err := vcs.dbCtxFactory(ctx)
	if err != nil {
		return "", err
	}

	lastestConfigVersion, err := dao.TConfigVersionOne(dbCtx, &dao.TConfigVersionParam{
		Name:    &css.Topic,
		OrderBy: lib.PString("version DESC"),
	})
	if err != nil {
		return "", err
	}

	if lastestConfigVersion != nil && lastestConfigVersion.DataSign == css.DataSignWithoutVersion {
		return lastestConfigVersion.Version, nil
	}

	version, err := css.CalculateVersion()
	if err != nil {
		return "", err
	}

	_, err = dao.TConfigVersionCreate(dbCtx, &dao.TConfigVersionParam{
		Name:     &css.Topic,
		DataSign: &css.DataSignWithoutVersion,
		Version:  &version,
	})
	return version, err
}
