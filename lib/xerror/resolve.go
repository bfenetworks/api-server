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

package xerror

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type ResolveResult struct {
	ErrNo int
	Type  string
	Msg   string

	err error
}

func (rr *ResolveResult) String() string {
	if rr == nil {
		return "<nil>"
	}

	bs, _ := json.Marshal(rr)
	return string(bs)
}

func (rr *ResolveResult) FullMsg() string {
	msg := fmt.Sprintf("%+v", rr.err)
	if i := strings.Index(msg, "\nnet/http.HandlerFunc.ServeHTTP"); i != -1 {
		msg = msg[:i]
	}
	return msg
}

func Resolve(err error) *ResolveResult {
	if err == nil {
		return nil
	}

	rr := &ResolveResult{
		ErrNo: 500,
		Msg:   fmt.Sprintf("%v", errors.Cause(err)),

		err: err,
	}

	msg := unwrapMsg(err)
	switch msg {
	case etParam:
		rr.ErrNo = 422
		rr.Type = "Param Illegal"
	case etModel:
		rr.Type = "Biz Exception"
	case etDirtyData:
		rr.Type = "Inner Dirty Data"
	case etDao:
		rr.Type = "Database Exception"
	case etNullData:
		rr.Type = "Record Not Exist"
		rr.ErrNo = 404
	case etExistedData:
		rr.Type = "Record Existed"
		rr.ErrNo = 555
	case etDependentUnReady:
		rr.Type = "Dependent Not Ready"
		rr.Msg = err.Error()
		rr.ErrNo = 510
	case etAuthenticateFail:
		rr.ErrNo = 401
		rr.Type = "Authenticate Fail"
	case etAuthorizateFail:
		rr.ErrNo = 402
		rr.Type = "Authorizate Fail"

	default: // never come here
		rr.Type = "Unknown Exception"
	}

	return rr
}
