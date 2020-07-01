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

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"tkestack.io/tke/api/monitor"
	"tkestack.io/tke/pkg/apiserver/authentication"
)

// FilterPrometheus is used to filter helm that do not belong
// to the tenant.
func FilterPrometheus(ctx context.Context, prom *monitor.Prometheus) error {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if prom.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("prometheus"), prom.ObjectMeta.Name)
	}
	return nil
}
