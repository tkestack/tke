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

// Code generated by client-gen. DO NOT EDIT.

package internalversion

import (
	"fmt"

	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
	authinternalversion "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	businessinternalversion "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	logagentinternalversion "tkestack.io/tke/api/client/clientset/internalversion/typed/logagent/internalversion"
	monitorinternalversion "tkestack.io/tke/api/client/clientset/internalversion/typed/monitor/internalversion"
	notifyinternalversion "tkestack.io/tke/api/client/clientset/internalversion/typed/notify/internalversion"
	platforminternalversion "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	registryinternalversion "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	Auth() authinternalversion.AuthInterface
	Business() businessinternalversion.BusinessInterface
	Logagent() logagentinternalversion.LogagentInterface
	Monitor() monitorinternalversion.MonitorInterface
	Notify() notifyinternalversion.NotifyInterface
	Platform() platforminternalversion.PlatformInterface
	Registry() registryinternalversion.RegistryInterface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	auth     *authinternalversion.AuthClient
	business *businessinternalversion.BusinessClient
	logagent *logagentinternalversion.LogagentClient
	monitor  *monitorinternalversion.MonitorClient
	notify   *notifyinternalversion.NotifyClient
	platform *platforminternalversion.PlatformClient
	registry *registryinternalversion.RegistryClient
}

// Auth retrieves the AuthClient
func (c *Clientset) Auth() authinternalversion.AuthInterface {
	return c.auth
}

// Business retrieves the BusinessClient
func (c *Clientset) Business() businessinternalversion.BusinessInterface {
	return c.business
}

// Logagent retrieves the LogagentClient
func (c *Clientset) Logagent() logagentinternalversion.LogagentInterface {
	return c.logagent
}

// Monitor retrieves the MonitorClient
func (c *Clientset) Monitor() monitorinternalversion.MonitorInterface {
	return c.monitor
}

// Notify retrieves the NotifyClient
func (c *Clientset) Notify() notifyinternalversion.NotifyInterface {
	return c.notify
}

// Platform retrieves the PlatformClient
func (c *Clientset) Platform() platforminternalversion.PlatformInterface {
	return c.platform
}

// Registry retrieves the RegistryClient
func (c *Clientset) Registry() registryinternalversion.RegistryInterface {
	return c.registry
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.auth, err = authinternalversion.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.business, err = businessinternalversion.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.logagent, err = logagentinternalversion.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.monitor, err = monitorinternalversion.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.notify, err = notifyinternalversion.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.platform, err = platforminternalversion.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.registry, err = registryinternalversion.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.auth = authinternalversion.NewForConfigOrDie(c)
	cs.business = businessinternalversion.NewForConfigOrDie(c)
	cs.logagent = logagentinternalversion.NewForConfigOrDie(c)
	cs.monitor = monitorinternalversion.NewForConfigOrDie(c)
	cs.notify = notifyinternalversion.NewForConfigOrDie(c)
	cs.platform = platforminternalversion.NewForConfigOrDie(c)
	cs.registry = registryinternalversion.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.auth = authinternalversion.New(c)
	cs.business = businessinternalversion.New(c)
	cs.logagent = logagentinternalversion.New(c)
	cs.monitor = monitorinternalversion.New(c)
	cs.notify = notifyinternalversion.New(c)
	cs.platform = platforminternalversion.New(c)
	cs.registry = registryinternalversion.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
