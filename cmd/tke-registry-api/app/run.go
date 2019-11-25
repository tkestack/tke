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

package app

import (
	"tkestack.io/tke/cmd/tke-registry-api/app/config"
	"tkestack.io/tke/pkg/util/log"
)

// Run runs the specified TKE notify apiserver. This should never exit.
func Run(cfg *config.Config, stopCh <-chan struct{}) error {
	log.Info("Starting Tencent Kubernetes Engine Registry API server")

	server, err := CreateServerChain(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run(stopCh)
}
