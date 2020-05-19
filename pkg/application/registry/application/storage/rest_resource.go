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
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"tkestack.io/tke/api/application"
	applicationinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/application/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	registryversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	"tkestack.io/tke/pkg/application/util"
)

// ResourceREST adapts a service registry into apiserver's RESTStorage model.
type ResourceREST struct {
	application       ApplicationStorage
	applicationClient *applicationinternalclient.ApplicationClient
	platformClient    platformversionedclient.PlatformV1Interface
	registryClient    registryversionedclient.RegistryV1Interface
}

// NewResourceREST returns a wrapper around the underlying generic storage and performs
// allocations and deallocations of various helm releases related resources like ports.
// TODO: all transactional behavior should be supported from within generic storage
//   or the strategy.
func NewResourceREST(
	application ApplicationStorage,
	applicationClient *applicationinternalclient.ApplicationClient,
	platformClient platformversionedclient.PlatformV1Interface,
	registryClient registryversionedclient.RegistryV1Interface,
) *ResourceREST {
	rest := &ResourceREST{
		application:       application,
		applicationClient: applicationClient,
		platformClient:    platformClient,
		registryClient:    registryClient,
	}
	return rest
}

// New returns an empty object that can be used with Create and Update after request data has been put into it.
func (rs *ResourceREST) New() runtime.Object {
	return rs.application.New()
}

// Get retrieves the object from the storage. It is required to support Patch.
func (rs *ResourceREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := rs.application.Get(ctx, name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	app := obj.(*application.App)

	client, err := util.NewHelmClient(ctx, rs.platformClient, app.Spec.TargetCluster, app.Namespace)
	if err != nil {
		return nil, err
	}
	resources, err := client.Status(&helmaction.StatusOptions{
		Namespace:   app.Namespace,
		ReleaseName: app.Spec.Name,
	})
	if err != nil {
		return nil, err
	}
	appResource := &application.AppResource{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: app.Namespace,
			Name:      app.Name,
		},
		Spec: application.AppResourceSpec{
			Type:          app.Spec.Type,
			TenantID:      app.Spec.TenantID,
			Name:          app.Spec.Name,
			TargetCluster: app.Spec.TargetCluster,
		},
	}
	res := make(application.Resources, len(resources.Resources))
	for k, v := range resources.Resources {
		res[k] = v[:]
	}
	appResource.Spec.Resources = res
	return appResource, nil
}
