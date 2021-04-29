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

	"tkestack.io/tke/pkg/util/addon"
	"tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/log"
)

const (
	DefaultVersion = "v1.0.0"
	AddonName      = "log-collector"
)

var defaultComponents = Components{struct {
	Name string
	Tag  string
}{Name: "log-collector", Tag: "v1.1.0"}}

type Components struct {
	LogCollector containerregistry.Image
}

// GetLatestVersion returns latest version
func GetLatestVersion() string {
	version, err := addon.GetLatestVersion(AddonName)
	if err != nil {
		log.Errorf("%v", err)
		return DefaultVersion
	}
	return version
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

func List() []string {
	versionMap, err := addon.GetVersionMap(AddonName)
	if err != nil {
		log.Errorf("get version map error: %v", err)
		return []string{DefaultVersion}
	}
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
	versionMap, err := addon.GetVersionMap(AddonName)
	if err != nil {
		return err
	}
	_, ok := versionMap[version]
	if !ok {
		return fmt.Errorf("the component version definition corresponding to version %s could not be found", version)
	}
	return nil
}

func Get(version string) Components {
	versionMap, err := addon.GetVersionMap(AddonName)
	if err != nil {
		log.Errorf("get version map error: %v", err)
		return defaultComponents
	}
	cv, ok := versionMap[version]
	if !ok {
		log.Errorf("the component version definition corresponding to version %s could not be found，return default(%+v) instead", version, defaultComponents)
		return defaultComponents
	}
	return Components{LogCollector: cv}
}
