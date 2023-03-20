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
	"helm.sh/helm/v3/pkg/action"
	"time"
)

// RollbackOptions is the options to a rollback call.
type RollbackOptions struct {
	Namespace   string
	ReleaseName string
	Revision    int64
	Timeout     time.Duration
	Wait        bool
	WaitForJobs bool
}

// Rollback roll back to the previous release.
func (c *Client) Rollback(options *RollbackOptions) error {
	actionConfig, err := c.buildActionConfig(options.Namespace)
	if err != nil {
		return err
	}
	client := action.NewRollback(actionConfig)
	client.Version = int(options.Revision)
	client.Timeout = options.Timeout
	client.Wait = options.Wait
	client.WaitForJobs = options.WaitForJobs
	return client.Run(options.ReleaseName)
}
