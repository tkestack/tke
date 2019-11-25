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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	cliset "tkestack.io/tke/api/client/clientset/versioned"
	"tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/controller/addon/gpumanager/operator/common"
)

type fakeOperator struct {
	configmapOperator      *common.ConfigmapOperator
	daemonSetOperator      *common.DaemonSetOperator
	deploymentOperator     *common.DeploymentOperator
	serviceAccountOperator *common.ServiceAccountOperator
	serviceOperator        *common.ServiceOperator
	serviceMetricOperator  *common.ServiceMetricOperator
}

var _ ObjectOperator = &fakeOperator{}

// NewFakeObjectOperator returns a fake implements of ObjectOperator for test
func NewFakeObjectOperator(tkeClient cliset.Interface, k8sClient kubernetes.Interface) ObjectOperator {
	op := &fakeOperator{
		configmapOperator: &common.ConfigmapOperator{},
		daemonSetOperator: &common.DaemonSetOperator{
			TkeCli:   tkeClient,
			K8sCli:   k8sClient,
			CrDomain: "test.com",
		},
		deploymentOperator: &common.DeploymentOperator{
			TkeCli:   tkeClient,
			K8sCli:   k8sClient,
			CrDomain: "test.com",
		},
		serviceAccountOperator: &common.ServiceAccountOperator{},
		serviceOperator:        &common.ServiceOperator{},
		serviceMetricOperator:  &common.ServiceMetricOperator{},
	}

	return op
}

func (op *fakeOperator) checkClusterExistence(clusterName string) error {
	_, err := op.daemonSetOperator.TkeCli.PlatformV1().Clusters().Get(clusterName, metav1.GetOptions{})
	return err
}

// CreateDaemonSet implements DaemonSetOperator
func (op *fakeOperator) CreateDaemonSet(obj *v1.GPUManager) error {
	if err := op.checkClusterExistence(obj.Spec.ClusterName); err != nil {
		return err
	}
	return op.daemonSetOperator.CreateDaemonSet(obj)
}

// DeleteDaemonSet implements DaemonSetOperator
func (op *fakeOperator) DeleteDaemonSet(obj *v1.GPUManager) error {
	if err := op.checkClusterExistence(obj.Spec.ClusterName); err != nil {
		return err
	}
	return op.daemonSetOperator.DeleteDaemonSet(obj)
}

// GetDaemonSet implements DaemonSetOperator
func (op *fakeOperator) GetDaemonSet(obj *v1.GPUManager) (gm *appsv1.DaemonSet, opErr error) {
	if err := op.checkClusterExistence(obj.Spec.ClusterName); err != nil {
		return nil, err
	}
	return op.daemonSetOperator.GetDaemonSet(obj)
}

// UpdateGPUManagerStatus implements ObjectOperator
func (op *fakeOperator) UpdateGPUManagerStatus(status *v1.GPUManager) error {
	return op.daemonSetOperator.UpdateGPUManagerStatus(status)
}

// UpdateDaemonSet implements DaemonSetOperator
func (op *fakeOperator) UpdateDaemonSet(obj *v1.GPUManager) error {
	if err := op.checkClusterExistence(obj.Spec.ClusterName); err != nil {
		return err
	}
	return op.daemonSetOperator.UpdateDaemonSet(obj)
}

// DiffDaemonSet implements DaemonSetOperator
func (op *fakeOperator) DiffDaemonSet(obj *v1.GPUManager) (changed bool, err error) {
	if err := op.checkClusterExistence(obj.Spec.ClusterName); err != nil {
		return false, err
	}
	return op.daemonSetOperator.DiffDaemonSet(obj)
}

// CreateServiceAccount implements ServiceAccountOperator
func (op *fakeOperator) CreateServiceAccount(clusterName string) error {
	if err := op.checkClusterExistence(clusterName); err != nil {
		return err
	}
	return op.serviceAccountOperator.CreateServiceAccount(clusterName)
}

// DeleteServiceAccount implements ServiceAccountOperator
func (op *fakeOperator) DeleteServiceAccount(clusterName string) error {
	if err := op.checkClusterExistence(clusterName); err != nil {
		return err
	}
	return op.serviceAccountOperator.DeleteServiceAccount(clusterName)
}

// CreateConfigmap implements ConfigmapOperator
func (op *fakeOperator) CreateConfigmap(clusterName string) error {
	if err := op.checkClusterExistence(clusterName); err != nil {
		return err
	}
	return op.configmapOperator.CreateConfigmap(clusterName)
}

// DeleteConfigmap implements ConfigmapOperator
func (op *fakeOperator) DeleteConfigmap(clusterName string) error {
	if err := op.checkClusterExistence(clusterName); err != nil {
		return err
	}
	return op.configmapOperator.DeleteConfigmap(clusterName)
}

// CreateService implements ServiceOperator
func (op *fakeOperator) CreateService(clusterName string) error {
	if err := op.checkClusterExistence(clusterName); err != nil {
		return err
	}
	return op.serviceOperator.CreateService(clusterName)
}

// DeleteService implements ServiceOperator
func (op *fakeOperator) DeleteService(clusterName string) error {
	if err := op.checkClusterExistence(clusterName); err != nil {
		return err
	}
	return op.serviceOperator.DeleteService(clusterName)
}

// CreateServiceMetric implements ServiceMetricOperator
func (op *fakeOperator) CreateServiceMetric(clusterName string) error {
	if err := op.checkClusterExistence(clusterName); err != nil {
		return err
	}
	return op.serviceMetricOperator.CreateServiceMetric(clusterName)
}

// DeleteServiceMetric implements ServiceMetricOperator
func (op *fakeOperator) DeleteServiceMetric(clusterName string) error {
	if err := op.checkClusterExistence(clusterName); err != nil {
		return err
	}
	return op.serviceMetricOperator.DeleteServiceMetric(clusterName)
}

// CreateDeployment implements DeploymentOperator
func (op *fakeOperator) CreateDeployment(obj *v1.GPUManager) error {
	if err := op.checkClusterExistence(obj.Spec.ClusterName); err != nil {
		return err
	}
	return op.deploymentOperator.CreateDeployment(obj)
}

// DeleteDeployment implements DeploymentOperator
func (op *fakeOperator) DeleteDeployment(obj *v1.GPUManager) error {
	if err := op.checkClusterExistence(obj.Spec.ClusterName); err != nil {
		return err
	}
	return op.deploymentOperator.DeleteDeployment(obj)
}

// GetDeployment implements DeploymentOperator
func (op *fakeOperator) GetDeployment(obj *v1.GPUManager) (gm *appsv1.Deployment, opErr error) {
	if err := op.checkClusterExistence(obj.Spec.ClusterName); err != nil {
		return nil, err
	}
	return op.deploymentOperator.GetDeployment(obj)
}

// UpdateDeployment implements DeploymentOperator
func (op *fakeOperator) UpdateDeployment(obj *v1.GPUManager) error {
	if err := op.checkClusterExistence(obj.Spec.ClusterName); err != nil {
		return err
	}
	return op.deploymentOperator.UpdateDeployment(obj)
}

// DiffDeployment implements DeploymentOperator
func (op *fakeOperator) DiffDeployment(obj *v1.GPUManager) (changed bool, err error) {
	if err := op.checkClusterExistence(obj.Spec.ClusterName); err != nil {
		return false, err
	}
	return op.deploymentOperator.DiffDeployment(obj)
}
