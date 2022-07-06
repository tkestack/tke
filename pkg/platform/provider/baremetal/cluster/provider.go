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

package cluster

import (
	"path"
	"strings"

	"github.com/AlekSi/pointer"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/server/mux"
	"k8s.io/client-go/tools/clientcmd"

	platformv1client "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/provider/baremetal/config"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	csioperatorimage "tkestack.io/tke/pkg/platform/provider/baremetal/phases/csioperator/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/validation"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/platform/types"
	"tkestack.io/tke/pkg/spec"
	"tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/log"
)

const (
	name = "Baremetal"
)

func RegisterProvider() {
	p, err := NewProvider()
	if err != nil {
		log.Errorf("init cluster provider error: %s", err)
		return
	}
	clusterprovider.Register(p.Name(), p)
}

type Provider struct {
	*clusterprovider.DelegateProvider

	Config *config.Config
}

var _ clusterprovider.Provider = &Provider{}

func NewProvider() (*Provider, error) {
	p := new(Provider)

	p.DelegateProvider = &clusterprovider.DelegateProvider{
		ProviderName: name,

		CreateHandlers: []clusterprovider.Handler{
			p.EnsureCopyFiles,
			p.EnsurePreClusterInstallHook,
			p.EnsurePreInstallHook,

			// configure system
			p.EnsureRegistryHosts,
			p.EnsureInitAPIServerHost,
			p.EnsureKernelModule,
			p.EnsureSysctl,
			p.EnsureDisableSwap,
			p.EnsurePreflight, // wait basic setting done

			p.EnsureClusterComplete,

			// install packages
			p.EnsureNvidiaDriver,
			p.EnsureNvidiaContainerRuntime,
			p.EnsureContainerRuntime,
			p.EnsureKubernetesImages,
			p.EnsureKubelet,
			p.EnsureCNIPlugins,
			p.EnsureConntrackTools,
			p.EnsureKubeadm,
			p.EnsureKeepalivedInit,
			p.EnsureThirdPartyHAInit,
			p.EnsureAuthzWebhook,
			p.EnsurePrepareForControlplane,

			p.EnsureKubeadmInitPhaseKubeletStart,
			p.EnsureKubeadmInitPhaseCerts,
			p.EnsureStoreCredential,
			p.EnsureKubeconfig, // for upload
			p.EnsureKubeadmInitPhaseKubeConfig,
			p.EnsureKubeadmInitPhaseControlPlane,
			p.EnsureKubeadmInitPhaseETCD,
			p.EnsureKubeadmInitPhaseWaitControlPlane,
			p.EnsureKubeadmInitPhaseUploadConfig,
			p.EnsureKubeadmInitPhaseUploadCerts,
			p.EnsureKubeadmInitPhaseBootstrapToken,
			p.EnsureKubeadmInitPhaseAddon,

			p.EnsureJoinPhasePreflight,
			p.EnsureJoinPhaseControlPlanePrepare,
			p.EnsureJoinPhaseKubeletStart,
			p.EnsureJoinPhaseControlPlaneJoinETCD,
			p.EnsureJoinPhaseControlPlaneJoinUpdateStatus,

			p.EnsureClusternetRegistration,

			p.EnsureGalaxy,
			p.EnsureCilium,

			p.EnsurePatchAnnotation, // wait rest master ready
			p.EnsureMarkControlPlane,
			p.EnsureKeepalivedWithLBOption,
			p.EnsureThirdPartyHA,
			p.EnsureModifyAPIServerHost,
			// deploy apps
			p.EnsureNvidiaDevicePlugin,
			p.EnsureGPUManager,
			p.EnsureCSIOperator,
			p.EnsureMetricsServer,

			p.EnsureCleanup,
			p.EnsureCreateClusterMark,
			p.EnsurePostInstallHook,
			p.EnsurePostClusterInstallHook,
		},
		UpdateHandlers: []clusterprovider.Handler{
			p.EnsureAPIServerCert,
			p.EnsureRenewCerts,
			p.EnsureStoreCredential,
			p.EnsureKeepalivedWithLBOption,
			p.EnsureThirdPartyHA,
		},
		UpgradeHandlers: []clusterprovider.Handler{
			p.EnsurePreClusterUpgradeHook,
			p.EnsureUpgradeCoreDNS,
			p.EnsureUpgradeControlPlaneNode,
			p.EnsurePostClusterUpgradeHook,
		},
		ScaleDownHandlers: []clusterprovider.Handler{
			p.EnsureRemoveETCDMember,
			p.EnsureRemoveNode,
		},
		DeleteHandlers: []clusterprovider.Handler{
			p.EnsureRemoveMachine,
			p.EnsureCleanClusterMark,
		},
	}
	p.ScaleUpHandlers = p.CreateHandlers

	cfg, err := config.New(constants.ConfigFile)
	if err != nil {
		return nil, err
	}
	p.Config = cfg

	containerregistry.Init(cfg.Registry.Domain, cfg.Registry.Namespace)

	// Run for compatibility with installer.
	// TODO: Installer reuse platform components
	if cfg.PlatformAPIClientConfig != "" {
		restConfig, err := clientcmd.BuildConfigFromFlags("", cfg.PlatformAPIClientConfig)
		if err != nil {
			log.Errorf("read PlatformAPIClientConfig error: %w", err)
		} else {
			p.PlatformClient, err = platformv1client.NewForConfig(restConfig)
			if err != nil {
				return nil, err
			}
		}
	}
	return p, nil
}

func (p *Provider) RegisterHandler(mux *mux.PathRecorderMux) {
	prefix := "/provider/" + strings.ToLower(p.Name())

	mux.HandleFunc(path.Join(prefix, "ping"), p.ping)
}

func (p *Provider) Validate(cluster *types.Cluster) field.ErrorList {
	return validation.ValidateCluster(p.PlatformClient, cluster)
}

func (p *Provider) ValidateUpdate(cluster *types.Cluster, oldCluster *types.Cluster) field.ErrorList {
	return validation.ValidateClusterUpdate(p.PlatformClient, cluster, oldCluster)
}

func (p *Provider) PreCreate(cluster *types.Cluster) error {
	if cluster.Spec.Version == "" {
		cluster.Spec.Version = spec.K8sVersions[0]
	}
	if cluster.Spec.ClusterCIDR == "" {
		cluster.Spec.ClusterCIDR = "10.244.0.0/16"
	}
	if cluster.Spec.NetworkDevice == "" {
		cluster.Spec.NetworkDevice = "eth0"

	}
	if cluster.Spec.Features.CSIOperator != nil {
		if cluster.Spec.Features.CSIOperator.Version == "" {
			cluster.Spec.Features.CSIOperator.Version = csioperatorimage.LatestVersion
		}
	}

	if p.Config.AuditEnabled() {
		if !cluster.AuthzWebhookEnabled() {
			cluster.Spec.Features.AuthzWebhookAddr = &platform.AuthzWebhookAddr{Builtin: &platform.
				BuiltinAuthzWebhookAddr{}}
		}
	}

	if p.Config.BusinessEnabled() {
		if !cluster.AuthzWebhookEnabled() {
			cluster.Spec.Features.AuthzWebhookAddr = &platform.AuthzWebhookAddr{Builtin: &platform.
				BuiltinAuthzWebhookAddr{}}
		}
	}

	if cluster.Spec.Properties.MaxClusterServiceNum == nil && cluster.Spec.ServiceCIDR == nil {
		cluster.Spec.Properties.MaxClusterServiceNum = pointer.ToInt32(256)
	}
	if cluster.Spec.Properties.MaxNodePodNum == nil {
		cluster.Spec.Properties.MaxNodePodNum = pointer.ToInt32(256)
	}
	// append SkipConditions when disable the cluster features.
	if cluster.Spec.Features.EnableCilium {
		cluster.Spec.Features.SkipConditions = append(cluster.Spec.Features.SkipConditions, "EnsureGalaxy")
	} else {
		cluster.Spec.Features.SkipConditions = append(cluster.Spec.Features.SkipConditions, "EnsureCilium")
	}
	if !cluster.Spec.Features.EnableMetricsServer {
		cluster.Spec.Features.SkipConditions = append(cluster.Spec.Features.SkipConditions, "EnsureMetricsServer")
	}
	if p.Config.Feature.SkipConditions != nil {
		cluster.Spec.Features.SkipConditions = append(cluster.Spec.Features.SkipConditions, p.Config.Feature.SkipConditions...)
	}

	if cluster.Spec.Etcd == nil {
		cluster.Spec.Etcd = &platform.Etcd{Local: &platform.LocalEtcd{}}
	}

	if cluster.Spec.Etcd.Local != nil {
		// reuse global etcd for tke components which create `etcd` service.
		cluster.Spec.Etcd.Local.ServerCertSANs = append(cluster.Spec.Etcd.Local.ServerCertSANs, "etcd")
		cluster.Spec.Etcd.Local.ServerCertSANs = append(cluster.Spec.Etcd.Local.ServerCertSANs, "etcd.kube-system")
	}

	return nil
}
