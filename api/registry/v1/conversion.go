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
		AddFieldLabelConversionsForNamespace,
		AddFieldLabelConversionsForRepository,
		AddFieldLabelConversionsForChartGroup,
		AddFieldLabelConversionsForChart,
	}
	for _, f := range funcs {
		if err := f(scheme); err != nil {
			return err
		}
	}

	return nil
}

// AddFieldLabelConversionsForNamespace adds a conversion function to convert
// field selectors of Namespace from the given version to internal version
// representation.
func AddFieldLabelConversionsForNamespace(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Namespace"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.name",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForRepository adds a conversion function to convert
// field selectors of Repository from the given version to internal version
// representation.
func AddFieldLabelConversionsForRepository(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Repository"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.name",
				"spec.namespaceName",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForChartGroup adds a conversion function to convert
// field selectors of ChartGroup from the given version to internal version
// representation.
func AddFieldLabelConversionsForChartGroup(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("ChartGroup"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.name",
				"spec.type",
				"spec.visibility",
				"repoType",  //custom label
				"projectID", //custom label
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForChart adds a conversion function to convert
// field selectors of Chart from the given version to internal version
// representation.
func AddFieldLabelConversionsForChart(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Chart"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.name",
				"spec.chartGroupName",
				"repoType",  //custom label
				"projectID", //custom label
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}
