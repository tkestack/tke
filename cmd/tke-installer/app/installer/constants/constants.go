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
	DataDir        = "data/"
	ClusterFile    = DataDir + "tke.json"
	ClusterLogFile = DataDir + "tke.log"
	ChartDirName   = "manifests/charts/"

	DefaultCustomResourceDir  = DataDir + "custom_upgrade_resource"
	DefaultCustomChartsName   = "custom.charts.tar.gz"
	CustomK8sImageDirName     = "images/"
	CustomK8sBinaryDirName    = "bins/"
	CustomK8sBinaryAmdDirName = "bins/linux-amd64/"
	CustomK8sBinaryArmDirName = "bins/linux-arm64/"

	ProviderConfigFile = "provider/baremetal/conf/config.yaml"

	HooksDir             = "hooks/"
	PreInstallHook       = HooksDir + "pre-install"
	PostClusterReadyHook = HooksDir + "post-cluster-ready"
	PostInstallHook      = HooksDir + "post-install"

	DockerCertsDir = "/etc/docker/certs.d"

	DefaultTeantID                 = "default"
	DefaultChartGroupName          = "public"
	DefaultExpansionChartGroupName = "expansion"
	DefaultCustomChartGroupName    = "custom"
	GlobalClusterName              = "global"

	DevRegistryDomain    = "docker.io"
	DevRegistryNamespace = "tkestack"
	ImagesFile           = "images.tar.gz"
	ImagesPattern        = DevRegistryNamespace + "/*"

	OIDCClientSecretFile = DataDir + "oidc_client_secret"
	CACrtFile            = DataDir + "ca.crt"
	CAKeyFile            = DataDir + "ca.key"
	FrontProxyCACrtFile  = DataDir + "front-proxy-ca.crt"
	ServerCrtFile        = DataDir + "server.crt"
	ServerKeyFile        = DataDir + "server.key"
	AdminCrtFile         = DataDir + "admin.crt"
	AdminKeyFile         = DataDir + "admin.key"
	WebhookCrtFile       = DataDir + "webhook.crt"
	WebhookKeyFile       = DataDir + "webhook.key"
	KubeconfigFile       = DataDir + "admin.kubeconfig"

	CACrtFileBaseName      = "ca.crt"
	CAKeyFileBaseName      = "ca.key"
	ServerCrtFileBaseName  = "server.crt"
	ServerKeyFileBaseName  = "server.key"
	AdminCrtFileBaseName   = "admin.crt"
	AdminKeyFileBaseName   = "admin.key"
	WebhookCrtFileBaseName = "webhook.crt"
	WebhookKeyFileBaseName = "webhook.key"
	KubeconfigFileBaseName = "admin.kubeconfig"

	AuthzWebhookNodePort = 31138

	DefaultApplicationInstallDriverType = "HelmV3"
	DefaultApplicationInstallValueType  = "yaml"
)

const (
	RegistryHTTPOptions = `-d \
--name registry-http \
--restart always \
-p 80:5000 \
-v /opt/tke-installer/registry:/var/lib/registry`

	RegistryHTTPSOptions = `-d \
--name registry-https \
--restart always \
-p 443:443 \
-v /opt/tke-installer/registry:/var/lib/registry \
-v registry-certs:/certs \
-e REGISTRY_HTTP_ADDR=0.0.0.0:443 \
-e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/server.crt \
-e REGISTRY_HTTP_TLS_KEY=/certs/server.key`
)

const (
	CPURequest                = 2
	MemoryRequest             = 4  // GiB
	FirstNodeDiskSpaceRequest = 30 // GiB

	PathForDiskSpaceRequest = "/var/lib"
)

const (
	CephRBDStorageClassName = "csi-rbd-sc"
	CephRBDChartReleaseName = "ceph-csi-rbd"
	CephFSStorageClassName  = "csi-cephfs-sc"
	CephFSChartReleaseName  = "ceph-csi-cephfs"
	NFSStorageClassName     = "nfs-sc"
	NFSChartReleaseName     = "nfs-subdir-external-provisioner"
)
