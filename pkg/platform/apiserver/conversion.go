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
	"fmt"
	"net/url"
	"unsafe"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
)

func registerKubernetesConversion(scheme *runtime.Scheme) error {
	// Add field conversion functions.
	if err := addFieldLabelConversionsForNode(scheme); err != nil {
		return err
	}
	if err := addFieldLabelConversionsForPod(scheme); err != nil {
		return err
	}
	if err := addFieldLabelConversionsForReplicationController(scheme); err != nil {
		return err
	}
	if err := addFieldLabelConversionsForEvent(scheme); err != nil {
		return err
	}
	if err := addFieldLabelConversionsForNamespace(scheme); err != nil {
		return err
	}
	return addFieldLabelConversionsForSecret(scheme)
}

func addFieldLabelConversionsForReplicationController(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(v1.SchemeGroupVersion.WithKind("ReplicationController"),
		func(label, value string) (string, string, error) {
			switch label {
			case "metadata.name",
				"metadata.namespace",
				"status.replicas":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

func addFieldLabelConversionsForNode(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(v1.SchemeGroupVersion.WithKind("Node"),
		func(label, value string) (string, string, error) {
			switch label {
			case "metadata.name":
				return label, value, nil
			case "spec.unschedulable":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		},
	)
}

func addFieldLabelConversionsForPod(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(v1.SchemeGroupVersion.WithKind("Pod"),
		func(label, value string) (string, string, error) {
			switch label {
			case "metadata.name",
				"metadata.namespace",
				"spec.nodeName",
				"spec.restartPolicy",
				"spec.schedulerName",
				"status.phase",
				"status.podIP",
				"status.nominatedNodeName":
				return label, value, nil
				// This is for backwards compatibility with old v1 clients which send spec.host
			case "spec.host":
				return "spec.nodeName", value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		},
	)
}

func addFieldLabelConversionsForEvent(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(v1.SchemeGroupVersion.WithKind("Event"),
		func(label, value string) (string, string, error) {
			switch label {
			case "involvedObject.kind",
				"involvedObject.namespace",
				"involvedObject.name",
				"involvedObject.uid",
				"involvedObject.apiVersion",
				"involvedObject.resourceVersion",
				"involvedObject.fieldPath",
				"reason",
				"source",
				"type",
				"metadata.namespace",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

func addFieldLabelConversionsForNamespace(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(v1.SchemeGroupVersion.WithKind("Namespace"),
		func(label, value string) (string, string, error) {
			switch label {
			case "status.phase",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

func addFieldLabelConversionsForSecret(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(v1.SchemeGroupVersion.WithKind("Secret"),
		func(label, value string) (string, string, error) {
			switch label {
			case "type",
				"metadata.namespace",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

func autoConvert_url_Values_To_v1_PodExecOptions(in *url.Values, out *v1.PodExecOptions, s conversion.Scope) error {
	// WARNING: Field TypeMeta does not have json tag, skipping.

	if values, ok := map[string][]string(*in)["stdin"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_bool(&values, &out.Stdin, s); err != nil {
			return err
		}
	} else {
		out.Stdin = false
	}
	if values, ok := map[string][]string(*in)["stdout"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_bool(&values, &out.Stdout, s); err != nil {
			return err
		}
	} else {
		out.Stdout = false
	}
	if values, ok := map[string][]string(*in)["stderr"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_bool(&values, &out.Stderr, s); err != nil {
			return err
		}
	} else {
		out.Stderr = false
	}
	if values, ok := map[string][]string(*in)["tty"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_bool(&values, &out.TTY, s); err != nil {
			return err
		}
	} else {
		out.TTY = false
	}
	if values, ok := map[string][]string(*in)["container"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_string(&values, &out.Container, s); err != nil {
			return err
		}
	} else {
		out.Container = ""
	}
	if values, ok := map[string][]string(*in)["command"]; ok && len(values) > 0 {
		out.Command = *(*[]string)(unsafe.Pointer(&values))
	} else {
		out.Command = nil
	}
	return nil
}

// Convert_url_Values_To_v1_PodExecOptions is an autogenerated conversion function.
func Convert_url_Values_To_v1_PodExecOptions(in *url.Values, out *v1.PodExecOptions, s conversion.Scope) error {
	return autoConvert_url_Values_To_v1_PodExecOptions(in, out, s)
}

func autoConvert_url_Values_To_v1_PodLogOptions(in *url.Values, out *v1.PodLogOptions, s conversion.Scope) error {
	// WARNING: Field TypeMeta does not have json tag, skipping.

	if values, ok := map[string][]string(*in)["container"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_string(&values, &out.Container, s); err != nil {
			return err
		}
	} else {
		out.Container = ""
	}
	if values, ok := map[string][]string(*in)["follow"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_bool(&values, &out.Follow, s); err != nil {
			return err
		}
	} else {
		out.Follow = false
	}
	if values, ok := map[string][]string(*in)["previous"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_bool(&values, &out.Previous, s); err != nil {
			return err
		}
	} else {
		out.Previous = false
	}
	if values, ok := map[string][]string(*in)["sinceSeconds"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_Pointer_int64(&values, &out.SinceSeconds, s); err != nil {
			return err
		}
	} else {
		out.SinceSeconds = nil
	}
	if values, ok := map[string][]string(*in)["sinceTime"]; ok && len(values) > 0 {
		if err := metav1.Convert_Slice_string_To_Pointer_v1_Time(&values, &out.SinceTime, s); err != nil {
			return err
		}
	} else {
		out.SinceTime = nil
	}
	if values, ok := map[string][]string(*in)["timestamps"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_bool(&values, &out.Timestamps, s); err != nil {
			return err
		}
	} else {
		out.Timestamps = false
	}
	if values, ok := map[string][]string(*in)["tailLines"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_Pointer_int64(&values, &out.TailLines, s); err != nil {
			return err
		}
	} else {
		out.TailLines = nil
	}
	if values, ok := map[string][]string(*in)["limitBytes"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_Pointer_int64(&values, &out.LimitBytes, s); err != nil {
			return err
		}
	} else {
		out.LimitBytes = nil
	}
	if values, ok := map[string][]string(*in)["insecureSkipTLSVerifyBackend"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_bool(&values, &out.InsecureSkipTLSVerifyBackend, s); err != nil {
			return err
		}
	} else {
		out.InsecureSkipTLSVerifyBackend = false
	}
	return nil
}

// Convert_url_Values_To_v1_PodLogOptions is an autogenerated conversion function.
func Convert_url_Values_To_v1_PodLogOptions(in *url.Values, out *v1.PodLogOptions, s conversion.Scope) error {
	return autoConvert_url_Values_To_v1_PodLogOptions(in, out, s)
}
