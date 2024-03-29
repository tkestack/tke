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
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/registry/clusteraddontype"
	"tkestack.io/tke/pkg/platform/util"
)

// AddonREST implements the REST endpoint.
type AddonREST struct {
	rest.Storage
	store          *registry.Store
	platformClient platforminternalclient.PlatformInterface
}

var _ = rest.Getter(&AddonREST{})

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *AddonREST) New() runtime.Object {
	return &platform.ClusterAddon{}
}

// Get finds a resource in the storage by name and returns it.
func (r *AddonREST) Get(ctx context.Context, clusterName string, options *metav1.GetOptions) (runtime.Object, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, options)
	if err != nil {
		return nil, err
	}
	cluster := clusterObject.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}
	af := newAddonFinder(clusterName, r.platformClient)
	return af.findAll(ctx)
}

type addonFinder struct {
	wg             sync.WaitGroup
	mutex          sync.Mutex
	clusterName    string
	platformClient platforminternalclient.PlatformInterface
	addons         []platform.ClusterAddon
	errors         []error
}

func newAddonFinder(clusterName string, platformClient platforminternalclient.PlatformInterface) *addonFinder {
	return &addonFinder{
		clusterName:    clusterName,
		platformClient: platformClient,
		errors:         make([]error, 0),
		addons:         make([]platform.ClusterAddon, 0),
	}
}

type addonFinderFunc func(ctx context.Context, a *addonFinder)

var (
	allAddonFinders = []addonFinderFunc{
		persistentEvent,
		tappcontroller,
		csiOperator,
		cronHPA,
	}
)

func (a *addonFinder) findAll(ctx context.Context) (*platform.ClusterAddonList, error) {
	a.wg.Add(len(allAddonFinders))
	for _, f := range allAddonFinders {
		go f(ctx, a)
	}
	a.wg.Wait()
	if len(a.errors) > 0 {
		return nil, utilerrors.NewAggregate(a.errors)
	}
	return &platform.ClusterAddonList{
		Items: a.addons,
	}, nil
}

func persistentEvent(ctx context.Context, a *addonFinder) {
	defer a.wg.Done()
	l, err := a.platformClient.PersistentEvents().List(ctx, metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("spec.clusterName", a.clusterName).String(),
	})
	if err != nil {
		a.mutex.Lock()
		a.errors = append(a.errors, err)
		a.mutex.Unlock()
		return
	}
	if len(l.Items) == 0 {
		return
	}
	a.mutex.Lock()
	a.addons = append(a.addons, platform.ClusterAddon{
		ObjectMeta: metav1.ObjectMeta{
			Name:              l.Items[0].ObjectMeta.Name,
			CreationTimestamp: l.Items[0].ObjectMeta.CreationTimestamp,
		},
		Spec: platform.ClusterAddonSpec{
			Type:    string(clusteraddontype.PersistentEvent),
			Level:   clusteraddontype.Types[clusteraddontype.PersistentEvent].Level,
			Version: l.Items[0].Spec.Version,
		},
		Status: platform.ClusterAddonStatus{
			Version: l.Items[0].Status.Version,
			Phase:   string(l.Items[0].Status.Phase),
			Reason:  l.Items[0].Status.Reason,
		},
	})
	a.mutex.Unlock()
}

func tappcontroller(ctx context.Context, a *addonFinder) {
	defer a.wg.Done()
	l, err := a.platformClient.TappControllers().List(ctx, metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("spec.clusterName", a.clusterName).String(),
	})
	if err != nil {
		a.mutex.Lock()
		a.errors = append(a.errors, err)
		a.mutex.Unlock()
		return
	}
	if len(l.Items) == 0 {
		return
	}
	a.mutex.Lock()
	a.addons = append(a.addons, platform.ClusterAddon{
		ObjectMeta: metav1.ObjectMeta{
			Name:              l.Items[0].ObjectMeta.Name,
			CreationTimestamp: l.Items[0].ObjectMeta.CreationTimestamp,
		},
		Spec: platform.ClusterAddonSpec{
			Type:    string(clusteraddontype.TappController),
			Level:   clusteraddontype.Types[clusteraddontype.TappController].Level,
			Version: l.Items[0].Spec.Version,
		},
		Status: platform.ClusterAddonStatus{
			Version: l.Items[0].Status.Version,
			Phase:   string(l.Items[0].Status.Phase),
			Reason:  l.Items[0].Status.Reason,
		},
	})
	a.mutex.Unlock()
}

func csiOperator(ctx context.Context, a *addonFinder) {
	defer a.wg.Done()
	l, err := a.platformClient.CSIOperators().List(ctx, metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("spec.clusterName", a.clusterName).String(),
	})
	if err != nil {
		a.mutex.Lock()
		a.errors = append(a.errors, err)
		a.mutex.Unlock()
		return
	}
	if len(l.Items) == 0 {
		return
	}
	a.mutex.Lock()
	a.addons = append(a.addons, platform.ClusterAddon{
		ObjectMeta: metav1.ObjectMeta{
			Name:              l.Items[0].ObjectMeta.Name,
			CreationTimestamp: l.Items[0].ObjectMeta.CreationTimestamp,
		},
		Spec: platform.ClusterAddonSpec{
			Type:    string(clusteraddontype.CSIOperator),
			Level:   clusteraddontype.Types[clusteraddontype.CSIOperator].Level,
			Version: l.Items[0].Spec.Version,
		},
		Status: platform.ClusterAddonStatus{
			Version: l.Items[0].Status.Version,
			Phase:   string(l.Items[0].Status.Phase),
			Reason:  l.Items[0].Status.Reason,
		},
	})
	a.mutex.Unlock()
}

func cronHPA(ctx context.Context, a *addonFinder) {
	defer a.wg.Done()
	l, err := a.platformClient.CronHPAs().List(ctx, metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("spec.clusterName", a.clusterName).String(),
	})
	if err != nil {
		a.mutex.Lock()
		a.errors = append(a.errors, err)
		a.mutex.Unlock()
		return
	}
	if len(l.Items) == 0 {
		return
	}
	a.mutex.Lock()
	a.addons = append(a.addons, platform.ClusterAddon{
		ObjectMeta: metav1.ObjectMeta{
			Name:              l.Items[0].ObjectMeta.Name,
			CreationTimestamp: l.Items[0].ObjectMeta.CreationTimestamp,
		},
		Spec: platform.ClusterAddonSpec{
			Type:    string(clusteraddontype.CronHPA),
			Level:   clusteraddontype.Types[clusteraddontype.CronHPA].Level,
			Version: l.Items[0].Spec.Version,
		},
		Status: platform.ClusterAddonStatus{
			Version: l.Items[0].Status.Version,
			Phase:   string(l.Items[0].Status.Phase),
			Reason:  l.Items[0].Status.Reason,
		},
	})
	a.mutex.Unlock()
}
