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
	"net/http"
	"sort"
	"strconv"
	"strings"

	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/monitor/services"
	"tkestack.io/tke/pkg/monitor/services/alertmanager"
	"tkestack.io/tke/pkg/monitor/services/prometheus"
	"tkestack.io/tke/pkg/monitor/services/rest"
	"tkestack.io/tke/pkg/util/log"

	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	alarmPolicyPrefix = "alarmpolicies"
	clustersPrefix    = "clusters"
)

type processor struct {
	platformClient        platformversionedclient.PlatformV1Interface
	prometheusProcessor   services.RuleProcessor
	alertmanagerProcessor services.RouteProcessor
}

// NewProcessor returns a processor to handle prometheus rules changes
func NewProcessor(platformClient platformversionedclient.PlatformV1Interface) services.BackendConfigProcessor {
	return &processor{
		platformClient:        platformClient,
		prometheusProcessor:   prometheus.NewProcessor(platformClient),
		alertmanagerProcessor: alertmanager.NewProcessor(platformClient),
	}
}

func (h *processor) RegisterWebService(ws *restful.WebService) {
	policiesPattern := strings.Join([]string{"", clustersPrefix, "{clusterName}", alarmPolicyPrefix}, "/")
	policyPattern := strings.Join([]string{policiesPattern, "{alarmPolicy}"}, "/")

	ws.Route(
		ws.POST(policiesPattern).
			To(h.Create).
			Param(ws.PathParameter("clusterName", "cluster name").DataType("string").Required(true)).
			Operation("createAlarmPolicy").
			Doc("Create a alarm policy of tke").
			Returns(http.StatusOK, "Created", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Consumes(restful.MIME_JSON).
			Produces(restful.MIME_JSON),
	)

	ws.Route(
		ws.DELETE(policyPattern).
			To(h.Delete).
			Param(ws.PathParameter("clusterName", "cluster name").DataType("string").Required(true)).
			Param(ws.PathParameter("alarmPolicy", "alarm policy name").DataType("string").Required(true)).
			Operation("deleteAlarmPolicy").
			Doc("Delete a alarm policy of tke").
			Returns(http.StatusOK, "Deleted", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)

	ws.Route(
		ws.PUT(policyPattern).
			To(h.Update).
			Param(ws.PathParameter("clusterName", "cluster name").DataType("string").Required(true)).
			Param(ws.PathParameter("alarmPolicy", "alarm policy name").DataType("string").Required(true)).
			Operation("replaceAlarmPolicy").
			Doc("Update a alarm policy of tke, this will replace all sub-alarmMetrics in this policy").
			Returns(http.StatusOK, "Updated", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Consumes(restful.MIME_JSON).
			Produces(restful.MIME_JSON),
	)

	ws.Route(
		ws.GET(policyPattern).
			To(h.Get).
			Param(ws.PathParameter("clusterName", "cluster name").DataType("string").Required(true)).
			Param(ws.PathParameter("alarmPolicy", "alarm policy name").DataType("string").Required(true)).
			Operation("getAlarmPolicy").
			Doc("Get alarm policy from tke").
			Returns(http.StatusOK, "Get", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)

	ws.Route(
		ws.GET(policiesPattern).
			To(h.List).
			Param(ws.QueryParameter("page", "page number").DataType("string").DefaultValue("1").Required(true)).
			Param(ws.QueryParameter("page_size", "page size").DataType("string").DefaultValue("10").Required(true)).
			Param(ws.PathParameter("clusterName", "cluster name").DataType("string").Required(true)).
			Operation("getAllAlarmPolicies").
			Doc("Get all alarm policies from tke, which belongs to the given cluster").
			Returns(http.StatusOK, "Get", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Error", rest.Response{}).
			Produces(restful.MIME_JSON),
	)

	log.Infof("Register monitor web service")
}

func writeResult(method string, clusterName, entityName string, status int, result *rest.Response, resp *restful.Response) {
	if status == http.StatusOK {
		log.Infof("Successfully %s %s(%s)", method, clusterName, entityName)
	} else {
		log.Errorf("failed to %s %s(%s) due to %s", method, clusterName, entityName, result.Err)
	}

	if strings.Contains(result.Err, "not found") {
		status = http.StatusNotFound
	}

	_ = resp.WriteHeaderAndEntity(status, result)
}

func (h *processor) Create(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest

	var (
		entityName  string
		clusterName string
	)

	defer func() {
		writeResult("create", clusterName, entityName, status, result, resp)
	}()

	clusterName = req.PathParameter("clusterName")
	if clusterName == "" {
		result.Err = "empty clusterName"
		return
	}

	alarmPolicy := new(rest.AlarmPolicy)
	err := req.ReadEntity(alarmPolicy)
	if err != nil {
		result.Err = errors.Wrapf(err, "decode request").Error()
		return
	}

	err = alarmPolicy.Validate()
	if err != nil {
		result.Err = errors.Wrapf(err, "validate alarmPolicy failed").Error()
		return
	}

	entityName = alarmPolicy.AlarmPolicySettings.AlarmPolicyName

	ruleGroup := &v1.RuleGroup{
		Name:     entityName,
		Interval: alarmPolicy.GetInterval(),
		Rules:    []v1.Rule{},
	}
	if alarmPolicy.AlarmPolicySettings != nil {
		for i := range alarmPolicy.AlarmPolicySettings.AlarmMetrics {
			a := alarmPolicy.AlarmPolicySettings.AlarmMetrics[i]
			rule := v1.Rule{
				Alert:       a.MetricName,
				Expr:        intstr.FromString(a.GetExpr(alarmPolicy)),
				For:         a.GetFor(alarmPolicy.AlarmPolicySettings.StatisticsPeriod),
				Labels:      a.GetLabels(entityName, "1"),
				Annotations: a.GetAnnotations(alarmPolicy),
			}
			ruleGroup.Rules = append(ruleGroup.Rules, rule)
		}
	}
	err = h.prometheusProcessor.CreateGroup(clusterName, entityName, ruleGroup)
	if err != nil {
		result.Err = err.Error()
		return
	}
	result.Result = true
	status = http.StatusOK
}

func (h *processor) Delete(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest

	var (
		entityName  string
		clusterName string
	)

	defer func() {
		writeResult("delete", clusterName, entityName, status, result, resp)
	}()

	clusterName = req.PathParameter("clusterName")
	if clusterName == "" {
		result.Err = "empty clusterName"
		return
	}

	alarmPolicyName := req.PathParameter("alarmPolicy")
	if alarmPolicyName == "" {
		result.Err = "empty alarmPolicy"
		return
	}

	entityName = alarmPolicyName

	err := h.prometheusProcessor.DeleteGroup(clusterName, alarmPolicyName)
	if err != nil {
		result.Err = err.Error()
		return
	}

	result.Result = true
	status = http.StatusOK
}

func (h *processor) Update(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest

	var (
		entityName  string
		clusterName string
	)

	defer func() {
		writeResult("update", clusterName, entityName, status, result, resp)
	}()

	clusterName = req.PathParameter("clusterName")
	if clusterName == "" {
		result.Err = "empty clusterName"
		return
	}

	alarmPolicyName := req.PathParameter("alarmPolicy")
	if alarmPolicyName == "" {
		result.Err = "empty alarmPolicy"
		return
	}

	entityName = alarmPolicyName

	oldRuleGroup, err := h.prometheusProcessor.GetGroup(clusterName, alarmPolicyName)
	if err != nil {
		result.Err = err.Error()
		return
	}

	version := "1"
	if oldRuleGroup != nil && len(oldRuleGroup.Rules) != 0 {
		if v, ok := oldRuleGroup.Rules[0].Labels[rest.VersionKey]; ok {
			i, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				version = fmt.Sprintf("%d", i+1)
			}
		}
	}

	alarmPolicy := new(rest.AlarmPolicy)
	err = req.ReadEntity(alarmPolicy)
	if err != nil {
		result.Err = errors.Wrapf(err, "decode request").Error()
		return
	}

	err = alarmPolicy.Validate()
	if err != nil {
		result.Err = errors.Wrapf(err, "validate alarmPolicy failed").Error()
		return
	}

	ruleGroup := &v1.RuleGroup{
		Name:     alarmPolicyName,
		Interval: alarmPolicy.GetInterval(),
		Rules:    []v1.Rule{},
	}

	if alarmPolicy.AlarmPolicySettings != nil {
		for i := range alarmPolicy.AlarmPolicySettings.AlarmMetrics {
			r := alarmPolicy.AlarmPolicySettings.AlarmMetrics[i]
			rule := v1.Rule{
				Alert:       r.MetricName,
				Expr:        intstr.FromString(r.GetExpr(alarmPolicy)),
				For:         r.GetFor(alarmPolicy.AlarmPolicySettings.StatisticsPeriod),
				Labels:      r.GetLabels(entityName, version),
				Annotations: r.GetAnnotations(alarmPolicy),
			}
			ruleGroup.Rules = append(ruleGroup.Rules, rule)
		}
	}

	err = h.prometheusProcessor.UpdateGroup(clusterName, alarmPolicyName, ruleGroup)
	if err != nil {
		result.Err = err.Error()
		return
	}

	result.Result = true
	status = http.StatusOK
}

func (h *processor) Get(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest

	var (
		entityName  string
		clusterName string
	)

	defer func() {
		writeResult("get", clusterName, entityName, status, result, resp)
	}()

	clusterName = req.PathParameter("clusterName")
	if clusterName == "" {
		result.Err = "empty clusterName"
		return
	}

	alarmPolicyName := req.PathParameter("alarmPolicy")
	if alarmPolicyName == "" {
		result.Err = "empty alarmPolicy"
		return
	}
	entityName = alarmPolicyName

	ruleGroup, err := h.prometheusProcessor.GetGroup(clusterName, alarmPolicyName)
	if err != nil {
		result.Err = err.Error()
		return
	}

	alarmPolicy := rest.NewAlarmPolicyFromRuleGroup(ruleGroup)

	result.Data = alarmPolicy
	result.Result = true
	status = http.StatusOK
}

func (h *processor) List(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest

	var (
		entityName  string
		clusterName string
		page        int64
		pageSize    int64
	)

	defer func() {
		writeResult("list", clusterName, entityName, status, result, resp)
	}()

	clusterName = req.PathParameter("clusterName")
	if clusterName == "" {
		result.Err = "empty clusterName"
		return
	}

	page, err := strconv.ParseInt(req.QueryParameter("page"), 10, 64)
	if err != nil {
		log.Infof("invalid page: %s", req.QueryParameter("page"))
		page = 1
	}
	pageSize, err = strconv.ParseInt(req.QueryParameter("page_size"), 10, 64)
	if err != nil {
		log.Infof("invalid page_size: %s", req.QueryParameter("page_size"))
		pageSize = 10
	}
	ruleGroups, err := h.prometheusProcessor.ListGroups(clusterName)
	if err != nil {
		result.Err = err.Error()
		return
	}

	alarmPolicies := rest.AlarmPolicies{}
	for i := range ruleGroups {
		rg := ruleGroups[i]
		alarmPolicy := rest.NewAlarmPolicyFromRuleGroup(rg)
		alarmPolicies = append(alarmPolicies, alarmPolicy)
	}
	sort.Sort(alarmPolicies)

	count := int64(len(alarmPolicies))
	if (page-1)*pageSize >= count {
		alarmPolicies = rest.AlarmPolicies{}
	} else if page*pageSize >= count {
		alarmPolicies = alarmPolicies[(page-1)*pageSize:]
	} else {
		alarmPolicies = alarmPolicies[(page-1)*pageSize : page*pageSize]
	}

	p := rest.AlarmPolicyPagination{
		Page:          page,
		PageSize:      pageSize,
		Total:         count,
		AlarmPolicies: alarmPolicies,
	}

	result.Data = p
	result.Result = true
	status = http.StatusOK
}
