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

package authentication

import (
	"context"

	"k8s.io/apiserver/pkg/endpoints/request"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
)

// UsernameAndTenantID implementation decomposition in the original
// kubernetes api server the user name obtained in *Userinfo is the actual
// username and tenant ID.
func UsernameAndTenantID(ctx context.Context) (username string, tenantID string) {
	userInfo, ok := request.UserFrom(ctx)
	if !ok {
		return "", ""
	}
	extra := userInfo.GetExtra()
	if len(extra) > 0 {
		if tenantIDs, ok := extra[oidc.TenantIDKey]; ok {
			if len(tenantIDs) > 0 {
				tenantID = tenantIDs[0]
			}
		}
	}
	return userInfo.GetName(), tenantID
}

func Groups(ctx context.Context) (groups []string) {
	userInfo, ok := request.UserFrom(ctx)
	if !ok {
		return nil
	}
	return userInfo.GetGroups()
}

// IsAdministrator check whether administrator
func IsAdministrator(ctx context.Context, privilegedUsername string) bool {
	username, tenantID := UsernameAndTenantID(ctx)
	return username == privilegedUsername && tenantID == ""
}
