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

package docker

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"tkestack.io/tke/pkg/util/template"

	"github.com/pkg/errors"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/res"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/pkg/util/supervisor"
)

type Option struct {
	InsecureRegistries string
	RegistryDomain     string
	Options            string
	IsGPU              bool
	ExtraArgs          map[string]string
}

const (
	dockerDaemonFile = "/etc/docker/daemon.json"
)

func Install(s ssh.Interface, option *Option) error {
	dstFile, err := res.Docker.CopyToNodeWithDefault(s)
	if err != nil {
		return err
	}

	cmd := "tar xvaf %s -C %s --strip-components=1"
	_, stderr, exit, err := s.Execf(cmd, dstFile, constants.DstBinDir)
	if err != nil {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	var args []string
	for k, v := range option.ExtraArgs {
		args = append(args, fmt.Sprintf(`--%s="%s"`, k, v))
	}
	err = s.WriteFile(strings.NewReader(fmt.Sprintf("DOCKER_EXTRA_ARGS=%s", strings.Join(args, " "))), "/etc/sysconfig/docker")
	if err != nil {
		return err
	}

	data, err := template.ParseFile(path.Join(constants.ConfDir, "docker/daemon.json"), option)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(data), dockerDaemonFile)
	if err != nil {
		return errors.Wrapf(err, "write %s error", dockerDaemonFile)
	}

	data, err = template.ParseFile(path.Join(constants.ConfDir, "docker/docker.service"), option)
	if err != nil {
		return err
	}
	ss := &supervisor.SystemdSupervisor{Name: "docker", SSH: s}
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
