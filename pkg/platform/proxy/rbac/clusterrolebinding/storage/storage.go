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
	rbacV1 "k8s.io/api/rbac/v1"
	rbacV1Alpha1 "k8s.io/api/rbac/v1alpha1"
	rbacV1Beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/util"
)

// Storage includes storage for resources.
type Storage struct {
	RuntimeClass *REST
}

// REST implements pkg/api/rest.StandardStorage.
type REST struct {
	*util.Store
}

// NewStorageV1Alpha1 returns a Storage object that will work against resources.
func NewStorageV1Alpha1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	rbacStore := &util.Store{
		NewFunc:        func() runtime.Object { return &rbacV1Alpha1.ClusterRoleBinding{} },
		NewListFunc:    func() runtime.Object { return &rbacV1Alpha1.ClusterRoleBindingList{} },
		Namespaced:     false,
		PlatformClient: platformClient,
	}

	return &Storage{
		RuntimeClass: &REST{rbacStore},
	}
}

// NewStorageV1Beta1 returns a Storage object that will work against resources.
func NewStorageV1Beta1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	rbacStore := &util.Store{
		NewFunc:        func() runtime.Object { return &rbacV1Beta1.ClusterRoleBinding{} },
		NewListFunc:    func() runtime.Object { return &rbacV1Beta1.ClusterRoleBindingList{} },
		Namespaced:     false,
		PlatformClient: platformClient,
	}

	return &Storage{
		RuntimeClass: &REST{rbacStore},
	}
}

// NewStorageV1 returns a Storage object that will work against resources.
func NewStorageV1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	rbacStore := &util.Store{
		NewFunc:        func() runtime.Object { return &rbacV1.ClusterRoleBinding{} },
		NewListFunc:    func() runtime.Object { return &rbacV1.ClusterRoleBindingList{} },
		Namespaced:     false,
		PlatformClient: platformClient,
	}

	return &Storage{
		RuntimeClass: &REST{rbacStore},
	}
}

var _ rest.ShortNamesProvider = &REST{}

// ShortNames returns short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{""}
}
