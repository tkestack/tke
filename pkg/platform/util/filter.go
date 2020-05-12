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
	"context"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/authentication"
)

// FilterCluster is used to filter clusters that do not belong to the tenant.
func FilterCluster(ctx context.Context, cluster *platform.Cluster) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if cluster.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("cluster"), cluster.ObjectMeta.Name)
	}
	return nil
}

// FilterClusterCredential is used to filter ClusterCredential that do not belong to the tenant.
func FilterClusterCredential(ctx context.Context, obj *platform.ClusterCredential) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if obj.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("ClusterCredential"), obj.ObjectMeta.Name)
	}
	return nil
}

// FilterMachine is used to filter machine that do not belong to the tenant.
func FilterMachine(ctx context.Context, machine *platform.Machine) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if machine.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("machine"), machine.ObjectMeta.Name)
	}
	return nil
}

// FilterRegistry is used to filter registry that do not belong to the tenant.
func FilterRegistry(ctx context.Context, registry *platform.Registry) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if registry.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("registry"), registry.ObjectMeta.Name)
	}
	return nil
}

// FilterPersistentEvent is used to filter persistent event that do not belong
// to the tenant.
func FilterPersistentEvent(ctx context.Context, pe *platform.PersistentEvent) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if pe.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("persistentevent"), pe.ObjectMeta.Name)
	}
	return nil
}

// FilterHelm is used to filter helm that do not belong
// to the tenant.
func FilterHelm(ctx context.Context, helm *platform.Helm) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if helm.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("helm"), helm.ObjectMeta.Name)
	}
	return nil
}

// FilterTappController is used to filter tapp controller that do not belong
// to the tenant.
func FilterTappController(ctx context.Context, tappController *platform.TappController) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if tappController.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("tappcontroller"), tappController.ObjectMeta.Name)
	}
	return nil
}

// FilterCSIOperator is used to filter csi operator that do not belong
// to the tenant.
func FilterCSIOperator(ctx context.Context, csiOperator *platform.CSIOperator) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if csiOperator.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("csioperator"), csiOperator.ObjectMeta.Name)
	}
	return nil
}

// FilterVolumeDecorator is used to filter volume decorator that do not belong
// to the tenant.
func FilterVolumeDecorator(ctx context.Context, decorator *platform.VolumeDecorator) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if decorator.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("volumedecorator"), decorator.ObjectMeta.Name)
	}
	return nil
}

// FilterLogCollector is used to filter log collector that do not belong
// to the tenant.
func FilterLogCollector(ctx context.Context, decorator *platform.LogCollector) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if decorator.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("logcollector"), decorator.ObjectMeta.Name)
	}
	return nil
}

// FilterCronHPA is used to filter CronHPA that do not belong
// to the tenant.
func FilterCronHPA(ctx context.Context, cronHPA *platform.CronHPA) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if cronHPA.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("cronhpa"), cronHPA.ObjectMeta.Name)
	}
	return nil
}

// FilterPrometheus is used to filter helm that do not belong
// to the tenant.
func FilterPrometheus(ctx context.Context, prom *platform.Prometheus) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if prom.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("prometheus"), prom.ObjectMeta.Name)
	}
	return nil
}

// FilterIPAM is used to filter ipam that do not belong
// to the tenant.
func FilterIPAM(ctx context.Context, ipam *platform.IPAM) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if ipam.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("ipam"), ipam.ObjectMeta.Name)
	}
	return nil
}

// FilterLBCF is used to filter LBCF that do not belong to the tenant.
func FilterLBCF(ctx context.Context, lbcf *platform.LBCF) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if lbcf.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("lbcf"), lbcf.ObjectMeta.Name)
	}
	return nil
}
