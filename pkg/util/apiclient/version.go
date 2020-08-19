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

package apiclient

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"k8s.io/client-go/kubernetes"
)

func ClusterVersionIsBefore19(client kubernetes.Interface) bool {
	result, err := CheckClusterVersion(client, "< 1.9")
	if err != nil {
		return false
	}

	return result
}

func ClusterVersionIsBefore116(client kubernetes.Interface) bool {
	result, err := CheckClusterVersion(client, "< 1.16")
	if err != nil {
		return false
	}

	return result
}

func CheckClusterVersion(client kubernetes.Interface, versionConstraint string) (bool, error) {
	version, err := GetClusterVersion(client)
	if err != nil {
		return false, err
	}

	return CheckVersion(version, versionConstraint)
}

func GetClusterVersion(client kubernetes.Interface) (string, error) {
	version, err := client.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v.%v", version.Major, strings.TrimSuffix(version.Minor, "+")), nil
}

func CheckVersion(version string, versionConstraint string) (bool, error) {
	c, err := semver.NewConstraint(versionConstraint)
	if err != nil {
		return false, err
	}
	v, err := semver.NewVersion(version)
	if err != nil {
		return false, err
	}

	return c.Check(v), nil
}

func CheckVersionOrDie(version string, versionConstraint string) bool {
	ok, err := CheckVersion(version, versionConstraint)
	if err != nil {
		panic(err)
	}

	return ok
}
