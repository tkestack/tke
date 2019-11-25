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

package resourcelock

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	businessv1 "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	monitorv1 "tkestack.io/tke/api/client/clientset/versioned/typed/monitor/v1"
	notifyv1 "tkestack.io/tke/api/client/clientset/versioned/typed/notify/v1"
	platformv1 "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
)

const (
	// LeaderElectionRecordAnnotationKey is key of leader election record annotation.
	LeaderElectionRecordAnnotationKey = "tke/leader"
)

// LeaderElectionRecord is the record that is stored in the leader election annotation.
// This information should be used for observational purposes only and could be replaced
// with a random string (e.g. UUID) with only slight modification of this code.
type LeaderElectionRecord struct {
	HolderIdentity       string      `json:"holderIdentity"`
	LeaseDurationSeconds int         `json:"leaseDurationSeconds"`
	AcquireTime          metav1.Time `json:"acquireTime"`
	RenewTime            metav1.Time `json:"renewTime"`
	LeaderTransitions    int         `json:"leaderTransitions"`
}

// Config common data that exists across different resource locks
type Config struct {
	Identity string
}

// NewPlatform will create a lock of a given type according to the input parameters
func NewPlatform(name string, client platformv1.PlatformV1Interface, rlc Config) Interface {
	return &PlatformConfigMapLock{
		ConfigMapMeta: metav1.ObjectMeta{
			Name: name,
		},
		Client:     client,
		LockConfig: rlc,
	}
}

// NewBusiness will create a lock of a given type according to the input parameters
func NewBusiness(name string, client businessv1.BusinessV1Interface, rlc Config) Interface {
	return &BusinessConfigMapLock{
		ConfigMapMeta: metav1.ObjectMeta{
			Name: name,
		},
		Client:     client,
		LockConfig: rlc,
	}
}

// NewNotify will create a lock of a given type according to the input parameters
func NewNotify(name string, client notifyv1.NotifyV1Interface, rlc Config) Interface {
	return &NotifyConfigMapLock{
		ConfigMapMeta: metav1.ObjectMeta{
			Name: name,
		},
		Client:     client,
		LockConfig: rlc,
	}
}

// NewMonitor will create a lock of a given type according to the input parameters
func NewMonitor(name string, client monitorv1.MonitorV1Interface, rlc Config) Interface {
	return &MonitorConfigMapLock{
		ConfigMapMeta: metav1.ObjectMeta{
			Name: name,
		},
		Client:     client,
		LockConfig: rlc,
	}
}
