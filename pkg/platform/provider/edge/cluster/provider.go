/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2022 Tencent. All Rights Reserved.
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
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/AlekSi/pointer"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/server/mux"
	"k8s.io/client-go/tools/clientcmd"

	platformv1client "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	bcluster "tkestack.io/tke/pkg/platform/provider/baremetal/cluster"
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
	name = "Edge"
)

func RegisterProvider() {
	p, err := NewProvider()
	if err != nil {
		log.Errorf("init edge cluster provider error: %s", err)
		return
	}
	clusterprovider.Register(p.Name(), p)
}

type Provider struct {
	*clusterprovider.DelegateProvider

	// batemetal
	bCluster *bcluster.Provider
	bconfig  *config.Config //todo edge cluster config
}

var _ clusterprovider.Provider = &Provider{}

func NewProvider() (*Provider, error) {
	p := new(Provider)

	var err error
	p.bCluster, err = bcluster.NewProvider()
	if err != nil {
		log.Errorf("init cluster provider error: %s", err)
		return nil, err
	}

	p.DelegateProvider = &clusterprovider.DelegateProvider{
		ProviderName: name,

		CreateHandlers: []clusterprovider.Handler{
			p.bCluster.EnsureCopyFiles,
			p.bCluster.EnsurePreClusterInstallHook,
			p.bCluster.EnsurePreInstallHook,

			// configure system
			p.bCluster.EnsureRegistryHosts,
			p.bCluster.EnsureInitAPIServerHost,
			p.bCluster.EnsureKernelModule,
			p.bCluster.EnsureSysctl,
			p.bCluster.EnsureDisableSwap,
			p.bCluster.EnsurePreflight, // wait basic setting done

			p.bCluster.EnsureClusterComplete,

			// install packages
			p.bCluster.EnsureNvidiaDriver,
			p.bCluster.EnsureNvidiaContainerRuntime,
			p.bCluster.EnsureContainerRuntime,
			p.bCluster.EnsureKubernetesImages,
			p.bCluster.EnsureKubelet,
			p.bCluster.EnsureCNIPlugins,
			p.bCluster.EnsureConntrackTools,
			p.bCluster.EnsureKubeadm,
			p.bCluster.EnsureKeepalivedInit,
			p.bCluster.EnsureThirdPartyHAInit,
			p.bCluster.EnsureAuthzWebhook,
			p.bCluster.EnsurePrepareForControlplane,

			p.bCluster.EnsureKubeadmInitPhaseKubeletStart,
			p.bCluster.EnsureKubeadmInitPhaseCerts,
			p.bCluster.EnsureStoreCredential,
			p.bCluster.EnsureKubeconfig, // for upload
			p.bCluster.EnsureKubeadmInitPhaseKubeConfig,
			p.bCluster.EnsureKubeadmInitPhaseControlPlane,
			p.bCluster.EnsureKubeadmInitPhaseETCD,
			p.bCluster.EnsureKubeadmInitPhaseWaitControlPlane,
			p.bCluster.EnsureKubeadmInitPhaseUploadConfig,
			p.bCluster.EnsureKubeadmInitPhaseUploadCerts,
			p.bCluster.EnsureKubeadmInitPhaseBootstrapToken,
			p.bCluster.EnsureKubeadmInitPhaseAddon,

			//p.bCluster.EnsureGalaxy,
			//p.bCluster.EnsureCilium,
			p.EnsureEdgeFlannel,

			p.bCluster.EnsureJoinPhasePreflight,
			p.bCluster.EnsureJoinPhaseControlPlanePrepare,
			p.bCluster.EnsureJoinPhaseKubeletStart,
			p.bCluster.EnsureJoinPhaseControlPlaneJoinETCD,
			p.bCluster.EnsureJoinPhaseControlPlaneJoinUpdateStatus,

			p.bCluster.EnsurePatchAnnotation, // wait rest master ready
			p.bCluster.EnsureMarkControlPlane,
			p.bCluster.EnsureKeepalivedWithLBOption,
			p.bCluster.EnsureThirdPartyHA,
			p.bCluster.EnsureModifyAPIServerHost,

			// deploy apps
			p.bCluster.EnsureNvidiaDevicePlugin,
			p.bCluster.EnsureGPUManager,
			p.bCluster.EnsureCSIOperator,
			p.bCluster.EnsureMetricsServer,

			// provider SuperEdge edge cluster
			p.EnsurePrepareEgdeCluster, // Prepare EgdeCluster
			p.EnsureApplyEdgeApps,      // Add on SuperEdge EdgeApps

			p.bCluster.EnsureCleanup,
			p.bCluster.EnsureCreateClusterMark,
			p.bCluster.EnsurePostInstallHook,
			p.bCluster.EnsurePostClusterInstallHook,
		},
		UpdateHandlers: []clusterprovider.Handler{
			p.bCluster.EnsureAPIServerCert,
			p.bCluster.EnsureRenewCerts,
			p.bCluster.EnsureStoreCredential,
			p.bCluster.EnsureKeepalivedWithLBOption,
			p.bCluster.EnsureThirdPartyHA,
		},
		UpgradeHandlers: []clusterprovider.Handler{
			p.bCluster.EnsurePreClusterUpgradeHook,
			p.bCluster.EnsureUpgradeCoreDNS,
			p.bCluster.EnsureUpgradeControlPlaneNode,
			p.bCluster.EnsurePostClusterUpgradeHook,
		},
		ScaleDownHandlers: []clusterprovider.Handler{
			p.bCluster.EnsureRemoveETCDMember,
			p.bCluster.EnsureRemoveNode,
		},
		DeleteHandlers: []clusterprovider.Handler{
			p.bCluster.EnsureRemoveMachine,
			p.bCluster.EnsureCleanClusterMark,
		},
	}
	p.bCluster.ScaleUpHandlers = p.bCluster.CreateHandlers //todo: edge ScaleUpHandlers

	cfg, err := config.New(constants.ConfigFile)
	if err != nil {
		return nil, err
	}
	p.bconfig = cfg

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

func (p *Provider) ping(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprint(resp, "pong")
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

	if p.bconfig.AuditEnabled() {
		if !cluster.AuthzWebhookEnabled() {
			cluster.Spec.Features.AuthzWebhookAddr = &platform.AuthzWebhookAddr{Builtin: &platform.
				BuiltinAuthzWebhookAddr{}}
		}
	}

	if p.bconfig.BusinessEnabled() {
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
	if p.bconfig.Feature.SkipConditions != nil {
		cluster.Spec.Features.SkipConditions = append(cluster.Spec.Features.SkipConditions, p.bconfig.Feature.SkipConditions...)
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
