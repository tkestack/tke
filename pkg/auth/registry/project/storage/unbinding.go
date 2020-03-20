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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"

	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/apiserver/filter"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
)

// UnBindingREST implements the REST endpoint.
type UnBindingREST struct {
	authClient authinternalclient.AuthInterface
}

var _ = rest.Creater(&UnBindingREST{})

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *UnBindingREST) New() runtime.Object {
	return &auth.ProjectPolicyBindingRequest{}
}

func (r *UnBindingREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	bind := obj.(*auth.ProjectPolicyBindingRequest)
	if len(bind.Policies) == 0 {
		return nil, errors.NewBadRequest("must specify policies")
	}

	if len(bind.Users) == 0 && len(bind.Groups) == 0 {
		return nil, errors.NewBadRequest("must specify users or groups")
	}

	projectID := filter.ProjectIDFrom(ctx)
	if projectID == "" {
		projectID = requestInfo.Name
	}

	projectPolicyList := &auth.ProjectPolicyBindingList{}
	var errs []error
	for _, policyID := range bind.Policies {

		policy, err := r.authClient.Policies().Get(policyID, metav1.GetOptions{})
		if err != nil {
			log.Error("Get policy failed", log.String("policy", policyID), log.Err(err))
			errs = append(errs, err)
			continue
		}

		if policy.Spec.Scope != auth.PolicyProject {
			errs = append(errs, fmt.Errorf("unable bind subject to platform-scoped policy %s in project", policyID))
			continue
		}

		projectPolicy, err := r.authClient.ProjectPolicyBindings().Get(util.ProjectPolicyName(projectID, policy.Name), metav1.GetOptions{})
		if err != nil && apierrors.IsNotFound(err) {
			// if projectPolicy not exist, create a new one
			projectPolicy, err = r.authClient.ProjectPolicyBindings().Create(&auth.ProjectPolicyBinding{
				Spec: auth.ProjectPolicyBindingSpec{
					TenantID:  policy.Spec.TenantID,
					ProjectID: projectID,
					PolicyID:  policy.Name,
				},
			})
			if err != nil {
				if apierrors.IsAlreadyExists(err) {
					projectPolicy, err = r.authClient.ProjectPolicyBindings().Get(util.ProjectPolicyName(projectID, policy.Name), metav1.GetOptions{})
				}
			}
		}

		if err != nil {
			log.Error("Create or get policy failed", log.String("policyID", policyID), log.Err(err))
			errs = append(errs, err)
			continue
		}

		remainedUsers := make([]auth.Subject, 0)
		for _, sub := range bind.Users {
			if !util.InSubjects(sub, bind.Users) {
				remainedUsers = append(remainedUsers, sub)
			}
		}
		projectPolicy.Spec.Users = remainedUsers

		remainedGroups := make([]auth.Subject, 0)
		for _, sub := range bind.Groups {
			if !util.InSubjects(sub, bind.Groups) {
				remainedGroups = append(remainedGroups, sub)
			}
		}
		projectPolicy.Spec.Groups = remainedGroups

		log.Info("unbind policy subjects", log.String("policy", projectPolicy.Name), log.Any("users", projectPolicy.Spec.Users), log.Any("groups", projectPolicy.Spec.Groups))
		projectPolicy, err = r.authClient.ProjectPolicyBindings().Update(projectPolicy)
		if err != nil {
			log.Error("Update project policy failed", log.String("policyID", projectPolicy.Name), log.Err(err))
			errs = append(errs, err)
		}

		projectPolicyList.Items = append(projectPolicyList.Items, *projectPolicy)

	}

	if len(errs) == 0 {
		return projectPolicyList, nil
	}

	return nil, utilerrors.NewAggregate(errs)
}
