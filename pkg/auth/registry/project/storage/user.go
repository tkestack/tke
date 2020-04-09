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
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"

	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/apiserver/filter"
	"tkestack.io/tke/pkg/auth/util"
	genericfilter "tkestack.io/tke/pkg/platform/apiserver/filter"
	genericutil "tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
)

// UserREST implements the REST endpoint.
type UserREST struct {
	Binding   *BindingREST
	UnBinding *UnBindingREST

	authClient authinternalclient.AuthInterface
}

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *UserREST) New() runtime.Object {
	return &auth.ProjectPolicyBindingRequest{}
}

// NewList returns an empty object that can be used with the List call.
func (r *UserREST) NewList() runtime.Object {
	return &auth.UserList{}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *UserREST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	keyword := genericfilter.FuzzyResourceFrom(ctx)
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	projectID := filter.ProjectIDFrom(ctx)
	if projectID == "" {
		projectID = requestInfo.Name
	}

	projectPolicyList, err := r.authClient.ProjectPolicyBindings().List(metav1.ListOptions{
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

	userPolicyMap := getUserPolicyMap(projectPolicyList)
	userList := &auth.UserList{}
	policyNameMap := map[string]string{}
	for userID, policyIDs := range userPolicyMap {
		user, err := r.authClient.Users().Get(util.CombineTenantAndName(tenantID, userID), metav1.GetOptions{})
		if err != nil {
			log.Error("Get user failed", log.String("id", userID), log.Err(err))
			continue
		}

		if keyword != "" && user.Spec.ID != keyword &&
			!strings.Contains(user.Spec.Name, keyword) &&
			!strings.Contains(user.Spec.DisplayName, keyword) {
			continue
		}

		m := make(map[string]string)
		for _, pid := range policyIDs {
			if name, ok := policyNameMap[pid]; ok {
				m[pid] = name
			} else {
				pol, err := r.authClient.Policies().Get(pid, metav1.GetOptions{})
				if err != nil {
					log.Error("Get policy failed", log.String("pid", pid), log.Err(err))
					continue
				}

				policyNameMap[pid] = pol.Spec.DisplayName
				m[pid] = pol.Spec.DisplayName
			}
		}

		b, err := json.Marshal(m)
		if err != nil {
			log.Error("Marshal policy map for user failed", log.String("user", userID), log.Err(err))
			continue
		}

		if user.Spec.Extra == nil {
			user.Spec.Extra = make(map[string]string)
		}
		user.Spec.Extra[util.PoliciesKey] = string(b)
		userList.Items = append(userList.Items, *user)
	}

	return userList, nil
}

func (r *UserREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	bind := obj.(*auth.ProjectPolicyBindingRequest)
	if len(bind.Users) == 0 {
		return nil, errors.NewBadRequest("must specify users")
	}
	projectID := filter.ProjectIDFrom(ctx)
	if projectID == "" {
		projectID = requestInfo.Name
	}

	projectPolicyList, err := r.authClient.ProjectPolicyBindings().List(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.projectID=%s", projectID),
	})
	if err != nil {
		log.Error("get project policies failed", log.String("project", projectID), log.Err(err))
		return nil, err
	}

	userPolicyMap := getUserPolicyMap(projectPolicyList)
	var (
		result = &auth.ProjectPolicyBindingList{}
		errs   []error
	)
	if len(bind.Policies) != 0 {
		obj, err := r.Binding.Create(ctx, obj, rest.ValidateAllObjectFunc, &metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
		result = obj.(*auth.ProjectPolicyBindingList)
	}

	for _, subj := range bind.Users {
		_, removed := genericutil.DiffStringSlice(userPolicyMap[subj.ID], bind.Policies)
		if len(removed) == 0 {
			continue
		}
		unbind := auth.ProjectPolicyBindingRequest{
			Policies: removed,
			Users: []auth.Subject{
				subj,
			},
		}

		_, err = r.UnBinding.Create(ctx, &unbind, rest.ValidateAllObjectFunc, &metav1.CreateOptions{})
		if err != nil {
			errs = append(errs, fmt.Errorf("unbind subj: %s with policies: %+v failed", subj.ID, removed))
		}
	}

	if len(errs) == 0 {
		return result, nil
	}

	return nil, apierrors.NewInternalError(utilerrors.NewAggregate(errs))
}

// GetUserPolicyMap get policies for members in project.
func getUserPolicyMap(policyList *auth.ProjectPolicyBindingList) map[string][]string {
	userPolicyMap := make(map[string][]string)
	for _, policy := range policyList.Items {
		for _, subj := range policy.Spec.Users {
			if _, ok := userPolicyMap[subj.ID]; ok {
				userPolicyMap[subj.ID] = append(userPolicyMap[subj.ID], policy.Spec.PolicyID)
			} else {
				userPolicyMap[subj.ID] = []string{policy.Spec.PolicyID}
			}
		}
	}

	return userPolicyMap
}
