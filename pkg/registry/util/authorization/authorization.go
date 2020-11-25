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

package authentication

import (
	"context"

	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/apiserver/filter"
	authfilter "tkestack.io/tke/pkg/auth/filter"
	"tkestack.io/tke/pkg/util/log"
)

// ConvertTKEAttributesForChart combines the API object information and the user.Info from the context to build a full authorizer.AttributesRecord for resource access.
func ConvertTKEAttributesForChart(ctx context.Context, u user.Info, verb string, cg registryv1.ChartGroup, chartName string) (in authorizer.Attributes, err error) {
	attribs := authorizer.AttributesRecord{}
	attribs.User = u

	// If ChartGroup belongs to project, set user extra group
	if cg.Spec.Visibility == registryv1.VisibilityProject && len(cg.Spec.Projects) > 0 {
		if userInfo, ok := u.(*user.DefaultInfo); ok {
			userInfo.Groups = append(userInfo.Groups, filter.GroupWithProject(cg.Spec.Projects[0]))
			attribs.User = userInfo
		}
	}

	// Start with common attributes that apply to resource and non-resource requests
	attribs.ResourceRequest = true
	attribs.Verb = verb

	attribs.APIGroup = registryv1.GroupName
	attribs.APIVersion = registryv1.Version
	attribs.Resource = "charts"
	attribs.Subresource = ""
	attribs.Namespace = cg.Name
	attribs.Name = chartName

	tkeAttributes := authfilter.ConvertTKEAttributes(ctx, &attribs)

	return tkeAttributes, nil
}

// ConvertTKEAttributesForChartGroup combines the API object information and the user.Info from the context to build a full authorizer.AttributesRecord for resource access.
func ConvertTKEAttributesForChartGroup(ctx context.Context, u user.Info, verb string, cg registryv1.ChartGroup) (in authorizer.Attributes, err error) {
	attribs := authorizer.AttributesRecord{}
	attribs.User = u

	// If ChartGroup belongs to project, set user extra group
	if cg.Spec.Visibility == registryv1.VisibilityProject && len(cg.Spec.Projects) > 0 {
		if userInfo, ok := u.(*user.DefaultInfo); ok {
			userInfo.Groups = append(userInfo.Groups, filter.GroupWithProject(cg.Spec.Projects[0]))
			attribs.User = userInfo
		}
	}

	// Start with common attributes that apply to resource and non-resource requests
	attribs.ResourceRequest = true
	attribs.Verb = verb

	attribs.APIGroup = registryv1.GroupName
	attribs.APIVersion = registryv1.Version
	attribs.Resource = "chartgroups"
	attribs.Subresource = ""
	attribs.Namespace = ""
	attribs.Name = cg.Name

	tkeAttributes := authfilter.ConvertTKEAttributes(ctx, &attribs)

	return tkeAttributes, nil
}

// AuthorizeForChart check if chart resource is authorized
func AuthorizeForChart(ctx context.Context, u user.Info, authzer authorizer.Authorizer, verb string, cg registryv1.ChartGroup, chartName string) (passed bool, err error) {
	tkeAttributes, err := ConvertTKEAttributesForChart(ctx, u, verb, cg, chartName)
	if err != nil {
		return false, err
	}
	authorized, reason, err := authzer.Authorize(ctx, tkeAttributes)
	if err != nil {
		return false, err
	}
	if authorized != authorizer.DecisionAllow {
		log.Warn(reason)
	}
	return authorized == authorizer.DecisionAllow, nil
}

// AuthorizeForChartGroup check if chartgroup resource is authorized
func AuthorizeForChartGroup(ctx context.Context, u user.Info, authzer authorizer.Authorizer, verb string, cg registryv1.ChartGroup) (passed bool, err error) {
	tkeAttributes, err := ConvertTKEAttributesForChartGroup(ctx, u, verb, cg)
	if err != nil {
		return false, err
	}
	authorized, reason, err := authzer.Authorize(ctx, tkeAttributes)
	if err != nil {
		return false, err
	}
	if authorized != authorizer.DecisionAllow {
		log.Warn(reason)
	}
	return authorized == authorizer.DecisionAllow, nil
}
