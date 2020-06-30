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
	"strings"

	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	"tkestack.io/tke/pkg/util/ssh"
)

type Option struct {
	Version        string
	RegistryDomain string
}

func PullKubernetesImages(s ssh.Interface, option *Option) error {
	images := images.ListKubernetesImageFullNamesWithVerion(option.Version)
	if len(images) == 0 {
		return fmt.Errorf("images is empty")
	}

	for _, image := range images {
		cmd := fmt.Sprintf("docker pull %s", image)
		_, err := s.CombinedOutput(cmd)
		if err != nil {
			if strings.Contains(err.Error(), "502 Bad Gateway") {
				cmd = fmt.Sprintf(" docker info | grep Proxy")
				output, _ := s.CombinedOutput(cmd)
				return fmt.Errorf(`pull image fail: %s. maybe set no_proxy for registry(%s,*.%s) in docker dameon.
					docker info:%s. see: https://docs.docker.com/config/daemon/systemd/#httphttps-proxy`,
					err, option.RegistryDomain, option.RegistryDomain, output)
			}

			return err
		}
	}

	return nil
}
