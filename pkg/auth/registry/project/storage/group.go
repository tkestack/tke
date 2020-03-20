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

	"tkestack.io/tke/pkg/apiserver/filter"

	"k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
)

// GroupREST implements the REST endpoint.
type GroupREST struct {
	authClient authinternalclient.AuthInterface
}

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *GroupREST) New() runtime.Object {
	return &auth.Group{}
}

// NewList returns an empty object that can be used with the List call.
func (r *GroupREST) NewList() runtime.Object {
	return &auth.GroupList{}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *GroupREST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}
	projectID := filter.ProjectIDFrom(ctx)
	if projectID == "" {
		projectID = requestInfo.Name
	}

	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	projectPolicyList, err := r.authClient.ProjectPolicyBindings().List(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.projectID=%s", projectID),
	})
	if err != nil {
		log.Error("get project policies failed", log.String("project", projectID), log.Err(err))
		return nil, err
	}

	groupPolicyMap := getGroupsPolicyMap(projectPolicyList)

	if tenantID == "" && len(projectPolicyList.Items) > 0 {
		tenantID = projectPolicyList.Items[0].Spec.TenantID
	}

	groupList := &auth.GroupList{}
	policyNameMap := map[string]string{}
	for groupID, policyIDs := range groupPolicyMap {
		group, err := r.authClient.Groups().Get(util.CombineTenantAndName(tenantID, groupID), metav1.GetOptions{})
		if err != nil {
			log.Error("Get group failed", log.String("id", groupID), log.Err(err))
			continue
		}

		group.Spec.Extra = make(map[string]string)
		for _, pid := range policyIDs {
			if name, ok := policyNameMap[pid]; ok {
				group.Spec.Extra[pid] = name
			} else {
				pol, err := r.authClient.Policies().Get(pid, metav1.GetOptions{})
				if err != nil {
					log.Error("Get policy failed", log.String("pid", pid), log.Err(err))
					continue
				}

				policyNameMap[pid] = pol.Spec.DisplayName
				group.Spec.Extra[pid] = pol.Spec.DisplayName
			}
		}

		groupList.Items = append(groupList.Items, *group)
	}

	return groupList, nil
}

// getGroupsPolicyMap get policies for groups in project.
func getGroupsPolicyMap(policyList *auth.ProjectPolicyBindingList) map[string][]string {
	groupPolicyMap := make(map[string][]string)
	for _, policy := range policyList.Items {
		for _, subj := range policy.Spec.Groups {
			if _, ok := groupPolicyMap[subj.ID]; ok {
				groupPolicyMap[subj.ID] = append(groupPolicyMap[subj.ID], policy.Spec.PolicyID)
			} else {
				groupPolicyMap[subj.ID] = []string{policy.Spec.PolicyID}
			}
		}
	}

	return groupPolicyMap
}
