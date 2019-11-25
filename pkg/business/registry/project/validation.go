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

package project

import (
	"fmt"

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
	_addNewQuotaErrorInfo  = "Please at first set quota"
	_clusterLimitErrorInfo = "parent project forbids this project to use cluster"
)

// ValidateProjectName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateProjectName = apimachineryvalidation.NameIsDNSLabel

// ValidateProject tests if required fields in the project are set.
func ValidateProject(project *business.Project, old *business.Project,
	objectGetter validation.BusinessObjectGetter, clusterGetter validation.ClusterGetter) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&project.ObjectMeta,
		false, ValidateProjectName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if project.Spec.DisplayName == "" {
		allErrs = append(allErrs,
			field.Required(fldSpecPath.Child("displayName"), "must specify display name"))
	}

	hardErrs := field.ErrorList{}
	fldHardPath := fldSpecPath.Child("hard")
	if len(project.Spec.Clusters) == 0 {
		allErrs = append(allErrs, field.Invalid(fldHardPath, nil, "must bind at least one cluster"))
	}
	for clusterName, clusterHard := range project.Spec.Clusters {
		for k, v := range clusterHard.Hard {
			hardErrs = append(hardErrs, validation.ValidateClusterVersioned(clusterGetter, clusterName, project.Spec.TenantID)...)
			resPath := fldHardPath.Key(clusterName + k)
			hardErrs = append(hardErrs, resource.ValidateResourceQuotaResourceName(k, resPath)...)
			hardErrs = append(hardErrs, resource.ValidateResourceQuantityValue(k, v, resPath)...)
		}
	}
	if len(hardErrs) != 0 {
		return append(allErrs, hardErrs...)
	}

	if len(project.Status.CalculatedChildProjects) != 0 || len(project.Status.CalculatedNamespaces) != 0 {
		allErrs = append(allErrs, validateAgainstChildren(project, objectGetter)...)
	}

	if project.Spec.ParentProjectName != "" {
		allErrs = append(allErrs, validateAgainstParent(project, old, objectGetter)...) /* need hardErrs == nil */
	}

	return allErrs
}

// ValidateProjectUpdate tests if required fields in the project are set during
// an update.
func ValidateProjectUpdate(project *business.Project, old *business.Project,
	objectGetter validation.BusinessObjectGetter, clusterGetter validation.ClusterGetter) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&project.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateProject(project, old, objectGetter, clusterGetter)...)

	if project.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "tenantID"),
			project.Spec.TenantID, "disallowed change the tenant"))
	}

	if project.Spec.ParentProjectName != old.Spec.ParentProjectName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "parentProjectName"),
			project.Spec.ParentProjectName, "disallowed change the parent project"))
	}

	fldPath := field.NewPath("status", "used")
	for clusterName, clusterUsed := range project.Status.Clusters {
		for k, v := range clusterUsed.Used {
			resPath := fldPath.Key(clusterName + "." + k)
			allErrs = append(allErrs, resource.ValidateResourceQuotaResourceName(k, resPath)...)
			allErrs = append(allErrs, resource.ValidateResourceQuantityValue(k, v, resPath)...)
		}
	}

	for clusterName, clusterHard := range project.Spec.Clusters {
		clusterUsed, clusterUsedExist := project.Status.Clusters[clusterName]
		if !clusterUsedExist {
			clusterUsed = business.UsedQuantity{
				Used: make(business.ResourceList),
			}
		}
		allErrs = append(allErrs,
			resource.ValidateUpdateResource(clusterHard.Hard, clusterUsed.Used,
				field.NewPath("spec", "hard"))...)
	}

	return allErrs
}

func validateAgainstChildren(project *business.Project, getter validation.BusinessObjectGetter) (allErrs field.ErrorList) {
	clusterTimes := map[string]int{}          /* cluster -> times */
	quotaTimes := map[string]map[string]int{} /* cluster -> (kind -> times) */

	fldChildrenProjectsPath := field.NewPath("status").Child("calculatedChildProjects")
	for _, name := range project.Status.CalculatedChildProjects {
		childProject, err := getter.Project(name, metav1.GetOptions{})
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldChildrenProjectsPath, name,
				fmt.Sprintf("failed to get child project '%s', for %s", name, err)))
			continue
		}

		for clusterName, clusterHard := range childProject.Spec.Clusters {
			_, has := quotaTimes[clusterName]
			if !has {
				quotaTimes[clusterName] = map[string]int{}
			}
			clusterTimes[clusterName] = clusterTimes[clusterName] + 1
			for k := range clusterHard.Hard {
				quotaTimes[clusterName][k] = quotaTimes[clusterName][k] + 1
			}
		}
	}

	fldChildrenNamespacesPath := field.NewPath("status").Child("calculatedNamespaces")
	for _, name := range project.Status.CalculatedNamespaces {
		childNamespace, err := getter.Namespace(project.Name, name, metav1.GetOptions{})
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldChildrenNamespacesPath, name,
				fmt.Sprintf("failed to get child namespace '%s', for %s", name, err)))
			continue
		}

		clusterName := childNamespace.Spec.ClusterName
		_, has := quotaTimes[clusterName]
		if !has {
			quotaTimes[clusterName] = map[string]int{}
		}
		clusterTimes[clusterName] = clusterTimes[clusterName] + 1
		for k := range childNamespace.Spec.Hard {
			quotaTimes[clusterName][k] = quotaTimes[clusterName][k] + 1
		}
	}

	fldHardPath := field.NewPath("spec").Child("hard")
	for clusterName, clusterHard := range project.Spec.Clusters {
		clusterQuota, has := quotaTimes[clusterName]
		if !has {
			continue /* skip clusters that have NOT been used by children */
		}
		for k := range clusterHard.Hard {
			if clusterQuota[k] < clusterTimes[clusterName] {
				allErrs = append(allErrs,
					field.Invalid(fldHardPath.Child(clusterName), k,
						fmt.Sprintf("%s '%s' of cluster '%s' for all children projects/namespaces", _addNewQuotaErrorInfo, k, clusterName)))
			}

		}
	}

	return allErrs
}

func validateAgainstParent(project *business.Project, old *business.Project, getter validation.BusinessObjectGetter) (allErrs field.ErrorList) {
	fldSpecPath := field.NewPath("spec")
	fldHardPath := fldSpecPath.Child("hard")
	fldParentProjectPath := fldSpecPath.Child("parentProjectName")

	parentProject, err := getter.Project(project.Spec.ParentProjectName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return append(allErrs, field.NotFound(fldParentProjectPath, project.Spec.ParentProjectName))
		}
		return append(allErrs, field.InternalError(fldParentProjectPath, err))
	}
	if parentProject.Spec.TenantID != project.Spec.TenantID {
		return append(allErrs, field.NotFound(fldParentProjectPath, project.Spec.ParentProjectName))
	}
	if parentProject.Status.Locked != nil && *parentProject.Status.Locked {
		return append(allErrs,
			field.Invalid(fldParentProjectPath,
				project.Spec.ParentProjectName,
				"parent project has been locked"))
	}
	for clusterName, needed := range project.Spec.Clusters {
		clusterHard, clusterHardExist := parentProject.Spec.Clusters[clusterName]
		if !clusterHardExist {
			allErrs = append(allErrs,
				field.Invalid(fldHardPath, clusterName,
					fmt.Sprintf("%s %s", _clusterLimitErrorInfo, clusterName)))
			continue
		}
		clusterUsed, clusterUsedExist := parentProject.Status.Clusters[clusterName]
		if !clusterUsedExist {
			clusterUsed = business.UsedQuantity{
				Used: make(business.ResourceList),
			}
		}

		var oldNeededHard business.ResourceList
		childProjectNameSet := sets.NewString(parentProject.Status.CalculatedChildProjects...)
		if childProjectNameSet.Has(project.ObjectMeta.Name) && old != nil {
			if oldNeed, hasOld := old.Spec.Clusters[clusterName]; hasOld {
				oldNeededHard = oldNeed.Hard
			}
		}
		allErrs = append(allErrs,
			resource.ValidateAllocatableResources(needed.Hard, oldNeededHard,
				clusterHard.Hard, clusterUsed.Used, fldHardPath)...)

	}

	return allErrs
}
