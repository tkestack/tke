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

package emigration

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/business"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	"tkestack.io/tke/pkg/business/registry/namespace"
)

// _validateNsEmigrationName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var _validateNsEmigrationName = apimachineryvalidation.NameIsDNSLabel

// ValidateNsEmigrationCreate tests if required fields in the NsEmigration are set correctly.
func ValidateNsEmigrationCreate(emigration *business.NsEmigration, businessClient *businessinternalclient.BusinessClient) field.ErrorList {
	allErrs := validateNsEmigration(emigration, businessClient)

	fldNamespace := field.NewPath("spec", "namespace")
	if emigration.Spec.Namespace == "" {
		allErrs = append(allErrs, field.Invalid(fldNamespace, emigration.Spec.Namespace, "empty business namespace name"))
	}
	ns, err := businessClient.Namespaces(emigration.Namespace).Get(context.Background(), emigration.Spec.Namespace, v1.GetOptions{})
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldNamespace, emigration.Spec.Namespace,
			fmt.Sprintf("failed to get this business namespace, %s", err)))
	} else if ns.Status.Phase != business.NamespaceAvailable {
		allErrs = append(allErrs,
			field.Invalid(fldNamespace, emigration.Spec.Namespace,
				fmt.Sprintf("in phase %s, can NOT emigrate", ns.Status.Phase)))
	}

	fldDestPrj := field.NewPath("spec", "destination")
	if emigration.Spec.Destination == "" {
		allErrs = append(allErrs, field.Invalid(fldDestPrj, emigration.Spec.Destination, "empty project name"))
	} else if emigration.Spec.Destination == emigration.Namespace {
		allErrs = append(allErrs, field.Invalid(fldDestPrj, emigration.Spec.Destination, "is still the current project"))
	}
	project, err := businessClient.Projects().Get(context.Background(), emigration.Spec.Destination, v1.GetOptions{})
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldDestPrj, emigration.Spec.Destination,
			fmt.Sprintf("failed to get this project, %s", err)))
	} else {
		fldProject := field.NewPath(fmt.Sprintf("destination(%s)", emigration.Spec.Destination))
		fldHard := field.NewPath(fmt.Sprintf("namespace(%s)", emigration.Spec.Namespace), "spec", "hard")
		allErrs = append(allErrs, namespace.ValidateAgainstProject(ns, nil, project, fldProject, fldHard)...)
	}
	_, err = businessClient.Namespaces(emigration.Spec.Destination).Get(context.Background(), emigration.Spec.Namespace, v1.GetOptions{})
	if err == nil {
		allErrs = append(allErrs, field.Invalid(fldDestPrj, emigration.Spec.Destination,
			fmt.Sprintf("already has a namespace with the name %s", emigration.Spec.Namespace)))
	} else if !errors.IsNotFound(err) {
		allErrs = append(allErrs, field.Invalid(fldDestPrj, emigration.Spec.Destination,
			fmt.Sprintf("failed to check whether there is a namespace with the name %s, %s", emigration.Spec.Namespace, err)))
	}

	return allErrs
}

// ValidateNsEmigrationUpdate tests if required fields in the NsEmigration are set during
// an update.
func ValidateNsEmigrationUpdate(emigration, old *business.NsEmigration, businessClient *businessinternalclient.BusinessClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&emigration.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, validateNsEmigration(emigration, businessClient)...)

	return allErrs
}

// validateNsEmigration tests if required fields in the NsEmigration are set.
func validateNsEmigration(emigration *business.NsEmigration, businessClient *businessinternalclient.BusinessClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&emigration.ObjectMeta,
		true, _validateNsEmigrationName, field.NewPath("metadata"))

	return allErrs
}
