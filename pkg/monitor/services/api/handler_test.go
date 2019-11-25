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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"tkestack.io/tke/pkg/monitor/services/rest"
	"tkestack.io/tke/pkg/monitor/util"
	prometheus_rule "tkestack.io/tke/pkg/platform/controller/addon/prometheus"

	"github.com/coreos/prometheus-operator/pkg/apis/monitoring"
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/coreos/prometheus-operator/pkg/client/versioned/fake"
	"github.com/emicklei/go-restful"
	"github.com/parnurzeal/gorequest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/pkg/util/log"
)

const (
	exampleAlarmPolicy = `{
    "AlarmPolicySettings":{
        "AlarmPolicyName":"test2",
        "AlarmPolicyDescription":"11111",
        "AlarmPolicyType":"pod",
        "StatisticsPeriod":60,
        "AlarmMetrics":[
            {
                "Measurement":"k8s_pod",
                "MetricName":"k8s_pod_rate_cpu_core_used_node",
                "Evaluator":{
                    "Type":"gt",
                    "Value":"80"
                },
                "ContinuePeriod":5,
                "Unit":"%"
            }
        ],
        "AlarmObjects":"elsanli-test-harbor-chartmuseum,elsanli-test-harbor-clair,elsanli-test-harbor-core,elsanli-test-harbor-jobservice",
        "AlarmObjectsType":"part"
    },
    "NotifySettings":{
        "ReceiverGroups":[
            "75061"
        ],
        "Receivers":[
            "75061"
        ],
        "NotifyWay":[
        {
            "ChannelName":"chan1",
            "TemplateName":"temp1"
        },
        {
            "ChannelName":"chan2",
            "TemplateName":"temp2"
        }
        ]
    },
    "Namespace":"default",
    "WorkloadType":"Deployment"
}
`
	testClusterName = "fake"
)

func TestProcessor_CreateGet(t *testing.T) {
	ch := make(chan struct{})
	_, addr, err := createProcessorServer(ch)
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	defer func() {
		close(ch)
	}()

	expectAlarmPolicy := getExpectAlarmPolicy(exampleAlarmPolicy)
	expectAlarmPolicyName := expectAlarmPolicy.AlarmPolicySettings.AlarmPolicyName
	url := fmt.Sprintf("http://%s/api/v1/monitor/%s/%s/%s", addr, clustersPrefix, testClusterName, alarmPolicyPrefix)

	t.Logf("With invalid request")
	client := gorequest.New().Post(url)
	setClient(client)
	resp, body, _ := client.SendString(`{}`).End()
	if resp.StatusCode == http.StatusOK {
		t.Errorf("creation should failed")
		return
	}

	t.Logf("With correct policy name")
	client = gorequest.New().Post(url)
	setClient(client)
	resp, body, _ = client.SendStruct(&expectAlarmPolicy).End()
	r := &rest.ResponseForTest{}
	_ = r.Decode(strings.NewReader(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("creation should success, code: %s, %s", resp.Status, r.Err)
		return
	}

	t.Logf("Validate insert policy information")
	client = gorequest.New().Get(url + "/" + expectAlarmPolicyName)
	setClient(client)
	resp, body, _ = client.End()

	r = &rest.ResponseForTest{}
	_ = r.Decode(strings.NewReader(body))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("get should success, code: %s, %s", resp.Status, r.Err)
		return
	}

	targetAlarmPolicy := &rest.AlarmPolicy{}
	err = json.Unmarshal([]byte(r.Data), targetAlarmPolicy)
	if err != nil {
		t.Errorf("can't decode result, %v", err)
		return
	}

	if !reflect.DeepEqual(targetAlarmPolicy, expectAlarmPolicy) {
		t.Errorf("alarm policy not equal, got %v, expect %v", targetAlarmPolicy, expectAlarmPolicy)
		return
	}
}

func TestProcessor_Update(t *testing.T) {
	ch := make(chan struct{})
	_, addr, err := createProcessorServer(ch)
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	defer func() {
		close(ch)
	}()

	expectAlarmPolicy := getExpectAlarmPolicy(exampleAlarmPolicy)
	expectAlarmPolicyName := expectAlarmPolicy.AlarmPolicySettings.AlarmPolicyName
	url := fmt.Sprintf("http://%s/api/v1/monitor/%s/%s/%s", addr, clustersPrefix, testClusterName, alarmPolicyPrefix)

	t.Logf("create policy")
	client := gorequest.New().Post(url)
	setClient(client)
	resp, body, _ := client.SendStruct(&expectAlarmPolicy).End()
	r := &rest.ResponseForTest{}
	_ = r.Decode(strings.NewReader(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("creation should success, code: %s, %s", resp.Status, r.Err)
		return
	}

	t.Logf("update policy")
	expectAlarmPolicy.AlarmPolicySettings.StatisticsPeriod = 120
	client = gorequest.New().Put(url + "/" + expectAlarmPolicyName)
	setClient(client)
	resp, body, _ = client.SendStruct(&expectAlarmPolicy).End()
	r = &rest.ResponseForTest{}
	_ = r.Decode(strings.NewReader(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("creation should success, code: %s, %s", resp.Status, r.Err)
		return
	}

	t.Logf("Validate insert policy information")
	client = gorequest.New().Get(url + "/" + expectAlarmPolicyName)
	setClient(client)
	resp, body, _ = client.End()

	r = &rest.ResponseForTest{}
	_ = r.Decode(strings.NewReader(body))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("get should success, code: %s, %s", resp.Status, r.Err)
		return
	}

	targetAlarmPolicy := &rest.AlarmPolicy{}
	err = json.Unmarshal([]byte(r.Data), targetAlarmPolicy)
	if err != nil {
		t.Errorf("can't decode result, %v", err)
		return
	}

	if !reflect.DeepEqual(targetAlarmPolicy, expectAlarmPolicy) {
		t.Errorf("alarm policy not equal, got %v, expect %v", targetAlarmPolicy, expectAlarmPolicy)
		return
	}
}

func TestProcessor_Delete(t *testing.T) {
	ch := make(chan struct{})
	_, addr, err := createProcessorServer(ch)
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	defer func() {
		close(ch)
	}()

	expectAlarmPolicy := getExpectAlarmPolicy(exampleAlarmPolicy)
	expectAlarmPolicyName := expectAlarmPolicy.AlarmPolicySettings.AlarmPolicyName
	url := fmt.Sprintf("http://%s/api/v1/monitor/%s/%s/%s", addr, clustersPrefix, testClusterName, alarmPolicyPrefix)

	t.Logf("create policy")
	client := gorequest.New().Post(url)
	setClient(client)
	resp, body, _ := client.SendStruct(&expectAlarmPolicy).End()
	r := &rest.ResponseForTest{}
	_ = r.Decode(strings.NewReader(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("creation should success, code: %s, %s", resp.Status, r.Err)
		return
	}

	t.Logf("delete policy")
	expectAlarmPolicy.AlarmPolicySettings.StatisticsPeriod = 120
	client = gorequest.New().Delete(url + "/" + expectAlarmPolicyName)
	setClient(client)
	resp, body, _ = client.End()
	r = &rest.ResponseForTest{}
	_ = r.Decode(strings.NewReader(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("creation should success, code: %s, %s", resp.Status, r.Err)
		return
	}

	t.Logf("Validate policy deletion")
	client = gorequest.New().Get(url + "/" + expectAlarmPolicyName)
	setClient(client)
	resp, body, _ = client.End()

	r = &rest.ResponseForTest{}
	_ = r.Decode(strings.NewReader(body))
	if resp.StatusCode == http.StatusOK {
		t.Errorf("get should fail, code: %s, %s", resp.Status, r.Err)
		return
	}
}

func TestProcessor_List(t *testing.T) {
	ch := make(chan struct{})
	_, addr, err := createProcessorServer(ch)
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	defer func() {
		close(ch)
	}()

	expectAlarmPolicy := getExpectAlarmPolicy(exampleAlarmPolicy)
	url := fmt.Sprintf("http://%s/api/v1/monitor/%s/%s/%s", addr, clustersPrefix, testClusterName, alarmPolicyPrefix)

	t.Logf("create policy")
	client := gorequest.New().Post(url)
	setClient(client)
	resp, body, _ := client.SendStruct(&expectAlarmPolicy).End()
	r := &rest.ResponseForTest{}
	_ = r.Decode(strings.NewReader(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("creation should success, code: %s, %s", resp.Status, r.Err)
		return
	}

	t.Logf("list policy")
	client = gorequest.New().Get(url).Param("page", "1").Param("page_size", "10")
	setClient(client)
	resp, body, _ = client.End()
	r = &rest.ResponseForTest{}
	_ = r.Decode(strings.NewReader(body))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("creation should success, code: %s, %s", resp.Status, r.Err)
		return
	}

	t.Logf("Validate list policy")
	targetAlarmPolicies := &rest.AlarmPolicyPagination{}
	err = json.Unmarshal([]byte(r.Data), targetAlarmPolicies)
	if err != nil {
		t.Errorf("can't decode result, %v", err)
		return
	}

	expectAlarmPolicies := &rest.AlarmPolicyPagination{
		Page:          1,
		PageSize:      10,
		Total:         1,
		AlarmPolicies: rest.AlarmPolicies{expectAlarmPolicy},
	}
	if !reflect.DeepEqual(targetAlarmPolicies, expectAlarmPolicies) {
		t.Errorf("alarm policy not equal, got %v, expect %v", targetAlarmPolicies, expectAlarmPolicies)
		return
	}
}

func init() {
	logOpts := log.NewOptions()
	logOpts.EnableCaller = true
	logOpts.Level = log.InfoLevel
	log.Init(logOpts)
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			log.Flush()
		}
	}()
}

func createProcessorServer(stopCh chan struct{}) (*fake.Clientset, string, error) {
	mClient := fake.NewSimpleClientset()
	prometheusRule := &monitoringv1.PrometheusRule{
		TypeMeta: metav1.TypeMeta{
			APIVersion: monitoring.GroupName + "/v1",
			Kind:       monitoringv1.PrometheusRuleKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheus_rule.PrometheusRuleAlert,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{prometheus_rule.PrometheusService: prometheus_rule.PrometheusCRDName, "role": "alert-rules"},
		},
		Spec: monitoringv1.PrometheusRuleSpec{Groups: []monitoringv1.RuleGroup{}},
	}
	_, err := mClient.MonitoringV1().Prometheuses(metav1.NamespaceSystem).List(metav1.ListOptions{})
	if err != nil {
		fmt.Printf("mclient err %s", err.Error())
	}
	util.ClusterNameToMonitor.Store(testClusterName, mClient)
	_, _ = mClient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Create(prometheusRule)
	// Because we have set kubernetes client, so set nil is ok
	p := NewProcessor(nil)

	ws := new(restful.WebService)
	ws.Path("/api/v1/monitor")
	p.RegisterWebService(ws)
	container := restful.NewContainer()
	container.Add(ws)

	srv := httptest.NewServer(container)
	go func() {
		<-stopCh
		srv.Close()
	}()

	return mClient, srv.Listener.Addr().String(), nil
}

func setClient(client *gorequest.SuperAgent) {
	client.Set(restful.HEADER_Accept, restful.MIME_JSON).
		Set(restful.HEADER_ContentType, restful.MIME_JSON)
}

func getExpectAlarmPolicy(data string) *rest.AlarmPolicy {
	alarmPolicy := &rest.AlarmPolicy{}
	err := json.Unmarshal([]byte(data), alarmPolicy)
	if err != nil {
		fmt.Printf("json unmarshal error: %v", err)
	}
	fmt.Printf("Got: %+v", alarmPolicy.AlarmPolicySettings)
	return alarmPolicy
}
