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

package util

import (
	"io"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	"tkestack.io/tke/api/monitor"
	platformv1 "tkestack.io/tke/api/platform/v1"

	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	alert_config "github.com/prometheus/alertmanager/config"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// GenericRuleOperator defines the generic rule operator for prometheus
type GenericRuleOperator interface {
	RuleOperator
	RuleGroupOperator
	PersistentOperator
	LoadPromRule([]v1.RuleGroup) error
	SavePromRule() []v1.RuleGroup
}

// RuleOperator defines the CRUD interface of rule record operator
type RuleOperator interface {
	InsertRule(groupName string, rule *v1.Rule) (int, *v1.Rule, error)
	DeleteRule(groupName, recordName string) (*v1.Rule, error)
	UpdateRule(groupName, recordName string, rev int, rule *v1.Rule) (int, *v1.Rule, error)
	GetRule(groupName, recordName string) (int, *v1.Rule, error)
	ListRule(groupName string) ([]*v1.Rule, error)
}

// RuleGroupOperator defines the CRUD interface of rule group operator
type RuleGroupOperator interface {
	InsertRuleGroup(group *v1.RuleGroup) (int, *v1.RuleGroup, error)
	DeleteRuleGroup(groupName string) (*v1.RuleGroup, error)
	UpdateRuleGroup(groupName string, rev int, group *v1.RuleGroup) (int, *v1.RuleGroup, error)
	GetRuleGroup(groupName string) (int, *v1.RuleGroup, error)
	ListGroup() ([]*v1.RuleGroup, error)
}

// PersistentOperator defined the persistent method of rule operator
type PersistentOperator interface {
	Load(r io.Reader) error
	Save(w io.Writer) error
}

// GenericRouteOperator defines the generic alert route for alertmanager
type GenericRouteOperator interface {
	RouteOperator
	PersistentOperator
}

// RouteOperator defines the CRUD interface of alert route
type RouteOperator interface {
	InsertRoute(route *alert_config.Route) (*alert_config.Route, error)
	DeleteRoute(alertValue string) (*alert_config.Route, error)
	UpdateRoute(alertValue string, route *alert_config.Route) (*alert_config.Route, error)
	GetRoute(alertValue string) (*alert_config.Route, error)
	ListRoute() ([]*alert_config.Route, error)
}

type ClusterClientSets map[string]*kubernetes.Clientset
type MetricServerClientSets map[string]*metricsv.Clientset
type DynamicClientSet map[string]dynamic.Interface
type ClusterSet map[string]*platformv1.Cluster
type ClusterCredentialSet map[string]*platformv1.ClusterCredential
type ClusterStatisticSet map[string]*monitor.ClusterStatistic

type WorkloadCounter struct {
	Deployment  int
	DaemonSet   int
	StatefulSet int
	TApp        int
}

func (w *WorkloadCounter) Total() int {
	return w.Deployment + w.DaemonSet + w.StatefulSet + w.TApp
}

type ResourceCounter struct {
	NodeTotal          int
	NodeAbnormal       int
	HasMetricServer    bool
	CPUUsed            float64
	CPURequest         float64
	CPULimit           float64
	CPUCapacity        float64
	CPUAllocatable     float64
	CPURequestRate     float64
	CPUAllocatableRate float64
	CPUUsage           float64
	MemUsed            int64
	MemRequest         int64
	MemLimit           int64
	MemCapacity        int64
	MemAllocatable     int64
	MemRequestRate     float64
	MemAllocatableRate float64
	MemUsage           float64
	PodCount           int
	CPUCapacityMap     map[string]map[string]float64
	CPUAllocatableMap  map[string]map[string]float64
	MemCapacityMap     map[string]map[string]int64
	MemAllocatableMap  map[string]map[string]int64
}

type ComponentHealth struct {
	Scheduler         bool
	ControllerManager bool
	Etcd              bool
}
