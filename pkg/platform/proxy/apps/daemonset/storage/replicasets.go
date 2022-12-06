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
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/proxy"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/page"
)

// ReplicaSetsREST implements the REST endpoint for find replicasets by a deployment.
type ReplicaSetsREST struct {
	rest.Storage
	platformClient platforminternalclient.PlatformInterface
}

var _ rest.GetterWithOptions = &ReplicaSetsREST{}
var _ rest.GroupVersionKindProvider = &ReplicaSetsREST{}

// GroupVersionKind is used to specify a particular GroupVersionKind to discovery.
func (r *ReplicaSetsREST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return appsv1.SchemeGroupVersion.WithKind("ReplicaSetList")
}

// New returns an empty object that can be used with Create and Update after
// request data has been put into it.
func (r *ReplicaSetsREST) New() runtime.Object {
	return &appsv1.ReplicaSetList{}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *ReplicaSetsREST) NewGetOptions() (runtime.Object, bool, string) {
	return &metav1.ListOptions{}, false, ""
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *ReplicaSetsREST) Get(ctx context.Context, name string, options runtime.Object) (runtime.Object, error) {
	client, err := proxy.ClientSet(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}
	listOpts := options.(*metav1.ListOptions)
	metaOptions := &metav1.GetOptions{}
	if listOpts.ResourceVersion != "" {
		metaOptions.ResourceVersion = listOpts.ResourceVersion
	}
	namespaceName, ok := request.NamespaceFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("a namespace must be specified")
	}

	if apiclient.ClusterVersionIsBefore19(client) {
		return listReplicaSetsByExtensions(ctx, client, namespaceName, name, metaOptions, listOpts)
	}
	return listReplicaSetsByApps(ctx, client, namespaceName, name, metaOptions, listOpts)
}

func listReplicaSetsByExtensions(ctx context.Context, client *kubernetes.Clientset, namespaceName, name string, options *metav1.GetOptions, listOpts *metav1.ListOptions) (runtime.Object, error) {
	daemonSet, err := client.ExtensionsV1beta1().DaemonSets(namespaceName).Get(ctx, name, *options)
	if err != nil {
		return nil, errors.NewNotFound(extensionsv1beta1.Resource("daemonSets/replicasets"), name)
	}

	selector, err := metav1.LabelSelectorAsSelector(daemonSet.Spec.Selector)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	listPodsOptions := listOpts.DeepCopy()
	listPodsOptions.Continue = ""
	listPodsOptions.Limit = 0
	if listPodsOptions.LabelSelector == "" {
		listPodsOptions.LabelSelector = selector.String()
	} else {
		listPodsOptions.LabelSelector = listPodsOptions.LabelSelector + "," + selector.String()
	}
	rsAllList, err := client.ExtensionsV1beta1().ReplicaSets(namespaceName).List(ctx, *listPodsOptions)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	rsList := &extensionsv1beta1.ReplicaSetList{
		Items: make([]extensionsv1beta1.ReplicaSet, 0),
	}
	for _, rs := range rsAllList.Items {
		for _, references := range rs.ObjectMeta.OwnerReferences {
			if references.UID == daemonSet.GetUID() {
				rsList.Items = append(rsList.Items, rs)
			}
		}
	}

	if listOpts.Continue != "" {
		start, limit, err := page.DecodeContinue(ctx, "DaemonSet/ReplicaSets", name, listOpts.Continue)
		if err != nil {
			return nil, err
		}
		newStart := start + limit
		if int(newStart+limit) < len(rsList.Items) {
			rsList.Continue, err = page.EncodeContinue(ctx, "DaemonSet/ReplicaSets", name, newStart, limit)
			if err != nil {
				return nil, err
			}
			items := rsList.Items[newStart : newStart+limit]
			rsList.Items = items
		} else {
			items := rsList.Items[newStart:len(rsList.Items)]
			rsList.Items = items
		}
	} else if listOpts.Limit != 0 {
		if int(listOpts.Limit) < len(rsList.Items) {
			rsList.Continue, err = page.EncodeContinue(ctx, "DaemonSet/ReplicaSets", name, 0, listOpts.Limit)
			if err != nil {
				return nil, err
			}
			items := rsList.Items[:listOpts.Limit]
			rsList.Items = items
		}
	}

	return rsList, nil
}

func listReplicaSetsByApps(ctx context.Context, client *kubernetes.Clientset, namespaceName, name string, options *metav1.GetOptions, listOpts *metav1.ListOptions) (runtime.Object, error) {
	daemonSet, err := client.AppsV1().DaemonSets(namespaceName).Get(ctx, name, *options)
	if err != nil {
		return nil, errors.NewNotFound(appsv1.Resource("daemonSets/replicasets"), name)
	}

	selector, err := metav1.LabelSelectorAsSelector(daemonSet.Spec.Selector)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	listPodsOptions := listOpts.DeepCopy()
	listPodsOptions.Continue = ""
	listPodsOptions.Limit = 0
	if listPodsOptions.LabelSelector == "" {
		listPodsOptions.LabelSelector = selector.String()
	} else {
		listPodsOptions.LabelSelector = listPodsOptions.LabelSelector + "," + selector.String()
	}
	rsAllList, err := client.AppsV1().ReplicaSets(namespaceName).List(ctx, *listPodsOptions)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	rsList := &appsv1.ReplicaSetList{
		Items: make([]appsv1.ReplicaSet, 0),
	}
	for _, rs := range rsAllList.Items {
		for _, references := range rs.ObjectMeta.OwnerReferences {
			if references.UID == daemonSet.GetUID() {
				rsList.Items = append(rsList.Items, rs)
			}
		}
	}

	if listOpts.Continue != "" {
		start, limit, err := page.DecodeContinue(ctx, "DaemonSet/ReplicaSets", name, listOpts.Continue)
		if err != nil {
			return nil, err
		}
		newStart := start + limit
		if int(newStart+limit) < len(rsList.Items) {
			rsList.Continue, err = page.EncodeContinue(ctx, "DaemonSet/ReplicaSets", name, newStart, limit)
			if err != nil {
				return nil, err
			}
			items := rsList.Items[newStart : newStart+limit]
			rsList.Items = items
		} else {
			items := rsList.Items[newStart:len(rsList.Items)]
			rsList.Items = items
		}
	} else if listOpts.Limit != 0 {
		if int(listOpts.Limit) < len(rsList.Items) {
			rsList.Continue, err = page.EncodeContinue(ctx, "DaemonSet/ReplicaSets", name, 0, listOpts.Limit)
			if err != nil {
				return nil, err
			}
			items := rsList.Items[:listOpts.Limit]
			rsList.Items = items
		}
	}

	return rsList, nil
}
