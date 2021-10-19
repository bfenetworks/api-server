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

package lib

import "time"

func PInt64(v int64) *int64 {
	return &v
}

func PString(v string) *string {
	return &v
}

func PUint32(v uint32) *uint32 {
	return &v
}

func PInt32(v int32) *int32 {
	return &v
}

func PInt8(v int8) *int8 {
	return &v
}

func PInt(v int) *int {
	return &v
}

func PBool(v bool) *bool {
	return &v
}

func PTime(v time.Time) *time.Time {
	return &v
}

func PTimeNow() *time.Time {
	return PTime(time.Now())
}
