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

package helm

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sync"
	"time"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	lister "tkestack.io/tke/api/client/listers/platform/v1"
	v1 "tkestack.io/tke/api/platform/v1"
	controllerUtil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/util/log"
)

const (
	checkingHcInterval  = 5 * time.Second
	runningHcInterval   = 2 * time.Minute
	unhealthyHcInternal = 20 * time.Second
	failedHcInternal    = 5 * time.Minute
	checkingTimeout     = 5 * time.Minute
	unhealthyTimeout    = 5 * time.Minute
)

// Prober is used for probing the instance status
type Prober interface {
	// Run means starts to run prober
	Run(ch <-chan struct{})
	// Exist means whether key is in the prober
	Exist(key string) bool
	// ExistByPhase means whether key is in the prober for the given Addon
	ExistByPhase(key string, phase v1.AddonPhase) bool
	// Set means sets a prober for the given Addon
	Set(key string, phase v1.AddonPhase)
	// Del means remove a prober of the given key
	Del(key string)
}

type prober struct {
	checking  sync.Map
	running   sync.Map
	unhealthy sync.Map
	failed    sync.Map
	lister    lister.HelmLister
	client    clientset.Interface
}

// NewHealthProber returns a prober
func NewHealthProber(lister lister.HelmLister, client clientset.Interface) Prober {
	prober := &prober{
		lister: lister,
		client: client,
	}
	return prober
}

func (p *prober) Exist(key string) bool {
	if p.ExistByPhase(key, v1.AddonPhaseChecking) || p.ExistByPhase(key, v1.AddonPhaseRunning) ||
		p.ExistByPhase(key, v1.AddonPhaseFailed) || p.ExistByPhase(key, v1.AddonPhaseUnhealthy) {
		return true
	}
	return false
}

// ExistByPhase tells whether key is in the prober according to the phase
func (p *prober) ExistByPhase(key string, phase v1.AddonPhase) bool {
	if _, ok := p.fetchSyncMap(phase).Load(key); !ok {
		return false
	}
	return true
}

func (p *prober) Set(key string, phase v1.AddonPhase) {
	p.fetchSyncMap(phase).Store(key, true)
}

func (p *prober) Del(key string) {
	p.delByPhase(key, v1.AddonPhaseChecking)
	p.delByPhase(key, v1.AddonPhaseRunning)
	p.delByPhase(key, v1.AddonPhaseFailed)
	p.delByPhase(key, v1.AddonPhaseUnhealthy)
}

func (p *prober) delByPhase(key string, phase v1.AddonPhase) {
	p.fetchSyncMap(phase).Delete(key)
}

func (p *prober) Run(ch <-chan struct{}) {
	go p.runProbe(v1.AddonPhaseChecking, checkingHcInterval, ch)
	go p.runProbe(v1.AddonPhaseRunning, runningHcInterval, ch)
	go p.runProbe(v1.AddonPhaseUnhealthy, unhealthyHcInternal, ch)
	go p.runProbe(v1.AddonPhaseFailed, failedHcInternal, ch)
}

func (p *prober) runProbe(phase v1.AddonPhase, duration time.Duration, ch <-chan struct{}) {
	log.Info(fmt.Sprintf("Helm %s prober begin running", phase))
	defer controllerUtil.CatchPanic(fmt.Sprintf("%s prober", phase), "Helm")

	p.periodHealthCheck(phase, duration, ch)
}

func (p *prober) periodHealthCheck(phase v1.AddonPhase, duration time.Duration, ch <-chan struct{}) {
	wait.Until(func() {
		p.fetchSyncMap(phase).Range(
			func(k, v interface{}) bool {
				key := k.(string)
				log.Info(fmt.Sprintf("Probe check"), log.String("helm", key), log.String("phase", string(phase)))
				curPhase, err := p.checkAddonHealthz(key, phase)
				if curPhase == "" || (curPhase != phase && err == nil) {
					p.fetchSyncMap(phase).Delete(key)
					return true
				}
				return true
			})
	}, duration, ch)
}

func (p *prober) fetchSyncMap(phase v1.AddonPhase) *sync.Map {
	switch phase {
	case v1.AddonPhaseRunning:
		return &p.running
	case v1.AddonPhaseChecking:
		return &p.checking
	case v1.AddonPhaseUnhealthy:
		return &p.unhealthy
	case v1.AddonPhaseFailed:
		return &p.failed
	default:
		log.Info(fmt.Sprintf("set unrecognized phase %v", phase))
	}
	return &sync.Map{}
}

func (p *prober) checkAddonHealthz(key string, oldPhase v1.AddonPhase) (v1.AddonPhase, error) {
	// 判断 key 没有被删除
	obj, err := p.lister.Get(key)
	if err != nil {
		// delete key from map
		return "", err
	}
	// 创建 provisioner 与用户集群交互
	var provisioner Provisioner
	if provisioner, err = createProvisioner(obj, p.client); err != nil {
		if errors.IsNotFound(err) {
			uptErr := updateHelmStatus(getUpdateObj(obj, v1.AddonPhaseFailed, err.Error()), p.client)
			return v1.AddonPhaseFailed, uptErr
		}
		if isCheckOrUnhealthyTimeout(obj.Status.LastReInitializingTimestamp.Time, oldPhase) {
			if err.Error() == "" {
				uptErr := updateHelmStatus(getUpdateObj(obj, v1.AddonPhaseFailed, string(oldPhase)+" timeout: fail to connect to tiller"), p.client)
				return v1.AddonPhaseFailed, uptErr
			}
			uptErr := updateHelmStatus(getUpdateObj(obj, v1.AddonPhaseFailed, err.Error()), p.client)
			return v1.AddonPhaseFailed, uptErr
		}
		if oldPhase == v1.AddonPhaseUnhealthy && obj.Status.Reason != err.Error() {
			uptErr := updateHelmStatus(getUpdateObj(obj, v1.AddonPhaseUnhealthy, err.Error()), p.client)
			return oldPhase, uptErr
		}
		return v1.AddonPhaseUnhealthy, nil
	}
	// 获取 addon 状态
	err = provisioner.GetStatus()
	if err != nil && errors.IsNotFound(err) {
		uptErr := updateHelmStatus(getUpdateObj(obj, v1.AddonPhaseFailed, err.Error()), p.client)
		return v1.AddonPhaseFailed, uptErr
	}
	if err != nil {
		// checking timeout or unhealthy timeout
		if isCheckOrUnhealthyTimeout(obj.Status.LastReInitializingTimestamp.Time, oldPhase) {
			if err.Error() == "" {
				uptErr := updateHelmStatus(getUpdateObj(obj, v1.AddonPhaseFailed, string(oldPhase)+" timeout: fail to connect to tiller"), p.client)
				return v1.AddonPhaseFailed, uptErr
			}
			uptErr := updateHelmStatus(getUpdateObj(obj, v1.AddonPhaseFailed, err.Error()), p.client)
			return v1.AddonPhaseFailed, uptErr
		}
		// running -> unhealthy
		if oldPhase == v1.AddonPhaseRunning {
			newObj := getUpdateObj(obj, v1.AddonPhaseUnhealthy, err.Error())
			newObj.Status.LastReInitializingTimestamp = metav1.Now()
			uptErr := updateHelmStatus(newObj, p.client)
			return v1.AddonPhaseUnhealthy, uptErr
		}
		// change reason, stay old phase
		if obj.Status.Reason != err.Error() {
			uptErr := updateHelmStatus(getUpdateObj(obj, oldPhase, err.Error()), p.client)
			return oldPhase, uptErr
		}
		return oldPhase, nil
	}
	// other phases -> running
	if oldPhase != v1.AddonPhaseRunning {
		uptErr := updateHelmStatus(getUpdateObj(obj, v1.AddonPhaseRunning, ""), p.client)
		return v1.AddonPhaseRunning, uptErr
	}
	// stay running
	return v1.AddonPhaseRunning, nil
}

func isCheckOrUnhealthyTimeout(recordTime time.Time, oldPhase v1.AddonPhase) bool {
	if (oldPhase == v1.AddonPhaseChecking && time.Since(recordTime) >= checkingTimeout) ||
		(oldPhase == v1.AddonPhaseUnhealthy && time.Since(recordTime) >= unhealthyTimeout) {
		return true
	}
	return false
}
