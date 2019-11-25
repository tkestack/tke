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
	"os"
	"path"

	"github.com/pkg/errors"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/ssh"
)

const (
	defaultVersion     = "v1.15.1"
	kubeadmConfigFile  = "kubeadm/kubeadm-config.yaml"
	kubeadmKubeletConf = "/usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf"

	joinControlePlaneCmd = `kubeadm join {{.ControlPlaneEndpoint}} \
--node-name={{.NodeName}} --token={{.BootstrapToken}} \
--control-plane --certificate-key={{.CertificateKey}} \
--skip-phases=control-plane-join/mark-control-plane \
--ignore-preflight-errors=ImagePull,Port-10250,FileContent--proc-sys-net-bridge-bridge-nf-call-iptables \
--discovery-token-unsafe-skip-ca-verification
`
	joinNodeCmd = `kubeadm join {{.ControlPlaneEndpoint}} \
--node-name={{.NodeName}} \
--token={{.BootstrapToken}} \
--ignore-preflight-errors=ImagePull,Port-10250,FileContent--proc-sys-net-bridge-bridge-nf-call-iptables \
--discovery-token-unsafe-skip-ca-verification
`
)

func Install(s ssh.Interface) error {
	basefile := fmt.Sprintf("kubeadm-%s.tar.gz", defaultVersion)
	srcFile := path.Join(constants.SrcDir, basefile)
	if _, err := os.Stat(srcFile); err != nil {
		return err
	}
	dstFile := path.Join(constants.DstTmpDir, basefile)
	err := s.CopyFile(srcFile, dstFile)
	if err != nil {
		return err
	}

	cmd := "tar xvaf %s -C %s "
	_, stderr, exit, err := s.Execf(cmd, dstFile, constants.DstBinDir)
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	data, err := util.ParseFileTemplate(path.Join(constants.ConfDir, "kubeadm/10-kubeadm.conf"), nil)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(data), kubeadmKubeletConf)
	if err != nil {
		return errors.Wrapf(err, "write %s error", kubeadmKubeletConf)
	}

	return nil
}

type InitOption struct {
	KubeadmConfigFileName string
	NodeName              string
	BootstrapToken        string
	CertificateKey        string

	ETCDImageTag         string
	CoreDNSImageTag      string
	KubernetesVersion    string
	ControlPlaneEndpoint string

	DNSDomain             string
	ServiceSubnet         string
	NodeCIDRMaskSize      int32
	ClusterCIDR           string
	ServiceClusterIPRange string

	APIServerExtraArgs         map[string]string
	ControllerManagerExtraArgs map[string]string
	SchedulerExtraArgs         map[string]string

	ImageRepository string
	ClusterName     string
}

func Init(s ssh.Interface, option *InitOption, extraCmd string) error {
	configData, err := util.ParseFileTemplate(path.Join(constants.ConfDir, kubeadmConfigFile), option)
	if err != nil {
		return errors.Wrap(err, "parse kubeadm config file error")
	}
	err = s.WriteFile(bytes.NewReader(configData), option.KubeadmConfigFileName)
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("kubeadm init phase %s --config=%s", extraCmd, option.KubeadmConfigFileName)
	stdout, stderr, exit, err := s.Exec(cmd)
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}
	log.Info(stdout)

	return nil
}

type JoinControlePlaneOption struct {
	NodeName             string
	BootstrapToken       string
	CertificateKey       string
	ControlPlaneEndpoint string
	OIDCCA               []byte
}

func JoinControlePlane(s ssh.Interface, option *JoinControlePlaneOption) error {
	if len(option.OIDCCA) != 0 { // ensure oidc ca exists becase kubeadm reset probably delete it!
		err := s.WriteFile(bytes.NewReader(option.OIDCCA), constants.OIDCCACertFile)
		if err != nil {
			return err
		}
	}

	cmd, err := util.ParseTemplate(joinControlePlaneCmd, option)
	if err != nil {
		return errors.Wrap(err, "parse joinControlePlaneCmd error")
	}
	stdout, stderr, exit, err := s.Exec(string(cmd))
	if err != nil || exit != 0 {
		_, _, _, _ = s.Exec("kubeadm reset -f")
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}
	log.Info(stdout)

	return nil
}

type JoinNodeOption struct {
	NodeName             string
	BootstrapToken       string
	ControlPlaneEndpoint string
}

func JoinNode(s ssh.Interface, option *JoinNodeOption) error {
	_, _, _, _ = s.Exec("kubeadm reset -f")

	cmd, err := util.ParseTemplate(joinNodeCmd, option)
	if err != nil {
		return errors.Wrap(err, "parse joinNodeCmd error")
	}
	stdout, stderr, exit, err := s.Exec(string(cmd))
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}
	log.Info(stdout)

	return nil
}
