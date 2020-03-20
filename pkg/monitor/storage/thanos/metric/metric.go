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

package metric

import (
	"context"
	"fmt"
	"strings"
	"time"
	"tkestack.io/tke/api/monitor"
	"tkestack.io/tke/pkg/monitor/storage/types"
	"tkestack.io/tke/pkg/util/log"

	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func (s *Thanos) Query(query *monitor.MetricQuery) (*types.MetricMergedResult, error) {
	// get start time
	if query.StartTime == nil {
		return nil, errors.NewInvalid(monitor.Kind("Metric"), "query", field.ErrorList{
			field.Required(field.NewPath("query", "startTime"), "must be specify"),
		})
	}
	startT := *query.StartTime / 1000

	// get end time
	if query.EndTime == nil {
		return nil, errors.NewInvalid(monitor.Kind("Metric"), "query", field.ErrorList{
			field.Required(field.NewPath("query", "endTime"), "must be specify"),
		})
	}
	endT := *query.EndTime / 1000

	var (
		timestamp               string
		step                    time.Duration
		err                     error
		groupByWithoutTimestamp []string
	)
	for i, group := range query.GroupBy {
		strs := strings.Split(group, "(")
		if strs[0] == "timestamp" && len(strs) >= 2 {
			timestamp = group
			step, err = time.ParseDuration(strings.Split(strs[1], ")")[0])
			if err != nil {
				return nil, errors.NewInvalid(monitor.Kind("Metric"), "query", field.ErrorList{
					field.Required(field.NewPath("query", "groupBy"), "invalid timestamp"),
				})
			}
			groupByWithoutTimestamp = append(query.GroupBy[:i], query.GroupBy[i+1:]...)
			break
		}
	}
	if timestamp == "" {
		timestamp = "timestamp(60s)"
		step, _ = time.ParseDuration("60s")
	}
	r := v1.Range{
		Start: time.Unix(startT, 0),
		End:   time.Unix(endT, 0),
		Step:  step,
	}

	v1api := v1.NewAPI(s.availableClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	promQuery, err := s.buildPromQuery(query)
	if err != nil {
		return nil, err
	}

	var results model.Matrix
	res, warnings, err := v1api.QueryRange(ctx, promQuery, r)
	if err != nil {
		log.Errorf("Error querying Prometheus: %v", err)
		return nil, errors.NewBadRequest(err.Error())
	}
	if len(warnings) > 0 {
		log.Warnf("Warnings: %v", warnings)
	}

	switch res.Type() {
	case model.ValMatrix:
		v, _ := res.(model.Matrix)
		results = v
	default:
		log.Errorf("unexpected res type: %v", res.Type())
		return nil, errors.NewBadRequest("unexpected res type")
	}

	return MergeResult(results, timestamp, groupByWithoutTimestamp, startT, endT, int64(step.Seconds()), query.Fields), nil
}

// assemble prometheus query string like {__name__=~"aaa|bbb|ccc",key1="value1",key2="value2"}
func (s *Thanos) buildPromQuery(query *monitor.MetricQuery) (string, error) {
	// format query condition string like "key1=\"value1\",key2=\"value2\""
	conditionStr := ""
	if len(query.Conditions) > 0 {
		for i, v := range query.Conditions {
			if v.Key == "tke_cluster_instance_id" {
				v.Key = "cluster_id"
			}
			conditionStr += fmt.Sprintf("%s%s\"%s\"", v.Key, v.Expr, v.Value)
			if i != len(query.Conditions)-1 {
				conditionStr += ","
			}
		}
	}

	// format query metric names like "\"aaa|bbb|ccc\""
	names := ""
	for i, fieldStr := range query.Fields {
		index := strings.Index(fieldStr, "(")
		if index == -1 {
			log.Warnf("Wrong field(%s), should be expr(resource_name)", fieldStr)
			continue
		}
		resource := fieldStr[index+1 : len(fieldStr)-1]
		names += resource
		if i != len(query.Fields)-1 {
			names += "|"
		}
	}
	names = fmt.Sprintf("\"%s\"", names)

	promQuery := ""
	if conditionStr != "" {
		promQuery = fmt.Sprintf("{__name__=~%s,%s}", names, conditionStr)
	} else {
		promQuery = fmt.Sprintf("{__name__=~%s}", names)
	}
	return promQuery, nil
}

func MergeResult(results model.Matrix, timestamp string, groupByWithoutTimestamp []string, startT int64, endT int64, step int64, fields []string) *types.MetricMergedResult {
	resultsLen := len(results)
	if resultsLen == 0 {
		return &types.MetricMergedResult{
			Columns: []string{},
			Data:    []interface{}{},
		}
	}

	tagIndexes := make(map[string]int)
	for _, result := range results {
		tags := []string{}
		for _, groupBy := range groupByWithoutTimestamp {
			if v, ok := result.Metric[model.LabelName(groupBy)]; ok {
				tags = append(tags, string(v))
			}
		}

		tagStr := strings.Join(tags, ",")
		_, ok := tagIndexes[tagStr]
		if !ok {
			tagIndexes[tagStr] = len(tagIndexes)
		}
	}
	size := ((endT-startT)/step + 1) * int64(len(tagIndexes))

	columns := []string{}
	columns = append(columns, timestamp)
	fieldIndexes := make(map[string]int)
	for i, fld := range fields {
		index := strings.Index(fld, "(")
		if index == -1 {
			continue
		}
		expr := fld[:index]
		resource := fld[index+1 : len(fld)-1]
		columns = append(columns, resource+"_"+expr)
		fieldIndexes[resource] = i
	}
	columns = append(columns, groupByWithoutTimestamp...)

	values := make([]interface{}, size)
	for i := int64(0); i < size; i++ {
		values[i] = make([]interface{}, len(columns))
	}

	for _, result := range results {
		metricName := string(result.Metric["__name__"])
		tags := []string{}
		for _, groupBy := range groupByWithoutTimestamp {
			if v, ok := result.Metric[model.LabelName(groupBy)]; ok {
				tags = append(tags, string(v))
			}
		}

		tagStr := strings.Join(tags, ",")
		tagIndex := tagIndexes[tagStr]
		fieldIndex := fieldIndexes[metricName]
		for _, sp := range result.Values {
			index := ((sp.Timestamp.Unix()-startT)/step)*int64(len(tagIndexes)) + int64(tagIndex)
			lengthNow := len(values)
			if int64(lengthNow) <= index {
				tmp := make([]interface{}, index+1)
				copy(tmp, values)
				values = tmp
				for i := lengthNow; int64(i) <= index; i++ {
					values[i] = make([]interface{}, len(columns))
				}
			}
			values[index].([]interface{})[0] = sp.Timestamp.Unix() * 1000
			values[index].([]interface{})[fieldIndex+1] = float64(sp.Value)
			for j, tag := range tags {
				values[index].([]interface{})[len(fieldIndexes)+1+j] = tag
			}
		}
	}

	// get rid of empty values as samples may not fulfill all the time gap
	for i, v := range values {
		if v.([]interface{})[0] == nil {
			continue
		} else {
			values = values[i:]
			break
		}
	}

	log.Debugf("columes %v", columns)
	log.Debugf("values %v", values)
	return &types.MetricMergedResult{
		Columns: columns,
		Data:    values,
	}
}
