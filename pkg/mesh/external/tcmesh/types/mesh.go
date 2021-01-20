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

package types

import (
	"time"

	"tkestack.io/tke/pkg/mesh/external/tcmesh/types/config"
	"tkestack.io/tke/pkg/mesh/external/tcmesh/types/state"
)

type Mesh struct {
	ID          int64             `json:"-"`
	TenantID    string            `json:"-"`
	Name        string            `json:"name"`
	Region      string            `json:"region"`
	Title       string            `json:"title"`
	Status      string            `json:"status"`
	Version     string            `json:"version"`
	LastVersion string            `json:"-"`
	Config      config.MeshConfig `json:"config"`
	Mode        string            `json:"mode"`
	State       state.MeshState   `json:"state"`
	CreatedAt   *time.Time        `json:"created_at"`
	UpdatedAt   *time.Time        `json:"updated_at"`
}

type Cluster struct {
	ID        int64  `json:"-"`
	Name      string `json:"name"`
	Region    string `json:"region"`
	MeshName  string
	Role      string     `json:"role"`
	Phase     string     `json:"phase"`
	CreatedAt *time.Time `json:"-"`
	LinkedAt  *time.Time `json:"linked_at"`
	Error     string     `json:"error"`
}

type LBResource struct {
	ID       int64
	Name     string
	MeshName string
	Status   string
}

type MeshSelector struct {
	Region string     `json:"region"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
	Filter meshFilter `json:"filter"`
}

type meshFilter struct {
	Names        []string `json:"names"`
	Titles       []string `json:"titles"`
	ClusterNames []string `json:"clusterNames"`
}
