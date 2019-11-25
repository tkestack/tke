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

package prometheus

import (
	"fmt"
	"strings"
	"time"

	"github.com/coreos/prometheus-operator/pkg/apis/monitoring"

	"tkestack.io/tke/pkg/monitor/services"
	"tkestack.io/tke/pkg/monitor/util"
	prometheus_rule "tkestack.io/tke/pkg/platform/controller/addon/prometheus"
	"tkestack.io/tke/pkg/util/log"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/coreos/prometheus-operator/pkg/client/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	exampleRuleStr = `groups:
- name: example
  interval: 1m
  rules:
  - alert: job:http_inprogress_requests:sum
    expr: sum(http_inprogress_requests) by (job)
    for: 1m
    labels:
      alert: test
    annotations:
      value: 1
`
	testClusterName = "fake"
)

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

func createProcessorServer() (*fake.Clientset, services.RuleProcessor, string, error) {
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

	return mClient, p, testClusterName, nil
}

func getExpectRule(data string) *monitoringv1.PrometheusRuleSpec {
	ruleSpec := &monitoringv1.PrometheusRuleSpec{}
	reader := strings.NewReader(data)
	if err := yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(ruleSpec); err != nil {
		fmt.Printf("unmarshal err %s", err.Error())
		return nil
	}
	fmt.Printf("ruleSpec: %v", ruleSpec)
	return ruleSpec
}
