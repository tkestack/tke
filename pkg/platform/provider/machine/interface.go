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
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"tkestack.io/tke/pkg/util/log"

	"github.com/thoas/go-funk"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
)

const (
	ReasonWaiting      = "Waiting"
	ReasonSkip         = "Skip"
	ReasonFailedInit   = "FailedInit"
	ReasonFailedUpdate = "FailedUpdate"
	ReasonFailedDelete = "FailedDelete"

	ConditionTypeDone = "EnsureDone"
)

// APIProvider APIProvider
type APIProvider interface {
	Validate(machine *platform.Machine) field.ErrorList
	ValidateUpdate(machine *platform.Machine, oldMachine *platform.Machine) field.ErrorList
	PreCreate(machine *platform.Machine) error
	AfterCreate(machine *platform.Machine) error
}

// ControllerProvider ControllerProvider
type ControllerProvider interface {
	OnCreate(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error
	OnUpdate(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error
	OnDelete(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error
}

// Provider defines a set of response interfaces for specific machine
// types in machine management.
type Provider interface {
	Name() string

	APIProvider
	ControllerProvider
}

var _ Provider = &DelegateProvider{}

type Handler func(context.Context, *platformv1.Machine, *typesv1.Cluster) error

func (h Handler) Name() string {
	name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	i := strings.LastIndex(name, ".")
	if i == -1 {
		return "Unknown"
	}
	return strings.TrimSuffix(name[i+1:], "-fm")
}

type DelegateProvider struct {
	ProviderName string

	ValidateFunc       func(machine *platform.Machine) field.ErrorList
	ValidateUpdateFunc func(machine *platform.Machine, oldMachine *platform.Machine) field.ErrorList
	PreCreateFunc      func(machine *platform.Machine) error
	AfterCreateFunc    func(machine *platform.Machine) error

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

func (p *DelegateProvider) ValidateUpdate(machine *platform.Machine, oldMachine *platform.Machine) field.ErrorList {
	if p.ValidateUpdateFunc != nil {
		return p.ValidateUpdateFunc(machine, oldMachine)
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

func (p *DelegateProvider) OnCreate(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	condition, err := p.getCreateCurrentCondition(machine)
	if err != nil {
		return err
	}

	if cluster.Spec.Features.SkipConditions != nil &&
		funk.ContainsString(cluster.Spec.Features.SkipConditions, condition.Type) {
		machine.SetCondition(platformv1.MachineCondition{
			Type:    condition.Type,
			Status:  platformv1.ConditionTrue,
			Reason:  ReasonSkip,
			Message: "Skip current condition",
		})
	} else {
		handler := p.getCreateHandler(condition.Type)
		if handler == nil {
			return fmt.Errorf("can't get handler by %s", condition.Type)
		}
		ctx := log.FromContext(ctx).WithName("MachineProvider.OnCreate").WithName(handler.Name()).WithContext(ctx)
		log.FromContext(ctx).Info("Doing")
		startTime := time.Now()
		err = handler(ctx, machine, cluster)
		log.FromContext(ctx).Info("Done", "error", err, "cost", time.Since(startTime).String())
		if err != nil {
			machine.SetCondition(platformv1.MachineCondition{
				Type:    condition.Type,
				Status:  platformv1.ConditionFalse,
				Message: err.Error(),
				Reason:  ReasonFailedInit,
			})
			return err
		}

		machine.SetCondition(platformv1.MachineCondition{
			Type:   condition.Type,
			Status: platformv1.ConditionTrue,
		})
	}

	nextConditionType := p.getNextConditionType(condition.Type)
	if nextConditionType == ConditionTypeDone {
		machine.Status.Phase = platformv1.MachineRunning
	} else {
		machine.SetCondition(platformv1.MachineCondition{
			Type:    nextConditionType,
			Status:  platformv1.ConditionUnknown,
			Message: "waiting execute",
			Reason:  ReasonWaiting,
		})
	}

	return nil
}

func (p *DelegateProvider) OnUpdate(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	if machine.Status.Phase != platformv1.MachineUpgrading {
		return nil
	}
	for _, handler := range p.UpdateHandlers {
		ctx := log.FromContext(ctx).WithName("MachineProvider.OnUpdate").WithName(handler.Name()).WithContext(ctx)
		log.FromContext(ctx).Info("Doing")
		startTime := time.Now()
		err := handler(ctx, machine, cluster)
		log.FromContext(ctx).Info("Done", "error", err, "cost", time.Since(startTime).String())
		if err != nil {
			machine.Status.Reason = ReasonFailedUpdate
			machine.Status.Message = fmt.Sprintf("%s error: %v", handler.Name(), err)
			return err
		}
	}
	machine.Status.Reason = ""
	machine.Status.Message = ""

	return nil
}

func (p *DelegateProvider) OnDelete(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	for _, handler := range p.DeleteHandlers {
		ctx := log.FromContext(ctx).WithName("MachineProvider.OnDelete").WithName(handler.Name()).WithContext(ctx)
		log.FromContext(ctx).Info("Doing")
		startTime := time.Now()
		err := handler(ctx, machine, cluster)
		log.FromContext(ctx).Info("Done", "error", err, "cost", time.Since(startTime).String())
		if err != nil {
			cluster.Status.Reason = ReasonFailedDelete
			cluster.Status.Message = fmt.Sprintf("%s error: %v", handler.Name(), err)
			return err
		}
	}
	cluster.Status.Reason = ""
	cluster.Status.Message = ""

	return nil
}

func (p *DelegateProvider) getNextConditionType(conditionType string) string {
	var (
		i       int
		handler Handler
	)
	for i, handler = range p.CreateHandlers {
		name := handler.Name()
		if name == conditionType {
			break
		}
	}
	if i == len(p.CreateHandlers)-1 {
		return ConditionTypeDone
	}
	next := p.CreateHandlers[i+1]

	return next.Name()
}

func (p *DelegateProvider) getCreateHandler(conditionType string) Handler {
	for _, f := range p.CreateHandlers {
		if conditionType == f.Name() {
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
			Type:    p.CreateHandlers[0].Name(),
			Status:  platformv1.ConditionUnknown,
			Message: "waiting process",
			Reason:  ReasonWaiting,
		}, nil
	}

	for _, condition := range c.Status.Conditions {
		if condition.Status == platformv1.ConditionFalse || condition.Status == platformv1.ConditionUnknown {
			return &condition, nil
		}
	}

	return nil, errors.New("no condition need process")
}
