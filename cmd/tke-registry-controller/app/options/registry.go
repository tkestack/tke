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

package options

import (
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	registryscheme "tkestack.io/tke/pkg/registry/apis/config/scheme"
	registryconfigv1 "tkestack.io/tke/pkg/registry/apis/config/v1"
)

// NewRegistryConfiguration will create a new RegistryConfiguration with default values
func NewRegistryConfiguration() (*registryconfig.RegistryConfiguration, error) {
	scheme, _, err := registryscheme.NewSchemeAndCodecs()
	if err != nil {
		return nil, err
	}
	versioned := &registryconfigv1.RegistryConfiguration{}
	scheme.Default(versioned)
	config := &registryconfig.RegistryConfiguration{}
	if err := scheme.Convert(versioned, config, nil); err != nil {
		return nil, err
	}
	return config, nil
}
