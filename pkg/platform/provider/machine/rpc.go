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

package machine

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

// Validate is used to respond to the input parameters of the machine when the
// machine is created.
func (g *RPC) Validate(machine platform.Machine) (resp field.ErrorList, err error) {
	err = g.client.Call("Plugin.Validate", map[string]interface{}{"machine": machine}, &resp)
	return
}

// OnInitialize is used to respond to the state of the machine when the machine
// is created.
func (g *RPC) OnInitialize(machine v1.Machine, cluster v1.Cluster, credential v1.ClusterCredential) (resp v1.Machine, err error) {
	err = g.client.Call("Plugin.OnInitialize",
		map[string]interface{}{
			"machine":    machine,
			"cluster":    cluster,
			"credential": credential,
		},
		&resp)
	return
}

// RPCServer is the RPC server that RPC talks to, conforming to the requirements of net/rpc
type RPCServer struct {
	Impl Provider
}

func (s *RPCServer) Name(_ interface{}, resp *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	*resp, err = s.Impl.Name()
	return
}

func (s *RPCServer) Init(args map[string]interface{}, _ *interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	err = s.Impl.Init(args["configFile"].(string))
	return
}

func (s *RPCServer) Validate(args map[string]interface{}, resp *field.ErrorList) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	*resp, err = s.Impl.Validate(args["machine"].(platform.Machine))
	return
}

func (s *RPCServer) OnInitialize(args map[string]interface{}, resp *v1.Machine) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[recovery] panic recovered:\nMessage: %s\n%s", r, debug.Stack())
		}
	}()
	*resp, err = s.Impl.OnInitialize(args["machine"].(v1.Machine), args["cluster"].(v1.Cluster), args["credential"].(v1.ClusterCredential))
	return
}
