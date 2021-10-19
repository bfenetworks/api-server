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

package ibasic

import (
	"context"
	"time"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/itxn"
)

const (
	ResourceProduct = "product"
)

type key string

var keyProduct key = "product"

func NewProductContext(ctx context.Context, product *Product) context.Context {
	return context.WithValue(ctx, keyProduct, product)
}

func MustGetProduct(ctx context.Context) (*Product, error) {
	obj := ctx.Value(keyProduct)
	if obj == nil {
		return nil, xerror.WrapParamErrorWithMsg("Fail To Get Product")
	}

	return obj.(*Product), nil
}

var (
	BuildinProduct = &Product{
		ID:   1,
		Name: "BFE",
	}
	ProxyProduct = &Product{
		ID:   1,
		Name: "proxy",
	}
)

type Product struct {
	ID                int64
	Name              string
	Description       string
	MailList          []string
	PhoneList         []string
	ContactPersonList []string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductFilter struct {
	ID   *int64
	NeID *int64
	IDs  []int64

	Name *string
}

type ProductParam struct {
	ID                *int64
	Name              *string
	Description       *string
	MailList          []string
	PhoneList         []string
	ContactPersonList []string
}

type ProductStorager interface {
	FetchProducts(context.Context, *ProductFilter) ([]*Product, error)
	DeleteProduct(context.Context, *Product) error
	CreateProduct(context.Context, *ProductParam) error
	UpdateProduct(context.Context, *Product, *ProductParam) error
}

type ProductManager struct {
	storager ProductStorager
	txn      itxn.TxnStorager
}

func NewProductManager(txn itxn.TxnStorager, storager ProductStorager) *ProductManager {
	return &ProductManager{
		txn:      txn,
		storager: storager,
	}
}

func (pm *ProductManager) FetchProducts(ctx context.Context, param *ProductFilter) (list []*Product, err error) {
	err = pm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err = pm.storager.FetchProducts(ctx, param)
		return err
	})

	return
}

func (pm *ProductManager) DeleteProduct(ctx context.Context, p *Product) (err error) {
	if p.ID == 1 {
		return xerror.WrapModelErrorWithMsg("Cant Delete Build-in Product")
	}
	err = pm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		return pm.storager.DeleteProduct(ctx, p)
	})

	return
}

func (pm *ProductManager) CreateProduct(ctx context.Context, p *ProductParam) (err error) {
	err = pm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err := pm.storager.FetchProducts(ctx, &ProductFilter{
			Name: p.Name,
		})
		if err != nil {
			return err
		}
		if len(list) != 0 {
			return xerror.WrapRecordExisted("Product")
		}

		return pm.storager.CreateProduct(ctx, p)
	})

	return
}

func (pm *ProductManager) UpdateProduct(ctx context.Context, p *Product, newVal *ProductParam) (err error) {
	if p.ID == 1 {
		return xerror.WrapModelErrorWithMsg("Cant Delete Build-in Product")
	}

	err = pm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		return pm.storager.UpdateProduct(ctx, p, newVal)
	})

	return
}

func ProductID2NameMap(list []*Product) map[int64]string {
	m := make(map[int64]string, len(list))
	for _, one := range list {
		m[one.ID] = one.Name
	}

	return m
}

func ProductIDMap(list []*Product) map[int64]*Product {
	m := make(map[int64]*Product, len(list))
	for _, one := range list {
		m[one.ID] = one
	}

	return m
}
