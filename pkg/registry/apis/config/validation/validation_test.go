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
	"testing"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
)

func TestValidateRegistryConfiguration(t *testing.T) {
	successCase := &registryconfig.RegistryConfiguration{
		Storage: registryconfig.Storage{
			FileSystem: &registryconfig.FileSystemStorage{
				RootDirectory: "fake_dir",
			},
		},
		Security: registryconfig.Security{
			TokenPrivateKeyFile: "fake_private_key.crt",
			TokenPublicKeyFile:  "fake_public_key.pem",
		},
		DefaultTenant: "default",
	}
	if allErrors := ValidateRegistryConfiguration(successCase); allErrors != nil {
		t.Errorf("expect no errors, got %v", allErrors)
	}

	errorCase := &registryconfig.RegistryConfiguration{
		Storage: registryconfig.Storage{
			FileSystem: &registryconfig.FileSystemStorage{
				RootDirectory: "",
			},
			S3: &registryconfig.S3Storage{
				Bucket: "",
			},
		},
		Redis: &registryconfig.Redis{},
	}
	const numErrs = 7
	if allErrors := ValidateRegistryConfiguration(errorCase); len(allErrors.(utilerrors.Aggregate).Errors()) != numErrs {
		t.Errorf("expect %d errors, got %v", numErrs, len(allErrors.(utilerrors.Aggregate).Errors()))
	}
}
