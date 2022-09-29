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

	"tkestack.io/tke/pkg/util/apiclient"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/proxy"
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
	client, err := proxy.ClientSet(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	namespaceName, ok := request.NamespaceFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("a namespace must be specified")
	}

	if apiclient.ClusterVersionIsBefore19(client) {
		return listEventsByExtensions(ctx, client, namespaceName, name, options)
	}
	return listEventsByApps(ctx, client, namespaceName, name, options)
}

func listEventsByExtensions(ctx context.Context, client *kubernetes.Clientset, namespaceName, name string, options *metav1.GetOptions) (runtime.Object, error) {
	deployment, err := client.ExtensionsV1beta1().Deployments(namespaceName).Get(ctx, name, *options)
	if err != nil {
		return nil, errors.NewNotFound(extensionsv1beta1.Resource("deployments/events"), name)
	}

	var resultEvents util.EventSlice

	events, errs := getAboutDeployEvents(ctx, client, deployment.Name, deployment.Namespace, string(deployment.UID))
	if len(errs) > 0 {
		return nil, utilerrors.NewAggregate(errs)
	}

	involvedObjectUIDMap := util.GetInvolvedObjectUIDMap(events)
	if v, ok := involvedObjectUIDMap[string(deployment.UID)]; ok {
		resultEvents = append(resultEvents, v...)
	}

	rsSelector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	rsListOptions := metav1.ListOptions{LabelSelector: rsSelector.String()}
	rsList, err := client.ExtensionsV1beta1().ReplicaSets(namespaceName).List(ctx, rsListOptions)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	for _, rs := range rsList.Items {
		// Skip replicaSets with same labels but owned by other deployments
		var ownedByDeployment bool
		ownedByDeployment = true
		if rs.ObjectMeta.OwnerReferences != nil {
			for _, references := range rs.ObjectMeta.OwnerReferences {
				if (references.Kind == "Deployment") && (references.Name != name) {
					ownedByDeployment = false
					break
				}
			}
		}
		if !ownedByDeployment {
			continue
		}

		if v, ok := involvedObjectUIDMap[string(rs.UID)]; ok {
			resultEvents = append(resultEvents, v...)
		}

		for _, references := range rs.ObjectMeta.OwnerReferences {
			if (references.Kind == "Deployment") && (references.Name == name) {
				podSelector, err := metav1.LabelSelectorAsSelector(rs.Spec.Selector)
				if err != nil {
					return nil, errors.NewInternalError(err)
				}
				podListOptions := metav1.ListOptions{LabelSelector: podSelector.String()}
				podListByRS, err := client.CoreV1().Pods(namespaceName).List(ctx, podListOptions)
				if err != nil {
					return nil, err
				}
				// Events cannot be queried for the deleted pods
				for _, pod := range podListByRS.Items {
					for _, podReferences := range pod.ObjectMeta.OwnerReferences {
						if (podReferences.Kind == "ReplicaSet") && (podReferences.Name == rs.Name) {
							if v, ok := involvedObjectUIDMap[string(pod.UID)]; ok {
								resultEvents = append(resultEvents, v...)
							}
						}
					}
				}
			}
		}
	}

	sort.Sort(resultEvents)

	return &corev1.EventList{
		Items: resultEvents,
	}, nil
}

func listEventsByApps(ctx context.Context, client *kubernetes.Clientset, namespaceName, name string, options *metav1.GetOptions) (runtime.Object, error) {
	deployment, err := client.AppsV1().Deployments(namespaceName).Get(ctx, name, *options)
	if err != nil {
		return nil, errors.NewNotFound(appsv1.Resource("deployments/events"), name)
	}

	var resultEvents util.EventSlice

	events, errs := getAboutDeployEvents(ctx, client, deployment.Name, deployment.Namespace, string(deployment.UID))
	if len(errs) > 0 {
		return nil, utilerrors.NewAggregate(errs)
	}

	involvedObjectUIDMap := util.GetInvolvedObjectUIDMap(events)
	if v, ok := involvedObjectUIDMap[string(deployment.UID)]; ok {
		resultEvents = append(resultEvents, v...)
	}

	rsSelector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	rsListOptions := metav1.ListOptions{LabelSelector: rsSelector.String()}
	rsList, err := client.AppsV1().ReplicaSets(namespaceName).List(ctx, rsListOptions)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	for _, rs := range rsList.Items {
		// Skip replicaSets with same labels but owned by other deployments
		var ownedByDeployment bool
		ownedByDeployment = true
		if rs.ObjectMeta.OwnerReferences != nil {
			for _, references := range rs.ObjectMeta.OwnerReferences {
				if (references.Kind == "Deployment") && (references.Name != name) {
					ownedByDeployment = false
					break
				}
			}
		}
		if !ownedByDeployment {
			continue
		}

		if v, ok := involvedObjectUIDMap[string(rs.UID)]; ok {
			resultEvents = append(resultEvents, v...)
		}

		for _, references := range rs.ObjectMeta.OwnerReferences {
			if (references.Kind == "Deployment") && (references.Name == name) {
				podSelector, err := metav1.LabelSelectorAsSelector(rs.Spec.Selector)
				if err != nil {
					return nil, errors.NewInternalError(err)
				}
				podListOptions := metav1.ListOptions{LabelSelector: podSelector.String()}
				podListByRS, err := client.CoreV1().Pods(namespaceName).List(ctx, podListOptions)
				if err != nil {
					return nil, err
				}
				for _, pod := range podListByRS.Items {
					for _, podReferences := range pod.ObjectMeta.OwnerReferences {
						if (podReferences.Kind == "ReplicaSet") && (podReferences.Name == rs.Name) {
							if v, ok := involvedObjectUIDMap[string(pod.UID)]; ok {
								resultEvents = append(resultEvents, v...)
							}
						}
					}
				}
			}
		}
	}

	sort.Sort(resultEvents)

	return &corev1.EventList{
		Items: resultEvents,
	}, nil
}

// getAboutDeployEvents Query all events in the namespace  or Query the Deployment ReplicaSet Pod asynchronously
func getAboutDeployEvents(ctx context.Context, client *kubernetes.Clientset, name, namespace, uid string) (util.EventSlice, []error) {
	return util.GetResourcesEvents(ctx, client, namespace, []metav1.ListOptions{
		{
			FieldSelector: fields.AndSelectors(
				fields.OneTermEqualSelector("involvedObject.uid", uid),
				fields.OneTermEqualSelector("involvedObject.name", name),
				fields.OneTermEqualSelector("involvedObject.namespace", namespace),
				fields.OneTermEqualSelector("involvedObject.kind", "Deployment")).String(),
			ResourceVersion: "0",
		},
		{
			FieldSelector: fields.AndSelectors(
				fields.OneTermEqualSelector("involvedObject.namespace", namespace),
				fields.OneTermEqualSelector("involvedObject.kind", "ReplicaSet")).String(),
			ResourceVersion: "0",
		},
		{
			FieldSelector: fields.AndSelectors(
				fields.OneTermEqualSelector("involvedObject.namespace", namespace),
				fields.OneTermEqualSelector("involvedObject.kind", "Pod")).String(),
			ResourceVersion: "0",
		},
	})
}
