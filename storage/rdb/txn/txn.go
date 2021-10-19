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

package txn

import (
	"context"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/model/itxn"
)

type RDBTxnStorager struct {
	dbCtxFactory lib.DBContextFactory
}

func NewRDBTxnStorager(dbCtxFactory lib.DBContextFactory) *RDBTxnStorager {
	return &RDBTxnStorager{
		dbCtxFactory: dbCtxFactory,
	}
}

var _ itxn.TxnStorager = &RDBTxnStorager{}

func (ps *RDBTxnStorager) AtomExecute(ctx context.Context, do func(context.Context) error) error {
	dbCtx, err := ps.dbCtxFactory(ctx, lib.OpenTxn())
	if err != nil {
		return err
	}

	return lib.RDBTxnExecute(dbCtx, do)
}
