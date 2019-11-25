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
	// TokenFile for store tke token auth
	TokenFile = "data/tke_token"
	// OIDCClientSecretFile for store oidc client secret
	OIDCClientSecretFile = "data/oidc_client_secret"
	// OIDCAdminPasswordFile for store
	OIDCAdminPasswordFile = "data/admin_password"
	// RegistryPasswordFile for registry
	RegistryPasswordFile = "data/registry_password"

	// CACrtFile for store
	CACrtFile = "data/ca.crt"
	// CAKeyFile for store
	CAKeyFile = "data/ca.key"
	// ServerCrtFile for store
	ServerCrtFile = "data/server.crt"
	// ServerKeyFile for store
	ServerKeyFile = "data/server.key"
	// AdminCrtFile for store
	AdminCrtFile = "data/admin.crt"
	// AdminKeyFile for store
	AdminKeyFile = "data/admin.key"

	KubeconfigFile = "data/admin.kubeconfig"
)
