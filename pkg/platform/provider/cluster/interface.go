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
	"k8s.io/apiserver/pkg/server/mux"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/types"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
)

const (
	ReasonWaiting      = "Waiting"
	ReasonSkip         = "Skip"
	ReasonFailedInit   = "FailedInit"
	ReasonFailedUpdate = "FailedUpdate"
	ReasonFailedDelete = "FailedDelete"

	ConditionTypeDone = "EnsureDone"
)

type APIProvider interface {
	RegisterHandler(mux *mux.PathRecorderMux)
	Validate(cluster *types.Cluster) field.ErrorList
	PreCreate(cluster *types.Cluster) error
	AfterCreate(cluster *types.Cluster) error
}

type ControllerProvider interface {
	// Setup called by controller to give an chance for plugin do some init work.
	Setup() error
	// Teardown called by controller for plugin do some clean job.
	Teardown() error

	OnCreate(ctx context.Context, cluster *v1.Cluster) error
	OnUpdate(ctx context.Context, cluster *v1.Cluster) error
	OnDelete(ctx context.Context, cluster *v1.Cluster) error

	// OnRunning call on first running.
	OnRunning(ctx context.Context, cluster *v1.Cluster) error
}

// Provider defines a set of response interfaces for specific cluster
// types in cluster management.
type Provider interface {
	Name() string

	APIProvider
	ControllerProvider
}

var _ Provider = &DelegateProvider{}

type Handler func(context.Context, *v1.Cluster) error

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

	ValidateFunc    func(cluster *types.Cluster) field.ErrorList
	PreCreateFunc   func(cluster *types.Cluster) error
	AfterCreateFunc func(cluster *types.Cluster) error

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

func (p *DelegateProvider) Setup() error {
	return nil
}

func (p *DelegateProvider) Teardown() error {
	return nil
}

func (p *DelegateProvider) RegisterHandler(mux *mux.PathRecorderMux) {
}

func (p *DelegateProvider) Validate(cluster *types.Cluster) field.ErrorList {
	if p.ValidateFunc != nil {
		return p.ValidateFunc(cluster)
	}

	return nil
}

func (p *DelegateProvider) PreCreate(cluster *types.Cluster) error {
	if p.PreCreateFunc != nil {
		return p.PreCreateFunc(cluster)
	}

	return nil
}

func (p *DelegateProvider) AfterCreate(cluster *types.Cluster) error {
	if p.AfterCreateFunc != nil {
		return p.AfterCreateFunc(cluster)
	}

	return nil
}

func (p *DelegateProvider) OnCreate(ctx context.Context, cluster *v1.Cluster) error {
	condition, err := p.getCreateCurrentCondition(cluster)
	if err != nil {
		return err
	}

	if cluster.Spec.Features.SkipConditions != nil &&
		funk.ContainsString(cluster.Spec.Features.SkipConditions, condition.Type) {
		cluster.SetCondition(platformv1.ClusterCondition{
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
		ctx = log.FromContext(ctx).WithName("ClusterProvider.OnCreate").WithName(handler.Name()).WithContext(ctx)
		log.FromContext(ctx).Info("Doing")
		startTime := time.Now()
		err = handler(ctx, cluster)
		log.FromContext(ctx).Info("Done", "error", err, "cost", time.Since(startTime).String())
		if err != nil {
			cluster.SetCondition(platformv1.ClusterCondition{
				Type:    condition.Type,
				Status:  platformv1.ConditionFalse,
				Message: err.Error(),
				Reason:  ReasonFailedInit,
			})
			return nil
		}

		cluster.SetCondition(platformv1.ClusterCondition{
			Type:   condition.Type,
			Status: platformv1.ConditionTrue,
		})
	}

	nextConditionType := p.getNextConditionType(condition.Type)
	if nextConditionType == ConditionTypeDone {
		cluster.Status.Phase = platformv1.ClusterRunning
		if err := p.OnRunning(ctx, cluster); err != nil {
			return fmt.Errorf("%s.OnRunning error: %w", p.Name(), err)
		}
	} else {
		cluster.SetCondition(platformv1.ClusterCondition{
			Type:    nextConditionType,
			Status:  platformv1.ConditionUnknown,
			Message: "waiting execute",
			Reason:  ReasonWaiting,
		})
	}

	return nil
}

func (p *DelegateProvider) OnUpdate(ctx context.Context, cluster *v1.Cluster) error {
	for _, handler := range p.UpdateHandlers {
		ctx := log.FromContext(ctx).WithName("ClusterProvider.OnUpdate").WithName(handler.Name()).WithContext(ctx)
		log.FromContext(ctx).Info("Doing")
		startTime := time.Now()
		err := handler(ctx, cluster)
		log.FromContext(ctx).Info("Done", "error", err, "cost", time.Since(startTime).String())
		if err != nil {
			cluster.Status.Reason = ReasonFailedUpdate
			cluster.Status.Message = fmt.Sprintf("%s error: %v", handler.Name(), err)
			return err
		}
	}
	cluster.Status.Reason = ""
	cluster.Status.Message = ""

	return nil
}

func (p *DelegateProvider) OnDelete(ctx context.Context, cluster *v1.Cluster) error {
	for _, handler := range p.DeleteHandlers {
		ctx := log.FromContext(ctx).WithName("ClusterProvider.OnDelete").WithName(handler.Name()).WithContext(ctx)
		log.FromContext(ctx).Info("Doing")
		startTime := time.Now()
		err := handler(ctx, cluster)
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

func (p *DelegateProvider) OnRunning(ctx context.Context, cluster *v1.Cluster) error {
	return nil
}

func (p *DelegateProvider) getNextConditionType(conditionType string) string {
	var (
		i       int
		handler Handler
	)
	for i, handler = range p.CreateHandlers {
		if handler.Name() == conditionType {
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
	for _, handler := range p.CreateHandlers {
		if conditionType == handler.Name() {
			return handler
		}
	}

	return nil
}

func (p *DelegateProvider) getCreateCurrentCondition(c *v1.Cluster) (*platformv1.ClusterCondition, error) {
	if c.Status.Phase == platformv1.ClusterRunning {
		return nil, errors.New("cluster phase is running now")
	}
	if len(p.CreateHandlers) == 0 {
		return nil, errors.New("no create handlers")
	}

	if len(c.Status.Conditions) == 0 {
		return &platformv1.ClusterCondition{
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
