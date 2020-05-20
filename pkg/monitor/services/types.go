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

package services

import (
	"context"

	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/emicklei/go-restful"
	alertconfig "github.com/prometheus/alertmanager/config"
)

// BackendConfigProcessor defines the interface of operation rules service of prometheus and alertmanager
type BackendConfigProcessor interface {
	RegisterWebService(ws *restful.WebService)
	ConfigProcessor
}

// ConfigProcessor defines the interface of operation rules service
type ConfigProcessor interface {
	Create(req *restful.Request, resp *restful.Response)
	Update(req *restful.Request, resp *restful.Response)
	Delete(req *restful.Request, resp *restful.Response)
	Get(req *restful.Request, resp *restful.Response)
	List(req *restful.Request, resp *restful.Response)
}

// RuleProcessor defines the interface of operation rules service of prometheus
type RuleProcessor interface {
	CreateGroup(ctx context.Context, clusterName, groupName string, ruleGroup *v1.RuleGroup) error
	DeleteGroup(ctx context.Context, clusterName, groupName string) error
	GetGroup(ctx context.Context, clusterName, groupName string) (*v1.RuleGroup, error)
	UpdateGroup(ctx context.Context, clusterName, groupName string, ruleGroup *v1.RuleGroup) error
	ListGroups(ctx context.Context, clusterName string) ([]*v1.RuleGroup, error)
	CreateRule(ctx context.Context, clusterName, groupName string, recordName string, rule *v1.Rule) error
	DeleteRule(ctx context.Context, clusterName, groupName string, recordName string) error
	GetRule(ctx context.Context, clusterName, groupName string, recordName string) (*v1.Rule, error)
	UpdateRule(ctx context.Context, clusterName, groupName string, recordName string, rule *v1.Rule) error
	ListRules(ctx context.Context, clusterName string, groupName string) ([]*v1.Rule, error)
}

// RouteProcessor defines the interface of operation route service of alertmanager
type RouteProcessor interface {
	Create(ctx context.Context, clusterName string, alertValue string, route *alertconfig.Route) error
	Delete(ctx context.Context, clusterName string, alertValue string) error
	Get(ctx context.Context, clusterName string, alertValue string) (*alertconfig.Route, error)
	List(ctx context.Context, clusterName string) ([]*alertconfig.Route, error)
	Update(ctx context.Context, clusterName string, alertValue string, route *alertconfig.Route) error
}
