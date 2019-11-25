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
	"fmt"
	"net/rpc"
	"runtime/debug"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	v1 "tkestack.io/tke/api/platform/v1"
)

// RPC is an implementation that talks over RPC
type RPC struct{ client *rpc.Client }

var _ Provider = &RPC{}

// Name returns the name of provider.
func (g *RPC) Name() (resp string, err error) {
	err = g.client.Call("Plugin.Name", new(interface{}), &resp)
	return
}

// Init to initialize the provider by given configuration file.
func (g *RPC) Init(configFile string) error {
	return g.client.Call("Plugin.Init",
		map[string]interface{}{
			"configFile": configFile,
		},
		new(interface{}),
	)
}

// Validate is used to respond to the input parameters of the cluster when the
// cluster is created.
func (g *RPC) Validate(cluster platform.Cluster) (resp field.ErrorList, err error) {
	err = g.client.Call("Plugin.Validate", map[string]interface{}{"cluster": cluster}, &resp)
	return
}

// PreCreate is used to respond to the input parameters of the verified cluster
// when the cluster is created, and will be written to the backend storage.
func (g *RPC) PreCreate(user UserInfo, cluster platform.Cluster) (resp platform.Cluster, err error) {
	err = g.client.Call("Plugin.PreCreate", map[string]interface{}{"user": user, "cluster": cluster}, &resp)
	return
}

// AfterCreate is used to respond to the time the cluster was created after it
// was written to the backend store.
func (g *RPC) AfterCreate(cluster platform.Cluster) (resp []interface{}, err error) {
	err = g.client.Call("Plugin.AfterCreate", map[string]interface{}{"cluster": cluster}, &resp)
	return
}

// ValidateUpdate is used to respond to the input parameters of the cluster when
// the cluster is updated.
func (g *RPC) ValidateUpdate(cluster platform.Cluster, oldCluster platform.Cluster) (resp field.ErrorList, err error) {
	err = g.client.Call("Plugin.ValidateUpdate", map[string]interface{}{"cluster": cluster, "oldCluster": oldCluster}, &resp)
	return
}

// OnInitialize is used to respond to the state of the cluster when the cluster
// is created.
func (g *RPC) OnInitialize(cluster Cluster) (resp Cluster, err error) {
	err = g.client.Call("Plugin.OnInitialize", map[string]interface{}{"cluster": cluster}, &resp)
	return
}

// OnUpdate is used to sync cluster status with spec
func (g *RPC) OnUpdate(cluster Cluster) (resp Cluster, err error) {
	err = g.client.Call("Plugin.OnUpdate", map[string]interface{}{"cluster": cluster}, &resp)
	return
}

// OnDelete is used to respond when the cluster will delete.
func (g *RPC) OnDelete(cluster v1.Cluster) error {
	return g.client.Call("Plugin.OnDelete", map[string]interface{}{"cluster": cluster}, new(interface{}))
}

// RPCServer is the RPC server that RPC talks to, conforming to the requirements of net/rpc
type RPCServer struct {
	Impl Provider
}

// Name returns the name of cluster provider.
func (s *RPCServer) Name(args interface{}, resp *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	*resp, err = s.Impl.Name()
	return
}

// Init to initialize the cluster provider by given configuration file.
func (s *RPCServer) Init(args map[string]interface{}, _ *interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	err = s.Impl.Init(args["configFile"].(string))
	return
}

// Validate is used to respond to the input parameters of the cluster when the
// cluster is created.
func (s *RPCServer) Validate(args map[string]interface{}, resp *field.ErrorList) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	*resp, err = s.Impl.Validate(args["cluster"].(platform.Cluster))
	return
}

// PreCreate is used to respond to the input parameters of the verified cluster
// when the cluster is created, and will be written to the backend storage.
func (s *RPCServer) PreCreate(args map[string]interface{}, resp *platform.Cluster) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	*resp, err = s.Impl.PreCreate(args["user"].(UserInfo), args["cluster"].(platform.Cluster))
	return
}

// AfterCreate is used to respond to the time the cluster was created after it
// was written to the backend store.
func (s *RPCServer) AfterCreate(args map[string]interface{}, resp *[]interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	*resp, err = s.Impl.AfterCreate(args["cluster"].(platform.Cluster))
	return
}

// ValidateUpdate is used to respond to the input parameters of the cluster when
// the cluster is updated.
func (s *RPCServer) ValidateUpdate(args map[string]interface{}, resp *field.ErrorList) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	*resp, err = s.Impl.ValidateUpdate(args["cluster"].(platform.Cluster), args["oldCluster"].(platform.Cluster))
	return
}

// OnInitialize is used to respond to the state of the cluster when the cluster
// is created.
func (s *RPCServer) OnInitialize(args map[string]interface{}, resp *Cluster) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	*resp, err = s.Impl.OnInitialize(args["cluster"].(Cluster))
	return
}

// OnUpdate is used to sync cluster status with spec
func (s *RPCServer) OnUpdate(args map[string]interface{}, resp *Cluster) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	*resp, err = s.Impl.OnUpdate(args["cluster"].(Cluster))
	return
}

// OnDelete is used to respond when the cluster will delete.
func (s *RPCServer) OnDelete(args map[string]interface{}, _ *interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	err = s.Impl.OnDelete(args["cluster"].(v1.Cluster))
	return
}
