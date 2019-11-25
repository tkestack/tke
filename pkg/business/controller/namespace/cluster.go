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

package namespace

import (
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "tkestack.io/tke/api/business/v1"
	"tkestack.io/tke/pkg/business/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/resource"
)

func calculateNamespaceUsed(kubeClient *kubernetes.Clientset, namespace *v1.Namespace) (message, reason string, list v1.ResourceList) {
	list = make(v1.ResourceList)
	resourceQuotaList, err := kubeClient.CoreV1().ResourceQuotas(namespace.Spec.Namespace).List(metav1.ListOptions{Limit: 1})
	if err != nil {
		message = "ConnectClusterError"
		reason = err.Error()
		return
	}
	if len(resourceQuotaList.Items) == 0 {
		message = "ResourceQuotaNotFound"
		reason = "ResourceQuota in the business cluster did not find"
		return
	}
	message = ""
	reason = ""
	list = resource.ConvertFromCoreV1ResourceList(resourceQuotaList.Items[0].Status.Used)
	return
}

func checkNamespaceOnCluster(kubeClient *kubernetes.Clientset, namespace *v1.Namespace) (message, reason string) {
	ns, err := kubeClient.CoreV1().Namespaces().Get(namespace.Spec.Namespace, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		message = "NamespaceNotFound"
		reason = "The namespace on the business cluster does not exist, it may have been deleted"
		return
	}
	if err != nil {
		message = "ConnectClusterError"
		reason = err.Error()
		return
	}
	projectName, ok := ns.ObjectMeta.Labels[util.LabelProjectName]
	if !ok {
		message = "NamespaceNoLabel"
		reason = "No project label were found on the namespace within the business cluster"
		return
	}
	if projectName != namespace.ObjectMeta.Namespace {
		message = "ConflictProject"
		reason = fmt.Sprintf("The namespace in the business cluster currently belongs to another project: %s", projectName)
		return
	}
	if ns.Status.Phase != corev1.NamespaceActive {
		message = "NamespaceNotActive"
		reason = "The current state of the namespace in the business cluster is abnormal."
		return
	}
	message = ""
	reason = ""
	return
}

func createNamespaceOnCluster(kubeClient *kubernetes.Clientset, namespace *v1.Namespace) error {
	ns, err := kubeClient.CoreV1().Namespaces().Get(namespace.Spec.Namespace, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		// create namespace
		nsOnCluster := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace.Spec.Namespace,
				Labels: map[string]string{
					util.LabelProjectName:   namespace.ObjectMeta.Namespace,
					util.LabelNamespaceName: namespace.ObjectMeta.Name,
				},
			},
		}
		_, err := kubeClient.CoreV1().Namespaces().Create(nsOnCluster)
		if err != nil {
			log.Error("Failed to create the namespace on cluster", log.String("namespaceName", namespace.ObjectMeta.Name), log.String("clusterName", namespace.Spec.ClusterName), log.Err(err))
			return err
		}
		return nil
	}
	if err != nil {
		log.Error("Failed to get the namespace on cluster", log.String("namespace", namespace.Spec.Namespace), log.String("namespaceName", namespace.ObjectMeta.Name), log.String("clusterName", namespace.Spec.ClusterName), log.Err(err))
		return err
	}
	projectName, ok := ns.ObjectMeta.Labels[util.LabelProjectName]
	if !ok {
		ns.Labels = make(map[string]string)
		ns.Labels[util.LabelProjectName] = namespace.ObjectMeta.Namespace
		ns.Labels[util.LabelNamespaceName] = namespace.ObjectMeta.Name
		_, err := kubeClient.CoreV1().Namespaces().Update(ns)
		if err != nil {
			log.Error("Failed to update the namespace on cluster", log.String("namespaceName", namespace.ObjectMeta.Name), log.String("clusterName", namespace.Spec.ClusterName), log.Err(err))
			return err
		}
		return nil
	}
	if projectName != namespace.ObjectMeta.Namespace {
		log.Error("The namespace in the cluster already belongs to another project and cannot be attributed to this project", log.String("clusterName", namespace.Spec.ClusterName), log.String("namespace", namespace.Spec.Namespace))
		return fmt.Errorf("namespace in the cluster already belongs to another project(%s) and cannot be attributed to this project(%s)", projectName, namespace.ObjectMeta.Namespace)
	}
	return nil
}

func createResourceQuotaOnCluster(kubeClient *kubernetes.Clientset, namespace *v1.Namespace) error {
	resourceQuota, err := kubeClient.CoreV1().ResourceQuotas(namespace.Spec.Namespace).Get(namespace.Spec.Namespace, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		// create resource quota
		resourceList := resource.ConvertToCoreV1ResourceList(namespace.Spec.Hard)
		rq := &corev1.ResourceQuota{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespace.Spec.Namespace,
				Namespace: namespace.Spec.Namespace,
			},
			Spec: corev1.ResourceQuotaSpec{
				Hard: resourceList,
			},
		}
		_, err := kubeClient.CoreV1().ResourceQuotas(namespace.Spec.Namespace).Create(rq)
		if err != nil {
			log.Error("Failed to create the resource quota on cluster", log.String("namespaceName", namespace.ObjectMeta.Name), log.String("clusterName", namespace.Spec.ClusterName), log.Err(err))
			return err
		}
		return nil
	}
	if err != nil {
		log.Error("Failed to get the resource quota on cluster", log.String("namespace", namespace.Spec.Namespace), log.String("namespaceName", namespace.ObjectMeta.Name), log.String("clusterName", namespace.Spec.ClusterName), log.Err(err))
		return err
	}
	resourceList := resource.ConvertToCoreV1ResourceList(namespace.Spec.Hard)
	if !reflect.DeepEqual(resourceQuota.Spec.Hard, resourceList) {
		resourceQuota.Spec.Hard = resourceList
		_, err := kubeClient.CoreV1().ResourceQuotas(namespace.Spec.Namespace).Update(resourceQuota)
		if err != nil {
			log.Error("Failed to update the resource quota on cluster", log.String("namespaceName", namespace.ObjectMeta.Name), log.String("clusterName", namespace.Spec.ClusterName), log.Err(err))
			return err
		}
	}
	return nil
}
