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

package xreq

import (
	"encoding/json"
	"net/http"

	"github.com/bfenetworks/api-server/lib"
	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/lib/xreq/internal"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	t_en "github.com/go-playground/validator/v10/translations/en"
	t_zh "github.com/go-playground/validator/v10/translations/zh"
	"github.com/gorilla/mux"
)

var (
	_zh = zh.New()
	_en = en.New()

	zhTrans, _   = ut.New(_zh, _zh).GetTranslator("zh")
	enTrans, _   = ut.New(_en, _en).GetTranslator("en")
	defaultTrans = enTrans

	langs = map[string]ut.Translator{
		"zh": zhTrans,
		"en": enTrans,
	}

	validate = validator.New()
)

func getTranslator(req *http.Request) ut.Translator {
	for _, one := range lib.AccepLanguages(req) {
		if l, ok := langs[one]; ok {
			return l
		}
	}

	return defaultTrans
}

func init() {
	if err := t_zh.RegisterDefaultTranslations(validate, zhTrans); err != nil {
		panic(err)
	}
	if err := t_en.RegisterDefaultTranslations(validate, enTrans); err != nil {
		panic(err)
	}
}

func JSONDeserializer(req *http.Request, data interface{}) error {
	if req == nil || req.Body == nil {
		return nil
	}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(data)
	if err != nil {
		return xerror.WrapParamError(err)
	}

	return nil
}

func ValidateData(data interface{}, lang ut.Translator) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		if lang == nil {
			lang = defaultTrans
		}
		return xerror.WrapParamErrorWithMsg(validationErrors[0].Translate(lang))
	}

	return xerror.WrapParamError(err)
}

func Bind(req *http.Request, data interface{}) error {
	if err := JSONDeserializer(req, data); err != nil {
		return err
	}

	vars := mux.Vars(req)
	if err := internal.MapUri(data, m2ms(vars)); err != nil {
		return err
	}
	if err := ValidateData(data, getTranslator(req)); err != nil {
		return err
	}

	return nil
}

func BindJSON(req *http.Request, data interface{}) error {
	if err := JSONDeserializer(req, data); err != nil {
		return err
	}

	if err := ValidateData(data, getTranslator(req)); err != nil {
		return err
	}

	return nil
}

func BindURI(req *http.Request, data interface{}) error {
	vars := mux.Vars(req)
	if err := internal.MapUri(data, m2ms(vars)); err != nil {
		return err
	}

	if err := ValidateData(data, getTranslator(req)); err != nil {
		return err
	}

	return nil
}

func BindForm(req *http.Request, data interface{}) error {
	req.ParseForm()
	if err := internal.MapForm(data, req.Form); err != nil {
		return err
	}

	if err := ValidateData(data, getTranslator(req)); err != nil {
		return err
	}

	return nil
}

func m2ms(m map[string]string) map[string][]string {
	r := map[string][]string{}
	for k, v := range m {
		r[k] = []string{v}
	}
	return r
}
