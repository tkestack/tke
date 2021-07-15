/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package v1

func (ca ClusterApps) Len() int {
	return len(ca)
}

func (ca ClusterApps) Swap(i, j int) {
	ca[i], ca[j] = ca[j], ca[i]
}

func (ca ClusterApps) Less(i, j int) bool {
	return ca[i].Priority < ca[j].Priority
}

func (ca ClusterApps) HasApp(appNamespace, appName string) bool {
	for _, app := range ca {
		if app.AppNamespace == appNamespace && app.App.Spec.Name == appName {
			return true
		}
	}
	return false
}
