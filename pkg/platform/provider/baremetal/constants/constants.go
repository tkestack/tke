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

const (
	// KubernetesDir is the directory Kubernetes owns for storing various configuration files
	KubernetesDir             = "/etc/kubernetes/"
	KubeletPodManifestDir     = KubernetesDir + "manifests/"
	SchedulerPolicyConfigFile = KubernetesDir + "scheduler-policy-config.json"

	EtcdPodManifestFile                  = KubeletPodManifestDir + "etcd.yaml"
	KubeAPIServerPodManifestFile         = KubeletPodManifestDir + "kube-apiserver.yaml"
	KubeControllerManagerPodManifestFile = KubeletPodManifestDir + "kube-controller-manager.yaml"
	KubeSchedulerPodManifestFile         = KubeletPodManifestDir + "kube-scheduler.yaml"

	DstTmpDir  = "/tmp/k8s/"
	DstBinDir  = "/usr/bin/"
	CNIBinDir  = "/opt/cni/bin/"
	CNIDataDir = "/var/lib/cni/"
	CNIConfDIr = "/etc/cni"

	CertificatesDir = KubernetesDir + "pki/"
	EtcdDataDir     = "/var/lib/etcd"

	TokenFile = KubernetesDir + "known_tokens.csv"

	KubectlConfigFile = "/root/.kube/config"

	OIDCCACertName = "oidc-ca.crt"
	OIDCCACertFile = CertificatesDir + OIDCCACertName

	// CACertName defines certificate name
	CACertName = CertificatesDir + "ca.crt"
	// EtcdCACertName defines etcd's CA certificate name
	EtcdCACertName = CertificatesDir + "etcd/ca.crt"
	// EtcdListenClientPort defines the port etcd listen on for client traffic
	EtcdListenClientPort = 2379
	// EtcdListenPeerPort defines the port etcd listen on for peer traffic
	EtcdListenPeerPort = 2380
	// APIServerEtcdClientCertName defines apiserver's etcd client certificate name
	APIServerEtcdClientCertName = CertificatesDir + "apiserver-etcd-client.crt"
	// APIServerEtcdClientKeyName defines apiserver's etcd client key name
	APIServerEtcdClientKeyName = CertificatesDir + "apiserver-etcd-client.key"

	KubeadmConfigFileName = KubernetesDir + "kubeadm-config.yaml"

	// LabelNodeRoleMaster specifies that a node is a control-plane
	// This is a duplicate definition of the constant in pkg/controller/service/service_controller.go
	LabelNodeRoleMaster = "node-role.kubernetes.io/master"

	ProviderDir  = "provider/baremetal/"
	SrcDir       = ProviderDir + "res/"
	ConfDir      = ProviderDir + "conf/"
	ManifestsDir = ProviderDir + "manifests/"
	ConfigFile   = ConfDir + "config.yaml"
)
