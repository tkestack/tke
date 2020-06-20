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

package installer

import (
	"net/http"
	"path"

	"github.com/emicklei/go-restful"
)

type AssetsResource struct {
	RootDir string
}

func NewAssertsResource() *AssetsResource {
	c := new(AssetsResource)
	c.RootDir = "assets/"

	return c
}

func (c *AssetsResource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/").To(func(req *restful.Request, resp *restful.Response) {
		http.ServeFile(
			resp.ResponseWriter,
			req.Request,
			c.RootDir)
	}))
	ws.Route(ws.GET("/{subpath:*}").To(func(req *restful.Request, resp *restful.Response) {
		http.ServeFile(
			resp.ResponseWriter,
			req.Request,
			path.Join(c.RootDir, req.PathParameter("subpath")))
	}))

	return ws
}
