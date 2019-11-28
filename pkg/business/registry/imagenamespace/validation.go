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

package imagenamespace

import (
	"fmt"

	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/business"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	registryversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
)

// _validateImageNamespaceName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var _validateImageNamespaceName = apimachineryvalidation.NameIsDNSLabel

// ValidateImageNamespaceCreate tests if required fields in the ImageNamespace are set correctly.
func ValidateImageNamespaceCreate(imageNamespace *business.ImageNamespace,
	businessClient *businessinternalclient.BusinessClient, registryClient registryversionedclient.RegistryV1Interface) field.ErrorList {
	allErrs := validateImageNamespace(imageNamespace, businessClient, registryClient)

	fldName := field.NewPath("spec", "name")
	namespaceList, err := registryClient.Namespaces().List(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", imageNamespace.Spec.TenantID, imageNamespace.Name),
	})
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldName, fmt.Errorf("failed to check imageNamespace conflicting, for %s", err)))
	} else if len(namespaceList.Items) != 0 {
		allErrs = append(allErrs, field.Invalid(fldName, imageNamespace.Name, "already has a registry namespace with the same name"))
	}

	return allErrs
}

// ValidateImageNamespaceUpdate tests if required fields in the ImageNamespace are set during
// an update.
func ValidateImageNamespaceUpdate(imageNamespace *business.ImageNamespace, old *business.ImageNamespace,
	businessClient *businessinternalclient.BusinessClient, registryClient registryversionedclient.RegistryV1Interface) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&imageNamespace.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, validateImageNamespace(imageNamespace, businessClient, registryClient)...)

	if imageNamespace.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "tenantID"),
				imageNamespace.Spec.TenantID, "disallowed change the tenant"))
	}

	if imageNamespace.Status.Phase == "" {
		allErrs = append(allErrs,
			field.Required(field.NewPath("status", "phase"), "must specify a phase"))
	}

	return allErrs
}

// validateImageNamespace tests if required fields in the ImageNamespace are set.
func validateImageNamespace(imageNamespace *business.ImageNamespace,
	businessClient *businessinternalclient.BusinessClient, registryClient registryversionedclient.RegistryV1Interface) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&imageNamespace.ObjectMeta,
		true, _validateImageNamespaceName, field.NewPath("metadata"))

	return allErrs
}
