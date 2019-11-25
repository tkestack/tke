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

package storage

import (
	"context"
	"sort"

	appsV1Beta1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/util"
)

// EventREST implements the REST endpoint for find events by a statefulset.
type EventREST struct {
	rest.Storage
	platformClient platforminternalclient.PlatformInterface
}

var _ rest.Getter = &EventREST{}
var _ rest.GroupVersionKindProvider = &EventREST{}

// GroupVersionKind is used to specify a particular GroupVersionKind to discovery.
func (r *EventREST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return corev1.SchemeGroupVersion.WithKind("EventList")
}

// New returns an empty object that can be used with Create and Update after
// request data has been put into it.
func (r *EventREST) New() runtime.Object {
	return &corev1.EventList{}
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *EventREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	client, err := util.ClientSet(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	namespaceName, ok := request.NamespaceFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("a namespace must be specified")
	}

	statefulSet, err := client.AppsV1beta1().StatefulSets(namespaceName).Get(name, *options)
	if err != nil {
		return nil, errors.NewNotFound(appsV1Beta1.Resource("statefulsets/events"), name)
	}

	selector := fields.AndSelectors(
		fields.OneTermEqualSelector("involvedObject.uid", string(statefulSet.UID)),
		fields.OneTermEqualSelector("involvedObject.name", statefulSet.Name),
		fields.OneTermEqualSelector("involvedObject.namespace", statefulSet.Namespace),
		fields.OneTermEqualSelector("involvedObject.kind", "StatefulSet"))
	listOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}
	statefulSetEvents, err := client.CoreV1().Events(namespaceName).List(listOptions)
	if err != nil {
		return nil, err
	}

	var events util.EventSlice
	for _, statefulSetEvent := range statefulSetEvents.Items {
		events = append(events, statefulSetEvent)
	}

	podSelector, err := metav1.LabelSelectorAsSelector(statefulSet.Spec.Selector)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	// list all of the pod, by stateful set labels
	podListOptions := metav1.ListOptions{LabelSelector: podSelector.String()}
	podAllList, err := client.CoreV1().Pods(namespaceName).List(podListOptions)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	for _, pod := range podAllList.Items {
		for _, podReferences := range pod.ObjectMeta.OwnerReferences {
			if (podReferences.Kind == "StatefulSet") && (podReferences.Name == statefulSet.Name) {
				podEventsSelector := fields.AndSelectors(
					fields.OneTermEqualSelector("involvedObject.uid", string(pod.UID)),
					fields.OneTermEqualSelector("involvedObject.name", pod.Name),
					fields.OneTermEqualSelector("involvedObject.namespace", pod.Namespace),
					fields.OneTermEqualSelector("involvedObject.kind", "Pod"))
				podEventsListOptions := metav1.ListOptions{
					FieldSelector: podEventsSelector.String(),
				}
				podEvents, err := client.CoreV1().Events(namespaceName).List(podEventsListOptions)
				if err != nil {
					return nil, err
				}

				for _, podEvent := range podEvents.Items {
					events = append(events, podEvent)
				}
			}
		}
	}

	sort.Sort(events)

	return &corev1.EventList{
		Items: events,
	}, nil
}
