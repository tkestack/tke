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
	"fmt"
	"sort"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/proxy"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/page"
)

// PodREST implements the REST endpoint for find pods by a node.
type PodREST struct {
	rest.Storage
	platformClient platforminternalclient.PlatformInterface
}

var _ rest.GetterWithOptions = &PodREST{}
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

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *PodREST) NewGetOptions() (runtime.Object, bool, string) {
	return &metav1.ListOptions{}, false, ""
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *PodREST) Get(ctx context.Context, name string, options runtime.Object) (runtime.Object, error) {
	client, err := proxy.ClientSet(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}
	listOpts := options.(*metav1.ListOptions)

	if apiclient.ClusterVersionIsBefore19(client) {
		return listPodByExtensions(ctx, client, name, listOpts)
	}

	return listPodByApps(ctx, client, name, listOpts)
}

func listPodByExtensions(ctx context.Context, client *kubernetes.Clientset, name string, listOpts *metav1.ListOptions) (runtime.Object, error) {
	if listOpts.FieldSelector == "" {
		listOpts.FieldSelector = fmt.Sprintf("spec.nodeName=%s", name)
	} else {
		listOpts.FieldSelector = listOpts.FieldSelector + "," + fmt.Sprintf("spec.nodeName=%s", name)
	}
	option := listOpts.DeepCopy()
	option.Limit = 0
	option.Continue = ""
	podList, err := client.CoreV1().Pods(metav1.NamespaceAll).List(ctx, *option)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	pods := util.NewPods(podList.Items)
	sort.Sort(pods)
	podList.Items = pods.GetPods()
	if listOpts.Continue != "" {
		start, limit, err := page.DecodeContinue(ctx, "Node", name, listOpts.Continue)
		if err != nil {
			return nil, err
		}
		newStart := start + limit
		if int(newStart+limit) < len(podList.Items) {
			podList.Continue, err = page.EncodeContinue(ctx, "Node", name, newStart, limit)
			if err != nil {
				return nil, err
			}
			items := podList.Items[newStart : newStart+limit]
			podList.Items = items
		} else {
			items := podList.Items[newStart:len(podList.Items)]
			podList.Items = items
		}
	} else if listOpts.Limit != 0 {
		if int(listOpts.Limit) < len(podList.Items) {
			podList.Continue, err = page.EncodeContinue(ctx, "Node", name, 0, listOpts.Limit)
			if err != nil {
				return nil, err
			}
			items := podList.Items[:listOpts.Limit]
			podList.Items = items
		}
	}

	return podList, nil
}

func listPodByApps(ctx context.Context, client *kubernetes.Clientset, name string, listOpts *metav1.ListOptions) (runtime.Object, error) {
	// list all of the pod, by node name
	if listOpts.FieldSelector == "" {
		listOpts.FieldSelector = fmt.Sprintf("spec.nodeName=%s", name)
	} else {
		listOpts.FieldSelector = listOpts.FieldSelector + "," + fmt.Sprintf("spec.nodeName=%s", name)
	}
	option := listOpts.DeepCopy()
	option.Limit = 0
	option.Continue = ""
	podList, err := client.CoreV1().Pods(metav1.NamespaceAll).List(ctx, *option)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	pods := util.NewPods(podList.Items)
	sort.Sort(pods)
	podList.Items = pods.GetPods()
	if listOpts.Continue != "" {
		start, limit, err := page.DecodeContinue(ctx, "Node", name, listOpts.Continue)
		if err != nil {
			return nil, err
		}
		newStart := start + limit
		if int(newStart+limit) < len(podList.Items) {
			podList.Continue, err = page.EncodeContinue(ctx, "Node", name, newStart, limit)
			if err != nil {
				return nil, err
			}
			items := podList.Items[newStart : newStart+limit]
			podList.Items = items
		} else {
			items := podList.Items[newStart:len(podList.Items)]
			podList.Items = items
		}
	} else if listOpts.Limit != 0 {
		if int(listOpts.Limit) < len(podList.Items) {
			podList.Continue, err = page.EncodeContinue(ctx, "Node", name, 0, listOpts.Limit)
			if err != nil {
				return nil, err
			}
			items := podList.Items[:listOpts.Limit]
			podList.Items = items
		}
	}

	return podList, nil
}
