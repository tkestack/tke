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
	"strings"

	utilhttp "tkestack.io/tke/pkg/util/http"
)

// TenantID according to the host name accessed by the http request, combined
// with the configured default tenant name and domain name suffix, returns the
// tenant name accessed by the user.
func TenantID(req *http.Request, domainSuffix string, defaultTenant string) string {
	if domainSuffix == "" {
		return defaultTenant
	}
	domain := utilhttp.DomainFromRequest(req)
	if strings.HasSuffix(domain, domainSuffix) {
		tenant := strings.TrimSuffix(domain, domainSuffix)
		if strings.HasSuffix(tenant, ".") {
			tenant = strings.TrimSuffix(tenant, ".")
		}
		if tenant != "" {
			return tenant
		}
	}
	return defaultTenant
}
