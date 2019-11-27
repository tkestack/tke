/*
 * Copyright 2019 THL A29 Limited, a Tencent company.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gpu

import (
	"fmt"
	"os"
	"path"

	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/util/ssh"

	clientset "k8s.io/client-go/kubernetes"
	"tkestack.io/tke/pkg/util/apiclient"
)

type NvidiaDriverOption struct {
	Version string
}

func InstallNvidiaDriver(s ssh.Interface, option *NvidiaDriverOption) error {
	basefile := fmt.Sprintf("NVIDIA-Linux-x86_64-%s.run", option.Version)
	srcFile := path.Join(constants.SrcDir, basefile)
	if _, err := os.Stat(srcFile); err != nil {
		return err
	}
	dstFile := path.Join(constants.DstTmpDir, basefile)
	err := s.CopyFile(srcFile, dstFile)
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("sh %s -s", dstFile)
	_, stderr, exit, err := s.Exec(cmd)
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	cmd = "nvidia-smi"
	_, stderr, exit, err = s.Exec(cmd)
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	return nil
}

type NvidiaContainerRuntimeOption struct {
	Version string
}

func InstallNvidiaContainerRuntime(s ssh.Interface, option *NvidiaContainerRuntimeOption) error {
	basefile := fmt.Sprintf("nvidia-container-runtime-%s.tgz", option.Version)
	srcFile := path.Join(constants.SrcDir, basefile)
	if _, err := os.Stat(srcFile); err != nil {
		return err
	}
	dstFile := path.Join(constants.DstTmpDir, basefile)
	err := s.CopyFile(srcFile, dstFile)
	if err != nil {
		return err
	}

	cmd := "tar xvaf %s -C /"
	_, stderr, exit, err := s.Execf(cmd, dstFile)
	if err != nil {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	cmd = "ln -sf /usr/bin/nvidia-container-toolkit /usr/bin/nvidia-container-runtime-hook"
	_, err = s.CombinedOutput(cmd)
	if err != nil {
		return fmt.Errorf("run cmd(%s) error:%s", cmd, err)
	}

	return nil
}

type NvidiaDevicePluginOption struct {
	Image string
}

func InstallNvidiaDevicePlugin(clientset clientset.Interface, option *NvidiaDevicePluginOption) error {
	err := apiclient.CreateResourceWithFile(clientset, constants.ManifestsDir+"gpu/nvidia-device-plugin.yaml", option)
	if err != nil {
		return err
	}
	return nil
}

func IsEnable(labels map[string]string) bool {
	return labels["nvidia-device-enable"] == "enable"
}

func MachineIsSupport(s ssh.Interface) bool {
	// https://wiki.debian.org/NvidiaGraphicsDrivers#NVIDIA_Proprietary_Driver
	_, err := s.CombinedOutput(`lspci -nn | egrep -i "3d|display|vga"`)
	return err == nil
}
