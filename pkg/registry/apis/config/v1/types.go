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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RegistryConfiguration contains the configuration for the Registry
type RegistryConfiguration struct {
	metav1.TypeMeta

	Storage  Storage  `json:"storage" yaml:"storage"`
	Security Security `json:"security" yaml:"security"`
	// +optional
	Redis         *Redis `json:"redis,omitempty" yaml:"redis,omitempty"`
	DefaultTenant string `json:"defaultTenant" yaml:"defaultTenant"`
	// +optional
	DomainSuffix  string `json:"domainSuffix,omitempty" yaml:"domainSuffix,omitempty"`
	HarborEnabled bool   `json:"harborEnabled,omitempty" yaml:"harborEnabled,omitempty"`
	HarborCAFile  string `json:"harborCAFile,omitempty" yaml:"harborCAFile,omitempty"`
}

type Storage struct {
	// +optional
	FileSystem *FileSystemStorage `json:"fileSystem,omitempty" yaml:"fileSystem,omitempty"`
	// +optional
	InMemory *InMemoryStorage `json:"inMemory,omitempty" yaml:"inMemory,omitempty"`
	// +optional
	S3 *S3Storage `json:"s3,omitempty" yaml:"s3,omitempty"`
	// +optional
	Delete *Delete `json:"delete,omitempty" yaml:"delete,omitempty"`
}

type FileSystemStorage struct {
	RootDirectory string `json:"rootDirectory" yaml:"rootDirectory"`
	// +optional
	MaxThreads *int64 `json:"maxThreads,omitempty" yaml:"maxThreads,omitempty"`
}

type InMemoryStorage struct{}

// S3StorageClass describes applied to each registry file.
type S3StorageClass string

const (
	S3StorageClassStandard          S3StorageClass = "STANDARD"
	S3StorageClassReducedRedundancy S3StorageClass = "REDUCED_REDUNDANCY"
)

type S3Storage struct {
	Bucket string `json:"bucket" yaml:"bucket"`
	Region string `json:"region" yaml:"region"`

	// +optional
	AccessKey *string `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	// +optional
	SecretKey *string `json:"secretKey,omitempty" yaml:"secretKey,omitempty"`
	// +optional
	RegionEndpoint *string `json:"regionEndpoint,omitempty" yaml:"regionEndpoint,omitempty"`
	// +optional
	Encrypt *bool `json:"encrypt,omitempty" yaml:"encrypt,omitempty"`
	// +optional
	KeyID *string `json:"keyID,omitempty" yaml:"keyID,omitempty"`
	// +optional
	Secure *bool `json:"secure,omitempty" yaml:"secure,omitempty"`
	// +optional
	SkipVerify *bool `json:"skipVerify,omitempty" yaml:"skipVerify,omitempty"`
	// +optional
	V4Auth *bool `json:"v4Auth,omitempty" yaml:"v4Auth,omitempty"`
	// +optional
	ChunkSize *int64 `json:"chunkSize,omitempty" yaml:"chunkSize,omitempty"`
	// +optional
	MultipartCopyChunkSize *int64 `json:"multipartCopyChunkSize,omitempty" yaml:"multipartCopyChunkSize,omitempty"`
	// +optional
	MultipartCopyMaxConcurrency *int64 `json:"multipartCopyMaxConcurrency,omitempty" yaml:"multipartCopyMaxConcurrency,omitempty"`
	// +optional
	MultipartCopyThresholdSize *int64 `json:"multipartCopyThresholdSize,omitempty" yaml:"multipartCopyThresholdSize,omitempty"`
	// +optional
	RootDirectory *string `json:"rootDirectory,omitempty" yaml:"rootDirectory,omitempty"`
	// +optional
	StorageClass *S3StorageClass `json:"storageClass,omitempty" yaml:"storageClass,omitempty"`
	// +optional
	UserAgent *string `json:"userAgent,omitempty" yaml:"userAgent,omitempty"`
	// +optional
	ObjectACL *string `json:"objectACL,omitempty" yaml:"objectACL,omitempty"`
}

// Delete cloud enable the deletion of image blobs and manifests by digest.
type Delete struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

type Security struct {
	TokenPrivateKeyFile string `json:"tokenPrivateKeyFile" yaml:"tokenPrivateKeyFile"`
	TokenPublicKeyFile  string `json:"tokenPublicKeyFile" yaml:"tokenPublicKeyFile"`
	// +optional
	TokenExpiredHours *int64 `json:"tokenExpiredHours,omitempty" yaml:"tokenExpiredHours,omitempty"`
	HTTPSecret        string `json:"httpSecret" yaml:"httpSecret"`
	AdminUsername     string `json:"adminUsername" yaml:"adminUsername"`
	AdminPassword     string `json:"adminPassword" yaml:"adminPassword"`
	// +optional
	EnableAnonymous *bool `json:"enableAnonymous" yaml:"enableAnonymous"`
}

// Redis configures the redis pool available to the registry cache.
type Redis struct {
	// Addr specifies the the redis instance available to the registry API server.
	Addr string `json:"addr" yaml:"addr"`
	// Password string to use when making a connection.
	Password string `json:"password" yaml:"password"`
	// DB specifies the database to connect to on the redis instance.
	DB int32 `json:"db" yaml:"db"`
	// +optional
	ReadTimeoutMillisecond *int64 `json:"readTimeoutMillisecond,omitempty" yaml:"readTimeoutMillisecond,omitempty"`
	// +optional
	DialTimeoutMillisecond *int64 `json:"dialTimeoutMillisecond,omitempty" yaml:"dialTimeoutMillisecond,omitempty"`
	// +optional
	WriteTimeoutMillisecond *int64 `json:"writeTimeoutMillisecond,omitempty" yaml:"writeTimeoutMillisecond,omitempty"`
	// PoolMaxIdle sets the maximum number of idle connections.
	// +optional
	PoolMaxIdle *int32 `json:"poolMaxIdle,omitempty" yaml:"poolMaxIdle,omitempty"`
	// PoolMaxActive sets the maximum number of connections that should be opened before
	// blocking a connection request.
	// +optional
	PoolMaxActive *int32 `json:"poolMaxActive,omitempty" yaml:"poolMaxActive,omitempty"`
	// PoolIdleTimeoutSeconds sets the amount time to wait before closing inactive connections.
	// +optional
	PoolIdleTimeoutSeconds *int64 `json:"poolIdleTimeoutSeconds,omitempty" yaml:"poolIdleTimeoutSeconds,omitempty"`
}
