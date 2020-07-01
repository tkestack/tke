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
	"tkestack.io/tke/pkg/auth/util"

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
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for configmap and all sub resources.
type Storage struct {
	Group  *REST
	Policy *PolicyREST
	Role   *RoleREST
}

// NewStorage returns a Storage object that will work against configmap.
func NewStorage(_ genericregistry.RESTOptionsGetter, authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer) *Storage {
	return &Storage{
		Group:  &REST{},
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
var _ rest.Lister = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"group"}
}

// New returns an empty object that can be used with Create and Update after request data has been put into it.
func (r *REST) New() runtime.Object {
	return &auth.Group{}
}

// NewList returns an empty object that can be used with the List call.
func (r *REST) NewList() runtime.Object {
	return &auth.GroupList{}
}

// ConvertToTable converts objects to metav1.Table objects using default table
// convertor.
func (r *REST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: convert role list to table
	tableConvertor := rest.NewDefaultTableConvertor(auth.Resource("groups"))
	return tableConvertor.ConvertToTable(ctx, object, tableOptions)
}

// Create creates a new version of a resource.
func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	return &auth.Group{}, nil
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID == "" {
		tenantID, name = util.ParseTenantAndName(name)
	}

	idp, ok := identityprovider.GetIdentityProvider(tenantID)
	if !ok {
		log.Error("Tenant has no related identity providers", log.String("tenantID", tenantID))
		return nil, apierrors.NewNotFound(auth.Resource("group"), name)
	}

	groupGetter, ok := idp.(identityprovider.GroupGetter)
	if !ok {
		log.Error("tenant %s related identity providers not implement UserLister interface", log.String("tenantID", tenantID))
		return nil, apierrors.NewNotFound(auth.Resource("group"), name)
	}

	return groupGetter.GetGroup(ctx, name, options)
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)

	if tenantID == "" {
		tenantID, _ = options.FieldSelector.RequiresExactMatch("spec.tenantID")
		if tenantID == "" {
			return &auth.GroupList{}, nil
		}
	}
	idp, ok := identityprovider.GetIdentityProvider(tenantID)
	if !ok {
		log.Error("Tenant has no related identity providers", log.String("tenantID", tenantID))
		return &auth.GroupList{}, nil
	}

	groupLister, ok := idp.(identityprovider.GroupLister)
	if !ok {
		log.Info("tenant %s related idp not implement UserLister interface", log.String("tenantID", tenantID))
		return &auth.GroupList{}, nil
	}

	users, err := groupLister.ListGroups(ctx, options)
	return users, err
}
