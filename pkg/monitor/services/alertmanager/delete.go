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

package alertmanager

import (
	"bytes"
	"tkestack.io/tke/pkg/util/log"

	"github.com/pkg/errors"
)

func (h *processor) Delete(clusterName string, alertValue string) error {
	h.Lock()
	defer h.Unlock()

	if clusterName == "" {
		return errors.New("empty clusterName")
	}

	if alertValue == "" {
		return errors.New("empty alertValue")
	}

	routeOp, err := h.loadConfig(clusterName)
	if err != nil {
		return errors.Wrapf(err, "route operator not found")
	}

	log.Infof("Start to delete route %s", alertValue)
	_, err = routeOp.DeleteRoute(alertValue)
	if err != nil {
		return errors.Wrapf(err, "failed to delete route %s", alertValue)
	}

	output := bytes.NewBufferString("")
	err = routeOp.Save(output)
	if err != nil {
		return errors.Wrapf(err, "failed to save")
	}

	err = h.saveConfig(clusterName, output.String())
	if err != nil {
		return errors.Wrapf(err, "failed to save configmap")
	}

	return nil
}
