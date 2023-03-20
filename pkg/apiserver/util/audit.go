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

package util

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/audit"
	"k8s.io/apiserver/pkg/audit/policy"
	genericapiserver "k8s.io/apiserver/pkg/server"
	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	genericapiserveroptions "k8s.io/apiserver/pkg/server/options"
	pluginbuffered "k8s.io/apiserver/plugin/pkg/audit/buffered"
	pluginlog "k8s.io/apiserver/plugin/pkg/audit/log"
	plugintruncate "k8s.io/apiserver/plugin/pkg/audit/truncate"
	pluginwebhook "k8s.io/apiserver/plugin/pkg/audit/webhook"
)

func SetupAuditConfig(genericAPIServerConfig *genericapiserver.Config, auditOptions *genericapiserveroptions.AuditOptions) error {
	if auditOptions.PolicyFile != "" {
		p, err := policy.LoadPolicyFromFile(auditOptions.PolicyFile)
		if err != nil {
			return fmt.Errorf("loading audit policy file: %v", err)
		}
		policyRuleEvaluator := policy.NewPolicyRuleEvaluator(p)

		logBackend := buildLogAuditBackend(auditOptions.LogOptions)
		webhookBackend, err := buildWebhookAuditBackend(auditOptions.WebhookOptions)
		if err != nil {
			return err
		}
		backend := logBackend
		if backend == nil && webhookBackend != nil {
			backend = webhookBackend
		} else if webhookBackend != nil {
			backend = audit.Union(backend, webhookBackend)
		}
		genericAPIServerConfig.AuditBackend = backend
		genericAPIServerConfig.AuditPolicyRuleEvaluator = policyRuleEvaluator
	}
	return nil
}

func buildLogAuditBackend(o apiserveroptions.AuditLogOptions) audit.Backend {
	if o.Path == "" {
		return nil
	}
	var w io.Writer = os.Stdout
	if o.Path != "-" {
		w = &lumberjack.Logger{
			Filename:   o.Path,
			MaxAge:     o.MaxAge,
			MaxBackups: o.MaxBackups,
			MaxSize:    o.MaxSize,
		}
	}
	groupVersion, _ := schema.ParseGroupVersion(o.GroupVersionString)
	logBackend := pluginlog.NewBackend(w, o.Format, groupVersion)
	logBackend = pluginbuffered.NewBackend(logBackend, o.BatchOptions.BatchConfig)
	return logBackend
}

func buildWebhookAuditBackend(o apiserveroptions.AuditWebhookOptions) (audit.Backend, error) {
	if o.ConfigFile == "" {
		return nil, nil
	}
	groupVersion, _ := schema.ParseGroupVersion(o.GroupVersionString)
	webhook, err := pluginwebhook.NewBackend(o.ConfigFile, groupVersion, wait.Backoff{
		Steps:    5,
		Duration: o.InitialBackoff,
		Factor:   1.0,
		Jitter:   0.1,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("initializing audit webhook: %v", err)
	}
	webhook = pluginbuffered.NewBackend(webhook, o.BatchOptions.BatchConfig)
	if o.TruncateOptions.Enabled {
		webhook = plugintruncate.NewBackend(webhook, o.TruncateOptions.TruncateConfig, groupVersion)
	}
	return webhook, nil
}
