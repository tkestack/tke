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
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	v1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
)

const conditionTypeHealthCheck = "HealthCheck"
const reasonHealthCheckFail = "HealthCheckFail"

type machineHealth struct {
	mu sync.Mutex
	m  map[string]*v1.Machine
}

func (s *machineHealth) Exist(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.m[name]
	return ok
}

func (s *machineHealth) Set(machine *v1.Machine) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[machine.Name] = machine
}

func (s *machineHealth) Del(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, name)
}

func (c *Controller) startMachineHealthCheck(key string, machine *v1.Machine) error {
	if !c.health.Exist(key) {
		c.health.Set(machine)
		go wait.PollImmediateUntil(1*time.Minute, c.watchHealth(machine.Name), c.stopCh)
		log.Info("Machine phase start new health check", log.String("name", key), log.String("phase", string(machine.Status.Phase)))
	} else {
		log.Info("Machine phase health check exit", log.String("name", key), log.String("phase", string(machine.Status.Phase)))
	}

	return nil
}

func (c *Controller) checkHealth(m *v1.Machine) error {
	clientset, err := util.BuildExternalClientSetWithName(c.client.PlatformV1(), m.Spec.ClusterName)
	if err != nil {
		m.Status.Phase = v1.MachineFailed
		m.Status.Message = err.Error()
		m.Status.Reason = reasonHealthCheckFail
		now := metav1.Now()
		c.addOrUpdateCondition(m, v1.MachineCondition{
			Type:               conditionTypeHealthCheck,
			Status:             v1.ConditionFalse,
			Message:            err.Error(),
			Reason:             reasonHealthCheckFail,
			LastTransitionTime: now,
			LastProbeTime:      now,
		})
		if err := c.persistUpdate(m); err != nil {
			log.Warn("Update machine status failed", log.String("name", m.Name), log.Err(err))
			return err
		}

		return err
	}

	node, err := clientset.CoreV1().Nodes().Get(m.Spec.IP, metav1.GetOptions{})
	if err != nil {
		m.Status.Phase = v1.MachineFailed
		m.Status.Message = err.Error()
		m.Status.Reason = reasonHealthCheckFail
		now := metav1.Now()
		c.addOrUpdateCondition(m, v1.MachineCondition{
			Type:               conditionTypeHealthCheck,
			Status:             v1.ConditionFalse,
			Message:            err.Error(),
			Reason:             reasonHealthCheckFail,
			LastTransitionTime: now,
			LastProbeTime:      now,
		})
		if err := c.persistUpdate(m); err != nil {
			log.Warn("Update machine status failed", log.String("name", m.Name), log.Err(err))
			return err
		}

		return err
	}

	m.Status.Phase = v1.MachineRunning
	m.Status.Message = ""
	m.Status.Reason = ""
	m.Status.MachineInfo = v1.MachineSystemInfo(node.Status.NodeInfo)
	c.addOrUpdateCondition(m, v1.MachineCondition{
		Type:          conditionTypeHealthCheck,
		Status:        v1.ConditionTrue,
		Message:       "",
		Reason:        "",
		LastProbeTime: metav1.Now(),
	})
	if err := c.persistUpdate(m); err != nil {
		log.Warn("Update machine status failed", log.String("name", m.Name), log.Err(err))
		return err
	}

	return nil
}

// for PollImmediateUntil, when return true ,an err while exit
func (c *Controller) watchHealth(name string) func() (bool, error) {
	return func() (bool, error) {
		log.Debug("Check machine health", log.String("name", name))

		m, err := c.client.PlatformV1().Machines().Get(name, metav1.GetOptions{})
		if err != nil {
			// if machine is not found,to exit the health check loop
			if errors.IsNotFound(err) {
				log.Warn("Machine not found, to exit the health check loop", log.String("name", name))
				return true, nil
			}
			log.Error("Check machine health, machine get failed", log.String("name", name), log.Err(err))
			return false, nil
		}

		// if status is terminated,to exit the  health check loop
		if m.Status.Phase == v1.MachineTerminating || m.Status.Phase == v1.MachineInitializing {
			log.Warn("Machine status is terminated, to exit the health check loop", log.String("name", name))
			return true, nil
		}

		_ = c.checkHealth(m)
		return false, nil
	}
}
