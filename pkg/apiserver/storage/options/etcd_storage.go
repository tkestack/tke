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

package options

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/healthz"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
)

const (
	flagETCDServersOverrides     = "etcd-servers-overrides"
	flagStorageMediaType         = "storage-media-type"
	flagDeleteCollectionWorkers  = "delete-collection-workers"
	flagEnableGarbageCollection  = "enable-garbage-collector"
	flagEnableWatchCache         = "watch-cache"
	flagDefaultWatchCacheSize    = "default-watch-cache-size"
	flagWatchCacheSizes          = "watch-cache-sizes"
	flagEncryptionProviderConfig = "encryption-provider-config"
)

const (
	configETCDServersOverrides     = "etcd.servers_overrides"
	configStorageMediaType         = "etcd.storage_media_type"
	configDeleteCollectionWorkers  = "etcd.delete_collection_workers"
	configEnableGarbageCollection  = "etcd.enable_garbage_collector"
	configEnableWatchCache         = "etcd.watch_cache"
	configDefaultWatchCacheSize    = "etcd.default_watch_cache_size"
	configWatchCacheSizes          = "etcd.watch_cache_sizes"
	configEncryptionProviderConfig = "etcd.encryption_provider_config"
)

// ETCDStorageOptions contains the options that storage backend by etcd.
type ETCDStorageOptions struct {
	*ETCDClientOptions
	ETCDServersOverrides             []string
	EncryptionProviderConfigFilePath string

	// To enable protobuf as storage format, it is enough
	// to set it to "application/vnd.kubernetes.protobuf".
	DefaultStorageMediaType string
	DeleteCollectionWorkers int
	EnableGarbageCollection bool

	// Set EnableWatchCache to false to disable all watch caches
	EnableWatchCache bool
	// Set DefaultWatchCacheSize to zero to disable watch caches for those resources that have no explicit cache size set
	DefaultWatchCacheSize int
	// WatchCacheSizes represents override to a given resource
	WatchCacheSizes []string
}

// NewETCDStorageOptions creates a Options object with default parameters.
func NewETCDStorageOptions(defaultETCDPathPrefix string) *ETCDStorageOptions {
	clientOptions := NewETCDClientOptions(defaultETCDPathPrefix)
	clientOptions.CountMetricPollPeriod = time.Minute
	return &ETCDStorageOptions{
		ETCDClientOptions:       clientOptions,
		DefaultStorageMediaType: "application/json",
		DeleteCollectionWorkers: 1,
		EnableGarbageCollection: true,
		EnableWatchCache:        true,
		DefaultWatchCacheSize:   100,
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *ETCDStorageOptions) AddFlags(fs *pflag.FlagSet) {
	o.ETCDClientOptions.AddFlags(fs)

	fs.StringSlice(flagETCDServersOverrides, o.ETCDServersOverrides, ""+
		"Per-resource etcd servers overrides, comma separated. The individual override "+
		"format: group/resource#servers, where servers are URLs, semicolon separated.")
	_ = viper.BindPFlag(configETCDServersOverrides, fs.Lookup(flagETCDServersOverrides))

	fs.String(flagEncryptionProviderConfig, o.EncryptionProviderConfigFilePath,
		"The file containing configuration for encryption providers to be used for storing secrets in etcd.")
	_ = viper.BindPFlag(configEncryptionProviderConfig, fs.Lookup(flagEncryptionProviderConfig))

	fs.String(flagStorageMediaType, o.DefaultStorageMediaType, ""+
		"The media type to use to store objects in storage. "+
		"Some resources or storage backends may only support a specific media type and will ignore this setting.")
	_ = viper.BindPFlag(configStorageMediaType, fs.Lookup(flagStorageMediaType))

	fs.Int(flagDeleteCollectionWorkers, o.DeleteCollectionWorkers,
		"Number of workers spawned for DeleteCollection call. These are used to speed up namespace cleanup.")
	_ = viper.BindPFlag(configDeleteCollectionWorkers, fs.Lookup(flagDeleteCollectionWorkers))

	fs.Bool(flagEnableGarbageCollection, o.EnableGarbageCollection, ""+
		"Enables the generic garbage collector. MUST be synced with the corresponding flag "+
		"of the kube-controller-manager.")
	_ = viper.BindPFlag(configEnableGarbageCollection, fs.Lookup(flagEnableGarbageCollection))

	fs.Bool(flagEnableWatchCache, o.EnableWatchCache,
		"Enable watch caching in the apiserver")
	_ = viper.BindPFlag(configEnableWatchCache, fs.Lookup(flagEnableWatchCache))

	fs.Int(flagDefaultWatchCacheSize, o.DefaultWatchCacheSize,
		"Default watch cache size. If zero, watch cache will be disabled for resources that do not have a default watch size set.")
	_ = viper.BindPFlag(configDefaultWatchCacheSize, fs.Lookup(flagDefaultWatchCacheSize))

	fs.StringSlice(flagWatchCacheSizes, o.WatchCacheSizes, ""+
		"List of watch cache sizes for every resource, comma separated. "+
		"The individual override format: resource[.group]#size, where resource is lowercase plural (no version), "+
		"group is optional, and size is a number. It takes effect when watch-cache is enabled. ")
	_ = viper.BindPFlag(configWatchCacheSizes, fs.Lookup(flagWatchCacheSizes))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *ETCDStorageOptions) ApplyFlags() []error {
	var errs []error

	errs = append(errs, o.ETCDClientOptions.ApplyFlags()...)

	o.ETCDServersOverrides = viper.GetStringSlice(configETCDServersOverrides)
	o.EncryptionProviderConfigFilePath = viper.GetString(configEncryptionProviderConfig)
	o.DeleteCollectionWorkers = viper.GetInt(configDeleteCollectionWorkers)
	o.EnableGarbageCollection = viper.GetBool(configEnableGarbageCollection)
	o.DefaultStorageMediaType = viper.GetString(configStorageMediaType)
	o.EnableWatchCache = viper.GetBool(configEnableWatchCache)
	o.DefaultWatchCacheSize = viper.GetInt(configDefaultWatchCacheSize)
	o.WatchCacheSizes = viper.GetStringSlice(configWatchCacheSizes)

	for _, override := range o.ETCDServersOverrides {
		tokens := strings.Split(override, "#")
		if len(tokens) != 2 {
			errs = append(errs, fmt.Errorf("--%s invalid, must be of format: group/resource#servers, where servers are URLs, semicolon separated", flagETCDServersOverrides))
			continue
		}

		apiResource := strings.Split(tokens[0], "/")
		if len(apiResource) != 2 {
			errs = append(errs, fmt.Errorf("--%s invalid, must be of format: group/resource#servers, where servers are URLs, semicolon separated", flagETCDServersOverrides))
			continue
		}

	}

	return errs
}

// ApplyWithStorageFactoryTo apply storage factory by etcd storage options.
func (o *ETCDStorageOptions) ApplyWithStorageFactoryTo(factory serverstorage.StorageFactory, c *server.Config) error {
	if err := o.addETCDHealthEndpoint(c); err != nil {
		return err
	}
	c.RESTOptionsGetter = &StorageFactoryRESTOptionsFactory{Options: *o, StorageFactory: factory}
	return nil
}

func (o *ETCDStorageOptions) addETCDHealthEndpoint(c *server.Config) error {
	healthCheck, err := o.ETCDClientOptions.NewHealthCheck()
	if err != nil {
		return err
	}
	c.HealthzChecks = append(c.HealthzChecks, healthz.NamedCheck("etcd", func(r *http.Request) error {
		return healthCheck()
	}))
	return nil
}

// StorageFactoryRESTOptionsFactory wrap the storage factory.
type StorageFactoryRESTOptionsFactory struct {
	Options        ETCDStorageOptions
	StorageFactory serverstorage.StorageFactory
}

// GetRESTOptions return the rest options.
func (f *StorageFactoryRESTOptionsFactory) GetRESTOptions(resource schema.GroupResource) (generic.RESTOptions, error) {
	storageConfig, err := f.StorageFactory.NewConfig(resource)
	if err != nil {
		return generic.RESTOptions{}, fmt.Errorf("unable to find storage destination for %v, due to %v", resource, err.Error())
	}

	ret := generic.RESTOptions{
		StorageConfig:           storageConfig,
		Decorator:               generic.UndecoratedStorage,
		DeleteCollectionWorkers: f.Options.DeleteCollectionWorkers,
		EnableGarbageCollection: f.Options.EnableGarbageCollection,
		ResourcePrefix:          f.StorageFactory.ResourcePrefix(resource),
		CountMetricPollPeriod:   f.Options.CountMetricPollPeriod,
	}
	if f.Options.EnableWatchCache {
		ret.Decorator = genericregistry.StorageWithCacher()
	}

	return ret, nil
}

// ParseWatchCacheSizes turns a list of cache size values into a map of group resources
// to requested sizes.
func ParseWatchCacheSizes(cacheSizes []string) (map[schema.GroupResource]int, error) {
	watchCacheSizes := make(map[schema.GroupResource]int)
	for _, c := range cacheSizes {
		tokens := strings.Split(c, "#")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("invalid value of watch cache size: %s", c)
		}

		size, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, fmt.Errorf("invalid size of watch cache size: %s", c)
		}
		if size < 0 {
			return nil, fmt.Errorf("watch cache size cannot be negative: %s", c)
		}

		watchCacheSizes[schema.ParseGroupResource(tokens[0])] = size
	}
	return watchCacheSizes, nil
}
