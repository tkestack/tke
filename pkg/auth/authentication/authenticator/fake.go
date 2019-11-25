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

package authenticator

import (
	genericauthenticator "k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
	"net/http"
	genericoidc "tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
)

// NewFakeAuthenticator creates new fake authenticator object.
func NewFakeAuthenticator() genericauthenticator.Request {
	return genericauthenticator.RequestFunc(func(req *http.Request) (*genericauthenticator.Response, bool, error) {
		auds, _ := genericauthenticator.AudiencesFrom(req.Context())
		return &genericauthenticator.Response{
			User: &user.DefaultInfo{
				Name: "fake-user",
				Extra: map[string][]string{
					genericoidc.TenantIDKey: {"default"},
				},
			},
			Audiences: auds,
		}, true, nil
	})
}
