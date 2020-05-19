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

// CustomPolicyBindingResourcesDeleterInterface to delete a policy with all resources in
// it.
type CustomPolicyBindingResourcesDeleterInterface interface {
	Delete(ctx context.Context, namespace, policyName string) error
}

// NewCustomPolicyBindingResourcesDeleter to create the customPolicyBindingResourcesDeleter that
// implement the CustomPolicyBindingResourcesDeleterInterface by given client and
// configure.
func NewCustomPolicyBindingResourcesDeleter(
	authClient v1clientset.AuthV1Interface,
	enforcer *casbin.SyncedEnforcer,
	finalizerToken v1.FinalizerName,
	deleteCustomPolicyWhenDone bool) CustomPolicyBindingResourcesDeleterInterface {
	d := &customPolicyBindingResourcesDeleter{
		authClient:                 authClient,
		enforcer:                   enforcer,
		finalizerToken:             finalizerToken,
		deleteCustomPolicyWhenDone: deleteCustomPolicyWhenDone,
	}
	return d
}

var _ CustomPolicyBindingResourcesDeleterInterface = &customPolicyBindingResourcesDeleter{}

// customPolicyBindingResourcesDeleter is used to delete all resources in a given binding.
type customPolicyBindingResourcesDeleter struct {
	authClient v1clientset.AuthV1Interface

	enforcer *casbin.SyncedEnforcer
	// The finalizer token that should be removed from the policy
	// when all resources in that policy have been deleted.
	finalizerToken v1.FinalizerName
	// Also delete the policy when all resources in the policy have been deleted.
	deleteCustomPolicyWhenDone bool
}

// Delete deletes all resources in the given binding.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   policy (does nothing if deletion timestamp is missing).
// * Verifies that the policy is in the "terminating" phase
//   (updates the policy phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given binding.
// * Deletes the policy if deleteCustomPolicyWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *customPolicyBindingResourcesDeleter) Delete(ctx context.Context, namespace, policyName string) error {
	// Multiple controllers may edit a policy during termination
	// first get the latest state of the policy before proceeding
	// if the policy was deleted already, don't do anything
	binding, err := d.authClient.CustomPolicyBindings(namespace).Get(ctx, policyName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if binding.DeletionTimestamp == nil {
		return nil
	}

	log.Infof("project policy controller - syncPolicy - policy: %s, finalizerToken: %s", binding.Name, d.finalizerToken)

	// ensure that the status is up to date on the policy
	// if we get a not found error, we assume the policy is truly gone
	binding, err = d.retryOnConflictError(ctx, binding, d.updatePolicyStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the policy asserts that policy is no longer deleting..
	if binding.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the policy if it is already finalized.
	if d.deleteCustomPolicyWhenDone && finalized(binding) {
		return d.deletePolicy(ctx, binding)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(ctx, binding)
	if err != nil {
		return err
	}

	// we have removed content, so mark it finalized by us
	binding, err = d.retryOnConflictError(ctx, binding, d.finalizeCustomPolicy)
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
	if d.deleteCustomPolicyWhenDone && finalized(binding) {
		return d.deletePolicy(ctx, binding)
	}

	return nil
}

// Deletes the given binding.
func (d *customPolicyBindingResourcesDeleter) deletePolicy(ctx context.Context, binding *v1.CustomPolicyBinding) error {
	var opts metav1.DeleteOptions
	uid := binding.UID
	if len(uid) > 0 {
		opts = metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.authClient.CustomPolicyBindings(binding.Namespace).Delete(ctx, binding.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		log.Error("error", log.Err(err))
		return err
	}
	return nil
}

// updateCustomPolicyFunc is a function that makes an update to a project policy
type updateCustomPolicyFunc func(ctx context.Context, binding *v1.CustomPolicyBinding) (*v1.CustomPolicyBinding, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in client code
func (d *customPolicyBindingResourcesDeleter) retryOnConflictError(ctx context.Context, binding *v1.CustomPolicyBinding, fn updateCustomPolicyFunc) (result *v1.CustomPolicyBinding, err error) {
	latestPolicy := binding
	for {
		result, err = fn(ctx, latestPolicy)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevPolicy := latestPolicy
		latestPolicy, err = d.authClient.CustomPolicyBindings(latestPolicy.Namespace).Get(ctx, latestPolicy.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevPolicy.UID != latestPolicy.UID {
			return nil, fmt.Errorf("policy uid has changed across retries")
		}
	}
}

// updatePolicyStatusFunc will verify that the status of the policy is correct
func (d *customPolicyBindingResourcesDeleter) updatePolicyStatusFunc(ctx context.Context, binding *v1.CustomPolicyBinding) (*v1.CustomPolicyBinding, error) {
	if binding.DeletionTimestamp.IsZero() || binding.Status.Phase == v1.BindingTerminating {
		return binding, nil
	}
	newPolicy := v1.CustomPolicyBinding{}
	newPolicy.ObjectMeta = binding.ObjectMeta
	newPolicy.Status = binding.Status
	newPolicy.Status.Phase = v1.BindingTerminating
	return d.authClient.CustomPolicyBindings(newPolicy.Namespace).UpdateStatus(ctx, &newPolicy, metav1.UpdateOptions{})
}

// finalized returns true if the binding.Spec.Finalizers is an empty list
func finalized(binding *v1.CustomPolicyBinding) bool {
	return len(binding.Spec.Finalizers) == 0
}

// finalizeCustomPolicy removes the specified finalizerToken and finalizes the policy
func (d *customPolicyBindingResourcesDeleter) finalizeCustomPolicy(ctx context.Context, binding *v1.CustomPolicyBinding) (*v1.CustomPolicyBinding, error) {
	policyFinalize := v1.CustomPolicyBinding{}
	policyFinalize.ObjectMeta = binding.ObjectMeta
	policyFinalize.Spec = binding.Spec
	finalizerSet := sets.NewString()
	for i := range binding.Spec.Finalizers {
		if binding.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(binding.Spec.Finalizers[i]))
		}
	}
	policyFinalize.Spec.Finalizers = make([]v1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		policyFinalize.Spec.Finalizers = append(policyFinalize.Spec.Finalizers, v1.FinalizerName(value))
	}

	updated := &v1.CustomPolicyBinding{}
	err := d.authClient.RESTClient().Put().
		Resource("CustomPolicyBindings").
		Namespace(policyFinalize.Namespace).
		Name(policyFinalize.Name).
		SubResource("finalize").
		Body(&policyFinalize).
		Do(ctx).
		Into(updated)

	if err != nil {
		return nil, err
	}
	return updated, err
}

type deleteResourceFunc func(ctx context.Context, deleter *customPolicyBindingResourcesDeleter, binding *v1.CustomPolicyBinding) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteRelatedRules,
}

// deleteAllContent will use the dynamic client to delete each resource identified in groupVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *customPolicyBindingResourcesDeleter) deleteAllContent(ctx context.Context, binding *v1.CustomPolicyBinding) error {
	log.Debug("CustomPolicyBinding controller - deleteAllContent", log.String("namespace", binding.Namespace), log.String("policyName", binding.Name))

	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(ctx, d, binding)
		if err != nil {
			// If there is an error, return directly, in case delete roles failed in next try if rule has been deleted.
			return err
		}
	}

	return nil
}

func deleteRelatedRules(ctx context.Context, deleter *customPolicyBindingResourcesDeleter, binding *v1.CustomPolicyBinding) error {
	log.Info("CustomPolicyBinding controller - deleteRelatedRules", log.String("namespace", binding.Namespace), log.String("policyName", binding.Name))
	var errs []error
	role := authutil.RoleKey(binding.Spec.RulePrefix, binding.Spec.PolicyID)

	expectCustomSubj := sets.String{}
	for _, subj := range binding.Spec.Users {
		expectCustomSubj.Insert(authutil.UserKey(binding.Spec.TenantID, subj.Name))
	}
	for _, subj := range binding.Spec.Groups {
		expectCustomSubj.Insert(authutil.GroupKey(binding.Spec.TenantID, subj.Name))
	}

	subjs := expectCustomSubj.UnsortedList()
	for _, subj := range subjs {
		if _, err := deleter.enforcer.DeleteRoleForUserInDomain(subj, role, binding.Spec.Domain); err != nil {
			errs = append(errs, err)
		}
	}

	rules := deleter.enforcer.GetFilteredGroupingPolicy(1, role, binding.Spec.Domain)
	if len(rules) == 0 {
		if _, err := deleter.enforcer.RemoveFilteredPolicy(0, role); err != nil {
			errs = append(errs, err)
		}
	}

	return utilerrors.NewAggregate(errs)
}
