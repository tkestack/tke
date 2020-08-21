/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"tkestack.io/tke/pkg/platform/proxy"

	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
	deploymentutil "k8s.io/kubectl/pkg/util/deployment"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
)

// RolloutUndoREST implements Creater
type RolloutUndoREST struct {
	rest.Storage
	platformClient platforminternalclient.PlatformInterface
}

var _ = rest.NamedCreater(&RolloutUndoREST{})
var _ = rest.GroupVersionKindProvider(&RolloutUndoREST{})

// GroupVersionKind is used to specify a particular GroupVersionKind to discovery.
func (r *RolloutUndoREST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return appsv1beta1.SchemeGroupVersion.WithKind("DeploymentRollback")
}

// New creates a new Rollback object
func (r *RolloutUndoREST) New() runtime.Object {
	return &appsv1beta1.DeploymentRollback{}
}

type rollbacker struct {
	client kubernetes.Interface
	gs     schema.GroupResource
}

func newDeploymentRollbacker(ctx context.Context, client kubernetes.Interface) *rollbacker {
	return &rollbacker{
		client: client,
		gs:     appsv1beta1.Resource("DeploymentRollback"),
	}
}

// Create inserts a new item according to the unique key from the object.
func (r *RolloutUndoREST) Create(ctx context.Context, name string, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	client, err := proxy.ClientSet(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	rollbackObj, ok := obj.(*appsv1beta1.DeploymentRollback)
	if !ok {
		return nil, errors.NewBadRequest(fmt.Sprintf("not a DeploymentRollback: %#v", obj))
	}
	if name != rollbackObj.Name {
		return nil, errors.NewBadRequest("name in URL does not match name in DeploymentRollback object")
	}

	if createValidation != nil {
		if err := createValidation(ctx, obj.DeepCopyObject()); err != nil {
			return nil, err
		}
	}

	rb := newDeploymentRollbacker(ctx, client)
	if err := rb.rollback(ctx, rollbackObj); err != nil {
		return nil, err
	}

	return &metav1.Status{
		Status:  metav1.StatusSuccess,
		Code:    http.StatusOK,
		Message: fmt.Sprintf("rollback request for deployment \"%s\" succeeded", rollbackObj.Name),
	}, nil
}

func (r *rollbacker) rollback(ctx context.Context, obj *appsv1beta1.DeploymentRollback) error {
	if errList := r.validateDeploymentRollback(obj); len(errList) != 0 {
		return errors.NewInvalid(obj.GroupVersionKind().GroupKind(), obj.Name, errList)
	}

	namespace, ok := request.NamespaceFrom(ctx)
	if !ok {
		return errors.NewBadRequest("a namespace must be specified")
	}

	deployment, err := r.client.AppsV1().Deployments(namespace).Get(ctx, obj.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	rsForRevision, err := r.deploymentRevision(deployment, obj.RollbackTo.Revision)
	if err != nil {
		return err
	}
	if deployment.Spec.Paused {
		return errors.NewConflict(
			r.gs,
			deployment.Name,
			fmt.Errorf("skipped rollback (deployment \"%s\" is paused)", deployment.Name))
	}

	// Skip if the revision already matches current Deployment
	if equalIgnoreHash(&rsForRevision.Spec.Template, &deployment.Spec.Template) {
		return errors.NewConflict(
			r.gs,
			deployment.Name,
			fmt.Errorf("skipped rollback (current template already matches revision %d)", obj.RollbackTo.Revision))
	}

	// remove hash label before patching back into the deployment
	delete(rsForRevision.Spec.Template.Labels, appsv1.DefaultDeploymentUniqueLabelKey)

	// compute deployment annotations
	annotations := map[string]string{}
	for k := range annotationsToSkip {
		if v, ok := deployment.Annotations[k]; ok {
			annotations[k] = v
		}
	}
	for k, v := range rsForRevision.Annotations {
		if !annotationsToSkip[k] {
			annotations[k] = v
		}
	}

	// make patch to restore
	patchType, patch, err := r.getDeploymentPatch(&rsForRevision.Spec.Template, annotations)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("failed restoring revision %d: %v", obj.RollbackTo.Revision, err))
	}

	// Restore revision
	if _, err = r.client.AppsV1().Deployments(namespace).Patch(ctx, deployment.Name, patchType, patch, metav1.PatchOptions{}); err != nil {
		return errors.NewInternalError(fmt.Errorf("failed restoring revision %d: %v", obj.RollbackTo.Revision, err))
	}

	return nil
}

func (r *rollbacker) validateDeploymentRollback(obj *appsv1beta1.DeploymentRollback) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateAnnotations(obj.UpdatedAnnotations, field.NewPath("updatedAnnotations"))
	if len(obj.Name) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("name"), "name is required"))
	}
	allErrs = append(allErrs, r.validateRollback(&obj.RollbackTo, field.NewPath("rollback"))...)
	return allErrs
}

func (r *rollbacker) validateRollback(rollback *appsv1beta1.RollbackConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	v := rollback.Revision
	allErrs = append(allErrs, apimachineryvalidation.ValidateNonnegativeField(v, fldPath.Child("version"))...)

	return allErrs
}

func (r *rollbacker) deploymentRevision(deployment *appsv1.Deployment, toRevision int64) (revision *appsv1.ReplicaSet, err error) {
	_, allOldRSs, newRS, err := deploymentutil.GetAllReplicaSets(deployment, r.client.AppsV1())
	if err != nil {
		return nil, err
	}
	allRSs := allOldRSs
	if newRS != nil {
		allRSs = append(allRSs, newRS)
	}

	var (
		latestReplicaSet   *appsv1.ReplicaSet
		latestRevision     = int64(-1)
		previousReplicaSet *appsv1.ReplicaSet
		previousRevision   = int64(-1)
	)
	for _, rs := range allRSs {
		if v, err := deploymentutil.Revision(rs); err == nil {
			if toRevision == 0 {
				if latestRevision < v {
					// newest one we've seen so far
					previousRevision = latestRevision
					previousReplicaSet = latestReplicaSet
					latestRevision = v
					latestReplicaSet = rs
				} else if previousRevision < v {
					// second newest one we've seen so far
					previousRevision = v
					previousReplicaSet = rs
				}
			} else if toRevision == v {
				return rs, nil
			}
		}
	}

	if toRevision > 0 {
		return nil, errors.NewConflict(r.gs, deployment.Name, fmt.Errorf("unable to find specified revision %v in history", toRevision))
	}

	if previousReplicaSet == nil {
		return nil, errors.NewConflict(r.gs, deployment.Name, fmt.Errorf("no rollout history found for deployment %q", deployment.Name))
	}
	return previousReplicaSet, nil
}

func (r *rollbacker) getDeploymentPatch(podTemplate *corev1.PodTemplateSpec, annotations map[string]string) (types.PatchType, []byte, error) {
	// Create a patch of the Deployment that replaces spec.template
	patch, err := json.Marshal([]interface{}{
		map[string]interface{}{
			"op":    "replace",
			"path":  "/spec/template",
			"value": podTemplate,
		},
		map[string]interface{}{
			"op":    "replace",
			"path":  "/metadata/annotations",
			"value": annotations,
		},
	})
	return types.JSONPatchType, patch, err
}

func equalIgnoreHash(template1, template2 *corev1.PodTemplateSpec) bool {
	t1Copy := template1.DeepCopy()
	t2Copy := template2.DeepCopy()
	// Remove hash labels from template.Labels before comparing
	delete(t1Copy.Labels, appsv1.DefaultDeploymentUniqueLabelKey)
	delete(t2Copy.Labels, appsv1.DefaultDeploymentUniqueLabelKey)
	return apiequality.Semantic.DeepEqual(t1Copy, t2Copy)
}

// annotationsToSkip lists the annotations that should be preserved from the deployment and not
// copied from the replicaset when rolling a deployment back
var annotationsToSkip = map[string]bool{
	corev1.LastAppliedConfigAnnotation:       true,
	deploymentutil.RevisionAnnotation:        true,
	deploymentutil.RevisionHistoryAnnotation: true,
	deploymentutil.DesiredReplicasAnnotation: true,
	deploymentutil.MaxReplicasAnnotation:     true,
	appsv1.DeprecatedRollbackTo:              true,
}
