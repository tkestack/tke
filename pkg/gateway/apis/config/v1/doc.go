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

// +k8s:deepcopy-gen=package
// +k8s:conversion-gen=tkestack.io/tke/pkg/gateway/apis/config
// +k8s:conversion-gen-external-types=tkestack.io/tke/pkg/gateway/apis/config/v1
// +k8s:defaulter-gen=TypeMeta
// +k8s:openapi-gen=true

// Package v1 is the v1 version of the API.
// +groupName=gateway.config.tkestack.io
package v1 // import "tkestack.io/tke/pkg/gateway/apis/config/v1"
