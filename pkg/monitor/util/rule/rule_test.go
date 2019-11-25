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

package rule

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v2"

	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"tkestack.io/tke/pkg/util/log"
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
)

func init() {
	logOpts := log.NewOptions()
	logOpts.EnableCaller = true
	logOpts.Level = log.ErrorLevel
	log.Init(logOpts)
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			log.Flush()
		}
	}()
}

func getExpectRule(data string) *v1.PrometheusRuleSpec {
	ruleSpec := &v1.PrometheusRuleSpec{}
	reader := strings.NewReader(data)
	if err := k8syaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(ruleSpec); err != nil {
		return nil
	}
	return ruleSpec
}

func TestNewGenericRuleOperator(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "rule")
	if err != nil {
		t.Errorf("create tmp file failed, %v", err)
		return
	}

	t.Logf("rule file: %s", file.Name())
	defer func() {
		_ = file.Close()
		if !t.Failed() {
			_ = os.Remove(file.Name())
		}
	}()

	expectRules := getExpectRule(exampleRuleStr)
	t.Logf("expectRules:%v", expectRules)
	expectStr, _ := yaml.Marshal(expectRules)
	_, err = file.Write(expectStr)
	if err != nil {
		t.Errorf("can't write tmp file, %v", err)
		return
	}

	_ = file.Sync()
	_, _ = file.Seek(0, io.SeekStart)

	operator := NewGenericRuleOperator(func(rule *v1.Rule) string {
		return rule.Record
	})

	t.Logf("Test load from disk")
	err = operator.Load(file)
	if err != nil {
		t.Errorf("can't load rules %v", err)
		return
	}

	t.Logf("Test get rule group")
	_, _, err = operator.GetRuleGroup("abc")
	if err == nil {
		t.Errorf("should not get rule")
		return
	}

	rev, targetGroup, err := operator.GetRuleGroup(expectRules.Groups[0].Name)
	if err != nil {
		t.Errorf("can't get rule %v", err)
		return
	}

	if rev != 1 {
		t.Errorf("rev should be 1, got %d", rev)
		return
	}

	if !reflect.DeepEqual(targetGroup, &expectRules.Groups[0]) {
		t.Errorf("rule not equal, got %v, expect %v", targetGroup, &expectRules.Groups[0])
		return
	}

	t.Logf("Test list groups")
	groups, err := operator.ListGroup()
	if err != nil {
		t.Errorf("can't list group, %v", err)
		return
	}

	if len(groups) != 1 {
		t.Errorf("group number should be 1, got %d", len(groups))
		return
	}

	if !reflect.DeepEqual(groups[0], &expectRules.Groups[0]) {
		t.Errorf("rule not equal, got %v, expect %v", groups[0], &expectRules.Groups[0])
		return
	}

	t.Logf("Test get rule")
	_, _, err = operator.GetRule(expectRules.Groups[0].Name, "test")
	if err == nil {
		t.Errorf("should not get rule")
		return
	}

	rev, targetRule, err := operator.GetRule(expectRules.Groups[0].Name, expectRules.Groups[0].Rules[0].Record)
	if err != nil {
		t.Errorf("can't get rule %v", err)
		return
	}

	if rev != 1 {
		t.Errorf("rev should be 1, got %d", rev)
		return
	}

	if !reflect.DeepEqual(targetRule, &expectRules.Groups[0].Rules[0]) {
		t.Errorf("rule not equal, got %v, expect %v", targetRule, &expectRules.Groups[0].Rules[0])
		return
	}

	t.Logf("Test list rule")
	_, err = operator.ListRule("abc")
	if err == nil {
		t.Errorf("shouldn't list abc group")
		return
	}

	rules, err := operator.ListRule(expectRules.Groups[0].Name)
	if err != nil {
		t.Errorf("can't list group rules, got %v", err)
		return
	}

	if len(rules) != 1 {
		t.Errorf("list rules should be 1, got %d", len(rules))
		return
	}

	if !reflect.DeepEqual(rules[0], &expectRules.Groups[0].Rules[0]) {
		t.Errorf("rule not equal, got %v, expect %v", rules[0], &expectRules.Groups[0].Rules[0])
		return
	}

	t.Logf("Test insert rule group")
	newRuleGroup := &v1.RuleGroup{
		Interval: time.Second.String(),
		Rules:    expectRules.Groups[0].Rules,
	}

	newRuleGroup.Name = expectRules.Groups[0].Name
	_, _, err = operator.InsertRuleGroup(newRuleGroup)
	if err == nil {
		t.Errorf("should not insert same rule")
		return
	}

	newRuleGroup.Name = "insertGroup"
	_, _, err = operator.InsertRuleGroup(newRuleGroup)
	if err != nil {
		t.Errorf("can't insert rule group %v", err)
		return
	}

	rev, targetGroup, err = operator.GetRuleGroup(newRuleGroup.Name)
	if err != nil {
		t.Errorf("can't get rule %v", err)
		return
	}

	if rev != 1 {
		t.Errorf("rev should be 1, got %v", rev)
		return
	}

	if !reflect.DeepEqual(targetGroup, newRuleGroup) {
		t.Errorf("rule not equal, got %v, expect %v", targetGroup, newRuleGroup)
		return
	}

	t.Logf("Test delete rule group")
	_, err = operator.DeleteRuleGroup("foo")
	if err == nil {
		t.Errorf("should not delete")
		return
	}

	_, err = operator.DeleteRuleGroup(newRuleGroup.Name)
	if err != nil {
		t.Errorf("can't delte rule, %v", err)
		return
	}

	_, _, err = operator.GetRuleGroup(newRuleGroup.Name)
	if err == nil {
		t.Errorf("should not get")
		return
	}

	t.Logf("Test insert rule")
	_, _, err = operator.InsertRule(expectRules.Groups[0].Name, &expectRules.Groups[0].Rules[0])
	if err == nil {
		t.Errorf("should not insert same rule")
		return
	}

	newRule := &v1.Rule{
		Record: "insertRule",
		Expr:   intstr.FromString("sum(abc) by (job)"),
		Labels: map[string]string{
			"alert": "abc",
		},
	}

	_, _, err = operator.InsertRule(expectRules.Groups[0].Name, newRule)
	if err != nil {
		t.Errorf("can't insert rule %v", err)
		return
	}

	ruleRev, targetRule, err := operator.GetRule(expectRules.Groups[0].Name, newRule.Record)
	if err != nil {
		t.Errorf("can't get rule %v", err)
		return
	}

	if ruleRev != 1 {
		t.Errorf("rule rev should be 2, got %d", rev)
		return
	}

	groupRev, _, err := operator.GetRuleGroup(expectRules.Groups[0].Name)
	if err != nil {
		t.Errorf("can't get rule group %v", err)
		return
	}

	if groupRev != 2 {
		t.Errorf("rev should be 2, got %d", rev)
		return
	}

	if !reflect.DeepEqual(targetRule, newRule) {
		t.Errorf("rule not equal, got %v, expect %v", targetRule, newRule)
		return
	}

	t.Logf("Test delete rule")
	_, err = operator.DeleteRule(expectRules.Groups[0].Name, "foo")
	if err == nil {
		t.Errorf("should not delete")
		return
	}

	_, err = operator.DeleteRule(expectRules.Groups[0].Name, newRule.Record)
	if err != nil {
		t.Errorf("can't delte rule, %v", err)
		return
	}

	_, _, err = operator.GetRule(expectRules.Groups[0].Name, newRule.Record)
	if err == nil {
		t.Errorf("should not get")
		return
	}

	t.Logf("Test save to disk")
	_, _ = file.Seek(0, io.SeekStart)
	_ = file.Truncate(0)
	err = operator.Save(file)
	if err != nil {
		t.Errorf("can't save rules %v", err)
		return
	}

	_, _ = file.Seek(0, io.SeekStart)
	targetBuffer := make([]byte, len(expectStr))
	_, err = file.Read(targetBuffer)
	if err != nil {
		t.Errorf("read tmp file %v", err)
		return
	}

	if string(targetBuffer) != string(expectStr) {
		t.Errorf("rule not equal, got %v, expect %v", string(targetBuffer), string(expectStr))
		return
	}

	t.Logf("Test update rule group")
	expectRules.Groups[0].Interval = time.Second.String()

	rev, _, err = operator.GetRuleGroup(expectRules.Groups[0].Name)
	if err != nil {
		t.Errorf("can't get rule %v", err)
		return
	}

	_, _, err = operator.UpdateRuleGroup("abc", rev, &expectRules.Groups[0])
	if err == nil {
		t.Errorf("should not update")
		return
	}

	_, _, err = operator.UpdateRuleGroup(expectRules.Groups[0].Name, rev, &expectRules.Groups[0])
	if err != nil {
		t.Errorf("can't update rule group, %v", err)
		return
	}

	_, targetGroup, err = operator.GetRuleGroup(expectRules.Groups[0].Name)
	if err != nil {
		t.Errorf("can't get rule %v", err)
		return
	}

	if !reflect.DeepEqual(targetGroup, &expectRules.Groups[0]) {
		t.Errorf("rule not equal, got %v, expect %v", targetGroup, &expectRules.Groups[0])
		return
	}

	t.Logf("Test update rule")
	expectRules.Groups[0].Rules[0].Labels = map[string]string{
		"alert": "test",
		"foo":   "bar",
	}

	rev, _, err = operator.GetRule(expectRules.Groups[0].Name, expectRules.Groups[0].Rules[0].Record)
	if err != nil {
		t.Errorf("can't get rule %v", err)
		return
	}

	_, _, err = operator.UpdateRule("abc", expectRules.Groups[0].Rules[0].Record, rev, &expectRules.Groups[0].Rules[0])
	if err == nil {
		t.Errorf("should not update")
		return
	}

	_, _, err = operator.UpdateRule(expectRules.Groups[0].Name, expectRules.Groups[0].Rules[0].Record, rev, &expectRules.Groups[0].Rules[0])
	if err != nil {
		t.Errorf("can't update rule, %v", err)
		return
	}

	_, _, err = operator.GetRule(expectRules.Groups[0].Name, expectRules.Groups[0].Rules[0].Record)
	if err != nil {
		t.Errorf("can't get rule %v", err)
		return
	}

	if !reflect.DeepEqual(targetGroup, &expectRules.Groups[0]) {
		t.Errorf("rule not equal, got %v, expect %v", targetGroup, &expectRules.Groups[0])
		return
	}
}
