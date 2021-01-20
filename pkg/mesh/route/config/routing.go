/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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
 *
 */

package config

import (
	"net/http"

	"github.com/emicklei/go-restful"
	meshconfig "tkestack.io/tke/pkg/mesh/apis/config"
	"tkestack.io/tke/pkg/mesh/services/rest"
)

type configHandler struct {
	meshConfig meshconfig.MeshConfiguration
}

func New(c meshconfig.MeshConfiguration) *configHandler {
	return &configHandler{
		meshConfig: c,
	}
}

func (c *configHandler) AddToWebService(ws *restful.WebService) {
	ws.Route(
		ws.GET("/config/istioversions").
		To(c.ListIstioSupportedVersions).
		// Supported Operation's Prefix
		//  "get", "log", "read", "replace", "patch", "delete", "deletecollection",
		//  "watch", "connect", "proxy", "list", "create", "patch"
		Operation("listIstioVersions").
		Doc("List istio supported versions").
		Returns(http.StatusOK, "List", rest.Response{}).
		Returns(http.StatusBadRequest, "Error", rest.Response{}).
		Returns(http.StatusNotFound, "Not Found", rest.Response{}).
		Produces(restful.MIME_JSON),
	)
}

func (c *configHandler) ListIstioSupportedVersions(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest

	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	result.Data = c.meshConfig.Istio.SupportedVersion
	result.Result = true
	status = http.StatusOK
}
