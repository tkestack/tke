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

package storage

import (
	"context"

	"github.com/casbin/casbin/v2"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"

	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for configmap and all sub resources.
type Storage struct {
	User   *REST
	Policy *PolicyREST
	Role   *RoleREST
}

// NewStorage returns a Storage object that will work against configmap.
func NewStorage(_ genericregistry.RESTOptionsGetter, authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer) *Storage {
	return &Storage{
		User:   &REST{},
		Policy: &PolicyREST{&REST{}, authClient, enforcer},
		Role:   &RoleREST{&REST{}, authClient, enforcer},
	}
}

// REST implements a RESTStorage for configmap against etcd.
type REST struct {
	rest.Storage
}

func (r *REST) NamespaceScoped() bool {
	return false
}

var _ rest.ShortNamesProvider = &REST{}
var _ rest.Creater = &REST{}
var _ rest.Scoper = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"user"}
}

// New returns an empty object that can be used with Create and Update after request data has been put into it.
func (r *REST) New() runtime.Object {
	return &auth.User{}
}

// NewList returns an empty object that can be used with the List call.
func (r *REST) NewList() runtime.Object {
	return &auth.UserList{}
}

// Create creates a new version of a resource.
func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	return &auth.User{}, nil
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		tenantID, name = util.ParseTenantAndName(name)
	}

	idp, ok := identityprovider.IdentityProvidersStore[tenantID]
	if !ok {
		log.Error("Tenant has no related identity providers", log.String("tenantID", tenantID))
		return nil, apierrors.NewNotFound(auth.Resource("user"), name)
	}

	userLister, ok := idp.(identityprovider.UserGetter)
	if !ok {
		log.Info("tenant %s related idp not implement UserLister interface", log.String("tenantID", tenantID))
		return nil, apierrors.NewNotFound(auth.Resource("user"), name)
	}

	return userLister.GetUser(ctx, name, options)
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		tenantID, _ = options.FieldSelector.RequiresExactMatch("spec.tenantID")
		if tenantID == "" {
			return nil, apierrors.NewBadRequest("List groups must specify tenantID")
		}
	}

	idp, ok := identityprovider.IdentityProvidersStore[tenantID]
	if !ok {
		log.Error("Tenant has no related identity providers", log.String("tenantID", tenantID))
		return &auth.UserList{}, nil
	}

	userLister, ok := idp.(identityprovider.UserLister)
	if !ok {
		log.Info("tenant %s related idp not implement UserLister interface", log.String("tenantID", tenantID))
		return &auth.UserList{}, nil
	}

	users, err := userLister.ListUsers(ctx, options)
	return users, err
}
