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

package webtty

import (
	"k8s.io/apiserver/pkg/server/mux"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
	"tkestack.io/tke/pkg/util/log"
)

// Path is the URL path for webtty.
const Path = "/webtty"

// RegisterRoute is used to register prefix path routing matches for all
// configured backend components.
func RegisterRoute(m *mux.PathRecorderMux, cfg *gatewayconfig.GatewayConfiguration) error {
	if cfg.Components.Platform == nil {
		log.Warn("WebTTY disabled because no platform component registered")
		return nil
	}
	address := cfg.Components.Platform.Address
	if address == "" {
		log.Warn("WebTTY disabled because platform component no address")
		return nil
	}
	handler, err := NewHandler(address)
	if err != nil {
		return err
	}
	m.Handle(Path, handler)
	return nil
}
