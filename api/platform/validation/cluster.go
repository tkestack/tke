/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package validation

import (
	"context"

	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/platform/types"
	utilvalidation "tkestack.io/tke/pkg/util/validation"
)

// ValidateCluster validates a given Cluster.
func ValidateCluster(ctx context.Context, cluster *types.Cluster) field.ErrorList {
	fldPath := field.NewPath("spec")
	allErrs := apimachineryvalidation.ValidateObjectMeta(&cluster.ObjectMeta, false, apimachineryvalidation.NameIsDNSLabel, field.NewPath("metadata"))

	allErrs = append(allErrs, ValidateClusteType(cluster.Spec.Type, fldPath.Child("type"))...)
	p, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		allErrs = append(allErrs, field.NotFound(fldPath, cluster.Spec.Type))
	}

	allErrs = append(allErrs, p.Validate(ctx, cluster)...)

	return allErrs
}

// ValidateClusterUpdate tests if an update to a cluster is valid.
func ValidateClusterUpdate(ctx context.Context, cluster *types.Cluster, oldCluster *types.Cluster) field.ErrorList {
	fldPath := field.NewPath("spec")

	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&cluster.ObjectMeta, &oldCluster.ObjectMeta, field.NewPath("metadata"))

	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.Type, oldCluster.Spec.Type, fldPath.Child("type"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.NetworkDevice, oldCluster.Spec.NetworkDevice, fldPath.Child("networkDevice"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.ClusterCIDR, oldCluster.Spec.ClusterCIDR, fldPath.Child("clusterCIDR"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.DNSDomain, oldCluster.Spec.DNSDomain, fldPath.Child("dnsDomain"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.DockerExtraArgs, oldCluster.Spec.DockerExtraArgs, fldPath.Child("dockerExtraArgs"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.KubeletExtraArgs, oldCluster.Spec.KubeletExtraArgs, fldPath.Child("kubeletExtraArgs"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.APIServerExtraArgs, oldCluster.Spec.APIServerExtraArgs, fldPath.Child("apiServerExtraArgs"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.ControllerManagerExtraArgs, oldCluster.Spec.ControllerManagerExtraArgs, fldPath.Child("controllerManagerExtraArgs"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.SchedulerExtraArgs, oldCluster.Spec.SchedulerExtraArgs, fldPath.Child("schedulerExtraArgs"))...)

	allErrs = append(allErrs, ValidateClusteType(cluster.Spec.Type, fldPath.Child("type"))...)
	p, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		allErrs = append(allErrs, field.NotFound(fldPath, cluster.Spec.Type))
	}

	allErrs = append(allErrs, p.ValidateUpdate(ctx, cluster, oldCluster)...)

	return allErrs
}

// ValidateClusteType validates a given type.
func ValidateClusteType(clusterType string, fldPath *field.Path) field.ErrorList {
	return utilvalidation.ValidateEnum(clusterType, fldPath, clusterprovider.Providers())
}
