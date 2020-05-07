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
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/apis/audit"
	"strconv"
	"strings"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/api/business"
	"tkestack.io/tke/api/monitor"
	"tkestack.io/tke/api/notify"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/api/registry"
	auditconfig "tkestack.io/tke/pkg/audit/apis/config"
	"tkestack.io/tke/pkg/audit/storage"
	"tkestack.io/tke/pkg/audit/storage/es"
	"tkestack.io/tke/pkg/audit/storage/types"
	"tkestack.io/tke/pkg/util/log"
)

// GroupName is the api group name for audit.
const GroupName = "audit.tkestack.io"

// Version is the api version for audit.
const Version = "v1"

// ClusterControlPlane is the cluster name the tkestack control-planes like tke-platform-api will use to report audit events
const ClusterControlPlane = "control-plane"

var controlPlaneGroups sets.String

func init() {
	controlPlaneGroups = sets.NewString(
		platform.GroupName,
		registry.GroupName,
		notify.GroupName,
		monitor.GroupName,
		business.GroupName,
		auth.GroupName,
	)
}

// RegisterRoute is used to register prefix path routing matches for all
// configured backend components.
func RegisterRoute(container *restful.Container, cfg *auditconfig.AuditConfiguration) error {
	return registerAuditRoute(container, cfg)
}

func registerAuditRoute(container *restful.Container, cfg *auditconfig.AuditConfiguration) error {
	ws := new(restful.WebService)
	ws.Path(fmt.Sprintf("/apis/%s/%s/events", GroupName, Version))
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)
	cli, err := es.NewStorage(cfg.Storage.ElasticSearch)
	if err != nil {
		return err
	}
	ws.Route(ws.POST("/sink/{clusterName}").To(sinkEvents(cli)).
		Operation("createEventsByCluster").
		Doc("Create new audit events").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/list").To(listEvents(cli)).
		Operation("listEvents").
		Doc("Create new audit events").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/listFieldValues").To(listFieldValues(cli)).
		Operation("listFieldValues").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	container.Add(ws)
	return nil
}

func sinkEvents(cli storage.AuditStorage) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		var eventList audit.EventList
		err := request.ReadEntity(&eventList)
		if err != nil {
			log.Infof("failed read events: %v", err)
			response.Write([]byte("failed"))
			return
		}
		clusterName := request.PathParameter("clusterName")
		events := types.ConvertEvents(eventList.Items)
		for _, event := range events {
			event.ClusterName = clusterName
		}
		events = eventsFilter(events)
		err = cli.Save(events)
		if err != nil {
			log.Errorf("failed save events: %v", err)
		}
		response.Write([]byte("success"))
	}
}

type filterFunc func(e *types.Event) bool

func controlPlaneFilter(e *types.Event) bool {
	if e.ClusterName == ClusterControlPlane {
		if !controlPlaneGroups.Has(e.APIGroup) {
			return false
		}
	}
	return true
}

func userKubeletFilter(e *types.Event) bool {
	if strings.HasPrefix(e.UserName, "system:node:") {
		return false
	}
	return true
}

func kubesystemServiceAccountFilter(e *types.Event) bool {
	if strings.HasPrefix(e.UserName, "system:serviceaccount:kube-system:") {
		return false
	}
	return true
}

func tkePlatformControllerFilter(e *types.Event) bool {
	if e.UserName == "admin" && e.Verb == "update" && e.Resource == "clusters" && strings.HasPrefix(e.UserAgent, "tke-platform-controller") {
		return false
	}
	return true
}

var eventFilters = []filterFunc{
	controlPlaneFilter,
	userKubeletFilter,
	kubesystemServiceAccountFilter,
	tkePlatformControllerFilter,
}

func eventFilter(e *types.Event) bool {
	for _, f := range eventFilters {
		if !f(e) {
			return false
		}
	}
	return true
}

func eventsFilter(events []*types.Event) []*types.Event {
	var result []*types.Event
	for i := range events {
		if eventFilter(events[i]) {
			result = append(result, events[i])
		}
	}
	return result
}

func listEvents(cli storage.AuditStorage) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		params := storage.QueryParameter{
			ClusterName: request.QueryParameter("cluster"),
			Namespace:   request.QueryParameter("namespace"),
			Resource:    request.QueryParameter("resource"),
			Name:        request.QueryParameter("name"),
			Query:       request.QueryParameter("query"),
			UserName:    request.QueryParameter("user"),
		}
		startTime := request.QueryParameter("startTime")
		if startTime != "" {
			if stime, err := strconv.ParseInt(startTime, 10, 64); err == nil {
				params.StartTime = stime
			}
		}
		endTime := request.QueryParameter("endTime")
		if endTime != "" {
			if etime, err := strconv.ParseInt(endTime, 10, 64); err == nil {
				params.EndTime = etime
			}
		}
		page := request.QueryParameter("pageIndex")
		size := request.QueryParameter("pageSize")
		if size != "" {
			if s, err := strconv.Atoi(size); err == nil && s > 0 {
				params.Size = s
			}
		}
		if params.Size == 0 {
			params.Size = 10
		}
		if page != "" {
			if p, err := strconv.Atoi(page); err == nil && p > 0 {
				params.Offset = (p - 1) * params.Size
			}
		}
		events, total, err := cli.Query(&params)
		if err != nil {
			log.Errorf("failed to query events: %v", err)
			response.WriteEntity(ResultStatus{Code: -1, Message: err.Error()})
		} else {
			response.WriteEntity(Pagination{ResultStatus: ResultStatus{Code: 0, Message: ""}, Total: total, Items: events})
		}
	}
}

func listFieldValues(cli storage.AuditStorage) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		result := cli.FieldValues()
		response.WriteEntity(result)
	}
}

// IgnoredAuthPathPrefixes returns a list of path prefixes that does not need to
// go through the built-in authentication and authorization middleware of apiserver.
func IgnoredAuthPathPrefixes() []string {
	return []string{
		fmt.Sprintf("/apis/%s/%s/events/sink", GroupName, Version),
	}
}

type ResultStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	ResultStatus `json:",inline"`
	Total        int            `json:"total"`
	Items        []*types.Event `json:"items"`
}
