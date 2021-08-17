/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

	"tkestack.io/tke/pkg/platform/proxy"
	"tkestack.io/tke/pkg/util/apiclient"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
)

// HPARest implements the REST endpoint for find hpalist by a statefulSet.
type HPARest struct {
	rest.Storage
	platformClient platforminternalclient.PlatformInterface
}

var _ rest.Getter = &HPARest{}
var _ rest.GroupVersionKindProvider = &HPARest{}

// GroupVersionKind is used to specify a particular GroupVersionKind to discovery.
func (r *HPARest) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return autoscalingv1.SchemeGroupVersion.WithKind("HorizontalPodAutoscalerList")
}

// New returns an empty object that can be used with Create and Update after
// request data has been put into it.
func (r *HPARest) New() runtime.Object {
	return &autoscalingv1.HorizontalPodAutoscalerList{}
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *HPARest) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	client, err := proxy.ClientSet(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	namespaceName, ok := request.NamespaceFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("a namespace must be specified")
	}

	if apiclient.ClusterVersionIsBefore19(client) {
		return listHPAsByAppsBeta(ctx, client, namespaceName, name, options)
	}
	return listHPAsByApps(ctx, client, namespaceName, name, options)
}

func listHPAsByAppsBeta(ctx context.Context, client *kubernetes.Clientset, namespaceName, name string, options *metav1.GetOptions) (runtime.Object, error) {
	statefulSet, err := client.AppsV1beta1().StatefulSets(namespaceName).Get(ctx, name, *options)
	if err != nil {
		return nil, errors.NewNotFound(extensionsv1beta1.Resource("statefulsets/horizontalpodautoscalers"), name)
	}

	hpas, err := client.AutoscalingV1().HorizontalPodAutoscalers(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	hpaList := &autoscalingv1.HorizontalPodAutoscalerList{
		Items: make([]autoscalingv1.HorizontalPodAutoscaler, 0),
	}
	for _, hpa := range hpas.Items {
		if hpa.Spec.ScaleTargetRef.Name == statefulSet.ObjectMeta.Name && hpa.Spec.ScaleTargetRef.Kind == "StatefulSet" {
			hpaList.Items = append(hpaList.Items, hpa)
		}
	}

	return hpaList, nil
}

func listHPAsByApps(ctx context.Context, client *kubernetes.Clientset, namespaceName, name string, options *metav1.GetOptions) (runtime.Object, error) {
	statefulSet, err := client.AppsV1().StatefulSets(namespaceName).Get(ctx, name, *options)
	if err != nil {
		return nil, errors.NewNotFound(appsv1.Resource("statefulsets/horizontalpodautoscalers"), name)
	}

	hpas, err := client.AutoscalingV1().HorizontalPodAutoscalers(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	hpaList := &autoscalingv1.HorizontalPodAutoscalerList{
		Items: make([]autoscalingv1.HorizontalPodAutoscaler, 0),
	}
	for _, hpa := range hpas.Items {
		if hpa.Spec.ScaleTargetRef.Name == statefulSet.ObjectMeta.Name && hpa.Spec.ScaleTargetRef.Kind == "StatefulSet" {
			hpaList.Items = append(hpaList.Items, hpa)
		}
	}

	return hpaList, nil
}
