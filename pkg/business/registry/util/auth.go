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
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	authv1 "tkestack.io/tke/api/auth/v1"
	"tkestack.io/tke/api/business"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	authversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/apiserver/filter"
	authutil "tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
)

func FilterWithUser(ctx context.Context,
	projectList *business.ProjectList,
	authClient authversionedclient.AuthV1Interface,
	businessClient *businessinternalclient.BusinessClient) (bool, *business.ProjectList, error) {

	userName, tenantID := authentication.GetUsernameAndTenantID(ctx)

	isAdmin := false
	listOpt := v1.ListOptions{FieldSelector: fmt.Sprintf("spec.tenantID=%s", tenantID)}
	platformList, err := businessClient.Platforms().List(ctx, listOpt)
	if err != nil {
		return false, nil, err
	}
	for _, platform := range platformList.Items {
		if util.InStringSlice(platform.Spec.Administrators, userName) {
			isAdmin = true
			break
		}
	}
	var userID string
	if authClient != nil {
		userList, err := authClient.Users().List(ctx, v1.ListOptions{
			FieldSelector: fields.AndSelectors(
				fields.OneTermEqualSelector("keyword", userName),
				fields.OneTermEqualSelector("policy", "true"),
				fields.OneTermEqualSelector("spec.tenantID", tenantID),
			).String(),
		})
		if err != nil {
			return false, nil, err
		}
		for _, user := range userList.Items {
			if user.Spec.Name == userName {
				log.Info("user", log.Any("user", user))
				userID = user.Name
				if authutil.IsPlatformAdministrator(user) {
					isAdmin = true
				}
				break
			}
		}
	}
	if projectList == nil || projectList.Items == nil {
		return isAdmin, projectList, nil
	}

	rawList := projectList.Items
	projectList.Items = nil
	picked := make(map[string]bool)
	if authClient == nil {
		for _, project := range rawList {
			if util.InStringSlice(project.Spec.Members, userName) && !picked[project.Name] {
				picked[project.Name] = true
				projectList.Items = append(projectList.Items, project)
			}
		}
		return isAdmin, projectList, nil
	}

	for _, project := range rawList {
		if len(project.Spec.Members) > 0 && project.Spec.Members[0] == userName && !picked[project.Name] {
			picked[project.Name] = true
			projectList.Items = append(projectList.Items, project)
		}
	}
	if userID == "" {
		return isAdmin, projectList, nil
	}

	belongs := &authv1.ProjectBelongs{}
	if err := authClient.RESTClient().Get().
		Resource("users").
		Name(userID).
		SubResource("projects").
		SetHeader(filter.HeaderTenantID, tenantID).
		Do(ctx).Into(belongs); err != nil {
		log.Error("Get user projects failed for tke-auth-api", log.String("user", userName), log.Err(err))
		return isAdmin, projectList, err
	}
	log.Debug("project belongs for user", log.String("user", userName), log.String("userID", userID), log.Any("belongs", belongs))
	for _, project := range rawList {
		if _, ok := belongs.MemberdProjects[project.Name]; ok && !picked[project.Name] {
			picked[project.Name] = true
			projectList.Items = append(projectList.Items, project)
		}
	}
	return isAdmin, projectList, nil
}
