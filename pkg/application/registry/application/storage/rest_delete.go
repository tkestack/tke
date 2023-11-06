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

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"tkestack.io/tke/api/application"
	v1 "tkestack.io/tke/api/application/v1"
	applicationversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	applicationstrategy "tkestack.io/tke/pkg/application/registry/application"
)

// CanDeleteREST adapts a service registry into apiserver's RESTStorage model.
type CanDeleteREST struct {
	application       ApplicationStorage
	applicationClient applicationversionedclient.ApplicationV1Interface
	platformClient    platformversionedclient.PlatformV1Interface
}

// NewCanDeleteREST returns a wrapper around the underlying generic storage and performs
// allocations and deallocations of various helm releases related resources like ports.
// TODO: all transactional behavior should be supported from within generic storage
//   or the strategy.
func NewCanDeleteREST(
	application ApplicationStorage,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
) *CanDeleteREST {
	rest := &CanDeleteREST{
		application:       application,
		applicationClient: applicationClient,
		platformClient:    platformClient,
	}
	return rest
}

// New returns an empty object that can be used with Create and Update after request data has been put into it.
func (rs *CanDeleteREST) New() runtime.Object {
	return rs.application.New()
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (rs *CanDeleteREST) ConnectMethods() []string {
	return []string{"GET"}
}

// Get retrieves the object from the storage. It is required to support Patch.
func (rs *CanDeleteREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := rs.application.Get(ctx, name, &metav1.GetOptions{})
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	app := obj.(*application.App)
	appv1 := &v1.App{}
	if err := v1.Convert_application_App_To_v1_App(app, appv1, nil); err != nil {
		return nil, err
	}
	hook := applicationstrategy.GetHooks(appv1)
	result, err := hook.CanDelete(ctx, rs.applicationClient, rs.platformClient, appv1)
	return &result, err
}
