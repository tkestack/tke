/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package util

import (
	"context"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/fields"
	"tkestack.io/tke/pkg/apiserver/authentication"
)

// PredicateListOptions determines the query options according to the tenant
// attribute of the request user.
func PredicateListOptions(ctx context.Context, options *metainternal.ListOptions) *metainternal.ListOptions {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return options
	}
	if options == nil {
		return &metainternal.ListOptions{
			FieldSelector: fields.OneTermEqualSelector("spec.tenantID", tenantID),
		}
	}
	if options.FieldSelector == nil {
		options.FieldSelector = fields.OneTermEqualSelector("spec.tenantID", tenantID)
		return options
	}
	options.FieldSelector = fields.AndSelectors(options.FieldSelector, fields.OneTermEqualSelector("spec.tenantID", tenantID))
	return options
}
