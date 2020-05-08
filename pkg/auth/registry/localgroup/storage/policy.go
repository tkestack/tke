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
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

// PolicyREST implements the REST endpoint, list policies bound to the user.
type PolicyREST struct {
	localGroupStore *registry.Store
	authClient      authinternalclient.AuthInterface
	enforcer        *casbin.SyncedEnforcer
}

var _ = rest.Lister(&PolicyREST{})

// NewList returns an empty object that can be used with the List call.
func (r *PolicyREST) NewList() runtime.Object {
	return &auth.PolicyList{}
}

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *PolicyREST) New() runtime.Object {
	return &auth.Policy{}
}

func (r *PolicyREST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	groupID := requestInfo.Name

	obj, err := r.localGroupStore.Get(ctx, groupID, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	localGroup := obj.(*auth.LocalGroup)
	projectID := filter.ProjectIDFrom(ctx)
	roles := r.enforcer.GetRolesForUserInDomain(util.GroupKey(localGroup.Spec.TenantID, localGroup.Name), projectID)

	var policyIDs []string
	for _, r := range roles {
		if strings.HasPrefix(r, "pol-") {
			policyIDs = append(policyIDs, r)
		}
	}

	var policyList = &auth.PolicyList{}
	for _, id := range policyIDs {
		pol, err := r.authClient.Policies().Get(id, metav1.GetOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			log.Error("Get pol failed", log.String("policy", id), log.Err(err))
			return nil, err
		}

		if err != nil {
			log.Warn("Pol has been deleted, but till in casbin", log.String("policy", id))
			continue
		}

		if projectID == "" || (projectID != "" && pol.Spec.Scope == auth.PolicyProject) {
			policyList.Items = append(policyList.Items, *pol)
		}
	}

	return policyList, nil
}
