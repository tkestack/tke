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

package util

import (
	"k8s.io/klog"
	"os/exec"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/test/util/cloudprovider"
)

func RunCmd(cmd string) (string, error) {
	klog.Info("Run cmd: ", cmd)
	command := exec.Command("bash", "-c", cmd)
	out, err := command.CombinedOutput()
	klog.Info("Cmd result: ", string(out))
	return string(out), err
}

func RunCmdOnNode(ins cloudprovider.Instance, cmd string) (string, error) {
	klog.Info("Run cmd: ", cmd, ". Node: ", ins.InstanceID)
	s, err := ssh.New(&ssh.Config{
		User:     ins.Username,
		Password: ins.Password,
		Host:     ins.PublicIP,
		Port:     int(ins.Port),
	})
	if err != nil {
		return "", err
	}
	out, err := s.CombinedOutput(cmd)
	klog.Info("Cmd result: ", string(out))
	return string(out), err
}
