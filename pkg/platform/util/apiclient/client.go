/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
	platformv1client "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
)

const (
	// APICallRetryInterval defines how long should wait before retrying a failed API operation
	APICallRetryInterval = 500 * time.Millisecond
	// PatchNodeTimeout specifies how long should wait for applying the label and taint on the master before timing out
	PatchNodeTimeout = 2 * time.Minute
)

// PatchMachine tries to patch a machine using patchFn for the actual mutating logic.
// Retries are provided by the wait package.
func PatchMachine(ctx context.Context, client platformv1client.PlatformV1Interface, machineName string, patchFn func(*platformv1.Machine)) error {
	// wait.Poll will rerun the condition function every interval function if
	// the function returns false. If the condition function returns an error
	// then the retries end and the error is returned.
	return wait.Poll(APICallRetryInterval, PatchNodeTimeout, PatchMachineOnce(ctx, client, machineName, patchFn))
}

// PatchMachineOnce executes patchFn on the machine object found by the machine name.
// This is a condition function meant to be used with wait.Poll. false, nil
// implies it is safe to try again, an error indicates no more tries should be
// made and true indicates success.
func PatchMachineOnce(ctx context.Context, client platformv1client.PlatformV1Interface, machineName string, patchFn func(*platformv1.Machine)) func() (bool, error) {
	return func() (bool, error) {
		// First get the machine object
		machine, err := client.Machines().Get(ctx, machineName, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		oldData, err := json.Marshal(machine)
		if err != nil {
			return false, errors.Wrapf(err, "failed to marshal unmodified machine %q into JSON", machine.Name)
		}

		// Execute the mutating function
		patchFn(machine)

		newData, err := json.Marshal(machine)
		if err != nil {
			return false, errors.Wrapf(err, "failed to marshal modified machine %q into JSON", machine.Name)
		}

		patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, corev1.Node{})
		if err != nil {
			return false, errors.Wrap(err, "failed to create two way merge patch")
		}

		if _, err := client.Machines().Patch(ctx, machine.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{}); err != nil {
			if apierrors.IsConflict(err) {
				return false, nil
			}
			return false, errors.Wrapf(err, "error patching machine %q through apiserver", machine.Name)
		}

		return true, nil
	}
}
