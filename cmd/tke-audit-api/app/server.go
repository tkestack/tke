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

package app

import (
	genericapiserver "k8s.io/apiserver/pkg/server"
	"tkestack.io/tke/cmd/tke-audit-api/app/config"
	"tkestack.io/tke/pkg/audit"
)

// CreateServerChain creates the audit connected via delegation.
func CreateServerChain(cfg *config.Config, stopCh <-chan struct{}) (*genericapiserver.GenericAPIServer, error) {
	auditConfig := createAuditConfig(cfg)
	auditServer, err := CreateAudit(auditConfig, genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	return auditServer.GenericAPIServer, nil
}

// CreateAudit creates and wires a workable tke-audit-api.
func CreateAudit(auditConfig *audit.Config, delegateAPIServer genericapiserver.DelegationTarget) (*audit.Audit, error) {
	return auditConfig.Complete().New(delegateAPIServer)
}

func createAuditConfig(cfg *config.Config) *audit.Config {
	return &audit.Config{
		GenericConfig: &genericapiserver.RecommendedConfig{
			Config: *cfg.GenericAPIServerConfig,
		},
		ExtraConfig: audit.ExtraConfig{
			ServerName:  cfg.ServerName,
			AuditConfig: cfg.AuditConfig,
		},
	}
}
