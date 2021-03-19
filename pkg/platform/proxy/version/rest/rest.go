/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

package rest

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/proxy"
)

type VersionProxyHandler struct {
	config *rest.Config
}

func NewVersionProxyHandler(config *rest.Config) *VersionProxyHandler {
	return &VersionProxyHandler{
		config: config,
	}
}

// Install install /version route
// Please confirm /version is not in pkg/apiserver/filter/authentication.go defaultIgnoreAuthPathPrefixes
func (s *VersionProxyHandler) Install(c *restful.Container) {
	versionWS := new(restful.WebService)
	versionWS.Path("/version")
	versionWS.Doc("git code version from which this is built")
	versionWS.Route(
		versionWS.GET("/").To(s.handle).
			Doc("get the code version").
			Operation("getVersion").
			Produces(restful.MIME_JSON).
			Consumes(restful.MIME_JSON))

	c.Add(versionWS)
}

// /version route process
func (s *VersionProxyHandler) handle(req *restful.Request, resp *restful.Response) {
	platformClient := platforminternalclient.NewForConfigOrDie(s.config)
	client, err := proxy.ClientSet(req.Request.Context(), platformClient)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, err.Error())
		return
	}
	version, err := client.Discovery().ServerVersion()

	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, err.Error())
		return
	}

	resp.WriteHeaderAndEntity(http.StatusOK, version)
}
