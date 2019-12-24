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
	// cmlogger "helm.sh/chartmuseum/pkg/chartmuseum/logger"
	// cmrouter "helm.sh/chartmuseum/pkg/chartmuseum/router"
	// "helm.sh/chartmuseum/pkg/chartmuseum/server/multitenant"
	"k8s.io/apiserver/pkg/server/mux"
	restclient "k8s.io/client-go/rest"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
)

const (
	// PathPrefix defines the path prefix for accessing the chartmuseum server.
	PathPrefix = "/chart/"
	// MaxUploadSize defines max size of post body (in bytes).
	MaxUploadSize = 20 * 1024 * 1024 * 1024
)

// IgnoredAuthPathPrefixes returns a list of path prefixes that does not need to
// go through the built-in authentication and authorization middleware of apiserver.
func IgnoredAuthPathPrefixes() []string {
	return []string{
		PathPrefix,
	}
}

type Options struct {
	RegistryConfig       *registryconfig.RegistryConfiguration
	LoopbackClientConfig *restclient.Config
	OIDCIssuerURL        string
	OIDCTokenReviewPath  string
	OIDCCAFile           string
	ExternalScheme       string
}

// RegisterRoute to register the chartmuseum server path prefix to apiserver.
func RegisterRoute(m *mux.PathRecorderMux, opts *Options) error {
	return nil
	// chartMuseumConfig, err := buildChartMuseumConfig(opts)
	// if err != nil {
	// 	log.Error("Failed to initialize chartmuseum server configuration", log.Err(err))
	// 	return err
	// }
	// multiTenantServer, err := multitenant.NewMultiTenantServer(*chartMuseumConfig)
	// if err != nil {
	// 	log.Error("Failed to create chartmuseum server", log.Err(err))
	// 	return err
	// }
	//
	// // add handler chain
	// var chainErr error
	// var handler http.Handler
	// handler, chainErr = authorization.WithAuthorization(multiTenantServer.Router, &authorization.Options{
	// 	AdminUsername:  opts.RegistryConfig.Security.AdminUsername,
	// 	ExternalScheme: opts.ExternalScheme,
	// 	LoopbackConfig: opts.LoopbackClientConfig,
	// })
	// if chainErr != nil {
	// 	return chainErr
	// }
	// handler, chainErr = authentication.WithAuthentication(handler, &authentication.Options{
	// 	SecurityConfig:  &opts.RegistryConfig.Security,
	// 	ExternalScheme:  opts.ExternalScheme,
	// 	OIDCIssuerURL:   opts.OIDCIssuerURL,
	// 	OIDCCAFile:      opts.OIDCCAFile,
	// 	TokenReviewPath: opts.OIDCTokenReviewPath,
	// })
	// if chainErr != nil {
	// 	return chainErr
	// }
	// handler = tenant.WithTenant(handler, PathPrefix, opts.RegistryConfig.DomainSuffix, opts.RegistryConfig.DefaultTenant)
	// handler = request.WithRequestID(handler)
	// m.HandlePrefix(PathPrefix, handler)
	//
	// return nil
}

//
// func buildChartMuseumConfig(opts *Options) (*multitenant.MultiTenantServerOptions, error) {
// 	// initialize logger
// 	zapLogger := log.ZapLogger()
// 	if zapLogger == nil {
// 		return nil, fmt.Errorf("logger has not been initialized")
// 	}
// 	logger := &cmlogger.Logger{SugaredLogger: zapLogger.Sugar()}
//
// 	// initialize router
// 	router := cmrouter.NewRouter(cmrouter.RouterOptions{
// 		Logger:        logger,
// 		ContextPath:   strings.TrimSuffix(PathPrefix, "/"),
// 		Depth:         2,
// 		EnableMetrics: true,
// 		MaxUploadSize: MaxUploadSize,
// 	})
//
// 	// initialize storage backend
// 	storageBackend, err := buildStorageConfiguration(opts)
// 	if err != nil {
// 		log.Error("Failed to create storage backend for charts", log.Err(err))
// 		return nil, err
// 	}
//
// 	// create server options
// 	return &multitenant.MultiTenantServerOptions{
// 		Router:              router,
// 		Logger:              logger,
// 		StorageBackend:      storageBackend,
// 		EnableAPI:           true,
// 		AllowForceOverwrite: true,
// 		AllowOverwrite:      false,
// 	}, nil
// }
//
// func buildStorageConfiguration(opts *Options) (storage.Backend, error) {
// 	var backend storage.Backend
// 	var err error
// 	storageCfg := &opts.RegistryConfig.Storage
// 	if storageCfg.FileSystem != nil {
// 		backend = storage.Backend(storage.NewLocalFilesystemBackend(storageCfg.FileSystem.RootDirectory))
// 	} else if storageCfg.S3 != nil {
// 		backend, err = buildS3StorageConfiguration(opts.RegistryConfig.Storage.S3)
// 	}
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if backend == nil {
// 		return nil, fmt.Errorf("no storage backend specified")
// 	}
// 	return backend, nil
// }
//
// func buildS3StorageConfiguration(cfg *registryconfig.S3Storage) (storage.Backend, error) {
// 	awsConfig := aws.NewConfig()
// 	sess, err := session.NewSession()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create new s3 storage session: %v", err)
// 	}
//
// 	var accessKey, secretKey string
// 	if cfg.AccessKey != nil {
// 		accessKey = *cfg.AccessKey
// 	}
// 	if cfg.SecretKey != nil {
// 		secretKey = *cfg.SecretKey
// 	}
//
// 	creds := credentials.NewChainCredentials([]credentials.Provider{
// 		&credentials.StaticProvider{
// 			Value: credentials.Value{
// 				AccessKeyID:     accessKey,
// 				SecretAccessKey: secretKey,
// 				SessionToken:    "",
// 			},
// 		},
// 		&credentials.EnvProvider{},
// 		&credentials.SharedCredentialsProvider{},
// 		&ec2rolecreds.EC2RoleProvider{Client: ec2metadata.New(sess)},
// 	})
//
// 	if cfg.RegionEndpoint != nil && *cfg.RegionEndpoint != "" {
// 		awsConfig.WithS3ForcePathStyle(true)
// 		awsConfig.WithEndpoint(*cfg.RegionEndpoint)
// 	}
//
// 	awsConfig.WithCredentials(creds)
// 	awsConfig.WithRegion(cfg.Region)
//
// 	secure := true
// 	if cfg.Secure != nil && !*cfg.Secure {
// 		secure = false
// 	}
// 	awsConfig.WithDisableSSL(!secure)
//
// 	var userAgent string
// 	var skipVerify bool
// 	if cfg.UserAgent != nil {
// 		userAgent = *cfg.UserAgent
// 	}
// 	if cfg.SkipVerify != nil && *cfg.SkipVerify {
// 		skipVerify = true
// 	}
// 	if userAgent != "" || skipVerify {
// 		httpTransport := http.DefaultTransport
// 		if skipVerify {
// 			httpTransport = &http.Transport{
// 				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 			}
// 		}
// 		if userAgent != "" {
// 			awsConfig.WithHTTPClient(&http.Client{
// 				Transport: transport.NewTransport(httpTransport, transport.NewHeaderRequestModifier(http.Header{http.CanonicalHeaderKey("User-Agent"): []string{userAgent}})),
// 			})
// 		} else {
// 			awsConfig.WithHTTPClient(&http.Client{
// 				Transport: transport.NewTransport(httpTransport),
// 			})
// 		}
// 	}
//
// 	sess, err = session.NewSession(awsConfig)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create new session with aws config: %v", err)
// 	}
// 	s3obj := s3.New(sess)
//
// 	// enable S3 compatible signature v2 signing instead
// 	if cfg.V4Auth != nil && !*cfg.V4Auth {
// 		setS3StorageV2Handlers(s3obj)
// 	}
//
// 	return &storage.AmazonS3Backend{
// 		Bucket:     cfg.Bucket,
// 		Client:     s3obj,
// 		Downloader: s3manager.NewDownloaderWithClient(s3obj),
// 		Prefix:     "",
// 		Uploader:   s3manager.NewUploaderWithClient(s3obj),
// 		SSE:        "",
// 	}, nil
// }
//
// // setS3StorageV2Handlers will setup v2 signature signing on the S3 driver
// func setS3StorageV2Handlers(svc *s3.S3) {
// 	svc.Handlers.Build.PushBack(func(r *s3request.Request) {
// 		parsedURL, err := url.Parse(r.HTTPRequest.URL.String())
// 		if err != nil {
// 			log.Fatalf("Failed to parse URL: %v", err)
// 		}
// 		r.HTTPRequest.URL.Opaque = parsedURL.Path
// 	})
//
// 	svc.Handlers.Sign.Clear()
// 	svc.Handlers.Sign.PushBack(distributions3.Sign)
// 	svc.Handlers.Sign.PushBackNamed(corehandlers.BuildContentLengthHandler)
// }
