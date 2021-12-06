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

package middleware

import (
	"github.com/rs/cors"
)

// NewCors enable Cross-Origin Resource Sharing
// detail see: https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Access_control_CORS
func NewCors() *cors.Cors {
	options := cors.Options{
		// allowed hosts
		AllowedOrigins: []string{"*"},
		// allowed HTTP methods
		AllowedMethods: []string{"POST", "GET", "PUT", "PATCH", "DELETE", "HEAD"},
		// allowed HTTP headers
		AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization", "Session_key", "Clientip"},
		// If set, allows to share auth credentials such as cookies
		AllowCredentials: true,
	}

	return cors.New(options)
}
