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

// LocalIdentitiedResourcesDeleterInterface to delete a localIdentity with all resources in
// it.
type LocalIdentitiedResourcesDeleterInterface interface {
	Delete(localIdentityName string) error
}

// NewLocalIdentitiedResourcesDeleterx to create the loalIdentitiedResourcesDeleter that
// implement the LocalIdentitiedResourcesDeleterInterface by given client and
// configure.
func NewLocalIdentitiedResourcesDeleter(localIdentityClient v1clientset.LocalIdentityInterface,
	authClient v1clientset.AuthV1Interface,
	enforcer *casbin.SyncedEnforcer,
	finalizerToken v1.FinalizerName,
	deleteLocalIdentityWhenDone bool) LocalIdentitiedResourcesDeleterInterface {
	d := &loalIdentitiedResourcesDeleter{
		localIdentityClient:         localIdentityClient,
		authClient:                  authClient,
		enforcer:                    enforcer,
		finalizerToken:              finalizerToken,
		deleteLocalIdentityWhenDone: deleteLocalIdentityWhenDone,
	}
	return d
}

var _ LocalIdentitiedResourcesDeleterInterface = &loalIdentitiedResourcesDeleter{}

// loalIdentitiedResourcesDeleter is used to delete all resources in a given localIdentity.
type loalIdentitiedResourcesDeleter struct {
	// Client to manipulate the localIdentity.
	localIdentityClient v1clientset.LocalIdentityInterface
	authClient          v1clientset.AuthV1Interface

	enforcer *casbin.SyncedEnforcer
	// The finalizer token that should be removed from the localIdentity
	// when all resources in that localIdentity have been deleted.
	finalizerToken v1.FinalizerName
	// Also delete the localIdentity when all resources in the localIdentity have been deleted.
	deleteLocalIdentityWhenDone bool
}

// Delete deletes all resources in the given localIdentity.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   localIdentity (does nothing if deletion timestamp is missing).
// * Verifies that the localIdentity is in the "terminating" phase
//   (updates the localIdentity phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given localIdentity.
// * Deletes the localIdentity if deleteLocalIdentityWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *loalIdentitiedResourcesDeleter) Delete(localIdentityName string) error {
	// Multiple controllers may edit a localIdentity during termination
	// first get the latest state of the localIdentity before proceeding
	// if the localIdentity was deleted already, don't do anything
	localIdentity, err := d.localIdentityClient.Get(localIdentityName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if localIdentity.DeletionTimestamp == nil {
		return nil
	}

	log.Infof("localIdentity controller - syncLocalIdentity - localIdentity: %s, finalizerToken: %s", localIdentity.Name, d.finalizerToken)

	// ensure that the status is up to date on the localIdentity
	// if we get a not found error, we assume the localIdentity is truly gone
	localIdentity, err = d.retryOnConflictError(localIdentity, d.updateLocalIdentityStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the localIdentity asserts that localIdentity is no longer deleting..
	if localIdentity.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the localIdentity if it is already finalized.
	if d.deleteLocalIdentityWhenDone && finalized(localIdentity) {
		return d.deleteLocalIdentity(localIdentity)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(localIdentity)
	if err != nil {
		return err
	}

	// we have removed content, so mark it finalized by us
	localIdentity, err = d.retryOnConflictError(localIdentity, d.finalizeLocalIdentity)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do localIdentity deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check if we can delete now.
	if d.deleteLocalIdentityWhenDone && finalized(localIdentity) {
		return d.deleteLocalIdentity(localIdentity)
	}

	return nil
}

// Deletes the given localIdentity.
func (d *loalIdentitiedResourcesDeleter) deleteLocalIdentity(localIdentity *v1.LocalIdentity) error {
	var opts *metav1.DeleteOptions
	uid := localIdentity.UID
	if len(uid) > 0 {
		opts = &metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	log.Info("localIdentity", log.Any("localIdentity", localIdentity))
	err := d.localIdentityClient.Delete(localIdentity.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		log.Error("error", log.Err(err))
		return err
	}
	return nil
}

// updateLocalIdentityFunc is a function that makes an update to a localIdentity
type updateLocalIdentityFunc func(localIdentity *v1.LocalIdentity) (*v1.LocalIdentity, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in client code
func (d *loalIdentitiedResourcesDeleter) retryOnConflictError(localIdentity *v1.LocalIdentity, fn updateLocalIdentityFunc) (result *v1.LocalIdentity, err error) {
	latestLocalIdentity := localIdentity
	for {
		result, err = fn(latestLocalIdentity)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevLocalIdentity := latestLocalIdentity
		latestLocalIdentity, err = d.localIdentityClient.Get(latestLocalIdentity.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevLocalIdentity.UID != latestLocalIdentity.UID {
			return nil, fmt.Errorf("localIdentity uid has changed across retries")
		}
	}
}

// updateLocalIdentityStatusFunc will verify that the status of the localIdentity is correct
func (d *loalIdentitiedResourcesDeleter) updateLocalIdentityStatusFunc(localIdentity *v1.LocalIdentity) (*v1.LocalIdentity, error) {
	if localIdentity.DeletionTimestamp.IsZero() || localIdentity.Status.Phase == v1.LocalIdentityDeleting {
		return localIdentity, nil
	}
	newLocalIdentity := v1.LocalIdentity{}
	newLocalIdentity.ObjectMeta = localIdentity.ObjectMeta
	newLocalIdentity.Status = localIdentity.Status
	newLocalIdentity.Status.Phase = v1.LocalIdentityDeleting
	return d.localIdentityClient.UpdateStatus(&newLocalIdentity)
}

// finalized returns true if the localIdentity.Spec.Finalizers is an empty list
func finalized(localIdentity *v1.LocalIdentity) bool {
	return len(localIdentity.Spec.Finalizers) == 0
}

// finalizeLocalIdentity removes the specified finalizerToken and finalizes the localIdentity
func (d *loalIdentitiedResourcesDeleter) finalizeLocalIdentity(localIdentity *v1.LocalIdentity) (*v1.LocalIdentity, error) {
	localIdentityFinalize := v1.LocalIdentity{}
	localIdentityFinalize.ObjectMeta = localIdentity.ObjectMeta
	localIdentityFinalize.Spec = localIdentity.Spec
	finalizerSet := sets.NewString()
	for i := range localIdentity.Spec.Finalizers {
		if localIdentity.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(localIdentity.Spec.Finalizers[i]))
		}
	}
	localIdentityFinalize.Spec.Finalizers = make([]v1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		localIdentityFinalize.Spec.Finalizers = append(localIdentityFinalize.Spec.Finalizers, v1.FinalizerName(value))
	}

	localIdentity = &v1.LocalIdentity{}
	err := d.authClient.RESTClient().Put().
		Resource("localidentities").
		Name(localIdentityFinalize.Name).
		SubResource("finalize").
		Body(&localIdentityFinalize).
		Do().
		Into(localIdentity)

	if err != nil {
		// it was removed already, so life is good
		if errors.IsNotFound(err) {
			return localIdentity, nil
		}
	}
	return localIdentity, err
}

type deleteResourceFunc func(deleter *loalIdentitiedResourcesDeleter, localIdentity *v1.LocalIdentity) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteRelatedRules,
	//deleteApikeys,
}

// deleteAllContent will use the dynamic client to delete each resource identified in groupVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *loalIdentitiedResourcesDeleter) deleteAllContent(localIdentity *v1.LocalIdentity) error {
	log.Debug("LocalIdentity controller - deleteAllContent", log.String("localIdentityName", localIdentity.ObjectMeta.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(d, localIdentity)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.Debug("LocalIdentity controller - deletedAllContent", log.String("localIdentityName", localIdentity.ObjectMeta.Name))
	return nil
}

func deleteRelatedRules(deleter *loalIdentitiedResourcesDeleter, localIdentity *v1.LocalIdentity) error {
	log.Debug("LocalIdentity controller - deleteRelatedRules", log.String("localIdentityName", localIdentity.ObjectMeta.Name))
	_, err := deleter.enforcer.DeleteRolesForUser(util.UserKey(localIdentity.Spec.TenantID, localIdentity.Spec.Username))
	return err
}
