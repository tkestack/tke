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

	"k8s.io/apimachinery/pkg/fields"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"

	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
)

// PolicyREST implements the REST endpoint, list policies bound to the user.
type PolicyREST struct {
	authClient authinternalclient.AuthInterface
}

var _ = rest.Lister(&PolicyREST{})
var _ = rest.Getter(&PolicyREST{})

// NewList returns an empty object that can be used with the List call.
func (r *PolicyREST) NewList() runtime.Object {
	return &auth.PolicyList{}
}

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *PolicyREST) New() runtime.Object {
	return &auth.Policy{}
}

// Get finds a resource in the storage by name and returns it.
func (r *PolicyREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	pol, err := r.authClient.Policies().Get(name, metav1.GetOptions{})
	if err != nil {
		log.Error("Get policy failed ", log.String("policy", name), log.Err(err))
		return nil, err
	}

	if err := util.FilterPolicy(ctx, pol); err != nil {
		return nil, err
	}
	return pol, nil
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *PolicyREST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	fieldSelector := fields.Nothing()
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID != "" {
		fieldSelector = fields.AndSelectors(fieldSelector,
			fields.OneTermEqualSelector("spec.tenantID", tenantID),
		)
	}

	fieldSelector = fields.AndSelectors(
		fieldSelector,
		fields.OneTermEqualSelector("spec.scope", string(auth.PolicyProject)))

	listOpts := metav1.ListOptions{
		FieldSelector: fieldSelector.String(),
	}

	policyList, err := r.authClient.Policies().List(listOpts)
	if err != nil {
		log.Error("List projected policy failed", log.String("selector", listOpts.FieldSelector), log.Err(err))
		return nil, err
	}

	return policyList, nil
}
