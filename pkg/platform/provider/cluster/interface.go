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
	"k8s.io/client-go/rest"
	platformv1client "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/types"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/platform/util/credential"
)

const (
	ReasonWaiting      = "Waiting"
	ReasonSkip         = "Skip"
	ReasonFailedInit   = "FailedInit"
	ReasonFailedUpdate = "FailedUpdate"
	ReasonFailedDelete = "FailedDelete"
	ConditionTypeDone  = "EnsureDone"
	defaultQPS         = -1
)

type APIProvider interface {
	RegisterHandler(mux *mux.PathRecorderMux)
	Validate(cluster *types.Cluster) field.ErrorList
	ValidateUpdate(cluster *types.Cluster, oldCluster *types.Cluster) field.ErrorList
	// for crd mode controller mutation
	MutateUpdate(cluster *types.Cluster, oldCluster *types.Cluster) ([]byte, field.ErrorList)
	PreCreate(cluster *types.Cluster) error
	AfterCreate(cluster *types.Cluster) error
}

type ControllerProvider interface {
	// Setup called by controller to give an chance for plugin do some init work.
	Setup() error
	// Teardown called by controller for plugin do some clean job.
	Teardown() error
	// NeedUpdate could be implemented by user to judge whether cluster need update or not.
	NeedUpdate(old, new *platformv1.Cluster) bool

	OnCreate(ctx context.Context, cluster *v1.Cluster) error
	OnUpdate(ctx context.Context, cluster *v1.Cluster) error
	OnDelete(ctx context.Context, cluster *v1.Cluster) error
	// OnFilter called by cluster controller informer for plugin
	// do the filter on the cluster obj for specific case:
	// return bool:
	//  false: drop the object to the queue
	//  true: add the object to queue, AddFunc and UpdateFunc will
	//  go through later
	OnFilter(ctx context.Context, cluster *platformv1.Cluster) bool
	// OnRunning call on first running.
	OnRunning(ctx context.Context, cluster *v1.Cluster) error
}

type RestConfigProvider interface {
	GetRestConfig(ctx context.Context, cluster *platformv1.Cluster, username string) (*rest.Config, error)
}

// Provider defines a set of response interfaces for specific cluster
// types in cluster management.
type Provider interface {
	Name() string

	APIProvider
	ControllerProvider
	RestConfigProvider
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

	ValidateFunc       func(cluster *types.Cluster) field.ErrorList
	ValidateUpdateFunc func(cluster *types.Cluster, oldCluster *types.Cluster) field.ErrorList
	MutateUpdateFunc   func(cluster *types.Cluster, oldCluster *types.Cluster) ([]byte, field.ErrorList)
	PreCreateFunc      func(cluster *types.Cluster) error
	AfterCreateFunc    func(cluster *types.Cluster) error

	CreateHandlers    []Handler
	DeleteHandlers    []Handler
	UpdateHandlers    []Handler
	UpgradeHandlers   []Handler
	ScaleUpHandlers   []Handler
	ScaleDownHandlers []Handler
	PlatformClient    platformv1client.PlatformV1Interface
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

func (p *DelegateProvider) ValidateUpdate(cluster *types.Cluster, oldCluster *types.Cluster) field.ErrorList {
	if p.ValidateUpdateFunc != nil {
		return p.ValidateUpdateFunc(cluster, oldCluster)
	}

	return nil
}

func (p *DelegateProvider) MutateUpdate(cluster *types.Cluster, oldCluster *types.Cluster) ([]byte, field.ErrorList) {
	if p.MutateUpdateFunc != nil {
		return p.MutateUpdateFunc(cluster, oldCluster)
	}

	return nil, nil
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

func (p *DelegateProvider) getUpdateReason(c *v1.Cluster) (reason string) {
	if c.Status.Phase == platformv1.ClusterUpgrading {
		return fmt.Sprintf("%s to kubernetes %s", platformv1.ClusterUpgrading, c.Spec.Version)
	}
	if c.Status.Phase == platformv1.ClusterUpscaling {
		var ips []string
		for _, machine := range c.Spec.ScalingMachines {
			ips = append(ips, machine.IP)
		}
		return fmt.Sprintf("%s on machine %s", platformv1.ClusterUpscaling, strings.Join(ips, ","))
	}
	return ""
}

func (p *DelegateProvider) OnCreate(ctx context.Context, cluster *v1.Cluster) error {
	condition, err := p.getCurrentCondition(cluster, platformv1.ClusterInitializing, p.CreateHandlers)
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
		}, false)
	} else {
		handler := p.getHandler(condition.Type, p.CreateHandlers)
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
			}, false)
			return nil
		}

		cluster.SetCondition(platformv1.ClusterCondition{
			Type:   condition.Type,
			Status: platformv1.ConditionTrue,
		}, false)
	}

	nextConditionType := p.getNextConditionType(condition.Type, p.CreateHandlers)
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
		}, false)
	}

	return nil
}

func (p *DelegateProvider) OnUpdate(ctx context.Context, cluster *v1.Cluster) error {
	handlers := []Handler{}
	phase := cluster.Status.Phase
	if phase == platformv1.ClusterRunning || phase == platformv1.ClusterFailed {
		handlers = p.UpdateHandlers
		return p.houseKeeping(ctx, cluster, handlers)
	}
	if phase == platformv1.ClusterUpgrading {
		handlers = p.UpgradeHandlers
	}
	if phase == platformv1.ClusterUpscaling {
		handlers = p.CreateHandlers
	}
	if phase == platformv1.ClusterDownscaling {
		handlers = p.ScaleDownHandlers
	}
	condition, err := p.getCurrentCondition(cluster, phase, handlers)
	if err != nil {
		return err
	}
	if condition == nil {
		return nil
	}
	if cluster.Spec.Features.SkipConditions != nil &&
		funk.ContainsString(cluster.Spec.Features.SkipConditions, condition.Type) {
		cluster.SetCondition(platformv1.ClusterCondition{
			Type:    condition.Type,
			Status:  platformv1.ConditionTrue,
			Reason:  ReasonSkip,
			Message: "Skip current condition",
		}, true)
	} else {
		handler := p.getHandler(condition.Type, handlers)
		if handler == nil {
			return fmt.Errorf("can't get handler by %s", condition.Type)
		}
		ctx := log.FromContext(ctx).WithName("ClusterProvider.OnUpdate").WithName(handler.Name()).WithContext(ctx)
		log.FromContext(ctx).Info("Doing")
		startTime := time.Now()
		err = handler(ctx, cluster)
		log.FromContext(ctx).Info("Done", "error", err, "cost", time.Since(startTime).String())
		if err != nil {
			cluster.SetCondition(platformv1.ClusterCondition{
				Type:    condition.Type,
				Status:  platformv1.ConditionFalse,
				Message: err.Error(),
				Reason:  ReasonFailedUpdate,
			}, true)
			return nil
		}
		cluster.SetCondition(platformv1.ClusterCondition{
			Type:   condition.Type,
			Status: platformv1.ConditionTrue,
			Reason: p.getUpdateReason(cluster),
		}, true)
	}

	nextConditionType := p.getNextConditionType(condition.Type, handlers)
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
		}, true)
	}

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

func (p *DelegateProvider) OnFilter(ctx context.Context, cluster *platformv1.Cluster) (pass bool) {
	return true
}

func (p *DelegateProvider) NeedUpdate(old, new *platformv1.Cluster) bool {
	return false
}

func (p *DelegateProvider) getNextConditionType(conditionType string, handlers []Handler) string {
	var (
		i       int
		handler Handler
	)
	for i, handler = range handlers {
		if handler.Name() == conditionType {
			break
		}
	}
	if i == len(handlers)-1 {
		return ConditionTypeDone
	}
	next := handlers[i+1]

	return next.Name()
}

func (p *DelegateProvider) getHandler(conditionType string, handlers []Handler) Handler {
	for _, handler := range handlers {
		if conditionType == handler.Name() {
			return handler
		}
	}

	return nil
}

func (p *DelegateProvider) houseKeeping(ctx context.Context, cluster *v1.Cluster, handlers []Handler) error {
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

func (p *DelegateProvider) getCurrentCondition(c *v1.Cluster, phase platformv1.ClusterPhase, handlers []Handler) (*platformv1.ClusterCondition, error) {
	if c.Status.Phase != phase {
		return nil, fmt.Errorf("cluster phase is %s now", phase)
	}
	if len(handlers) == 0 {
		return nil, fmt.Errorf("no handlers")
	}

	if len(c.Status.Conditions) == 0 {
		return &platformv1.ClusterCondition{
			Type:    handlers[0].Name(),
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
	if c.Status.Phase == platformv1.ClusterUpgrading ||
		c.Status.Phase == platformv1.ClusterUpscaling ||
		c.Status.Phase == platformv1.ClusterDownscaling ||
		c.Status.Phase == platformv1.ClusterRunning {
		return &platformv1.ClusterCondition{
			Type:    handlers[0].Name(),
			Status:  platformv1.ConditionUnknown,
			Message: "waiting process",
			Reason:  ReasonWaiting,
		}, nil
	}
	return nil, errors.New("no condition need process")
}

// GetRestConfig returns the cluster's rest config
func (p *DelegateProvider) GetRestConfig(ctx context.Context, cluster *platformv1.Cluster, username string) (*rest.Config, error) {
	cc, err := credential.GetClusterCredentialV1(ctx, p.PlatformClient, cluster, username)
	if err != nil {
		return nil, err
	}
	config := &rest.Config{}
	if cc != nil {
		config = cc.RESTConfig(cluster)
	}
	if config.QPS == 0 {
		config.QPS = defaultQPS
	}
	return config, nil
}
