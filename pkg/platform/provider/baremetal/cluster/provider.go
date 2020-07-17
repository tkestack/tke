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
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/config"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/validation"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/spec"
	"tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/log"
)

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
	clusterprovider.Register(p.Name(), p)
}

type Provider struct {
	config         *config.Config
	createHandlers []Handler
}

var _ clusterprovider.Provider = &Provider{}

type Handler func(*Cluster) error

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

		p.EnsureRegistryHosts,
		p.EnsureKernelModule,
		p.EnsureSysctl,
		p.EnsureDisableSwap,

		p.EnsurePreflight, // wait basic setting done

		p.EnsureClusterComplete,

		p.EnsureNvidiaDriver,
		p.EnsureNvidiaContainerRuntime,
		p.EnsureDocker,
		p.EnsureKubelet,
		p.EnsureCNIPlugins,
		p.EnsureKubeadm,

		p.EnsurePrepareForControlplane,

		p.EnsureKubeadmInitKubeletStartPhase,
		p.EnsureKubeadmInitCertsPhase,
		p.EnsureStoreCredential,
		p.EnsureKubeconfig,
		p.EnsureKubeadmInitKubeConfigPhase,
		p.EnsureKubeadmInitControlPlanePhase,
		p.EnsureKubeadmInitEtcdPhase,
		p.EnsureKubeadmInitWaitControlPlanePhase,
		p.EnsureKubeadmInitUploadConfigPhase,
		p.EnsureKubeadmInitUploadCertsPhase,
		p.EnsureKubeadmInitBootstrapTokenPhase,
		p.EnsureKubeadmInitAddonPhase,
		p.EnsureGalaxy,

		p.EnsureJoinControlePlane,
		p.EnsurePatchAnnotation, // wait rest master ready
		p.EnsureMarkControlPlane,

		p.EnsureNvidiaDevicePlugin,

		p.EnsureCleanup,

		p.EnsurePostInstallHook,
	}

	return p, nil
}

func (p *Provider) Name() string {
	return "Baremetal"
}

func (p *Provider) ValidateCredential(cluster clusterprovider.InternalCluster) field.ErrorList {
	return nil
}

func (p *Provider) Validate(c *platform.Cluster) field.ErrorList {
	return validation.ValidateCluster(c)
}

func (p *Provider) PreCreate(user clusterprovider.UserInfo, cluster platform.Cluster) (platform.Cluster, error) {
	if cluster.Spec.Version == "" {
		cluster.Spec.Version = spec.K8sVersions[0]
	}
	if cluster.Spec.ClusterCIDR == "" {
		cluster.Spec.ClusterCIDR = "10.244.0.0/16"
	}
	if cluster.Spec.NetworkDevice == "" {
		cluster.Spec.NetworkDevice = "eth0"
	}
	if cluster.Spec.Features.IPVS == nil {
		cluster.Spec.Features.IPVS = pointer.ToBool(true)
	}
	if cluster.Spec.Properties.MaxClusterServiceNum == nil && cluster.Spec.ServiceCIDR == nil {
		cluster.Spec.Properties.MaxClusterServiceNum = pointer.ToInt32(256)
	}
	if cluster.Spec.Properties.MaxNodePodNum == nil {
		cluster.Spec.Properties.MaxNodePodNum = pointer.ToInt32(256)
	}

	return cluster, nil
}

func (p *Provider) AfterCreate(cluster platform.Cluster) ([]interface{}, error) {
	return nil, nil
}

func (p *Provider) OnDelete(cluster platformv1.Cluster) error {
	return nil
}

func (p *Provider) OnInitialize(args clusterprovider.Cluster) (clusterprovider.Cluster, error) {
	c, err := NewCluster(args, p.config)
	if err != nil {
		log.Warn("NewCluster error", log.Err(err))
		return c.Cluster, err
	}

	err = p.create(c)

	return c.Cluster, err
}

func (p *Provider) OnUpdate(args clusterprovider.Cluster) (clusterprovider.Cluster, error) {
	return args, nil
}

func (p *Provider) create(c *Cluster) error {
	condition, err := p.getCreateCurrentCondition(c)
	if err != nil {
		return err
	}

	skipConditions := c.Spec.Features.SkipConditions
	if p.config.Feature.SkipConditions != nil {
		skipConditions = append(skipConditions, p.config.Feature.SkipConditions...)
	}

	now := metav1.Now()
	if funk.ContainsString(skipConditions, condition.Type) {
		c.SetCondition(platformv1.ClusterCondition{
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
		err = f(c)
		if err != nil {
			log.Warn(err.Error())
			c.SetCondition(platformv1.ClusterCondition{
				Type:          condition.Type,
				Status:        platformv1.ConditionFalse,
				LastProbeTime: now,
				Message:       err.Error(),
				Reason:        ReasonFailedProcess,
			})
			c.Status.Reason = ReasonFailedProcess
			c.Status.Message = err.Error()
			return nil
		}

		c.SetCondition(platformv1.ClusterCondition{
			Type:               condition.Type,
			Status:             platformv1.ConditionTrue,
			LastProbeTime:      now,
			LastTransitionTime: now,
			Reason:             ReasonSuccessfulProcess,
		})
	}

	nextConditionType := p.getNextConditionType(condition.Type)
	if nextConditionType == ConditionTypeDone {
		c.Status.Phase = platformv1.ClusterRunning

		log.Info("all done")
	} else {
		c.SetCondition(platformv1.ClusterCondition{
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

func (p *Provider) getCreateCurrentCondition(c *Cluster) (*platformv1.ClusterCondition, error) {
	if c.Status.Phase == platformv1.ClusterRunning {
		return nil, errors.New("cluster phase is running now")
	}
	if len(p.createHandlers) == 0 {
		return nil, errors.New("no create handlers")
	}

	if len(c.Status.Conditions) == 0 {
		return &platformv1.ClusterCondition{
			Type:          p.createHandlers[0].name(),
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
