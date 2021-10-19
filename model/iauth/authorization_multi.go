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

package iauth

import (
	"context"
)

type Op string

const (
	OpAnd Op = "&&"
)

type MultiAuthorizer struct {
	as []Authorizer
	op Op

	// some authorizer may not init before be used
	asFactories []func() Authorizer
}

func NewMultiAuthorizer(as []Authorizer, op Op) *MultiAuthorizer {
	return &MultiAuthorizer{
		as: as,
		op: op,
	}
}

func NewMultiAuthorizerWithFactory(asFactories []func() Authorizer, op Op) *MultiAuthorizer {
	return &MultiAuthorizer{
		asFactories: asFactories,
		op:          op,
	}
}

func (a *MultiAuthorizer) init() {
	if a.as != nil {
		return
	}

	a.as = []Authorizer{}
	for _, one := range a.asFactories {
		a.as = append(a.as, one())
	}
}

func (a *MultiAuthorizer) Authorizate(ctx context.Context, user *User) (ok bool, err error) {
	a.init()

	for _, one := range a.as {
		ok, err = one.Authorizate(ctx, user)
		if !ok || err != nil {
			return
		}
	}

	return true, nil
}

func (a *MultiAuthorizer) Grant(ctx context.Context, user *User) (err error) {
	a.init()

	for _, one := range a.as {
		if err := one.Grant(ctx, user); err != nil {
			return err
		}
	}

	return nil
}

func (a *MultiAuthorizer) Revoke(ctx context.Context, user *User) (err error) {
	a.init()

	for _, one := range a.as {
		if err := one.Revoke(ctx, user); err != nil {
			return err
		}
	}

	return nil
}
