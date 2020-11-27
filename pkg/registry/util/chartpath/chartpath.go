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

package chartpath

import (
	"tkestack.io/tke/api/registry"
	v1 "tkestack.io/tke/api/registry/v1"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	"tkestack.io/tke/pkg/registry/config"
	chartpathv1 "tkestack.io/tke/pkg/registry/util/chartpath/v1"
)

// BuildChartPathBasicOptions will judge chartgroup type and return well-structured ChartPathOptions
func BuildChartPathBasicOptions(repo config.RepoConfiguration, cg registry.ChartGroup) (opt helmaction.ChartPathOptions, err error) {
	var v1ChartGroup = &v1.ChartGroup{}
	err = v1.Convert_registry_ChartGroup_To_v1_ChartGroup(&cg, v1ChartGroup, nil)
	if err != nil {
		return opt, err
	}

	return chartpathv1.BuildChartPathBasicOptions(repo, *v1ChartGroup)
}
