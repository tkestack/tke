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

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"

	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/util/log"
)

// UnbindingREST implements the REST endpoint.
type UnbindingREST struct {
	roleStore *registry.Store

	authClient authinternalclient.AuthInterface
}

var _ = rest.Creater(&UnbindingREST{})

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *UnbindingREST) New() runtime.Object {
	return &auth.Binding{}
}

func (r *UnbindingREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	bind := obj.(*auth.Binding)
	polObj, err := r.roleStore.Get(ctx, requestInfo.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	role := polObj.(*auth.Role)
	remainedUsers := make([]auth.Subject, 0)
	for _, sub := range role.Status.Users {
		if !util.InSubjects(sub, bind.Users) {
			remainedUsers = append(remainedUsers, sub)
		}
	}

	role.Status.Users = remainedUsers

	remainedGroups := make([]auth.Subject, 0)
	for _, sub := range role.Status.Groups {
		if !util.InSubjects(sub, bind.Groups) {
			remainedGroups = append(remainedGroups, sub)
		}
	}

	role.Status.Groups = remainedGroups
	log.Info("unbind role subjects", log.String("role", role.Name), log.Any("users", role.Status.Users), log.Any("groups", role.Status.Groups))
	return r.authClient.Roles().UpdateStatus(role)
}
