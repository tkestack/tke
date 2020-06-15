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

package namespace

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/business"
	"tkestack.io/tke/pkg/platform/util/validation"
	"tkestack.io/tke/pkg/util/resource"
)

const (
	_clusterLimitErrorInfo = "parent project forbids its namespaces to use cluster"
)

var (
	_forbiddenNamespaces = map[string]bool{
		"kube-system":     true,
		"kube-public":     true,
		"kube-node-lease": true,
	}
)

// ValidateNamespaceName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateNamespaceName = apimachineryvalidation.NameIsDNSLabel

// ValidateAgainstProject validates namespace against project.
func ValidateAgainstProject(namespace, old *business.Namespace, project *business.Project, fldProject, fldHard *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if project.Spec.TenantID != namespace.Spec.TenantID {
		return append(allErrs, field.NotFound(fldProject, namespace.ObjectMeta.Namespace))
	}
	if project.Status.Locked != nil && *project.Status.Locked {
		return append(allErrs, field.Invalid(fldProject, namespace.ObjectMeta.Namespace, "project has been locked"))
	}
	if clusterHard, clusterHardExist := project.Spec.Clusters[namespace.Spec.ClusterName]; !clusterHardExist {
		allErrs = append(
			allErrs,
			field.Invalid(fldProject,
				namespace.ObjectMeta.Namespace,
				fmt.Sprintf("%s %s", _clusterLimitErrorInfo, namespace.Spec.ClusterName)))
	} else {
		clusterUsed, clusterUsedExist := project.Status.Clusters[namespace.Spec.ClusterName]
		if !clusterUsedExist {
			clusterUsed = business.UsedQuantity{
				Used: make(business.ResourceList),
			}
		}

		var oldSpecHard business.ResourceList
		childNamespaceNameSet := sets.NewString(project.Status.CalculatedNamespaces...)
		if childNamespaceNameSet.Has(namespace.ObjectMeta.Name) && old != nil {
			oldSpecHard = old.Spec.Hard
		}
		allErrs = append(allErrs,
			resource.ValidateAllocatableResources(namespace.Spec.Hard, oldSpecHard,
				clusterHard.Hard, clusterUsed.Used, fldHard)...)
	}
	return allErrs
}

// ValidateNamespace tests if required fields in the namespace are set.
func ValidateNamespace(ctx context.Context, namespace *business.Namespace, old *business.Namespace,
	objectGetter validation.BusinessObjectGetter, clusterGetter validation.ClusterGetter) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&namespace.ObjectMeta,
		true, ValidateNamespaceName, field.NewPath("metadata"))
	if old == nil { // Only validate cluster when creating namespaces.
		allErrs = append(allErrs,
			validation.ValidateClusterVersioned(ctx, clusterGetter, namespace.Spec.ClusterName, namespace.Spec.TenantID)...)
	}

	fldSpecPath := field.NewPath("spec")

	hardErrs := field.ErrorList{}
	fldHardPath := fldSpecPath.Child("hard")
	for k, v := range namespace.Spec.Hard {
		resPath := fldHardPath.Key(k)
		hardErrs = append(hardErrs, resource.ValidateResourceQuotaResourceName(k, resPath)...)
		hardErrs = append(hardErrs, resource.ValidateResourceQuantityValue(k, v, resPath)...)
	}
	if len(hardErrs) != 0 {
		return append(allErrs, hardErrs...)
	}

	fldProject := field.NewPath("metadata", "namespace")
	project, err := objectGetter.Project(ctx, namespace.ObjectMeta.Namespace, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return append(allErrs, field.NotFound(fldProject, namespace.ObjectMeta.Namespace))
		}
		return append(allErrs, field.InternalError(fldProject, err))
	}

	allErrs = append(allErrs, ValidateAgainstProject(namespace, old, project, fldProject, fldHardPath)...)

	fldNamespacePath := fldSpecPath.Child("namespace")
	if namespace.Spec.Namespace == "" {
		allErrs = append(allErrs, field.Required(fldNamespacePath, "must specify namespace name"))
	} else {
		ns := strings.ToLower(namespace.Spec.Namespace)
		if _, hit := _forbiddenNamespaces[ns]; hit {
			allErrs = append(allErrs,
				field.Invalid(fldNamespacePath, namespace.Spec.Namespace,
					fmt.Sprintf("cannot allocate the cluster's `%s` namespace in the project", ns)))
		}
	}

	return allErrs
}

// ValidateNamespaceUpdate tests if required fields in the namespace are set during
// an update.
func ValidateNamespaceUpdate(ctx context.Context, namespace *business.Namespace, old *business.Namespace,
	objectGetter validation.BusinessObjectGetter, clusterGetter validation.ClusterGetter) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&namespace.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateNamespace(ctx, namespace, old, objectGetter, clusterGetter)...)

	if namespace.Spec.ClusterName != old.Spec.ClusterName {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "clusterName"),
				namespace.Spec.ClusterName, "disallowed change the cluster name"))
	}

	if namespace.Spec.Namespace != old.Spec.Namespace {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "namespace"),
				namespace.Spec.Namespace, "disallowed change the namespace"))
	}

	if namespace.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "tenantID"),
				namespace.Spec.TenantID, "disallowed change the tenant"))
	}

	if namespace.Status.Phase == "" {
		allErrs = append(allErrs,
			field.Required(field.NewPath("status", "phase"), "must specify a phase"))
	}

	if old.Status.Phase != business.NamespaceAvailable && namespace.Status.Phase == business.NamespaceLocked {
		allErrs = append(allErrs,
			field.Required(field.NewPath("status", "phase"),
				fmt.Sprintf("can only lock namespaces in phase %s", business.NamespaceAvailable)))
	}

	fldPath := field.NewPath("status", "used")
	for k, v := range namespace.Status.Used {
		resPath := fldPath.Key(k)
		allErrs = append(allErrs, resource.ValidateResourceQuotaResourceName(k, resPath)...)
		allErrs = append(allErrs, resource.ValidateResourceQuantityValue(k, v, resPath)...)
	}

	allErrs = append(allErrs,
		resource.ValidateUpdateResource(namespace.Spec.Hard, old.Status.Used,
			field.NewPath("spec", "hard"))...)

	return allErrs
}

func ValidateNamespaceEmigrate(namespace *business.Namespace, old *business.Namespace) field.ErrorList {
	if namespace.Status.Phase != business.NamespaceLocked {
		return append(field.ErrorList{},
			field.Invalid(field.NewPath("status", "phase"), namespace.Status.Phase,
				fmt.Sprintf("valid value should be %s", business.NamespaceLocked)))
	}

	if old.Status.Phase != business.NamespaceAvailable {
		var msg string
		if old.Status.Phase == business.NamespaceLocked {
			msg = fmt.Sprintf("namespace is already in %s, can NOT emigrate again!", old.Status.Phase)
		} else {
			msg = fmt.Sprintf("namespace is in phase %s, can NOT emigrate!", old.Status.Phase)
		}
		return append(field.ErrorList{}, field.Invalid(field.NewPath("status", "phase"), old.Status.Phase, msg))
	}

	return nil
}
