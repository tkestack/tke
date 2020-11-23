/*
 * Tencent is pleased to support the open source community by making TKEStack available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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
	"fmt"

	"github.com/casbin/casbin/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

func IsPlatformAdmin(ctx context.Context, username string, tenantID string, authClient authinternalclient.AuthInterface,
	enforcer *casbin.SyncedEnforcer) (bool, error) {
	idp, err := authClient.IdentityProviders().Get(ctx, tenantID, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	administrators := idp.Spec.Administrators
	for _, admin := range administrators {
		if admin == username {
			return true, nil
		}
	}

	// Use implicit roles to check admin
	roles, err := enforcer.GetImplicitRolesForUser(UserKey(tenantID, username), DefaultDomain)
	if err != nil {
		return false, err
	}
	for _, r := range roles {
		if r == fmt.Sprintf(AdministratorPolicyPattern, tenantID) {
			return true, nil
		}
	}

	return false, nil
}
