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
 *
 */

package tcmesh

import (
	"context"
	"fmt"

	restclient "k8s.io/client-go/rest"
	meshconfig "tkestack.io/tke/pkg/mesh/apis/config"
	"tkestack.io/tke/pkg/mesh/external/tcmesh/types"
	"tkestack.io/tke/pkg/mesh/util/simplerest"
	"tkestack.io/tke/pkg/util/log"
)

const (
	TcmAPIPrefix = "/api/meshes"
)

var (
	TcmAuthHeader = tcmHeader{
		Key:   "X-Remote-TenantID",
		Value: "default",
	}
)

type tcmHeader struct {
	Key   string
	Value string
}

type Client struct {
	client simplerest.Interface
	cache  *Cacher
}

func BuildClientConfig(config *meshconfig.MeshManagerConfig) (*restclient.Config, error) {
	return &restclient.Config{
		Host: config.Address,
	}, nil
}

func NewForConfig(c *restclient.Config) (*Client, error) {

	config := *c

	client, err := simplerest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	tcmClient := &Client{client: client}
	tcmClient.cache = GetCacher(context.TODO(), tcmClient)

	return tcmClient, nil
}

func (c *Client) RESTClient() simplerest.Interface {
	if c == nil {
		return nil
	}
	return c.client
}

func (c *Client) Cache() *Cacher {
	return c.cache
}

func (c *Client) ReloadMeshes() {
	c.cache.ReloadMeshes()
}

func (c *Client) GetMeshes(ctx context.Context) ([]types.Mesh, error) {
	uri := TcmAPIPrefix
	var response []types.Mesh
	err := c.RESTClient().
		Get().
		SetHeader(TcmAuthHeader.Key, TcmAuthHeader.Value).
		RequestURI(uri).
		Do(ctx).
		Into(&response)
	if err != nil {
		log.Errorf(fmt.Sprintf("Error getting mesh cluster list from mesh-manager. Err: %v", err))
		return nil, err
	}

	return response, nil
}

func (c *Client) GetMesh(ctx context.Context, meshName string) (*types.Mesh, error) {
	uri := TcmAPIPrefix + "/" + meshName
	var response types.Mesh
	err := c.RESTClient().
		Get().
		SetHeader(TcmAuthHeader.Key, TcmAuthHeader.Value).
		RequestURI(uri).
		Do(ctx).
		Into(&response)
	if err != nil {
		log.Errorf(fmt.Sprintf("Error getting mesh cluster [%s] from mesh-manager. Err: %v", meshName, err))
		return nil, err
	}

	return &response, nil
}
