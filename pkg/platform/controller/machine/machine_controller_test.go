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
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	machineprovider "tkestack.io/tke/pkg/platform/provider/machine"
)

func newMachineForTest(resourcesVersion string, spec *platformv1.MachineSpec, phase platformv1.MachinePhase, conditions []platformv1.MachineCondition) *platformv1.Machine {
	mc := &platformv1.Machine{
		ObjectMeta: v1.ObjectMeta{ResourceVersion: resourcesVersion},
		Spec: platformv1.MachineSpec{
			TenantID:    "default",
			ClusterName: "global",
			Type:        "Baremetal",
			IP:          "127.0.0.1",
			Port:        22,
			Username:    "root",
		},
		Status: platformv1.MachineStatus{
			Phase: platformv1.MachineRunning,
			Conditions: []platformv1.MachineCondition{
				{
					Type:          machineprovider.ConditionTypeHealthCheck,
					Status:        platformv1.ConditionTrue,
					LastProbeTime: v1.Now(),
				},
			},
		},
	}
	if spec != nil {
		mc.Spec = *spec
	}
	if len(phase) != 0 {
		mc.Status.Phase = phase
	}
	if conditions != nil {
		mc.Status.Conditions = conditions
	}
	return mc
}

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
				old: newMachineForTest("old", &platformv1.MachineSpec{IP: "127.0.0.1"}, platformv1.MachinePhase(""), nil),
				new: newMachineForTest("new", &platformv1.MachineSpec{IP: "localhost"}, platformv1.MachinePhase(""), nil),
			},
			want: true,
		},
		{
			name: "Initializing to Running",
			args: args{
				old: newMachineForTest("old", nil, platformv1.MachineInitializing, nil),
				new: newMachineForTest("new", nil, platformv1.MachineRunning, nil),
			},
			want: true,
		},
		{
			name: "Initializing to Failed",
			args: args{
				old: newMachineForTest("old", nil, platformv1.MachineInitializing, nil),
				new: newMachineForTest("new", nil, platformv1.MachineFailed, nil),
			},
			want: true,
		},
		{
			name: "Running to Failed",
			args: args{
				old: newMachineForTest("old", nil, platformv1.MachineRunning, nil),
				new: newMachineForTest("new", nil, platformv1.MachineFailed, nil),
			},
			want: true,
		},
		{
			name: "Running to Terminating",
			args: args{
				old: newMachineForTest("old", nil, platformv1.MachineRunning, nil),
				new: newMachineForTest("new", nil, platformv1.MachineTerminating, nil),
			},
			want: true,
		},
		{
			name: "Failed to Initializing",
			args: args{
				old: newMachineForTest("old", nil, platformv1.MachineFailed, nil),
				new: newMachineForTest("new", nil, platformv1.MachineInitializing, nil),
			},
			want: true,
		},
		{
			name: "Failed to Running",
			args: args{
				old: newMachineForTest("old", nil, platformv1.MachineFailed, nil),
				new: newMachineForTest("new", nil, platformv1.MachineRunning, nil),
			},
			want: true,
		},
		{
			name: "Initializing last conditon unkonwn to false",
			args: args{
				old: newMachineForTest("old", nil, platformv1.MachineInitializing, []platformv1.MachineCondition{{Status: platformv1.ConditionUnknown}}),
				new: newMachineForTest("new", nil, platformv1.MachineInitializing, []platformv1.MachineCondition{{Status: platformv1.ConditionFalse}}),
			},
			want: false,
		},
		{
			name: "last conditon unkonwn to true",
			args: args{
				old: newMachineForTest("old", nil, platformv1.MachinePhase(""), []platformv1.MachineCondition{{Status: platformv1.ConditionUnknown}}),
				new: newMachineForTest("new", nil, platformv1.MachinePhase(""), []platformv1.MachineCondition{{Status: platformv1.ConditionFalse}}),
			},
			want: true,
		},
		{
			name: "Initializing last conditon false retrun true if resync",
			args: func() args {
				// resource version equal
				new := newMachineForTest("new", nil, platformv1.MachineInitializing, []platformv1.MachineCondition{{Status: platformv1.ConditionFalse}})
				return args{new, new}
			}(),
			want: true,
		},
		{
			name: "last conditon true to unknown",
			args: args{
				old: newMachineForTest("old", nil, platformv1.MachinePhase(""), []platformv1.MachineCondition{{Status: platformv1.ConditionTrue}}),
				new: newMachineForTest("new", nil, platformv1.MachinePhase(""), []platformv1.MachineCondition{{Status: platformv1.ConditionUnknown}}),
			},
			want: true,
		},
		{
			name: "last conditon false to unknown",
			args: args{
				old: newMachineForTest("old", nil, platformv1.MachinePhase(""), []platformv1.MachineCondition{{Status: platformv1.ConditionFalse}}),
				new: newMachineForTest("new", nil, platformv1.MachinePhase(""), []platformv1.MachineCondition{{Status: platformv1.ConditionUnknown}}),
			},
			want: true,
		},
		{
			name: "health check is not long enough",
			args: func() args {
				new := newMachineForTest("new", nil, platformv1.MachinePhase(""), []platformv1.MachineCondition{{
					Type:          machineprovider.ConditionTypeHealthCheck,
					Status:        platformv1.ConditionTrue,
					LastProbeTime: v1.NewTime(time.Now().Add(-resyncInternal / 2))}})
				return args{new, new}
			}(),
			want: false,
		},
		{
			name: "health check is long enough",
			args: func() args {
				new := newMachineForTest("new", nil, platformv1.MachinePhase(""), []platformv1.MachineCondition{{
					Type:          machineprovider.ConditionTypeHealthCheck,
					Status:        platformv1.ConditionTrue,
					LastProbeTime: v1.NewTime(time.Now().Add(-resyncInternal - 1))}})
				return args{new, new}
			}(),
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
