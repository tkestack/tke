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

	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/trace"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	"tkestack.io/tke/pkg/platform/controller/addon/gpumanager/template"
	"tkestack.io/tke/pkg/util/log"
)

const (
	// DefaultServiceName defines the name of service for the metric service of GPUManager
	DefaultServiceName = "gpu-quota-admission"
)

// ServiceOperator is used to operate remote service account and role-binding
type ServiceOperator struct {
	TkeCli clientset.Interface
	K8sCli kubernetes.Interface
}

// CreateService create a service to help collecting metrics of GPUManager
func (s *ServiceOperator) CreateService(clusterName string) error {
	t := trace.New(fmt.Sprintf("CreateService-%s", clusterName))
	defer t.LogIfLong(time.Second)

	reader := strings.NewReader(template.ServiceTemplate)

	servicePayload := &corev1.Service{}
	err := yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(servicePayload)
	if err != nil {
		return err
	}
	servicePayload.Name = DefaultServiceName
	cluster, err := s.TkeCli.PlatformV1().Clusters().Get(clusterName, metav1.GetOptions{})
	if err == nil {
		if ip, ok := cluster.Annotations[constants.GPUQuotaAdmissionIPAnnotaion]; ok {
			servicePayload.Spec.ClusterIP = ip
		}
	}
	t.Step("create service")
	_, err = s.K8sCli.CoreV1().Services(metav1.NamespaceSystem).Create(servicePayload)
	if err != nil {
		log.Error(fmt.Sprintf("can't create service for %s, got %v", clusterName, err))
		return err
	}

	return nil
}

// DeleteService delete a service of metric of GPUManager
func (s *ServiceOperator) DeleteService(clusterName string) error {
	t := trace.New(fmt.Sprintf("DeleteService-%s", clusterName))
	defer t.LogIfLong(time.Second)

	t.Step("delete service metric")
	err := s.K8sCli.CoreV1().Services(metav1.NamespaceSystem).Delete(DefaultServiceName, metav1.NewDeleteOptions(0))
	if err != nil && !apierror.IsNotFound(err) {
		log.Error(fmt.Sprintf("can't delete service for %s, got %v", clusterName, err))
		return err
	}

	return nil
}
