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
	"reflect"
	"testing"

	prometheus_rule "tkestack.io/tke/pkg/platform/controller/addon/prometheus"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestProcessor_UpdateGroup(t *testing.T) {
	mClient, p, clusterName, err := createProcessorServer()
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	expectRuleGroup := getExpectRule(exampleRuleStr)

	t.Logf("With non-existed group")
	err = p.UpdateGroup(clusterName, "non-exist-group", &expectRuleGroup.Groups[0])
	if err == nil {
		t.Errorf("update should failed")
		return
	}

	err = p.CreateGroup(clusterName, expectRuleGroup.Groups[0].Name, &expectRuleGroup.Groups[0])
	if err != nil {
		t.Errorf("creation should success, %s", err)
		return
	}

	t.Logf("With correct group name")
	expectRuleGroup.Groups[0].Rules = nil
	err = p.UpdateGroup(clusterName, expectRuleGroup.Groups[0].Name, &expectRuleGroup.Groups[0])
	if err != nil {
		t.Errorf("update should success, code: %s", err)
		return
	}

	t.Logf("Get group")
	targetGroup, err := p.GetGroup(clusterName, expectRuleGroup.Groups[0].Name)
	if err != nil {
		t.Errorf("get should success, code: %s", err)
		return
	}

	if !reflect.DeepEqual(targetGroup, &expectRuleGroup.Groups[0]) {
		t.Errorf("rule group not equal, got %v, expect %v", targetGroup, &expectRuleGroup.Groups[0])
		return
	}

	t.Logf("Validate persistent data")
	prometheusRule, err := mClient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Get(prometheus_rule.PrometheusRuleAlert, metav1.GetOptions{})
	if err != nil {
		t.Errorf("can't get persistent data, %v", err)
		return
	}
	groups := prometheusRule.Spec.Groups

	if !reflect.DeepEqual(&groups[0], &expectRuleGroup.Groups[0]) {
		t.Errorf("rule group not equal, got %v, expect %v", &groups[0], &expectRuleGroup.Groups[0])
		return
	}
}

func TestProcessor_UpdateRule(t *testing.T) {
	mClient, p, clusterName, err := createProcessorServer()
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	expectRuleGroup := getExpectRule(exampleRuleStr)
	recordName := expectRuleGroup.Groups[0].Rules[0].Alert

	t.Logf("With non-existed record")
	err = p.UpdateRule(clusterName, expectRuleGroup.Groups[0].Name, "non-exist-record", &expectRuleGroup.Groups[0].Rules[0])
	if err == nil {
		t.Errorf("update should failed")
		return
	}

	t.Logf("Create record")
	err = p.CreateGroup(clusterName, expectRuleGroup.Groups[0].Name, &expectRuleGroup.Groups[0])
	if err != nil {
		t.Errorf("creation should success, %s", err)
		return
	}

	t.Logf("With correct record name")
	expectRuleGroup.Groups[0].Rules[0].Labels = map[string]string{
		"alert": "test",
		"foo":   "bar",
	}
	err = p.UpdateRule(clusterName, expectRuleGroup.Groups[0].Name, recordName, &expectRuleGroup.Groups[0].Rules[0])
	if err != nil {
		t.Errorf("update should success, code: %s", err)
		return
	}

	targetRule, err := p.GetRule(clusterName, expectRuleGroup.Groups[0].Name, recordName)
	if err != nil {
		t.Errorf("get should success, code: %s", err)
		return
	}

	if !reflect.DeepEqual(targetRule, &expectRuleGroup.Groups[0].Rules[0]) {
		t.Errorf("rule group not equal, got %v, expect %v", targetRule, &expectRuleGroup.Groups[0].Rules[0])
		return
	}

	t.Logf("Validate persistent data")
	prometheusRule, err := mClient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Get(prometheus_rule.PrometheusRuleAlert, metav1.GetOptions{})
	if err != nil {
		t.Errorf("can't get persistent data, %v", err)
		return
	}
	groups := prometheusRule.Spec.Groups

	if !reflect.DeepEqual(&groups[0].Rules[0], &expectRuleGroup.Groups[0].Rules[0]) {
		t.Errorf("rule group not equal, got %v, expect %v", groups[0].Rules[0], &expectRuleGroup.Groups[0].Rules[0])
		return
	}
}
