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

package project

import "github.com/gogo/protobuf/proto"

// Sample defines one metrics sample
type Sample struct {
	Value     float64 `protobuf:"fixed64,1,opt,name=value,proto3" json:"value,omitempty"`
	Timestamp int64   `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

// TimeSeries defines time series metrics that prepared to be send
type TimeSeries struct {
	Labels  []Label  `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels"`
	Samples []Sample `protobuf:"bytes,2,rep,name=samples,proto3" json:"samples"`
}

// Label defines label of time series metrics
type Label struct {
	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

// WriteRequest defines request for prom-beat
type WriteRequest struct {
	TimeSeries []TimeSeries `protobuf:"bytes,1,rep,name=timeseries,proto3" json:"timeseries"`
}

func (m *WriteRequest) Reset()         { *m = WriteRequest{} }
func (m *WriteRequest) String() string { return proto.CompactTextString(m) }
func (*WriteRequest) ProtoMessage()    {}
