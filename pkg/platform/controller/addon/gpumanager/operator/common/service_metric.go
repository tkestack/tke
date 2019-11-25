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
	// DefaultServiceMetricName defines the name of service for the metric service of GPUManager
	DefaultServiceMetricName = "gpu-manager-metric"
)

// ServiceMetricOperator is used to operate remote service account and role-binding
type ServiceMetricOperator struct {
	K8sCli kubernetes.Interface
}

// CreateServiceMetric create a service to help collecting metrics of GPUManager
func (s *ServiceMetricOperator) CreateServiceMetric(clusterName string) error {
	t := trace.New(fmt.Sprintf("CreateServiceMetric-%s", clusterName))
	defer t.LogIfLong(time.Second)

	reader := strings.NewReader(template.MetricServiceTemplate)

	servicePayload := &corev1.Service{}
	err := yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(servicePayload)
	if err != nil {
		return err
	}
	servicePayload.Name = DefaultServiceMetricName

	t.Step("create service metric")
	_, err = s.K8sCli.CoreV1().Services(metav1.NamespaceSystem).Create(servicePayload)
	if err != nil {
		log.Error(fmt.Sprintf("can't create service metric for %s, got %v", clusterName, err))
		return err
	}

	return nil
}

// DeleteServiceMetric delete a service of metric of GPUManager
func (s *ServiceMetricOperator) DeleteServiceMetric(clusterName string) error {
	t := trace.New(fmt.Sprintf("DeleteServiceMetric-%s", clusterName))
	defer t.LogIfLong(time.Second)

	t.Step("delete service metric")
	err := s.K8sCli.CoreV1().Services(metav1.NamespaceSystem).Delete(DefaultServiceMetricName, metav1.NewDeleteOptions(0))
	if err != nil && !apierror.IsNotFound(err) {
		log.Error(fmt.Sprintf("can't delete service metric for %s, got %v", clusterName, err))
		return err
	}

	return nil
}
