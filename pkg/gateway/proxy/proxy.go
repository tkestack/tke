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

package proxy

import (
	"fmt"
	"k8s.io/apiserver/pkg/server/mux"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/api/business"
	"tkestack.io/tke/api/monitor"
	"tkestack.io/tke/api/notify"
	"tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
	"tkestack.io/tke/pkg/gateway/proxy/handler/frontproxy"
	"tkestack.io/tke/pkg/gateway/proxy/handler/passthrough"
	platformapiserver "tkestack.io/tke/pkg/platform/apiserver"
	"tkestack.io/tke/pkg/registry/distribution"
	"tkestack.io/tke/pkg/util/log"
)

const apiPrefix = "/apis"

type moduleName string

const (
	moduleNamePlatform moduleName = "platform"
	moduleNameBusiness moduleName = "business"
	moduleNameNotify   moduleName = "notify"
	moduleNameRegistry moduleName = "registry"
	moduleNameAuth     moduleName = "auth"
	moduleNameMonitor  moduleName = "monitor"
)

type modulePath struct {
	prefix    string
	protected bool
}

func componentPrefix() map[moduleName][]modulePath {
	componentPrefixMap := map[moduleName][]modulePath{
		moduleNamePlatform: {},
		moduleNameBusiness: {
			modulePath{
				prefix:    fmt.Sprintf("%s/%s/", apiPrefix, business.GroupName),
				protected: true,
			},
		},
		moduleNameNotify: {
			modulePath{
				prefix:    "/webhook/",
				protected: false,
			},
			modulePath{
				prefix:    fmt.Sprintf("%s/%s/", apiPrefix, notify.GroupName),
				protected: true,
			},
		},
		moduleNameAuth: {
			modulePath{
				prefix:    "/api/authv1/",
				protected: true,
			},
			modulePath{
				prefix:    fmt.Sprintf("%s/%s/", apiPrefix, auth.GroupName),
				protected: true,
			},
		},
		moduleNameMonitor: {
			modulePath{
				prefix:    "/apis/v1/monitor/",
				protected: true,
			},
			modulePath{
				prefix:    fmt.Sprintf("%s/%s/", apiPrefix, monitor.GroupName),
				protected: true,
			},
		},
		moduleNameRegistry: {
			modulePath{
				prefix:    distribution.PathPrefix,
				protected: false,
			},
			modulePath{
				prefix:    distribution.APIPrefix,
				protected: false,
			},
			modulePath{
				prefix:    fmt.Sprintf("%s/%s/", apiPrefix, registry.GroupName),
				protected: true,
			},
		},
	}
	platformResourceConfig := platformapiserver.DefaultAPIResourceConfigSource()
	if platformResourceConfig != nil {
		platformModulePath := componentPrefixMap[moduleNamePlatform]
		for k := range platformResourceConfig.GroupVersionConfigs {
			prefix := fmt.Sprintf("%s/%s/%s/", apiPrefix, k.Group, k.Version)
			if k.Group == "" {
				prefix = fmt.Sprintf("%s/%s/", "/api", k.Version)
			}
			platformModulePath = append(platformModulePath, modulePath{
				prefix:    prefix,
				protected: true,
			})
		}
		componentPrefixMap[moduleNamePlatform] = platformModulePath
	}
	return componentPrefixMap
}

// RegisterRoute is used to register prefix path routing matches for all
// configured backend components.
func RegisterRoute(m *mux.PathRecorderMux, cfg *gatewayconfig.GatewayConfiguration, oidcAuthenticator *oidc.Authenticator) error {
	pathPrefixProxyMap := prefixProxy(cfg)
	for pathPrefix, proxyComponent := range pathPrefixProxyMap {
		if proxyComponent.FrontProxy != nil {
			handler, err := frontproxy.NewHandler(proxyComponent.Address, proxyComponent.FrontProxy, oidcAuthenticator, pathPrefix.protected)
			if err != nil {
				return err
			}
			log.Info("Registered reverse proxy of front proxy mode for backend component", log.String("pathPrefix", pathPrefix.prefix), log.Bool("protected", pathPrefix.protected), log.String("address", proxyComponent.Address))
			m.HandlePrefix(pathPrefix.prefix, handler)
		} else if proxyComponent.Passthrough != nil {
			handler, err := passthrough.NewHandler(proxyComponent.Address, proxyComponent.Passthrough, pathPrefix.protected)
			if err != nil {
				return err
			}
			log.Info("Registered reverse proxy of front proxy mode for backend component", log.String("pathPrefix", pathPrefix.prefix), log.Bool("protected", pathPrefix.protected), log.String("address", proxyComponent.Address))
			m.HandlePrefix(pathPrefix.prefix, handler)
		}
	}
	return nil
}

func prefixProxy(cfg *gatewayconfig.GatewayConfiguration) map[modulePath]gatewayconfig.Component {
	componentPrefixMap := componentPrefix()
	pathPrefixProxyMap := make(map[modulePath]gatewayconfig.Component)
	// platform
	if cfg.Components.Platform != nil {
		if prefixes, ok := componentPrefixMap[moduleNamePlatform]; ok {
			for _, prefix := range prefixes {
				pathPrefixProxyMap[prefix] = *cfg.Components.Platform
			}
		}
	}
	// business
	if cfg.Components.Business != nil {
		if prefixes, ok := componentPrefixMap[moduleNameBusiness]; ok {
			for _, prefix := range prefixes {
				pathPrefixProxyMap[prefix] = *cfg.Components.Business
			}
		}
	}
	// notify
	if cfg.Components.Notify != nil {
		if prefixes, ok := componentPrefixMap[moduleNameNotify]; ok {
			for _, prefix := range prefixes {
				pathPrefixProxyMap[prefix] = *cfg.Components.Notify
			}
		}
	}
	// monitor
	if cfg.Components.Monitor != nil {
		if prefixes, ok := componentPrefixMap[moduleNameMonitor]; ok {
			for _, prefix := range prefixes {
				pathPrefixProxyMap[prefix] = *cfg.Components.Monitor
			}
		}
	}
	// auth
	if cfg.Components.Auth != nil {
		if prefixes, ok := componentPrefixMap[moduleNameAuth]; ok {
			for _, prefix := range prefixes {
				pathPrefixProxyMap[prefix] = *cfg.Components.Auth
			}
		}
	}
	// registry
	if cfg.Components.Registry != nil {
		if prefixes, ok := componentPrefixMap[moduleNameRegistry]; ok {
			for _, prefix := range prefixes {
				pathPrefixProxyMap[prefix] = *cfg.Components.Registry
			}
		}
	}
	return pathPrefixProxyMap
}
