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
	// DefaultConfigmapName defines the name of service for the metric service of GPUManager
	DefaultConfigmapName = "gpu-quota-admission"
)

// ConfigmapOperator is used to operate remote service account and role-binding
type ConfigmapOperator struct {
	K8sCli kubernetes.Interface
}

// CreateConfigmap create a service to help collecting metrics of GPUManager
func (s *ConfigmapOperator) CreateConfigmap(clusterName string) error {
	t := trace.New(fmt.Sprintf("CreateConfigmap-%s", clusterName))
	defer t.LogIfLong(time.Second)

	reader := strings.NewReader(template.ConfigMapTemplate)

	configmapPayload := &corev1.ConfigMap{}
	err := yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(configmapPayload)
	if err != nil {
		return err
	}
	configmapPayload.Name = DefaultConfigmapName

	t.Step("create configmap")
	_, err = s.K8sCli.CoreV1().ConfigMaps(metav1.NamespaceSystem).Create(configmapPayload)
	if err != nil {
		log.Error(fmt.Sprintf("can't create configmap for %s, got %v", clusterName, err))
		return err
	}

	return nil
}

// DeleteConfigmap delete a configmap of metric of GPUManager
func (s *ConfigmapOperator) DeleteConfigmap(clusterName string) error {
	t := trace.New(fmt.Sprintf("DeleteConfigmap-%s", clusterName))
	defer t.LogIfLong(time.Second)

	t.Step("delete configmap")
	err := s.K8sCli.CoreV1().ConfigMaps(metav1.NamespaceSystem).Delete(DefaultConfigmapName, metav1.NewDeleteOptions(0))
	if err != nil && !apierror.IsNotFound(err) {
		log.Error(fmt.Sprintf("can't delete configmap for %s, got %v", clusterName, err))
		return err
	}

	return nil
}
