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
	"tkestack.io/tke/pkg/auth/util"

	"tkestack.io/tke/pkg/util/log"

	"k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"

	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

// BindingREST implements the REST endpoint.
type BindingREST struct {
	*registry.Store

	authClient authinternalclient.AuthInterface
}

var _ = rest.Creater(&BindingREST{})

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *BindingREST) New() runtime.Object {
	return &auth.Binding{}
}

func (r *BindingREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	log.Info("requestinfo", log.Any("requestInfo", requestInfo))

	bind := obj.(*auth.Binding)
	polObj, err := r.Get(ctx, requestInfo.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	policy := polObj.(*auth.Policy)

	for _, sub := range bind.Subjects {
		if sub.Name != "" {
			if !util.InSubjects(sub, policy.Status.Subjects) {
				policy.Status.Subjects = append(policy.Status.Subjects, sub)
			}
		}
	}

	log.Info("bind policy subjects", log.String("policy", policy.Name), log.Any("subjects", policy.Status.Subjects))

	return r.authClient.Policies().UpdateStatus(policy)
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *BindingREST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	polObj, err := r.Get(ctx, requestInfo.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	policy := polObj.(*auth.Policy)

	localIdentityList := &auth.LocalIdentityList{}
	for _, subj := range policy.Status.Subjects {
		var localIdentity *auth.LocalIdentity
		if subj.ID != "" {
			localIdentity, err = r.authClient.LocalIdentities().Get(subj.ID, metav1.GetOptions{})
			if err != nil {
				log.Error("Get localIdentity failed", log.String("id", subj.ID), log.Err(err))
				localIdentity = constructLocalIdentity(subj.ID, subj.Name)
			}
		} else {
			localIdentity = constructLocalIdentity(subj.ID, subj.Name)
		}

		localIdentity.Spec.HashedPassword = ""
		localIdentityList.Items = append(localIdentityList.Items, *localIdentity)
	}

	return localIdentityList, nil
}

func constructLocalIdentity(userID, username string) *auth.LocalIdentity {
	return &auth.LocalIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name: userID,
		},
		Spec: auth.LocalIdentitySpec{
			Username: username,
		},
	}
}
