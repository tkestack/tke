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

package validation

import (
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/validation"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	utilvalidation "tkestack.io/tke/pkg/util/validation"
)

// ValidateClusterCredential validates a given ClusterCredential.
func ValidateClusterCredential(credential *platform.ClusterCredential, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := validation.ValidateObjectMeta(&credential.ObjectMeta, false, apimachineryvalidation.NameIsDNSLabel, field.NewPath("metadata"))

	if credential.ClusterName != "" {
		cluster, err := platformClient.Clusters().Get(credential.ClusterName, metav1.GetOptions{})
		if err != nil {
			return allErrs
		}

		// Deprecated: will remove in next release
		if cluster.Spec.Type == "Imported" {
			if credential.Token == nil && credential.ClientKey == nil && credential.ClientCert == nil {
				allErrs = append(allErrs, field.Required(field.NewPath(""),
					"must specify at least one of token or client certificate authentication"))

				return allErrs
			}

			if credential.ClientCert == nil && credential.ClientKey != nil ||
				credential.ClientCert != nil && credential.ClientKey == nil {
				allErrs = append(allErrs, field.Required(field.NewPath("clientCert"),
					"`clientCert` and `clientKey` must provide togther"))
			}

			host, err := cluster.Host()
			if err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("clusterName"), credential.ClusterName, err.Error()))
			} else {
				restConfig := &rest.Config{
					Host:    host,
					Timeout: 5 * time.Second,
				}
				if credential.CACert != nil {
					restConfig.CAData = credential.CACert
					if err = utilvalidation.ValidateRESTConfig(restConfig); err != nil {
						if !apierrors.IsUnauthorized(err) {
							allErrs = append(allErrs, field.Invalid(field.NewPath("caCert"), "", err.Error()))
						}
					}
				} else {
					restConfig.Insecure = true
				}
				if credential.Token != nil {
					config := rest.CopyConfig(restConfig)
					config.BearerToken = *credential.Token
					if err = utilvalidation.ValidateRESTConfig(config); err != nil {
						if apierrors.IsUnauthorized(err) {
							allErrs = append(allErrs, field.Invalid(field.NewPath("token"), *credential.Token, err.Error()))
						} else {
							allErrs = append(allErrs, field.InternalError(field.NewPath("token"), err))
						}
					}
				}
				if credential.ClientCert != nil && credential != nil {
					config := rest.CopyConfig(restConfig)
					config.TLSClientConfig.CertData = credential.ClientCert
					config.TLSClientConfig.KeyData = credential.ClientKey
					if err = utilvalidation.ValidateRESTConfig(config); err != nil {
						if apierrors.IsUnauthorized(err) {
							allErrs = append(allErrs, field.Invalid(field.NewPath("clientCert"), "", err.Error()))
						} else {
							allErrs = append(allErrs, field.InternalError(field.NewPath("clientCert"), err))
						}
					}
				}
			}
		}
	}

	return allErrs
}

// ValidateUpdateClusterCredential tests if an update to a ClusterCredential is valid.
func ValidateUpdateClusterCredential(newObj *platform.ClusterCredential, oldObj *platform.ClusterCredential, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&newObj.ObjectMeta, &oldObj.ObjectMeta, field.NewPath("metadata"))

	return allErrs
}
