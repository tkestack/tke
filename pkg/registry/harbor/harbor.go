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

package harbor

import (
	"k8s.io/apiserver/pkg/server/mux"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"

	"tkestack.io/tke/pkg/registry/harbor/handler"

	// import filesystem driver to store images
	_ "github.com/docker/distribution/registry/storage/driver/filesystem"
	// import in memory driver to store images
	_ "github.com/docker/distribution/registry/storage/driver/inmemory"
	// import s3 object storage driver to store images
	_ "github.com/docker/distribution/registry/storage/driver/s3-aws"
	// import token authentication controller
	_ "tkestack.io/tke/pkg/registry/distribution/auth/token"
)

// PathPrefix defines the path prefix for accessing the docker registry v2 server.
const PathPrefix = "/v2/"
const AuthPrefix = "/service/"

type Options struct {
	RegistryConfig *registryconfig.RegistryConfiguration
	ExternalHost   string
}

// IgnoredAuthPathPrefixes returns a list of path prefixes that does not need to
// go through the built-in authentication and authorization middleware of apiserver.
func IgnoreAuthPathPrefixes() []string {
	return []string{
		PathPrefix,
		AuthPrefix,
	}
}

// RegisterRoute to register the docker distribution server path prefix to apiserver.
func RegisterRoute(m *mux.PathRecorderMux, opts *Options) error {

	handler, err := handler.NewHandler("https://"+opts.RegistryConfig.DomainSuffix, opts.RegistryConfig.HarborCAFile, opts.ExternalHost)
	if err != nil {
		return err
	}

	m.HandlePrefix(PathPrefix, handler)
	m.HandlePrefix(AuthPrefix, handler)

	return nil
}
