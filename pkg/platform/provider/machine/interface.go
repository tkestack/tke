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
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/thoas/go-funk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/log"
)

const (
	ReasonFailedProcess     = "FailedProcess"
	ReasonWaitingProcess    = "WaitingProcess"
	ReasonSuccessfulProcess = "SuccessfulProcess"
	ReasonSkipProcess       = "SkipProcess"

	ConditionTypeDone = "EnsureDone"
)

// Provider defines a set of response interfaces for specific machine
// types in machine management.
type Provider interface {
	Name() string

	Validate(machine *platform.Machine) field.ErrorList

	PreCreate(machine *platform.Machine) error
	AfterCreate(machine *platform.Machine) error

	OnCreate(machine *platformv1.Machine, cluster *typesv1.Cluster) error
	OnUpdate(machine *platformv1.Machine, cluster *typesv1.Cluster) error
	OnDelete(machine *platformv1.Machine, cluster *typesv1.Cluster) error
}

var _ Provider = &DelegateProvider{}

type Handler func(*platformv1.Machine, *typesv1.Cluster) error

type DelegateProvider struct {
	ProviderName string

	ValidateFunc    func(machine *platform.Machine) field.ErrorList
	PreCreateFunc   func(machine *platform.Machine) error
	AfterCreateFunc func(machine *platform.Machine) error

	CreateHandlers []Handler
	DeleteHandlers []Handler
	UpdateHandlers []Handler
}

func (p *DelegateProvider) Name() string {
	if p.ProviderName == "" {
		return "unknown"
	}
	return p.ProviderName
}

func (p *DelegateProvider) Validate(machine *platform.Machine) field.ErrorList {
	if p.ValidateFunc != nil {
		return p.ValidateFunc(machine)
	}

	return nil
}

func (p *DelegateProvider) PreCreate(machine *platform.Machine) error {
	if p.PreCreateFunc != nil {
		return p.PreCreateFunc(machine)
	}

	return nil
}

func (p *DelegateProvider) AfterCreate(machine *platform.Machine) error {
	if p.AfterCreateFunc != nil {
		return p.AfterCreateFunc(machine)
	}

	return nil
}

func (p *DelegateProvider) OnCreate(machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	condition, err := p.getCreateCurrentCondition(machine)
	if err != nil {
		return err
	}

	now := metav1.Now()
	if cluster.Spec.Features.SkipConditions != nil &&
		funk.ContainsString(cluster.Spec.Features.SkipConditions, condition.Type) {
		machine.SetCondition(platformv1.MachineCondition{
			Type:               condition.Type,
			Status:             platformv1.ConditionTrue,
			LastProbeTime:      now,
			LastTransitionTime: now,
			Reason:             ReasonSkipProcess,
		})
	} else {
		f := p.getCreateHandler(condition.Type)
		if f == nil {
			return fmt.Errorf("can't get handler by %s", condition.Type)
		}
		log.Infow("OnCreate", "handler", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
			"machineName", machine.Name)
		err = f(machine, cluster)
		if err != nil {
			machine.SetCondition(platformv1.MachineCondition{
				Type:          condition.Type,
				Status:        platformv1.ConditionFalse,
				LastProbeTime: now,
				Message:       err.Error(),
				Reason:        ReasonFailedProcess,
			})
			machine.Status.Reason = ReasonFailedProcess
			machine.Status.Message = err.Error()
			return nil
		}

		machine.SetCondition(platformv1.MachineCondition{
			Type:               condition.Type,
			Status:             platformv1.ConditionTrue,
			LastProbeTime:      now,
			LastTransitionTime: now,
			Reason:             ReasonSuccessfulProcess,
		})
	}

	nextConditionType := p.getNextConditionType(condition.Type)
	if nextConditionType == ConditionTypeDone {
		machine.Status.Phase = platformv1.MachineRunning
	} else {
		machine.SetCondition(platformv1.MachineCondition{
			Type:               nextConditionType,
			Status:             platformv1.ConditionUnknown,
			LastProbeTime:      now,
			LastTransitionTime: now,
			Message:            "waiting process",
			Reason:             ReasonWaitingProcess,
		})
	}
	return nil
}

func (p *DelegateProvider) OnUpdate(machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	for _, f := range p.UpdateHandlers {
		log.Infow("OnUpdate", "handler", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
			"machineName", machine.Name)
		err := f(machine, cluster)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *DelegateProvider) OnDelete(machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	for _, f := range p.DeleteHandlers {
		log.Infow("OnDelete", "handler", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
			"machineName", machine.Name)
		err := f(machine, cluster)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h Handler) name() string {
	name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	i := strings.Index(name, "Ensure")
	if i == -1 {
		return ""
	}
	return strings.TrimSuffix(name[i:], "-fm")
}

func (p *DelegateProvider) getNextConditionType(conditionType string) string {
	var (
		i int
		f Handler
	)
	for i, f = range p.CreateHandlers {
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		if strings.Contains(name, conditionType) {
			break
		}
	}
	if i == len(p.CreateHandlers)-1 {
		return ConditionTypeDone
	}
	next := p.CreateHandlers[i+1]

	return next.name()
}

func (p *DelegateProvider) getCreateHandler(conditionType string) Handler {
	for _, f := range p.CreateHandlers {
		if conditionType == f.name() {
			return f
		}
	}

	return nil
}

func (p *DelegateProvider) getCreateCurrentCondition(c *platformv1.Machine) (*platformv1.MachineCondition, error) {
	if c.Status.Phase == platformv1.MachineRunning {
		return nil, errors.New("machine phase is running now")
	}
	if len(p.CreateHandlers) == 0 {
		return nil, errors.New("no create handlers")
	}

	if len(c.Status.Conditions) == 0 {
		return &platformv1.MachineCondition{
			Type:          p.CreateHandlers[0].name(),
			Status:        platformv1.ConditionUnknown,
			LastProbeTime: metav1.Now(),
			Message:       "waiting process",
			Reason:        ReasonWaitingProcess,
		}, nil
	}

	for _, condition := range c.Status.Conditions {
		if condition.Status == platformv1.ConditionFalse || condition.Status == platformv1.ConditionUnknown {
			return &condition, nil
		}
	}

	return nil, errors.New("no condition need process")
}
