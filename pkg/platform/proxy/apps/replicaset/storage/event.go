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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/util"
)

// EventREST implements the REST endpoint for find events by a daemonset.
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

	if controllerutil.IsClusterVersionBefore1_9(client) {
		return listEventsByExtensions(client, namespaceName, name, options)
	}
	return listEventsByApps(client, namespaceName, name, options)
}

func listEventsByExtensions(client *kubernetes.Clientset, namespaceName, name string, options *metav1.GetOptions) (runtime.Object, error) {
	replicaSet, err := client.ExtensionsV1beta1().ReplicaSets(namespaceName).Get(name, *options)
	if err != nil {
		return nil, errors.NewNotFound(extensionsv1beta1.Resource("replicasets/events"), name)
	}

	selector := fields.AndSelectors(
		fields.OneTermEqualSelector("involvedObject.uid", string(replicaSet.UID)),
		fields.OneTermEqualSelector("involvedObject.name", replicaSet.Name),
		fields.OneTermEqualSelector("involvedObject.namespace", replicaSet.Namespace),
		fields.OneTermEqualSelector("involvedObject.kind", "ReplicaSet"))
	listOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}
	return client.CoreV1().Events(namespaceName).List(listOptions)
}

func listEventsByApps(client *kubernetes.Clientset, namespaceName, name string, options *metav1.GetOptions) (runtime.Object, error) {
	replicaSet, err := client.AppsV1().ReplicaSets(namespaceName).Get(name, *options)
	if err != nil {
		return nil, errors.NewNotFound(appsv1.Resource("replicasets/events"), name)
	}

	selector := fields.AndSelectors(
		fields.OneTermEqualSelector("involvedObject.uid", string(replicaSet.UID)),
		fields.OneTermEqualSelector("involvedObject.name", replicaSet.Name),
		fields.OneTermEqualSelector("involvedObject.namespace", replicaSet.Namespace),
		fields.OneTermEqualSelector("involvedObject.kind", "ReplicaSet"))
	listOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}
	return client.CoreV1().Events(namespaceName).List(listOptions)
}
