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
	storageV1 "k8s.io/api/storage/v1"
	storageV1Beta1 "k8s.io/api/storage/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/proxy"
)

// Storage includes storage for resources.
type Storage struct {
	StorageClass *REST
	Events       *EventREST
}

// REST implements pkg/api/rest.StandardStorage.
type REST struct {
	*proxy.Store
}

// NewStorageV1Beta1 returns a Storage object that will work against resources.
func NewStorageV1Beta1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	storageClassStore := &proxy.Store{
		NewFunc:        func() runtime.Object { return &storageV1Beta1.StorageClass{} },
		NewListFunc:    func() runtime.Object { return &storageV1Beta1.StorageClassList{} },
		Namespaced:     false,
		PlatformClient: platformClient,
	}

	return &Storage{
		StorageClass: &REST{storageClassStore},
		Events: &EventREST{
			platformClient: platformClient,
		},
	}
}

// NewStorageV1 returns a Storage object that will work against resources.
func NewStorageV1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	storageClassStore := &proxy.Store{
		NewFunc:        func() runtime.Object { return &storageV1.StorageClass{} },
		NewListFunc:    func() runtime.Object { return &storageV1.StorageClassList{} },
		Namespaced:     false,
		PlatformClient: platformClient,
	}

	return &Storage{
		StorageClass: &REST{storageClassStore},
		Events: &EventREST{
			platformClient: platformClient,
		},
	}
}

var _ rest.ShortNamesProvider = &REST{}

// ShortNames is used by cli to have short names representation of resources.
func (r *REST) ShortNames() []string {
	return []string{"sc"}
}
