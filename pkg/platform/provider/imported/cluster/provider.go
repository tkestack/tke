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
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	delegatecluster "tkestack.io/tke/pkg/platform/provider/delegate/cluster"
	"tkestack.io/tke/pkg/platform/provider/imported/validation"
	"tkestack.io/tke/pkg/platform/types"
	"tkestack.io/tke/pkg/util/log"
)

func init() {
	p, err := NewProvider()
	if err != nil {
		log.Errorf("init cluster provider error: %s", err)
		return
	}
	clusterprovider.Register(p.Name(), p)
}

type Provider struct {
	*delegatecluster.DelegateProvider
}

var _ clusterprovider.Provider = &Provider{}

func NewProvider() (*Provider, error) {
	p := new(Provider)

	p.DelegateProvider = &delegatecluster.DelegateProvider{
		ProviderName: "Imported",
		CreateHandlers: []delegatecluster.Handler{
			p.EnsureCreateClusterMark,
		},
		DeleteHandlers: []delegatecluster.Handler{
			p.EnsureCleanClusterMark,
		},
	}
	return p, nil
}

func (p *Provider) Validate(ctx context.Context, cluster *types.Cluster) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, p.DelegateProvider.Validate(ctx, cluster)...)
	allErrs = append(allErrs, validation.ValidateCluster(ctx, cluster)...)

	return allErrs
}
