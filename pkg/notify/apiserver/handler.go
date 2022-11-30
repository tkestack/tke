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

package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/server/mux"
	restclient "k8s.io/client-go/rest"
	notifyinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/notify/internalversion"
	"tkestack.io/tke/api/notify"
	"tkestack.io/tke/pkg/util/log"
)

const (
	alertNameKey          = "alertName"
	startsAtKey           = "startsAt"
	alarmPolicyTypeKey    = "alarmPolicyType"
	alarmPolicyNameKey    = "alarmPolicyName"
	clusterIDKey          = "clusterID"
	clusterDisplayNameKey = "clusterDisplayName"
	valueKey              = "value"
	workloadKindKey       = "workloadKind"
	namespaceKey          = "namespace"
	workloadNameKey       = "workloadName"
	virtualMachineKey     = "virtualMachine"
	podNameKey            = "podName"
	nodeNameKey           = "nodeName"
	nodeRoleKey           = "nodeRole"
	unitKey               = "unit"
	evaluateTypeKey       = "evaluateType"
	evaluateValueKey      = "evaluateValue"
	metricDisplayNameKey  = "metricDisplayName"
	summaryKey            = "summary"
)

// Request response struct
type responseMsg struct {
	StatusCode int    `json:"status"`
	Msg        string `json:"message"`
}

// Alert indicates the alert infos
type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    time.Time         `json:"startsAt"`
	EndsAt      time.Time         `json:"endsAt"`
	Status      string            `json:"status"`
}

// Notification indicates the notification for alertmanager of prometheus
type Notification struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
}

func registerAlarmWebhook(m *mux.PathRecorderMux, loopbackClientConfig *restclient.Config) {
	m.HandleFunc("/webhook", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Check method
		if req.Method != "POST" {
			setErrResponse("Should use POST method", http.StatusMethodNotAllowed, w)
			return
		}
		if req.Body == nil {
			setErrResponse("Body is nil", http.StatusBadRequest, w)
			return
		}
		var bodyBytes []byte
		notifyInfo := &Notification{}
		bodyBytes, _ = ioutil.ReadAll(req.Body)
		err := json.Unmarshal(bodyBytes, notifyInfo)
		if err != nil {
			setErrResponse(fmt.Sprintf("Invalid body: %v", err.Error()), http.StatusBadRequest, w)
			return
		}
		if len(notifyInfo.Alerts) == 0 {
			setErrResponse("Alerts is nil", http.StatusBadRequest, w)
			return
		}
		log.Infof("Receive alerts: %+v", notifyInfo.Alerts)
		for _, alert := range notifyInfo.Alerts {
			annotations := alert.Annotations
			notifyWay, ok := annotations["notifyWay"]
			if !ok {
				setErrResponse("The notifyWay does not exist", http.StatusBadRequest, w)
				return
			}
			ways := strings.Split(notifyWay, ",")
			if len(ways) == 0 {
				setErrResponse("notifyWay is nil", http.StatusBadRequest, w)
				return
			}
			variables := getVariables(alert)
			for _, way := range ways {
				channelAndTemplate := strings.Split(way, ":")
				if len(channelAndTemplate) != 2 {
					setErrResponse("Invalid notifyWay", http.StatusBadRequest, w)
					return
				}
				channel := channelAndTemplate[0]
				template := channelAndTemplate[1]

				var receivers []string
				receiversList, ok := annotations["receivers"]
				if ok && receiversList != "" {
					receivers = strings.Split(receiversList, ",")
				}

				var receiverGroups []string
				receiverGroupsList, ok := annotations["receiverGroups"]
				if ok && receiverGroupsList != "" {
					receiverGroups = strings.Split(receiverGroupsList, ",")
				}

				if len(receivers) == 0 && len(receiverGroups) == 0 {
					setErrResponse("receivers and receiverGroups are nil", http.StatusBadRequest, w)
					return
				}
				messageRequest := newMessageRequest(channel, template, receivers, receiverGroups, variables, alert.Status)

				notifyClient := notifyinternalclient.NewForConfigOrDie(loopbackClientConfig)
				_, err = notifyClient.MessageRequests(messageRequest.ObjectMeta.Namespace).Create(req.Context(), messageRequest, metav1.CreateOptions{})
				if err != nil {
					setErrResponse(err.Error(), http.StatusInternalServerError, w)
					return
				}
				log.Infof("messageRequest created: %+v", messageRequest.Spec)
			}
		}
		response := &responseMsg{
			StatusCode: http.StatusOK,
			Msg:        "Success",
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			setErrResponse(err.Error(), http.StatusInternalServerError, w)
			return
		}
	})
}

func setErrResponse(msg string, statusCode int, w http.ResponseWriter) {
	response := &responseMsg{
		StatusCode: statusCode,
		Msg:        msg,
	}
	jsonMsg, _ := json.Marshal(response)
	http.Error(w, string(jsonMsg), statusCode)
}

func newMessageRequest(channel string, template string, receivers []string, receiverGroups []string, variables map[string]string, status string) *notify.MessageRequest {
	return &notify.MessageRequest{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: channel,
		},
		Spec: notify.MessageRequestSpec{
			TemplateName:   template,
			Receivers:      receivers,
			ReceiverGroups: receiverGroups,
			Variables:      variables,
		},

		Status: notify.MessageRequestStatus{
			AlertStatus: status,
		},
	}
}

func getVariables(alert Alert) map[string]string {
	variables := make(map[string]string)
	labels := alert.Labels
	annotations := alert.Annotations

	summary := fmt.Sprintf("发生时间：%s", processStartTime(alert.StartsAt))

	alarmPolicyTypeValue, ok := annotations["alarmPolicyType"]
	if ok {
		summary = fmt.Sprintf("%s\n告警策略类型：%s", summary, alarmPolicyTypeValue)
	}

	metricDisplayNameValue, ok := annotations[metricDisplayNameKey]
	if ok {
		summary = fmt.Sprintf("%s\n告警项：%s", summary, metricDisplayNameValue)
	}

	valueValue, ok := annotations["value"]
	if ok {
		summary = fmt.Sprintf("%s\n指标值：%s", summary, valueValue)
	}

	unitValue, ok := annotations[unitKey]
	if ok {
		summary = fmt.Sprintf("%s%s", summary, unitValue)
	}

	evaluateTypeValue, ok := annotations[evaluateTypeKey]
	if ok {
		summary = fmt.Sprintf("%s %s", summary, evaluateTypeValue)
	}

	evaluateValue, ok := annotations[evaluateValueKey]
	if ok {
		summary = fmt.Sprintf("%s %s%s", summary, evaluateValue, unitValue)
	}

	alarmPolicyNameValue, ok := labels["alarmPolicyName"]
	if ok {
		summary = fmt.Sprintf("%s\n告警策略名：%s", summary, alarmPolicyNameValue)
	}

	alertNameValue, ok := labels["alertname"]
	if ok {
		summary = fmt.Sprintf("%s\n指标名：%s", summary, alertNameValue)
	}

	clusterIDValue, ok := labels["cluster_id"]
	if ok {
		summary = fmt.Sprintf("%s\n集群ID：%s", summary, clusterIDValue)
	}

	clusterDisplayNameValue, ok := labels["cluster_display_name"]
	if ok {
		summary = fmt.Sprintf("%s\n集群名称：%s", summary, clusterDisplayNameValue)
	}

	workloadKindValue, ok := labels["workload_kind"]
	if ok {
		summary = fmt.Sprintf("%s\n工作负载类型：%s", summary, workloadKindValue)
	}

	workloadNameValue, ok := labels["workload_name"]
	if ok {
		summary = fmt.Sprintf("%s\n工作负载名称：%s", summary, workloadNameValue)
	}

	virtualMachineName, ok := labels["name"]
	if ok {
		summary = fmt.Sprintf("%s\n虚拟机名称：%s", summary, virtualMachineName)
	}

	namespaceValue, ok := labels["namespace"]
	if ok {
		summary = fmt.Sprintf("%s\n命名空间：%s", summary, namespaceValue)
	}

	podNameValue, ok := labels["pod_name"]
	if ok {
		summary = fmt.Sprintf("%s\nPOD名称：%s", summary, podNameValue)
	}

	nodeNameValue, ok := labels["node"]
	if ok {
		summary = fmt.Sprintf("%s\n节点名称：%s", summary, nodeNameValue)
	}

	nodeRoleValue, ok := labels["node_role"]
	if ok {
		summary = fmt.Sprintf("%s\n节点类型：%s", summary, nodeRoleValue)
	}

	variables[startsAtKey] = processStartTime(alert.StartsAt)
	variables[alarmPolicyTypeKey] = alarmPolicyTypeValue
	variables[alarmPolicyNameKey] = alarmPolicyNameValue
	variables[valueKey] = valueValue
	variables[alertNameKey] = alertNameValue
	variables[clusterIDKey] = clusterIDValue
	variables[clusterDisplayNameKey] = clusterDisplayNameValue
	variables[workloadKindKey] = workloadKindValue
	variables[workloadNameKey] = workloadNameValue
	variables[namespaceKey] = namespaceValue
	variables[podNameKey] = podNameValue
	variables[nodeNameKey] = nodeNameValue
	variables[nodeRoleKey] = nodeRoleValue
	variables[unitKey] = unitValue
	variables[evaluateTypeKey] = evaluateTypeValue
	variables[evaluateValueKey] = evaluateValue
	variables[metricDisplayNameKey] = metricDisplayNameValue
	variables[summaryKey] = summary
	variables[virtualMachineKey] = virtualMachineName

	return variables
}

func processStartTime(t time.Time) string {
	localTime := t.Local()
	return localTime.Format(time.RFC3339)
}
