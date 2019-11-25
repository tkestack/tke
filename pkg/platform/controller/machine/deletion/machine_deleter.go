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

package deletion

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
)

// MachineDeleterInterface to delete a machine with all resources in it.
type MachineDeleterInterface interface {
	Delete(name string) error
}

// NewMachineDeleter creates the machineeleter object and returns it.
func NewMachineDeleter(machineClient v1clientset.MachineInterface,
	platformClient v1clientset.PlatformV1Interface,
	finalizerToken v1.FinalizerName,
	deleteWhenDone bool) MachineDeleterInterface {
	d := &machineDeleter{
		machineClient:  machineClient,
		platformClient: platformClient,
		deleteWhenDone: deleteWhenDone,
		finalizerToken: finalizerToken,
	}
	return d
}

var _ MachineDeleterInterface = &machineDeleter{}

// machineDeleter is used to delete all resources in a given machine.
type machineDeleter struct {
	// Client to manipulate the machine.
	machineClient  v1clientset.MachineInterface
	platformClient v1clientset.PlatformV1Interface
	// The finalizer token that should be removed from the machine
	// when all resources in that machine have been deleted.
	finalizerToken v1.FinalizerName
	// Also delete the machine when all resources in the machine have been deleted.
	deleteWhenDone bool
}

// Delete deletes all resources in the given machine.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   machine (does nothing if deletion timestamp is missing).
// * Verifies that the machine is in the "terminating" phase
//   (updates the machine phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given machine.
// * Deletes the machine if deleteWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *machineDeleter) Delete(name string) error {
	// Multiple controllers may edit a machine during termination
	// first get the latest state of the machine before proceeding
	// if the machine was deleted already, don't do anything
	machine, err := d.machineClient.Get(name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if machine.DeletionTimestamp == nil {
		return nil
	}

	log.Info("Machine controller - machine deleter", log.String("name", machine.Name), log.String("finalizerToken", string(d.finalizerToken)))

	// ensure that the status is up to date on the machine
	// if we get a not found error, we assume the machine is truly gone
	machine, err = d.retryOnConflictError(machine, d.updateMachineStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the machine asserts that machine is no longer deleting..
	if machine.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the machine if it is already finalized.
	if d.deleteWhenDone && finalized(machine) {
		return d.deleteMachine(machine)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(machine)
	if err != nil {
		return err
	}

	// we have removed content, so mark it finalized by us
	machine, err = d.retryOnConflictError(machine, d.finalizeMachine)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do machine deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check if we can delete now.
	if d.deleteWhenDone && finalized(machine) {
		return d.deleteMachine(machine)
	}
	return nil
}

// Deletes the given machine.
func (d *machineDeleter) deleteMachine(machine *v1.Machine) error {
	var opts *metav1.DeleteOptions
	uid := machine.UID
	if len(uid) > 0 {
		opts = &metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.machineClient.Delete(machine.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// updateMachineFunc is a function that makes an update to a namespace
type updateMachineFunc func(machine *v1.Machine) (*v1.Machine, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in client code
func (d *machineDeleter) retryOnConflictError(machine *v1.Machine, fn updateMachineFunc) (result *v1.Machine, err error) {
	latestMachine := machine
	for {
		result, err = fn(latestMachine)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevMachine := latestMachine
		latestMachine, err = d.machineClient.Get(latestMachine.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevMachine.UID != latestMachine.UID {
			return nil, fmt.Errorf("machine uid has changed across retries")
		}
	}
}

// updateMachineStatusFunc will verify that the status of the machine is correct
func (d *machineDeleter) updateMachineStatusFunc(machine *v1.Machine) (*v1.Machine, error) {
	if machine.DeletionTimestamp.IsZero() || machine.Status.Phase == v1.MachineTerminating {
		return machine, nil
	}
	newMachine := v1.Machine{}
	newMachine.ObjectMeta = machine.ObjectMeta
	newMachine.Status = machine.Status
	newMachine.Status.Phase = v1.MachineTerminating
	return d.machineClient.UpdateStatus(&newMachine)
}

// finalized returns true if the machine.Spec.Finalizers is an empty list
func finalized(machine *v1.Machine) bool {
	return len(machine.Spec.Finalizers) == 0
}

// finalizeMachine removes the specified finalizerToken and finalizes the machine
func (d *machineDeleter) finalizeMachine(machine *v1.Machine) (*v1.Machine, error) {
	machineFinalize := v1.Machine{}
	machineFinalize.ObjectMeta = machine.ObjectMeta
	machineFinalize.Spec = machine.Spec

	finalizerSet := sets.NewString()
	for i := range machine.Spec.Finalizers {
		if machine.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(machine.Spec.Finalizers[i]))
		}
	}
	machineFinalize.Spec.Finalizers = make([]v1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		machineFinalize.Spec.Finalizers = append(machineFinalize.Spec.Finalizers, v1.FinalizerName(value))
	}

	machine = &v1.Machine{}
	err := d.platformClient.RESTClient().Put().
		Resource("machines").
		Name(machineFinalize.Name).
		SubResource("finalize").
		Body(&machineFinalize).
		Do().
		Into(machine)

	if err != nil {
		// it was removed already, so life is good
		if errors.IsNotFound(err) {
			return machine, nil
		}
	}
	return machine, err
}

type deleteResourceFunc func(deleter *machineDeleter, machine *v1.Machine) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteNode,
}

// deleteAllContent will use the client to delete each resource identified in machine.
func (d *machineDeleter) deleteAllContent(machine *v1.Machine) error {
	log.Debug("Machine controller - deleteAllContent", log.String("machineName", machine.ObjectMeta.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(d, machine)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.Debug("Machine controller - deletedAllContent", log.String("machineName", machine.ObjectMeta.Name))
	return nil
}

func deleteNode(deleter *machineDeleter, machine *v1.Machine) error {
	log.Debug("Machine controller - deleteNode", log.String("machineName", machine.ObjectMeta.Name))

	clientset, err := util.BuildExternalClientSetWithName(deleter.platformClient, machine.Spec.ClusterName)
	if err != nil {
		return err
	}

	err = clientset.CoreV1().Nodes().Delete(machine.Spec.IP, &metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
