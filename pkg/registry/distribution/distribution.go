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

package distribution

import (
	"fmt"
	"github.com/docker/distribution/configuration"
	"github.com/docker/distribution/registry/handlers"
	"k8s.io/apiserver/pkg/server/mux"
	restclient "k8s.io/client-go/rest"
	"net/http"
	"time"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	"tkestack.io/tke/pkg/registry/distribution/auth"
	rcontext "tkestack.io/tke/pkg/registry/distribution/context"
	"tkestack.io/tke/pkg/registry/distribution/notification"
	"tkestack.io/tke/pkg/registry/distribution/tenant"
	"tkestack.io/tke/pkg/util/transport"

	// import filesystem driver to store images
	_ "github.com/docker/distribution/registry/storage/driver/filesystem"
	// import in memory driver to store images
	_ "github.com/docker/distribution/registry/storage/driver/inmemory"
	// import s3 object storage driver to store images
	_ "github.com/docker/distribution/registry/storage/driver/s3-aws"
	// import token authentication controller
	_ "tkestack.io/tke/pkg/registry/distribution/auth/token"
)

const PathPrefix = "/v2/"
const APIPrefix = "/registry/"

const (
	notificationName = "tke"
)

type Options struct {
	RegistryConfig       *registryconfig.RegistryConfiguration
	ExternalScheme       string
	LoopbackClientConfig *restclient.Config
	OIDCIssuerURL        string
	OIDCTokenReviewPath  string
	OIDCCAFile           string
}

func IgnoreAuthPathPrefixes() []string {
	return []string{
		PathPrefix,
		auth.Path,
		notification.Path,
	}
}

func RegisterRoute(m *mux.PathRecorderMux, opts *Options) error {
	distConfig, err := buildDistributionConfig(opts)
	if err != nil {
		return err
	}

	distCtx := rcontext.BuildDistributionContext()
	distHandler := handlers.NewApp(distCtx, distConfig)
	wrappedDistHandler := tenant.WithTenant(distHandler, PathPrefix, opts.RegistryConfig.DomainSuffix, opts.RegistryConfig.DefaultTenant)
	m.HandlePrefix(PathPrefix, wrappedDistHandler)

	authHandler, err := auth.NewHandler(&auth.Options{
		SecurityConfig:  &opts.RegistryConfig.Security,
		OIDCIssuerURL:   opts.OIDCIssuerURL,
		OIDCCAFile:      opts.OIDCCAFile,
		TokenReviewPath: opts.OIDCTokenReviewPath,
		DomainSuffix:    opts.RegistryConfig.DomainSuffix,
		DefaultTenant:   opts.RegistryConfig.DefaultTenant,
		LoopbackConfig:  opts.LoopbackClientConfig,
	})
	if err != nil {
		return err
	}
	m.Handle(auth.Path, authHandler)

	notificationHandler, err := notification.NewHandler(opts.LoopbackClientConfig)
	if err != nil {
		return err
	}
	m.Handle(notification.Path, notificationHandler)

	return nil
}

func buildDistributionConfig(opts *Options) (*configuration.Configuration, error) {
	dist := &configuration.Configuration{}

	dist.Storage = buildStorageConfiguration(opts)
	dist.Auth = buildAuthConfiguration(opts)

	endpoints, err := buildNotificationsConfiguration(opts)
	if err != nil {
		return nil, err
	}
	dist.Notifications.Endpoints = endpoints
	dist.HTTP.Secret = opts.RegistryConfig.Security.HTTPSecret
	dist.Compatibility.Schema1.Enabled = false

	if opts.RegistryConfig.Redis != nil {
		redisCfg := opts.RegistryConfig.Redis
		dist.Redis.Addr = redisCfg.Addr
		dist.Redis.Password = redisCfg.Password
		dist.Redis.DB = int(redisCfg.DB)
		if redisCfg.DialTimeoutMillisecond != nil {
			dist.Redis.DialTimeout = time.Duration(*redisCfg.DialTimeoutMillisecond) * time.Millisecond
		}
		if redisCfg.ReadTimeoutMillisecond != nil {
			dist.Redis.ReadTimeout = time.Duration(*redisCfg.ReadTimeoutMillisecond) * time.Millisecond
		}
		if redisCfg.WriteTimeoutMillisecond != nil {
			dist.Redis.WriteTimeout = time.Duration(*redisCfg.WriteTimeoutMillisecond) * time.Millisecond
		}
		if redisCfg.PoolMaxIdle != nil {
			dist.Redis.Pool.MaxIdle = int(*redisCfg.PoolMaxIdle)
		}
		if redisCfg.PoolMaxActive != nil {
			dist.Redis.Pool.MaxActive = int(*redisCfg.PoolMaxActive)
		}
		if redisCfg.PoolIdleTimeoutSeconds != nil {
			dist.Redis.Pool.IdleTimeout = time.Duration(*redisCfg.PoolIdleTimeoutSeconds) * time.Second
		}
	}

	return dist, nil
}

func buildNotificationsConfiguration(opts *Options) ([]configuration.Endpoint, error) {
	url := fmt.Sprintf("%s%s", opts.LoopbackClientConfig.Host, notification.Path)
	tlsConfig, err := restclient.TLSConfigFor(opts.LoopbackClientConfig)
	if err != nil {
		return nil, err
	}
	http.DefaultTransport = transport.Transport(tlsConfig)
	return []configuration.Endpoint{
		{
			Name:      notificationName,
			Disabled:  false,
			URL:       url,
			Timeout:   3 * time.Second,
			Threshold: 5,
			Backoff:   1 * time.Second,
		},
	}, nil
}

func buildAuthConfiguration(opts *Options) map[string]configuration.Parameters {
	authConfig := make(map[string]configuration.Parameters)
	authParams := map[string]interface{}{
		"issuer":         auth.Issuer,
		"service":        auth.Service,
		"scheme":         opts.ExternalScheme,
		"authpath":       auth.Path,
		"rootcertbundle": opts.RegistryConfig.Security.TokenPublicKeyFile,
	}
	authConfig["tkestack"] = authParams
	return authConfig
}

func buildStorageConfiguration(opts *Options) map[string]configuration.Parameters {
	storage := make(map[string]configuration.Parameters)
	if opts.RegistryConfig.Redis != nil {
		cache := make(map[string]interface{}, 1)
		cache["blobdescriptor"] = "redis"
		storage["cache"] = cache
	}

	storageCfg := &opts.RegistryConfig.Storage
	if storageCfg.FileSystem != nil {
		fileSystem := make(map[string]interface{})
		fileSystem["rootdirectory"] = storageCfg.FileSystem.RootDirectory

		if storageCfg.FileSystem.MaxThreads != nil {
			fileSystem["maxthreads"] = *storageCfg.FileSystem.MaxThreads
		}
		storage["filesystem"] = fileSystem
	}

	if storageCfg.InMemory != nil {
		inMemory := make(map[string]interface{})
		storage["inmemory"] = inMemory
	}

	if storageCfg.S3 != nil {
		s3 := make(map[string]interface{})
		s3["bucket"] = storageCfg.S3.Bucket
		s3["region"] = storageCfg.S3.Region

		if storageCfg.S3.AccessKey != nil {
			s3["accesskey"] = *storageCfg.S3.AccessKey
		}
		if storageCfg.S3.SecretKey != nil {
			s3["secretkey"] = *storageCfg.S3.SecretKey
		}
		if storageCfg.S3.RegionEndpoint != nil {
			s3["regionendpoint"] = *storageCfg.S3.RegionEndpoint
		}
		if storageCfg.S3.Encrypt != nil {
			s3["encrypt"] = *storageCfg.S3.Encrypt
		}
		if storageCfg.S3.KeyID != nil {
			s3["keyid"] = *storageCfg.S3.KeyID
		}
		if storageCfg.S3.Secure != nil {
			s3["secure"] = *storageCfg.S3.Secure
		}
		if storageCfg.S3.SkipVerify != nil {
			s3["skipverify"] = *storageCfg.S3.SkipVerify
		}
		if storageCfg.S3.V4Auth != nil {
			s3["v4auth"] = *storageCfg.S3.V4Auth
		}
		if storageCfg.S3.ChunkSize != nil {
			s3["chunksize"] = *storageCfg.S3.ChunkSize
		}
		if storageCfg.S3.MultipartCopyChunkSize != nil {
			s3["multipartcopychunksize"] = *storageCfg.S3.MultipartCopyChunkSize
		}
		if storageCfg.S3.MultipartCopyMaxConcurrency != nil {
			s3["multipartcopymaxconcurrency"] = *storageCfg.S3.MultipartCopyMaxConcurrency
		}
		if storageCfg.S3.MultipartCopyThresholdSize != nil {
			s3["multipartcopythresholdsize"] = *storageCfg.S3.MultipartCopyThresholdSize
		}
		if storageCfg.S3.RootDirectory != nil {
			s3["rootdirectory"] = *storageCfg.S3.RootDirectory
		}
		if storageCfg.S3.StorageClass != nil {
			s3["storageclass"] = string(*storageCfg.S3.StorageClass)
		}
		if storageCfg.S3.UserAgent != nil {
			s3["useragent"] = *storageCfg.S3.UserAgent
		}
		if storageCfg.S3.ObjectACL != nil {
			s3["objectacl"] = *storageCfg.S3.ObjectACL
		}

		storage["s3"] = s3
	}
	return storage
}
