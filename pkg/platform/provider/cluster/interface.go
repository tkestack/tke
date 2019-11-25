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

package cluster

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	v1 "tkestack.io/tke/api/platform/v1"
)

const (
	pluginName = "clusterProvider"
)

// handshakeConfig are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "CLUSTER_PROVIDER_PLUGIN",
	MagicCookieValue: "tke",
}

// pluginMaps is the map of plugins we can dispense.
var pluginMaps = map[string]plugin.Plugin{
	pluginName: &Plugin{},
}

// UserInfo describes a user that has been authenticated to the system.
type UserInfo struct {
	// Name returns the name that uniquely identifies this user among all
	// other active users.
	Name string
	// Groups returns the names of the groups the user is a member of
	Groups []string
	// GetExtra can contain any additional information that the authenticator
	// thought was interesting.  One example would be scopes on a token.
	// Keys in this map should be namespaced to the authenticator or
	// authenticator/authorizer pair making use of them.
	Extra map[string][]string
}

// Cluster wrap cluster and it's credential
type Cluster struct {
	v1.Cluster
	ClusterCredential v1.ClusterCredential
}

// Provider defines a set of response interfaces for specific cluster
// types in cluster management.
type Provider interface {
	Name() (string, error)
	Init(configFile string) error
	Validate(cluster platform.Cluster) (field.ErrorList, error)
	PreCreate(user UserInfo, cluster platform.Cluster) (platform.Cluster, error)
	AfterCreate(cluster platform.Cluster) ([]interface{}, error)
	ValidateUpdate(cluster platform.Cluster, oldCluster platform.Cluster) (field.ErrorList, error)
	OnInitialize(cluster Cluster) (Cluster, error)
	OnUpdate(cluster Cluster) (Cluster, error)
	OnDelete(cluster v1.Cluster) error
}

// Plugin is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a RPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return RPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type Plugin struct {
	Impl Provider
}

// Server returns the server side implement.
func (p *Plugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

// Client returns the client side implement.
func (Plugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPC{client: c}, nil
}
