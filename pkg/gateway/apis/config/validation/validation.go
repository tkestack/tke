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

package validation

import (
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
)

// ValidateGatewayConfiguration validates `gc` and returns an error if it is invalid
func ValidateGatewayConfiguration(gc *gatewayconfig.GatewayConfiguration) error {
	var allErrors []error

	fld := field.NewPath("components")
	if gc.Components.Platform != nil {
		allErrors = append(allErrors, validateComponent(gc.Components.Platform, fld.Child("platform"))...)
	}

	if gc.Components.Business != nil {
		allErrors = append(allErrors, validateComponent(gc.Components.Business, fld.Child("business"))...)
	}

	if gc.Components.Notify != nil {
		allErrors = append(allErrors, validateComponent(gc.Components.Notify, fld.Child("notify"))...)
	}

	if gc.Components.Monitor != nil {
		allErrors = append(allErrors, validateComponent(gc.Components.Monitor, fld.Child("monitor"))...)
	}

	if gc.Components.Auth != nil {
		allErrors = append(allErrors, validateComponent(gc.Components.Auth, fld.Child("auth"))...)
	}

	if gc.Components.Registry != nil {
		allErrors = append(allErrors, validateComponent(gc.Components.Registry, fld.Child("registry"))...)
	}

	return utilerrors.NewAggregate(allErrors)
}

func validateComponent(c *gatewayconfig.Component, fld *field.Path) []error {
	var allErrors []error

	if c.Address == "" {
		allErrors = append(allErrors, field.Required(fld.Child("address"), "must be specify"))
	}

	authMode := 0
	if c.FrontProxy != nil {
		authMode++
		subFld := fld.Child("frontProxy")
		if c.FrontProxy.ClientKeyFile == "" {
			allErrors = append(allErrors, field.Required(subFld.Child("clientKeyFile"), "must be specify"))
		}
		if c.FrontProxy.ClientCertFile == "" {
			allErrors = append(allErrors, field.Required(subFld.Child("clientCertFile"), "must be specify"))
		}
		if c.FrontProxy.UsernameHeader == "" {
			allErrors = append(allErrors, field.Required(subFld.Child("usernameHeader"), "must be specify"))
		}
		if c.FrontProxy.GroupsHeader == "" {
			allErrors = append(allErrors, field.Required(subFld.Child("groupsHeader"), "must be specify"))
		}
		if c.FrontProxy.ExtraPrefixHeader == "" {
			allErrors = append(allErrors, field.Required(subFld.Child("extraPrefixHeader"), "must be specify"))
		}
	}

	if c.Passthrough != nil {
		authMode++
	}

	if authMode == 0 {
		allErrors = append(allErrors, field.Required(fld, "at least 1 authentication mode is required"))
	} else if authMode > 1 {
		allErrors = append(allErrors, field.Required(fld, "authentication mode can only specify at most one"))
	}

	return allErrors
}
