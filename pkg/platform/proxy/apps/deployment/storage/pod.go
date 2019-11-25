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
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/util"
)

// PodREST implements the REST endpoint for find pods by a deployment.
type PodREST struct {
	rest.Storage
	platformClient platforminternalclient.PlatformInterface
}

var _ rest.Getter = &PodREST{}
var _ rest.GroupVersionKindProvider = &PodREST{}

// GroupVersionKind is used to specify a particular GroupVersionKind to discovery.
func (r *PodREST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return corev1.SchemeGroupVersion.WithKind("PodList")
}

// New returns an empty object that can be used with Create and Update after
// request data has been put into it.
func (r *PodREST) New() runtime.Object {
	return &corev1.PodList{}
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *PodREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
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
		return nil, errors.NewNotFound(extensionsv1beta1.Resource("deployments/pods"), name)
	}

	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	listOptions := metav1.ListOptions{LabelSelector: selector.String()}
	rsList, err := client.ExtensionsV1beta1().ReplicaSets(namespaceName).List(listOptions)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	// list all of the pod, by deployment labels
	podAllList, err := client.CoreV1().Pods(namespaceName).List(listOptions)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	podList := &corev1.PodList{
		Items: make([]corev1.Pod, 0),
	}
	for _, rs := range rsList.Items {
		for _, references := range rs.ObjectMeta.OwnerReferences {
			if (references.Kind == "Deployment") && (references.Name == name) {
				for _, value := range podAllList.Items {
					for _, podReferences := range value.ObjectMeta.OwnerReferences {
						if (podReferences.Kind == "ReplicaSet") && (podReferences.Name == rs.Name) {
							podList.Items = append(podList.Items, value)
						}
					}
				}
			}
		}
	}

	return podList, nil
}
