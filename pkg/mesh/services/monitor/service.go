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
	"context"

	tcmeshclient "tkestack.io/tke/pkg/mesh/external/tcmesh"
	"tkestack.io/tke/pkg/mesh/services"
	"tkestack.io/tke/pkg/mesh/services/rest"
)

const (
	IstioServiceIngressGateway = "istio-ingressgateway_istio-system"
)

type monitorService struct {
	tcmeshClient *tcmeshclient.Client
}

var _ services.MonitorService = &monitorService{}

func New(tcmeshClient *tcmeshclient.Client) services.MonitorService {

	return &monitorService{
		tcmeshClient: tcmeshClient,
	}
}

func (s *monitorService) GetTracingMetricsData(
	ctx context.Context, query rest.MetricQuery,
) (*rest.MetricData, error) {

	uri := "/api/metrics/monitor/query"
	var response *rest.MetricData
	err := s.tcmeshClient.RESTClient().Post().RequestURI(uri).Body(query).Do(ctx).Into(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *monitorService) GetMonitorMetricsData(
	ctx context.Context, query *rest.MetricQuery,
) (*rest.MetricData, error) {

	uri := "/api/metrics/tracing/query"
	var response *rest.MetricData
	err := s.tcmeshClient.RESTClient().Post().RequestURI(uri).Body(query).Do(ctx).Into(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *monitorService) GetMonitorTopologyData(
	ctx context.Context, query *rest.TopoQuery,
) (*rest.TopoData, error) {

	uri := "/api/metrics/monitor/topology"
	var response *rest.TopoData
	err := s.tcmeshClient.RESTClient().Post().RequestURI(uri).Body(query).Do(ctx).Into(response)
	if err != nil {
		return nil, err
	}

	// filter by app
	result, err := s.filterTopologyByApp(response, query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *monitorService) filterTopologyByApp(
	topoData *rest.TopoData, request *rest.TopoQuery,
) (*rest.TopoData, error) {

	result := &rest.TopoData{
		Nodes: make(map[string]rest.Node),
		Edges: []rest.Edge{},
	}

	if request.App == "" {
		return topoData, nil
	}

	return result, nil
}

//func inApp(serviceNames []string, name string) bool {
//	tmpSvcName := getServiceName(name)
//	for _, svcName := range serviceNames {
//		if tmpSvcName == svcName {
//			return true
//		}
//	}
//
//	return false
//}
//
//// omitted the namespace name
//func getServiceName(origin string) string {
//	pos := strings.LastIndex(origin, "_")
//	return origin[0:pos]
//}
//
//func getSvcName(node rest.Node) string {
//	var name string
//	if node.Type == "service" {
//		name = node.Name
//	} else {
//		// delete the "service_" prefix
//		name = node.ServiceNodeID[8:]
//	}
//	return name
//}
