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

	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/auth/registry/localidentity"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"

	"github.com/casbin/casbin/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

// PasswordREST implements the REST endpoint.
type PasswordREST struct {
	localIdentityStore *registry.Store
	authClient         authinternalclient.AuthInterface
	enforcer           *casbin.SyncedEnforcer
}

var _ = rest.Creater(&PasswordREST{})

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *PasswordREST) New() runtime.Object {
	return &auth.PasswordReq{}
}

// Create used to update password of the local identity.
func (r *PasswordREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	username, tenantID := authentication.UsernameAndTenantID(ctx)
	isPlatformAdmin, err := util.IsPlatformAdmin(ctx, username, tenantID, r.authClient, r.enforcer)
	if err != nil {
		return nil, err
	}

	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, apierrors.NewBadRequest("unable to get request info from context")
	}

	userID := requestInfo.Name
	localIdentityObj, err := r.localIdentityStore.Get(ctx, userID, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	localIdentity := localIdentityObj.(*auth.LocalIdentity)
	if !isPlatformAdmin && (localIdentity.Spec.Username != username || localIdentity.Spec.TenantID != tenantID) {
		return nil, fmt.Errorf("you are not a administrator, and you are not allowd to change other users' password")
	}

	passwordReq := obj.(*auth.PasswordReq)

	if err := localidentity.ValidateLocalIdentityPasswordUpdate(localIdentity, passwordReq); err != nil {
		log.Error("Update password for localIdentity failed", log.String("localIdentity", userID), log.Err(err))
		return nil, apierrors.NewBadRequest(err.Error())
	}

	objUpdated, _, err := r.localIdentityStore.Update(ctx, userID, rest.DefaultUpdatedObjectInfo(localIdentity), rest.ValidateAllObjectFunc, rest.ValidateAllObjectUpdateFunc, false, &metav1.UpdateOptions{})
	return objUpdated, err
}
