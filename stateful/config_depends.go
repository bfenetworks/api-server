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
	"io/ioutil"
	"path/filepath"

	"github.com/bfenetworks/api-server/lib"
)

type DependsConfig struct {
	NavTreeFile string `validate:"required,min=1"`
	I18nDir     string `validate:"required,min=1"`
	UIIcon      string
	UILogo      string

	role2Nav map[string]*NavTree
}

type I18nConfig struct {
	Lang    string
	Mapping map[string]string
}

func (dc *DependsConfig) Role2Nav() map[string]*NavTree {
	return dc.role2Nav
}

func (d *DependsConfig) Init() error {
	tree := &NavTree{}
	err := lib.LoadConfAuto(d.NavTreeFile, tree)
	if err != nil {
		return err
	}
	tree.init()

	role2Nav := map[string]*NavTree{}
	for role := range roles {
		role2Nav[role], err = tree.deriveByRole(role)
		if err != nil {
			return err
		}
	}

	d.role2Nav = role2Nav

	i18nfiles, err := ioutil.ReadDir(d.I18nDir)
	if err != nil {
		return err
	}

	langs := []*I18nConfig{}
	for _, f := range i18nfiles {
		config := &I18nConfig{}
		if err := lib.LoadConfAuto(filepath.Join(d.I18nDir, f.Name()), &config); err != nil {
			return err
		}

		langs = append(langs, config)
	}

	return InitI18n(langs)
}

var (
	roles = map[string]bool{}
)

type NavTree struct {
	ID         string     `json:"id" toml:"id"`
	Text       string     `json:"text" toml:"text"`
	AllowRoles []string   `json:"allowed_roles,omitempty" toml:"allowed_roles"`
	Children   []*NavTree `json:"children,omitempty" toml:"children"`
}

func (an *NavTree) init() {
	for _, role := range an.AllowRoles {
		roles[role] = true
	}

	for _, child := range an.Children {
		child.init()
	}
}

func (an *NavTree) matchRole(role string) bool {
	if an.AllowRoles == nil {
		return true
	}
	for _, one := range an.AllowRoles {
		if one == role {
			return true
		}
	}

	return false
}

func (an *NavTree) deriveByRole(role string) (*NavTree, error) {
	if !an.matchRole(role) {
		return nil, nil
	}

	var children []*NavTree
	for _, child := range an.Children {
		childDerive, err := child.deriveByRole(role)
		if err != nil {
			return nil, err
		}
		if childDerive != nil {
			children = append(children, childDerive)
		}
	}

	if (children) == nil && an.AllowRoles == nil {
		return nil, nil
	}

	return &NavTree{
		Text:     an.Text,
		ID:       an.ID,
		Children: children,
	}, nil
}
