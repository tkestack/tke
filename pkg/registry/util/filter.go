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
	"k8s.io/apimachinery/pkg/api/errors"
	"tkestack.io/tke/api/registry"
	"tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
)

// FilterNamespace is used to filter namespaces that do not belong to the tenant.
func FilterNamespace(ctx context.Context, namespace *registry.Namespace) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if namespace.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("namespace"), namespace.ObjectMeta.Name)
	}
	return nil
}

// FilterRepository is used to filter repositories that do not belong to the tenant.
func FilterRepository(ctx context.Context, repository *registry.Repository) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if repository.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("repository"), repository.ObjectMeta.Name)
	}
	return nil
}
