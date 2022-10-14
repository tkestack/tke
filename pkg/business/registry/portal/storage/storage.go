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
	registryUtil "tkestack.io/tke/pkg/business/registry/util"
	"tkestack.io/tke/pkg/util/log"
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
var _ rest.Lister = &REST{}

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

// ConvertToTable converts objects to metav1.Table objects using default table
// convertor.
func (r *REST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*v1.Table, error) {
	tableConvertor := rest.NewDefaultTableConvertor(business.Resource("portals"))
	return tableConvertor.ConvertToTable(ctx, object, tableOptions)
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, _ *metainternal.ListOptions) (runtime.Object, error) {
	log.Debugf("business portal list, ctx %v", ctx)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID == "" || tenantID == "default" {
		return &business.Portal{
			Administrator: true,
			Projects:      make(map[string]string),
		}, nil
	}
	listOpt := v1.ListOptions{FieldSelector: fmt.Sprintf("spec.tenantID=%s", tenantID)}
	projectList, err := r.businessClient.Projects().List(ctx, listOpt)
	if err != nil {
		return nil, err
	}

	log.Debugf("business portal list, before FilterWithUser: %v", projectList)
	isAdmin, projectList, err := registryUtil.FilterWithUser(ctx, projectList, r.authClient, r.businessClient)
	log.Debugf("business portal list, after FilterWithUser: %v, isAdmin '%b'", projectList, isAdmin)
	if err != nil {
		return nil, err
	}

	projects := make(map[string]string)
	extension := make(map[string]business.PortalProject)
	for _, project := range projectList.Items {
		projects[project.ObjectMeta.Name] = project.Spec.DisplayName
		extension[project.ObjectMeta.Name] = business.PortalProject{
			Phase:  (string)(project.Status.Phase),
			Parent: project.Spec.ParentProjectName,
		}
	}

	return &business.Portal{
		Administrator: isAdmin,
		Projects:      projects,
		Extension:     extension,
	}, nil
}
