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
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"net/http"
	genericoidc "tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
)

// RequestUser according to the basic auth credentials carried in the http
// request, use the password as an APIKey to call the authentication method,
// and return the user's username and tenant id.
func RequestUser(req *http.Request, apiKeyAuthenticator authenticator.Password) (username, tenantID string, authenticated bool) {
	user, password, ok := req.BasicAuth()
	if !ok {
		return
	}
	res, authOk, err := apiKeyAuthenticator.AuthenticatePassword(req.Context(), user, password)
	if err != nil || !authOk || res == nil {
		return
	}
	username = res.User.GetName()
	extra := res.User.GetExtra()
	if len(extra) > 0 {
		if t, ok := extra[genericoidc.TenantIDKey]; ok {
			if len(t) > 0 {
				tenantID = t[0]
			}
		}
	}
	authenticated = true
	return
}
