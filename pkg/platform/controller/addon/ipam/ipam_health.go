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

package ipam

import (
	"sync"

	"tkestack.io/tke/api/platform/v1"
)

type ipamHealth struct {
	mu      sync.Mutex
	ipamMap map[string]*v1.IPAM
}

func (s *ipamHealth) Exist(ipamName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.ipamMap[ipamName]
	return ok
}

func (s *ipamHealth) Set(ipam *v1.IPAM) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ipamMap[ipam.Name] = ipam
}

func (s *ipamHealth) Del(ipamName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ipamMap, ipamName)
}
