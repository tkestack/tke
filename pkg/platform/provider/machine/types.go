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

package machine

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

	gob.Register(metav1.TypeMeta{})
	gob.Register(metav1.ObjectMeta{})
	gob.Register(metav1.OwnerReference{})
	gob.Register(metav1.Time{})
	gob.Register(metav1.Status{})
	gob.Register(metav1.Timestamp{})
	gob.Register(metav1.Duration{})
	gob.Register(field.ErrorList{})

	gob.Register(platform.Machine{})
	gob.Register(platform.MachineSpec{})
	gob.Register(platform.MachineStatus{})
	gob.Register(platform.MachineCondition{})
	gob.Register(platform.MachineAddress{})

	gob.Register(v1.Machine{})
	gob.Register(v1.MachineSpec{})
	gob.Register(v1.MachineStatus{})
	gob.Register(v1.MachineCondition{})
	gob.Register(v1.MachineAddress{})

	gob.Register(v1.Cluster{})
	gob.Register(v1.ClusterCredential{})
}
