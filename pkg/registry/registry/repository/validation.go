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

package repository

import (
	"fmt"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	"tkestack.io/tke/api/registry"
)

// ValidateRepositoryName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateRepositoryName = apimachineryvalidation.NameIsDNSLabel

// ValidateRepository tests if required fields in the message are set.
func ValidateRepository(repository *registry.Repository, registryClient *registryinternalclient.RegistryClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&repository.ObjectMeta, true, ValidateRepositoryName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if repository.Spec.Name == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("name"), "must specify name"))
	}
	if repository.Spec.NamespaceName == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("namespaceName"), "must specify namespaceName"))
	}

	if repository.Spec.NamespaceName != "" && repository.Spec.Name != "" {
		namespaceList, err := registryClient.Namespaces().List(metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", repository.Spec.TenantID, repository.Spec.NamespaceName),
		})
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldSpecPath.Child("namespaceName"), err))
		} else if len(namespaceList.Items) == 0 {
			allErrs = append(allErrs, field.NotFound(fldSpecPath.Child("namespaceName"), repository.Spec.NamespaceName))
		} else {
			namespace := namespaceList.Items[0]
			if repository.ObjectMeta.Namespace != namespace.ObjectMeta.Name {
				allErrs = append(allErrs, field.NotFound(field.NewPath("metadata", "namespace"), repository.ObjectMeta.Namespace))
			}

			repoList, err := registryClient.Repositories(namespace.ObjectMeta.Name).List(metav1.ListOptions{
				FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s,spec.namespaceName=%s", repository.Spec.TenantID, repository.Spec.Name, repository.Spec.NamespaceName),
			})
			if err != nil {
				allErrs = append(allErrs, field.InternalError(fldSpecPath.Child("name"), err))
			} else if len(repoList.Items) > 0 {
				allErrs = append(allErrs, field.Duplicate(fldSpecPath.Child("name"), repository.Spec.TenantID))
			}
		}
	}

	visibilities := sets.NewString(string(registry.VisibilityPrivate), string(registry.VisibilityPublic))
	if !visibilities.Has(string(repository.Spec.Visibility)) {
		allErrs = append(allErrs, field.NotSupported(fldSpecPath.Child("visibility"), repository.Spec.Visibility, visibilities.List()))
	}

	return allErrs
}

// ValidateRepositoryUpdate tests if required fields in the repository are set
// during an update.
func ValidateRepositoryUpdate(repository *registry.Repository, old *registry.Repository, registryClient *registryinternalclient.RegistryClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&repository.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))

	if repository.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "tenantID"), "disallowed change the tenant"))
	}

	if repository.Spec.Name != old.Spec.Name {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "name"), repository.Spec.Name, "disallowed change the name"))
	}

	if repository.Spec.NamespaceName != old.Spec.NamespaceName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "namespaceName"), repository.Spec.Name, "disallowed change the namespace"))
	}

	return allErrs
}
