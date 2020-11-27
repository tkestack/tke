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
	"tkestack.io/tke/api/application"
	v1 "tkestack.io/tke/api/application/v1"
	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/application/config"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	chartpathv1 "tkestack.io/tke/pkg/application/util/chartpath/v1"
)

// FullfillChartInfo will full fill chart info by chartgroup
func FullfillChartInfo(appChart application.Chart, cg registryv1.ChartGroup) (application.Chart, error) {
	var v1Chart = &v1.Chart{}
	err := v1.Convert_application_Chart_To_v1_Chart(&appChart, v1Chart, nil)
	if err != nil {
		return appChart, err
	}

	v1Info, err := chartpathv1.FullfillChartInfo(*v1Chart, cg)
	if err != nil {
		return appChart, err
	}

	var chart = &application.Chart{}
	err = v1.Convert_v1_Chart_To_application_Chart(&v1Info, chart, nil)
	if err != nil {
		return appChart, err
	}

	return *chart, nil
}

// BuildChartPathBasicOptions will judge chartgroup type and return well-structured ChartPathOptions
func BuildChartPathBasicOptions(repo config.RepoConfiguration, appChart application.Chart) (opt helmaction.ChartPathOptions, err error) {
	var v1Chart = &v1.Chart{}
	err = v1.Convert_application_Chart_To_v1_Chart(&appChart, v1Chart, nil)
	if err != nil {
		return opt, err
	}

	return chartpathv1.BuildChartPathBasicOptions(repo, *v1Chart)
}
