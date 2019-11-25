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

package provider

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	machineprovider "tkestack.io/tke/pkg/platform/provider/machine"
)

// LoadMachineProvider load the specify type machine provider from map
func LoadMachineProvider(machineProviders *sync.Map, pluginType string) (machineprovider.Provider, error) {
	pluginName := fmt.Sprintf("%s-machine", strings.ToLower(pluginType))
	machineProviderClient, ok := machineProviders.Load(pluginName)
	if !ok {
		return nil, errors.Errorf("can't get %q provider", pluginName)
	}

	p, ok := machineProviderClient.(machineprovider.Provider)
	if !ok {
		return nil, errors.New("provider type assertion error")
	}

	return p, nil
}

// LoadClusterProvider load the specify type cluster provider from map
func LoadClusterProvider(clusterProviders *sync.Map, pluginType string) (clusterprovider.Provider, error) {
	pluginName := fmt.Sprintf("%s-cluster", strings.ToLower(pluginType))
	client, ok := clusterProviders.Load(pluginName)
	if !ok {
		return nil, errors.Errorf("can't get %q provider", pluginName)
	}

	p, ok := client.(clusterprovider.Provider)
	if !ok {
		return nil, errors.New("provider type assertion error")
	}

	return p, nil
}
