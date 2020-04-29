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
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/provider/imported/util/mark"
	"tkestack.io/tke/pkg/platform/types"
	"tkestack.io/tke/pkg/spec"
	"tkestack.io/tke/pkg/util/apiclient"
	utilvalidation "tkestack.io/tke/pkg/util/validation"
)

// ValidateCluster validates a given Cluster.
func ValidateCluster(cluster *types.Cluster) field.ErrorList {
	allErrs := ValidatClusterAddresses(cluster.Status.Addresses, field.NewPath("status", "addresses"))

	if cluster.Spec.ClusterCredentialRef != nil {
		allErrs = append(allErrs, ValidatClusterCredentialRef(cluster, field.NewPath("spec", "clusterCredentialRef"))...)

		client, err := cluster.Clientset()
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("name"), cluster.Name, fmt.Sprintf("get clientset error: %w", err)))
		}
		if cluster.Status.Phase == platform.ClusterInitializing {
			allErrs = append(allErrs, ValidateClusterMark(cluster.Name, field.NewPath("name"), client)...)
			allErrs = append(allErrs, ValidateClusterVersion(field.NewPath("name"), client)...)
		}
	}

	return allErrs
}

// ValidatClusterAddresses validates a given ClusterAddresses.
func ValidatClusterAddresses(addresses []platform.ClusterAddress, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(addresses) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("status", "addresses"), "must specify at least one obj access address"))
	} else {
		for i, address := range addresses {
			fldPath := fldPath.Index(i)
			allErrs = utilvalidation.ValidateEnum(address.Type, fldPath.Child("type"), []interface{}{
				platform.AddressAdvertise,
				platform.AddressReal,
			})
			if address.Host == "" {
				allErrs = append(allErrs, field.Required(fldPath.Child("host"), "must specify host"))
			}
			for _, msg := range validation.IsValidPortNum(int(address.Port)) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("port"), address.Port, msg))
			}
			err := utilvalidation.IsHTTPSReachle(address.Host, address.Port, 5*time.Second)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath, address, err.Error()))
			}
		}
	}
	return allErrs
}

// ValidatClusterCredentialRef validates cluster.Spec.ClusterCredentialRef.
func ValidatClusterCredentialRef(cluster *types.Cluster, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	credential := cluster.ClusterCredential
	if credential.Token == nil && credential.ClientKey == nil && credential.ClientCert == nil {
		allErrs = append(allErrs, field.Required(field.NewPath(""),
			"must specify at least one of token or client certificate authentication"))

		return allErrs
	}

	if credential.ClientCert == nil && credential.ClientKey != nil ||
		credential.ClientCert != nil && credential.ClientKey == nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("clientCert"),
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

	return allErrs
}

// ValidateClusterMark validates a given cluster had imported already.
func ValidateClusterMark(clusterName string, fldPath *field.Path, client kubernetes.Interface) field.ErrorList {
	allErrs := field.ErrorList{}

	_, err := mark.Get(client)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			allErrs = append(allErrs, field.InternalError(fldPath, err))
		}
	} else {
		allErrs = append(allErrs, field.Invalid(fldPath, clusterName,
			fmt.Sprintf("can't imported same cluster, you can use `kubectl -n%s delete configmap %s`", mark.Namespace, mark.Name)))
	}

	return allErrs

}

// ValidateClusterMark validates a given cluster's version.
func ValidateClusterVersion(fldPath *field.Path, client kubernetes.Interface) field.ErrorList {
	allErrs := field.ErrorList{}

	v, err := apiclient.GetClusterVersion(client)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath, err))
		return allErrs
	}

	result, err := apiclient.CheckVersion(v, spec.K8sVersionConstraint)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath, err))
		return allErrs
	}
	if !result {
		allErrs = append(allErrs, field.Invalid(fldPath, v, fmt.Sprintf("cluster version must %s", spec.K8sVersionConstraint)))
	}

	return allErrs

}
