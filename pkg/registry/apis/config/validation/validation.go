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

package validation

import (
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
)

// ValidateRegistryConfiguration validates `rc` and returns an error if it is invalid
func ValidateRegistryConfiguration(rc *registryconfig.RegistryConfiguration) error {
	var allErrors []error

	storageCount := 0
	storageFld := field.NewPath("storage")

	if rc.Storage.S3 != nil {
		storageCount++
		subFld := storageFld.Child("s3")

		if rc.Storage.S3.Region == "" {
			allErrors = append(allErrors, field.Required(subFld.Child("region"), "must be specify"))
		}

		if rc.Storage.S3.Bucket == "" {
			allErrors = append(allErrors, field.Required(subFld.Child("bucket"), "must be specify"))
		}
	}

	if rc.Storage.FileSystem != nil {
		storageCount++
		subFld := storageFld.Child("fileSystem")

		if rc.Storage.FileSystem.RootDirectory == "" {
			allErrors = append(allErrors, field.Required(subFld.Child("rootDirectory"), "must be specify"))
		}
	}

	if rc.Storage.InMemory != nil {
		storageCount++
	}

	if storageCount == 0 {
		allErrors = append(allErrors, field.Required(storageFld, "at least 1 storage driver is required"))
	}

	securityFld := field.NewPath("security")
	if rc.Security.TokenPrivateKeyFile == "" {
		allErrors = append(allErrors, field.Required(securityFld.Child("tokenPrivateKeyFile"), "must be specify"))
	}
	if rc.Security.TokenPublicKeyFile == "" {
		allErrors = append(allErrors, field.Required(securityFld.Child("tokenPublicKeyFile"), "must be specify"))
	}

	if rc.DefaultTenant == "" {
		allErrors = append(allErrors, field.Required(field.NewPath("defaultTenant"), "must be specify"))
	}

	if rc.Redis != nil {
		redisFld := field.NewPath("redis")
		if rc.Redis.Addr == "" {
			allErrors = append(allErrors, field.Required(redisFld.Child("addr"), "must be specify"))
		}
	}

	return utilerrors.NewAggregate(allErrors)
}
