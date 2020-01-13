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
	"fmt"

	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/business"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	authversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util"
)

// Storage includes storage for portal information and all sub resources.
type Storage struct {
	Portal *REST
}

// NewStorage returns a Storage object that will work against projects.
func NewStorage(_ genericregistry.RESTOptionsGetter, businessClient *businessinternalclient.BusinessClient, authClient authversionedclient.AuthV1Interface) *Storage {
	return &Storage{
		Portal: &REST{
			businessClient: businessClient,
			authClient:     authClient,
		},
	}
}

// REST implements a RESTStorage for user setting.
type REST struct {
	businessClient *businessinternalclient.BusinessClient
	authClient     authversionedclient.AuthV1Interface
}

var _ rest.ShortNamesProvider = &REST{}
var _ rest.Scoper = &REST{}
var _ rest.Storage = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"pt"}
}

// NamespaceScoped returns true if the storage is namespaced.
func (r *REST) NamespaceScoped() bool {
	return false
}

// New returns an empty object that can be used with Create and Update after
// request data has been put into it.
func (r *REST) New() runtime.Object {
	return &business.Portal{}
}

// NewList returns an empty object that can be used with the List call.
func (r *REST) NewList() runtime.Object {
	return &business.Portal{}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	username, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return &business.Portal{
			Administrator: true,
			Projects:      make(map[string]string),
		}, nil
	}
	listOpt := v1.ListOptions{FieldSelector: fmt.Sprintf("spec.tenantID=%s", tenantID)}
	platformList, err := r.businessClient.Platforms().List(listOpt)
	if err != nil {
		return nil, err
	}
	administrator := false
	for _, platform := range platformList.Items {
		if util.InStringSlice(platform.Spec.Administrators, username) {
			administrator = true
			break
		}
	}

	if !administrator && r.authClient != nil {
		idp, err := r.authClient.IdentityProviders().Get(tenantID, v1.GetOptions{})
		if err != nil {
			return nil, err
		}

		if util.InStringSlice(idp.Spec.Administrators, username) {
			administrator = true
		}
	}

	projectList, err := r.businessClient.Projects().List(listOpt)
	if err != nil {
		return nil, err
	}
	projects := make(map[string]string)
	for _, project := range projectList.Items {
		if util.InStringSlice(project.Spec.Members, username) {
			projects[project.ObjectMeta.Name] = project.Spec.DisplayName
		}
	}
	return &business.Portal{
		Administrator: administrator,
		Projects:      projects,
	}, nil
}
