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
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/trace"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	v1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/controller/addon/gpumanager/template"
	"tkestack.io/tke/pkg/util/log"
)

const (
	// DefaultDaemonSetName defines the name of GPUManager daemonset
	DefaultDaemonSetName = "gpu-manager-daemonset"
	maxUpdateRetry       = 5
	// LabelNodeRoleMaster defines the label name of node role
	LabelNodeRoleMaster = "node-role.kubernetes.io/master"
)

// DaemonSetOperator is used to operate remote GPUManager daemonset
type DaemonSetOperator struct {
	TkeCli      clientset.Interface
	K8sCli      kubernetes.Interface
	CrDomain    string
	CrNamespace string
}

// CreatePayload returns the daemonset which holds GPUManager information
func CreatePayload(version string) (ret *appsv1.DaemonSet, err error) {
	reader := strings.NewReader(template.GPUManagerDaemonSetTemplate)

	// TODO: add OwnerReference for auto cleanup
	// We have to set image and a label which indicates the owner of the GPUManager
	payload := &appsv1.DaemonSet{}
	err = yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(payload)
	if err != nil {
		return nil, err
	}
	payload.Name = DefaultDaemonSetName
	payload.Spec.Template.Spec.Containers[0].Image = images.Get(version).GPUManager.FullName()
	payload.Spec.Template.Spec.ServiceAccountName = DefaultServiceAccountName
	payload.Spec.Template.Spec.Tolerations[0].Key = LabelNodeRoleMaster
	payload.Spec.Template.Spec.Tolerations[0].Effect = corev1.TaintEffectNoSchedule

	return payload, nil
}

// CreateDaemonSet implements DaemonSetOperator
func (c *DaemonSetOperator) CreateDaemonSet(obj *v1.GPUManager) error {
	t := trace.New(fmt.Sprintf("CreateDaemonSet-%s", obj.Spec.ClusterName))
	defer t.LogIfLong(time.Second)

	log.Info(fmt.Sprintf("Start to create daemonset for %s", obj.Spec.ClusterName))

	t.Step("create daemonset to remote")
	payload, err := CreatePayload(obj.Spec.Version)
	if err != nil {
		return err
	}
	_, err = c.K8sCli.AppsV1().DaemonSets(metav1.NamespaceSystem).Create(payload)
	if err != nil && !apierror.IsAlreadyExists(err) {
		log.Error(fmt.Sprintf("create daemonset for %s, got %s", obj.Spec.ClusterName, err))
	}

	return err
}

// DeleteDaemonSet implements DaemonSetOperator
func (c *DaemonSetOperator) DeleteDaemonSet(obj *v1.GPUManager) error {
	t := trace.New(fmt.Sprintf("DeleteDaemonSet-%s", obj.Spec.ClusterName))
	defer t.LogIfLong(time.Second)

	log.Info(fmt.Sprintf("Start to delete daemonset for %s", obj.Spec.ClusterName))
	t.Step("delete daemonset from remote")
	err := c.K8sCli.AppsV1().DaemonSets(metav1.NamespaceSystem).Delete(DefaultDaemonSetName, metav1.NewDeleteOptions(0))
	if err != nil && !apierror.IsNotFound(err) {
		log.Error(fmt.Sprintf("delete daemonset for %s, got %s", obj.Spec.ClusterName, err))
		return err
	}

	return nil
}

// GetDaemonSet implements DaemonSetOperator
func (c *DaemonSetOperator) GetDaemonSet(obj *v1.GPUManager) (ds *appsv1.DaemonSet, opErr error) {
	t := trace.New(fmt.Sprintf("GetDaemonSet-%s", obj.Spec.ClusterName))
	defer t.LogIfLong(time.Second)

	log.Info(fmt.Sprintf("Start to get daemonset for %s", obj.Spec.ClusterName))
	t.Step("get daemonset from remote")
	ds, opErr = c.K8sCli.AppsV1().DaemonSets(metav1.NamespaceSystem).Get(DefaultDaemonSetName, metav1.GetOptions{})
	if opErr != nil {
		log.Error(fmt.Sprintf("get daemonset for %s, got %s", obj.Spec.ClusterName, opErr))
	}

	return
}

// UpdateGPUManagerStatus implements Operator
func (c *DaemonSetOperator) UpdateGPUManagerStatus(obj *v1.GPUManager) error {
	t := trace.New(fmt.Sprintf("UpdateGPUManagerStatus-%s", obj.Spec.ClusterName))
	defer t.LogIfLong(maxUpdateRetry * time.Second)

	log.Info(fmt.Sprintf("Start to update status to %s for %s", obj.Status.Phase, obj.Spec.ClusterName))
	t.Step("update GPUManager status to remote")
	return wait.PollImmediate(time.Second, maxUpdateRetry*time.Second, func() (done bool, err error) {
		_, err = c.TkeCli.PlatformV1().GPUManagers().UpdateStatus(obj)
		if err == nil {
			return true, nil
		}

		if apierror.IsNotFound(err) {
			log.Error(fmt.Sprintf("GPUManager in cluster %s is not found, ignore this update", obj.Spec.ClusterName))
			return true, nil
		}

		if apierror.IsConflict(err) {
			return false, nil
		}

		log.Error(fmt.Sprintf("update status to %s for %s, got %s", obj.Status.Phase, obj.Spec.ClusterName, err))
		return false, err
	})
}

// UpdateDaemonSet implements DaemonSetOperator
func (c *DaemonSetOperator) UpdateDaemonSet(obj *v1.GPUManager) error {
	t := trace.New(fmt.Sprintf("UpdateDaemonSet-%s", obj.Spec.ClusterName))
	defer t.LogIfLong(time.Second)

	log.Info(fmt.Sprintf("Start to update GPUManager for %s", obj.Spec.ClusterName))

	t.Step("patch image version to remote")
	patch := fmt.Sprintf(template.ImagePatchTemplate, images.Get(obj.Spec.Version).GPUManager.FullName())
	patched, err := c.K8sCli.AppsV1().DaemonSets(metav1.NamespaceSystem).Patch(DefaultDaemonSetName, types.JSONPatchType, []byte(patch))
	if err != nil {
		log.Error(fmt.Sprintf("patch image version for %s, got %s", obj.Spec.ClusterName, err))
		return err
	}

	log.Info(fmt.Sprintf("GPUManager %s patched version: %s", obj.Spec.ClusterName, patched))
	return nil
}

// DiffDaemonSet implements DaemonSetOperator
func (c *DaemonSetOperator) DiffDaemonSet(obj *v1.GPUManager) (needUpgrade bool, err error) {
	t := trace.New(fmt.Sprintf("DiffDaemonSet-%s", obj.Spec.ClusterName))
	defer t.LogIfLong(time.Second)

	log.Info(fmt.Sprintf("Diff GPUManager for %s", obj.Spec.ClusterName))
	t.Step("get daemonset from remote")
	oldDs, err := c.GetDaemonSet(obj)
	if err != nil {
		log.Error(fmt.Sprintf("diff daemonset for %s, got %s", obj.Spec.ClusterName, err))
		return false, err
	}

	oldImage := oldDs.Spec.Template.Spec.Containers[0].Image
	newImage := images.Get(obj.Spec.Version).GPUManager.FullName()

	log.Info(fmt.Sprintf("GPUManager %s oldImage:%s, newImage: %s", obj.Spec.ClusterName, oldImage, newImage))
	return oldImage != newImage, nil
}
