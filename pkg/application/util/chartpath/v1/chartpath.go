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

package v1

import (
	"fmt"
	"net/url"

	v1 "tkestack.io/tke/api/application/v1"
	registryv1 "tkestack.io/tke/api/registry/v1"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	registryconfig "tkestack.io/tke/pkg/registry/config"
	registryutil "tkestack.io/tke/pkg/registry/util"
)

// FullfillChartInfo will full fill chart info by chartgroup
func FullfillChartInfo(appChart v1.Chart, cg registryv1.ChartGroup) (v1.Chart, error) {
	if cg.Spec.Type == registryv1.RepoTypeImported {
		appChart.ImportedRepo = true
		appChart.RepoURL = cg.Spec.ImportedInfo.Addr
		appChart.RepoUsername = cg.Spec.ImportedInfo.Username
		appChart.RepoPassword = cg.Spec.ImportedInfo.Password
	}
	return appChart, nil
}

// BuildChartPathBasicOptions will judge chartgroup type and return well-structured ChartPathOptions
func BuildChartPathBasicOptions(repo registryconfig.RepoConfiguration, appChart v1.Chart) (opt helmaction.ChartPathOptions, err error) {
	if appChart.ImportedRepo {
		password, err := registryutil.VerifyDecodedPassword(appChart.RepoPassword)
		if err != nil {
			return opt, err
		}

		opt.RepoURL = appChart.RepoURL
		opt.Username = appChart.RepoUsername
		opt.Password = password
	} else {
		loc := &url.URL{
			Scheme: repo.Scheme,
			Host:   registryutil.BuildTenantRegistryDomain(repo.DomainSuffix, appChart.TenantID),
			Path:   fmt.Sprintf("/chart/%s", appChart.ChartGroupName),
		}
		opt.CaFile = repo.CaFile
		opt.RepoURL = loc.String()
		opt.Username = repo.Admin
		opt.Password = repo.AdminPassword
	}

	opt.ChartRepo = appChart.TenantID + "/" + appChart.ChartGroupName
	opt.Chart = appChart.ChartName
	opt.Version = appChart.ChartVersion
	return opt, nil
}
