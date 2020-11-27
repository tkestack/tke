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
	"tkestack.io/tke/api/application"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/api/business"
	"tkestack.io/tke/api/logagent"
	"tkestack.io/tke/api/monitor"
	"tkestack.io/tke/api/notify"
	"tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	auditapi "tkestack.io/tke/pkg/audit/api"
	authapiserver "tkestack.io/tke/pkg/auth/apiserver"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
	"tkestack.io/tke/pkg/gateway/proxy/handler/frontproxy"
	"tkestack.io/tke/pkg/gateway/proxy/handler/passthrough"
	"tkestack.io/tke/pkg/gateway/proxy/handler/rewriteproxy"
	platformapiserver "tkestack.io/tke/pkg/platform/apiserver"
	"tkestack.io/tke/pkg/registry/chartmuseum"
	"tkestack.io/tke/pkg/registry/distribution"
	"tkestack.io/tke/pkg/registry/harbor"
	"tkestack.io/tke/pkg/util/log"
)

const (
	apiPrefix      = "/apis"
	openapiPrefix  = "/openapi"
	openapiVersion = "v2"
)

type moduleName string

const (
	moduleNamePlatform    moduleName = "platform"
	moduleNameBusiness    moduleName = "business"
	moduleNameNotify      moduleName = "notify"
	moduleNameRegistry    moduleName = "registry"
	moduleNameAuth        moduleName = "auth"
	moduleNameMonitor     moduleName = "monitor"
	moduleNameLogagent    moduleName = "logagent"
	moduleNameAudit       moduleName = "audit"
	moduleNameApplication moduleName = "application"
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
				prefix:    fmt.Sprintf("%s/%s/", apiPrefix, notify.GroupName),
				protected: true,
			},
		},
		moduleNameAuth: {
			modulePath{
				prefix:    fmt.Sprintf("%s/%s/", apiPrefix, auth.GroupName),
				protected: true,
			},
			modulePath{
				prefix:    authapiserver.AuthPath,
				protected: false,
			},
			modulePath{
				prefix:    fmt.Sprintf("%s/", authapiserver.APIKeyPasswordPath),
				protected: false,
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
				prefix:    chartmuseum.PathPrefix,
				protected: false,
			},
			modulePath{
				prefix:    harbor.RegistryPrefix,
				protected: false,
			},
			modulePath{
				prefix:    harbor.ChartAPIPrefix,
				protected: false,
			},
			modulePath{
				prefix:    harbor.AuthPrefix,
				protected: false,
			},
			modulePath{
				prefix:    harbor.ChartPrefix,
				protected: false,
			},
			modulePath{
				prefix:    fmt.Sprintf("%s/%s/", apiPrefix, registry.GroupName),
				protected: true,
			},
		},
		moduleNameLogagent: {
			modulePath{
				prefix:    fmt.Sprintf("%s/%s/", apiPrefix, logagent.GroupName),
				protected: true,
			},
		},
		moduleNameAudit: {
			modulePath{
				prefix:    fmt.Sprintf("%s/%s/%s/events/sink/", apiPrefix, auditapi.GroupName, auditapi.Version),
				protected: false,
			},
			modulePath{
				prefix:    fmt.Sprintf("%s/%s/%s/events/list/", apiPrefix, auditapi.GroupName, auditapi.Version),
				protected: true,
			},
			modulePath{
				prefix:    fmt.Sprintf("%s/%s/%s/events/listFieldValues/", apiPrefix, auditapi.GroupName, auditapi.Version),
				protected: true,
			},
		},
		moduleNameApplication: {
			modulePath{
				prefix:    fmt.Sprintf("%s/%s/", apiPrefix, application.GroupName),
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
			log.Info("Registered reverse proxy of passthrough mode for backend component", log.String("pathPrefix", pathPrefix.prefix), log.Bool("protected", pathPrefix.protected), log.String("address", proxyComponent.Address))
			m.HandlePrefix(pathPrefix.prefix, handler)
		}
	}
	pathOpenapi := openapiProxy(cfg)
	for path, proxyComponent := range pathOpenapi {
		if proxyComponent.Passthrough != nil {
			handler, err := rewriteproxy.NewHandler(
				proxyComponent.Address,
				proxyComponent.Passthrough,
				path.protected,
				func(string) string { return fmt.Sprintf("%s/%s", openapiPrefix, openapiVersion) },
			)
			if err != nil {
				return err
			}
			log.Info("Registered openapi proxy for backend component", log.String("path", path.prefix), log.Bool("protected", path.protected), log.String("address", proxyComponent.Address))
			m.Handle(path.prefix, handler)
		}
	}
	// proxy /webhook to tke-notify-api for alert
	if cfg.Components.Notify != nil && cfg.Components.Notify.Passthrough != nil {
		handler, err := passthrough.NewHandler(cfg.Components.Notify.Address, cfg.Components.Notify.Passthrough, false)
		if err != nil {
			return err
		}
		log.Info("Registered reverse proxy of passthrough mode for backend component", log.String("path", "/webhook"), log.Bool("protected", false), log.String("address", cfg.Components.Notify.Address))
		m.Handle("/webhook", handler)
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
	//log agent
	if cfg.Components.LogAgent != nil {
		if prefixes, ok := componentPrefixMap[moduleNameLogagent]; ok {
			for _, prefix := range prefixes {
				pathPrefixProxyMap[prefix] = *cfg.Components.LogAgent
			}
		}
	}
	// audit
	if cfg.Components.Audit != nil {
		if prefixes, ok := componentPrefixMap[moduleNameAudit]; ok {
			for _, prefix := range prefixes {
				pathPrefixProxyMap[prefix] = *cfg.Components.Audit
			}
		}
	}
	// application
	if cfg.Components.Application != nil {
		if prefixes, ok := componentPrefixMap[moduleNameApplication]; ok {
			for _, prefix := range prefixes {
				pathPrefixProxyMap[prefix] = *cfg.Components.Application
			}
		}
	}
	return pathPrefixProxyMap
}

func openapiProxy(cfg *gatewayconfig.GatewayConfiguration) map[modulePath]gatewayconfig.Component {
	newOpenapiPath := func(name moduleName) modulePath {
		return modulePath{
			prefix:    fmt.Sprintf("/tke-%s-api", name),
			protected: false,
		}
	}
	openapiProxyMap := make(map[modulePath]gatewayconfig.Component)
	// platform
	if cfg.Components.Platform != nil {
		openapiProxyMap[newOpenapiPath(moduleNamePlatform)] = *cfg.Components.Platform
	}
	// business
	if cfg.Components.Business != nil {
		openapiProxyMap[newOpenapiPath(moduleNameBusiness)] = *cfg.Components.Business
	}
	// notify
	if cfg.Components.Notify != nil {
		openapiProxyMap[newOpenapiPath(moduleNameNotify)] = *cfg.Components.Notify
	}
	// monitor
	if cfg.Components.Monitor != nil {
		openapiProxyMap[newOpenapiPath(moduleNameMonitor)] = *cfg.Components.Monitor
	}
	// auth
	if cfg.Components.Auth != nil {
		openapiProxyMap[newOpenapiPath(moduleNameAuth)] = *cfg.Components.Auth
	}
	// registry
	if cfg.Components.Registry != nil {
		openapiProxyMap[newOpenapiPath(moduleNameRegistry)] = *cfg.Components.Registry
	}
	// logagent
	if cfg.Components.LogAgent != nil {
		openapiProxyMap[newOpenapiPath(moduleNameLogagent)] = *cfg.Components.LogAgent
	}
	// audit
	if cfg.Components.Audit != nil {
		openapiProxyMap[newOpenapiPath(moduleNameAudit)] = *cfg.Components.Audit
	}
	return openapiProxyMap
}
