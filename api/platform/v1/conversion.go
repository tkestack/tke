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
		AddFieldLabelConversionsForCluster,
		AddFieldLabelConversionsForClusterCredential,
		AddFieldLabelConversionsForMachine,
		AddFieldLabelConversionsForRegistry,
		AddFieldLabelConversionsForPersistentEvent,
		AddFieldLabelConversionsForHelm,
		AddFieldLabelConversionsForTappController,
		AddFieldLabelConversionsForCSIOperator,
		AddFieldLabelConversionsForVolumeDecorator,
		AddFieldLabelConversionsForLogCollector,
		AddFieldLabelConversionsForCronHPA,
		AddFieldLabelConversionsForPrometheus,
		AddFieldLabelConversionsForIPAM,
		AddFieldLabelConversionsForLBCF,
	}
	for _, f := range funcs {
		if err := f(scheme); err != nil {
			return err
		}
	}

	return nil
}

// AddFieldLabelConversionsForCluster adds a conversion function to convert
// field selectors of Cluster from the given version to internal version
// representation.
func AddFieldLabelConversionsForCluster(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Cluster"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.version",
				"spec.type",
				"status.locked",
				"status.version",
				"status.phase",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForClusterCredential adds a conversion function to convert
// field selectors of ClusterCredential from the given version to internal version
// representation.
func AddFieldLabelConversionsForClusterCredential(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("ClusterCredential"),
		func(label, value string) (string, string, error) {
			switch label {
			case "tenantID",
				"clusterName",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForMachine adds a conversion function to convert
// field selectors of Cluster from the given version to internal version
// representation.
func AddFieldLabelConversionsForMachine(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Machine"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.clusterName",
				"spec.type",
				"spec.ip",
				"status.phase",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForPersistentEvent adds a conversion function to convert
// field selectors of Project from the given version to internal version
// representation.
func AddFieldLabelConversionsForPersistentEvent(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("PersistentEvent"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.clusterName",
				"spec.version",
				"status.phase",
				"status.version",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForRegistry adds a conversion function to convert
// field selectors of Registry from the given version to internal version
// representation.
func AddFieldLabelConversionsForRegistry(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Registry"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.clusterName",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForHelm adds a conversion function to convert
// field selectors of Helm from the given version to internal version
// representation.
func AddFieldLabelConversionsForHelm(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Helm"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.clusterName",
				"spec.version",
				"status.phase",
				"status.version",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForTappController adds a conversion function to convert
// field selectors of TappController from the given version to internal version
// representation.
func AddFieldLabelConversionsForTappController(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("TappController"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.clusterName",
				"spec.version",
				"status.phase",
				"status.version",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForCSIOperator adds a conversion function to convert
// field selectors of CSIOperator from the given version to internal version
// representation.
func AddFieldLabelConversionsForCSIOperator(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("CSIOperator"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.clusterName",
				"spec.version",
				"status.phase",
				"status.version",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForVolumeDecorator adds a conversion function to convert
// field selectors of VolumeDecorator from the given version to internal version
// representation.
func AddFieldLabelConversionsForVolumeDecorator(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("VolumeDecorator"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.clusterName",
				"spec.version",
				"status.phase",
				"status.version",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForLogCollector adds a conversion function to convert
// field selectors of VolumeDecorator from the given version to internal version
// representation.
func AddFieldLabelConversionsForLogCollector(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("LogCollector"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.clusterName",
				"spec.version",
				"status.phase",
				"status.version",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForCronHPA adds a conversion function to convert
// field selectors of CronHPA from the given version to internal version
// representation.
func AddFieldLabelConversionsForCronHPA(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("CronHPA"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.clusterName",
				"spec.version",
				"status.phase",
				"status.version",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForPrometheus adds a conversion function to convert
// field selectors of Prometheus from the given version to internal version
// representation.
func AddFieldLabelConversionsForPrometheus(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Prometheus"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.clusterName",
				"spec.version",
				"status.phase",
				"status.version",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForIPAM adds a conversion function to convert
// field selectors of Project from the given version to internal version
// representation.
func AddFieldLabelConversionsForIPAM(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("IPAM"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.clusterName",
				"spec.version",
				"status.phase",
				"status.version",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForLBCF adds a conversion function to convert
// field selectors of LBCF from the given version to internal version
// representation.
func AddFieldLabelConversionsForLBCF(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("LBCF"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.clusterName",
				"spec.version",
				"status.phase",
				"status.version",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}
