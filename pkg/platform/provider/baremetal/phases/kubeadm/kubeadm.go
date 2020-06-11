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

package kubeadm

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	kubeadmv1beta2 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeadm/v1beta2"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/res"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/pkg/util/template"
)

const (
	kubeadmKubeletConf = "/usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf"

	initCmd = `kubeadm init phase {{.Phase}} --config={{.Config}}`
	joinCmd = `kubeadm join phase {{.Phase}} --config={{.Config}}`
)

var (
	ignoreErrors = []string{
		"ImagePull",
		"Port-10250",
		"FileContent--proc-sys-net-bridge-bridge-nf-call-iptables",
		"DirAvailable--etc-kubernetes-manifests",
	}
)

func Install(s ssh.Interface, version string) error {
	dstFile, err := res.Kubeadm.CopyToNode(s, version)
	if err != nil {
		return err
	}

	cmd := "tar xvaf %s -C %s "
	_, stderr, exit, err := s.Execf(cmd, dstFile, constants.DstBinDir)
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	data, err := template.ParseFile(path.Join(constants.ConfDir, "kubeadm/10-kubeadm.conf"), nil)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(data), kubeadmKubeletConf)
	if err != nil {
		return errors.Wrapf(err, "write %s error", kubeadmKubeletConf)
	}

	return nil
}

func Init(s ssh.Interface, kubeadmConfig *InitConfig, phase string) error {
	configData, err := kubeadmConfig.Marshal()
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(configData), constants.KubeadmConfigFileName)
	if err != nil {
		return err
	}

	cmd, err := template.ParseString(initCmd, map[string]interface{}{
		"Phase":  phase,
		"Config": constants.KubeadmConfigFileName,
	})
	if err != nil {
		return errors.Wrap(err, "parse initCmd error")
	}
	out, err := s.CombinedOutput(string(cmd))
	if err != nil {
		return fmt.Errorf("kubeadm.Init error: %w", err)
	}
	log.Info(string(out))

	return nil
}

func Join(s ssh.Interface, config *kubeadmv1beta2.JoinConfiguration, phase string) error {
	configData, err := MarshalToYAML(config)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(configData), constants.KubeadmConfigFileName)
	if err != nil {
		return err
	}
	if phase == "preflight" {
		phase = fmt.Sprintf("preflight --ignore-preflight-errors=%s", strings.Join(ignoreErrors, ","))
	}

	cmd, err := template.ParseString(joinCmd, map[string]interface{}{
		"Phase":  phase,
		"Config": constants.KubeadmConfigFileName,
	})
	if err != nil {
		return errors.Wrap(err, "parse joinCmd error")
	}
	out, err := s.CombinedOutput(string(cmd))
	if err != nil {
		return fmt.Errorf("kubeadm.Join error: %w", err)
	}
	log.Info(string(out))

	return nil
}

func RenewCerts(s ssh.Interface) error {
	err := fixKubeadmBug1753(s)
	if err != nil {
		return fmt.Errorf("fixKubeadmBug1753(https://github.com/kubernetes/kubeadm/issues/1753) error: %w", err)
	}

	cmd := fmt.Sprintf("kubeadm alpha certs renew all --config=%s", constants.KubeadmConfigFileName)
	_, err = s.CombinedOutput(cmd)
	if err != nil {
		return err
	}

	err = RestartControlPlane(s)
	if err != nil {
		return err
	}

	return nil
}

// https://github.com/kubernetes/kubeadm/issues/1753
func fixKubeadmBug1753(s ssh.Interface) error {
	needUpdate := false

	data, err := s.ReadFile(constants.KubeletKubeConfigFileName)
	if err != nil {
		return err
	}
	kubeletKubeconfig, err := clientcmd.Load(data)
	if err != nil {
		return err
	}
	for _, info := range kubeletKubeconfig.AuthInfos {
		if info.ClientKeyData == nil && info.ClientCertificateData == nil {
			continue
		}

		info.ClientKeyData = []byte{}
		info.ClientCertificateData = []byte{}
		info.ClientKey = constants.KubeletClientCurrent
		info.ClientCertificate = constants.KubeletClientCurrent

		needUpdate = true
	}

	if needUpdate {
		data, err := runtime.Encode(clientcmdlatest.Codec, kubeletKubeconfig)
		if err != nil {
			return err
		}
		err = s.WriteFile(bytes.NewReader(data), constants.KubeletKubeConfigFileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func RestartControlPlane(s ssh.Interface) error {
	targets := []string{"kube-apiserver", "kube-controller-manager", "kube-scheduler"}
	for _, one := range targets {
		err := RestartContainerByFilter(s, DockerFilterForControlPlane(one))
		if err != nil {
			return err
		}
	}

	return nil
}

func DockerFilterForControlPlane(name string) string {
	return fmt.Sprintf("label=io.kubernetes.container.name=%s", name)
}

func RestartContainerByFilter(s ssh.Interface, filter string) error {
	cmd := fmt.Sprintf("docker rm -f $(docker ps -q -f '%s')", filter)
	_, err := s.CombinedOutput(cmd)
	if err != nil {
		return err
	}

	err = wait.PollImmediate(5*time.Second, 5*time.Minute, func() (bool, error) {
		cmd = fmt.Sprintf("docker ps -q -f '%s'", filter)
		output, err := s.CombinedOutput(cmd)
		if err != nil {
			return false, nil
		}
		if len(output) == 0 {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("restart container(%s) error: %w", filter, err)
	}

	return nil
}
