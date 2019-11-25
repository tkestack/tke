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
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	runtimeutil "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	// Scheme is the default instance of runtime.Scheme to which types in the TKE API are already registered.
	Scheme = runtime.NewScheme()
	// Codecs provides access to encoding and decoding for the scheme
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	Install(Scheme)
}

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	runtimeutil.Must(registerMeta(scheme))
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
	// ListOptions is the only options struct which needs conversion (it exposes labels and fields
	// as selectors for convenience). The other types have only a single representation today.
	scheme.AddKnownTypes(schemaGroupVersion,
		&metav1.Status{},
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
