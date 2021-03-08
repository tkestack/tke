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
 *
 */

package clusters

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	meshconfig "tkestack.io/tke/pkg/mesh/apis/config"
	"tkestack.io/tke/pkg/mesh/route"
	"tkestack.io/tke/pkg/mesh/services"
	"tkestack.io/tke/pkg/mesh/services/rest"
	"tkestack.io/tke/pkg/mesh/util/web"
)

// clusterHandler cluster rest handler implements
type clusterHandler struct {
	clusterService services.ClusterService
	istioService   services.IstioService
	config         meshconfig.MeshConfiguration
}

// New new ClusterHandler
func New(config meshconfig.MeshConfiguration, clusterService services.ClusterService,
	istioService services.IstioService) route.ClusterHandler {

	return &clusterHandler{
		clusterService: clusterService,
		istioService:   istioService,
		config:         config,
	}
}

func (h *clusterHandler) AddToWebService(ws *restful.WebService) {
	// ========== K8S Proxy API START ===========
	// "get", "log", "read", "replace", "patch", "delete", "deletecollection",
	// "watch", "connect", "proxy", "list", "create", "patch"
	var HTTPMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	proxyPath := fmt.Sprintf("/%s/{%s}/proxy/{path:*}", web.ClusterPrefix, web.ClusterKey)
	for _, method := range HTTPMethods {

		route := ws.Method(method).Path(proxyPath)
		route.
			To(h.Proxy).
			Operation(fmt.Sprintf("proxy%sKubeResources", method)).
			Doc("proxy kubernetes resources").
			Returns(http.StatusOK, "List", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Param(ws.PathParameter(web.ClusterKey, "cluster name").DataType("string")).
			Param(ws.PathParameter("path", "kube apiserver path").DataType("string"))

		ws.Route(route)
	}
	// ========== K8S Proxy API END ===========

	// ========== K8S API START ==========
	ws.Route(
		ws.GET(fmt.Sprintf("/%s", web.ClusterPrefix)).
			To(h.List).
			Operation("listCluster").
			Doc("List cluster from tke").
			Returns(http.StatusOK, "List", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
	ws.Route(
		ws.GET(fmt.Sprintf("/%s/{%s}", web.ClusterPrefix, web.ClusterKey)).
			To(h.Get).
			Param(ws.PathParameter(web.ClusterKey, "cluster name").DataType("string").Required(true)).
			Operation("getCluster").
			Doc("Get cluster from tke").
			Returns(http.StatusOK, "Get", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
	ws.Route(
		ws.GET(fmt.Sprintf("/%s/{%s}/%s", web.ClusterPrefix, web.ClusterKey, web.NamespacePrefix)).
			To(h.ListNamespaces).
			Param(ws.PathParameter(web.ClusterKey, "cluster name").DataType("string").Required(true)).
			Operation("listClusterNamespaces").
			Doc("list cluster namespaces").
			Returns(http.StatusOK, "List", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
	// ========== K8S API END ==========

	// ========== Istio API START ==========
	ws.Route(
		ws.GET(
			fmt.Sprintf("/%s/{%s}/%s/{%s}/all",
				web.ClusterPrefix, web.ClusterKey,
				web.NamespacePrefix, web.NamespaceKey,
			)).
			To(h.ListAll).
			Param(ws.PathParameter(web.ClusterKey, "cluster name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.NamespaceKey, "namespace name").DataType("string").Required(true)).
			Param(ws.QueryParameter(web.LabelSelectorQuery, "resource labels selector").DataType("string").
				DefaultValue("app=test,version=1.0.0")).
			Operation("listAllIstioResources").
			Doc("list all istio resources").
			Returns(http.StatusOK, "Get", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)

	ws.Route(
		ws.GET(
			fmt.Sprintf("/%s/{%s}/%s/{%s}/{%s}",
				web.ClusterPrefix, web.ClusterKey,
				web.NamespacePrefix, web.NamespaceKey, web.IstioResourceKind,
			)).
			To(h.ListIstioResources).
			Param(ws.PathParameter(web.ClusterKey, "cluster name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.NamespaceKey, "namespace name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.IstioResourceKind, "istio resource kind").DataType("string").Required(true).
				AllowableValues(map[string]string{
					web.GatewayPrefix:         web.GatewayPrefix,
					web.VirtualservicePrefix:  web.VirtualservicePrefix,
					web.DestinationrulePrefix: web.DestinationrulePrefix,
					web.ServiceentriePrefix:   web.ServiceentriePrefix,
					web.SidecarPrefix:         web.SidecarPrefix,
					web.EnvoyfilterPrefix:     web.EnvoyfilterPrefix,
					web.WorkloadentryPrefix:   web.WorkloadentryPrefix,
				})).
			Param(ws.QueryParameter(web.LabelSelectorQuery, "resource labels selector").DataType("string").
				DefaultValue("app=test,version=1.0.0")).
			Operation("listIstioResources").
			Doc("list istio resources").
			Returns(http.StatusOK, "Get", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
	ws.Route(
		ws.GET(
			fmt.Sprintf("/%s/{%s}/%s/{%s}/{%s}/{%s}",
				web.ClusterPrefix, web.ClusterKey,
				web.NamespacePrefix, web.NamespaceKey,
				web.IstioResourceKind, web.IstioResourceKey,
			)).
			To(h.GetIstioResource).
			Param(ws.PathParameter(web.ClusterKey, "cluster name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.NamespaceKey, "namespace name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.IstioResourceKind, "istio resource kind").DataType("string").Required(true).
				AllowableValues(map[string]string{
					web.GatewayPrefix:         web.GatewayPrefix,
					web.VirtualservicePrefix:  web.VirtualservicePrefix,
					web.DestinationrulePrefix: web.DestinationrulePrefix,
					web.ServiceentriePrefix:   web.ServiceentriePrefix,
					web.SidecarPrefix:         web.SidecarPrefix,
					web.EnvoyfilterPrefix:     web.EnvoyfilterPrefix,
					web.WorkloadentryPrefix:   web.WorkloadentryPrefix,
				})).
			Param(ws.PathParameter(web.IstioResourceKey, "istio resource name").DataType("string").Required(true)).
			Operation("getIstioResource").
			Doc("get istio resource").
			Returns(http.StatusOK, "Get", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
	ws.Route(
		ws.POST(
			fmt.Sprintf("/%s/{%s}/%s/{%s}/{%s}",
				web.ClusterPrefix, web.ClusterKey,
				web.NamespacePrefix, web.NamespaceKey,
				web.IstioResourceKind,
			)).
			To(h.CreateIstioResource).
			Param(ws.PathParameter(web.ClusterKey, "cluster name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.NamespaceKey, "namespace name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.IstioResourceKind, "istio resource kind").DataType("string").Required(true).
				AllowableValues(map[string]string{
					web.GatewayPrefix:         web.GatewayPrefix,
					web.VirtualservicePrefix:  web.VirtualservicePrefix,
					web.DestinationrulePrefix: web.DestinationrulePrefix,
					web.ServiceentriePrefix:   web.ServiceentriePrefix,
					web.SidecarPrefix:         web.SidecarPrefix,
					web.EnvoyfilterPrefix:     web.EnvoyfilterPrefix,
					web.WorkloadentryPrefix:   web.WorkloadentryPrefix,
				})).
			Operation("createIstioResource").
			Doc("create istio resource").
			Returns(http.StatusOK, "Created", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
	ws.Route(
		ws.DELETE(
			fmt.Sprintf("/%s/{%s}/%s/{%s}/{%s}/{%s}",
				web.ClusterPrefix, web.ClusterKey,
				web.NamespacePrefix, web.NamespaceKey,
				web.IstioResourceKind, web.IstioResourceKey,
			)).
			To(h.DeleteIstioResource).
			Param(ws.PathParameter(web.ClusterKey, "cluster name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.NamespaceKey, "namespace name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.IstioResourceKind, "istio resource kind").DataType("string").Required(true).
				AllowableValues(map[string]string{
					web.GatewayPrefix:         web.GatewayPrefix,
					web.VirtualservicePrefix:  web.VirtualservicePrefix,
					web.DestinationrulePrefix: web.DestinationrulePrefix,
					web.ServiceentriePrefix:   web.ServiceentriePrefix,
					web.SidecarPrefix:         web.SidecarPrefix,
					web.EnvoyfilterPrefix:     web.EnvoyfilterPrefix,
					web.WorkloadentryPrefix:   web.WorkloadentryPrefix,
				})).
			Param(ws.PathParameter(web.IstioResourceKey, "istio resource name").DataType("string").Required(true)).
			Operation("deleteIstioResource").
			Doc("delete istio resource").
			Returns(http.StatusOK, "Deleted", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
	// ========== Istio API END ==========

	// ========== Mesh API START ==========
	ws.Route(
		ws.POST(
			fmt.Sprintf("/%s/{%s}/%s/{%s}/northtraffic/%s",
				web.ClusterPrefix, web.ClusterKey,
				web.NamespacePrefix, web.NamespaceKey,
				web.GatewayPrefix,
			)).
			To(h.GetNorthTrafficGateway).
			Param(ws.PathParameter(web.ClusterKey, "cluster name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.NamespaceKey, "namespace name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.IstioResourceKey, "istio resource name").DataType("string").Required(true)).
			Operation("getNorthTrafficGateway").
			Doc("get istio north traffic gateway").
			Returns(http.StatusOK, "Get", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
	ws.Route(
		ws.POST(
			fmt.Sprintf("/%s/{%s}/%s/{%s}/northtraffic/%s",
				web.ClusterPrefix, web.ClusterKey,
				web.NamespacePrefix, web.NamespaceKey,
				web.GatewayPrefix,
			)).
			To(h.CreateNorthTrafficGateway).
			Param(ws.PathParameter(web.ClusterKey, "cluster name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.NamespaceKey, "namespace name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.IstioResourceKey, "istio resource name").DataType("string").Required(true)).
			Operation("createNorthTrafficGateway").
			Doc("create istio north traffic gateway").
			Returns(http.StatusOK, "Created", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
	ws.Route(
		ws.PATCH(
			fmt.Sprintf("/%s/{%s}/%s/{%s}/northtraffic/%s/{%s}",
				web.ClusterPrefix, web.ClusterKey,
				web.NamespacePrefix, web.NamespaceKey,
				web.GatewayPrefix, web.GatewayKey,
			)).
			To(h.UpdateNorthTrafficGateway).
			Param(ws.PathParameter(web.ClusterKey, "cluster name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.NamespaceKey, "namespace name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.GatewayKey, "istio gateway name").DataType("string").Required(true)).
			Operation("patchNorthTrafficGateway").
			Doc("update istio north traffic gateway").
			Returns(http.StatusOK, "Updated", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
	ws.Route(
		ws.DELETE(
			fmt.Sprintf("/%s/{%s}/%s/{%s}/northtraffic/%s/{%s}",
				web.ClusterPrefix, web.ClusterKey,
				web.NamespacePrefix, web.NamespaceKey,
				web.GatewayPrefix, web.GatewayKey,
			)).
			To(h.DeleteNorthTrafficGateway).
			Param(ws.PathParameter(web.ClusterKey, "cluster name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.NamespaceKey, "namespace name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.GatewayKey, "istio gateway name").DataType("string").Required(true)).
			Operation("deleteNorthTrafficGateway").
			Doc("delete istio north traffic gateway").
			Returns(http.StatusOK, "Updated", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
	// ========== Mesh API END ==========
}
