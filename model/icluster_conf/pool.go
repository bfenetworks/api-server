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

package icluster_conf

import (
	"context"
	"github.com/bfenetworks/api-server/model/ibasic"
)

var (
	PoolTagBFE     int8 = 1
	PoolTagProduct int8 = 2
)

type PoolFilter struct {
	Name      *string
	IDs       []int64
	ID        *int64
	ProductID *int64
}

type PoolParam struct {
	ID        *int64
	Name      *string
	ProductID *int64
	Type      *int8
	Tag       *int8
}

type Pool struct {
	ID      int64
	Name    string
	Type    int8
	Ready   bool
	Tag     int8
	Product *ibasic.Product

	instances []Instance
}

func (p *Pool) SetDefaultInstances(is []Instance) {
	p.instances = is
}

func (p *Pool) GetDefaultPool() *InstancePool {
	return &InstancePool{
		Name:      p.Name,
		Instances: p.instances,
	}
}

type PoolStorage interface {
	FetchPool(ctx context.Context, name string) (*Pool, error)
	FetchPools(ctx context.Context, param *PoolFilter) ([]*Pool, error)

	CreatePool(ctx context.Context, product *ibasic.Product, data *PoolParam) (*Pool, error)
	UpdatePool(ctx context.Context, oldData *Pool, diff *PoolParam) error
	DeletePool(ctx context.Context, pool *Pool) error
}
