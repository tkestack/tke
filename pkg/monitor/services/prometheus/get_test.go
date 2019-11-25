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

	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
)

func TestProcessor_GetGroup(t *testing.T) {
	_, p, clusterName, err := createProcessorServer()
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	expectRuleGroup := getExpectRule(exampleRuleStr)
	expectRuleGroup.Groups[0].Rules = make([]v1.Rule, 0)

	t.Logf("With non-existed group")
	_, err = p.GetGroup(clusterName, "non-exist-group")
	if err == nil {
		t.Errorf("get should failed")
		return
	}

	err = p.CreateGroup(clusterName, expectRuleGroup.Groups[0].Name, &expectRuleGroup.Groups[0])
	if err != nil {
		t.Errorf("creation should success, code: %s", err)
		return
	}

	t.Logf("With correct group name")
	targetGroup, err := p.GetGroup(clusterName, expectRuleGroup.Groups[0].Name)
	if err != nil {
		t.Errorf("get should success, code: %s", err)
		return
	}

	if !reflect.DeepEqual(targetGroup, &expectRuleGroup.Groups[0]) {
		t.Errorf("rule group not equal, got %v, expect %v", targetGroup, &expectRuleGroup.Groups[0])
		return
	}
}

func TestProcessor_GetRule(t *testing.T) {
	_, p, clusterName, err := createProcessorServer()
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	expectRuleGroup := getExpectRule(exampleRuleStr)
	recordName := expectRuleGroup.Groups[0].Rules[0].Alert

	t.Logf("With non-existed record")
	_, err = p.GetRule(clusterName, expectRuleGroup.Groups[0].Name, "non-exist-group")
	if err == nil {
		t.Errorf("get should failed")
		return
	}

	err = p.CreateGroup(clusterName, expectRuleGroup.Groups[0].Name, &expectRuleGroup.Groups[0])
	if err != nil {
		t.Errorf("creation should success, code: %s", err)
		return
	}

	t.Logf("With existed record")
	targetRule, err := p.GetRule(clusterName, expectRuleGroup.Groups[0].Name, recordName)
	if err != nil {
		t.Errorf("get should success, code: %s", err)
		return
	}

	if !reflect.DeepEqual(targetRule, &expectRuleGroup.Groups[0].Rules[0]) {
		t.Errorf("rule group not equal, got %v, expect %v", targetRule, &expectRuleGroup.Groups[0].Rules[0])
		return
	}
}
