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
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"tkestack.io/tke/pkg/auth/types"
)

// ResourceAttributesFrom combines the API object information and the user.Info from the context to build a full authorizer.AttributesRecord for resource access.
func ResourceAttributesFrom(user user.Info, in types.ResourceAttributes) authorizer.AttributesRecord {
	return authorizer.AttributesRecord{
		User:            user,
		Verb:            in.Verb,
		Namespace:       in.Namespace,
		APIGroup:        in.Group,
		APIVersion:      in.Version,
		Resource:        in.Resource,
		Subresource:     in.Subresource,
		Name:            in.Name,
		ResourceRequest: true,
	}
}

// NonResourceAttributesFrom combines the API object information and the user.Info from the context to build a full authorizer.AttributesRecord for non resource access.
// Tke-auth considers non-resource path as the resource field.
func NonResourceAttributesFrom(user user.Info, in types.NonResourceAttributes) authorizer.AttributesRecord {
	return authorizer.AttributesRecord{
		User:            user,
		ResourceRequest: false,
		Resource:        in.Path,
		Verb:            in.Verb,
	}
}

// AuthorizationAttributesFrom takes a spec and returns the proper authz attributes to check it.
func AuthorizationAttributesFrom(spec types.SubjectAccessReviewSpec) authorizer.AttributesRecord {
	userToCheck := &user.DefaultInfo{
		Name:  spec.User,
		UID:   spec.UID,
		Extra: spec.Extra,
	}

	var authorizationAttributes authorizer.AttributesRecord
	if spec.ResourceAttributes != nil {
		authorizationAttributes = ResourceAttributesFrom(userToCheck, *spec.ResourceAttributes)
	} else {
		authorizationAttributes = NonResourceAttributesFrom(userToCheck, *spec.NonResourceAttributes)
	}

	return authorizationAttributes
}

// AuthorizationAttributesListFrom takes a spec and returns the proper authz attribute list to check it.
func AuthorizationAttributesListFrom(spec types.SubjectAccessReviewSpec) []authorizer.AttributesRecord {
	userToCheck := &user.DefaultInfo{
		Name:  spec.User,
		UID:   spec.UID,
		Extra: spec.Extra,
	}

	var authorizationAttributesList []authorizer.AttributesRecord
	for _, resAttr := range spec.ResourceAttributesList {
		attr := ResourceAttributesFrom(userToCheck, *resAttr)
		authorizationAttributesList = append(authorizationAttributesList, attr)
	}

	return authorizationAttributesList
}
