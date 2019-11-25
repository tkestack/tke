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

package common

import (
	"fmt"
	"strings"
	"time"

	"tkestack.io/tke/pkg/platform/controller/addon/gpumanager/images"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/trace"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	v1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/controller/addon/gpumanager/template"
	"tkestack.io/tke/pkg/util/log"
)

const (
	// DefaultDeploymentName defines the name of GPUManager deployment
	DefaultDeploymentName = "gpu-quota-admission"
)

// DeploymentOperator is used to operate remote GPUManager deployment
type DeploymentOperator struct {
	TkeCli      clientset.Interface
	K8sCli      kubernetes.Interface
	CrDomain    string
	CrNamespace string
}

// CreateDeploymentPayload returns the deployment which holds GPUManager information
func CreateDeploymentPayload(version string) (ret *appsv1.Deployment, err error) {
	reader := strings.NewReader(template.DeploymentTemplate)

	// TODO: add OwnerReference for auto cleanup
	// We have to set image and a label which indicates the owner of the GPUManager
	payload := &appsv1.Deployment{}
	err = yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(payload)
	if err != nil {
		return nil, err
	}
	payload.Name = DefaultDeploymentName
	payload.Spec.Template.Spec.InitContainers[0].Image = images.Get(version).Busybox.FullName()
	payload.Spec.Template.Spec.Containers[0].Image = images.Get(version).GPUQuotaAdmission.FullName()
	payload.Spec.Template.Spec.ServiceAccountName = DefaultServiceAccountName
	payload.Spec.Template.Spec.Tolerations[0].Key = LabelNodeRoleMaster
	payload.Spec.Template.Spec.Tolerations[0].Effect = corev1.TaintEffectNoSchedule

	return payload, nil
}

// CreateDeployment implements DeploymentOperator
func (c *DeploymentOperator) CreateDeployment(obj *v1.GPUManager) error {
	t := trace.New(fmt.Sprintf("CreateDeployment-%s", obj.Spec.ClusterName))
	defer t.LogIfLong(time.Second)

	log.Info(fmt.Sprintf("Start to create deployment for %s", obj.Spec.ClusterName))

	t.Step("create deployment to remote")
	payload, err := CreateDeploymentPayload(obj.Spec.Version)
	if err != nil {
		return err
	}
	_, err = c.K8sCli.AppsV1().Deployments(metav1.NamespaceSystem).Create(payload)
	if err != nil && !apierror.IsAlreadyExists(err) {
		log.Error(fmt.Sprintf("create deployment for %s, got %s", obj.Spec.ClusterName, err))
	}

	return err
}

// DeleteDeployment implements DeploymentOperator
func (c *DeploymentOperator) DeleteDeployment(obj *v1.GPUManager) error {
	t := trace.New(fmt.Sprintf("DeleteDeployment-%s", obj.Spec.ClusterName))
	defer t.LogIfLong(time.Second)

	log.Info(fmt.Sprintf("Start to delete deployment for %s", obj.Spec.ClusterName))
	t.Step("delete deployment from remote")
	err := c.K8sCli.AppsV1().Deployments(metav1.NamespaceSystem).Delete(DefaultDeploymentName, metav1.NewDeleteOptions(0))
	if err != nil && !apierror.IsNotFound(err) {
		log.Error(fmt.Sprintf("delete deployment for %s, got %s", obj.Spec.ClusterName, err))
		return err
	}

	return nil
}

// GetDeployment implements DeploymentOperator
func (c *DeploymentOperator) GetDeployment(obj *v1.GPUManager) (ds *appsv1.Deployment, opErr error) {
	t := trace.New(fmt.Sprintf("GetDeployment-%s", obj.Spec.ClusterName))
	defer t.LogIfLong(time.Second)

	log.Info(fmt.Sprintf("Start to get deployment for %s", obj.Spec.ClusterName))
	t.Step("get deployment from remote")
	ds, opErr = c.K8sCli.AppsV1().Deployments(metav1.NamespaceSystem).Get(DefaultDeploymentName, metav1.GetOptions{})
	if opErr != nil {
		log.Error(fmt.Sprintf("get deployment for %s, got %s", obj.Spec.ClusterName, opErr))
	}

	return
}

// UpdateDeployment implements DeploymentOperator
func (c *DeploymentOperator) UpdateDeployment(obj *v1.GPUManager) error {
	t := trace.New(fmt.Sprintf("UpdateDeployment-%s", obj.Spec.ClusterName))
	defer t.LogIfLong(time.Second)

	log.Info(fmt.Sprintf("Start to update quota admission for %s", obj.Spec.ClusterName))

	t.Step("patch admission image version to remote")
	patch := fmt.Sprintf(template.ImagePatchTemplate, images.Get(obj.Spec.Version).GPUQuotaAdmission.FullName())
	patched, err := c.K8sCli.AppsV1().Deployments(metav1.NamespaceSystem).Patch(DefaultDeploymentName, types.JSONPatchType, []byte(patch))
	if err != nil {
		log.Error(fmt.Sprintf("patch image version for %s, got %s", obj.Spec.ClusterName, err))
		return err
	}

	log.Info(fmt.Sprintf("GPUManager %s patched version: %s", obj.Spec.ClusterName, patched))
	return nil
}

// DiffDeployment implements DeploymentOperator
func (c *DeploymentOperator) DiffDeployment(obj *v1.GPUManager) (needUpgrade bool, err error) {
	t := trace.New(fmt.Sprintf("DiffDeployment-%s", obj.Spec.ClusterName))
	defer t.LogIfLong(time.Second)

	log.Info(fmt.Sprintf("Diff GPUManager for %s", obj.Spec.ClusterName))
	t.Step("get deployment from remote")
	oldDs, err := c.GetDeployment(obj)
	if err != nil {
		log.Error(fmt.Sprintf("diff deployment for %s, got %s", obj.Spec.ClusterName, err))
		return false, err
	}

	oldImage := oldDs.Spec.Template.Spec.Containers[0].Image
	newImage := images.Get(obj.Spec.Version).GPUQuotaAdmission.FullName()

	log.Info(fmt.Sprintf("GPUManager %s oldImage:%s, newImage: %s", obj.Spec.ClusterName, oldImage, newImage))
	return oldImage != newImage, nil
}
