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
	"strings"
	"time"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// UninstallOptions is the option for uninstalling releases.
type UninstallOptions struct {
	DryRun       bool
	KeepHistory  bool
	DisableHooks bool
	Timeout      time.Duration
	Description  string

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

	// 需要执行hook
	if !options.DisableHooks {
		// 先禁用删除的hook，通过调用函数的方式执行hook
		client.DisableHooks = true
		client.KeepHistory = true
		// 手动调用
		rels, err := actionConfig.Releases.History(options.ReleaseName)
		if err != nil {
			return nil, errors.Wrapf(err, "uninstall: Release not loaded: %s", options.ReleaseName)
		}
		if len(rels) < 1 {
			return nil, errors.New("no release provided")
		}

		releaseutil.SortByRevision(rels)
		rel := rels[len(rels)-1]
		err = actionConfig.ExecHook(rel, release.HookPreDelete, options.Timeout)
		if err != nil {
			return nil, err
		}
		// release记录已经是删除状态了，就不用再调用uninstall了，因为设置了keepHistory，会保留release记录
		if rel.Info.Status != release.StatusUninstalled {
			_, err = client.Run(options.ReleaseName)
			if err != nil {
				if !strings.Contains(err.Error(), "release: not found") || !k8serrors.IsNotFound(err) {
					return nil, err
				}
			}
		}
		err = actionConfig.ExecHook(rel, release.HookPostDelete, options.Timeout)
		if err != nil {
			return nil, err
		}

		// 如果不保留release，保证hook成功之后再删除release
		if !options.KeepHistory {
			for _, rel := range rels {
				_, err = actionConfig.Releases.Delete(rel.Name, rel.Version)
				if err != nil {
					return nil, errors.Wrapf(err, "uninstall: Release not loaded: %s", options.ReleaseName)
				}
			}
		}
		return nil, err
	} else {
		return client.Run(options.ReleaseName)
	}
}
