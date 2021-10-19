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
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var Cause = errors.Cause

const (
	etParam            = "PARAM"
	etModel            = "Model"
	etDirtyData        = "Model.DirtyData"
	etNullData         = "Model.NullData"
	etDependentUnReady = "Model.DependentUnReady"
	etExistedData      = "Model.ExistedData"
	etDao              = "DAO"

	etAuthenticateFail = "Authentication.Fail"
	etAuthorizateFail  = "Authorization.Fail"
)

func unwrapMsg(err error) string {
	if err == nil {
		return ""
	}

	return strings.SplitN(err.Error(), ": ", 2)[0]
}

// WrapParamErrorWithMsg  Just endpoints invoke
func WrapParamErrorWithMsg(tip string, args ...interface{}) error {
	return WrapParamError(fmt.Errorf(tip, args...))
}

// WrapParamError  Just endpoints layout invoke
func WrapParamError(err error) error {
	if err == nil {
		return nil
	}
	if unwrapMsg(err) == etParam {
		return err
	}

	return errors.WithMessage(err, etParam)
}

// WrapDaoError just storage layout invoke
func WrapDaoError(err error) error {
	if err == nil {
		return nil
	}
	if unwrapMsg(err) == etDao {
		return err
	}
	return errors.Wrap(err, etDao)
}

// WrapModelError Just Service layout invoke
func WrapModelError(err error) error {
	if err == nil {
		return nil
	}
	if unwrapMsg(err) == etModel {
		return err
	}
	return errors.Wrap(err, etModel)
}

func WrapAuthorizateFailErrorWithMsg(msg string, args ...interface{}) error {
	return errors.Wrap(fmt.Errorf(msg, args...), etAuthorizateFail)
}

func WrapAuthenticateFailErrorWithMsg(msg string, args ...interface{}) error {
	return errors.Wrap(fmt.Errorf(msg, args...), etAuthenticateFail)
}

// WrapDependentUnReadyErrorWithMsg Just Service layout invoke
func WrapDependentUnReadyErrorWithMsg(msg string, args ...interface{}) error {
	return errors.Wrap(fmt.Errorf(msg, args...), etDependentUnReady)
}

// WrapModelErrorWithMsg Just Service layout invoke
func WrapModelErrorWithMsg(msg string, args ...interface{}) error {
	return errors.Wrap(fmt.Errorf(msg, args...), etModel)
}

func WrapRecordNotExist(topic ...string) error {
	msg := "Record Not Exist"
	if len(topic) == 1 {
		msg = topic[0] + " " + msg
	}
	return errors.Wrap(fmt.Errorf(msg), etNullData)
}

// WrapRecordExisted Record Existed
func WrapRecordExisted(topic ...string) error {
	msg := "Record Existed"

	if len(topic) == 1 {
		msg = topic[0] + " " + msg
	}
	return errors.Wrap(fmt.Errorf(msg), etExistedData)
}

// WrapDirtyDataError Just Service layout invoke
func WrapDirtyDataError(err error) error {
	if err == nil {
		return nil
	}
	if unwrapMsg(err) == etDirtyData {
		return err
	}
	return errors.Wrap(err, etDirtyData)
}

// WrapDirtyDataErrorWithMsg Just Service layout invoke
func WrapDirtyDataErrorWithMsg(msg string, args ...interface{}) error {
	return errors.Wrap(fmt.Errorf(msg, args...), etDirtyData)
}
