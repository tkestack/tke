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
	"fmt"
	"reflect"
	"sort"

	"tkestack.io/tke/pkg/util/containerregistry"
)

const (
	// LatestVersion is latest version of addon.
	LatestVersion = "v1.0.0"
)

type Components struct {
	PrometheusService                containerregistry.Image
	KubeStateService                 containerregistry.Image
	NodeExporterService              containerregistry.Image
	AlertManagerService              containerregistry.Image
	ConfigMapReloadWorkLoad          containerregistry.Image
	PrometheusOperatorService        containerregistry.Image
	PrometheusConfigReloaderWorkload containerregistry.Image
	PrometheusBeatWorkLoad           containerregistry.Image
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

var versionMap = map[string]Components{
	LatestVersion: {
		PrometheusService:                containerregistry.Image{Name: "prometheus", Tag: "v2.11.0"},
		KubeStateService:                 containerregistry.Image{Name: "kube-state-metrics", Tag: "v1.6.0"},
		NodeExporterService:              containerregistry.Image{Name: "node-exporter", Tag: "v0.15.2"},
		AlertManagerService:              containerregistry.Image{Name: "alertmanager", Tag: "v0.18.0"},
		ConfigMapReloadWorkLoad:          containerregistry.Image{Name: "configmap-reload", Tag: "v0.1"},
		PrometheusOperatorService:        containerregistry.Image{Name: "prometheus-operator", Tag: "v0.31.1"},
		PrometheusConfigReloaderWorkload: containerregistry.Image{Name: "prometheus-config-reloader", Tag: "v0.31.1"},
		PrometheusBeatWorkLoad:           containerregistry.Image{Name: "prometheusbeat", Tag: "6.4.1"},
	},
}

func List() []string {
	items := make([]string, 0, len(versionMap))
	keys := make([]string, 0, len(versionMap))
	for key := range versionMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		v := reflect.ValueOf(versionMap[key])
		for i := 0; i < v.NumField(); i++ {
			v, _ := v.Field(i).Interface().(containerregistry.Image)
			items = append(items, v.BaseName())
		}
	}

	return items
}

func Validate(version string) error {
	_, ok := versionMap[version]
	if !ok {
		return fmt.Errorf("the component version definition corresponding to version %s could not be found", version)
	}
	return nil
}

func Get(version string) Components {
	cv, ok := versionMap[version]
	if !ok {
		panic(fmt.Sprintf("the component version definition corresponding to version %s could not be found", version))
	}
	return cv
}
