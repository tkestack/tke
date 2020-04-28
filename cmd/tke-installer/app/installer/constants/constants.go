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

	ProviderConfigFile = "provider/baremetal/conf/config.yaml"

	HooksDir             = "hooks/"
	PreInstallHook       = HooksDir + "pre-install"
	PostClusterReadyHook = HooksDir + "post-cluster-ready"
	PostInstallHook      = HooksDir + "post-install"

	DockerCertsDir = "/etc/docker/certs.d"

	DefaultTeantID = "default"

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

	KubeconfigFile = DataDir + "admin.kubeconfig"
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
	CPURequest    = 8
	MemoryRequest = 15 // GiB
)
