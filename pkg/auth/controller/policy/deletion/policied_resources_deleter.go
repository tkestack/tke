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

	"github.com/casbin/casbin/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	v1 "tkestack.io/tke/api/auth/v1"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	"tkestack.io/tke/pkg/util/log"
)

// PoliciedResourcesDeleterInterface to delete a policy with all resources in
// it.
type PoliciedResourcesDeleterInterface interface {
	Delete(policyName string) error
}

// NewPoliciedResourcesDeleter to create the policiedResourcesDeleter that
// implement the PoliciedResourcesDeleterInterface by given client and
// configure.
func NewPoliciedResourcesDeleter(pilicyClient v1clientset.PolicyInterface,
	authClient v1clientset.AuthV1Interface,
	enforcer *casbin.SyncedEnforcer,
	finalizerToken v1.FinalizerName,
	deletePolicyWhenDone bool) PoliciedResourcesDeleterInterface {
	d := &policiedResourcesDeleter{
		policyClient:         pilicyClient,
		authClient:           authClient,
		enforcer:             enforcer,
		finalizerToken:       finalizerToken,
		deletePolicyWhenDone: deletePolicyWhenDone,
	}
	return d
}

var _ PoliciedResourcesDeleterInterface = &policiedResourcesDeleter{}

// policiedResourcesDeleter is used to delete all resources in a given policy.
type policiedResourcesDeleter struct {
	// Client to manipulate the policy.
	policyClient v1clientset.PolicyInterface
	authClient   v1clientset.AuthV1Interface

	enforcer *casbin.SyncedEnforcer
	// The finalizer token that should be removed from the policy
	// when all resources in that policy have been deleted.
	finalizerToken v1.FinalizerName
	// Also delete the policy when all resources in the policy have been deleted.
	deletePolicyWhenDone bool
}

// Delete deletes all resources in the given policy.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   policy (does nothing if deletion timestamp is missing).
// * Verifies that the policy is in the "terminating" phase
//   (updates the policy phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given policy.
// * Deletes the policy if deletePolicyWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *policiedResourcesDeleter) Delete(policyName string) error {
	// Multiple controllers may edit a policy during termination
	// first get the latest state of the policy before proceeding
	// if the policy was deleted already, don't do anything
	policy, err := d.policyClient.Get(policyName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if policy.DeletionTimestamp == nil {
		return nil
	}

	log.Infof("policy controller - syncPolicy - policy: %s, finalizerToken: %s", policy.Name, d.finalizerToken)

	// ensure that the status is up to date on the policy
	// if we get a not found error, we assume the policy is truly gone
	policy, err = d.retryOnConflictError(policy, d.updatePolicyStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the policy asserts that policy is no longer deleting..
	if policy.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the policy if it is already finalized.
	if d.deletePolicyWhenDone && finalized(policy) {
		return d.deletePolicy(policy)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(policy)
	if err != nil {
		return err
	}

	// we have removed content, so mark it finalized by us
	policy, err = d.retryOnConflictError(policy, d.finalizePolicy)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do policy deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check if we can delete now.
	if d.deletePolicyWhenDone && finalized(policy) {
		return d.deletePolicy(policy)
	}

	return nil
}

// Deletes the given policy.
func (d *policiedResourcesDeleter) deletePolicy(policy *v1.Policy) error {
	var opts *metav1.DeleteOptions
	uid := policy.UID
	if len(uid) > 0 {
		opts = &metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	log.Info("policy", log.Any("policy", policy))
	err := d.policyClient.Delete(policy.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		log.Error("error", log.Err(err))
		return err
	}
	return nil
}

// updatePolicyFunc is a function that makes an update to a policy
type updatePolicyFunc func(policy *v1.Policy) (*v1.Policy, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in client code
func (d *policiedResourcesDeleter) retryOnConflictError(policy *v1.Policy, fn updatePolicyFunc) (result *v1.Policy, err error) {
	latestPolicy := policy
	for {
		result, err = fn(latestPolicy)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevPolicy := latestPolicy
		latestPolicy, err = d.policyClient.Get(latestPolicy.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevPolicy.UID != latestPolicy.UID {
			return nil, fmt.Errorf("policy uid has changed across retries")
		}
	}
}

// updatePolicyStatusFunc will verify that the status of the policy is correct
func (d *policiedResourcesDeleter) updatePolicyStatusFunc(policy *v1.Policy) (*v1.Policy, error) {
	if policy.DeletionTimestamp.IsZero() || policy.Status.Phase == v1.PolicyTerminating {
		return policy, nil
	}
	newPolicy := v1.Policy{}
	newPolicy.ObjectMeta = policy.ObjectMeta
	newPolicy.Status = policy.Status
	newPolicy.Status.Phase = v1.PolicyTerminating
	return d.policyClient.UpdateStatus(&newPolicy)
}

// finalized returns true if the policy.Spec.Finalizers is an empty list
func finalized(policy *v1.Policy) bool {
	return len(policy.Spec.Finalizers) == 0
}

// finalizePolicy removes the specified finalizerToken and finalizes the policy
func (d *policiedResourcesDeleter) finalizePolicy(policy *v1.Policy) (*v1.Policy, error) {
	policyFinalize := v1.Policy{}
	policyFinalize.ObjectMeta = policy.ObjectMeta
	policyFinalize.Spec = policy.Spec
	finalizerSet := sets.NewString()
	for i := range policy.Spec.Finalizers {
		if policy.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(policy.Spec.Finalizers[i]))
		}
	}
	policyFinalize.Spec.Finalizers = make([]v1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		policyFinalize.Spec.Finalizers = append(policyFinalize.Spec.Finalizers, v1.FinalizerName(value))
	}

	policy = &v1.Policy{}
	err := d.authClient.RESTClient().Put().
		Resource("policies").
		Name(policyFinalize.Name).
		SubResource("finalize").
		Body(&policyFinalize).
		Do().
		Into(policy)

	if err != nil {
		// it was removed already, so life is good
		if errors.IsNotFound(err) {
			return policy, nil
		}
	}
	return policy, err
}

type deleteResourceFunc func(deleter *policiedResourcesDeleter, policy *v1.Policy) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteRelatedRules,
}

// deleteAllContent will use the dynamic client to delete each resource identified in groupVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *policiedResourcesDeleter) deleteAllContent(policy *v1.Policy) error {
	log.Debug("Policy controller - deleteAllContent", log.String("policyName", policy.ObjectMeta.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(d, policy)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.Debug("Policy controller - deletedAllContent", log.String("policyName", policy.ObjectMeta.Name))
	return nil
}

func deleteRelatedRules(deleter *policiedResourcesDeleter, policy *v1.Policy) error {
	log.Info("Policy controller - deleteRelatedRules", log.String("policyName", policy.ObjectMeta.Name))
	_, err := deleter.enforcer.DeleteRole(policy.Name)
	return err
}
