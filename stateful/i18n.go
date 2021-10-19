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

package stateful

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/bfenetworks/api-server/lib"
)

var langs2languagePacks = map[string]map[string]*LangMapping{}

type LangMapping struct {
	From string
	To   string

	placeholderCount int
	fromRegex        *regexp.Regexp
}

func NewLangMapping(from, to string) (*LangMapping, error) {
	mm := &LangMapping{
		From: from,
		To:   to,
	}

	var err error
	mm.fromRegex, err = regexp.Compile(mm.From)
	if err != nil {
		return nil, fmt.Errorf("regex fail, string: %s, err: %v", mm.From, err)
	}
	mm.placeholderCount = strings.Count(mm.To, "%s")

	return mm, nil
}

func (mm *LangMapping) TryTrans(raw string) (string, bool) {
	// if len(typeAndMsg)
	if mm.fromRegex == nil {
		return raw, false
	}

	params := mm.fromRegex.FindStringSubmatch(raw)
	if params == nil {
		return raw, false
	}

	if mm.placeholderCount == 0 {
		return mm.To, true
	}

	list := make([]interface{}, mm.placeholderCount)
	for i := range list {
		if i+1 < len(params) {
			list[i] = params[i+1]
		} else {
			list[i] = ""
		}
	}
	return fmt.Sprintf(mm.To, list...), true
}

func InitI18n(i18nConfig []*I18nConfig) error {
	mm := map[string]map[string]*LangMapping{}
	for _, lang := range i18nConfig {
		for k, v := range lang.Mapping {
			m, err := NewLangMapping(k, v)
			if err != nil {
				return nil
			}

			if _, ok := mm[lang.Lang]; !ok {
				mm[lang.Lang] = map[string]*LangMapping{}
			}

			mm[lang.Lang][m.From] = m
		}
	}

	langs2languagePacks = mm
	return nil
}

func getLangPack(req *http.Request) (map[string]*LangMapping, bool) {
	for _, lang := range lib.AccepLanguages(req) {
		langConfig, ok := langs2languagePacks[lang]
		if ok {
			return langConfig, ok
		}

	}

	return nil, false

}

func TryMappingErrMsg(req *http.Request, errMsg string) string {
	if errMsg == "" {
		return errMsg
	}

	langConfig, ok := getLangPack(req)
	if !ok {
		return errMsg
	}

	// errMsg look like {type}: {msg}
	typeAndMsg := strings.SplitN(errMsg, ": ", 2)
	var t, m string
	if len(typeAndMsg) == 2 {
		t, m = typeAndMsg[0], typeAndMsg[1]
	} else {
		m = typeAndMsg[0]
	}
	if t != "" {
		newT, ok := langConfig[t]
		if ok {
			t = newT.To
		}
	}

	for _, mapping := range langConfig {
		m, ok = mapping.TryTrans(m)
		if ok {
			break
		}
	}

	if len(typeAndMsg) == 2 {
		return t + ": " + m
	}

	return m
}
