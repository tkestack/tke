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

package cluster

import (
	"context"
	"encoding/gob"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/api/platform/v1"
)

func init() {
	gob.Register(context.DeadlineExceeded)

	gob.Register(UserInfo{})
	gob.Register(Cluster{})

	gob.Register(metav1.TypeMeta{})
	gob.Register(metav1.ObjectMeta{})
	gob.Register(metav1.OwnerReference{})
	gob.Register(metav1.Time{})
	gob.Register(metav1.Status{})
	gob.Register(metav1.Timestamp{})
	gob.Register(metav1.Duration{})
	gob.Register(field.ErrorList{})

	gob.Register(platform.Cluster{})
	gob.Register(platform.ClusterSpec{})
	gob.Register(platform.ClusterStatus{})
	gob.Register(platform.ClusterFeature{})
	gob.Register(platform.ClusterProperty{})
	gob.Register(platform.ClusterCondition{})
	gob.Register(platform.ClusterAddress{})
	gob.Register(platform.ClusterAddon{})
	gob.Register(platform.ClusterCredential{})
	gob.Register(platform.ClusterComponentReplicas{})
	gob.Register(platform.ClusterComponent{})

	gob.Register(v1.Cluster{})
	gob.Register(v1.ClusterSpec{})
	gob.Register(v1.ClusterStatus{})
	gob.Register(v1.ClusterFeature{})
	gob.Register(v1.ClusterProperty{})
	gob.Register(v1.ClusterCondition{})
	gob.Register(v1.ClusterAddress{})
	gob.Register(v1.ClusterAddon{})
	gob.Register(v1.ClusterCredential{})
	gob.Register(v1.ClusterComponentReplicas{})
	gob.Register(v1.ClusterComponent{})
}
