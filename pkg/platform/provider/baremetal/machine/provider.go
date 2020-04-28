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
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/config"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/validation"
	machineprovider "tkestack.io/tke/pkg/platform/provider/machine"
	"tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/log"
)

const providerName = "Baremetal"

const (
	ReasonFailedProcess     = "FailedProcess"
	ReasonWaitingProcess    = "WaitingProcess"
	ReasonSuccessfulProcess = "SuccessfulProcess"
	ReasonSkipProcess       = "SkipProcess"

	ConditionTypeDone = "EnsureDone"
)

func init() {
	p, err := NewProvider()
	if err != nil {
		panic(err)
	}
	machineprovider.Register(p.Name(), p)
}

type Provider struct {
	config         *config.Config
	createHandlers []Handler
}

func NewProvider() (*Provider, error) {
	p := new(Provider)

	cfg, err := config.New(constants.ConfigFile)
	if err != nil {
		return nil, err
	}
	p.config = cfg

	containerregistry.Init(cfg.Registry.Domain, cfg.Registry.Namespace)

	p.createHandlers = []Handler{
		p.EnsureCopyFiles,
		p.EnsurePreInstallHook,

		p.EnsureClean,
		p.EnsureRegistryHosts,
		p.EnsureKernelModule,
		p.EnsureSysctl,
		p.EnsureDisableSwap,

		p.EnsurePreflight, // wait basic setting done

		p.EnsureNvidiaDriver,
		p.EnsureNvidiaContainerRuntime,
		p.EnsureDocker,
		p.EnsureKubelet,
		p.EnsureCNIPlugins,
		p.EnsureKubeadm,

		p.EnsureJoinNode,
		p.EnsureKubeconfig,
		p.EnsureMarkNode,
		p.EnsureNodeReady,

		p.EnsurePostInstallHook,
	}

	return p, nil
}

type Handler func(*Machine) error

var _ machineprovider.Provider = &Provider{}

func (p *Provider) Name() string {
	return providerName
}

func (p *Provider) Validate(machine *platform.Machine) field.ErrorList {
	return validation.ValidateMachine(machine)
}

func (p *Provider) OnInitialize(tkev1Machine platformv1.Machine, tkev1Cluster platformv1.Cluster, credential platformv1.ClusterCredential) (platformv1.Machine, error) {
	m, err := NewMachine(tkev1Machine, &tkev1Cluster, &credential, p.config)
	if err != nil {
		log.Warn("NewMachine error", log.Err(err))
		return m.Machine, err
	}

	err = p.create(m)

	return m.Machine, err
}

func (p *Provider) create(m *Machine) error {
	condition, err := p.getCreateCurrentCondition(m)
	if err != nil {
		return err
	}

	skipConditions := m.Cluster.Spec.Features.SkipConditions
	if skipConditions == nil {
		skipConditions = p.config.Feature.SkipConditions
	}
	now := metav1.Now()
	if funk.ContainsString(skipConditions, condition.Type) {
		m.SetCondition(platformv1.MachineCondition{
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
		err = f(m)
		now := metav1.Now()
		if err != nil {
			log.Warn(err.Error())
			m.SetCondition(platformv1.MachineCondition{
				Type:          condition.Type,
				Status:        platformv1.ConditionFalse,
				LastProbeTime: now,
				Message:       err.Error(),
				Reason:        ReasonFailedProcess,
			})
			m.Status.Reason = ReasonFailedProcess
			m.Status.Message = err.Error()
			return nil
		}

		m.SetCondition(platformv1.MachineCondition{
			Type:               condition.Type,
			Status:             platformv1.ConditionTrue,
			LastProbeTime:      now,
			LastTransitionTime: now,
			Reason:             ReasonSuccessfulProcess,
		})
	}

	nextConditionType := p.getNextConditionType(condition.Type)
	if nextConditionType == ConditionTypeDone {
		m.Status.Phase = platformv1.MachineRunning

		log.Info("all done")
	} else {
		m.SetCondition(platformv1.MachineCondition{
			Type:               nextConditionType,
			Status:             platformv1.ConditionUnknown,
			LastProbeTime:      now,
			LastTransitionTime: now,
			Message:            "waiting process",
			Reason:             ReasonWaitingProcess,
		})

		log.Infof("%s is done, next is %s", condition.Type, nextConditionType)
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

func (p *Provider) getNextConditionType(conditionType string) string {
	var (
		i int
		f Handler
	)
	for i, f = range p.createHandlers {
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		if strings.Contains(name, conditionType) {
			break
		}
	}
	if i == len(p.createHandlers)-1 {
		return ConditionTypeDone
	}
	next := p.createHandlers[i+1]

	return next.name()
}

func (p *Provider) getCreateHandler(conditionType string) Handler {
	for _, f := range p.createHandlers {
		if conditionType == f.name() {
			return f
		}
	}

	return nil
}

func (p *Provider) getCreateCurrentCondition(m *Machine) (*platformv1.MachineCondition, error) {
	if m.Status.Phase == platformv1.MachineRunning {
		return nil, errors.New("machine phase is running now")
	}
	if len(p.createHandlers) == 0 {
		return nil, errors.New("no create handlers")
	}

	if len(m.Status.Conditions) == 0 {
		return &platformv1.MachineCondition{
			Type:          p.createHandlers[0].name(),
			Status:        platformv1.ConditionUnknown,
			LastProbeTime: metav1.Now(),
			Message:       "waiting process",
			Reason:        ReasonWaitingProcess,
		}, nil
	}

	for _, condition := range m.Status.Conditions {
		if condition.Status == platformv1.ConditionFalse || condition.Status == platformv1.ConditionUnknown {
			return &condition, nil
		}
	}

	return nil, errors.New("no condition need process")
}
