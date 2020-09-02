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

package util

import (
	"helm.sh/helm/v3/pkg/release"
)

// Filter returning a release of matches.
func Filter(rels []*release.Release, namespace, name string) (*release.Release, bool) {
	if len(rels) == 0 {
		return nil, false
	}
	for _, v := range rels {
		if v.Namespace == namespace && v.Name == name {
			return v, true
		}
	}
	return nil, false
}
