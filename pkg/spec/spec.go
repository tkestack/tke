/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package spec

import (
	"github.com/thoas/go-funk"
)

var (
	Archs         = []string{"amd64", "arm64"}
	Arm64         = "arm64"
	Arm64Variants = []string{"v8", "unknown"}
	OSs           = []string{"linux"}

	K8sVersions = []string{"1.18.2", "1.16.9", "1.14.10"}
	// K8sValidVersions for backward compatibility.
	K8sValidVersions = append(K8sVersions, []string{"1.16.6"}...)
	K8sVersionsWithV = funk.Map(K8sVersions, func(s string) string {
		return "v" + s
	}).([]string)
	KubeadmVersions = []string{K8sVersionsWithV[len(K8sVersionsWithV)-1]}

	K8sVersionConstraint = ">= 1.10"

	DockerVersions                 = []string{"18.09.9"}
	CNIPluginsVersions             = []string{"v0.8.5"}
	NvidiaDriverVersions           = []string{"440.31"}
	NvidiaContainerRuntimeVersions = []string{"3.1.4"}
)
