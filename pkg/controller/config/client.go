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

package config

import (
	"fmt"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"tkestack.io/tke/pkg/controller/options"
)

// BuildClientConfig to build the rest config by given options.
func BuildClientConfig(opts *options.APIServerClientOptions) (cfg *restclient.Config, ok bool, err error) {
	if opts.Required {
		if opts.Server == "" && opts.ServerClientConfig == "" {
			err = fmt.Errorf("either %s or %s should be specified",
				options.FlagAPIClientServer(opts.Name),
				options.FlagAPIClientServerClientConfig(opts.Name))
			return
		}
	}
	if opts.Server == "" && opts.ServerClientConfig == "" {
		return
	}

	cfg, err = clientcmd.BuildConfigFromFlags(opts.Server, opts.ServerClientConfig)
	if err != nil {
		return
	}
	cfg.ContentConfig.ContentType = opts.ContentType
	cfg.QPS = opts.QPS
	cfg.Burst = int(opts.Burst)
	ok = true
	return
}
