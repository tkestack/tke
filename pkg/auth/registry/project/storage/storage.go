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
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/request"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/apiserver/filter"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"

	"github.com/casbin/casbin/v2"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/auth"

	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

// Storage includes storage for configmap and all sub resources.
type Storage struct {
	Project *REST
	User    *UserREST
	Group   *GroupREST
	Policy  *PolicyREST

	Binding   *BindingREST
	UnBinding *UnBindingREST
}

// NewStorage returns a Storage object that will work against configmap.
func NewStorage(_ genericregistry.RESTOptionsGetter, authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer) *Storage {
	return &Storage{
		Project:   &REST{authClient: authClient},
		User:      &UserREST{&BindingREST{authClient}, &UnBindingREST{authClient}, authClient},
		Group:     &GroupREST{&BindingREST{authClient}, &UnBindingREST{authClient}, authClient},
		Policy:    &PolicyREST{authClient},
		Binding:   &BindingREST{authClient},
		UnBinding: &UnBindingREST{authClient},
	}
}

// REST implements a RESTStorage for configmap against etcd.
type REST struct {
	rest.Storage

	authClient authinternalclient.AuthInterface
}

func (r *REST) NamespaceScoped() bool {
	return false
}

var _ rest.Scoper = &REST{}
var _ rest.Lister = &REST{}
var _ rest.Getter = &REST{}

// New returns an empty object that can be used with Create and Update after request data has been put into it.
func (r *REST) New() runtime.Object {
	return &auth.Project{}
}

// NewList returns an empty object that can be used with the List call.
func (r *REST) NewList() runtime.Object {
	return &auth.ProjectList{}
}

// ConvertToTable converts objects to metav1.Table objects using default table
// convertor.
func (r *REST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: convert role list to table
	tableConvertor := rest.NewDefaultTableConvertor(auth.Resource("projects"))
	return tableConvertor.ConvertToTable(ctx, object, tableOptions)
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, projectName string, options *metav1.GetOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	projectID := filter.ProjectIDFrom(ctx)
	if projectID == "" {
		projectID = requestInfo.Name
	}

	projectPolicyList, err := r.authClient.ProjectPolicyBindings().List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.projectID=%s", projectID),
	})
	if err != nil {
		log.Error("get project policies failed", log.String("project", projectID), log.Err(err))
		return nil, err
	}

	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" && len(projectPolicyList.Items) > 0 {
		tenantID = projectPolicyList.Items[0].Spec.TenantID
	}

	project := auth.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: projectID,
		},
		TenantID: tenantID,
		Users:    make(map[string]string),
		Groups:   make(map[string]string),
	}

	for _, policy := range projectPolicyList.Items {
		for _, subj := range policy.Spec.Users {
			if _, ok := project.Users[subj.ID]; !ok {
				project.Users[subj.ID] = subj.Name
			}
		}

		for _, subj := range policy.Spec.Groups {
			if _, ok := project.Groups[subj.ID]; !ok {
				project.Groups[subj.ID] = subj.Name
			}
		}
	}

	return &project, nil
}

func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	v1opts := &v1.ListOptions{}
	if tenantID != "" {
		v1opts = util.PredicateV1ListOptions(tenantID, options)
	}

	projectPolicyList, err := r.authClient.ProjectPolicyBindings().List(ctx, *v1opts)
	if err != nil {
		log.Error("list project policies failed", log.Err(err))
		return nil, err
	}

	projectMap := make(map[string]*auth.Project)

	for _, policy := range projectPolicyList.Items {
		project, ok := projectMap[policy.Spec.ProjectID]
		if !ok {
			project = &auth.Project{
				ObjectMeta: metav1.ObjectMeta{
					Name: policy.Spec.ProjectID,
				},
				TenantID: policy.Spec.TenantID,
				Users:    make(map[string]string),
				Groups:   make(map[string]string),
			}
		}

		for _, subj := range policy.Spec.Users {
			if _, ok := project.Users[subj.ID]; !ok {
				project.Users[subj.ID] = subj.Name
			}
		}

		for _, subj := range policy.Spec.Groups {
			if _, ok := project.Groups[subj.ID]; !ok {
				project.Groups[subj.ID] = subj.Name
			}
		}

		projectMap[policy.Spec.ProjectID] = project
	}

	projectList := auth.ProjectList{}
	for _, item := range projectMap {
		projectList.Items = append(projectList.Items, *item)
	}

	return &projectList, nil
}
