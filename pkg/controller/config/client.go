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
func BuildClientConfig(opts *options.APIServerClientOptions) (*restclient.Config, error) {
	if opts.Server == "" && opts.ServerClientConfig == "" {
		return nil, fmt.Errorf("either --api-server or --api-server-client-config should be specified")
	}
	apiServerClientConfig, err := clientcmd.BuildConfigFromFlags(opts.Server, opts.ServerClientConfig)
	if err != nil {
		return nil, err
	}
	apiServerClientConfig.ContentConfig.ContentType = opts.ContentType
	apiServerClientConfig.QPS = opts.QPS
	apiServerClientConfig.Burst = int(opts.Burst)
	return apiServerClientConfig, nil
}
