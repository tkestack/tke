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
	"sort"

	"tkestack.io/tke/pkg/util/containerregistry"
)

const (
	// LatestVersion is latest version of addon.
	LatestVersion = "v1.0.2"
)

var versionMap = map[string]containerregistry.Image{
	LatestVersion: containerregistry.Image{Name: "csi-operator", Tag: LatestVersion},
}

func ListImages() []string {
	items := make([]string, 0, len(versionMap))
	keys := make([]string, 0, len(versionMap))
	for key := range versionMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		items = append(items, versionMap[key].FullName())
	}

	return items
}

func ListVersions() []string {
	keys := make([]string, 0, len(versionMap))
	for key := range versionMap {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func GetImage(version string) *containerregistry.Image {
	cv, ok := versionMap[version]
	if !ok {
		panic(fmt.Sprintf("the component version definition corresponding to version %s could not be found", version))
	}
	return &cv
}
