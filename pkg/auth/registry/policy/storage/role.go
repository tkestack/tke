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
	"strings"
	"tkestack.io/tke/pkg/apiserver/filter"

	"k8s.io/apiserver/pkg/registry/generic/registry"

	"github.com/casbin/casbin/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/endpoints/request"
	"tkestack.io/tke/pkg/util/log"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

// RoleREST implements the REST endpoint, list policies bound to the user.
type RoleREST struct {
	policyStore *registry.Store
	authClient  authinternalclient.AuthInterface
	enforcer    *casbin.SyncedEnforcer
}

var _ = rest.Lister(&RoleREST{})

// NewList returns an empty object that can be used with the List call.
func (r *RoleREST) NewList() runtime.Object {
	return &auth.RoleList{}
}

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *RoleREST) New() runtime.Object {
	return &auth.Role{}
}

// ConvertToTable converts objects to metav1.Table objects using default table
// convertor.
func (r *RoleREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: convert role list to table
	tableConvertor := rest.NewDefaultTableConvertor(auth.Resource("roles"))
	return tableConvertor.ConvertToTable(ctx, object, tableOptions)
}

func (r *RoleREST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	obj, err := r.policyStore.Get(ctx, requestInfo.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	policy := obj.(*auth.Policy)

	projectID := filter.ProjectIDFrom(ctx)
	roles := r.enforcer.GetUsersForRoleInDomain(policy.Name, projectID)
	var roleIDs []string
	for _, r := range roles {
		if strings.HasPrefix(r, "rol-") {
			roleIDs = append(roleIDs, r)
		}
	}

	var roleList = &auth.RoleList{}
	for _, id := range roleIDs {
		rol, err := r.authClient.Roles().Get(ctx, id, metav1.GetOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			log.Error("Get role failed", log.String("role", id), log.Err(err))
			return nil, err
		}

		if err != nil {
			log.Warn("role has been deleted, but till in casbin", log.String("role", id))
			continue
		}

		roleList.Items = append(roleList.Items, *rol)
	}

	return roleList, nil
}
