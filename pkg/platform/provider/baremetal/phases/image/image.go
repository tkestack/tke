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

package image

import (
	"fmt"
	"k8s.io/klog"
	"strings"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/ssh"
)

type Option struct {
	Version        string
	RegistryDomain string
}

func PullKubernetesImages(c *v1.Cluster, s ssh.Interface, option *Option) error {
	images := images.ListKubernetesImageFullNamesWithVerion(option.Version)
	if len(images) == 0 {
		return fmt.Errorf("images is empty")
	}

	for _, image := range images {
		cmd := ""
		output, err := s.CombinedOutput("env | grep DOCKERHUBACTIONAUTH= | awk -F 'DOCKERHUBACTIONAUTH=' '{print $2}'")
		klog.Info("=====AUTH ENV STRING")
		klog.Info(output)
		if err != nil {
			return err
		}
		AUTHENV := strings.TrimSpace(string(output))
		if c.Cluster.Spec.Features.ContainerRuntime == platformv1.Docker {
			cmd = fmt.Sprintf("docker pull %s", image)
		} else {
			if len(AUTHENV) == 0 {
				cmd = fmt.Sprintf("crictl pull %s", image)
			} else {
				cmd = fmt.Sprintf("crictl pull --auth=%s %s", AUTHENV, image)
			}
		}
		_, err = s.CombinedOutput(cmd)
		if err != nil {
			str := err.Error()
			str = strings.ReplaceAll(str, AUTHENV, "")
			return fmt.Errorf("pull kubernetes image failed with error: %s", str)
		}
	}
	return nil
}
