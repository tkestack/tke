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

package apikey

import (
	"fmt"
	"time"

	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/validation/field"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"

	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/auth/util/sign"
	"tkestack.io/tke/pkg/util/log"
)

var (
	minExpire = time.Second
	maxExpire = 100 * 365 * 24 * time.Hour

	defaultAPIKeyTimeout = metav1.Duration{Duration: 7 * 24 * time.Hour}
)

// ValidateAPIkey tests if required fields in the signing key are set.
func ValidateAPIkey(apiKey *auth.APIKey, keySigner sign.KeySigner, privilegedUsername string) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&apiKey.ObjectMeta, false, apiMachineryValidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")

	if apiKey.Spec.APIkey == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("apiKey"), "must specify apiKey"))
	}

	if claims, err := keySigner.Verify(apiKey.Spec.APIkey); err != nil {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("apiKey"), apiKey.Spec.APIkey, err.Error()))
	} else {
		// if not super admin, must specify tenantID
		if apiKey.Spec.TenantID == "" && claims.UserName != privilegedUsername {
			allErrs = append(allErrs, field.Required(fldSpecPath.Child("tenantID"), "must specify tenantID"))
		}

		if apiKey.Spec.Username != claims.UserName {
			allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("username"), apiKey.Spec.Username, "must be same with username of  apikey"))
		}

		if claims.IssuedAt != apiKey.Spec.IssueAt.Unix() {
			allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("issue_at"), apiKey.Spec.IssueAt, "must be same with issue time of apiKey"))
		}

		if claims.ExpiresAt != apiKey.Spec.ExpireAt.Unix() {
			allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("expire_at"), apiKey.Spec.IssueAt, "must be same with expire time of apiKey"))
		}
	}

	return allErrs
}

// ValidateAPIKeyUpdate tests if required fields in the session are set during
// an update.
func ValidateAPIKeyUpdate(apiKey *auth.APIKey, oldAPIKey *auth.APIKey) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, apiMachineryValidation.ValidateObjectMetaUpdate(&apiKey.ObjectMeta, &oldAPIKey.ObjectMeta, field.NewPath("metadata"))...)

	fldSpecPath := field.NewPath("spec")
	if apiKey.Spec.APIkey != oldAPIKey.Spec.APIkey {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("apiKey"), apiKey.Spec.APIkey, "disallowed change the apiKey"))
	}

	if apiKey.Spec.TenantID != oldAPIKey.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("tenantID"), apiKey.Spec.TenantID, "disallowed change the tenantID"))
	}

	if apiKey.Spec.IssueAt != oldAPIKey.Spec.IssueAt {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("issue_at"), apiKey.Spec.IssueAt, "disallowed change the issue_at"))
	}

	if apiKey.Spec.ExpireAt != oldAPIKey.Spec.ExpireAt {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("expire_at"), apiKey.Spec.ExpireAt, "disallowed change the expire_at"))
	}

	if apiKey.Spec.Username != oldAPIKey.Spec.Username {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("username"), apiKey.Spec.ExpireAt, "disallowed change the username"))
	}

	return allErrs
}

func validateAPIKeyExpire(expire metav1.Duration) error {
	if expire.Duration < minExpire || expire.Duration > maxExpire {
		return fmt.Errorf("expire %v must not shorter than %v or longer than %v", expire, minExpire, maxExpire)
	}
	return nil
}

// ValidateAPIKeyReq tests if required fields in the signing key are set.
func ValidateAPIKeyReq(apiKeyReq *auth.APIKeyReq) error {
	allErrs := field.ErrorList{}

	if apiKeyReq.Expire.Duration == 0 {
		apiKeyReq.Expire = defaultAPIKeyTimeout
	}

	fldPath := field.NewPath("")
	if err := validateAPIKeyExpire(apiKeyReq.Expire); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("expire"), apiKeyReq.Expire, err.Error()))
	}

	return allErrs.ToAggregate()
}

// ValidateAPIkeyPassword tests if required fields in the signing key are set.
func ValidateAPIkeyPassword(apiKeyPass *auth.APIKeyReqPassword, authClient authinternalclient.AuthInterface) error {
	allErrs := field.ErrorList{}

	if apiKeyPass.Expire.Duration == 0 {
		apiKeyPass.Expire = defaultAPIKeyTimeout
	}

	fldPath := field.NewPath("")
	if err := validateAPIKeyExpire(apiKeyPass.Expire); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("expire"), apiKeyPass.Expire, err.Error()))
	}
	tenantUserSelector := fields.AndSelectors(
		fields.OneTermEqualSelector("spec.tenantID", apiKeyPass.TenantID),
		fields.OneTermEqualSelector("spec.username", apiKeyPass.Username))

	localIdentityList, err := authClient.LocalIdentities().List(metav1.ListOptions{FieldSelector: tenantUserSelector.String()})
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath.Child("username"), err))
	} else if len(localIdentityList.Items) == 0 {
		allErrs = append(allErrs, field.NotFound(fldPath.Child("username"), apiKeyPass.Username))
	} else {
		if len(localIdentityList.Items) > 1 {
			log.Warn("More than one local identity have the same name", log.String("tenantID", apiKeyPass.TenantID), log.String("userName", apiKeyPass.Username))
		}

		localIdentity := localIdentityList.Items[0]
		if err := util.VerifyDecodedPassword(apiKeyPass.Password, localIdentity.Spec.HashedPassword); err != nil {
			log.Error("Invalid password", log.ByteString("input password", []byte(apiKeyPass.Password)), log.String("store password", localIdentity.Spec.HashedPassword), log.Err(err))
			allErrs = append(allErrs, field.InternalError(fldPath.Child("password"), err))
		}
	}

	return allErrs.ToAggregate()
}
