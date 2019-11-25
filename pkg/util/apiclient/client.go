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

package apiclient

import (
	"errors"
	"net/http"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	tkeclientset "tkestack.io/tke/api/client/clientset/versioned"
)

// GetClientset return clientset
func GetClientset(masterEndpoint string, token string, caCert []byte) (*kubernetes.Clientset, error) {
	restConfig := &rest.Config{
		Host:        masterEndpoint,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: caCert,
		},
		Timeout: 5 * time.Second,
	}

	return kubernetes.NewForConfig(restConfig)
}

// GetPlatformClientset return clientset
func GetPlatformClientset(masterEndpoint string, token string, caCert []byte) (tkeclientset.Interface, error) {
	restConfig := &rest.Config{
		Host:        masterEndpoint,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: caCert,
		},
		Timeout: 5 * time.Second,
	}

	return tkeclientset.NewForConfig(restConfig)
}

// CheckAPIHealthz check healthz
func CheckAPIHealthz(client rest.Interface) bool {
	healthStatus := 0
	client.Get().AbsPath("/healthz").Do().StatusCode(&healthStatus)
	return healthStatus == http.StatusOK
}

// CheckDeployment check Deployment current replicas is equal to desired and all pods are running
func CheckDeployment(client kubernetes.Interface, namespace string, name string) (bool, error) {
	deployment, err := client.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if *deployment.Spec.Replicas != deployment.Status.Replicas {
		return false, errors.New("deployment.Spec.Replicas != deployment.Status.Replicas")
	}

	labelSelector := metav1.FormatLabelSelector(deployment.Spec.Selector)
	pods, err := client.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return false, err
	}
	for _, pod := range pods.Items {
		if !IsPodReady(&pod) {
			return false, nil
		}
	}

	return true, nil
}

// CheckStatefulSet check StatefulSet current replicas is equal to desired and all pods are running
func CheckStatefulSet(client kubernetes.Interface, namespace string, name string) (bool, error) {
	statefulSet, err := client.AppsV1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if *statefulSet.Spec.Replicas != statefulSet.Status.Replicas {
		return false, errors.New("statefulSet.Spec.Replicas != statefulSet.Status.Replicas")
	}

	labelSelector := metav1.FormatLabelSelector(statefulSet.Spec.Selector)
	pods, err := client.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return false, err
	}
	for _, pod := range pods.Items {
		if !IsPodReady(&pod) {
			return false, nil
		}
	}

	return true, nil
}

// CheckDaemonset check daemonset current replicas is equal to desired and all pods are running
func CheckDaemonset(client kubernetes.Interface, namespace string, name string) (bool, error) {
	daemonSet, err := client.AppsV1().DaemonSets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if daemonSet.Status.DesiredNumberScheduled != daemonSet.Status.NumberReady {
		return false, errors.New("daemonSet.Status.DesiredNumberScheduled != daemonSet.Status.NumberReady")
	}

	return true, nil
}

// IsPodReady returns true if a pod is ready; false otherwise.
func IsPodReady(pod *corev1.Pod) bool {
	return isPodReadyConditionTrue(pod.Status)
}

// IsPodReadyConditionTrue returns true if a pod is ready; false otherwise.
func isPodReadyConditionTrue(status corev1.PodStatus) bool {
	condition := getPodReadyCondition(status)
	return condition != nil && condition.Status == corev1.ConditionTrue
}

// GetPodReadyCondition extracts the pod ready condition from the given status and returns that.
// Returns nil if the condition is not present.
func getPodReadyCondition(status corev1.PodStatus) *corev1.PodCondition {
	_, condition := getPodCondition(&status, corev1.PodReady)
	return condition
}

// GetPodCondition extracts the provided condition from the given status and returns that.
// Returns nil and -1 if the condition is not present, and the index of the located condition.
func getPodCondition(status *corev1.PodStatus, conditionType corev1.PodConditionType) (int, *corev1.PodCondition) {
	if status == nil {
		return -1, nil
	}
	return getPodConditionFromList(status.Conditions, conditionType)
}

// GetPodConditionFromList extracts the provided condition from the given list of condition and
// returns the index of the condition and the condition. Returns -1 and nil if the condition is not present.
func getPodConditionFromList(conditions []corev1.PodCondition, conditionType corev1.PodConditionType) (int, *corev1.PodCondition) {
	if conditions == nil {
		return -1, nil
	}
	for i := range conditions {
		if conditions[i].Type == conditionType {
			return i, &conditions[i]
		}
	}
	return -1, nil
}
