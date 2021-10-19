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

package iroute_conf

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/bfenetworks/bfe/bfe_config/bfe_route_conf/host_rule_conf"

	"github.com/bfenetworks/api-server/lib/xerror"
	"github.com/bfenetworks/api-server/model/ibasic"
	"github.com/bfenetworks/api-server/model/itxn"
)

type Domain struct {
	ID        int64
	ProductID int64
	Name      string

	UsingAdvancedRedirect int8
	UsingAdvancedHsts     int8
}

type DomainFilter struct {
	Product *ibasic.Product
	Name    *string
}

type DomainParam struct {
	ProductID *int64
	Name      *string

	UsingAdvancedRedirect *int8
	UsingAdvancedHsts     *int8
}

func newHostTableConf(version string, productMapID2Name map[int64]string,
	domains []*Domain) *host_rule_conf.HostTableConf {
	tagNameHandler := func(productName, host string) string {
		return strings.ToLower(productName)
	}

	tag2hosts, product2tags := map[string][]string{}, map[string]map[string]bool{
		defaultProduct: {},
	}

	sort.Slice(domains, func(i, j int) bool {
		return domains[i].ProductID < domains[j].ProductID
	})
	for _, domain := range domains {
		productName := productMapID2Name[domain.ProductID]
		tagName := tagNameHandler(productName, domain.Name)
		tag2hosts[tagName] = append(tag2hosts[tagName], domain.Name)

		tagMap := product2tags[productName]
		if tagMap == nil {
			tagMap = map[string]bool{}
		}
		tagMap[tagName] = true
		product2tags[productName] = tagMap
	}

	_tag2host := host_rule_conf.HostTagToHost{}
	hostTags := host_rule_conf.ProductToHostTag{}

	for tag, hosts := range tag2hosts {
		tmp := host_rule_conf.HostnameList(hosts)
		_tag2host[tag] = &tmp
	}

	for product, tags := range product2tags {
		tagList := host_rule_conf.HostTagList{}
		for tag := range tags {
			tagList = append(tagList, tag)
		}

		hostTags[product] = &tagList
	}

	return &host_rule_conf.HostTableConf{
		Version:        &version,
		DefaultProduct: &defaultProduct,
		Hosts:          &_tag2host,
		HostTags:       &hostTags,
	}
}

type DomainStorager interface {
	FetchDomains(ctx context.Context, param *DomainFilter) ([]*Domain, error)
	CreateDomain(ctx context.Context, product *ibasic.Product, param *DomainParam) error
	DeleteDomain(ctx context.Context, product *ibasic.Product, domain *Domain) error
}

type DomainManager struct {
	txn itxn.TxnStorager

	storager         DomainStorager
	routeRuleManager *RouteRuleManager
}

func NewDomainManager(txn itxn.TxnStorager, storager DomainStorager,
	routeRuleManager *RouteRuleManager) *DomainManager {

	return &DomainManager{
		txn:              txn,
		storager:         storager,
		routeRuleManager: routeRuleManager,
	}
}

func (m *DomainManager) DomainList(ctx context.Context, param *DomainFilter) (list []*Domain, err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err = m.storager.FetchDomains(ctx, param)
		return err
	})

	return
}

func (m *DomainManager) CreateDomain(ctx context.Context, product *ibasic.Product, param *DomainParam) (err error) {
	err = m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		param.ProductID = &product.ID
		domains, err := m.storager.FetchDomains(ctx, nil)
		if err != nil {
			return err
		}

		// copy golang/src/crypto/x509/verify.go matchHostnames() func
		domainMatch := func(pattern string, domain string) bool {
			domain = strings.TrimSuffix(domain, ".")
			pattern = strings.TrimSuffix(pattern, ".")

			if len(pattern) == 0 || len(domain) == 0 {
				return false
			}

			patternParts := strings.Split(pattern, ".")
			hostParts := strings.Split(domain, ".")

			if len(patternParts) != len(hostParts) {
				return false
			}

			for i, patternPart := range patternParts {
				if i == 0 && patternPart == "*" {
					continue
				}
				if patternPart != hostParts[i] {
					return false
				}
			}

			return true
		}

		domainName := *param.Name
		for _, oldDomain := range domains {
			if oldDomain.Name == domainName {
				return xerror.WrapRecordExisted("Domain")
			}

			a, b := domainName, oldDomain.Name
			if !strings.Contains(domainName, "*") {
				a, b = oldDomain.Name, domainName
			}
			// check if oldDomains is covered by domainName
			if domainMatch(a, b) {
				return xerror.WrapModelErrorWithMsg("Domain Name %s Be Covered By Wildcard Domain %s", b, a)
			}
		}

		return m.storager.CreateDomain(ctx, product, param)
	})

	return
}

func (m *DomainManager) DeleteDomain(ctx context.Context, product *ibasic.Product, domain *Domain) (err error) {
	dbui, err := m.BeUsed(ctx, product, domain)
	if err != nil {
		return err
	}
	if dbui != nil {
		return xerror.WrapDependentUnReadyErrorWithMsg(dbui.String())
	}

	return m.txn.AtomExecute(ctx, func(ctx context.Context) error {
		return m.storager.DeleteDomain(ctx, product, domain)
	})
}

func (m *DomainManager) BeUsed(ctx context.Context, product *ibasic.Product, domain *Domain) (*DomainBeUsedInfo, error) {
	if domain.UsingAdvancedHsts != 0 || domain.UsingAdvancedRedirect != 0 {
		return &DomainBeUsedInfo{
			domain:         domain,
			hasHTTPSConfig: true,
		}, nil
	}

	rule, err := m.routeRuleManager.FetchProductRule(ctx, product)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, nil
	}
	if useInfo := rule.HostBeUsed(domain.Name); useInfo != nil {
		return &DomainBeUsedInfo{
			domain:   domain,
			RoutRule: useInfo,
		}, nil
	}

	return nil, nil
}

type DomainBeUsedInfo struct {
	domain         *Domain
	hasHTTPSConfig bool

	RoutRule *HostUsedInfo
}

func (dbui *DomainBeUsedInfo) String() string {
	domainName := dbui.domain.Name
	if d := dbui.RoutRule; d != nil {
		return fmt.Sprintf("Domain %s Be Used By %s Rule %s", domainName, d.Type, d.Detail)
	}

	if dbui.hasHTTPSConfig {
		return fmt.Sprintf("Domain %s Be Used By HTTPS Config", domainName)
	}

	return ""
}

func (dbui *DomainBeUsedInfo) Dependent() (typ, name string) {
	if d := dbui.RoutRule; d != nil {
		return d.Type, d.Detail
	}

	if dbui.hasHTTPSConfig {
		return "DomainHttpsConfig", ""
	}

	return "", ""
}
