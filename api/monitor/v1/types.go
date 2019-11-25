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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=create
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Metric defines the structure for querying monitoring data requests and results.
type Metric struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// +optional
	Query MetricQuery `json:"query,omitempty" protobuf:"bytes,2,opt,name=query"`
	// +optional
	JSONResult string `json:"jsonResult,omitempty" protobuf:"bytes,3,opt,name=jsonResult"`
}

type MetricQuery struct {
	Table string `json:"table" protobuf:"bytes,1,opt,name=table"`
	// +optional
	StartTime *int64 `json:"startTime,omitempty" protobuf:"varint,2,opt,name=startTime"`
	// +optional
	EndTime    *int64                 `json:"endTime,omitempty" protobuf:"varint,3,opt,name=endTime"`
	Fields     []string               `json:"fields" protobuf:"bytes,4,rep,name=fields"`
	Conditions []MetricQueryCondition `json:"conditions" protobuf:"bytes,5,rep,name=conditions"`
	// +optional
	OrderBy string `json:"orderBy,omitempty" protobuf:"bytes,6,opt,name=orderBy"`
	// +optional
	Order   string   `json:"order,omitempty" protobuf:"bytes,7,opt,name=order"`
	GroupBy []string `json:"groupBy" protobuf:"bytes,8,rep,name=groupBy"`
	Limit   int32    `json:"limit" protobuf:"varint,9,opt,name=limit"`
	Offset  int32    `json:"offset" protobuf:"varint,10,opt,name=offset"`
}

type MetricQueryCondition struct {
	Key   string `json:"key" protobuf:"bytes,1,opt,name=key"`
	Expr  string `json:"expr" protobuf:"bytes,2,opt,name=expr"`
	Value string `json:"value" protobuf:"bytes,3,opt,name=value"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMap holds configuration data for tke to consume.
type ConfigMap struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Data contains the configuration data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// Values with non-UTF-8 byte sequences must use the BinaryData field.
	// The keys stored in Data must not overlap with the keys in
	// the BinaryData field, this is enforced during validation process.
	// +optional
	Data map[string]string `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`

	// BinaryData contains the binary data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// BinaryData can contain byte sequences that are not in the UTF-8 range.
	// The keys stored in BinaryData must not overlap with the ones in
	// the Data field, this is enforced during validation process.
	// +optional
	BinaryData map[string][]byte `json:"binaryData,omitempty" protobuf:"bytes,3,rep,name=binaryData"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMapList is a resource containing a list of ConfigMap objects.
type ConfigMapList struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of ConfigMaps.
	Items []ConfigMap `json:"items" protobuf:"bytes,2,rep,name=items"`
}
