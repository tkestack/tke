/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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
 *
 */

package rest

import (
	"encoding/json"
	"io"

	istionetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	autoscaling "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Response defines the structure of http response of
type Response struct {
	Result bool        `json:"result"`
	Err    string      `json:"err,omitempty"`
	Rev    int         `json:"rev,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

// ResponseForTest leaves data as json.RawMessage to unmarshal to struct we want, just for unit test
type ResponseForTest struct {
	Result bool            `json:"result"`
	Err    string          `json:"err,omitempty"`
	Rev    int             `json:"rev,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

// Request defines the structure of http request
type Request struct {
	Data string `json:"data,omitempty"`
}

// NewResult returns a response for http response
func NewResult(result bool, errMsg string) *Response {
	return &Response{
		Result: result,
		Err:    errMsg,
	}
}

// Decode decodes the reader content to response object
func (r *Response) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(r)
}

// Decode decodes the reader content to request object
func (r *Request) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(r)
}

// Decode decodes the reader content to ResponseForTest object
func (r *ResponseForTest) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(r)
}

type MetricRequest struct {
	Table      string      `json:"table"`
	StartTime  int64       `json:"startTime"`
	EndTime    int64       `json:"endTime"`
	Fields     []string    `json:"fields"`
	Conditions []Condition `json:"conditions"`
	OrderBy    string      `json:"orderBy"`
	GroupBy    []string    `json:"groupBy"`
	Order      string      `json:"order"`
	Limit      int         `json:"limit"`
}

type Condition []interface{}

type MetricData struct {
	Columns []string        `json:"columns"`
	Data    [][]interface{} `json:"data"`
}

type MetricQuery struct {
	Table string
	// +optional
	StartTime *int64
	// +optional
	EndTime    *int64
	Fields     []string
	Conditions []MetricQueryCondition
	// +optional
	OrderBy string
	// +optional
	Order   string
	GroupBy []string
	Limit   int32
	Offset  int32
}

type MetricQueryCondition struct {
	Key   string
	Expr  string
	Value interface{}
}

type TopoQuery struct {
	AppID       string   `json:"appId"`
	StartTime   *int64   `json:"startTime"`
	EndTime     *int64   `json:"endTime"`
	MeshID      string   `json:"meshId"`
	TopoType    string   `json:"topoType"`
	Namespaces  []string `json:"namespaces"`
	MeshVersion string   `json:"meshVersion"`
	App         string   `json:"app"`
}

type TopoData struct {
	Edges []Edge          `json:"edges"`
	Nodes map[string]Node `json:"nodes"`
}

type Edge struct {
	Did  string `json:"did"`
	Sid  string `json:"sid"`
	Type string `json:"type"`
}

type Node struct {
	Name           string      `json:"name"`
	Type           string      `json:"type"`
	ServiceNodeID  string      `json:"serviceNodeId,omitempty"`
	WorkloadNodeID string      `json:"workloadNodeId,omitempty"`
	HTTPMetric     *HTTPMetric `json:"http_metric,omitempty"`
	TCPMetric      *TCPMetric  `json:"tcp_metric,omitempty"`
}

type HTTPMetric struct {
	Count        float32 `json:"count"`
	DurationAvg  float32 `json:"durationAvg"`
	Rps          float32 `json:"rps"`
	SuccessCount float32 `json:"successCount"`
	FailedCount  float32 `json:"failedCount"`
}

type TCPMetric struct {
	ConnectionReceivedBytesTotal float32 `json:"connectionReceivedBytesTotal"`
	ConnectionSentBytesTotal     float32 `json:"connectionSentBytesTotal"`
	// Count                        float32 `json:"count"`
}

type MicroService struct {
	Cluster          Cluster                           `json:"cluster"`
	Service          corev1.Service                    `json:"service"`
	Endpoints        corev1.Endpoints                  `json:"endpoints"`
	Workloads        []Workload                        `json:"workloads"`
	VirtualServices  []istionetworking.VirtualService  `json:"virtualServices"`
	DestinationRules []istionetworking.DestinationRule `json:"destinationRules"`
	Pods             []corev1.Pod                      `json:"pods"`
}

type Cluster struct {
	MeshName    string `json:"mesh"`
	ClusterName string `json:"cluster"`
	Role        string `json:"role,omitempty"`
	Region      string `json:"region,omitempty"`
}

type Workload struct {
	runtime.Object `json:",inline"`
	HPA            autoscaling.HorizontalPodAutoscaler `json:"hpa,omitempty"`
}
