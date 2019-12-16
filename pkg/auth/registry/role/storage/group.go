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

// GroupREST implements the REST endpoint.
type GroupREST struct {
	roleStore *registry.Store

	authClient authinternalclient.AuthInterface
}

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *GroupREST) New() runtime.Object {
	return &auth.Binding{}
}

// NewList returns an empty object that can be used with the List call.
func (r *GroupREST) NewList() runtime.Object {
	return &auth.LocalGroupList{}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *GroupREST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	rolObj, err := r.roleStore.Get(ctx, requestInfo.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	role := rolObj.(*auth.Role)

	groupList := &auth.LocalGroupList{}
	for _, subj := range role.Status.Groups {
		var group *auth.LocalGroup
		if subj.ID != "" {
			group, err = r.authClient.LocalGroups().Get(subj.ID, metav1.GetOptions{})
			if err != nil {
				log.Error("Get group failed", log.String("id", subj.ID), log.Err(err))
				group = constructgroup(subj.ID, subj.Name)
			}
		} else {
			group = constructgroup(subj.ID, subj.Name)
		}

		groupList.Items = append(groupList.Items, *group)
	}

	return groupList, nil
}

func constructgroup(userID, groupName string) *auth.LocalGroup {
	return &auth.LocalGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: userID,
		},
		Spec: auth.LocalGroupSpec{
			Username: groupName,
		},
	}
}
