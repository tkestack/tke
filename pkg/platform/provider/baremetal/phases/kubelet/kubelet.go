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

package kubelet

import (
	"bytes"
	"fmt"
	"path"

	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/res"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/pkg/util/supervisor"
	"tkestack.io/tke/pkg/util/template"
)

func Install(s ssh.Interface, version string) (err error) {
	dstFile, err := res.KubernetesNode.CopyToNode(s, version)
	if err != nil {
		return err
	}

	for _, file := range []string{"kubelet", "kubectl"} {
		file = path.Join(constants.DstBinDir, file)
		if ok, err := s.Exist(file); err == nil && ok {
			backupFile, err := ssh.BackupFile(s, file)
			if err != nil {
				return fmt.Errorf("backup file %q error: %w", file, err)
			}
			defer func() {
				if err == nil {
					return
				}
				if err = ssh.RestoreFile(s, backupFile); err != nil {
					err = fmt.Errorf("restore file %q error: %w", backupFile, err)
				}
			}()
		}
	}

	cmd := "tar xvaf %s -C %s --strip-components=3"
	_, stderr, exit, err := s.Execf(cmd, dstFile, constants.DstBinDir)
	if err != nil {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	serviceData, err := template.ParseFile(path.Join(constants.ConfDir, "kubelet/kubelet.service"), nil)
	if err != nil {
		return err
	}

	ss := &supervisor.SystemdSupervisor{Name: "kubelet", SSH: s}
	err = ss.Deploy(bytes.NewReader(serviceData))
	if err != nil {
		return err
	}

	err = ss.Start()
	if err != nil {
		return err
	}

	cmd = "kubectl completion bash > /etc/bash_completion.d/kubectl"
	_, err = s.CombinedOutput(cmd)
	if err != nil {
		return err
	}

	return nil
}
