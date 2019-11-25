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
	autoscalingv1 "k8s.io/api/autoscaling/v1"
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

// HPARest implements the REST endpoint for find hpalist by a deployment.
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
		return nil, errors.NewNotFound(extensionsv1beta1.Resource("deployments/horizontalpodautoscalers"), name)
	}

	hpas, err := client.AutoscalingV1().HorizontalPodAutoscalers(namespaceName).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	hpaList := &autoscalingv1.HorizontalPodAutoscalerList{
		Items: make([]autoscalingv1.HorizontalPodAutoscaler, 0),
	}
	for _, hpa := range hpas.Items {
		if hpa.Spec.ScaleTargetRef.Name == deployment.ObjectMeta.Name && hpa.Spec.ScaleTargetRef.Kind == "Deployment" {
			hpaList.Items = append(hpaList.Items, hpa)
		}
	}

	return hpaList, nil
}
