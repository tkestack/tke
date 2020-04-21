/*
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
	"fmt"
	"strings"

	"tkestack.io/tke/pkg/auth/util"

	"github.com/casbin/casbin/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	v1 "tkestack.io/tke/api/auth/v1"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	"tkestack.io/tke/pkg/util/log"
)

// GroupedResourcesDeleterInterface to delete a group with all resources in
// it.
type GroupedResourcesDeleterInterface interface {
	Delete(groupName string) error
}

// NewGroupedResourcesDeleter to create the groupedResourcesDeleter that
// implement the GroupedResourcesDeleterInterface by given client and
// configure.
func NewGroupedResourcesDeleter(groupClient v1clientset.LocalGroupInterface,
	authClient v1clientset.AuthV1Interface,
	enforcer *casbin.SyncedEnforcer,
	finalizerToken v1.FinalizerName,
	deleteGroupWhenDone bool) GroupedResourcesDeleterInterface {
	d := &groupedResourcesDeleter{
		groupClient:         groupClient,
		authClient:          authClient,
		enforcer:            enforcer,
		finalizerToken:      finalizerToken,
		deleteGroupWhenDone: deleteGroupWhenDone,
	}
	return d
}

var _ GroupedResourcesDeleterInterface = &groupedResourcesDeleter{}

// groupedResourcesDeleter is used to delete all resources in a given group.
type groupedResourcesDeleter struct {
	// Client to manipulate the group.
	groupClient v1clientset.LocalGroupInterface
	authClient  v1clientset.AuthV1Interface

	enforcer *casbin.SyncedEnforcer
	// The finalizer token that should be removed from the group
	// when all resources in that group have been deleted.
	finalizerToken v1.FinalizerName
	// Also delete the group when all resources in the group have been deleted.
	deleteGroupWhenDone bool
}

// Delete deletes all resources in the given group.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   group (does nothing if deletion timestamp is missing).
// * Verifies that the group is in the "terminating" phase
//   (updates the group phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given group.
// * Deletes the group if deleteGroupWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *groupedResourcesDeleter) Delete(groupName string) error {
	// Multiple controllers may edit a group during termination
	// first get the latest state of the group before proceeding
	// if the group was deleted already, don't do anything
	group, err := d.groupClient.Get(groupName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if group.DeletionTimestamp == nil {
		return nil
	}

	log.Infof("group controller - syncGroup - group: %s, finalizerToken: %s", group.Name, d.finalizerToken)

	// ensure that the status is up to date on the group
	// if we get a not found error, we assume the group is truly gone
	group, err = d.retryOnConflictError(group, d.updateGroupStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the group asserts that group is no longer deleting..
	if group.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the group if it is already finalized.
	if d.deleteGroupWhenDone && finalized(group) {
		return d.deleteGroup(group)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(group)
	if err != nil {
		return err
	}

	// we have removed content, so mark it finalized by us
	group, err = d.retryOnConflictError(group, d.finalizeGroup)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do group deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check if we can delete now.
	if d.deleteGroupWhenDone && finalized(group) {
		return d.deleteGroup(group)
	}

	return nil
}

// Deletes the given group.
func (d *groupedResourcesDeleter) deleteGroup(group *v1.LocalGroup) error {
	var opts *metav1.DeleteOptions
	uid := group.UID
	if len(uid) > 0 {
		opts = &metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}

	err := d.groupClient.Delete(group.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		log.Error("error", log.Err(err))
		return err
	}
	return nil
}

// updateGroupFunc is a function that makes an update to a group
type updateGroupFunc func(group *v1.LocalGroup) (*v1.LocalGroup, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in client code
func (d *groupedResourcesDeleter) retryOnConflictError(group *v1.LocalGroup, fn updateGroupFunc) (result *v1.LocalGroup, err error) {
	latestGroup := group
	for {
		result, err = fn(latestGroup)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevGroup := latestGroup
		latestGroup, err = d.groupClient.Get(latestGroup.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevGroup.UID != latestGroup.UID {
			return nil, fmt.Errorf("group uid has changed across retries")
		}
	}
}

// updateGroupStatusFunc will verify that the status of the group is correct
func (d *groupedResourcesDeleter) updateGroupStatusFunc(group *v1.LocalGroup) (*v1.LocalGroup, error) {
	if group.DeletionTimestamp.IsZero() || group.Status.Phase == v1.GroupTerminating {
		return group, nil
	}
	newGroup := v1.LocalGroup{}
	newGroup.ObjectMeta = group.ObjectMeta
	newGroup.Status = group.Status
	newGroup.Status.Phase = v1.GroupTerminating
	return d.groupClient.UpdateStatus(&newGroup)
}

// finalized returns true if the group.Spec.Finalizers is an empty list
func finalized(group *v1.LocalGroup) bool {
	return len(group.Spec.Finalizers) == 0
}

// finalizeGroup removes the specified finalizerToken and finalizes the group
func (d *groupedResourcesDeleter) finalizeGroup(group *v1.LocalGroup) (*v1.LocalGroup, error) {
	groupFinalize := v1.LocalGroup{}
	groupFinalize.ObjectMeta = group.ObjectMeta
	groupFinalize.Spec = group.Spec
	finalizerSet := sets.NewString()
	for i := range group.Spec.Finalizers {
		if group.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(group.Spec.Finalizers[i]))
		}
	}
	groupFinalize.Spec.Finalizers = make([]v1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		groupFinalize.Spec.Finalizers = append(groupFinalize.Spec.Finalizers, v1.FinalizerName(value))
	}

	updated := &v1.LocalGroup{}
	err := d.authClient.RESTClient().Put().
		Resource("localgroups").
		Name(groupFinalize.Name).
		SubResource("finalize").
		Body(&groupFinalize).
		Do().
		Into(updated)

	if err != nil {
		return nil, err
	}
	return updated, err
}

type deleteResourceFunc func(deleter *groupedResourcesDeleter, group *v1.LocalGroup) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteRelatedRoles,
	deleteRelatedRules,
}

// deleteAllContent will use the dynamic client to delete each resource identified in groupVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *groupedResourcesDeleter) deleteAllContent(group *v1.LocalGroup) error {
	log.Debug("LocalGroup controller - deleteAllContent", log.String("groupName", group.Name))

	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(d, group)
		if err != nil {
			// If there is an error, return directly, in case delete roles failed in next try if rule has been deleted.
			log.Error("delete content for group failed", log.String("group", group.Name), log.Err(err))
			return err
		}
	}

	return nil
}

func deleteRelatedRoles(deleter *groupedResourcesDeleter, group *v1.LocalGroup) error {
	log.Debug("LocalGroup controller - deleteRelatedRoles", log.String("group", group.Name))

	subj := util.GroupKey(group.Spec.TenantID, group.Name)
	roles := deleter.enforcer.GetRolesForUserInDomain(subj, util.DefaultDomain)
	log.Info("Try removing related rules for group", log.String("group", group.Name), log.Strings("rules", roles))

	binding := v1.Binding{}
	binding.Groups = append(binding.Groups, v1.Subject{ID: group.Name, Name: group.Spec.DisplayName})
	var errs []error
	for _, role := range roles {
		switch {
		case strings.HasPrefix(role, "pol-"):
			pol := &v1.Policy{}
			err := deleter.authClient.RESTClient().Post().
				Resource("policies").
				Name(role).
				SubResource("unbinding").
				Body(&binding).
				Do().Into(pol)
			if err != nil {
				log.Error("Unbind policy for group failed", log.String("group", group.Name),
					log.String("policy", role), log.Err(err))
				errs = append(errs, err)
			}
		case strings.HasPrefix(role, "rol-"):
			rol := &v1.Role{}
			err := deleter.authClient.RESTClient().Post().
				Resource("roles").
				Name(role).
				SubResource("unbinding").
				Body(&binding).
				Do().Into(rol)
			if err != nil {
				log.Error("Unbind role for group failed", log.String("group", group.Name),
					log.String("policy", role), log.Err(err))
				errs = append(errs, err)
			}
		default:
			log.Error("Unknown role name for group, remove it", log.String("group", group.Name), log.String("role", role))
			_, err := deleter.enforcer.DeleteRoleForUser(subj, role)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func deleteRelatedRules(deleter *groupedResourcesDeleter, group *v1.LocalGroup) error {
	log.Info("LocalGroup controller - deleteRelatedRules", log.String("groupName", group.Name))
	_, err := deleter.enforcer.DeleteRole(util.GroupKey(group.Spec.TenantID, group.Name))
	return err
}
