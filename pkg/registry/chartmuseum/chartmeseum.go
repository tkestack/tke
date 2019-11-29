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

package chartmuseum

import (
	"fmt"
	"github.com/chartmuseum/storage"
	cmlogger "helm.sh/chartmuseum/pkg/chartmuseum/logger"
	cmrouter "helm.sh/chartmuseum/pkg/chartmuseum/router"
	"helm.sh/chartmuseum/pkg/chartmuseum/server/multitenant"
	"k8s.io/apiserver/pkg/server/mux"
	"strings"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	"tkestack.io/tke/pkg/registry/chartmuseum/request"
	"tkestack.io/tke/pkg/registry/chartmuseum/tenant"
	"tkestack.io/tke/pkg/util/log"
)

const (
	// PathPrefix defines the path prefix for accessing the chartmuseum server.
	PathPrefix = "/chart/"
	// MaxUploadSize defines max size of post body (in bytes).
	MaxUploadSize = 20 * 1024 * 1024
)

// IgnoredAuthPathPrefixes returns a list of path prefixes that does not need to
// go through the built-in authentication and authorization middleware of apiserver.
func IgnoredAuthPathPrefixes() []string {
	return []string{
		PathPrefix,
	}
}

type Options struct {
	RegistryConfig *registryconfig.RegistryConfiguration
}

// RegisterRoute to register the chartmuseum server path prefix to apiserver.
func RegisterRoute(m *mux.PathRecorderMux, opts *Options) error {
	chartMuseumConfig, err := buildChartMuseumConfig(opts)
	if err != nil {
		log.Error("Failed to initialize chartmuseum server configuration", log.Err(err))
		return err
	}
	multiTenantServer, err := multitenant.NewMultiTenantServer(*chartMuseumConfig)
	if err != nil {
		log.Error("Failed to create chartmuseum server", log.Err(err))
		return err
	}
	wrappedChartHandler := tenant.WithTenant(multiTenantServer.Router, PathPrefix, opts.RegistryConfig.DomainSuffix, opts.RegistryConfig.DefaultTenant)
	wrappedChartHandler = request.WithRequestID(wrappedChartHandler)
	m.HandlePrefix(PathPrefix, wrappedChartHandler)
	return nil
}

func buildChartMuseumConfig(opts *Options) (*multitenant.MultiTenantServerOptions, error) {
	// initialize logger
	zapLogger := log.ZapLogger()
	if zapLogger == nil {
		return nil, fmt.Errorf("logger has not been initialized")
	}
	logger := &cmlogger.Logger{SugaredLogger: zapLogger.Sugar()}

	// initialize router
	router := cmrouter.NewRouter(cmrouter.RouterOptions{
		Logger:        logger,
		ContextPath:   strings.TrimSuffix(PathPrefix, "/"),
		Depth:         2,
		EnableMetrics: true,
		MaxUploadSize: MaxUploadSize,
	})

	// initialize storage backend
	storageBackend, err := buildStorageConfiguration(opts)
	if err != nil {
		log.Error("Failed to create storage backend for charts", log.Err(err))
		return nil, err
	}

	// create server options
	return &multitenant.MultiTenantServerOptions{
		Router:              router,
		Logger:              logger,
		StorageBackend:      storageBackend,
		EnableAPI:           true,
		AllowForceOverwrite: true,
		AllowOverwrite:      false,
	}, nil
}

func buildStorageConfiguration(opts *Options) (storage.Backend, error) {
	var backend storage.Backend
	storageCfg := &opts.RegistryConfig.Storage
	if storageCfg.FileSystem != nil {
		backend = storage.Backend(storage.NewLocalFilesystemBackend(storageCfg.FileSystem.RootDirectory))
	}

	if backend == nil {
		return nil, fmt.Errorf("no storage backend specified")
	}
	return backend, nil
}
