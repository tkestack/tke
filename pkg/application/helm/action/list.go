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
	"helm.sh/helm/v3/pkg/release"
)

// ListOptions is the query options to a list call.
type ListOptions struct {
	// All ignores the limit/offset
	All bool
	// AllNamespaces searches across namespaces
	AllNamespaces bool
	// Overrides the default lexicographic sorting
	ByDate      bool
	SortReverse bool
	// Limit is the number of items to return per Run()
	Limit int
	// Offset is the starting index for the Run() call
	Offset int
	// Filter is a filter that is applied to the results
	Filter       string
	Short        bool
	Uninstalled  bool
	Superseded   bool
	Uninstalling bool
	Deployed     bool
	Failed       bool
	Pending      bool

	Namespace string
}

// List returning a set of matches.
func (c *Client) List(options *ListOptions) ([]*release.Release, error) {
	var actionConfig *action.Configuration
	var err error
	if options.AllNamespaces {
		actionConfig, err = c.buildActionConfig("")
	} else {
		actionConfig, err = c.buildActionConfig(options.Namespace)
	}
	if err != nil {
		return []*release.Release{}, err
	}
	client := action.NewList(actionConfig)
	client.All = options.All
	client.AllNamespaces = options.AllNamespaces
	client.ByDate = options.ByDate
	client.SortReverse = options.SortReverse
	client.Limit = options.Limit
	client.Offset = options.Offset
	client.Filter = options.Filter
	client.Short = options.Short
	client.Uninstalled = options.Uninstalled
	client.Uninstalling = options.Uninstalling
	client.Superseded = options.Superseded
	client.Deployed = options.Deployed
	client.Failed = options.Failed
	client.Pending = options.Pending
	client.SetStateMask()
	return client.Run()
}
