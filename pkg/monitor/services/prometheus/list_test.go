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
)

func TestProcessor_ListGroup(t *testing.T) {
	_, p, clusterName, err := createProcessorServer()
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	expectRuleGroup := getExpectRule(exampleRuleStr)
	err = p.CreateGroup(clusterName, expectRuleGroup.Groups[0].Name, &expectRuleGroup.Groups[0])
	if err != nil {
		t.Errorf("creation should success, code: %s", err)
		return
	}

	t.Logf("List all groups")
	targetGroups, err := p.ListGroups(clusterName)
	if err != nil {
		t.Errorf("list should success, code: %s", err)
		return
	}

	if len(targetGroups) != 1 {
		t.Errorf("group number should be 1, got %d", len(targetGroups))
		return
	}

	if !reflect.DeepEqual(targetGroups[0], &expectRuleGroup.Groups[0]) {
		t.Errorf("rule group not equal, got %v, expect %v", targetGroups[0], &expectRuleGroup.Groups[0])
		return
	}
}

func TestProcessor_ListRules(t *testing.T) {
	_, p, clusterName, err := createProcessorServer()
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	expectRuleGroup := getExpectRule(exampleRuleStr)
	err = p.CreateGroup(clusterName, expectRuleGroup.Groups[0].Name, &expectRuleGroup.Groups[0])
	if err != nil {
		t.Errorf("creation should success, code: %s", err)
		return
	}

	t.Logf("List non-existed group")
	_, err = p.ListRules(clusterName, "non-exist-group")
	if err == nil {
		t.Errorf("list should failed")
		return
	}

	t.Logf("List correct group")
	targetRules, err := p.ListRules(clusterName, expectRuleGroup.Groups[0].Name)
	if err != nil {
		t.Errorf("list should success, code: %s", err)
		return
	}

	if len(targetRules) != 1 {
		t.Errorf("rules number should be 1, got %d", len(targetRules))
		return
	}

	if !reflect.DeepEqual(targetRules[0], &expectRuleGroup.Groups[0].Rules[0]) {
		t.Errorf("rule group not equal, got %v, expect %v", targetRules[0], &expectRuleGroup.Groups[0].Rules[0])
		return
	}
}
