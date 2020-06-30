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

package machine

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/tools/clientcmd"
	platformv1client "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/provider/baremetal/config"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/validation"
	machineprovider "tkestack.io/tke/pkg/platform/provider/machine"
	"tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/log"
)

const (
	name = "Baremetal"
)

func init() {
	p, err := NewProvider()
	if err != nil {
		log.Errorf("init machine provider error: %s", err)
		return
	}
	machineprovider.Register(p.Name(), p)
}

type Provider struct {
	*machineprovider.DelegateProvider

	config         *config.Config
	platformClient platformv1client.PlatformV1Interface
}

func NewProvider() (*Provider, error) {
	p := new(Provider)

	p.DelegateProvider = &machineprovider.DelegateProvider{
		ProviderName: name,

		CreateHandlers: []machineprovider.Handler{
			p.EnsureCopyFiles,
			p.EnsurePreInstallHook,

			p.EnsureClean,
			p.EnsureRegistryHosts,
			p.EnsureKernelModule,
			p.EnsureSysctl,
			p.EnsureDisableSwap,

			p.EnsurePreflight, // wait basic setting done

			p.EnsureNvidiaDriver,
			p.EnsureNvidiaContainerRuntime,
			p.EnsureDocker,
			p.EnsureKubelet,
			p.EnsureCNIPlugins,
			p.EnsureConntrackTools,
			p.EnsureKubeadm,

			p.EnsureJoinPhasePreflight,
			p.EnsureJoinPhaseKubeletStart,

			p.EnsureKubeconfig,
			p.EnsureMarkNode,
			p.EnsureNodeReady,

			p.EnsurePostInstallHook,
		},
		UpdateHandlers: []machineprovider.Handler{
			p.EnsureUpgrade,
		},
	}

	cfg, err := config.New(constants.ConfigFile)
	if err != nil {
		return nil, err
	}
	p.config = cfg

	containerregistry.Init(cfg.Registry.Domain, cfg.Registry.Namespace)

	// Run for compatibility with installer.
	// TODO: Installer reuse platform components
	if cfg.PlatformAPIClientConfig != "" {
		restConfig, err := clientcmd.BuildConfigFromFlags("", cfg.PlatformAPIClientConfig)
		if err != nil {
			return nil, err
		}
		p.platformClient, err = platformv1client.NewForConfig(restConfig)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

var _ machineprovider.Provider = &Provider{}

func (p *Provider) Validate(machine *platform.Machine) field.ErrorList {
	return validation.ValidateMachine(machine)
}
