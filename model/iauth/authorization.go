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

/*
SomeOne Can READ/UPDATE/CREATE/DELETE/EXPORT some Feature
*/

type Authorizer interface {
	Authorizate(context.Context, *User) (bool, error)
	Grant(context.Context, *User) error
	Revoke(context.Context, *User) error
}

var (
	_ Authorizer = &FeatureAuthorizer{}
	_ Authorizer = &MultiAuthorizer{}
	_ Authorizer = &ProductAuthorizateManager{}
)

func Authorizate(ctx context.Context, authoriztor Authorizer, user *User) (ok bool, err error) {
	return authoriztor.Authorizate(ctx, user)
}

func Grant(ctx context.Context, authoriztor Authorizer, user *User) (err error) {
	return authoriztor.Grant(ctx, user)
}

func Revoke(ctx context.Context, authoriztor Authorizer, user *User) (err error) {
	return authoriztor.Revoke(ctx, user)
}
