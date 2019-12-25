/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package util

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	"tkestack.io/tke/pkg/util/log"

	"k8s.io/apimachinery/pkg/fields"
	authv1 "tkestack.io/tke/api/auth/v1"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	authv1client "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	authv1lister "tkestack.io/tke/api/client/listers/auth/v1"
)

// RestAdapter is the policy storage adapter for Casbin. With this library, Casbin can load policy
// from kubernetes rest storage and save policy to it. Rest adapter support the Auto-Save feature for Casbin policy.
// This means it can support adding a single policy rule to the storage, or removing a single policy
// rule from the storage.
type RestAdapter struct {
	ruleClient authv1client.RuleInterface
	lister     authv1lister.RuleLister
}

// NewAdapter creates a new adaptor instance.
func NewAdapter(ruleClient authv1client.RuleInterface, ruleLister authv1lister.RuleLister) *RestAdapter {
	adapter := &RestAdapter{
		ruleClient: ruleClient,
		lister:     ruleLister,
	}

	return adapter
}

// LoadPolicy loads all of policys from backend
func (a *RestAdapter) LoadPolicy(model model.Model) error {
	rules, err := a.lister.List(labels.Everything())
	if err != nil {
		return fmt.Errorf("list all rules failed: %v", err)
	}

	// Uncomment this line to see when the policy is loaded.
	// log.Info("List rules", log.Int("rules", len(rules)))

	for _, rule := range rules {
		a.loadPolicy(rule, model)
	}

	return nil
}

func (a *RestAdapter) loadPolicy(rule *authv1.Rule, model model.Model) {
	casRule := rule.Spec
	lineText := casRule.PType
	if casRule.V0 != "" {
		lineText += ", " + casRule.V0
	}
	if casRule.V1 != "" {
		lineText += ", " + casRule.V1
	}
	if casRule.V2 != "" {
		lineText += ", " + casRule.V2
	}
	if casRule.V3 != "" {
		lineText += ", " + casRule.V3
	}
	if casRule.V4 != "" {
		lineText += ", " + casRule.V4
	}
	if casRule.V5 != "" {
		lineText += ", " + casRule.V5
	}
	if casRule.V6 != "" {
		lineText += ", " + casRule.V6
	}

	persist.LoadPolicyLine(lineText, model)
}

// SavePolicy will rewrite all of policies in ETCD with the current data in Casbin
func (a *RestAdapter) SavePolicy(model model.Model) error {
	// clean old rule data
	err := a.destroy()
	if err != nil {
		return err
	}

	var rules []authv1.Rule

	for ptype, ast := range model["p"] {
		for _, line := range ast.Policy {
			rules = append(rules, ConvertRule(ptype, line))
		}
	}

	for ptype, ast := range model["g"] {
		for _, line := range ast.Policy {
			rules = append(rules, ConvertRule(ptype, line))
		}
	}

	return a.savePolicy(rules)
}

// destroy or clean all of policy
func (a *RestAdapter) destroy() error {
	err := a.ruleClient.DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{})
	return err
}

func ConvertRule(ptype string, line []string) (rule authv1.Rule) {
	rule = authv1.Rule{}
	rule.Spec.PType = ptype

	if len(line) > 0 {
		rule.Spec.V0 = line[0]
	}
	if len(line) > 1 {
		rule.Spec.V1 = line[1]
	}
	if len(line) > 2 {
		rule.Spec.V2 = line[2]
	}
	if len(line) > 3 {
		rule.Spec.V3 = line[3]
	}
	if len(line) > 4 {
		rule.Spec.V4 = line[4]
	}
	if len(line) > 5 {
		rule.Spec.V5 = line[5]
	}
	if len(line) > 6 {
		rule.Spec.V6 = line[6]
	}

	return rule
}

func (a *RestAdapter) savePolicy(rules []authv1.Rule) error {
	for _, rule := range rules {
		if _, err := a.ruleClient.Create(&rule); err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

// AddPolicy adds a policy rule to the storage.
// Part of the Auto-Save feature.
func (a *RestAdapter) AddPolicy(sec string, ptype string, line []string) error {
	rule := ConvertRule(ptype, line)
	if _, err := a.ruleClient.Create(&rule); !apierrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// RemovePolicy removes a policy rule from the storage.
// Part of the Auto-Save feature.
func (a *RestAdapter) RemovePolicy(sec string, ptype string, line []string) error {
	rule := ConvertRule(ptype, line)
	filter := a.constructRemoveSelector(rule)

	return a.removeFilteredPolicy(filter)
}

func (a *RestAdapter) constructRemoveSelector(rule authv1.Rule) string {
	ruleFieldSet := fields.Set{}
	ruleFieldSet["spec.ptype"] = rule.Spec.PType
	ruleFieldSet["spec.v0"] = rule.Spec.V0
	ruleFieldSet["spec.v1"] = rule.Spec.V1
	ruleFieldSet["spec.v2"] = rule.Spec.V2
	ruleFieldSet["spec.v3"] = rule.Spec.V3
	ruleFieldSet["spec.v4"] = rule.Spec.V4
	ruleFieldSet["spec.v5"] = rule.Spec.V5
	ruleFieldSet["spec.v6"] = rule.Spec.V6
	return fields.SelectorFromSet(ruleFieldSet).String()
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// Part of the Auto-Save feature.
func (a *RestAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	rule := authv1.Rule{}

	rule.Spec.PType = ptype
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		rule.Spec.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		rule.Spec.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		rule.Spec.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		rule.Spec.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		rule.Spec.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		rule.Spec.V5 = fieldValues[5-fieldIndex]
	}
	if fieldIndex <= 6 && 6 < fieldIndex+len(fieldValues) {
		rule.Spec.V6 = fieldValues[6-fieldIndex]
	}

	filter := a.constructFilterSelector(rule)

	return a.removeFilteredPolicy(filter)
}

func (a *RestAdapter) constructFilterSelector(rule authv1.Rule) string {

	ruleFieldSet := fields.Set{}
	if rule.Spec.PType != "" {
		ruleFieldSet["spec.ptype"] = rule.Spec.PType
	}

	if rule.Spec.V0 != "" {
		ruleFieldSet["spec.v0"] = rule.Spec.V0
	}

	if rule.Spec.V1 != "" {
		ruleFieldSet["spec.v1"] = rule.Spec.V1
	}

	if rule.Spec.V2 != "" {
		ruleFieldSet["spec.v2"] = rule.Spec.V2
	}

	if rule.Spec.V3 != "" {
		ruleFieldSet["spec.v3"] = rule.Spec.V3
	}

	if rule.Spec.V4 != "" {
		ruleFieldSet["spec.v4"] = rule.Spec.V4
	}

	if rule.Spec.V5 != "" {
		ruleFieldSet["spec.v5"] = rule.Spec.V5
	}

	if rule.Spec.V6 != "" {
		ruleFieldSet["spec.v6"] = rule.Spec.V6
	}

	return fields.SelectorFromSet(ruleFieldSet).String()
}

func (a *RestAdapter) removeFilteredPolicy(filter string) error {
	log.Info("RemoveFilterPolicy", log.String("filter", filter))
	return a.ruleClient.DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{FieldSelector: filter})
}
