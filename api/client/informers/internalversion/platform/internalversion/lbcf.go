/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

// Code generated by informer-gen. DO NOT EDIT.

package internalversion

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	clientsetinternalversion "tkestack.io/tke/api/client/clientset/internalversion"
	internalinterfaces "tkestack.io/tke/api/client/informers/internalversion/internalinterfaces"
	internalversion "tkestack.io/tke/api/client/listers/platform/internalversion"
	platform "tkestack.io/tke/api/platform"
)

// LBCFInformer provides access to a shared informer and lister for
// LBCFs.
type LBCFInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() internalversion.LBCFLister
}

type lBCFInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewLBCFInformer constructs a new informer for LBCF type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewLBCFInformer(client clientsetinternalversion.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredLBCFInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredLBCFInformer constructs a new informer for LBCF type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredLBCFInformer(client clientsetinternalversion.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Platform().LBCFs().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Platform().LBCFs().Watch(context.TODO(), options)
			},
		},
		&platform.LBCF{},
		resyncPeriod,
		indexers,
	)
}

func (f *lBCFInformer) defaultInformer(client clientsetinternalversion.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredLBCFInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *lBCFInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&platform.LBCF{}, f.defaultInformer)
}

func (f *lBCFInformer) Lister() internalversion.LBCFLister {
	return internalversion.NewLBCFLister(f.Informer().GetIndexer())
}
