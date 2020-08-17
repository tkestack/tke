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

package action

import (
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

// UninstallOptions is the option for uninstalling releases.
type UninstallOptions struct {
	DryRun      bool
	KeepHistory bool
	Timeout     time.Duration
	Description string

	ReleaseName string
	Namespace   string
}

// Uninstall provides the implementation of 'helm uninstall'.
func (c *Client) Uninstall(options *UninstallOptions) (*release.UninstallReleaseResponse, error) {
	actionConfig, err := c.buildActionConfig(options.Namespace)
	if err != nil {
		return nil, err
	}
	client := action.NewUninstall(actionConfig)
	client.DryRun = options.DryRun
	client.KeepHistory = options.KeepHistory
	client.Timeout = options.Timeout
	client.Description = options.Description
	return client.Run(options.ReleaseName)
}
