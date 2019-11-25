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
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/trace"
	"strings"
	"time"
	"tkestack.io/tke/pkg/platform/controller/addon/gpumanager/template"
	"tkestack.io/tke/pkg/util/log"
)

const (
	// DefaultServiceAccountName defines the service account and role-binding name of GPUManager
	DefaultServiceAccountName = "gpu-manager"
)

// ServiceAccountOperator is used to operate remote service account and role-binding
type ServiceAccountOperator struct {
	K8sCli kubernetes.Interface
}

// CreateServiceAccount create a service account and role-binding of clusterName
func (sa *ServiceAccountOperator) CreateServiceAccount(clusterName string) error {
	t := trace.New(fmt.Sprintf("CreateServiceAccount-%s", clusterName))
	defer t.LogIfLong(time.Second)

	reader := strings.NewReader(template.ServiceAccountTemplate)

	serviceAccountPayload := &corev1.ServiceAccount{}
	err := yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(serviceAccountPayload)
	if err != nil {
		return err
	}
	serviceAccountPayload.Name = DefaultServiceAccountName

	t.Step("create service account")
	_, err = sa.K8sCli.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(serviceAccountPayload)
	if err != nil {
		log.Error(fmt.Sprintf("can't create service account for %s, got %v", clusterName, err))
		return err
	}

	reader = strings.NewReader(template.RoleTemplate)
	rolePayload := &rbacv1.ClusterRole{}
	err = yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(rolePayload)
	if err != nil {
		return err
	}
	rolePayload.Name = DefaultServiceAccountName

	t.Step("create role")
	_, err = sa.K8sCli.RbacV1().ClusterRoles().Create(rolePayload)

	if err != nil {
		log.Error(fmt.Sprintf("can't create role for %s, got %v", clusterName, err))
		return err
	}

	reader = strings.NewReader(template.RoleBindingTemplate)
	roleBindingPayload := &rbacv1.ClusterRoleBinding{}
	err = yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(roleBindingPayload)
	if err != nil {
		return err
	}
	roleBindingPayload.Name = DefaultServiceAccountName

	t.Step("create role binding")
	_, err = sa.K8sCli.RbacV1().ClusterRoleBindings().Create(roleBindingPayload)

	if err != nil {
		log.Error(fmt.Sprintf("can't create role binding for %s, got %v", clusterName, err))
		return err
	}

	return nil
}

// DeleteServiceAccount delete a service account and role-binding of clusterName
func (sa *ServiceAccountOperator) DeleteServiceAccount(clusterName string) error {
	t := trace.New(fmt.Sprintf("DeleteServiceAccount-%s", clusterName))
	defer t.LogIfLong(time.Second)

	t.Step("delete role binding")
	err := sa.K8sCli.RbacV1().ClusterRoleBindings().Delete(DefaultServiceAccountName, metav1.NewDeleteOptions(0))
	if err != nil && !apierror.IsNotFound(err) {
		log.Error(fmt.Sprintf("can't delete role binding for %s, got %v", clusterName, err))
		return err
	}

	t.Step("delete role")
	err = sa.K8sCli.RbacV1().ClusterRoles().Delete(DefaultServiceAccountName, metav1.NewDeleteOptions(0))
	if err != nil && !apierror.IsNotFound(err) {
		log.Error(fmt.Sprintf("can't delete role for %s, got %v", clusterName, err))
		return err
	}

	t.Step("delete service account")
	err = sa.K8sCli.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(DefaultServiceAccountName, metav1.NewDeleteOptions(0))
	if err != nil && !apierror.IsNotFound(err) {
		log.Error(fmt.Sprintf("can't delete service account for %s, got %v", clusterName, err))
		return err
	}

	return nil
}
