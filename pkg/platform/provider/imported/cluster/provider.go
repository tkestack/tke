/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package cluster

import (
	"fmt"
	"net"
	"reflect"
	"runtime"
	"strings"
	"time"

	"tkestack.io/tke/pkg/platform/util"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/util/log"
)

const providerName = "Imported"

const (
	ReasonFailedProcess     = "FailedProcess"
	ReasonWaitingProcess    = "WaitingProcess"
	ReasonSuccessfulProcess = "SuccessfulProcess"

	ConditionTypeDone = "EnsureDone"
)

func init() {
	p, err := newProvider()
	if err != nil {
		panic(err)
	}
	clusterprovider.Register(p.Name(), p)
}

type Provider struct {
	createHandlers []Handler
}

var _ clusterprovider.Provider = &Provider{}

type Handler func(*Cluster) error

func newProvider() (*Provider, error) {
	p := new(Provider)

	p.createHandlers = []Handler{
		p.EnsureClusterReady,
	}

	return p, nil
}

func (p *Provider) Name() string {
	return providerName
}

func (p *Provider) ValidateCredential(cluster clusterprovider.InternalCluster) (field.ErrorList, error) {
	var allErrs field.ErrorList

	credential := cluster.ClusterCredential

	if len(credential.CACert) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("caCert"), "must specify CA root certificate"))
	}

	if credential.Token == nil && credential.ClientKey == nil && credential.ClientCert == nil {
		allErrs = append(allErrs, field.Required(field.NewPath(""), "must specify at least one of token or client certificate authentication"))
	} else {
		if credential.ClientCert == nil && credential.ClientKey != nil {
			allErrs = append(allErrs, field.Required(field.NewPath("clientCert"), "must specify both the public and private keys of the client certificate"))
		}

		if credential.ClientCert != nil && credential.ClientKey == nil {
			allErrs = append(allErrs, field.Required(field.NewPath("clientKey"), "must specify both the public and private keys of the client certificate"))
		}

		clientset, err := util.BuildClientSet(&cluster.Cluster, &cluster.ClusterCredential)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(field.NewPath(""), err))
		}
		_, err = clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath(""), credential.ClusterName, fmt.Sprintf("invalid credential:%s", err)))
		}
	}

	return allErrs, nil
}

func (p *Provider) Validate(c platform.Cluster) (field.ErrorList, error) {
	var allErrs field.ErrorList

	if len(c.Status.Addresses) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("status", "addresses"), "must specify at least one obj access address"))
	} else {
		for _, address := range c.Status.Addresses {
			if address.Host == "" {
				allErrs = append(allErrs, field.Required(field.NewPath("status", "addresses", string(address.Type), "host"), "must specify the ip of address"))
			}
			if address.Port == 0 {
				allErrs = append(allErrs, field.Required(field.NewPath("status", "addresses", string(address.Type), "port"), "must specify the port of address"))
			}
			_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", address.Host, address.Port), 5*time.Second)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("status", "addresses"), address, err.Error()))
			}
		}
	}
	return allErrs, nil
}

func (p *Provider) PreCreate(user clusterprovider.UserInfo, cluster platform.Cluster) (platform.Cluster, error) {
	return cluster, nil
}

func (p *Provider) AfterCreate(cluster platform.Cluster) ([]interface{}, error) {
	return nil, nil
}

func (p *Provider) ValidateUpdate(cluster platform.Cluster, oldCluster platform.Cluster) (field.ErrorList, error) {
	var allErrs field.ErrorList
	return allErrs, nil
}

func (p *Provider) OnDelete(cluster platformv1.Cluster) error {
	return nil
}

func (p *Provider) OnInitialize(args clusterprovider.Cluster) (clusterprovider.Cluster, error) {
	c, err := NewCluster(args)
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

	f := p.getCreateHandler(condition.Type)
	if f == nil {
		return fmt.Errorf("can't get handler by %s", condition.Type)
	}
	err = f(c)
	now := metav1.Now()
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
