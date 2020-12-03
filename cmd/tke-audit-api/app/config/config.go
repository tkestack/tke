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

package config

import (
	"fmt"
	"path/filepath"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/kube-openapi/pkg/common"
	generatedopenapi "tkestack.io/tke/api/openapi"
	"tkestack.io/tke/cmd/tke-audit-api/app/options"
	"tkestack.io/tke/pkg/apiserver"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/apiserver/authorization"
	"tkestack.io/tke/pkg/apiserver/handler"
	"tkestack.io/tke/pkg/apiserver/openapi"
	audit "tkestack.io/tke/pkg/audit/api"
	auditconfig "tkestack.io/tke/pkg/audit/apis/config"
	"tkestack.io/tke/pkg/audit/apis/config/validation"
	"tkestack.io/tke/pkg/audit/config/configfiles"
	auditopenapi "tkestack.io/tke/pkg/audit/openapi"
	utilfs "tkestack.io/tke/pkg/util/filesystem"
	"tkestack.io/tke/pkg/util/log"
)

const (
	license = "Apache 2.0"
	title   = "Tencent Kubernetes Engine Audit API"
)

// Config is the running configuration structure of the TKE controller manager.
type Config struct {
	ServerName             string
	GenericAPIServerConfig *genericapiserver.Config
	AuditConfig            *auditconfig.AuditConfiguration
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE apiserver command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	auditConfig, err := options.NewAuditConfiguration()
	if err != nil {
		log.Error("Failed create default audit configuration", log.Err(err))
		return nil, err
	}

	// load config file, if provided
	if configFile := opts.AuditConfig; len(configFile) > 0 {
		auditConfig, err = loadConfigFile(configFile)
		if err != nil {
			log.Error("Failed to load audit configuration file", log.String("configFile", configFile), log.Err(err))
			return nil, err
		}
	}
	if err := validation.ValidateAuditConfiguration(auditConfig); err != nil {
		log.Error("Failed to validate audit configuration", log.Err(err))
		return nil, err
	}

	genericAPIServerConfig := genericapiserver.NewConfig(apiserver.Codecs)
	var ignoredAuthPathPrefixes []string
	ignoredAuthPathPrefixes = append(ignoredAuthPathPrefixes, audit.IgnoredAuthPathPrefixes()...)
	genericAPIServerConfig.BuildHandlerChainFunc = handler.BuildHandlerChain(ignoredAuthPathPrefixes, nil, nil)
	genericAPIServerConfig.EnableIndex = false
	genericAPIServerConfig.EnableDiscovery = false
	genericAPIServerConfig.EnableProfiling = false

	if err := opts.Generic.ApplyTo(genericAPIServerConfig); err != nil {
		return nil, err
	}
	if err := opts.SecureServing.ApplyTo(&genericAPIServerConfig.SecureServing, &genericAPIServerConfig.LoopbackClientConfig); err != nil {
		return nil, err
	}

	openapi.SetupOpenAPI(genericAPIServerConfig, func(callback common.ReferenceCallback) map[string]common.OpenAPIDefinition {
		result := make(map[string]common.OpenAPIDefinition)
		generated := generatedopenapi.GetOpenAPIDefinitions(callback)
		for k, v := range generated {
			result[k] = v
		}
		customs := auditopenapi.GetOpenAPIDefinitions(callback)
		for k, v := range customs {
			result[k] = v
		}
		return result
	}, title, license, opts.Generic.ExternalHost, opts.Generic.ExternalPort)

	if err := authentication.SetupAuthentication(genericAPIServerConfig, opts.Authentication); err != nil {
		return nil, err
	}

	if err := authorization.SetupAuthorization(genericAPIServerConfig, opts.Authorization); err != nil {
		return nil, err
	}

	return &Config{
		ServerName:             serverName,
		GenericAPIServerConfig: genericAPIServerConfig,
		AuditConfig:            auditConfig,
	}, nil
}

func loadConfigFile(name string) (*auditconfig.AuditConfiguration, error) {
	const errFmt = "failed to load audit config file %s, error %v"
	// compute absolute path based on current working dir
	auditConfigFile, err := filepath.Abs(name)
	if err != nil {
		return nil, fmt.Errorf(errFmt, name, err)
	}
	loader, err := configfiles.NewFsLoader(utilfs.DefaultFs{}, auditConfigFile)
	if err != nil {
		return nil, fmt.Errorf(errFmt, name, err)
	}
	kc, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf(errFmt, name, err)
	}
	return kc, err
}
