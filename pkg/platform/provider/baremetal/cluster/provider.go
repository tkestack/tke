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
	"context"
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
	delegatecluster "tkestack.io/tke/pkg/platform/provider/delegate/cluster"
	"tkestack.io/tke/pkg/platform/types"
	"tkestack.io/tke/pkg/spec"
	"tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/log"
)

const (
	name = "Baremetal"
)

func init() {
	p, err := NewProvider()
	if err != nil {
		log.Errorf("init cluster provider error: %s", err)
		return
	}
	clusterprovider.Register(p.Name(), p)
}

type Provider struct {
	*delegatecluster.DelegateProvider

	config         *config.Config
	platformClient platformv1client.PlatformV1Interface
}

var _ clusterprovider.Provider = &Provider{}

func NewProvider() (*Provider, error) {
	p := new(Provider)

	p.DelegateProvider = &delegatecluster.DelegateProvider{
		ProviderName: name,

		CreateHandlers: []delegatecluster.Handler{
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

			p.EnsureGalaxy,
			p.EnsureCilium,

			p.EnsureJoinPhasePreflight,
			p.EnsureJoinPhaseControlPlanePrepare,
			p.EnsureJoinPhaseKubeletStart,
			p.EnsureJoinPhaseControlPlaneJoinETCD,
			p.EnsureJoinPhaseControlPlaneJoinUpdateStatus,

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
			p.EnsureDisableOffloading, // will remove it when upgrade to k8s v1.18.5
			p.EnsurePostInstallHook,
			p.EnsurePostClusterInstallHook,
		},
		UpdateHandlers: []delegatecluster.Handler{
			p.EnsureAPIServerCert,
			p.EnsureRenewCerts,
			p.EnsureStoreCredential,
			p.EnsureKeepalivedWithLBOption,
			p.EnsureThirdPartyHA,
		},
		UpgradeHandlers: []delegatecluster.Handler{
			p.EnsurePreClusterUpgradeHook,
			p.EnsureUpgradeCoreDNS,
			p.EnsureUpgradeControlPlaneNode,
			p.EnsurePostClusterUpgradeHook,
		},
		ScaleDownHandlers: []delegatecluster.Handler{
			p.EnsureRemoveETCDMember,
			p.EnsureRemoveNode,
		},
		DeleteHandlers: []delegatecluster.Handler{
			p.EnsureCleanClusterMark,
		},
	}
	p.ScaleUpHandlers = p.CreateHandlers

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
			log.Errorf("read PlatformAPIClientConfig error: %w", err)
		} else {
			p.platformClient, err = platformv1client.NewForConfig(restConfig)
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

func (p *Provider) Validate(ctx context.Context, cluster *types.Cluster) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, p.DelegateProvider.Validate(ctx, cluster)...)
	allErrs = append(allErrs, validation.ValidateCluster(p.platformClient, cluster)...)

	return allErrs
}

func (p *Provider) ValidateUpdate(ctx context.Context, cluster *types.Cluster, oldCluster *types.Cluster) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, p.DelegateProvider.ValidateUpdate(ctx, cluster, oldCluster)...)
	allErrs = append(allErrs, ValidateClusterScale(cluster.Cluster, oldCluster.Cluster, field.NewPath("spec"))...)

	return allErrs
}

func (p *Provider) PreCreate(ctx context.Context, cluster *types.Cluster) error {
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

	if p.config.AuditEnabled() {
		if !cluster.AuthzWebhookEnabled() {
			cluster.Spec.Features.AuthzWebhookAddr = &platform.AuthzWebhookAddr{Builtin: &platform.
				BuiltinAuthzWebhookAddr{}}
		}
	}

	if p.config.BusinessEnabled() {
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
	if p.config.Feature.SkipConditions != nil {
		cluster.Spec.Features.SkipConditions = append(cluster.Spec.Features.SkipConditions, p.config.Feature.SkipConditions...)
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

// ValidateClusterScale tests if master scale up/down to a cluster is valid.
func ValidateClusterScale(cluster *platform.Cluster, oldCluster *platform.Cluster, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(cluster.Spec.Machines) == len(oldCluster.Spec.Machines) {
		return allErrs
	}
	ha := cluster.Spec.Features.HA
	if ha == nil {
		allErrs = append(allErrs, field.Invalid(fldPath, cluster.Spec.Machines, "HA configuration should enabled for master scale"))
		return allErrs
	}
	if ha.TKEHA == nil && ha.ThirdPartyHA == nil {
		allErrs = append(allErrs, field.Invalid(fldPath, cluster.Spec.Machines, "tkestack HA or third party HA should enabled for master scale"))
		return allErrs
	}
	_, err := PrepareClusterScale(cluster, oldCluster)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, cluster.Spec.Machines, err.Error()))
	}
	return allErrs
}
