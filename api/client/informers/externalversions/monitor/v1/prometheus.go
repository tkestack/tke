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

package v1

import (
	"context"
	time "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	versioned "tkestack.io/tke/api/client/clientset/versioned"
	internalinterfaces "tkestack.io/tke/api/client/informers/externalversions/internalinterfaces"
	v1 "tkestack.io/tke/api/client/listers/monitor/v1"
	monitorv1 "tkestack.io/tke/api/monitor/v1"
)

// PrometheusInformer provides access to a shared informer and lister for
// Prometheuses.
type PrometheusInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.PrometheusLister
}

type prometheusInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewPrometheusInformer constructs a new informer for Prometheus type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewPrometheusInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredPrometheusInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredPrometheusInformer constructs a new informer for Prometheus type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredPrometheusInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MonitorV1().Prometheuses().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MonitorV1().Prometheuses().Watch(context.TODO(), options)
			},
		},
		&monitorv1.Prometheus{},
		resyncPeriod,
		indexers,
	)
}

func (f *prometheusInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredPrometheusInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *prometheusInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&monitorv1.Prometheus{}, f.defaultInformer)
}

func (f *prometheusInformer) Lister() v1.PrometheusLister {
	return v1.NewPrometheusLister(f.Informer().GetIndexer())
}
