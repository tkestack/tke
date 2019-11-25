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
	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/pkg/errors"
	"tkestack.io/tke/pkg/util/log"
)

func (h *processor) UpdateGroup(clusterName, groupName string, entity *v1.RuleGroup) error {
	h.Lock()
	defer h.Unlock()

	ruleOp, err := h.loadRule(clusterName)
	if err != nil {
		return errors.Wrapf(err, "rule operator not found")
	}

	log.Infof("Start to update ruleGroup %s", entity.Name)

	_, _, err = ruleOp.UpdateRuleGroup(groupName, 0, entity)
	if err != nil {
		return errors.Wrapf(err, "failed to update")
	}

	groups := ruleOp.SavePromRule()

	err = h.saveRule(clusterName, groups)
	if err != nil {
		return errors.Wrapf(err, "failed to save configmap")
	}

	return nil
}

func (h *processor) UpdateRule(clusterName, groupName, recordName string, entity *v1.Rule) error {
	h.Lock()
	defer h.Unlock()

	ruleOp, err := h.loadRule(clusterName)
	if err != nil {
		return errors.Wrapf(err, "rule operator not found")
	}

	log.Infof("Start to update rule into %s", groupName)

	_, _, err = ruleOp.UpdateRule(groupName, recordName, 0, entity)
	if err != nil {
		return errors.Wrapf(err, "failed to update %s(%s)", groupName, recordName)
	}

	groups := ruleOp.SavePromRule()

	err = h.saveRule(clusterName, groups)
	if err != nil {
		return errors.Wrapf(err, "failed to save configmap")
	}

	return nil
}
