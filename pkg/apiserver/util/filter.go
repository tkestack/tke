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
	"tkestack.io/tke/pkg/apiserver/filter"
)

// PredicateListOptions determines the query options according to the tenant
// attribute of the request user.
func PredicateListOptions(ctx context.Context, options *metainternal.ListOptions) *metainternal.ListOptions {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID == "" {
		tenantID = filter.TenantIDFrom(ctx)
		if tenantID == "" {
			return options
		}
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

// FullListOptionsFieldSelector fill options fieldSelector.
func FullListOptionsFieldSelector(options *metainternal.ListOptions, fieldSelector fields.Selector) *metainternal.ListOptions {
	if options == nil {
		return &metainternal.ListOptions{
			FieldSelector: fieldSelector,
		}
	}
	if options.FieldSelector == nil {
		options.FieldSelector = fieldSelector
		return options
	}
	options.FieldSelector = fields.AndSelectors(options.FieldSelector, fieldSelector)
	return options
}

// InterceptFuzzyResourceNameFromListOptions determines the query options according to the fuzzyResourceName.
func InterceptFuzzyResourceNameFromListOptions(options *metainternal.ListOptions, fuzzyResourceName string) (*metainternal.ListOptions, string) {
	return InterceptCustomSelectorFromListOptions(options, "metadata.name", fuzzyResourceName)
}

// InterceptCustomSelectorFromListOptions determines the query options according to the selector.
func InterceptCustomSelectorFromListOptions(options *metainternal.ListOptions, selector, defaultValue string) (*metainternal.ListOptions, string) {
	if options != nil && options.FieldSelector != nil {
		if name, ok := options.FieldSelector.RequiresExactMatch(selector); ok {
			options.FieldSelector, _ = options.FieldSelector.Transform(func(k, v string) (string, string, error) {
				if k == selector {
					return "", "", nil
				}
				return k, v, nil
			})
			defaultValue = name
		}
	}
	return options, defaultValue
}
