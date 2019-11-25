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
	"github.com/pkg/errors"
	"tkestack.io/tke/pkg/util/log"
)

func (h *processor) DeleteGroup(clusterName, groupName string) error {
	h.Lock()
	defer h.Unlock()

	ruleOp, err := h.loadRule(clusterName)
	if err != nil {
		return errors.Wrapf(err, "rule operator not found")
	}

	log.Infof("Start to delete ruleGroup %s", groupName)

	_, err = ruleOp.DeleteRuleGroup(groupName)
	if err != nil {
		return errors.Wrapf(err, "failed to delete")
	}

	groups := ruleOp.SavePromRule()

	err = h.saveRule(clusterName, groups)
	if err != nil {
		return errors.Wrapf(err, "failed to save configmap")
	}

	return nil
}

func (h *processor) DeleteRule(clusterName, groupName, recordName string) error {
	h.Lock()
	defer h.Unlock()

	ruleOp, err := h.loadRule(clusterName)
	if err != nil {
		return errors.Wrapf(err, "rule operator not found")
	}

	log.Infof("Start to delete rule into %s(%s)", groupName, recordName)

	_, err = ruleOp.DeleteRule(groupName, recordName)
	if err != nil {
		return errors.Wrapf(err, "failed to delete group %s(%s)", groupName, recordName)
	}

	groups := ruleOp.SavePromRule()

	err = h.saveRule(clusterName, groups)
	if err != nil {
		return errors.Wrapf(err, "failed to save configmap")
	}

	return nil
}
