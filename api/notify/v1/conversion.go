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
		AddFieldLabelConversionsForChannel,
		AddFieldLabelConversionsForTemplate,
		AddFieldLabelConversionsForReceiver,
		AddFieldLabelConversionsForReceiverGroup,
		AddFieldLabelConversionsForMessageRequest,
		AddFieldLabelConversionsForMessage,
	}
	for _, f := range funcs {
		if err := f(scheme); err != nil {
			return err
		}
	}

	return nil
}

// AddFieldLabelConversionsForChannel adds a conversion function to convert
// field selectors of Channel from the given version to internal version
// representation.
func AddFieldLabelConversionsForChannel(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Channel"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"statuts.phase",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForTemplate adds a conversion function to convert
// field selectors of Template from the given version to internal version
// representation.
func AddFieldLabelConversionsForTemplate(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Template"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForReceiver adds a conversion function to convert
// field selectors of Receiver from the given version to internal version
// representation.
func AddFieldLabelConversionsForReceiver(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Receiver"),
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

// AddFieldLabelConversionsForReceiverGroup adds a conversion function to convert
// field selectors of ReceiverGroup from the given version to internal version
// representation.
func AddFieldLabelConversionsForReceiverGroup(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("ReceiverGroup"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForMessageRequest adds a conversion function to convert
// field selectors of MessageRequest from the given version to internal version
// representation.
func AddFieldLabelConversionsForMessageRequest(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("MessageRequest"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"status.phase",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}

// AddFieldLabelConversionsForMessage adds a conversion function to convert
// field selectors of Message from the given version to internal version
// representation.
func AddFieldLabelConversionsForMessage(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("Message"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.tenantID",
				"spec.receiverName",
				"spec.username",
				"spec.channelMessageID",
				"status.phase",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}
