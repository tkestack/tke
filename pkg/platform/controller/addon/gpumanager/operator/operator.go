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
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	v1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/controller/addon/gpumanager/operator/common"
	"tkestack.io/tke/pkg/platform/util"
)

type k8sOperator struct {
	configmapOperator      *common.ConfigmapOperator
	daemonSetOperator      *common.DaemonSetOperator
	deploymentOperator     *common.DeploymentOperator
	serviceAccountOperator *common.ServiceAccountOperator
	serviceOperator        *common.ServiceOperator
	serviceMetricOperator  *common.ServiceMetricOperator
}

var _ ObjectOperator = &k8sOperator{}

// NewObjectOperator returns a real-world instance of ObjectOperator
func NewObjectOperator(TkeCli clientset.Interface) ObjectOperator {
	op := &k8sOperator{
		configmapOperator: &common.ConfigmapOperator{},
		daemonSetOperator: &common.DaemonSetOperator{
			TkeCli: TkeCli,
		},
		deploymentOperator: &common.DeploymentOperator{
			TkeCli: TkeCli,
		},
		serviceAccountOperator: &common.ServiceAccountOperator{},
		serviceOperator:        &common.ServiceOperator{},
		serviceMetricOperator:  &common.ServiceMetricOperator{},
	}

	return op
}

func (op *k8sOperator) getClusterClient(clusterName string) (*kubernetes.Clientset, error) {
	cluster, err := op.daemonSetOperator.TkeCli.PlatformV1().Clusters().Get(clusterName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return util.BuildExternalClientSet(cluster, op.daemonSetOperator.TkeCli.PlatformV1())
}

// CreateDaemonSet implements DaemonSetOperator
func (op *k8sOperator) CreateDaemonSet(created *v1.GPUManager) error {
	kubeClient, err := op.getClusterClient(created.Spec.ClusterName)
	if err != nil {
		return err
	}

	op.daemonSetOperator.K8sCli = kubeClient
	return op.daemonSetOperator.CreateDaemonSet(created)
}

// DeleteDaemonSet implements DaemonSetOperator
func (op *k8sOperator) DeleteDaemonSet(deleted *v1.GPUManager) error {
	kubeClient, err := op.getClusterClient(deleted.Spec.ClusterName)
	if err != nil && !apierror.IsNotFound(err) {
		return err
	}

	if apierror.IsNotFound(err) {
		return nil
	}

	op.daemonSetOperator.K8sCli = kubeClient
	return op.daemonSetOperator.DeleteDaemonSet(deleted)
}

// GetDaemonSet implements DaemonSetOperator
func (op *k8sOperator) GetDaemonSet(get *v1.GPUManager) (gm *appsv1.DaemonSet, opErr error) {
	kubeClient, err := op.getClusterClient(get.Spec.ClusterName)
	if err != nil {
		return nil, err
	}

	op.daemonSetOperator.K8sCli = kubeClient
	return op.daemonSetOperator.GetDaemonSet(get)
}

// UpdateGPUManagerStatus implements ObjectOperator
func (op *k8sOperator) UpdateGPUManagerStatus(status *v1.GPUManager) error {
	return op.daemonSetOperator.UpdateGPUManagerStatus(status)
}

// UpdateDaemonSet implements DaemonSetOperator
func (op *k8sOperator) UpdateDaemonSet(updated *v1.GPUManager) error {
	kubeClient, err := op.getClusterClient(updated.Spec.ClusterName)
	if err != nil {
		return err
	}

	op.daemonSetOperator.K8sCli = kubeClient
	return op.daemonSetOperator.UpdateDaemonSet(updated)
}

// DiffDaemonSet implements DaemonSetOperator
func (op *k8sOperator) DiffDaemonSet(diff *v1.GPUManager) (changed bool, err error) {
	kubeClient, err := op.getClusterClient(diff.Spec.ClusterName)
	if err != nil {
		return false, err
	}

	op.daemonSetOperator.K8sCli = kubeClient
	return op.daemonSetOperator.DiffDaemonSet(diff)
}

// CreateServiceAccount implements ServiceAccountOperator
func (op *k8sOperator) CreateServiceAccount(clusterName string) error {
	kubeClient, err := op.getClusterClient(clusterName)
	if err != nil {
		return err
	}

	op.serviceAccountOperator.K8sCli = kubeClient
	return op.serviceAccountOperator.CreateServiceAccount(clusterName)
}

// DeleteServiceAccount implements ServiceAccountOperator
func (op *k8sOperator) DeleteServiceAccount(clusterName string) error {
	kubeClient, err := op.getClusterClient(clusterName)
	if err != nil {
		return err
	}

	op.serviceAccountOperator.K8sCli = kubeClient
	return op.serviceAccountOperator.DeleteServiceAccount(clusterName)
}

// CreateConfigmap implements ConfigmapOperator
func (op *k8sOperator) CreateConfigmap(clusterName string) error {
	kubeClient, err := op.getClusterClient(clusterName)
	if err != nil {
		return err
	}

	op.configmapOperator.K8sCli = kubeClient
	return op.configmapOperator.CreateConfigmap(clusterName)
}

// DeleteConfigmap implements ConfigmapOperator
func (op *k8sOperator) DeleteConfigmap(clusterName string) error {
	kubeClient, err := op.getClusterClient(clusterName)
	if err != nil {
		return err
	}

	op.configmapOperator.K8sCli = kubeClient
	return op.configmapOperator.DeleteConfigmap(clusterName)
}

// CreateService implements ServiceOperator
func (op *k8sOperator) CreateService(clusterName string) error {
	kubeClient, err := op.getClusterClient(clusterName)
	if err != nil {
		return err
	}

	op.serviceOperator.K8sCli = kubeClient
	return op.serviceOperator.CreateService(clusterName)
}

// DeleteService implements ServiceOperator
func (op *k8sOperator) DeleteService(clusterName string) error {
	kubeClient, err := op.getClusterClient(clusterName)
	if err != nil {
		return err
	}

	op.serviceOperator.K8sCli = kubeClient
	return op.serviceOperator.DeleteService(clusterName)
}

// CreateServiceMetric implements ServiceMetricOperator
func (op *k8sOperator) CreateServiceMetric(clusterName string) error {
	kubeClient, err := op.getClusterClient(clusterName)
	if err != nil {
		return err
	}

	op.serviceMetricOperator.K8sCli = kubeClient
	return op.serviceMetricOperator.CreateServiceMetric(clusterName)
}

// DeleteServiceMetric implements ServiceMetricOperator
func (op *k8sOperator) DeleteServiceMetric(clusterName string) error {
	kubeClient, err := op.getClusterClient(clusterName)
	if err != nil {
		return err
	}

	op.serviceMetricOperator.K8sCli = kubeClient
	return op.serviceMetricOperator.DeleteServiceMetric(clusterName)
}

// CreateDeployment implements DeploymentOperator
func (op *k8sOperator) CreateDeployment(created *v1.GPUManager) error {
	kubeClient, err := op.getClusterClient(created.Spec.ClusterName)
	if err != nil {
		return err
	}

	op.deploymentOperator.K8sCli = kubeClient
	return op.deploymentOperator.CreateDeployment(created)
}

// DeleteDeployment implements DeploymentOperator
func (op *k8sOperator) DeleteDeployment(deleted *v1.GPUManager) error {
	kubeClient, err := op.getClusterClient(deleted.Spec.ClusterName)
	if err != nil && !apierror.IsNotFound(err) {
		return err
	}

	if apierror.IsNotFound(err) {
		return nil
	}

	op.deploymentOperator.K8sCli = kubeClient
	return op.deploymentOperator.DeleteDeployment(deleted)
}

// GetDeployment implements DeploymentOperator
func (op *k8sOperator) GetDeployment(get *v1.GPUManager) (gm *appsv1.Deployment, opErr error) {
	kubeClient, err := op.getClusterClient(get.Spec.ClusterName)
	if err != nil {
		return nil, err
	}

	op.deploymentOperator.K8sCli = kubeClient
	return op.deploymentOperator.GetDeployment(get)
}

// UpdateDeployment implements DeploymentOperator
func (op *k8sOperator) UpdateDeployment(updated *v1.GPUManager) error {
	kubeClient, err := op.getClusterClient(updated.Spec.ClusterName)
	if err != nil {
		return err
	}

	op.deploymentOperator.K8sCli = kubeClient
	return op.deploymentOperator.UpdateDeployment(updated)
}

// DiffDeployment implements DeploymentOperator
func (op *k8sOperator) DiffDeployment(diff *v1.GPUManager) (changed bool, err error) {
	kubeClient, err := op.getClusterClient(diff.Spec.ClusterName)
	if err != nil {
		return false, err
	}

	op.deploymentOperator.K8sCli = kubeClient
	return op.deploymentOperator.DiffDeployment(diff)
}
