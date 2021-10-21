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

package machine

import (
	"testing"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
)

func TestController_needsUpdate(t *testing.T) {
	type fields struct {
		// queue          workqueue.RateLimitingInterface
		// lister         platformv1lister.MachineLister
		// listerSynced   cache.InformerSynced
		// log            log.Logger
		// platformClient platformversionedclient.PlatformV1Interface
		// deleter        deletion.MachineDeleterInterface
	}
	type args struct {
		old *platformv1.Machine
		new *platformv1.Machine
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "change spec",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Spec:       platformv1.MachineSpec{Type: "old"},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Spec:       platformv1.MachineSpec{Type: "nes"},
				},
			},
			want: true,
		},
		{
			name: "Initializing to Running",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineInitializing},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineInitializing},
				},
			},
			want: true,
		},
		{
			name: "Initializing to Failed",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineInitializing},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineFailed},
				},
			},
			want: true,
		},
		{
			name: "Running to Failed",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineRunning},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineFailed},
				},
			},
			want: true,
		},
		{
			name: "Running to Terminating",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineRunning},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineTerminating},
				},
			},
			want: true,
		},
		{
			name: "Failed to Initializing",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineFailed},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineInitializing},
				},
			},
			want: true,
		},
		{
			name: "Failed to Running",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineFailed},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineRunning},
				},
			},
			want: true,
		},
		{
			name: "Failed to Running",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineFailed},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status:     platformv1.MachineStatus{Phase: platformv1.MachineRunning},
				},
			},
			want: true,
		},
		{
			name: "last conditon unkonwn to false",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status: platformv1.MachineStatus{Phase: platformv1.MachineInitializing,
						Conditions: []platformv1.MachineCondition{{Status: platformv1.ConditionUnknown}}},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status: platformv1.MachineStatus{Phase: platformv1.MachineInitializing,
						Conditions: []platformv1.MachineCondition{{Status: platformv1.ConditionFalse}}},
				},
			},
			want: false,
		},
		{
			name: "last conditon unkonwn to true",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status: platformv1.MachineStatus{Phase: platformv1.MachineInitializing,
						Conditions: []platformv1.MachineCondition{{Status: platformv1.ConditionUnknown}}},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status: platformv1.MachineStatus{Phase: platformv1.MachineInitializing,
						Conditions: []platformv1.MachineCondition{{Status: platformv1.ConditionTrue}}},
				},
			},
			want: true,
		},
		{
			name: "last conditon unkonwn to true resync",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status: platformv1.MachineStatus{Phase: platformv1.MachineInitializing,
						Conditions: []platformv1.MachineCondition{{Status: platformv1.ConditionUnknown}}},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status: platformv1.MachineStatus{Phase: platformv1.MachineInitializing,
						Conditions: []platformv1.MachineCondition{{Status: platformv1.ConditionUnknown}}},
				},
			},
			want: true,
		},
		{
			name: "last conditon true to unknown",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status: platformv1.MachineStatus{Phase: platformv1.MachineInitializing,
						Conditions: []platformv1.MachineCondition{{Status: platformv1.ConditionTrue}}},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status: platformv1.MachineStatus{Phase: platformv1.MachineInitializing,
						Conditions: []platformv1.MachineCondition{{Status: platformv1.ConditionUnknown}}},
				},
			},
			want: true,
		},
		{
			name: "last conditon false to unknown",
			args: args{
				old: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "old"},
					Status: platformv1.MachineStatus{Phase: platformv1.MachineInitializing,
						Conditions: []platformv1.MachineCondition{{Status: platformv1.ConditionFalse}}},
				},
				new: &platformv1.Machine{
					ObjectMeta: v1.ObjectMeta{ResourceVersion: "new"},
					Status: platformv1.MachineStatus{Phase: platformv1.MachineInitializing,
						Conditions: []platformv1.MachineCondition{{Status: platformv1.ConditionUnknown}}},
				},
			},
			want: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				// queue:          tt.fields.queue,
				// lister:         tt.fields.lister,
				// listerSynced:   tt.fields.listerSynced,
				// log:            tt.fields.log,
				// platformClient: tt.fields.platformClient,
				// deleter:        tt.fields.deleter,
			}
			if got := c.needsUpdate(tt.args.old, tt.args.new); got != tt.want {
				t.Errorf("Controller.needsUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
