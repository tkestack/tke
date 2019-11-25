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
	"tkestack.io/tke/cmd/tke-notify-controller/app/config"
	"tkestack.io/tke/cmd/tke-notify-controller/app/options"
	"tkestack.io/tke/pkg/app"
	"tkestack.io/tke/pkg/util/log"
)

const commandDesc = `The notify controller manager is a daemon that embeds the core control loops. In 
applications of robotics and automation, a control loop is a non-terminating 
loop that regulates the state of the system. In TKE, a controller is a control 
loop that watches the shared state of the project and cluster through the 
apiserver and makes changes attempting to move the current state towards the 
desired state.`

// NewApp creates a App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.NewOptions(basename, KnownControllers(), ControllersDisabledByDefault.List())
	application := app.NewApp("Tencent Kubernetes Engine Notify Controller Manager",
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

		cfg, err := config.CreateConfigFromOptions(basename, opts)
		if err != nil {
			return err
		}

		stopCh := commonapiserver.SetupSignalHandler()
		return Run(cfg, stopCh)
	}
}
