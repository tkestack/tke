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

package route

import (
	dexserver "github.com/dexidp/dex/server"
	"k8s.io/apiserver/pkg/server/mux"
	"tkestack.io/tke/pkg/auth/types"
)

// RegisterOIDCRoute to install route for OIDC http handler.
func RegisterOIDCRoute(m *mux.PathRecorderMux, server *dexserver.Server) {
	m.HandlePrefix("/"+types.IssuerName+"/", server)
}
