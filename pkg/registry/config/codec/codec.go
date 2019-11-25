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

package codec

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"tkestack.io/tke/pkg/apiserver"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	registryscheme "tkestack.io/tke/pkg/registry/apis/config/scheme"
)

// EncodeRegistryConfig encodes an internal RegistryConfiguration to an external YAML representation
func EncodeRegistryConfig(internal *registryconfig.RegistryConfiguration, targetVersion schema.GroupVersion) ([]byte, error) {
	encoder, err := NewRegistryConfigYAMLEncoder(targetVersion)
	if err != nil {
		return nil, err
	}
	// encoder will convert to external version
	data, err := runtime.Encode(encoder, internal)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// NewRegistryConfigYAMLEncoder returns an encoder that can write objects in the registryconfig API group to YAML
func NewRegistryConfigYAMLEncoder(targetVersion schema.GroupVersion) (runtime.Encoder, error) {
	_, codecs, err := registryscheme.NewSchemeAndCodecs()
	if err != nil {
		return nil, err
	}
	mediaType := "application/yaml"
	info, ok := runtime.SerializerInfoForMediaType(codecs.SupportedMediaTypes(), mediaType)
	if !ok {
		return nil, fmt.Errorf("unsupported media type %q", mediaType)
	}
	return codecs.EncoderForVersion(info.Serializer, targetVersion), nil
}

// NewYAMLEncoder generates a new runtime.Encoder that encodes objects to YAML
func NewYAMLEncoder(groupName string) (runtime.Encoder, error) {
	// encode to YAML
	mediaType := "application/yaml"
	info, ok := runtime.SerializerInfoForMediaType(apiserver.Codecs.SupportedMediaTypes(), mediaType)
	if !ok {
		return nil, fmt.Errorf("unsupported media type %q", mediaType)
	}

	versions := apiserver.Scheme.PrioritizedVersionsForGroup(groupName)
	if len(versions) == 0 {
		return nil, fmt.Errorf("no enabled versions for group %q", groupName)
	}

	// the "best" version supposedly comes first in the list returned from apiserver.Registry.EnabledVersionsForGroup
	return apiserver.Codecs.EncoderForVersion(info.Serializer, versions[0]), nil
}

// DecodeRegistryConfiguration decodes a serialized RegistryConfiguration to the internal type
func DecodeRegistryConfiguration(registryCodecs *serializer.CodecFactory, data []byte) (*registryconfig.RegistryConfiguration, error) {
	// the UniversalDecoder runs defaulting and returns the internal type by default
	obj, gvk, err := registryCodecs.UniversalDecoder().Decode(data, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decode, error: %v", err)
	}

	internalKC, ok := obj.(*registryconfig.RegistryConfiguration)
	if !ok {
		return nil, fmt.Errorf("failed to cast object to RegistryConfiguration, unexpected type: %v", gvk)
	}

	return internalKC, nil
}
