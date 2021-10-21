/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
	"testing"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
)

func TestController_needsUpdate(t *testing.T) {
	// type fields struct {
	// 	queue             workqueue.RateLimitingInterface
	// 	lister            platformv1lister.ClusterLister
	// 	listerSynced      cache.InformerSynced
	// 	log               log.Logger
	// 	platformClient    platformversionedclient.PlatformV1Interface
	// 	deleter           deletion.ClusterDeleterInterface
	// 	healthCheckPeriod time.Duration
	// }
	type args struct {
		old *platformv1.Cluster
		new *platformv1.Cluster
	}
	tests := []struct {
		name string
		// fields fields
		args args
		want bool
	}{
		{
			name: "change spec",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Spec:       platformv1.ClusterSpec{DisplayName: "old"},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Spec:       platformv1.ClusterSpec{DisplayName: "nes"},
				},
			},
			want: true,
		},
		{
			name: "Initializing to Running",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterRunning},
				},
			},
			want: true,
		},
		{
			name: "Initializing to Failed",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterFailed},
				},
			},
			want: true,
		},
		{
			name: "Running to Failed",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterRunning},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterFailed},
				},
			},
			want: true,
		},
		{
			name: "Running to Terminating",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterRunning},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterTerminating},
				},
			},
			want: true,
		},
		{
			name: "Failed to Terminating",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterFailed},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterTerminating},
				},
			},
			want: true,
		},
		{
			name: "Failed to Running",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterFailed},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterRunning},
				},
			},
			want: true,
		},
		{
			name: "Failed to Initializing",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterFailed},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing},
				},
			},
			want: true,
		},
		{
			name: "last conditon unkonwn to false",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status: platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing,
						Conditions: []platformv1.ClusterCondition{{Status: platformv1.ConditionUnknown}}},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status: platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing,
						Conditions: []platformv1.ClusterCondition{{Status: platformv1.ConditionFalse}}},
				},
			},
			want: false,
		},
		{
			name: "last conditon unkonwn to true",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status: platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing,
						Conditions: []platformv1.ClusterCondition{{Status: platformv1.ConditionUnknown}}},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status: platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing,
						Conditions: []platformv1.ClusterCondition{{Status: platformv1.ConditionTrue}}},
				},
			},
			want: true,
		},
		{
			name: "last conditon unkonwn to true resync",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status: platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing,
						Conditions: []platformv1.ClusterCondition{{Status: platformv1.ConditionUnknown}}},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status: platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing,
						Conditions: []platformv1.ClusterCondition{{Status: platformv1.ConditionUnknown}}},
				},
			},
			want: true,
		},
		{
			name: "last conditon true to unknown",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status: platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing,
						Conditions: []platformv1.ClusterCondition{{Status: platformv1.ConditionTrue}}},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status: platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing,
						Conditions: []platformv1.ClusterCondition{{Status: platformv1.ConditionUnknown}}},
				},
			},
			want: true,
		},
		{
			name: "last conditon false to unknown",
			args: args{
				old: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status: platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing,
						Conditions: []platformv1.ClusterCondition{{Status: platformv1.ConditionFalse}}},
				},
				new: &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status: platformv1.ClusterStatus{Phase: platformv1.ClusterInitializing,
						Conditions: []platformv1.ClusterCondition{{Status: platformv1.ConditionUnknown}}},
				},
			},
			want: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				// queue:             tt.fields.queue,
				// lister:            tt.fields.lister,
				// listerSynced:      tt.fields.listerSynced,
				// log:               tt.fields.log,
				// platformClient:    tt.fields.platformClient,
				// deleter:           tt.fields.deleter,
				healthCheckPeriod: time.Second,
			}
			if got := c.needsUpdate(tt.args.old, tt.args.new); got != tt.want {
				t.Errorf("Controller.needsUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
