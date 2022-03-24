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

package multiclusterrolebinding

import (
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/authz"
)

var ValidateMultiClusterRoleBindingName = apimachineryvalidation.NameIsDNSLabel

// ValidateMultiClusterRoleBinding tests if required fields in the cluster are set.
func ValidateMultiClusterRoleBinding(mcrb *authz.MultiClusterRoleBinding) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&mcrb.ObjectMeta, true, ValidateMultiClusterRoleBindingName, field.NewPath("metadata"))
	clusters := mcrb.Spec.Clusters
	if len(clusters) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "clusters"), "empty clusters"))
	}
	if len(clusters) > 1 {
		for _, cls := range clusters {
			if cls == "*" {
				allErrs = append(allErrs, field.Required(field.NewPath("spec", "clusters"), "cluster '*' is invalidate"))
				break
			}
		}
	}
	return allErrs
}

// ValidateMultiClusterRoleBindingUpdate tests if required fields in the namespace set are
// set during an update.
func ValidateMultiClusterRoleBindingUpdate(clusterroletemplatebinding *authz.MultiClusterRoleBinding, old *authz.MultiClusterRoleBinding) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&clusterroletemplatebinding.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateMultiClusterRoleBinding(clusterroletemplatebinding)...)
	return allErrs
}
