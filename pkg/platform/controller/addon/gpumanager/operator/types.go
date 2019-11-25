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

package operator

import (
	appsv1 "k8s.io/api/apps/v1"
	"tkestack.io/tke/api/platform/v1"
)

// ObjectOperator defines a series of methods to manager GPUManager and its components
type ObjectOperator interface {
	ConfigmapOperator
	DaemonSetOperator
	DeploymentOperator
	ServiceAccountOperator
	ServiceOperator
	ServiceMetricOperator
	// UpdateGPUManagerStatus means update GPUManager status to the given object
	UpdateGPUManagerStatus(status *v1.GPUManager) error
}

// DaemonSetOperator defines a series of methods to operate daemonset which belongs to GPUManager
type DaemonSetOperator interface {
	// CreateDaemonSet means create a GPUManager daemonset to specified cluster
	CreateDaemonSet(created *v1.GPUManager) error
	// DeleteDaemonSet means remove a GPUManager daemonset from specified cluster
	DeleteDaemonSet(deleted *v1.GPUManager) error
	// GetDaemonSet means retrieve the GPUManager daemonset from specified cluster
	GetDaemonSet(get *v1.GPUManager) (*appsv1.DaemonSet, error)
	// DiffDaemonSet returns if we need to upgrade daemonset
	DiffDaemonSet(diff *v1.GPUManager) (bool, error)
	// UpdateDaemonSet means the GPUManager daemonset in the cluster need a perform change
	UpdateDaemonSet(updated *v1.GPUManager) error
}

// DeploymentOperator defines a series of methods to operate daemonset which belongs to GPUManager
type DeploymentOperator interface {
	// CreateDeployment means create a GPUManager daemonset to specified cluster
	CreateDeployment(created *v1.GPUManager) error
	// DeleteDeployment means remove a GPUManager daemonset from specified cluster
	DeleteDeployment(deleted *v1.GPUManager) error
	// GetDeployment means retrieve the GPUManager daemonset from specified cluster
	GetDeployment(get *v1.GPUManager) (*appsv1.Deployment, error)
	// DiffDeployment returns if we need to upgrade daemonset
	DiffDeployment(diff *v1.GPUManager) (bool, error)
	// UpdateDeployment means the GPUManager daemonset in the cluster need a perform change
	UpdateDeployment(updated *v1.GPUManager) error
}

// ServiceAccountOperator defines a series of methods to operate service account and role-binding which belongs to GPUManager
type ServiceAccountOperator interface {
	// CreateServiceAccount means create a service account and role-binding of GPUManager
	CreateServiceAccount(clusterName string) error
	// DeleteServiceAccount means delete a service account and role-binding of GPUManager
	DeleteServiceAccount(clusterName string) error
}

// ServiceOperator defines a series of methods to operate service which belongs to GPUManager
type ServiceOperator interface {
	// CreateService means create a service for metric of GPUManager
	CreateService(clusterName string) error
	// DeleteService means delete a service for metric of GPUManager
	DeleteService(clusterName string) error
}

// ServiceMetricOperator defines a series of methods to operate service which belongs to GPUManager
type ServiceMetricOperator interface {
	// CreateServiceMetric means create a service for metric of GPUManager
	CreateServiceMetric(clusterName string) error
	// DeleteServiceMetric means delete a service for metric of GPUManager
	DeleteServiceMetric(clusterName string) error
}

// ConfigmapOperator defines a series of methods to operate the configmap
type ConfigmapOperator interface {
	// CreateConfigmap means create a service for metric of GPUManager
	CreateConfigmap(clusterName string) error
	// DeleteConfigmap means delete a service for metric of GPUManager
	DeleteConfigmap(clusterName string) error
}
