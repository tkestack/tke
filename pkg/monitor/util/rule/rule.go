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

	"tkestack.io/tke/pkg/monitor/util"

	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type rule struct {
	keyFunc KeyFunc
	data    map[string]*ruleGroupWrapper
}

type ruleGroupWrapper struct {
	v *v1.RuleGroup
	// ruleGroup revision
	rev     int
	ruleMap map[string]*ruleMetadata
}

type ruleMetadata struct {
	rev int
	idx int
}

var _ util.GenericRuleOperator = &rule{}

// KeyFunc defines a function to generate key from v1.Rule
type KeyFunc func(rule *v1.Rule) string

// NewGenericRuleOperator returns a instance to operate prometheus rules
func NewGenericRuleOperator(keyFunc KeyFunc) util.GenericRuleOperator {
	return &rule{
		keyFunc: keyFunc,
		data:    make(map[string]*ruleGroupWrapper),
	}
}

func (r *rule) InsertRule(groupName string, rule *v1.Rule) (int, *v1.Rule, error) {
	ruleGroup := r.getWrapperRuleGroup(groupName)
	if ruleGroup == nil {
		return 0, nil, errors.New("group not found")
	}

	key := r.keyFunc(rule)
	if _, ok := ruleGroup.ruleMap[key]; ok {
		return 0, nil, errors.New("record already existed")
	}

	if err := util.ValidateLabels(rule.Labels); err != nil {
		return 0, nil, err
	}

	lastIdx := len(ruleGroup.v.Rules)
	ruleGroup.v.Rules = append(ruleGroup.v.Rules, *rule)
	ruleGroup.rev++

	metadata := &ruleMetadata{
		rev: 1,
		idx: lastIdx,
	}
	ruleGroup.ruleMap[key] = metadata

	return metadata.rev, rule, nil
}

func (r *rule) DeleteRule(groupName, recordName string) (*v1.Rule, error) {
	ruleGroup := r.getWrapperRuleGroup(groupName)
	if ruleGroup == nil {
		return nil, errors.New("group not found")
	}

	metadata, ok := ruleGroup.ruleMap[recordName]
	if !ok {
		return nil, errors.New("record not found")
	}

	found := &ruleGroup.v.Rules[metadata.idx]

	rebuildRuleGroup := make([]v1.Rule, 0)
	rebuildRuleGroup = append(rebuildRuleGroup, ruleGroup.v.Rules[:metadata.idx]...)
	if metadata.idx < len(ruleGroup.ruleMap)-1 {
		rebuildRuleGroup = append(rebuildRuleGroup, ruleGroup.v.Rules[metadata.idx+1:]...)
	}
	ruleGroup.v.Rules = rebuildRuleGroup
	ruleGroup.rev++
	delete(ruleGroup.ruleMap, recordName)

	// Fix idx
	for _, v := range ruleGroup.ruleMap {
		if v.idx > metadata.idx {
			v.idx--
		}
	}

	return found, nil
}

func (r *rule) UpdateRule(groupName, recordName string, rev int, rule *v1.Rule) (int, *v1.Rule, error) {
	ruleGroup := r.getWrapperRuleGroup(groupName)
	if ruleGroup == nil {
		return 0, nil, errors.New("group not found")
	}

	metadata, ok := ruleGroup.ruleMap[recordName]
	if !ok {
		return 0, nil, errors.New("record not found")
	}

	if err := util.ValidateLabels(rule.Labels); err != nil {
		return 0, nil, err
	}

	// temporarily disable revision check
	// if metadata.rev != rev {
	//	return 0, nil, errors.Errorf("rev conflict %d", metadata.rev)
	// }

	ruleGroup.v.Rules[metadata.idx] = *rule
	metadata.rev++
	ruleGroup.rev++

	return metadata.rev, rule, nil
}

func (r *rule) GetRule(groupName, recordName string) (int, *v1.Rule, error) {
	ruleGroup := r.getWrapperRuleGroup(groupName)
	if ruleGroup == nil {
		return 0, nil, errors.New("group not found")
	}

	metadata, ok := ruleGroup.ruleMap[recordName]
	if !ok {
		return 0, nil, errors.New("record not found")
	}

	return metadata.rev, &ruleGroup.v.Rules[metadata.idx], nil
}

func (r *rule) ListRule(groupName string) ([]*v1.Rule, error) {
	ruleGroup := r.getWrapperRuleGroup(groupName)
	if ruleGroup == nil {
		return nil, errors.New("group not found")
	}

	rules := make([]*v1.Rule, len(ruleGroup.v.Rules))

	for i := range ruleGroup.v.Rules {
		rules[i] = &ruleGroup.v.Rules[i]
	}

	return rules, nil
}

func (r *rule) InsertRuleGroup(group *v1.RuleGroup) (int, *v1.RuleGroup, error) {
	if _, ok := r.data[group.Name]; ok {
		return 0, nil, errors.New("group already existed")
	}

	for _, rule := range group.Rules {
		if err := util.ValidateLabels(rule.Labels); err != nil {
			return 0, nil, err
		}
	}

	inserted := r.insertRuleGroup(group)
	return inserted.rev, group, nil
}

func (r *rule) DeleteRuleGroup(groupName string) (*v1.RuleGroup, error) {
	found, ok := r.data[groupName]
	if !ok {
		return nil, errors.New("group not found")
	}

	delete(r.data, groupName)

	return found.v, nil
}

func (r *rule) GetRuleGroup(groupName string) (int, *v1.RuleGroup, error) {
	ruleGroup := r.getWrapperRuleGroup(groupName)
	if ruleGroup == nil {
		return 0, nil, errors.New("group not found")
	}

	return ruleGroup.rev, ruleGroup.v, nil
}

func (r *rule) UpdateRuleGroup(groupName string, rev int, data *v1.RuleGroup) (int, *v1.RuleGroup, error) {
	ruleGroup := r.getWrapperRuleGroup(groupName)
	if ruleGroup == nil {
		return 0, nil, errors.New("group not found")
	}

	for _, rule := range data.Rules {
		if err := util.ValidateLabels(rule.Labels); err != nil {
			return 0, nil, err
		}
	}

	// temporarily disable revision check
	// if ruleGroup.rev != rev {
	//	return 0, nil, errors.Errorf("rev conflict %d", ruleGroup.rev)
	// }

	delete(r.data, groupName)
	inserted := r.insertRuleGroup(data)
	// inserted.rev = rev + 1

	return inserted.rev, data, nil
}

func (r *rule) ListGroup() ([]*v1.RuleGroup, error) {
	groups := make([]*v1.RuleGroup, len(r.data))

	i := 0
	for _, data := range r.data {
		groups[i] = data.v
		i++
	}

	return groups, nil
}

func (r *rule) LoadPromRule(groups []v1.RuleGroup) error {
	for i := range groups {
		r.insertRuleGroup(&groups[i])
	}
	return nil
}

func (r *rule) SavePromRule() []v1.RuleGroup {
	var groups []v1.RuleGroup
	for _, groupData := range r.data {
		group := *groupData.v
		groups = append(groups, group)
	}
	return groups
}

func (r *rule) Load(reader io.Reader) error {
	data := &v1.PrometheusRuleSpec{}
	err := yaml.NewDecoder(reader).Decode(data)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	for i := range data.Groups {
		r.insertRuleGroup(&data.Groups[i])
	}

	return nil
}

func (r *rule) Save(writer io.Writer) error {
	out := &v1.PrometheusRuleSpec{Groups: make([]v1.RuleGroup, 0)}

	for _, groupData := range r.data {
		out.Groups = append(out.Groups, *groupData.v)
	}

	return yaml.NewEncoder(writer).Encode(out)
}

func (r *rule) insertRuleGroup(group *v1.RuleGroup) *ruleGroupWrapper {
	groupRule := &ruleGroupWrapper{
		ruleMap: make(map[string]*ruleMetadata),
	}
	r.data[group.Name] = groupRule
	groupRule.v = group
	groupRule.rev = 1

	for i := range group.Rules {
		groupRule.ruleMap[r.keyFunc(&group.Rules[i])] = &ruleMetadata{
			rev: 1,
			idx: i,
		}
	}

	return groupRule
}

func (r *rule) getWrapperRuleGroup(groupName string) *ruleGroupWrapper {
	groupRule, ok := r.data[groupName]
	if ok {
		return groupRule
	}

	return nil
}
