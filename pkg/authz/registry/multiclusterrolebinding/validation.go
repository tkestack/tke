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
	"context"
	"fmt"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/tools/cache"
	"tkestack.io/tke/api/authz"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
)

var ValidateMultiClusterRoleBindingName = apimachineryvalidation.NameIsDNSLabel

// ValidateMultiClusterRoleBinding tests if required fields in the cluster are set.
func ValidateMultiClusterRoleBinding(mcrb *authz.MultiClusterRoleBinding, roleGetter rest.Getter, platformClient platformversionedclient.PlatformV1Interface) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&mcrb.ObjectMeta, true, ValidateMultiClusterRoleBindingName, field.NewPath("metadata"))
	if len(mcrb.Spec.TenantID) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "tenantID"), "empty tenantID"))
		return allErrs
	}
	if len(mcrb.Spec.Username) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "username"), "empty username"))
		return allErrs
	}
	clusters := mcrb.Spec.Clusters
	for _, cls := range clusters {
		if cls == "*" {
			if len(clusters) != 1 {
				allErrs = append(allErrs, field.Required(field.NewPath("spec", "clusters"), "cluster '*' is invalidate"))
				return allErrs
			}
		} else {
			if _, err := platformClient.Clusters().Get(context.TODO(), cls, metav1.GetOptions{ResourceVersion: "0"}); err != nil {
				allErrs = append(allErrs, field.Required(field.NewPath("spec", "clusters"), fmt.Sprintf("get cluster '%s' failed, err '%v'", cls, err)))
				return allErrs
			}
		}
	}
	roleNs, roleName, err := cache.SplitMetaNamespaceKey(mcrb.Spec.RoleName)
	if err != nil {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "roleName"), "roleName invalidate"))
		return allErrs
	}
	if roleNs != "" && roleNs != "default" && roleNs != mcrb.Namespace {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "roleName"), "roleName invalidate"))
		return allErrs
	}
	if roleNs == "" {
		roleNs = "default"
	}
	ctx := request.WithNamespace(request.NewContext(), roleNs)
	if _, err := roleGetter.Get(ctx, roleName, &metav1.GetOptions{}); err != nil {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "roleName"), fmt.Sprintf("roleName '%s' not exist", mcrb.Spec.RoleName)))
	}
	return allErrs
}

// ValidateMultiClusterRoleBindingUpdate tests if required fields in the namespace set are
// set during an update.
func ValidateMultiClusterRoleBindingUpdate(clusterroletemplatebinding *authz.MultiClusterRoleBinding, old *authz.MultiClusterRoleBinding, roleGetter rest.Getter, platformClient platformversionedclient.PlatformV1Interface) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&clusterroletemplatebinding.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateMultiClusterRoleBinding(clusterroletemplatebinding, roleGetter, platformClient)...)
	return allErrs
}
