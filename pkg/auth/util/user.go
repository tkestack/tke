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

package util

import (
	"context"
	"fmt"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

func GetLocalIdentity(ctx context.Context, authClient authinternalclient.AuthInterface, tenantID, username string) (auth.LocalIdentity, error) {
	tenantUserSelector := fields.AndSelectors(
		fields.OneTermEqualSelector("spec.tenantID", tenantID),
		fields.OneTermEqualSelector("spec.username", username))

	localIdentityList, err := authClient.LocalIdentities().List(ctx, v1.ListOptions{FieldSelector: tenantUserSelector.String()})
	if err != nil {
		return auth.LocalIdentity{}, err
	}

	if len(localIdentityList.Items) == 0 {
		return auth.LocalIdentity{}, apierrors.NewNotFound(auth.Resource("localIdentity"), username)
	}

	return localIdentityList.Items[0], nil
}

func UserKey(tenantID string, name string) string {
	return fmt.Sprintf("%s%s", UserPrefix(tenantID), name)
}

func UserPrefix(tenantID string) string {
	return fmt.Sprintf("%s::user::", tenantID)
}

func GroupKey(tenantID string, name string) string {
	return fmt.Sprintf("%s%s", GroupPrefix(tenantID), name)
}

func GroupPrefix(tenantID string) string {
	return fmt.Sprintf("%s::group::", tenantID)
}

func GetGroupsForUser(ctx context.Context, authClient authinternalclient.AuthInterface, userID string) (auth.LocalGroupList, error) {
	groupList := auth.LocalGroupList{}
	err := authClient.RESTClient().Get().
		Resource("localidentities").
		Name(userID).
		SubResource("groups").Do(ctx).Into(&groupList)

	return groupList, err
}

func ParseTenantAndName(str string) (string, string) {
	parts := strings.Split(str, "::")
	if len(parts) > 1 {
		return parts[0], parts[1]
	}

	return "", str
}

func CombineTenantAndName(tenantID, name string) string {
	return fmt.Sprintf("%s::%s", tenantID, name)
}
