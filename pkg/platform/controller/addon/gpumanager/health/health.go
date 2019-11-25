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

package health

import (
	"fmt"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"sync"
	"time"
	gmlister "tkestack.io/tke/api/client/listers/platform/v1"
	"tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/controller/addon/gpumanager/operator"
	"tkestack.io/tke/pkg/platform/controller/addon/gpumanager/utils"
	"tkestack.io/tke/pkg/util/log"
)

type gmProber struct {
	mu     sync.Mutex
	store  sets.String
	lister gmlister.GPUManagerLister
	op     operator.ObjectOperator
}

// NewHealthProber returns a prober for GPUManager
func NewHealthProber(lister gmlister.GPUManagerLister, operator operator.ObjectOperator) Prober {
	prober := &gmProber{
		store:  sets.NewString(),
		lister: lister,
		op:     operator,
	}

	return prober
}

// Exist tells you whether GPUManager's key is in the prober
func (p *gmProber) Exist(key string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.store.Has(key)
}

func (p *gmProber) Set(v *v1.GPUManager) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.store.Insert(v.Name)
}

func (p *gmProber) Del(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.store.Delete(key)
}

func (p *gmProber) getAllKeys() []string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.store.UnsortedList()
}

func (p *gmProber) Run(ch <-chan struct{}) {
	go wait.Until(func() {
		for _, key := range p.getAllKeys() {
			log.Info(fmt.Sprintf("Start to probe GPUManager %s in cluster", key))
			gm, err := p.lister.Get(key)
			if err != nil {
				log.Error(fmt.Sprintf("can't find GPUManager %s in cluster, err %s", key, err))
				return
			}

			ds, err := p.op.GetDaemonSet(gm)
			if err != nil && !apierror.IsNotFound(err) {
				log.Error(fmt.Sprintf("can't find GPUManager %s daemonsets, err %s", key, err))
				return
			}

			// Don't resubmit a daemonset if user has delete it, just mark it to failed
			if err != nil {
				_ = p.op.UpdateGPUManagerStatus(utils.ForUpdateItem(gm, v1.AddonPhaseFailed, "daemonset is not found"))
				return
			}

			availableNum, unavailableNum, desireNum := ds.Status.NumberAvailable, ds.Status.NumberUnavailable, ds.Status.DesiredNumberScheduled
			log.Info(fmt.Sprintf("GPUManager %s available: %d, unavaliable: %d, desired: %d", key, availableNum, unavailableNum, desireNum))

			switch gm.Status.Phase {
			case v1.AddonPhaseUpgrading, v1.AddonPhaseUnhealthy:
				if availableNum == desireNum {
					_ = p.op.UpdateGPUManagerStatus(utils.ForUpdateItem(gm, v1.AddonPhaseRunning, ""))
				}
			case v1.AddonPhaseRunning:
				if unavailableNum > 0 {
					_ = p.op.UpdateGPUManagerStatus(utils.ForUpdateItem(gm, v1.AddonPhaseUnhealthy, ""))
				}
			default:
			}
		}
	}, time.Second*10, ch)
}
