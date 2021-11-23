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

package prometheus

import (
	"context"
	"time"

	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"tkestack.io/tke/pkg/monitor/util"
	"tkestack.io/tke/pkg/monitor/util/rule"
	"tkestack.io/tke/pkg/util/log"
)

const (
	prometheusService   = "prometheus"
	prometheusRuleAlert = "prometheus-alerts"
	prometheusCRDName   = "k8s"
)

func (h *processor) loadRule(ctx context.Context, clusterName string) (util.GenericRuleOperator, error) {
	monitoringClient, err := util.GetMonitoringClient(ctx, clusterName, h.platformClient)
	if err != nil {
		return nil, err
	}

	promRule, err := monitoringClient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Get(ctx, prometheusRuleAlert, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	log.Infof("Load rule from prometheusRule %s(%s)", clusterName, prometheusRuleAlert)
	ruleOp := rule.NewGenericRuleOperator(func(rule *v1.Rule) string {
		return rule.Alert
	})
	err = ruleOp.LoadPromRule(promRule.Spec.Groups)
	if err != nil {
		return nil, err
	}
	return ruleOp, nil
}

func (h *processor) saveRule(ctx context.Context, clusterName string, groups []v1.RuleGroup) error {
	monitoringClient, err := util.GetMonitoringClient(ctx, clusterName, h.platformClient)
	if err != nil {
		return err
	}

	log.Infof("Save rule to prometheusRule %s(%s)", clusterName, prometheusRuleAlert)

	return wait.PollImmediate(time.Second, time.Second*5, func() (done bool, err error) {
		promRule, getErr := monitoringClient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Get(ctx, prometheusRuleAlert, metav1.GetOptions{})
		if getErr != nil {
			return false, getErr
		}

		promRule.Spec.Groups = groups
		_, err = monitoringClient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Update(ctx, promRule, metav1.UpdateOptions{})
		if err == nil {
			return true, nil
		}

		if apierror.IsConflict(err) {
			return false, nil
		}

		return false, err
	})
}
