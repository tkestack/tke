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
	"github.com/hashicorp/go-plugin"
	"tkestack.io/tke/pkg/app"
	cluster2 "tkestack.io/tke/pkg/platform/provider/baremetal/cluster"
	"tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/log/hclog"
)

const commandDesc = `Baremetal cluster provider for TKE`

// NewApp creates a App object with default parameters.
func NewApp(basename string) *app.App {
	application := app.NewApp("Baremetal Cluster Provider",
		basename,
		app.WithSilence(),
		app.WithDescription(commandDesc),
		app.WithRunFunc(run()),
	)
	return application
}

func run() app.RunFunc {
	return func(basename string) error {
		serveConfig := cluster.NewServeConfig(new(cluster2.Provider), hclog.NewHCLogger(log.ZapLogger()))
		plugin.Serve(serveConfig)
		return nil
	}
}
