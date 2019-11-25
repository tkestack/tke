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
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type cachedCluster struct {
	// The cached state of the cluster
	state *v1.Cluster
}

type clusterCache struct {
	mu         sync.Mutex // protects clusterMap
	clusterMap map[string]*cachedCluster
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *clusterCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.clusterMap))
	for k := range s.clusterMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the clusterMap under the given key
func (s *clusterCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.clusterMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *clusterCache) get(clusterName string) (*cachedCluster, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cluster, ok := s.clusterMap[clusterName]
	return cluster, ok
}

func (s *clusterCache) Exist(clusterName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.clusterMap[clusterName]
	return ok
}

func (s *clusterCache) getOrCreate(clusterName string) *cachedCluster {
	s.mu.Lock()
	defer s.mu.Unlock()
	cluster, ok := s.clusterMap[clusterName]
	if !ok {
		cluster = &cachedCluster{}
		s.clusterMap[clusterName] = cluster
	}
	return cluster
}

func (s *clusterCache) set(clusterName string, cluster *cachedCluster) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clusterMap[clusterName] = cluster
}

func (s *clusterCache) delete(clusterName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clusterMap, clusterName)
}
