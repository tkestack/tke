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

package constants

import (
	"time"

	platformv1 "tkestack.io/tke/api/platform/v1"
)

const (
	AuditPolicyConfigName  = "audit-policy.yaml"
	AuthzWebhookConfigName = "tke-authz-webhook.yaml"
	OIDCCACertName         = "oidc-ca.crt"
	AdminCertName          = "admin.crt"
	AdminKeyName           = "admin.key"
	WebhookCertName        = "webhook.crt"
	WebhookKeyName         = "webhook.key"
	// Kubernetes Config
	KubernetesDir                       = "/etc/kubernetes/"
	KubernetesSchedulerPolicyConfigFile = KubernetesDir + "scheduler-policy-config.json"
	KubernetesAuditWebhookConfigFile    = KubernetesDir + "audit-api-client-config.yaml"
	TokenFile                           = KubernetesDir + "known_tokens.csv"
	KubernetesAuditPolicyConfigFile     = KubernetesDir + AuditPolicyConfigName
	KubernetesAuthzWebhookConfigFile    = KubernetesDir + AuthzWebhookConfigName
	KubeadmConfigFileName               = KubernetesDir + "kubeadm-config.yaml"
	KubeletKubeConfigFileName           = KubernetesDir + "kubelet.conf"

	KubeletPodManifestDir                = KubernetesDir + "manifests/"
	EtcdPodManifestFile                  = KubeletPodManifestDir + "etcd.yaml"
	KubeAPIServerPodManifestFile         = KubeletPodManifestDir + "kube-apiserver.yaml"
	KubeControllerManagerPodManifestFile = KubeletPodManifestDir + "kube-controller-manager.yaml"
	KubeSchedulerPodManifestFile         = KubeletPodManifestDir + "kube-scheduler.yaml"
	KeepavlivedManifestFile              = KubeletPodManifestDir + "keepalived.yaml"

	KubeadmPathInNodePackge = "kubernetes/node/bin/kubeadm"
	KubeletPathInNodePackge = "kubernetes/node/bin/kubelet"
	KubectlPathInNodePackge = "kubernetes/node/bin/kubectl"

	DstTmpDir  = "/tmp/k8s/"
	DstBinDir  = "/usr/bin/"
	CNIBinDir  = "/opt/cni/bin/"
	CNIDataDir = "/var/lib/cni/"
	CNIConfDIr = "/etc/cni/net.d"
	AppCertDir = "/app/certs/"

	// AppCert
	AppAdminCertFile = AppCertDir + AdminCertName
	AppAdminKeyFile  = AppCertDir + AdminKeyName

	// ETC
	EtcdDataDir          = "/var/lib/etcd"
	KubectlConfigFile    = "/root/.kube/config"
	KeepavliedConfigFile = "/etc/keepalived/keepalived.conf"

	// PKI
	CertificatesDir = KubernetesDir + "pki/"
	OIDCCACertFile  = CertificatesDir + OIDCCACertName
	WebhookCertFile = CertificatesDir + WebhookCertName
	WebhookKeyFile  = CertificatesDir + WebhookKeyName
	AdminCertFile   = CertificatesDir + AdminCertName
	AdminKeyFile    = CertificatesDir + AdminKeyName

	// CACertName defines certificate name
	CACertName = CertificatesDir + "ca.crt"
	// CAKeyName defines certificate name
	CAKeyName = CertificatesDir + "ca.key"
	// APIServerCertName defines API's server certificate name
	APIServerCertName = CertificatesDir + "apiserver.crt"
	// APIServerKeyName defines API's server key name
	APIServerKeyName = CertificatesDir + "apiserver.key"
	// KubeletClientCurrent defines kubelet rotate certificates
	KubeletClientCurrent = "/var/lib/kubelet/pki/kubelet-client-current.pem"
	// EtcdCACertName defines etcd's CA certificate name
	EtcdCACertName = CertificatesDir + "etcd/ca.crt"
	// EtcdCAKeyName defines etcd's CA key name
	EtcdCAKeyName = CertificatesDir + "etcd/ca.key"
	// EtcdListenClientPort defines the port etcd listen on for client traffic
	EtcdListenClientPort = 2379
	// EtcdListenPeerPort defines the port etcd listen on for peer traffic
	EtcdListenPeerPort = 2380
	// APIServerEtcdClientCertName defines apiserver's etcd client certificate name
	APIServerEtcdClientCertName = CertificatesDir + "apiserver-etcd-client.crt"
	// APIServerEtcdClientKeyName defines apiserver's etcd client key name
	APIServerEtcdClientKeyName = CertificatesDir + "apiserver-etcd-client.key"

	// LabelNodeRoleMaster specifies that a node is a control-plane
	// This is a duplicate definition of the constant in pkg/controller/service/service_controller.go
	LabelNodeRoleMaster = "node-role.kubernetes.io/master"

	// LabelNodeNeedUpgrade specifies that a node need upgrade.
	LabelNodeNeedUpgrade = platformv1.GroupName + "/need-upgrade"

	// Provider
	ProviderDir           = "provider/baremetal/"
	SrcDir                = ProviderDir + "res/"
	ConfDir               = ProviderDir + "conf/"
	ConfigFile            = ConfDir + "config.yaml"
	AuditPolicyConfigFile = ConfDir + AuditPolicyConfigName
	OIDCConfigFile        = ConfDir + OIDCCACertName
	ManifestsDir          = ProviderDir + "manifests/"
	GPUManagerManifest    = SrcDir + "gpu-manager/gpu-manager.yaml"
	CSIOperatorManifest   = ManifestsDir + "csi-operator/csi-operator.yaml"
	MetricsServerManifest = ManifestsDir + "metrics-server/metrics-server.yaml"
	CiliumManifest        = SrcDir + "cilium/*.yaml"

	KUBERNETES                   = 1
	DNSIPIndex                   = 10
	GPUQuotaAdmissionIPIndex     = 9
	GPUQuotaAdmissionIPAnnotaion = platformv1.GroupName + "/gpu-quota-admission-ip"

	// RenewCertsTimeThreshold control how long time left to renew certs
	RenewCertsTimeThreshold = 30 * 24 * time.Hour

	// MinNumCPU mininum cpu number.
	MinNumCPU = 2

	APIServerHostName = "api.tke.com"

	// include itself
	NeedUpgradeCoreDNSLowerK8sVersion = "1.19.0"
	// not include itself
	NeedUpgradeCoreDNSUpperK8sVersion = "1.21.0"
)
