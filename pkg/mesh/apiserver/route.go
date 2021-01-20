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

package apiserver

import (
	"github.com/emicklei/go-restful"
	clustersroute "tkestack.io/tke/pkg/mesh/route/clusters"
	configroute "tkestack.io/tke/pkg/mesh/route/config"
	"tkestack.io/tke/pkg/mesh/route/meshes"
	monitorroute "tkestack.io/tke/pkg/mesh/route/monitor"
	clustersservice "tkestack.io/tke/pkg/mesh/services/clusters"
	istioservice "tkestack.io/tke/pkg/mesh/services/istio"
	"tkestack.io/tke/pkg/mesh/services/mesh"
)

const MeshAPIPrefix = "/apis/v1/mesh"

type Route struct {
	config completedConfig
}

func (r *Route) WebService() *restful.WebService {
	// go-restful tracing logs
	// restful.EnableTracing(true)

	ws := new(restful.WebService)

	ws.Path(MeshAPIPrefix)

	var (
		config         = *r.config.ExtraConfig.MeshConfig
		platformClient = r.config.ExtraConfig.PlatformClient
		clusterClients = r.config.ExtraConfig.ClusterClients
		tcmClient      = r.config.ExtraConfig.TCMeshClient

		clusterService = clustersservice.New(platformClient, clusterClients)
		istioService   = istioservice.New(config, clusterClients)
		meshService    = mesh.New(config, tcmClient, clusterClients)
	)

	configroute.New(config).AddToWebService(ws)
	clustersroute.New(config, clusterService, istioService).AddToWebService(ws)
	monitorroute.New(tcmClient).AddToWebService(ws)
	meshes.New(config, clusterService, istioService, meshService).AddToWebService(ws)

	return ws
}
