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

package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/gogo/protobuf/proto"
	"github.com/influxdata/influxdb1-client/models"
	"github.com/pkg/errors"
)

const (
	equalStr           = "eq"
	equalExpr          = "=="
	greaterStr         = "gt"
	greaterExpr        = ">"
	lessStr            = "lt"
	lessExpr           = "<"
	unitKey            = "unit"
	isBoolKey          = "isBool"
	receiverGroupKey   = "receiverGroups"
	receiverKey        = "receivers"
	notifyWayKey       = "notifyWay"
	defaultLabelKey    = "alert"
	alarmPolicyNameKey = "alarmPolicyName"
	// VersionKey is a key of version
	VersionKey             = "version"
	notifySettingSep       = ","
	channelTemplateSep     = ":"
	alarmPolicyTypeCluster = "cluster"
	alarmPolicyTypeNode    = "node"
	alarmPolicyTypePod     = "pod"
	alarmPolicyTypeKey     = "alarmPolicyType"
	alarmObjectsTypeKey    = "alarmObjectsType"
	alarmObjectsTypePart   = "part"
	filterNamespaceKey     = "namespace"
	filterWorkloadKindKey  = "workload_kind"
	filterWorkloadNameKey  = "workload_name"
	measurementKey         = "measurement"
	valueKey               = "value"
	valueStr               = "{{ $value }}"
	evaluateTypeKey        = "evaluateType"
	evaluateValueKey       = "evaluateValue"
	metricDisplayNameKey   = "metricDisplayName"
)

// Response defines the structure of http response of prometheus and alertmanager
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

// Request defines the structure of http request of prometheus and alertmanager
type Request struct {
	Data string `json:"data,omitempty"`
}

// NewResult returns a response for http response of prometheus and alertmanager
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

// MetricsRequest defines the metrics request format, setting default request value
type MetricsRequest struct {
	Table      string          `json:"table"`
	StartTime  *int64          `json:"startTime"`
	EndTime    *int64          `json:"endTime"`
	Fields     []string        `json:"fields"`
	Conditions [][]interface{} `json:"conditions"`
	OrderBy    string          `json:"orderBy,omitempty"`
	Order      string          `json:"order,omitempty"`
	GroupBy    []string        `json:"groupBy"`
	Limit      int             `json:"limit,omitempty"`
	Offset     int             `json:"offset,omitempty"`
}

// MetricsResponse defines the metrics result
type MetricsResponse struct {
	Results []*MetricResult `json:"results"`
}

// MetricResult defines the single metric result
type MetricResult struct {
	MaxValue *float64     `json:"maxValue,omitempty"`
	AvgValue *float64     `json:"avgValue,omitempty"`
	Series   []models.Row `json:"series"`
}

// MergedResponse defines the merged result
type MergedResponse struct {
	Response *MergedResult `json:"response"`
}

// MergedResult defines the result of columns and data
type MergedResult struct {
	Columns []string      `json:"columns"`
	Data    []interface{} `json:"data"`
}

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
	Timeseries []TimeSeries `protobuf:"bytes,1,rep,name=timeseries,proto3" json:"timeseries"`
}

func (m *WriteRequest) Reset()         { *m = WriteRequest{} }
func (m *WriteRequest) String() string { return proto.CompactTextString(m) }
func (*WriteRequest) ProtoMessage()    {}

// MeshRequest defines request for mesh serviceCallRelation
type MeshRequest struct {
	StartTime  int64         `json:"startTime"`
	EndTime    int64         `json:"endTime"`
	MeshID     string        `json:"meshId"`
	Namespaces []interface{} `json:"namespaces"`
	TopoType   string        `json:"topoType"`
	AppID      int           `json:"AppId"`
}

// AlarmPolicy defines alarm policy info
type AlarmPolicy struct {
	AlarmPolicySettings *AlarmPolicySettings `json:"AlarmPolicySettings"`
	NotifySettings      *NotifySettings      `json:"NotifySettings"`
	Namespace           string               `json:"Namespace"`
	WorkloadType        string               `json:"WorkloadType"`
}

// AlarmPolicySettings defines alarm policy settings
type AlarmPolicySettings struct {
	AlarmPolicyName  string         `json:"AlarmPolicyName"`
	AlarmPolicyType  string         `json:"AlarmPolicyType"`
	AlarmMetrics     []*AlarmMetric `json:"AlarmMetrics"`
	AlarmObjects     string         `json:"AlarmObjects"`
	AlarmObjectsType string         `json:"AlarmObjectsType"`
	StatisticsPeriod int64          `json:"statisticsPeriod"`
}

// AlarmMetric defines alarm metric info
type AlarmMetric struct {
	Measurement       string     `json:"Measurement"`
	MetricName        string     `json:"MetricName"`
	MetricDisplayName string     `json:"MetricDisplayName"`
	ContinuePeriod    int64      `json:"ContinuePeriod"`
	Evaluator         *Evaluator `json:"Evaluator"`
	Unit              string     `json:"Unit"`
}

// Evaluator contains type and value to form expr
type Evaluator struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

// NotifyWay contains notification channels and templates
type NotifyWay struct {
	ChannelName  string `json:"ChannelName"`
	TemplateName string `json:"TemplateName"`
}

// NotifySettings contains notification info of alarm policy
type NotifySettings struct {
	ReceiverGroups []string    `json:"ReceiverGroups"`
	Receivers      []string    `json:"Receivers"`
	NotifyWay      []NotifyWay `json:"NotifyWay"`
}

// MetricFilter contains metric filter of expr
type MetricFilter struct {
	Namespace    string `json:"namespace"`
	WorkloadKind string `json:"workload_kind"`
	WorkloadName string `json:"workload_name"`
}

// AlarmPolicyPagination defines alarm policy response with pagination ability
type AlarmPolicyPagination struct {
	Page          int64         `json:"page"`
	PageSize      int64         `json:"pageSize"`
	Total         int64         `json:"total"`
	AlarmPolicies AlarmPolicies `json:"alarmPolicies"`
}

// AlarmPolicies defines slice of AlarmPolicy
type AlarmPolicies []*AlarmPolicy

// Len implements sort interface for AlarmPolicies
func (a AlarmPolicies) Len() int {
	return len(a)
}

// Less implements sort interface for AlarmPolicies
func (a AlarmPolicies) Less(i, j int) bool {
	if a[i].AlarmPolicySettings != nil && a[j].AlarmPolicySettings != nil {
		return a[i].AlarmPolicySettings.AlarmPolicyName < a[j].AlarmPolicySettings.AlarmPolicyName
	}
	return false
}

// Swap implements sort interface for AlarmPolicies
func (a AlarmPolicies) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// ToStr converts MetricFilter to string, which is used in expr
func (f *MetricFilter) ToStr() string {
	filter := fmt.Sprintf("%s=\"%s\",%s=\"%s\",%s=~\"%s\"",
		filterNamespaceKey, f.Namespace,
		filterWorkloadKindKey, f.WorkloadKind,
		filterWorkloadNameKey, f.WorkloadName)
	return filter
}

// Validate check if the policy is available
func (p *AlarmPolicy) Validate() error {
	if p.AlarmPolicySettings == nil {
		return errors.New("empty AlarmPolicySettings")
	}

	if p.NotifySettings == nil {
		return errors.New("empty NotifySettings")
	}

	if p.AlarmPolicySettings.AlarmPolicyName == "" {
		return errors.New("empty alarmPolicy name")
	}

	if p.AlarmPolicySettings.AlarmMetrics == nil {
		return errors.New("empty AlarmMetric")
	}

	if p.AlarmPolicySettings.StatisticsPeriod == 0 {
		return errors.New("zero StatisticsPeriod")
	}

	return nil
}

// GetInterval computes interval of ruleGroup from AlarmPolicySettings.StatisticsPeriod
func (p *AlarmPolicy) GetInterval() string {
	return fmt.Sprintf("%ds", p.AlarmPolicySettings.StatisticsPeriod)
}

// GetExpr builds expr of prometheus rule from AlarmMetric and AlarmPolicy
func (r *AlarmMetric) GetExpr(alarmPolicy *AlarmPolicy) string {
	var (
		op     string
		metric string
	)
	switch r.Evaluator.Type {
	case greaterStr:
		op = greaterExpr
	case equalStr:
		op = equalExpr
	case lessStr:
		op = lessExpr
	}

	switch alarmPolicy.AlarmPolicySettings.AlarmPolicyType {
	case alarmPolicyTypeCluster:
		metric = r.MetricName
	case alarmPolicyTypeNode:
		metric = r.MetricName
	case alarmPolicyTypePod:
		if alarmPolicy.AlarmPolicySettings.AlarmObjectsType == alarmObjectsTypePart {
			alarmObjects := strings.Split(alarmPolicy.AlarmPolicySettings.AlarmObjects, ",")
			filter := MetricFilter{
				Namespace:    alarmPolicy.Namespace,
				WorkloadKind: alarmPolicy.WorkloadType,
				WorkloadName: strings.Join(alarmObjects, "|"),
			}
			metric = fmt.Sprintf("%s{%s}", r.MetricName, filter.ToStr())
		} else {
			metric = r.MetricName
		}
	}
	var value string
	b, err := parseBool(r.Evaluator.Value)
	if err == nil {
		if b {
			value = "1"
		} else {
			value = "0"
		}
		op = equalExpr
	} else {
		value = r.Evaluator.Value
	}
	return fmt.Sprintf("%s %s %s", metric, op, value)
}

// GetFor computes For of prometheus rule from statisticsPeriod and ContinuePeriod
func (r *AlarmMetric) GetFor(statisticsPeriod int64) string {
	return fmt.Sprintf("%ds", statisticsPeriod*r.ContinuePeriod)
}

// GetAnnotations builds annotations of prometheus rule from AlarmMetric and AlarmPolicy
func (r *AlarmMetric) GetAnnotations(alarmPolicy *AlarmPolicy) map[string]string {
	annotations := make(map[string]string)
	notifySettings := alarmPolicy.NotifySettings
	alarmPolicySettings := alarmPolicy.AlarmPolicySettings
	var op string
	switch r.Evaluator.Type {
	case greaterStr:
		op = greaterExpr
	case equalStr:
		op = equalExpr
	case lessStr:
		op = lessExpr
	}
	v, err := parseBool(r.Evaluator.Value)
	if err == nil {
		annotations[isBoolKey] = "true"
		annotations[evaluateValueKey] = fmt.Sprintf("%t", v)
		op = equalExpr
	} else {
		annotations[isBoolKey] = "false"
		annotations[evaluateValueKey] = r.Evaluator.Value
	}
	annotations[evaluateTypeKey] = op
	annotations[unitKey] = r.Unit
	annotations[receiverGroupKey] = strings.Join(notifySettings.ReceiverGroups, notifySettingSep)
	annotations[receiverKey] = strings.Join(notifySettings.Receivers, notifySettingSep)
	var notifyWay []string
	for _, n := range notifySettings.NotifyWay {
		notifyWay = append(notifyWay, fmt.Sprintf("%s%s%s", n.ChannelName, channelTemplateSep, n.TemplateName))
	}
	annotations[notifyWayKey] = strings.Join(notifyWay, notifySettingSep)
	annotations[alarmPolicyTypeKey] = alarmPolicySettings.AlarmPolicyType
	annotations[alarmObjectsTypeKey] = alarmPolicySettings.AlarmObjectsType
	annotations[measurementKey] = r.Measurement
	annotations[valueKey] = valueStr
	annotations[metricDisplayNameKey] = r.MetricDisplayName

	return annotations
}

// GetLabels builds labels of prometheus rule from AlarmMetric
func (r *AlarmMetric) GetLabels(alarmPolicyName string, version string) map[string]string {
	labels := make(map[string]string)
	labels[defaultLabelKey] = r.MetricName
	labels[alarmPolicyNameKey] = alarmPolicyName
	labels[VersionKey] = version
	return labels
}

// NewMetricFilterFromExpr creates MetricFilter from expr
func NewMetricFilterFromExpr(expr string) *MetricFilter {
	filter := &MetricFilter{}
	metric := strings.Split(expr, " ")[0]

	if !strings.Contains(metric, "{") {
		return filter
	}
	filterStr := strings.Split(strings.Split(metric, "{")[1], "}")[0]
	for _, str := range strings.Split(filterStr, ",") {
		switch {
		case strings.HasPrefix(str, filterNamespaceKey):
			filter.Namespace = strings.Split(str, "\"")[1]
		case strings.HasPrefix(str, filterWorkloadKindKey):
			filter.WorkloadKind = strings.Split(str, "\"")[1]
		case strings.HasPrefix(str, filterWorkloadNameKey):
			filter.WorkloadName = strings.Split(str, "\"")[1]
		}
	}
	return filter
}

// NewAlarmPolicyFromRuleGroup creates AlarmPolicy from ruleGroup
func NewAlarmPolicyFromRuleGroup(ruleGroup *v1.RuleGroup) *AlarmPolicy {
	interval, _ := time.ParseDuration(ruleGroup.Interval)
	alarmPolicySettings := &AlarmPolicySettings{
		AlarmPolicyName:  ruleGroup.Name,
		AlarmMetrics:     []*AlarmMetric{},
		StatisticsPeriod: int64(interval.Seconds()),
	}

	alarmPolicy := &AlarmPolicy{
		AlarmPolicySettings: alarmPolicySettings,
	}

	for i := range ruleGroup.Rules {
		r := ruleGroup.Rules[i]
		alarmMetric := NewAlarmMetricFromRule(ruleGroup, r)
		alarmPolicy.AlarmPolicySettings.AlarmMetrics = append(alarmPolicy.AlarmPolicySettings.AlarmMetrics, alarmMetric)
		if alarmPolicy.NotifySettings == nil {
			alarmPolicy.NotifySettings = NewNotifySettingsFromRuleAnnotations(r.Annotations)
		}
		if alarmPolicy.AlarmPolicySettings.AlarmPolicyType == "" {
			if v, ok := r.Annotations[alarmPolicyTypeKey]; ok {
				alarmPolicy.AlarmPolicySettings.AlarmPolicyType = v
			}
		}
		if alarmPolicy.AlarmPolicySettings.AlarmObjectsType == "" {
			if v, ok := r.Annotations[alarmObjectsTypeKey]; ok {
				alarmPolicy.AlarmPolicySettings.AlarmObjectsType = v
			}
		}

		if alarmPolicy.AlarmPolicySettings.AlarmObjectsType == alarmObjectsTypePart &&
			(alarmPolicy.Namespace == "" || alarmPolicy.WorkloadType == "" || alarmPolicy.AlarmPolicySettings.AlarmObjects == "") {
			filter := NewMetricFilterFromExpr(r.Expr.String())
			alarmPolicy.Namespace = filter.Namespace
			alarmPolicy.WorkloadType = filter.WorkloadKind
			if filter.WorkloadName != "" {
				alarmPolicy.AlarmPolicySettings.AlarmObjects = strings.Join(strings.Split(filter.WorkloadName, "|"), ",")
			}
		}
	}

	return alarmPolicy
}

// NewNotifySettingsFromRuleAnnotations creates NotifySettings from rule annotation
func NewNotifySettingsFromRuleAnnotations(annotation map[string]string) *NotifySettings {
	ns := &NotifySettings{}
	if v, ok := annotation[receiverGroupKey]; ok {
		ns.ReceiverGroups = strings.Split(v, notifySettingSep)
	}

	if v, ok := annotation[receiverKey]; ok {
		ns.Receivers = strings.Split(v, notifySettingSep)
	}

	if v, ok := annotation[notifyWayKey]; ok {
		ns.NotifyWay = []NotifyWay{}
		for _, n := range strings.Split(v, notifySettingSep) {
			cts := strings.Split(n, channelTemplateSep)
			if len(cts) != 2 {
				continue
			}
			ns.NotifyWay = append(ns.NotifyWay, NotifyWay{ChannelName: cts[0], TemplateName: cts[1]})
		}
	}

	return ns
}

// NewAlarmMetricFromRule creates AlarmMetric from rule and ruleGroup
func NewAlarmMetricFromRule(ruleGroup *v1.RuleGroup, rule v1.Rule) *AlarmMetric {
	var interval time.Duration
	var rulefor time.Duration
	var isBool bool
	if ruleGroup.Interval == "" {
		interval = 60 * time.Second
	} else {
		interval, _ = time.ParseDuration(ruleGroup.Interval)
	}
	if v, ok := rule.Annotations[isBoolKey]; ok {
		b, err := strconv.ParseBool(v)
		if err == nil {
			isBool = b
		}
	}
	rulefor, _ = time.ParseDuration(rule.For)
	alarmMetric := &AlarmMetric{
		MetricName:     rule.Alert,
		ContinuePeriod: (rulefor / interval).Nanoseconds(),
		Evaluator:      NewEvaluatorFromExpr(rule.Expr.String(), isBool),
	}

	if v, ok := rule.Annotations[unitKey]; ok {
		alarmMetric.Unit = v
	}
	if v, ok := rule.Annotations[measurementKey]; ok {
		alarmMetric.Measurement = v
	}
	if v, ok := rule.Annotations[metricDisplayNameKey]; ok {
		alarmMetric.MetricDisplayName = v
	}
	return alarmMetric
}

// NewEvaluatorFromExpr creates Evaluator from expr
func NewEvaluatorFromExpr(expr string, isBool bool) *Evaluator {
	strs := strings.Split(expr, " ")
	var etype string
	var value string
	switch strs[1] {
	case equalExpr:
		etype = equalStr
	case greaterExpr:
		etype = greaterStr
	case lessExpr:
		etype = lessStr
	}
	if isBool {
		if strs[2] == "0" {
			value = "false"
		} else {
			value = "true"
		}
	} else {
		value = strs[2]
	}
	evaluator := &Evaluator{
		Type:  etype,
		Value: value,
	}
	return evaluator
}

func parseBool(str string) (bool, error) {
	_, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return false, errors.New("not bool")
	}

	b, err := strconv.ParseBool(str)
	if err == nil {
		return b, nil
	}

	return false, errors.New("not bool")
}
