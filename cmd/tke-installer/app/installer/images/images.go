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

package images

import (
	"reflect"
	"sort"

	"tkestack.io/tke/pkg/app/version"

	"tkestack.io/tke/pkg/util/containerregistry"
)

type Components struct {
	Busybox    containerregistry.Image
	Keepalived containerregistry.Image
	InfluxDB   containerregistry.Image

	ProviderRes              containerregistry.Image
	BaremetalClusterProvider containerregistry.Image
	BaremetalMachineProvider containerregistry.Image

	TKEGateway            containerregistry.Image
	TKEAuth               containerregistry.Image
	TKEBusinessAPI        containerregistry.Image
	TKEBusinessController containerregistry.Image
	TKEMonitorAPI         containerregistry.Image
	TKEMonitorController  containerregistry.Image
	TKENotifyAPI          containerregistry.Image
	TKENotifyController   containerregistry.Image
	TKEPlatformAPI        containerregistry.Image
	TKEPlatformController containerregistry.Image
	TKERegistryAPI        containerregistry.Image
}

func (c Components) Get(name string) *containerregistry.Image {
	v := reflect.ValueOf(c)
	for i := 0; i < v.NumField(); i++ {
		v, _ := v.Field(i).Interface().(containerregistry.Image)
		if v.Name == name {
			return &v
		}
	}
	return nil
}

var Version = version.Get().GitVersion

var components = Components{
	Busybox:    containerregistry.Image{Name: "busybox", Tag: "1.31.0"},
	Keepalived: containerregistry.Image{Name: "keepalived", Tag: "2.0.16-r0"},
	InfluxDB:   containerregistry.Image{Name: "influxdb", Tag: "1.7.6-alpine"},

	ProviderRes: containerregistry.Image{Name: "provider-res", Tag: "v1.14.6-1"},

	BaremetalClusterProvider: containerregistry.Image{Name: "baremetal-cluster-provider", Tag: Version},
	BaremetalMachineProvider: containerregistry.Image{Name: "baremetal-machine-provider", Tag: Version},

	TKEAuth:               containerregistry.Image{Name: "tke-auth", Tag: Version},
	TKEBusinessAPI:        containerregistry.Image{Name: "tke-business-api", Tag: Version},
	TKEBusinessController: containerregistry.Image{Name: "tke-business-controller", Tag: Version},
	TKEGateway:            containerregistry.Image{Name: "tke-gateway", Tag: Version},
	TKEMonitorAPI:         containerregistry.Image{Name: "tke-monitor-api", Tag: Version},
	TKEMonitorController:  containerregistry.Image{Name: "tke-monitor-controller", Tag: Version},
	TKENotifyAPI:          containerregistry.Image{Name: "tke-notify-api", Tag: Version},
	TKENotifyController:   containerregistry.Image{Name: "tke-notify-controller", Tag: Version},
	TKEPlatformAPI:        containerregistry.Image{Name: "tke-platform-api", Tag: Version},
	TKEPlatformController: containerregistry.Image{Name: "tke-platform-controller", Tag: Version},
	TKERegistryAPI:        containerregistry.Image{Name: "tke-registry-api", Tag: Version},
}

func List() []string {
	var items []string
	v := reflect.ValueOf(components)
	for i := 0; i < v.NumField(); i++ {
		v, _ := v.Field(i).Interface().(containerregistry.Image)
		items = append(items, v.BaseName())
	}
	sort.Strings(items)

	return items
}

func Get() Components {
	return components
}
