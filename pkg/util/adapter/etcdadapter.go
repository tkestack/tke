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

package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	client "github.com/coreos/etcd/clientv3"
)

const (
	// requestTimeout the timeout for failing to operate etcd object.
	requestTimeout = 5 * time.Second

	// placeHolder represent the NULL value in the Casbin Rule.
	placeHolder = "_"

	// defaultKey is the root path in ETCD, if not provided.
	defaultKey = "casbin_policy"
)

// casbinRule represents the struct stored into etcd backend.
type casbinRule struct {
	Key   string `json:"key"`
	PType string `json:"ptype"`
	V0    string `json:"v0"`
	V1    string `json:"v1"`
	V2    string `json:"v2"`
	V3    string `json:"v3"`
	V4    string `json:"v4"`
	V5    string `json:"v5"`
	V6    string `json:"v6"`
}

// Adapter is the policy storage adapter for Casbin. With this library, Casbin can load policy
// from ETCD and save policy to it. ETCD adapter support the Auto-Save feature for Casbin policy.
// This means it can support adding a single policy rule to the storage, or removing a single policy
// rule from the storage. See: https://github.com/sebastianliu/etcd-adapter.
type Adapter struct {
	key string
	// etcd connection client
	conn *client.Client
}

// NewAdapter creates a new adaptor instance.
func NewAdapter(client *client.Client, key string) *Adapter {
	return newAdapter(client, key)
}

func newAdapter(client *client.Client, key string) *Adapter {
	if key == "" {
		key = defaultKey
	}
	a := &Adapter{
		conn: client,
		key:  key,
	}

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)

	return a
}

// finalizer is the destructor for Adapter.
func finalizer(a *Adapter) {
	_ = a.conn.Close()
}

// LoadPolicy loads all of policys from ETCD
func (a *Adapter) LoadPolicy(model model.Model) error {
	var rule casbinRule
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	getResp, err := a.conn.Get(ctx, a.getRootKey(), client.WithPrefix())
	if err != nil {
		return err
	}
	if len(getResp.Kvs) == 0 {
		// there is no policy
		return errors.New("there is no policy in ETCD for the moment")
	}
	for _, kv := range getResp.Kvs {
		err = json.Unmarshal(kv.Value, &rule)
		if err != nil {
			return err
		}
		a.loadPolicy(rule, model)
	}
	return nil
}

func (a *Adapter) getRootKey() string {
	return fmt.Sprintf("/%s", a.key)
}

func (a *Adapter) loadPolicy(rule casbinRule, model model.Model) {
	lineText := rule.PType
	if rule.V0 != "" {
		lineText += ", " + rule.V0
	}
	if rule.V1 != "" {
		lineText += ", " + rule.V1
	}
	if rule.V2 != "" {
		lineText += ", " + rule.V2
	}
	if rule.V3 != "" {
		lineText += ", " + rule.V3
	}
	if rule.V4 != "" {
		lineText += ", " + rule.V4
	}
	if rule.V5 != "" {
		lineText += ", " + rule.V5
	}
	if rule.V6 != "" {
		lineText += ", " + rule.V6
	}

	persist.LoadPolicyLine(lineText, model)
}

// SavePolicy will rewrite all of policies in ETCD with the current data in Casbin
func (a *Adapter) SavePolicy(model model.Model) error {
	// clean old rule data
	_ = a.destroy()

	var rules []casbinRule

	for ptype, ast := range model["p"] {
		for _, line := range ast.Policy {
			rules = append(rules, a.convertRule(ptype, line))
		}
	}

	for ptype, ast := range model["g"] {
		for _, line := range ast.Policy {
			rules = append(rules, a.convertRule(ptype, line))
		}
	}

	return a.savePolicy(rules)
}

// destroy or clean all of policy
func (a *Adapter) destroy() error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := a.conn.Delete(ctx, a.getRootKey(), client.WithPrefix())
	return err
}

func (a *Adapter) convertRule(ptype string, line []string) (rule casbinRule) {
	rule = casbinRule{}
	rule.PType = ptype
	policys := []string{ptype}
	length := len(line)

	if len(line) > 0 {
		rule.V0 = line[0]
		policys = append(policys, line[0])
	}
	if len(line) > 1 {
		rule.V1 = line[1]
		policys = append(policys, line[1])
	}
	if len(line) > 2 {
		rule.V2 = line[2]
		policys = append(policys, line[2])
	}
	if len(line) > 3 {
		rule.V3 = line[3]
		policys = append(policys, line[3])
	}
	if len(line) > 4 {
		rule.V4 = line[4]
		policys = append(policys, line[4])
	}
	if len(line) > 5 {
		rule.V5 = line[5]
		policys = append(policys, line[5])
	}

	if len(line) > 6 {
		rule.V6 = line[6]
		policys = append(policys, line[6])
	}

	for i := 0; i < 7-length; i++ {
		policys = append(policys, placeHolder)
	}

	rule.Key = strings.Join(policys, "::")

	return rule
}

func (a *Adapter) savePolicy(rules []casbinRule) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	for _, rule := range rules {
		ruleData, _ := json.Marshal(rule)
		_, err := a.conn.Put(ctx, a.constructPath(rule.Key), string(ruleData))
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Adapter) constructPath(key string) string {
	return fmt.Sprintf("/%s/%s", a.key, key)
}

// AddPolicy adds a policy rule to the storage.
// Part of the Auto-Save feature.
func (a *Adapter) AddPolicy(sec string, ptype string, line []string) error {
	rule := a.convertRule(ptype, line)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	ruleData, _ := json.Marshal(rule)
	_, err := a.conn.Put(ctx, a.constructPath(rule.Key), string(ruleData))
	return err
}

// RemovePolicy removes a policy rule from the storage.
// Part of the Auto-Save feature.
func (a *Adapter) RemovePolicy(sec string, ptype string, line []string) error {
	rule := a.convertRule(ptype, line)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := a.conn.Delete(ctx, a.constructPath(rule.Key))
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// Part of the Auto-Save feature.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	rule := casbinRule{}

	rule.PType = ptype
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		rule.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		rule.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		rule.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		rule.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		rule.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		rule.V5 = fieldValues[5-fieldIndex]
	}
	if fieldIndex <= 6 && 6 < fieldIndex+len(fieldValues) {
		rule.V6 = fieldValues[6-fieldIndex]
	}

	filter := a.constructFilter(rule)

	return a.removeFilteredPolicy(filter)
}

func (a *Adapter) constructFilter(rule casbinRule) string {
	var filter string
	if rule.PType != "" {
		filter = fmt.Sprintf("/%s/%s", a.key, rule.PType)
	} else {
		filter = fmt.Sprintf("/%s/.*", a.key)
	}

	if rule.V0 != "" {
		filter = fmt.Sprintf("%s::%s", filter, normalize(rule.V0))
	} else {
		filter = fmt.Sprintf("%s::.*", filter)
	}

	if rule.V1 != "" {
		filter = fmt.Sprintf("%s::%s", filter, normalize(rule.V1))
	} else {
		filter = fmt.Sprintf("%s::.*", filter)
	}

	if rule.V2 != "" {
		filter = fmt.Sprintf("%s::%s", filter, normalize(rule.V2))
	} else {
		filter = fmt.Sprintf("%s::.*", filter)
	}

	if rule.V3 != "" {
		filter = fmt.Sprintf("%s::%s", filter, normalize(rule.V3))
	} else {
		filter = fmt.Sprintf("%s::.*", filter)
	}

	if rule.V4 != "" {
		filter = fmt.Sprintf("%s::%s", filter, normalize(rule.V4))
	} else {
		filter = fmt.Sprintf("%s::.*", filter)
	}

	if rule.V5 != "" {
		filter = fmt.Sprintf("%s::%s", filter, normalize(rule.V5))
	} else {
		filter = fmt.Sprintf("%s::.*", filter)
	}

	if rule.V6 != "" {
		filter = fmt.Sprintf("%s::%s", filter, normalize(rule.V6))
	} else {
		filter = fmt.Sprintf("%s::.*", filter)
	}

	return filter
}

func (a *Adapter) removeFilteredPolicy(filter string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	// get all policy key
	getResp, err := a.conn.Get(ctx, a.constructPath(""), client.WithPrefix(), client.WithKeysOnly())
	if err != nil {
		return err
	}
	var filteredKeys []string
	for _, kv := range getResp.Kvs {
		matched, err := regexp.MatchString(filter, string(kv.Key))
		if err != nil {
			return err
		}
		if matched {
			filteredKeys = append(filteredKeys, string(kv.Key))
		}
	}
	for _, key := range filteredKeys {
		_, err := a.conn.Delete(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func normalize(str string) string {
	return strings.Replace(str, "*", "\\*", -1)
}
