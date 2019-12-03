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
	"github.com/casbin/casbin/v2"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/util/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BindingREST implements the REST endpoint.
type BindingREST struct {
	*registry.Store
	enforcer    *casbin.SyncedEnforcer
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

	bind := obj.(*auth.Binding)
	polObj, err := r.Get(ctx, requestInfo.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	policy := polObj.(*auth.Policy)
	var errs []error
	for _, user := range bind.Subjects {
		if _, err := r.enforcer.Enforcer.AddRoleForUser(fmt.Sprintf("%s::%s", policy.Spec.TenantID, user.Name), requestInfo.Name); err != nil {
			log.Error("Add policy for user failed", log.String("pol", requestInfo.Name), log.String("user", user.Name), log.Err(err))
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil, nil
	}

	return nil, errors.NewInternalError(utilerrors.NewAggregate(errs))
}
