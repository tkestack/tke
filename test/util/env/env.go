/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package env

import (
	"os"
	"path"
	"strings"
	"tkestack.io/tke/cmd/tke-installer/app/installer/images"
	"tkestack.io/tke/pkg/spec"

	"github.com/joho/godotenv"
)

const (
	envFile = "tke.env"
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	godotenv.Overload(path.Join(home, envFile), envFile) // for local dev
}

const (
	VERSION              = "VERSION"
	PROVIDERRESVERSION   = "PROVIDERRESVERSION"
	K8SVERSION           = "K8SVERSION"
	KUBECONFIG           = "KUBECONFIG"
	SECRETID             = "SECRET_ID"
	SECRETKEY            = "SECRET_KEY"
	REGION               = "REGION"
	CREATEINSTANCESPARAM = "CREATE_INSTANCES_PARAM"
	PASSWORD             = "PASSWORD"
	NEEDDELETE           = "NEED_DELETE"
)

func ImageVersion() string {
	return os.Getenv(VERSION)
}

func Kubeconfig() string {
	return os.Getenv(KUBECONFIG)
}

func SecretID() string {
	return os.Getenv(SECRETID)
}

func SecretKey() string {
	return os.Getenv(SECRETKEY)
}

func Region() string {
	return os.Getenv(REGION)
}

func Password() string {
	return os.Getenv(PASSWORD)
}

func CreateInstancesParam() string {
	return os.Getenv(CREATEINSTANCESPARAM)
}

func NeedDelete() bool {
	return strings.ToLower(os.Getenv(NEEDDELETE)) == "true"
}

func ProviderResImageVersion() string {
	v := os.Getenv(PROVIDERRESVERSION)
	if v == "" {
		v = images.Get().ProviderRes.Tag
	}
	return v
}

func K8sVersion() string {
	v := os.Getenv(K8SVERSION)
	if v == "" {
		v = spec.K8sVersions[0]
	}
	return v
}
