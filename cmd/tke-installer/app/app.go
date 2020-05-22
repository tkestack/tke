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
	"tkestack.io/tke/cmd/tke-installer/app/config"
	"tkestack.io/tke/cmd/tke-installer/app/options"
	"tkestack.io/tke/pkg/app"
	"tkestack.io/tke/pkg/util/log"
)

const commandDesc = `The TKE Installer is used to setup the first kubernetes cluster.`

// NewApp creates a App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.NewOptions(basename)
	application := app.NewApp("Tencent Kubernetes Engine Installer",
		basename,
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

		return Run(cfg)
	}
}
