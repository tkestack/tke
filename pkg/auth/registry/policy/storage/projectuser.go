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

	"k8s.io/apiserver/pkg/registry/rest"

	"tkestack.io/tke/pkg/apiserver/filter"
	"tkestack.io/tke/pkg/auth/util"

	"k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic/registry"

	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/util/log"
)

// ProjectUserREST implements the REST endpoint.
type ProjectUserREST struct {
	policyStore *registry.Store

	authClient authinternalclient.AuthInterface
}

var _ rest.Lister = &ProjectUserREST{}

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *ProjectUserREST) New() runtime.Object {
	return &auth.ProjectPolicyBinding{}
}

// NewList returns an empty object that can be used with the List call.
func (r *ProjectUserREST) NewList() runtime.Object {
	return &auth.UserList{}
}

// ConvertToTable converts objects to metav1.Table objects using default table
// convertor.
func (r *ProjectUserREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: convert role list to table
	tableConvertor := rest.NewDefaultTableConvertor(auth.Resource("projectusers"))
	return tableConvertor.ConvertToTable(ctx, object, tableOptions)
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *ProjectUserREST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	policyID := requestInfo.Name
	polObj, err := r.policyStore.Get(ctx, requestInfo.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	policy := polObj.(*auth.Policy)
	if policy.Spec.Scope != auth.PolicyProject {
		return nil, errors.NewBadRequest("unable bind subject to platform-scoped policy, please use binding api")
	}

	// TODO if projectID is empty, list all users fro related project
	projectID := filter.ProjectIDFrom(ctx)
	if projectID == "" {
		return nil, errors.NewBadRequest("must specify projectID")
	}

	proBinding, err := r.authClient.ProjectPolicyBindings().Get(ctx, util.ProjectPolicyName(projectID, policyID), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	userList := &auth.UserList{}
	for _, subj := range proBinding.Spec.Users {
		var user *auth.User
		if subj.ID != "" {
			user, err = r.authClient.Users().Get(ctx, util.CombineTenantAndName(policy.Spec.TenantID, subj.ID), metav1.GetOptions{})
			if err != nil {
				log.Error("Get user failed", log.String("id", subj.ID), log.Err(err))
				user = constructUser(subj.ID, subj.Name)
			}
		} else {
			user = constructUser(subj.ID, subj.Name)
		}

		userList.Items = append(userList.Items, *user)
	}

	return userList, nil
}
