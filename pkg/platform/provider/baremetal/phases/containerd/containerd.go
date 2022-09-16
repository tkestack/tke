/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

package containerd

import (
	"bytes"
	"fmt"
	"path"

	"tkestack.io/tke/pkg/util/template"

	"github.com/pkg/errors"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/res"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/pkg/util/supervisor"
)

type Option struct {
	InsecureRegistries []string
	IsGPU              bool
	Root               string
	SandboxImage       string
	RegistryMirror     string
}

const (
	containerdConfigFile = "/etc/containerd/config.toml"
)

func Install(s ssh.Interface, option *Option) error {
	// add path to sudoers
	cmd := `sed -i '$a\Defaults secure_path = /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin' /etc/sudoers`
	_, err := s.CombinedOutput(cmd)
	if err != nil {
		return err
	}

	dstFile, err := res.Containerd.CopyToNodeWithDefault(s)
	if err != nil {
		return err
	}

	// Install containerd exclude cni binaries and cni config file.
	cmd = "tar xvaf %s -C %s --exclude=etc/cni --exclude=opt"
	_, stderr, exit, err := s.Execf(cmd, dstFile, "/")
	if err != nil {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	dstFile, err = res.Nerdctl.CopyToNodeWithDefault(s)
	if err != nil {
		return err
	}

	cmd = "tar xvaf %s -C /usr/local/bin/ "
	_, stderr, exit, err = s.Execf(cmd, dstFile)
	if err != nil {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	data, err := template.ParseFile(path.Join(constants.SrcDir, "containerd/config.toml"), option)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(data), containerdConfigFile)
	if err != nil {
		return errors.Wrapf(err, "write %s error", containerdConfigFile)
	}

	data, err = template.ParseFile(path.Join(constants.SrcDir, "containerd/containerd.service"), option)
	if err != nil {
		return err
	}
	ss := &supervisor.SystemdSupervisor{Name: "containerd", SSH: s}
	err = ss.Deploy(bytes.NewReader(data))
	if err != nil {
		return err
	}

	err = ss.Start()
	if err != nil {
		return err
	}

	return nil
}
