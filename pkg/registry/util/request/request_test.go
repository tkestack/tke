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

package request

import (
	"net/http"
	"testing"
)

func TestTenantID(t *testing.T) {
	var tests = []struct {
		name          string
		req           *http.Request
		defaultTenant string
		domainSuffix  string
		expectTenant  string
	}{
		{
			name: "domain suffix with dot",
			req: &http.Request{
				Host: "t.example.com:80",
			},
			defaultTenant: "default",
			domainSuffix:  ".example.com",
			expectTenant:  "t",
		}, {
			name: "domain suffix without dot",
			req: &http.Request{
				Host: "t.example.com:443",
			},
			defaultTenant: "default",
			domainSuffix:  "example.com",
			expectTenant:  "t",
		}, {
			name: "request host without domain",
			req: &http.Request{
				Host: "127.0.0.1:80",
			},
			defaultTenant: "default",
			domainSuffix:  ".example.com",
			expectTenant:  "default",
		}, {
			name: "request host with port",
			req: &http.Request{
				Host: "t.example.com:8080",
			},
			defaultTenant: "default",
			domainSuffix:  ".example.com",
			expectTenant:  "t",
		}, {
			name: "domain is default tenant",
			req: &http.Request{
				Host: "default.example.com:443",
			},
			defaultTenant: "default",
			domainSuffix:  "example.com",
			expectTenant:  "default",
		},
	}

	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			tenantID := TenantID(rt.req, rt.domainSuffix, rt.defaultTenant)
			if tenantID != rt.expectTenant {
				t.Errorf("failed %s: expected: %s actual: %s", rt.name, rt.expectTenant, tenantID)
			}
		})
	}
}
