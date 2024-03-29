/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package cluster

import (
	"context"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/tools/clientcmd"
	platformv1client "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/platform/provider/imported/config"
	"tkestack.io/tke/pkg/platform/provider/imported/validation"
	"tkestack.io/tke/pkg/platform/types"
	"tkestack.io/tke/pkg/util/log"
)

func RegisterProvider() {
	p, err := NewProvider()
	if err != nil {
		log.Errorf("init cluster provider error: %s", err)
		return
	}
	clusterprovider.Register(p.Name(), p)
}

type Provider struct {
	*clusterprovider.DelegateProvider
	config *config.Config
}

var _ clusterprovider.Provider = &Provider{}

func NewProvider() (*Provider, error) {
	p := new(Provider)

	p.DelegateProvider = &clusterprovider.DelegateProvider{
		ProviderName: "Imported",
		CreateHandlers: []clusterprovider.Handler{
			p.EnsureCreateClusterMark,
		},
		DeleteHandlers: []clusterprovider.Handler{
			p.EnsureCleanClusterMark,
		},
	}
	cfg, err := config.New(constants.ConfigFile)
	if err != nil {
		return nil, err
	}
	p.config = cfg
	if cfg.PlatformAPIClientConfig != "" {
		restConfig, err := clientcmd.BuildConfigFromFlags("", cfg.PlatformAPIClientConfig)
		if err != nil {
			log.Errorf("read PlatformAPIClientConfig error: %w", err)
		} else {
			p.PlatformClient, err = platformv1client.NewForConfig(restConfig)
			if err != nil {
				return nil, err
			}
		}
	}
	return p, nil
}

func (p *Provider) Validate(cluster *types.Cluster) field.ErrorList {
	return validation.ValidateCluster(context.Background(), cluster)
}
