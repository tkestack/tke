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
	"context"
	"reflect"
	"testing"

	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	prometheusrule "tkestack.io/tke/pkg/platform/controller/addon/prometheus"
)

func TestProcessor_DeleteGroup(t *testing.T) {
	mClient, p, clusterName, err := createProcessorServer()
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	expectRuleGroup := getExpectRule(exampleRuleStr)
	expectRuleGroup.Groups[0].Rules = make([]v1.Rule, 0)

	t.Logf("With non-existed group")
	err = p.DeleteGroup(context.Background(), clusterName, "non-exist-group")
	if err == nil {
		t.Errorf("deletion should failed")
		return
	}

	err = p.CreateGroup(context.Background(), clusterName, expectRuleGroup.Groups[0].Name, &expectRuleGroup.Groups[0])
	if err != nil {
		t.Errorf("creation should success, code: %s", err)
		return
	}

	t.Logf("With correct group name")
	err = p.DeleteGroup(context.Background(), clusterName, expectRuleGroup.Groups[0].Name)
	if err != nil {
		t.Errorf("delete should success, code: %s", err)
		return
	}

	t.Logf("Validate deletion")
	_, err = p.GetGroup(context.Background(), clusterName, expectRuleGroup.Groups[0].Name)
	if err == nil {
		t.Errorf("get should failed")
		return
	}

	t.Logf("Validate persistent data")
	prometheusRule, err := mClient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Get(context.Background(), prometheusrule.PrometheusRuleAlert, metav1.GetOptions{})
	if err != nil {
		t.Errorf("can't get persistent data, %v", err)
		return
	}
	groups := prometheusRule.Spec.Groups

	expectRuleGroup.Groups = nil
	if !reflect.DeepEqual(groups, expectRuleGroup.Groups) {
		t.Errorf("rule group is not equal, got %v, expect %v", groups, expectRuleGroup.Groups)
		return
	}
}

func TestProcessor_DeleteRule(t *testing.T) {
	mClient, p, clusterName, err := createProcessorServer()
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	expectRuleGroup := getExpectRule(exampleRuleStr)
	recordName := expectRuleGroup.Groups[0].Rules[0].Alert

	t.Logf("With non-existed record name")
	err = p.DeleteRule(context.Background(), clusterName, expectRuleGroup.Groups[0].Name, "non-exist-record")
	if err == nil {
		t.Errorf("delete should failed")
		return
	}

	err = p.CreateGroup(context.Background(), clusterName, expectRuleGroup.Groups[0].Name, &expectRuleGroup.Groups[0])
	if err != nil {
		t.Errorf("creation should success, code: %s", err)
		return
	}

	t.Logf("With existed record")
	err = p.DeleteRule(context.Background(), clusterName, expectRuleGroup.Groups[0].Name, recordName)
	if err != nil {
		t.Errorf("should not fail, code: %s", err)
		return
	}

	t.Logf("Validate deletion")
	_, err = p.GetRule(context.Background(), clusterName, expectRuleGroup.Groups[0].Name, recordName)
	if err == nil {
		t.Errorf("get should failed")
		return
	}

	t.Logf("Validate persistent data")
	prometheusRule, err := mClient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Get(context.Background(), prometheusrule.PrometheusRuleAlert, metav1.GetOptions{})
	if err != nil {
		t.Errorf("can't get persistent data, %v", err)
		return
	}
	groups := prometheusRule.Spec.Groups

	expectRuleGroup.Groups[0].Rules = make([]v1.Rule, 0)
	if !reflect.DeepEqual(groups, expectRuleGroup.Groups) {
		t.Errorf("rule group is not equal, got %v, expect %v", groups, expectRuleGroup.Groups)
		return
	}
}
