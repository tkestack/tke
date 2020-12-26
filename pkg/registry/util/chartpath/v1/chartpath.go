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

	v1 "tkestack.io/tke/api/registry/v1"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	registryconfig "tkestack.io/tke/pkg/registry/config"
	"tkestack.io/tke/pkg/registry/util"
)

// BuildChartPathBasicOptions will judge chartgroup type and return well-structured ChartPathOptions
func BuildChartPathBasicOptions(repo registryconfig.RepoConfiguration, cg v1.ChartGroup) (opt helmaction.ChartPathOptions, err error) {
	if cg.Spec.Type == v1.RepoTypeImported {
		password, err := util.VerifyDecodedPassword(cg.Spec.ImportedInfo.Password)
		if err != nil {
			return opt, err
		}

		opt.RepoURL = cg.Spec.ImportedInfo.Addr
		opt.Username = cg.Spec.ImportedInfo.Username
		opt.Password = password
	} else {
		loc := &url.URL{
			Scheme: repo.Scheme,
			Host:   util.BuildTenantRegistryDomain(repo.DomainSuffix, cg.Spec.TenantID),
			Path:   fmt.Sprintf("/chart/%s", cg.Spec.Name),
		}
		opt.CaFile = repo.CaFile
		opt.RepoURL = loc.String()
		opt.Username = repo.Admin
		opt.Password = repo.AdminPassword
	}
	opt.ChartRepo = cg.Spec.TenantID + "/" + cg.Spec.Name
	opt.Chart = ""
	opt.Version = ""
	return opt, nil
}
