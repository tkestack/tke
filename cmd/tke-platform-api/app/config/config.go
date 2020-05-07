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
	"io"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/audit"
	"k8s.io/apiserver/pkg/audit/policy"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/filters"
	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	pluginbuffered "k8s.io/apiserver/plugin/pkg/audit/buffered"
	pluginlog "k8s.io/apiserver/plugin/pkg/audit/log"
	plugintruncate "k8s.io/apiserver/plugin/pkg/audit/truncate"
	pluginwebhook "k8s.io/apiserver/plugin/pkg/audit/webhook"
	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	generatedopenapi "tkestack.io/tke/api/openapi"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/cmd/tke-platform-api/app/options"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/apiserver/authorization"
	"tkestack.io/tke/pkg/apiserver/debug"
	"tkestack.io/tke/pkg/apiserver/handler"
	"tkestack.io/tke/pkg/apiserver/openapi"
	"tkestack.io/tke/pkg/apiserver/storage"
	"tkestack.io/tke/pkg/platform/apiserver"
)

const (
	license = "Apache 2.0"
	title   = "Tencent Kubernetes Engine Platform API"
)

// Config is the running configuration structure of the TKE platform apiserver.
type Config struct {
	ServerName                     string
	GenericAPIServerConfig         *genericapiserver.Config
	VersionedSharedInformerFactory versionedinformers.SharedInformerFactory
	StorageFactory                 *serverstorage.DefaultStorageFactory
	PrivilegedUsername             string
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE platform apiserver command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	genericAPIServerConfig := genericapiserver.NewConfig(platform.Codecs)
	genericAPIServerConfig.BuildHandlerChainFunc = handler.BuildHandlerChain(nil)
	genericAPIServerConfig.LongRunningFunc = filters.BasicLongRunningRequestCheck(
		sets.NewString("watch", "proxy"),
		sets.NewString("attach", "exec", "proxy", "log", "portforward"),
	)
	genericAPIServerConfig.MergedResourceConfig = apiserver.DefaultAPIResourceConfigSource()
	genericAPIServerConfig.EnableIndex = false
	genericAPIServerConfig.EnableProfiling = false

	if opts.Audit.PolicyFile != "" {
		p, err := policy.LoadPolicyFromFile(opts.Audit.PolicyFile)
		if err != nil {
			return nil, fmt.Errorf("loading audit policy file: %v", err)
		}
		checker := policy.NewChecker(p)

		logBackend := buildLogAuditBackend(opts.Audit.LogOptions)
		webhookBackend, err := buildWebhookAuditBackend(opts.Audit.WebhookOptions)
		if err != nil {
			return nil, err
		}
		backend := logBackend
		if backend == nil && webhookBackend != nil {
			backend = webhookBackend
		} else if webhookBackend != nil {
			backend = audit.Union(backend, webhookBackend)
		}
		genericAPIServerConfig.AuditBackend = backend
		genericAPIServerConfig.AuditPolicyChecker = checker
	}

	if err := opts.Generic.ApplyTo(genericAPIServerConfig); err != nil {
		return nil, err
	}
	if err := opts.SecureServing.ApplyTo(&genericAPIServerConfig.SecureServing, &genericAPIServerConfig.LoopbackClientConfig); err != nil {
		return nil, err
	}

	openapi.SetupOpenAPI(genericAPIServerConfig, generatedopenapi.GetOpenAPIDefinitions, title, license, opts.Generic.ExternalHost, opts.Generic.ExternalPort)

	// storageFactory
	storageFactoryConfig := storage.NewFactoryConfig(platform.Codecs, platform.Scheme)
	storageFactoryConfig.APIResourceConfig = genericAPIServerConfig.MergedResourceConfig
	completedStorageFactoryConfig, err := storageFactoryConfig.Complete(opts.ETCD)
	if err != nil {
		return nil, err
	}
	storageFactory, err := completedStorageFactoryConfig.New()
	if err != nil {
		return nil, err
	}
	if err := opts.ETCD.ApplyWithStorageFactoryTo(storageFactory, genericAPIServerConfig); err != nil {
		return nil, err
	}

	// client config
	genericAPIServerConfig.LoopbackClientConfig.ContentConfig.ContentType = "application/vnd.kubernetes.protobuf"

	kubeClientConfig := genericAPIServerConfig.LoopbackClientConfig
	clientgoExternalClient, err := versionedclientset.NewForConfig(kubeClientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create real external clientset: %v", err)
	}
	versionedInformers := versionedinformers.NewSharedInformerFactory(clientgoExternalClient, 10*time.Minute)

	debug.SetupDebug(genericAPIServerConfig, opts.Debug)

	if err := authentication.SetupAuthentication(genericAPIServerConfig, opts.Authentication); err != nil {
		return nil, err
	}

	if err := authorization.SetupAuthorization(genericAPIServerConfig, opts.Authorization); err != nil {
		return nil, err
	}

	return &Config{
		ServerName:                     serverName,
		GenericAPIServerConfig:         genericAPIServerConfig,
		VersionedSharedInformerFactory: versionedInformers,
		StorageFactory:                 storageFactory,
		PrivilegedUsername:             opts.Authentication.PrivilegedUsername,
	}, nil
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
	webhook, err := pluginwebhook.NewBackend(o.ConfigFile, groupVersion, o.InitialBackoff)
	if err != nil {
		return nil, fmt.Errorf("initializing audit webhook: %v", err)
	}
	webhook = pluginbuffered.NewBackend(webhook, o.BatchOptions.BatchConfig)
	if o.TruncateOptions.Enabled {
		webhook = plugintruncate.NewBackend(webhook, o.TruncateOptions.TruncateConfig, groupVersion)
	}
	return webhook, nil
}
