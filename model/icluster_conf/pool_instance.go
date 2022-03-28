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
	"fmt"

	"github.com/bfenetworks/api-server/lib/xerror"
)

type Instance struct {
	HostName string            `json:"Name"`
	IP       string            `json:"Addr"`
	Port     int               `json:"Port"`
	Ports    map[string]int    `json:"Ports,omitempty"`
	Tags     map[string]string `json:"tags,omitempty"`
	Weight   int64             `json:"Weight"`
	Disable  bool              `json:"Disable"`
}

func (i *Instance) IPWithPort() string {
	if i.Port == 0 {
		i.Port = i.Ports["Default"]
	}

	return fmt.Sprintf("%s:%d", i.IP, i.Port)
}

type InstancePool struct {
	Name      string
	Instances []Instance
}

const (
	InstancePoolTypeRDB   int8 = 1
	InstancePoolTypeNacos int8 = 2
)

type InstancePoolStorager interface {
	UpdateInstances(context.Context, *Pool, *InstancePool) error

	BatchFetchInstances(context.Context, []*Pool) (map[string]*InstancePool, error)
}

type InstancePoolManager struct {
	instancePoolStorages map[int8]InstancePoolStorager
}

func NewInstancePoolManager(instancePoolStorages map[int8]InstancePoolStorager) *InstancePoolManager {
	return &InstancePoolManager{
		instancePoolStorages: instancePoolStorages,
	}
}

func (pim *InstancePoolManager) BatchFetchInstances(ctx context.Context, pools []*Pool) (map[string]*InstancePool, error) {
	type2InstancePoolList := map[int8][]*Pool{}
	for _, one := range pools {
		type2InstancePoolList[one.Type] = append(type2InstancePoolList[one.Type], one)
	}

	for typ := range type2InstancePoolList {
		_, ok := pim.instancePoolStorages[typ]
		if !ok {
			return nil, xerror.WrapModelErrorWithMsg("Type %d not register Storager", typ)
		}
	}

	rst := map[string]*InstancePool{}
	for typ, pisList := range type2InstancePoolList {
		storager := pim.instancePoolStorages[typ]
		r, err := storager.BatchFetchInstances(ctx, pisList)
		if err != nil {
			return nil, err
		}

		for name, pis := range r {
			rst[name] = pis
		}
	}

	return rst, nil
}

func (pim *InstancePoolManager) UpdateInstances(ctx context.Context, pool *Pool, pis *InstancePool) error {
	storager, ok := pim.instancePoolStorages[pool.Type]
	if !ok {
		return xerror.WrapModelErrorWithMsg("Type %d not register Storager", pool.Type)
	}

	return storager.UpdateInstances(ctx, pool, pis)
}

func PoolMap2List(m map[string]*Pool) []*Pool {
	var r []*Pool
	for _, one := range m {
		r = append(r, one)
	}

	return r
}
