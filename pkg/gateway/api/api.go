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

package api

import (
	"github.com/emicklei/go-restful"
	"golang.org/x/oauth2"
	"net/http"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
)

// GroupName is the api group name for gateway.
const GroupName = "gateway.tkestack.io"

// Version is the api version for gateway.
const Version = "v1"

// RegisterRoute is used to register prefix path routing matches for all
// configured backend components.
func RegisterRoute(container *restful.Container, cfg *gatewayconfig.GatewayConfiguration, oauthConfig *oauth2.Config, oidcHTTPClient *http.Client, oidcAuthenticator *oidc.Authenticator) error {
	registerTokenRoute(container, oauthConfig, oidcHTTPClient, oidcAuthenticator, cfg.DisableOIDCProxy)
	registerSysInfoRoute(container, cfg)
	registerLogoutRoute(container)
	return nil
}
