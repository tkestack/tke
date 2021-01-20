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

package route

import (
	"github.com/emicklei/go-restful"
)

// Routing routing interface
type Routing interface {
	AddToWebService(ws *restful.WebService)
}

// ClusterHandler cluster rest handler
type ClusterHandler interface {
	Routing
	Get(req *restful.Request, resp *restful.Response)
	List(req *restful.Request, resp *restful.Response)
}

// ConfigHandler config rest handler
type ConfigHandler interface {
	Routing
	ListIstioSupportedVersions(req *restful.Request, resp *restful.Response)
}

// MonitorHandler monitor rest handler
type MonitorHandler interface {
	Routing
	GetTracingMetricsData(req *restful.Request, resp *restful.Response)
	GetMonitorMetricsData(req *restful.Request, resp *restful.Response)
	GetMonitorTopologyData(req *restful.Request, resp *restful.Response)
}