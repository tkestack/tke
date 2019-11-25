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

package storage

import (
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/util"
)

// Storage includes storage for resources.
type Storage struct {
	NetworkPolicy *REST
}

// REST implements pkg/api/rest.StandardStorage.
type REST struct {
	*util.Store
}

// NewStorageV1 returns a Storage object that will work against resources.
func NewStorageV1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	networkPolicyStore := &util.Store{
		NewFunc:        func() runtime.Object { return &networkingV1.NetworkPolicy{} },
		NewListFunc:    func() runtime.Object { return &networkingV1.NetworkPolicyList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	return &Storage{
		NetworkPolicy: &REST{networkPolicyStore},
	}
}

// NewStorageExtensionsV1Beta1 returns a Storage object that will work against resources.
func NewStorageExtensionsV1Beta1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	networkPolicyStore := &util.Store{
		NewFunc:        func() runtime.Object { return &extensionsv1beta1.NetworkPolicy{} },
		NewListFunc:    func() runtime.Object { return &extensionsv1beta1.NetworkPolicyList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	return &Storage{
		NetworkPolicy: &REST{networkPolicyStore},
	}
}
