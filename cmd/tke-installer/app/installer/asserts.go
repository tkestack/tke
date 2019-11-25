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
	"tkestack.io/tke/cmd/tke-installer/assets"

	"github.com/emicklei/go-restful"
)

// AssertsResource is the REST layer
type AssertsResource struct {
}

// NewAssertsResource create a AssertsResource
func NewAssertsResource() *AssertsResource {
	return new(AssertsResource)
}

// WebService creates a new service that can handle REST requests for AssertsResource
func (c *AssertsResource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/").To(func(req *restful.Request, resp *restful.Response) {
		assets.ServeHTTP(resp.ResponseWriter, req.Request)
	}))
	ws.Route(ws.GET("/{path:*}").To(func(req *restful.Request, resp *restful.Response) {
		assets.ServeHTTP(resp.ResponseWriter, req.Request)
	}))

	return ws
}
