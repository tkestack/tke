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

package v1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
)

func addConversionFuncs(scheme *runtime.Scheme) error {
	funcs := []func(scheme *runtime.Scheme) error{
		AddFieldLabelConversionsForLocalIdentity,
		AddFieldLabelConversionsForAPIKey,
		AddFieldLabelConversionsForPolicy,
		AddFieldLabelConversionsForRule,
		AddFieldLabelConversionsForCategory,
		AddFieldLabelConversionsForLocalGroup,
		AddFieldLabelConversionsForRole,
		AddFieldLabelConversionsForUser,
		AddFieldLabelConversionsForGroup,
		AddFieldLabelConversionsForIdentityProvider,
	}
	for _, f := range funcs {
		if err := f(scheme); err != nil {
			return err
		}
	}

	return nil
}

// AddFieldLabelConversionsForLocalIdentity adds a conversion function to convert
// field selectors of LocalIdentify from the given version to internal version
// representation.
func AddFieldLabelConversionsForLocalIdentity(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("LocalIdentity"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.username",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForAPIKey adds a conversion function to convert
// field selectors of APIKey from the given version to internal version
// representation.
func AddFieldLabelConversionsForAPIKey(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("APIKey"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.apiKey",
				"spec.username",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForPolicy adds a conversion function to convert
// field selectors of Policy from the given version to internal version
// representation.
func AddFieldLabelConversionsForPolicy(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Policy"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.username",
				"spec.category",
				"spec.displayName",
				"spec.type",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForRule adds a conversion function to convert
// field selectors of Rule from the given version to internal version
// representation.
func AddFieldLabelConversionsForRule(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Rule"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.ptype",
				"spec.v0",
				"spec.v1",
				"spec.v2",
				"spec.v3",
				"spec.v4",
				"spec.v5",
				"spec.v6",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForCategory adds a conversion function to convert
// field selectors of Category from the given version to internal version
// representation.
func AddFieldLabelConversionsForCategory(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Category"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.username",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForIdentityProvider adds a conversion function to convert
// field selectors of IdentityProvider from the given version to internal version
// representation.
func AddFieldLabelConversionsForIdentityProvider(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("IdentityProvider"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.type",
				"spec.name",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForLocalGroup adds a conversion function to convert
// field selectors of LocalGroup from the given version to internal version
// representation.
func AddFieldLabelConversionsForLocalGroup(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("LocalGroup"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.displayName",
				"spec.tenantID",
				"spec.username",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForRole adds a conversion function to convert
// field selectors of Role from the given version to internal version
// representation.
func AddFieldLabelConversionsForRole(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Role"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.displayName",
				"spec.tenantID",
				"spec.username",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForUser adds a conversion function to convert
// field selectors of User from the given version to internal version
// representation.
func AddFieldLabelConversionsForUser(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("User"),
		func(label, value string) (string, string, error) {
			switch label {
			case "keyword",
				"spec.tenantID":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForGroup adds a conversion function to convert
// field selectors of IdentityProvider from the given version to internal version
// representation.
func AddFieldLabelConversionsForGroup(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Group"),
		func(label, value string) (string, string, error) {
			switch label {
			case "keyword",
				"spec.tenantID":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}
