/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
	"sync"

	typescfg "tkestack.io/tke/pkg/mesh/external/tcmesh/types/config"

	"tkestack.io/tke/pkg/util/log"
)

var (
	mcCacher    *Cacher
	mcCacheOnce sync.Once
)

const (
	ClusterRoleMaster = "master"
	ClusterRoleRemote = "remote"
)

type Cacher struct {
	meshToClustersMap        map[string]*[]string
	meshToMainClustersMap    map[string][]string
	clusterToMeshMap         map[string]string
	clusterToMainClustersMap map[string][]string
	mcMappingsMutex          sync.RWMutex
	ctx                      context.Context
}

var client *Client

func GetCacher(ctx context.Context, tcmClient *Client) *Cacher {
	mcCacheOnce.Do(func() {
		if mcCacher == nil {
			mcCacher = &Cacher{
				meshToClustersMap:        map[string]*[]string{},
				meshToMainClustersMap:    map[string][]string{},
				clusterToMeshMap:         map[string]string{},
				clusterToMainClustersMap: map[string][]string{},
				ctx:                      ctx,
				mcMappingsMutex:          sync.RWMutex{},
			}
		}
		if client == nil {
			client = tcmClient
		}
	})
	return mcCacher
}

func (c *Cacher) ReloadMeshes() {
	log.Infof("Ready to sync meshes from mesh-manager")
	meshes, err := client.GetMeshes(c.ctx)

	if err != nil {
		log.Errorf("Error occur while sync meshes from mesh-manager. Detail: %v", err)
		return
	}

	if meshes == nil {
		return
	}

	meshToClustersMap := map[string]*[]string{}
	meshToMainClusterMap := map[string][]string{}
	clusterToMeshMap := map[string]string{}
	clusterToMainClusterMap := map[string][]string{}

	for _, mesh := range meshes {

		for _, cluster := range mesh.Config.Clusters {
			if clustersList, exists := meshToClustersMap[mesh.Name]; exists {
				*clustersList = append(*clustersList, cluster.Name)
			} else {
				meshToClustersMap[mesh.Name] = &[]string{cluster.Name}
			}

			if c.isMaster(cluster) {
				if mainList, exists := meshToMainClusterMap[mesh.Name]; exists {
					meshToMainClusterMap[mesh.Name] = append(mainList, cluster.Name)
				} else {
					meshToMainClusterMap[mesh.Name] = []string{cluster.Name}
				}
			}
			clusterToMeshMap[cluster.Name] = mesh.Name

		}
	}

	for clusterName, meshName := range clusterToMeshMap {
		mainClusterName := meshToMainClusterMap[meshName]
		clusterToMainClusterMap[clusterName] = mainClusterName
	}

	c.mcMappingsMutex.Lock()
	defer c.mcMappingsMutex.Unlock()

	c.meshToClustersMap = meshToClustersMap
	c.meshToMainClustersMap = meshToMainClusterMap
	c.clusterToMeshMap = clusterToMeshMap
	c.clusterToMainClustersMap = clusterToMainClusterMap

	log.Debugf("Sync meshes end. \n"+
		"meshToClusters --> %v;\n"+
		"meshToMainClusters --> %v;\n"+
		"clusterToMesh --> %v;\n"+
		"clusterToMainClusters --> %v;\n",
		c.meshToClustersMap,
		c.meshToMainClustersMap,
		c.clusterToMeshMap,
		c.clusterToMainClustersMap)
}

//func (c *Cacher) preSync() {
//	log.Infof("Pre-sync mesh info while starting...")
//	c.ReloadMeshes()
//}

func (c *Cacher) Clusters(meshName string) []string {
	clusters := c.getClusters(meshName)

	if len(clusters) != 0 {
		return clusters
	}
	c.ReloadMeshes()
	return c.getClusters(meshName)
}

func (c *Cacher) getClusters(meshName string) []string {
	c.mcMappingsMutex.RLock()
	defer c.mcMappingsMutex.RUnlock()
	if clusters, hit := c.meshToClustersMap[meshName]; hit {
		return *clusters
	}
	return nil
}

func (c *Cacher) MainClusters(meshName string) []string {
	mainCluster := c.getMainClusters(meshName)

	if len(mainCluster) > 0 {
		return mainCluster
	}
	c.ReloadMeshes()
	return c.getMainClusters(meshName)
}

func (c *Cacher) MainClustersMap(meshName string) map[string]struct{} {
	ret := make(map[string]struct{})
	mainClusters := c.MainClusters(meshName)
	for _, cls := range mainClusters {
		ret[cls] = struct{}{}
	}
	return ret
}

func (c *Cacher) getMainClusters(meshName string) []string {
	c.mcMappingsMutex.RLock()
	defer c.mcMappingsMutex.RUnlock()
	if cluster, hit := c.meshToMainClustersMap[meshName]; hit {
		return cluster
	}
	return []string{}
}

func (c *Cacher) MemberClusters(meshName string) []string {
	clusters := c.getMemberClusters(meshName)

	if len(clusters) > 0 {
		return clusters
	}
	c.ReloadMeshes()
	return c.getMemberClusters(meshName)
}

func (c *Cacher) getMemberClusters(meshName string) []string {
	c.mcMappingsMutex.RLock()
	defer c.mcMappingsMutex.RUnlock()
	if clusters, hit := c.meshToClustersMap[meshName]; hit {
		return *clusters
	}
	return []string{}
}

func (c *Cacher) GetMeshByCluster(clusterName string) string {
	mesh := c.getMeshByCluster(clusterName)

	if mesh != "" {
		return mesh
	}
	c.ReloadMeshes()
	return c.getMeshByCluster(clusterName)
}

func (c *Cacher) getMeshByCluster(clusterName string) string {
	c.mcMappingsMutex.RLock()
	defer c.mcMappingsMutex.RUnlock()
	if mesh, hit := c.clusterToMeshMap[clusterName]; hit {
		return mesh
	}
	return ""
}

func (c *Cacher) GetMainClusterByMemberCluster(clusterName string) []string {
	mainCluster := c.getMainClusterByMemberCluster(clusterName)

	if len(mainCluster) > 0 {
		return mainCluster
	}
	c.ReloadMeshes()
	return c.getMainClusterByMemberCluster(clusterName)
}

func (c *Cacher) getMainClusterByMemberCluster(clusterName string) []string {
	c.mcMappingsMutex.RLock()
	defer c.mcMappingsMutex.RUnlock()
	if cluster, hit := c.clusterToMainClustersMap[clusterName]; hit {
		return cluster
	}
	return []string{}
}

func (c *Cacher) isMaster(cluster *typescfg.ClusterConfig) bool {
	return cluster.Role == ClusterRoleMaster
}
