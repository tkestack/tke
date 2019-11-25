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

package app

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	cacheddiscovery "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/restmapper"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	"tkestack.io/tke/cmd/tke-platform-controller/app/config"
	"tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/controller/util"
)

// InitFunc is used to launch a particular controller.  It may run additional "should I activate checks".
// Any error returned will cause the controller process to `Fatal`
// The bool indicates whether the controller was enabled.
type InitFunc func(ctx ControllerContext) (debuggingHandler http.Handler, enabled bool, err error)

// ControllerContext represents the context of controller.
type ControllerContext struct {
	// ClientBuilder will provide a client for this controller to use
	ClientBuilder controller.ClientBuilder

	// InformerFactory gives access to informers for the controller.
	InformerFactory versionedinformers.SharedInformerFactory

	// DeferredDiscoveryRESTMapper is a RESTMapper that will defer
	// initialization of the RESTMapper until the first mapping is
	// requested.
	RESTMapper *restmapper.DeferredDiscoveryRESTMapper

	// AvailableResources is a map listing currently available resources
	AvailableResources map[schema.GroupVersionResource]bool

	// Stop is the stop channel
	Stop <-chan struct{}

	// InformersStarted is closed after all of the controllers have been initialized and are running.  After this point it is safe,
	// for an individual controller to start the shared informers. Before it is closed, they should not.
	InformersStarted chan struct{}

	// ResyncPeriod generates a duration each time it is invoked; this is so that
	// multiple controllers don't get into lock-step and all hammer the apiserver
	// with list requests simultaneously.
	ResyncPeriod            func() time.Duration
	ControllerStartInterval time.Duration
	ClusterProviders        *sync.Map
	MachineProviders        *sync.Map

	// Remote write/read address for prometheus
	RemoteAddresses []string
	RemoteType      string
}

// IsControllerEnabled returns whether the controller has been enabled
func (c ControllerContext) IsControllerEnabled(name string) bool {
	return util.IsControllerEnabled(name, ControllersDisabledByDefault)
}

// CreateControllerContext creates a context struct containing references to resources needed by the
// controllers such as the cloud provider and clientBuilder. rootClientBuilder is only used for
// the shared-informers client and token controller.
func CreateControllerContext(cfg *config.Config, rootClientBuilder controller.ClientBuilder, stop <-chan struct{}) (ControllerContext, error) {
	versionedClient := rootClientBuilder.ClientOrDie("shared-informers")
	sharedInformers := versionedinformers.NewSharedInformerFactory(versionedClient, controller.ResyncPeriod(&cfg.Component)())

	// If apiserver is not running we should wait for some time and fail only then. This is particularly
	// important when we start apiserver and controller manager at the same time.
	if err := controller.WaitForAPIServer(versionedClient, 10*time.Second); err != nil {
		return ControllerContext{}, fmt.Errorf("failed to wait for apiserver being healthy: %v", err)
	}

	// Use a discovery client capable of being refreshed.
	discoveryClient := rootClientBuilder.ClientOrDie("controller-discovery")
	cachedClient := cacheddiscovery.NewMemCacheClient(discoveryClient.Discovery())
	restMapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedClient)
	go wait.Until(func() {
		restMapper.Reset()
	}, 5*time.Minute, stop)

	availableResources, err := controller.GetAvailableResources(rootClientBuilder)
	if err != nil {
		return ControllerContext{}, err
	}

	ctx := ControllerContext{
		ClientBuilder:           rootClientBuilder,
		InformerFactory:         sharedInformers,
		RESTMapper:              restMapper,
		AvailableResources:      availableResources,
		Stop:                    stop,
		InformersStarted:        make(chan struct{}),
		ResyncPeriod:            controller.ResyncPeriod(&cfg.Component),
		ControllerStartInterval: cfg.Component.ControllerStartInterval,
		ClusterProviders:        cfg.Provider.ClusterProviders,
		MachineProviders:        cfg.Provider.MachineProviders,
		RemoteAddresses:         cfg.Features.MonitorStorageAddresses,
		RemoteType:              cfg.Features.MonitorStorageType,
	}
	return ctx, nil
}
