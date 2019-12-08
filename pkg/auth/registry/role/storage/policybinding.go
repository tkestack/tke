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
	"tkestack.io/tke/pkg/util"

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

// PolicyBindingREST implements the REST endpoint.
type PolicyBindingREST struct {
	*registry.Store

	authClient authinternalclient.AuthInterface
}

var _ = rest.Creater(&PolicyBindingREST{})

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *PolicyBindingREST) New() runtime.Object {
	return &auth.PolicyBinding{}
}

func (r *PolicyBindingREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	log.Info("requestinfo", log.Any("requestInfo", requestInfo))

	bind := obj.(*auth.PolicyBinding)
	polObj, err := r.Get(ctx, requestInfo.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	role := polObj.(*auth.Role)

	for _, pid := range bind.Policies {
		if !util.InStringSlice(role.Spec.Policies, pid) {
			role.Spec.Policies = append(role.Spec.Policies, pid)
		}
	}

	log.Info("role policies", log.String("role", role.Name), log.Any("policies", role.Spec.Policies))

	return r.authClient.Roles().Update(role)
}
