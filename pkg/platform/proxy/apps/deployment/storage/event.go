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
	"sync"

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

type eventsFinder struct {
	wg             sync.WaitGroup
	mutex          sync.Mutex
	namespaceName  string
	platformClient platforminternalclient.PlatformInterface
	client         kubernetes.Clientset
	ctx            context.Context
	events         util.EventSlice
	errors         []error
}

func newEventsFinder(ctx context.Context, namespaceName string, client kubernetes.Clientset, platformClient platforminternalclient.PlatformInterface) *eventsFinder {
	return &eventsFinder{
		platformClient: platformClient,
		ctx:            ctx,
		client:         client,
		namespaceName:  namespaceName,
		events:         nil,
		errors:         make([]error, 0),
	}
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

	ef := newEventsFinder(ctx, namespaceName, *client, r.platformClient)

	if apiclient.ClusterVersionIsBefore19(client) {
		return ef.listEventsByExtensions(client, namespaceName, name, options)
	}
	return ef.listEventsByApps(client, namespaceName, name, options)
}

func (ef *eventsFinder) getEvents(listOptions metav1.ListOptions) {
	defer ef.wg.Done()

	events, err := ef.client.CoreV1().Events(ef.namespaceName).List(listOptions)
	if err != nil {
		ef.mutex.Lock()
		ef.errors = append(ef.errors, err)
		ef.mutex.Unlock()
		return
	}
	if len(events.Items) == 0 {
		return
	}

	ef.mutex.Lock()
	for _, event := range events.Items {
		ef.events = append(ef.events, event)
	}
	ef.mutex.Unlock()
}

func (ef *eventsFinder) listEventsByExtensions(client *kubernetes.Clientset, namespaceName, name string, options *metav1.GetOptions) (runtime.Object, error) {
	deployment, err := client.ExtensionsV1beta1().Deployments(namespaceName).Get(name, *options)
	if err != nil {
		return nil, errors.NewNotFound(extensionsv1beta1.Resource("deployments/events"), name)
	}

	deploymentSelector := fields.AndSelectors(
		fields.OneTermEqualSelector("involvedObject.uid", string(deployment.UID)),
		fields.OneTermEqualSelector("involvedObject.name", deployment.Name),
		fields.OneTermEqualSelector("involvedObject.namespace", deployment.Namespace),
		fields.OneTermEqualSelector("involvedObject.kind", "Deployment"))
	deploymentListOptions := metav1.ListOptions{
		FieldSelector: deploymentSelector.String(),
	}

	ef.wg.Add(1)
	go ef.getEvents(deploymentListOptions)

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

		rsEventsSelector := fields.AndSelectors(
			fields.OneTermEqualSelector("involvedObject.uid", string(rs.UID)),
			fields.OneTermEqualSelector("involvedObject.name", rs.Name),
			fields.OneTermEqualSelector("involvedObject.namespace", rs.Namespace),
			fields.OneTermEqualSelector("involvedObject.kind", "ReplicaSet"))
		rsEventsListOptions := metav1.ListOptions{
			FieldSelector: rsEventsSelector.String(),
		}
		ef.wg.Add(1)
		go ef.getEvents(rsEventsListOptions)

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
						if (podReferences.Kind == "ReplicaSet") && (podReferences.Name == rs.Name) {
							podEventsSelector := fields.AndSelectors(
								fields.OneTermEqualSelector("involvedObject.uid", string(pod.UID)),
								fields.OneTermEqualSelector("involvedObject.name", pod.Name),
								fields.OneTermEqualSelector("involvedObject.namespace", pod.Namespace),
								fields.OneTermEqualSelector("involvedObject.kind", "Pod"))
							podEventsListOptions := metav1.ListOptions{
								FieldSelector: podEventsSelector.String(),
							}
							ef.wg.Add(1)
							go ef.getEvents(podEventsListOptions)
						}
					}
				}
			}
		}
	}

	ef.wg.Wait()
	if len(ef.errors) > 0 {
		return nil, utilerrors.NewAggregate(ef.errors)
	}

	sort.Sort(ef.events)

	return &corev1.EventList{
		Items: ef.events,
	}, nil
}

func (ef *eventsFinder) listEventsByApps(client *kubernetes.Clientset, namespaceName, name string, options *metav1.GetOptions) (runtime.Object, error) {
	deployment, err := client.AppsV1().Deployments(namespaceName).Get(name, *options)
	if err != nil {
		return nil, errors.NewNotFound(appsv1.Resource("deployments/events"), name)
	}

	deploymentSelector := fields.AndSelectors(
		fields.OneTermEqualSelector("involvedObject.uid", string(deployment.UID)),
		fields.OneTermEqualSelector("involvedObject.name", deployment.Name),
		fields.OneTermEqualSelector("involvedObject.namespace", deployment.Namespace),
		fields.OneTermEqualSelector("involvedObject.kind", "Deployment"))
	deploymentListOptions := metav1.ListOptions{
		FieldSelector: deploymentSelector.String(),
	}

	ef.wg.Add(1)
	go ef.getEvents(deploymentListOptions)

	rsSelector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	rsListOptions := metav1.ListOptions{LabelSelector: rsSelector.String()}
	rsList, err := client.AppsV1().ReplicaSets(namespaceName).List(rsListOptions)
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

		rsEventsSelector := fields.AndSelectors(
			fields.OneTermEqualSelector("involvedObject.uid", string(rs.UID)),
			fields.OneTermEqualSelector("involvedObject.name", rs.Name),
			fields.OneTermEqualSelector("involvedObject.namespace", rs.Namespace),
			fields.OneTermEqualSelector("involvedObject.kind", "ReplicaSet"))
		rsEventsListOptions := metav1.ListOptions{
			FieldSelector: rsEventsSelector.String(),
		}
		ef.wg.Add(1)
		go ef.getEvents(rsEventsListOptions)

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
						if (podReferences.Kind == "ReplicaSet") && (podReferences.Name == rs.Name) {
							podEventsSelector := fields.AndSelectors(
								fields.OneTermEqualSelector("involvedObject.uid", string(pod.UID)),
								fields.OneTermEqualSelector("involvedObject.name", pod.Name),
								fields.OneTermEqualSelector("involvedObject.namespace", pod.Namespace),
								fields.OneTermEqualSelector("involvedObject.kind", "Pod"))
							podEventsListOptions := metav1.ListOptions{
								FieldSelector: podEventsSelector.String(),
							}
							ef.wg.Add(1)
							go ef.getEvents(podEventsListOptions)
						}
					}
				}
			}
		}
	}

	ef.wg.Wait()
	if len(ef.errors) > 0 {
		return nil, utilerrors.NewAggregate(ef.errors)
	}

	sort.Sort(ef.events)

	return &corev1.EventList{
		Items: ef.events,
	}, nil
}
