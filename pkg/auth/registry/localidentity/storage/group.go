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

	"k8s.io/apiserver/pkg/registry/generic/registry"

	"github.com/casbin/casbin/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/endpoints/request"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

// GroupREST implements the REST endpoint, list policies bound to the user.
type GroupREST struct {
	localIdentityStore *registry.Store
	authClient         authinternalclient.AuthInterface
	enforcer           *casbin.SyncedEnforcer
}

var _ = rest.Lister(&GroupREST{})

// NewList returns an empty object that can be used with the List call.
func (r *GroupREST) NewList() runtime.Object {
	return &auth.LocalGroupList{}
}

// New returns an empty object that can be used with Create and Update after
// request data has been put into it.
func (r *GroupREST) New() runtime.Object {
	return &auth.LocalGroup{}
}

// ConvertToTable converts objects to metav1.Table objects using default table
// convertor.
func (r *GroupREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: convert role list to table
	tableConvertor := rest.NewDefaultTableConvertor(auth.Resource("groups"))
	return tableConvertor.ConvertToTable(ctx, object, tableOptions)
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *GroupREST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	userID := requestInfo.Name

	obj, err := r.localIdentityStore.Get(ctx, userID, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	localIdentity := obj.(*auth.LocalIdentity)

	roles := r.enforcer.GetRolesForUserInDomain(util.UserKey(localIdentity.Spec.TenantID, localIdentity.Spec.Username), util.DefaultDomain)

	var groupIDs []string
	for _, r := range roles {
		if strings.HasPrefix(r, util.GroupPrefix(localIdentity.Spec.TenantID)) {
			groupIDs = append(groupIDs, strings.TrimPrefix(r, util.GroupPrefix(localIdentity.Spec.TenantID)))
		}
	}

	var groupList = &auth.GroupList{}
	for _, id := range groupIDs {
		grp, err := r.authClient.Groups().Get(ctx, util.CombineTenantAndName(localIdentity.Spec.TenantID, id), metav1.GetOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			log.Error("Get group failed", log.String("group", id), log.Err(err))
			return nil, err
		}

		if err != nil {
			_, _ = r.enforcer.DeleteRoleForUserInDomain(util.UserKey(localIdentity.Spec.TenantID, localIdentity.Spec.Username),
				util.GroupKey(localIdentity.Spec.TenantID, id), util.DefaultDomain)
			log.Warn("group has been deleted, but till in casbin", log.String("group", id))
			continue
		}

		groupList.Items = append(groupList.Items, *grp)
	}

	return groupList, nil
}
