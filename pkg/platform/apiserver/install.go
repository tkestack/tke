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

package apiserver

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	runtimeutil "k8s.io/apimachinery/pkg/util/runtime"
	"tkestack.io/tke/api/platform"
	// register project group api scheme for api server.
	admissionv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	certV1beta1 "k8s.io/api/certificates/v1beta1"
	coordinationv1 "k8s.io/api/coordination/v1"
	coordinationv1beta1 "k8s.io/api/coordination/v1beta1"
	corev1 "k8s.io/api/core/v1"
	eventsv1beta1 "k8s.io/api/events/v1beta1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1Beta1 "k8s.io/api/networking/v1beta1"
	_ "tkestack.io/tke/api/platform/install"

	nodev1alpha1 "k8s.io/api/node/v1alpha1"
	nodev1beta1 "k8s.io/api/node/v1beta1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1alpha1 "k8s.io/api/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"

	schedulingv1alpha1 "k8s.io/api/scheduling/v1alpha1"
	schedulingv1beta1 "k8s.io/api/scheduling/v1beta1"

	schedulingv1 "k8s.io/api/scheduling/v1"
	settingsv1alpha1 "k8s.io/api/settings/v1alpha1"
	storagev1 "k8s.io/api/storage/v1"
	storagev1alpha1 "k8s.io/api/storage/v1alpha1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
)

func init() {
	Install(platform.Scheme)
}

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	runtimeutil.Must(corev1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(corev1.SchemeGroupVersion))

	runtimeutil.Must(admissionv1beta1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(admissionv1beta1.SchemeGroupVersion))

	runtimeutil.Must(autoscalingv1.AddToScheme(scheme))
	runtimeutil.Must(autoscalingv2beta1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(autoscalingv1.SchemeGroupVersion, autoscalingv2beta1.SchemeGroupVersion))

	runtimeutil.Must(appsv1.AddToScheme(scheme))
	runtimeutil.Must(appsv1beta1.AddToScheme(scheme))
	runtimeutil.Must(appsv1beta2.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(appsv1.SchemeGroupVersion, appsv1beta2.SchemeGroupVersion, appsv1beta1.SchemeGroupVersion))

	runtimeutil.Must(extensionsv1beta1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(extensionsv1beta1.SchemeGroupVersion))

	runtimeutil.Must(eventsv1beta1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(eventsv1beta1.SchemeGroupVersion))

	runtimeutil.Must(batchv1.AddToScheme(scheme))
	runtimeutil.Must(batchv1beta1.AddToScheme(scheme))
	runtimeutil.Must(batchv2alpha1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(batchv1.SchemeGroupVersion, batchv1beta1.SchemeGroupVersion, batchv2alpha1.SchemeGroupVersion))

	runtimeutil.Must(networkingv1.AddToScheme(scheme))
	runtimeutil.Must(networkingv1Beta1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(networkingv1.SchemeGroupVersion, networkingv1Beta1.SchemeGroupVersion))

	runtimeutil.Must(coordinationv1.AddToScheme(scheme))
	runtimeutil.Must(coordinationv1beta1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(coordinationv1.SchemeGroupVersion, coordinationv1beta1.SchemeGroupVersion))

	runtimeutil.Must(certV1beta1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(certV1beta1.SchemeGroupVersion))

	runtimeutil.Must(policyv1beta1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(policyv1beta1.SchemeGroupVersion))

	runtimeutil.Must(rbacv1.AddToScheme(scheme))
	runtimeutil.Must(rbacv1beta1.AddToScheme(scheme))
	runtimeutil.Must(rbacv1alpha1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(rbacv1.SchemeGroupVersion, rbacv1alpha1.SchemeGroupVersion, rbacv1beta1.SchemeGroupVersion))

	runtimeutil.Must(schedulingv1alpha1.AddToScheme(scheme))
	runtimeutil.Must(schedulingv1beta1.AddToScheme(scheme))
	runtimeutil.Must(schedulingv1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(schedulingv1beta1.SchemeGroupVersion, schedulingv1alpha1.SchemeGroupVersion, schedulingv1.SchemeGroupVersion))

	runtimeutil.Must(nodev1alpha1.AddToScheme(scheme))
	runtimeutil.Must(nodev1beta1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(nodev1beta1.SchemeGroupVersion, nodev1alpha1.SchemeGroupVersion))

	runtimeutil.Must(settingsv1alpha1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(settingsv1alpha1.SchemeGroupVersion))

	runtimeutil.Must(storagev1alpha1.AddToScheme(scheme))
	runtimeutil.Must(storagev1beta1.AddToScheme(scheme))
	runtimeutil.Must(storagev1.AddToScheme(scheme))
	runtimeutil.Must(scheme.SetVersionPriority(storagev1.SchemeGroupVersion, storagev1beta1.SchemeGroupVersion, storagev1alpha1.SchemeGroupVersion))

	runtimeutil.Must(registerMeta(scheme))
	runtimeutil.Must(registerKubernetesConversion(scheme))
}

func registerMeta(scheme *runtime.Scheme) error {
	schemaGroupVersion := schema.GroupVersion{Group: metainternal.GroupName, Version: runtime.APIVersionInternal}

	if err := scheme.AddIgnoredConversionType(&metav1.TypeMeta{}, &metav1.TypeMeta{}); err != nil {
		return err
	}
	_ = scheme.AddConversionFuncs(
		metav1.Convert_string_To_labels_Selector,
		metav1.Convert_labels_Selector_To_string,

		metav1.Convert_string_To_fields_Selector,
		metav1.Convert_fields_Selector_To_string,

		metainternal.Convert_v1_List_To_internalversion_List,
		metainternal.Convert_internalversion_List_To_v1_List,

		metainternal.Convert_internalversion_ListOptions_To_v1_ListOptions,
		metainternal.Convert_v1_ListOptions_To_internalversion_ListOptions,
	)
	// ListOptions is the onl y options struct which needs conversion (it exposes labels and fields
	// as selectors for convenience). The other types have only a single representation today.
	scheme.AddKnownTypes(schemaGroupVersion,
		&metav1.GetOptions{},
		&metav1.ExportOptions{},
		&metav1.DeleteOptions{},
		&metav1.ListOptions{},
	)
	scheme.AddKnownTypes(schemaGroupVersion,
		&metav1beta1.Table{},
		&metav1beta1.TableOptions{},
		&metav1beta1.PartialObjectMetadata{},
		&metav1beta1.PartialObjectMetadataList{},
	)
	scheme.AddKnownTypes(metav1beta1.SchemeGroupVersion,
		&metav1beta1.Table{},
		&metav1beta1.TableOptions{},
		&metav1beta1.PartialObjectMetadata{},
		&metav1beta1.PartialObjectMetadataList{},
	)
	// Allow delete options to be decoded across all version in this scheme (we may want to be more clever than this)
	scheme.AddUnversionedTypes(schemaGroupVersion, &metav1.DeleteOptions{})
	metav1.AddToGroupVersion(scheme, metav1.SchemeGroupVersion)
	_ = metainternal.AddToScheme(scheme)
	return nil
}
