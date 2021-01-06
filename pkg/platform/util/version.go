/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

package util

import (
	"context"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPlatformVersionsFromClusterInfo(ctx context.Context, client kubernetes.Interface) (tkeVersion string, k8sValidVersions []string, err error) {
	k8sValidVersions = []string{}
	clusterInfo, err := client.CoreV1().ConfigMaps("kube-public").Get(ctx, "cluster-info", metav1.GetOptions{})
	if err != nil {
		return
	}
	tkeVersion, ok := clusterInfo.Data["tkeVersion"]
	if !ok {
		return
	}
	err = json.Unmarshal([]byte(clusterInfo.Data["k8sValidVersions"]), &k8sValidVersions)
	return
}
