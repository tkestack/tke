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

	Storage  Storage  `json:"storage"`
	Security Security `json:"security"`
	// +optional
	Redis         *Redis `json:"redis,omitempty"`
	DefaultTenant string `json:"defaultTenant"`
	// +optional
	DomainSuffix string `json:"domainSuffix,omitempty"`
}

type Storage struct {
	// +optional
	FileSystem *FileSystemStorage `json:"fileSystem,omitempty"`
	// +optional
	InMemory *InMemoryStorage `json:"inMemory,omitempty"`
	// +optional
	S3 *S3Storage `json:"s3,omitempty"`
}

type FileSystemStorage struct {
	RootDirectory string `json:"rootDirectory"`
	// +optional
	MaxThreads *int64 `json:"maxThreads,omitempty"`
}

type InMemoryStorage struct{}

// S3StorageClass describes applied to each registry file.
type S3StorageClass string

const (
	S3StorageClassStandard          S3StorageClass = "STANDARD"
	S3StorageClassReducedRedundancy S3StorageClass = "REDUCED_REDUNDANCY"
)

type S3Storage struct {
	Bucket string `json:"bucket"`
	Region string `json:"region"`

	// +optional
	AccessKey *string `json:"accessKey,omitempty"`
	// +optional
	SecretKey *string `json:"secretKey,omitempty"`
	// +optional
	RegionEndpoint *string `json:"regionEndpoint,omitempty"`
	// +optional
	Encrypt *bool `json:"encrypt,omitempty"`
	// +optional
	KeyID *string `json:"keyID,omitempty"`
	// +optional
	Secure *bool `json:"secure,omitempty"`
	// +optional
	SkipVerify *bool `json:"skipVerify,omitempty"`
	// +optional
	V4Auth *bool `json:"v4Auth,omitempty"`
	// +optional
	ChunkSize *int64 `json:"chunkSize,omitempty"`
	// +optional
	MultipartCopyChunkSize *int64 `json:"multipartCopyChunkSize,omitempty"`
	// +optional
	MultipartCopyMaxConcurrency *int64 `json:"multipartCopyMaxConcurrency,omitempty"`
	// +optional
	MultipartCopyThresholdSize *int64 `json:"multipartCopyThresholdSize,omitempty"`
	// +optional
	RootDirectory *string `json:"rootDirectory,omitempty"`
	// +optional
	StorageClass *S3StorageClass `json:"storageClass,omitempty"`
	// +optional
	UserAgent *string `json:"userAgent,omitempty"`
	// +optional
	ObjectACL *string `json:"objectACL,omitempty"`
}

type Security struct {
	TokenPrivateKeyFile string `json:"tokenPrivateKeyFile"`
	TokenPublicKeyFile  string `json:"tokenPublicKeyFile"`
	// +optional
	TokenExpiredHours *int64 `json:"tokenExpiredHours,omitempty"`
	HTTPSecret        string `json:"httpSecret"`
	AdminUsername     string `json:"adminUsername"`
	AdminPassword     string `json:"adminPassword"`
}

// Redis configures the redis pool available to the registry cache.
type Redis struct {
	// Addr specifies the the redis instance available to the registry API server.
	Addr string `json:"addr"`
	// Password string to use when making a connection.
	Password string `json:"password"`
	// DB specifies the database to connect to on the redis instance.
	DB int32 `json:"db"`
	// +optional
	ReadTimeoutMillisecond *int64 `json:"readTimeoutMillisecond,omitempty"`
	// +optional
	DialTimeoutMillisecond *int64 `json:"dialTimeoutMillisecond,omitempty"`
	// +optional
	WriteTimeoutMillisecond *int64 `json:"writeTimeoutMillisecond,omitempty"`
	// PoolMaxIdle sets the maximum number of idle connections.
	// +optional
	PoolMaxIdle *int32 `json:"poolMaxIdle,omitempty"`
	// PoolMaxActive sets the maximum number of connections that should be opened before
	// blocking a connection request.
	// +optional
	PoolMaxActive *int32 `json:"poolMaxActive,omitempty"`
	// PoolIdleTimeoutSeconds sets the amount time to wait before closing inactive connections.
	// +optional
	PoolIdleTimeoutSeconds *int64 `json:"poolIdleTimeoutSeconds,omitempty"`
}
