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

	"k8s.io/apiserver/pkg/registry/rest"

	"tkestack.io/tke/pkg/apiserver/filter"

	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/endpoints/request"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/auth/util"
	genericutil "tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
)

// GroupREST implements the REST endpoint.
type GroupREST struct {
	Binding   *BindingREST
	UnBinding *UnBindingREST

	authClient authinternalclient.AuthInterface
}

var _ rest.Creater = &GroupREST{}
var _ rest.Lister = &GroupREST{}

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *GroupREST) New() runtime.Object {
	return &auth.Group{}
}

// NewList returns an empty object that can be used with the List call.
func (r *GroupREST) NewList() runtime.Object {
	return &auth.GroupList{}
}

// ConvertToTable converts objects to metav1.Table objects using default table
// convertor.
func (r *GroupREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: convert role list to table
	tableConvertor := rest.NewDefaultTableConvertor(auth.Resource("groups"))
	return tableConvertor.ConvertToTable(ctx, object, tableOptions)
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
	projectPolicyList, err := r.authClient.ProjectPolicyBindings().List(ctx, metav1.ListOptions{
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
		group, err := r.authClient.Groups().Get(ctx, util.CombineTenantAndName(tenantID, groupID), metav1.GetOptions{})
		if err != nil {
			log.Error("Get group failed", log.String("id", groupID), log.Err(err))
			continue
		}

		m := make(map[string]string)
		for _, pid := range policyIDs {
			if name, ok := policyNameMap[pid]; ok {
				m[pid] = name
			} else {
				pol, err := r.authClient.Policies().Get(ctx, pid, metav1.GetOptions{})
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
			log.Error("Marshal policy map for group failed", log.String("group", groupID), log.Err(err))
			continue
		}

		if group.Spec.Extra == nil {
			group.Spec.Extra = make(map[string]string)
		}
		group.Spec.Extra[util.PoliciesKey] = string(b)
		groupList.Items = append(groupList.Items, *group)
	}

	return groupList, nil
}

func (r *GroupREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
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

	_, tenantID := authentication.GetUsernameAndTenantID(ctx)

	projectPolicyList, err := r.authClient.ProjectPolicyBindings().List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.projectID=%s", projectID),
	})

	if err != nil {
		log.Error("get project policies failed", log.String("project", projectID), log.Err(err))
		return nil, err
	}
	groupPolicyMap := getGroupsPolicyMap(projectPolicyList)
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

	for _, subj := range bind.Groups {
		_, removed := genericutil.DiffStringSlice(groupPolicyMap[subj.ID], bind.Policies)
		if len(removed) == 0 {
			continue
		}
		unbind := auth.ProjectPolicyBindingRequest{
			TenantID: tenantID,
			Policies: removed,
			Groups: []auth.Subject{
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
