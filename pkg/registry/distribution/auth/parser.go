/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package auth

import (
	"fmt"
	"net/url"
	"strings"
	"tkestack.io/tke/pkg/registry/distribution/tenant"
)

type imageParser interface {
	parse(s string) (*image, error)
}

type image struct {
	tenantID  string
	namespace string
	repo      string
	tag       string
}

type basicParser struct{}

func (b basicParser) parse(s string) (*image, error) {
	return parseImg(s)
}

// build Image accepts a string like library/ubuntu:14.04 and build a image struct
func parseImg(s string) (*image, error) {
	repo := strings.SplitN(s, "/", 2)
	if len(repo) < 2 {
		return nil, fmt.Errorf("unable to parse image from string: %s", s)
	}

	var tenantID, namespace string
	tns := strings.SplitN(repo[0], "-", 2)
	switch len(tns) {
	case 1:
		if tns[0] != tenant.CrossTenantNamespace {
			return nil, fmt.Errorf("image `%s` belongs to an unknown tenant", repo[0])
		}
		tenantID = ""
		namespace = tenant.CrossTenantNamespace
	case 2:
		tenantID = tns[0]
		namespace = tns[1]
	}

	i := strings.SplitN(repo[1], ":", 2)

	res := &image{
		tenantID:  tenantID,
		namespace: namespace,
		repo:      i[0],
	}
	if len(i) == 2 {
		res.tag = i[1]
	}
	return res, nil
}

func parseScopes(u *url.URL) []string {
	var result []string
	for _, sector := range u.Query()["scope"] {
		result = append(result, strings.Split(sector, " ")...)
	}
	return result
}
