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
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	platformv1client "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/config"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/validation"
	machineprovider "tkestack.io/tke/pkg/platform/provider/machine"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/strategicpatch"
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
			p.EnsureManifestDir,

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
			p.EnsureDisableOffloading, // will remove it when upgrade to k8s v1.18.5
			p.EnsurePostInstallHook,
		},
		UpdateHandlers: []machineprovider.Handler{
			p.EnsurePreUpgradeHook,
			p.EnsureUpgrade,
			p.EnsurePostUpgradeHook,
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

func (p *Provider) OnUpdate(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	p.ensureSyncMachineNodeLabel(ctx, machine)

	return p.DelegateProvider.OnUpdate(ctx, machine, cluster)
}

func (p *Provider) ensureSyncMachineNodeLabel(ctx context.Context, machine *platformv1.Machine) {
	cluster, err := typesv1.GetClusterByName(ctx, p.platformClient, machine.Spec.ClusterName)
	if err != nil {
		log.FromContext(ctx).Error(err, "sync Machine node label error")
		return
	}

	client, err := cluster.Clientset()
	if err != nil {
		log.FromContext(ctx).Error(err, "sync Machine node label error")
		return
	}

	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		node, err := client.CoreV1().Nodes().Get(ctx, machine.Spec.IP, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return nil
			}
			return err
		}

		labels := node.GetLabels()
		_, ok := labels[string(apiclient.LabelMachineIPV4)]
		if ok {
			return nil
		}

		oldNode := node.DeepCopy()
		labels[string(apiclient.LabelMachineIPV4)] = machine.Spec.IP
		node.SetLabels(labels)

		patchBytes, err := strategicpatch.GetPatchBytes(oldNode, node)
		if err != nil {
			return fmt.Errorf("GetPatchBytes for node error: %w", err)
		}

		_, err = client.CoreV1().Nodes().Patch(ctx, node.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
		return err
	})

	if err != nil {
		log.FromContext(ctx).Error(err, "sync Machine node label error")
	}
}
