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

	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
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

// EventREST implements the REST endpoint for find events by a deployment.
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

	deployment, err := client.ExtensionsV1beta1().Deployments(namespaceName).Get(name, *options)
	if err != nil {
		return nil, errors.NewNotFound(extensionsv1beta1.Resource("deployments/events"), name)
	}

	selector := fields.AndSelectors(
		fields.OneTermEqualSelector("involvedObject.uid", string(deployment.UID)),
		fields.OneTermEqualSelector("involvedObject.name", deployment.Name),
		fields.OneTermEqualSelector("involvedObject.namespace", deployment.Namespace),
		fields.OneTermEqualSelector("involvedObject.kind", "Deployment"))
	listOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}
	deploymentEvents, err := client.CoreV1().Events(namespaceName).List(listOptions)
	if err != nil {
		return nil, err
	}

	var events util.EventSlice
	for _, deploymentEvent := range deploymentEvents.Items {
		events = append(events, deploymentEvent)
	}

	rsSelector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	rsListOptions := metav1.ListOptions{LabelSelector: rsSelector.String()}
	rsList, err := client.ExtensionsV1beta1().ReplicaSets(namespaceName).List(rsListOptions)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	for _, rs := range rsList.Items {
		rsEventsSelector := fields.AndSelectors(
			fields.OneTermEqualSelector("involvedObject.uid", string(rs.UID)),
			fields.OneTermEqualSelector("involvedObject.name", rs.Name),
			fields.OneTermEqualSelector("involvedObject.namespace", rs.Namespace),
			fields.OneTermEqualSelector("involvedObject.kind", "ReplicaSet"))
		rsEventsListOptions := metav1.ListOptions{
			FieldSelector: rsEventsSelector.String(),
		}
		rsEvents, err := client.CoreV1().Events(namespaceName).List(rsEventsListOptions)
		if err != nil {
			return nil, err
		}

		for _, rsEvent := range rsEvents.Items {
			events = append(events, rsEvent)
		}

		for _, references := range rs.ObjectMeta.OwnerReferences {
			if (references.Kind == "Deployment") && (references.Name == name) {
				podSelector, err := metav1.LabelSelectorAsSelector(rs.Spec.Selector)
				if err != nil {
					return nil, errors.NewInternalError(err)
				}
				podListOptions := metav1.ListOptions{LabelSelector: podSelector.String()}
				podListByRS, err := client.CoreV1().Pods(namespaceName).List(podListOptions)
				if err != nil {
					return nil, err
				}
				for _, pod := range podListByRS.Items {
					for _, podReferences := range pod.ObjectMeta.OwnerReferences {
						if (podReferences.Kind == "ReplicaSet") || (podReferences.Name == rs.Name) {
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
			}
		}
	}

	sort.Sort(events)

	return &corev1.EventList{
		Items: events,
	}, nil
}
