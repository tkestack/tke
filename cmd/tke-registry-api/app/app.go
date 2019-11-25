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
	commonapiserver "k8s.io/apiserver/pkg/server"
	"tkestack.io/tke/cmd/tke-registry-api/app/config"
	"tkestack.io/tke/cmd/tke-registry-api/app/options"
	"tkestack.io/tke/pkg/app"
	"tkestack.io/tke/pkg/util/log"
)

const commandDesc = `The Registry API server validates and configures data for the api objects which 
include namespace and images. The API Server services REST operations and 
provides the frontend to the registry's shared state through which all other 
components interact. At the same time, the registry also provides a container 
image warehouse service implemented in accordance with the Docker Registry V2 
protocol.`

// NewApp creates a App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.NewOptions(basename)
	application := app.NewApp("Tencent Kubernetes Engine Registry API Server", basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithRunFunc(run(opts)),
	)
	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()

		if err := opts.Complete(); err != nil {
			return err
		}

		cfg, err := config.CreateConfigFromOptions(basename, opts)
		if err != nil {
			return err
		}

		stopCh := commonapiserver.SetupSignalHandler()
		return Run(cfg, stopCh)
	}
}
