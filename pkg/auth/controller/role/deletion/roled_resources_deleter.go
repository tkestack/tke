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

package deletion

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	v1 "tkestack.io/tke/api/auth/v1"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	authutil "tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
)

// RoledResourcesDeleterInterface to delete a role with all resources in
// it.
type RoledResourcesDeleterInterface interface {
	Delete(roleName string) error
}

// NewRoleedResourcesDeleter to create the roledResourcesDeleter that
// implement the RoledResourcesDeleterInterface by given client and
// configure.
func NewRoleedResourcesDeleter(pilicyClient v1clientset.RoleInterface,
	authClient v1clientset.AuthV1Interface,
	enforcer *casbin.SyncedEnforcer,
	finalizerToken v1.FinalizerName,
	deleteRoleWhenDone bool) RoledResourcesDeleterInterface {
	d := &roledResourcesDeleter{
		roleClient:         pilicyClient,
		authClient:         authClient,
		enforcer:           enforcer,
		finalizerToken:     finalizerToken,
		deleteRoleWhenDone: deleteRoleWhenDone,
	}
	return d
}

var _ RoledResourcesDeleterInterface = &roledResourcesDeleter{}

// roledResourcesDeleter is used to delete all resources in a given role.
type roledResourcesDeleter struct {
	// Client to manipulate the role.
	roleClient v1clientset.RoleInterface
	authClient v1clientset.AuthV1Interface

	enforcer *casbin.SyncedEnforcer
	// The finalizer token that should be removed from the role
	// when all resources in that role have been deleted.
	finalizerToken v1.FinalizerName
	// Also delete the role when all resources in the role have been deleted.
	deleteRoleWhenDone bool
}

// Delete deletes all resources in the given role.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   role (does nothing if deletion timestamp is missing).
// * Verifies that the role is in the "terminating" phase
//   (updates the role phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given role.
// * Deletes the role if deleteRoleWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *roledResourcesDeleter) Delete(roleName string) error {
	// Multiple controllers may edit a role during termination
	// first get the latest state of the role before proceeding
	// if the role was deleted already, don't do anything
	role, err := d.roleClient.Get(context.Background(), roleName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if role.DeletionTimestamp == nil {
		return nil
	}

	log.Infof("role controller - syncRole - role: %s, finalizerToken: %s", role.Name, d.finalizerToken)

	// ensure that the status is up to date on the role
	// if we get a not found error, we assume the role is truly gone
	role, err = d.retryOnConflictError(role, d.updateRoleStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the role asserts that role is no longer deleting..
	if role.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the role if it is already finalized.
	if d.deleteRoleWhenDone && finalized(role) {
		return d.deleteRole(role)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(role)
	if err != nil {
		return err
	}

	// we have removed content, so mark it finalized by us
	role, err = d.retryOnConflictError(role, d.finalizeRole)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do role deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check if we can delete now.
	if d.deleteRoleWhenDone && finalized(role) {
		return d.deleteRole(role)
	}

	return nil
}

// Deletes the given role.
func (d *roledResourcesDeleter) deleteRole(role *v1.Role) error {
	var opts metav1.DeleteOptions
	uid := role.UID
	if len(uid) > 0 {
		opts = metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.roleClient.Delete(context.Background(), role.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		log.Error("error", log.Err(err))
		return err
	}
	return nil
}

// updateRoleFunc is a function that makes an update to a role
type updateRoleFunc func(role *v1.Role) (*v1.Role, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in client code
func (d *roledResourcesDeleter) retryOnConflictError(role *v1.Role, fn updateRoleFunc) (result *v1.Role, err error) {
	latestRole := role
	for {
		result, err = fn(latestRole)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevRole := latestRole
		latestRole, err = d.roleClient.Get(context.Background(), latestRole.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevRole.UID != latestRole.UID {
			return nil, fmt.Errorf("role uid has changed across retries")
		}
	}
}

// updateRoleStatusFunc will verify that the status of the role is correct
func (d *roledResourcesDeleter) updateRoleStatusFunc(role *v1.Role) (*v1.Role, error) {
	if role.DeletionTimestamp.IsZero() || role.Status.Phase == v1.RoleTerminating {
		return role, nil
	}
	newRole := v1.Role{}
	newRole.ObjectMeta = role.ObjectMeta
	newRole.Status = role.Status
	newRole.Status.Phase = v1.RoleTerminating
	return d.roleClient.UpdateStatus(context.Background(), &newRole, metav1.UpdateOptions{})
}

// finalized returns true if the role.Spec.Finalizers is an empty list
func finalized(role *v1.Role) bool {
	return len(role.Spec.Finalizers) == 0
}

// finalizeRole removes the specified finalizerToken and finalizes the role
func (d *roledResourcesDeleter) finalizeRole(role *v1.Role) (*v1.Role, error) {
	roleFinalize := v1.Role{}
	roleFinalize.ObjectMeta = role.ObjectMeta
	roleFinalize.Spec = role.Spec
	finalizerSet := sets.NewString()
	for i := range role.Spec.Finalizers {
		if role.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(role.Spec.Finalizers[i]))
		}
	}
	roleFinalize.Spec.Finalizers = make([]v1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		roleFinalize.Spec.Finalizers = append(roleFinalize.Spec.Finalizers, v1.FinalizerName(value))
	}

	updated := &v1.Role{}
	err := d.authClient.RESTClient().Put().
		Resource("roles").
		Name(roleFinalize.Name).
		SubResource("finalize").
		Body(&roleFinalize).
		Do(context.Background()).
		Into(updated)

	if err != nil {
		return nil, err
	}
	return updated, err
}

type deleteResourceFunc func(deleter *roledResourcesDeleter, role *v1.Role) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteRelatedRules,
}

// deleteAllContent will use the dynamic client to delete each resource identified in roleVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *roledResourcesDeleter) deleteAllContent(role *v1.Role) error {
	log.Debug("Role controller - deleteAllContent", log.String("roleName", role.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(d, role)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	return nil
}

func deleteRelatedRules(deleter *roledResourcesDeleter, role *v1.Role) error {
	log.Info("Role controller - deleteRelatedRules", log.String("roleName", role.Name))
	projectID := authutil.DefaultDomain
	if role.Spec.ProjectID != "" {
		projectID = role.Spec.ProjectID
	}

	users := deleter.enforcer.GetUsersForRoleInDomain(role.Name, projectID)
	log.Info("Try removing related rules for role", log.String("role", role.Name), log.Strings("rules", users))
	_, err := deleter.enforcer.DeleteRole(role.Name)
	return err
}
