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

// Portions Copyright 2014 The Kubernetes Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/server/options/encryptionconfig"
	"k8s.io/apiserver/pkg/server/resourceconfig"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"strings"
	storageoptions "tkestack.io/tke/pkg/apiserver/storage/options"
)

// SpecialDefaultResourcePrefixes are prefixes compiled into Kubernetes.
var SpecialDefaultResourcePrefixes = map[schema.GroupResource]string{}

// NewFactoryConfig creates the default FactoryConfig object.
func NewFactoryConfig(codecs runtime.StorageSerializer, scheme *runtime.Scheme) *FactoryConfig {
	var resources []schema.GroupVersionResource
	return &FactoryConfig{
		Serializer:                codecs,
		DefaultResourceEncoding:   serverstorage.NewDefaultResourceEncodingConfig(scheme),
		ResourceEncodingOverrides: resources,
	}
}

// FactoryConfig represents the configuration of etcd backend storage.
type FactoryConfig struct {
	StorageConfig                    storagebackend.Config
	APIResourceConfig                *serverstorage.ResourceConfig
	DefaultResourceEncoding          *serverstorage.DefaultResourceEncodingConfig
	DefaultStorageMediaType          string
	Serializer                       runtime.StorageSerializer
	ResourceEncodingOverrides        []schema.GroupVersionResource
	ETCDServersOverrides             []string
	EncryptionProviderConfigFilePath string
}

// Complete takes the command arguments and factory and infers any remaining
// options.
func (c *FactoryConfig) Complete(etcdOptions *storageoptions.ETCDStorageOptions) (*CompletedFactoryConfig, error) {
	c.StorageConfig = storagebackend.Config{
		Type:   storagebackend.StorageTypeETCD3,
		Prefix: etcdOptions.Prefix,
		Transport: storagebackend.TransportConfig{
			ServerList: etcdOptions.ServerList,
			KeyFile:    etcdOptions.KeyFile,
			CertFile:   etcdOptions.CertFile,
			CAFile:     etcdOptions.CAFile,
		},
		Paging:                etcdOptions.Paging,
		Codec:                 etcdOptions.Codec,
		EncodeVersioner:       etcdOptions.EncodeVersioner,
		Transformer:           etcdOptions.Transformer,
		CompactionInterval:    etcdOptions.CompactionInterval,
		CountMetricPollPeriod: etcdOptions.CountMetricPollPeriod,
	}
	c.DefaultStorageMediaType = etcdOptions.DefaultStorageMediaType
	c.ETCDServersOverrides = etcdOptions.ETCDServersOverrides
	c.EncryptionProviderConfigFilePath = etcdOptions.EncryptionProviderConfigFilePath
	return &CompletedFactoryConfig{c}, nil
}

// CompletedFactoryConfig represents the configuration of etcd backend storage.
type CompletedFactoryConfig struct {
	*FactoryConfig
}

// New creates the DefaultStorageFactory object and returns it.
func (c *CompletedFactoryConfig) New() (*serverstorage.DefaultStorageFactory, error) {
	resourceEncodingConfig := resourceconfig.MergeResourceEncodingConfigs(c.DefaultResourceEncoding, c.ResourceEncodingOverrides)
	storageFactory := serverstorage.NewDefaultStorageFactory(
		c.StorageConfig,
		c.DefaultStorageMediaType,
		c.Serializer,
		resourceEncodingConfig,
		c.APIResourceConfig,
		SpecialDefaultResourcePrefixes)

	for _, override := range c.ETCDServersOverrides {
		tokens := strings.Split(override, "#")
		apiresource := strings.Split(tokens[0], "/")

		group := apiresource[0]
		resource := apiresource[1]
		groupResource := schema.GroupResource{Group: group, Resource: resource}

		servers := strings.Split(tokens[1], ";")
		storageFactory.SetEtcdLocation(groupResource, servers)
	}
	if len(c.EncryptionProviderConfigFilePath) != 0 {
		transformerOverrides, err := encryptionconfig.GetTransformerOverrides(c.EncryptionProviderConfigFilePath)
		if err != nil {
			return nil, err
		}
		for groupResource, transformer := range transformerOverrides {
			storageFactory.SetTransformer(groupResource, transformer)
		}
	}
	return storageFactory, nil
}
