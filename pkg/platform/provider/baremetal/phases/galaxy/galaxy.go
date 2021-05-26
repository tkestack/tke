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

package galaxy

import (
	"context"
	errs "errors"
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"net"
	"os/exec"
	"strings"
	"time"
	baremetalconstants "tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/galaxy/images"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/log"
)

const (
	BackendProviderFlannel = "flannel"
	BackendProviderCalico  = "calico"
	BackendTypeVxLAN       = "vxlan"
)

// Option for galaxy
type Option struct {
	Version          string
	NodeCIDR         string
	NetDevice        string
	BackendType      string
	BackendProvider  string
	NodeCIDRMaskSize int32
}

// Install to install the galaxy workload
func Install(ctx context.Context, clientset apiclient.KubeInterfaces, option *Option) error {

	if err := installNetworkProvider(ctx, clientset, option); err != nil {
		return err
	}
	return installGalaxy(ctx, clientset, option)
}

func installNetworkProvider(ctx context.Context, clientset apiclient.KubeInterfaces, option *Option) error {
	switch provider := option.BackendProvider; provider {
	case BackendProviderFlannel:
		return installFlannel(ctx, clientset, option)
	case BackendProviderCalico:
		return installCalico(ctx, clientset, option)
	default:
		return errs.New(fmt.Sprintf("unknown network provider: %s", provider))
	}
}

func installFlannel(ctx context.Context, clientset apiclient.KubeInterfaces, option *Option) error {
	// old flannel interface should be deleted
	if err := cleanFlannelInterfaces(); err != nil {
		return err
	}

	err := apiclient.CreateResourceWithDir(ctx, clientset, baremetalconstants.FlannelManifest,
		map[string]interface{}{
			"BackendType":  option.BackendType,
			"Network":      option.NodeCIDR,
			"FlannelImage": images.Get(option.Version).Flannel.FullName(),
		})
	if err != nil {
		return err
	}
	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDaemonset(ctx, clientset, "kube-system", "flannel")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func cleanFlannelInterfaces() error {
	var err error
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, iface := range ifaces {
		if strings.Contains(iface.Name, "flannel") {
			cmd := exec.Command("ip", "link", "delete", iface.Name)
			if err := cmd.Run(); err != nil {
				log.Errorf("fail to delete link %s : %v", iface.Name, err)
			}
		}
	}
	return err
}

func installCalico(ctx context.Context, clientset apiclient.KubeInterfaces, option *Option) error {
	err := apiclient.CreateResourceWithDir(ctx, clientset, baremetalconstants.CalicoManifest,
		map[string]interface{}{
			"CalicoCNIImage":            images.Get(option.Version).CalicoCNI.FullName(),
			"CalicoNodeImage":           images.Get(option.Version).CalicoNode.FullName(),
			"CalicoFlexvolDriverImage":  images.Get(option.Version).CalicoFlexvolDriver.FullName(),
			"CalicoKubeControllerImage": images.Get(option.Version).CalicoKubeControllers.FullName(),
			"ClusterCIDR":               option.NodeCIDR,
			"BackendType":               option.BackendType,
			"NodeCIDRMaskSize":          option.NodeCIDRMaskSize,
		})
	if err != nil {
		return err
	}
	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDaemonset(ctx, clientset, "kube-system", "calico-node")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func installGalaxy(ctx context.Context, clientset apiclient.KubeInterfaces, option *Option) error {
	err := apiclient.CreateResourceWithDir(ctx, clientset, baremetalconstants.GalaxyManifest,
		map[string]interface{}{
			"BackendProvider":   option.BackendProvider,
			"DeviceName":        option.NetDevice,
			"GalaxyDaemonImage": images.Get(option.Version).GalaxyDaemon.FullName(),
		})
	if err != nil {
		return err
	}
	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDaemonset(ctx, clientset, "kube-system", "galaxy")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}
