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

	"tkestack.io/tke/pkg/apiserver/filter"

	"tkestack.io/tke/pkg/auth/util"

	"tkestack.io/tke/pkg/util/log"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"

	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

// ProjectUnBindingREST implements the REST endpoint.
type ProjectUnBindingREST struct {
	policyStore *registry.Store

	authClient authinternalclient.AuthInterface
}

var _ = rest.Creater(&ProjectUnBindingREST{})

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *ProjectUnBindingREST) New() runtime.Object {
	return &auth.Binding{}
}

func (r *ProjectUnBindingREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
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
		return nil, errors.NewBadRequest("unable unbind subject to platform-scoped policy, please use binding api")
	}

	projectID := filter.ProjectIDFrom(ctx)
	if projectID == "" {
		return nil, errors.NewBadRequest("must specify policyID and projectID")
	}

	bind := obj.(*auth.Binding)
	projectPolicyBinding, err := r.authClient.ProjectPolicyBindings().Get(ctx, util.ProjectPolicyName(projectID, policyID), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	remainedUsers := make([]auth.Subject, 0)
	for _, sub := range projectPolicyBinding.Spec.Users {
		if !util.InSubjects(sub, bind.Users) {
			remainedUsers = append(remainedUsers, sub)
		}
	}

	projectPolicyBinding.Spec.Users = remainedUsers
	remainedGroups := make([]auth.Subject, 0)
	for _, sub := range projectPolicyBinding.Spec.Groups {
		if !util.InSubjects(sub, bind.Groups) {
			remainedGroups = append(remainedGroups, sub)
		}
	}

	projectPolicyBinding.Spec.Groups = remainedGroups
	log.Info("unbind policy subjects", log.String("policy", projectPolicyBinding.Name), log.Any("users", projectPolicyBinding.Spec.Users), log.Any("groups", projectPolicyBinding.Spec.Groups))
	return r.authClient.ProjectPolicyBindings().Update(ctx, projectPolicyBinding, metav1.UpdateOptions{})
}
