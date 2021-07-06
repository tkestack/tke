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

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"tkestack.io/tke/pkg/spec"
)

func main() {
	env := os.Environ()
	env = append(env, fmt.Sprintf("ARCHS=%s", strings.Join(spec.Archs, " ")))
	env = append(env, fmt.Sprintf("OSS=%s", strings.Join(spec.OSs, " ")))
	env = append(env, fmt.Sprintf("K8S_VERSIONS=%s", strings.Join(spec.K8sVersionsWithV, " ")))
	env = append(env, fmt.Sprintf("DOCKER_VERSIONS=%s", strings.Join(spec.DockerVersions, " ")))
	env = append(env, fmt.Sprintf("CNI_PLUGINS_VERSIONS=%s", strings.Join(spec.CNIPluginsVersions, " ")))
	env = append(env, fmt.Sprintf("NERDCTL_VERSIONS=%s", strings.Join(spec.NerdctlVersions, " ")))
	env = append(env, fmt.Sprintf("NVIDIA_DRIVER_VERSIONS=%s", strings.Join(spec.NvidiaDriverVersions, " ")))
	env = append(env, fmt.Sprintf("NVIDIA_CONTAINER_RUNTIME_VERSIONS=%s", strings.Join(spec.NvidiaContainerRuntimeVersions, " ")))

	for _, one := range env {
		fmt.Println(one)

	}

	if len(os.Args) == 1 {
		return
	}

	for _, args := range os.Args[1:] {
		cmd := exec.Command("sh", "-c", args)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = env
		cmd.Run()
	}
}
