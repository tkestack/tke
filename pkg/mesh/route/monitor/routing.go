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

package monitor

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	tcmeshclient "tkestack.io/tke/pkg/mesh/external/tcmesh"
	"tkestack.io/tke/pkg/mesh/route"
	"tkestack.io/tke/pkg/mesh/services"
	monitorservice "tkestack.io/tke/pkg/mesh/services/monitor"
	"tkestack.io/tke/pkg/mesh/services/rest"
)

type monitorHandler struct {
	monitorService services.MonitorService
}

func New(tcmeshClient *tcmeshclient.Client) route.MonitorHandler {
	return &monitorHandler{
		monitorService: monitorservice.New(tcmeshClient),
	}
}

func (h *monitorHandler) AddToWebService(ws *restful.WebService) {
	ws.Route(
		ws.POST("/meshes/{meshName}/namespaces/{namespaceName}/tracings").
			To(h.GetTracingMetricsData).
			Param(ws.PathParameter("meshName", "mesh name").
				DataType("string").Required(true)).
			Param(ws.PathParameter("namespaceName", "namespace name").
				DataType("string").Required(true)).
			Operation("getTracingMetricsData").
			Doc("Get tracing metrics data").
			Returns(http.StatusOK, "Get", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)

	ws.Route(
		ws.POST("/meshes/{meshName}/namespaces/{namespaceName}/metrics").
			To(h.GetMonitorMetricsData).
			Param(ws.PathParameter("meshName", "mesh name").
				DataType("string").Required(true)).
			Param(ws.PathParameter("namespaceName", "namespace name").
				DataType("string").Required(true)).
			Operation("getMonitorMetricsData").
			Doc("Get monitor metrics data").
			Returns(http.StatusOK, "Get", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)

	ws.Route(
		ws.POST("/meshes/{meshName}/namespaces/{namespaceName}/topologies").
			To(h.GetMonitorTopologyData).
			Param(ws.PathParameter("meshName", "mesh name").
				DataType("string").Required(true)).
			Param(ws.PathParameter("namespaceName", "namespace name").
				DataType("string").Required(true)).
			Operation("getMonitorTopologyData").
			Doc("Get monitor topology data").
			Returns(http.StatusOK, "Get", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
}

func (h *monitorHandler) GetTracingMetricsData(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest

	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	metricQuery := rest.MetricQuery{}

	err := req.ReadEntity(&metricQuery)
	if err != nil {
		result.Err = errors.Wrapf(err, "decode request").Error()
		return
	}

	data, err := h.monitorService.GetTracingMetricsData(req.Request.Context(), metricQuery)
	if err != nil {
		result.Err = err.Error()
		return
	}

	result.Data = data
	result.Result = true
	status = http.StatusOK
}

func (h *monitorHandler) GetMonitorMetricsData(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest

	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	metricQuery := new(rest.MetricQuery)

	err := req.ReadEntity(metricQuery)
	if err != nil {
		result.Err = errors.Wrapf(err, "decode request").Error()
		return
	}

	data, err := h.monitorService.GetMonitorMetricsData(req.Request.Context(), metricQuery)
	if err != nil {
		result.Err = err.Error()
		return
	}

	result.Data = data
	result.Result = true
	status = http.StatusOK
}

func (h *monitorHandler) GetMonitorTopologyData(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest

	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	topoQuery := new(rest.TopoQuery)

	err := req.ReadEntity(topoQuery)
	if err != nil {
		result.Err = errors.Wrapf(err, "decode request").Error()
		return
	}

	data, err := h.monitorService.GetMonitorTopologyData(req.Request.Context(), topoQuery)
	if err != nil {
		result.Err = err.Error()
		return
	}

	result.Data = data
	result.Result = true
	status = http.StatusOK
}
