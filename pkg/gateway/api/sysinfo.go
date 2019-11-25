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
	"fmt"
	"github.com/emicklei/go-restful"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
	"tkestack.io/tke/api/auth"
	authv1 "tkestack.io/tke/api/auth/v1"
	"tkestack.io/tke/api/business"
	businessv1 "tkestack.io/tke/api/business/v1"
	"tkestack.io/tke/api/monitor"
	monitorv1 "tkestack.io/tke/api/monitor/v1"
	"tkestack.io/tke/api/notify"
	notifyv1 "tkestack.io/tke/api/notify/v1"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	registryv1 "tkestack.io/tke/api/registry/v1"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
)

type GroupVersion struct {
	GroupName string `json:"groupName"`
	Version   string `json:"version"`
}

// Components contains the service address of each module in TKE
type Components struct {
	Platform *GroupVersion `json:"platform,omitempty"`
	Business *GroupVersion `json:"business,omitempty"`
	Monitor  *GroupVersion `json:"monitor,omitempty"`
	Notify   *GroupVersion `json:"notify,omitempty"`
	Auth     *GroupVersion `json:"auth,omitempty"`
	Registry *GroupVersion `json:"registry,omitempty"`
}

type SysInfo struct {
	Components Components `json:"components"`
	Registry   *Registry  `json:"registry,omitempty"`
	Auth       *Auth      `json:"auth,omitempty"`
}

type Registry struct {
	DefaultTenant string `json:"defaultTenant,omitempty"`
	DomainSuffix  string `json:"domainSuffix,omitempty"`
}

type Auth struct {
	DefaultTenant string `json:"defaultTenant,omitempty"`
}

func registerSysInfoRoute(container *restful.Container, cfg *gatewayconfig.GatewayConfiguration) {
	ws := new(restful.WebService)
	ws.Path(fmt.Sprintf("/apis/%s/%s/sysinfo", GroupName, Version))
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)
	ws.Route(ws.
		GET("/").
		Doc("get system information of TKE").
		Operation("getSysInfo").
		Returns(http.StatusOK, "Ok", SysInfo{}).
		To(handleGetSysInfoFunc(cfg)))
	container.Add(ws)
}

func handleGetSysInfoFunc(cfg *gatewayconfig.GatewayConfiguration) func(*restful.Request, *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		cmp := Components{}
		if cfg.Components.Platform != nil {
			cmp.Platform = &GroupVersion{
				GroupName: platform.GroupName,
				Version:   platformv1.Version,
			}
		}
		if cfg.Components.Business != nil {
			cmp.Business = &GroupVersion{
				GroupName: business.GroupName,
				Version:   businessv1.Version,
			}
		}
		if cfg.Components.Notify != nil {
			cmp.Notify = &GroupVersion{
				GroupName: notify.GroupName,
				Version:   notifyv1.Version,
			}
		}
		if cfg.Components.Auth != nil {
			cmp.Auth = &GroupVersion{
				GroupName: auth.GroupName,
				Version:   authv1.Version,
			}
		}
		if cfg.Components.Monitor != nil {
			cmp.Monitor = &GroupVersion{
				GroupName: monitor.GroupName,
				Version:   monitorv1.Version,
			}
		}
		if cfg.Components.Registry != nil {
			cmp.Registry = &GroupVersion{
				GroupName: registryv1.GroupName,
				Version:   registryv1.Version,
			}
		}
		sysInfo := SysInfo{
			Components: cmp,
		}
		if cfg.Auth != nil {
			sysInfo.Auth = &Auth{
				DefaultTenant: cfg.Auth.DefaultTenant,
			}
		}
		if cfg.Registry != nil {
			sysInfo.Registry = &Registry{
				DefaultTenant: cfg.Registry.DefaultTenant,
				DomainSuffix:  cfg.Registry.DomainSuffix,
			}
		}
		responsewriters.WriteRawJSON(http.StatusOK, sysInfo, response.ResponseWriter)
	}
}
