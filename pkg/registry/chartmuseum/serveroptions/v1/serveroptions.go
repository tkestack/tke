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
	"helm.sh/chartmuseum/pkg/chartmuseum/server/multitenant"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	registryconfigv1 "tkestack.io/tke/pkg/registry/apis/config/v1"
	"tkestack.io/tke/pkg/registry/chartmuseum/serveroptions"
	"tkestack.io/tke/pkg/util/log"
)

// BuildChartMuseumConfig build chartmuseum options from registry config
func BuildChartMuseumConfig(v1RegistryConfig *registryconfigv1.RegistryConfiguration, pathPrefix string, maxUploadSize int) (*multitenant.MultiTenantServerOptions, error) {
	registryConfig := &registryconfig.RegistryConfiguration{}
	if err := registryconfigv1.Convert_v1_RegistryConfiguration_To_config_RegistryConfiguration(v1RegistryConfig, registryConfig, nil); err != nil {
		log.Error("Failed to convert registry configuration", log.Err(err))
		return nil, err
	}
	return serveroptions.BuildChartMuseumConfig(registryConfig, pathPrefix, maxUploadSize)
}
