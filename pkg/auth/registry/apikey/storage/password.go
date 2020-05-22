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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/auth/registry/apikey"
	"tkestack.io/tke/pkg/auth/util"
)

// PasswordREST implements the REST endpoint.
type PasswordREST struct {
	apiKeyStore *registry.Store
	keySigner   util.KeySigner

	authClient authinternalclient.AuthInterface
}

var _ = rest.Creater(&PasswordREST{})

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *PasswordREST) New() runtime.Object {
	return &auth.APIKeyReqPassword{}
}

func (r *PasswordREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {

	apikeyReq := obj.(*auth.APIKeyReqPassword)

	if err := apikey.ValidateAPIkeyPassword(ctx, apikeyReq, r.authClient); err != nil {
		return nil, apierrors.NewBadRequest(err.Error())
	}
	apiKey, err := r.keySigner.Generate(ctx, apikeyReq.Username, apikeyReq.TenantID, apikeyReq.Expire.Duration)
	if err != nil {
		return nil, apierrors.NewBadRequest(err.Error())
	}
	apiKey.Spec.Description = apikeyReq.Description

	return r.apiKeyStore.Create(ctx, apiKey, createValidation, options)
}
