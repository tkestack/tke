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

package util

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
)

// EventSlice implements sort.Interface for []Event based on the EventTime field.
type EventSlice []corev1.Event

func (e EventSlice) Len() int {
	return len(e)
}

func (e EventSlice) Less(i, j int) bool {
	return e[i].LastTimestamp.Before(&e[j].LastTimestamp)
}

func (e EventSlice) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// GetEvents list the resource events by resource namespace and name.
func GetEvents(client *kubernetes.Clientset, uid, namespace, name, kind string) (*corev1.EventList, error) {
	selector := fields.AndSelectors(
		fields.OneTermEqualSelector("involvedObject.uid", uid),
		fields.OneTermEqualSelector("involvedObject.name", name),
		fields.OneTermEqualSelector("involvedObject.namespace", namespace),
		fields.OneTermEqualSelector("involvedObject.kind", kind))
	listOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}
	return client.CoreV1().Events(namespace).List(listOptions)
}
