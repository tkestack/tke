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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RegistryConfiguration contains the configuration for the Registry
type RegistryConfiguration struct {
	metav1.TypeMeta

	Storage  Storage
	Security Security
	// +optional
	Redis         *Redis
	DefaultTenant string
	// +optional
	DomainSuffix string
}

type Storage struct {
	// +optional
	FileSystem *FileSystemStorage
	// +optional
	InMemory *InMemoryStorage
	// +optional
	S3 *S3Storage
}

type FileSystemStorage struct {
	RootDirectory string
	// +optional
	MaxThreads *int64
}

type InMemoryStorage struct{}

// S3StorageClass describes applied to each registry file.
type S3StorageClass string

const (
	S3StorageClassStandard          S3StorageClass = "STANDARD"
	S3StorageClassReducedRedundancy S3StorageClass = "REDUCED_REDUNDANCY"
)

type S3Storage struct {
	Bucket string
	Region string

	// +optional
	AccessKey *string
	// +optional
	SecretKey *string
	// +optional
	RegionEndpoint *string
	// +optional
	Encrypt *bool
	// +optional
	KeyID *string
	// +optional
	Secure *bool
	// +optional
	SkipVerify *bool
	// +optional
	V4Auth *bool
	// +optional
	ChunkSize *int64
	// +optional
	MultipartCopyChunkSize *int64
	// +optional
	MultipartCopyMaxConcurrency *int64
	// +optional
	MultipartCopyThresholdSize *int64
	// +optional
	RootDirectory *string
	// +optional
	StorageClass *S3StorageClass
	// +optional
	UserAgent *string
	// +optional
	ObjectACL *string
}

type Security struct {
	TokenPrivateKeyFile string
	TokenPublicKeyFile  string
	// +optional
	TokenExpiredHours *int64
	HTTPSecret        string
	AdminUsername     string
	AdminPassword     string
}

// Redis configures the redis pool available to the registry cache.
type Redis struct {
	// Addr specifies the the redis instance available to the registry API server.
	Addr string
	// Password string to use when making a connection.
	Password string
	// DB specifies the database to connect to on the redis instance.
	DB int32
	// +optional
	ReadTimeoutMillisecond *int64
	// +optional
	DialTimeoutMillisecond *int64
	// +optional
	WriteTimeoutMillisecond *int64
	// PoolMaxIdle sets the maximum number of idle connections.
	// +optional
	PoolMaxIdle *int32
	// PoolMaxActive sets the maximum number of connections that should be opened before
	// blocking a connection request.
	// +optional
	PoolMaxActive *int32
	// PoolIdleTimeoutSeconds sets the amount time to wait before closing inactive connections.
	// +optional
	PoolIdleTimeoutSeconds *int64
}
