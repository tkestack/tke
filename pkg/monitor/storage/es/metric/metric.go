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
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"strings"
	"time"
	"tkestack.io/tke/api/monitor"
	"tkestack.io/tke/pkg/monitor/storage/es/client"
	"tkestack.io/tke/pkg/monitor/storage/types"
	"tkestack.io/tke/pkg/util/log"
)

const (
	defaultLimit = 10000
	defaultIndex = "prometheusbeat-6.4.1"
)

func (s *ES) Query(query *monitor.MetricQuery) (*types.MetricMergedResult, error) {
	// convert query condition
	conditions := make([]interface{}, len(query.Conditions))
	for i, condition := range query.Conditions {
		tmp := make([]interface{}, 3)
		tmp[0] = condition.Key
		tmp[1] = condition.Expr
		tmp[2] = condition.Value
		if tmp[0] == "tke_cluster_instance_id" {
			tmp[0] = "cluster_id"
		}
		conditions[i] = tmp
	}

	// get start time
	startT := query.StartTime
	if startT == nil {
		return nil, errors.NewInvalid(monitor.Kind("Metric"), "query", field.ErrorList{
			field.Required(field.NewPath("query", "startTime"), "must be specify"),
		})
	}

	// get end time
	endT := query.EndTime
	if endT == nil {
		now := time.Now().Unix() * 1e3
		endT = &now
		log.Debug("the EndTime is nil")
	}

	// get table(index) without date
	table := query.Table
	if table == "" {
		table = defaultIndex
	}

	// get all indices
	indices, err := s.availableClient.Indices()
	if err != nil {
		log.Errorf("Failed to get indices", log.Err(err))
		return nil, errors.NewBadRequest(err.Error())
	}

	table = client.GetTablesMonitor(indices, table, *startT, *endT)
	if table == "" {
		log.Errorf("Get es indices error: the index is not exist")
		return nil, errors.NewBadRequest("the index is not exist")
	}
	log.Debugf("Get es table %s", table)

	// get orderBy
	var orderBy string
	if query.OrderBy == "" {
		orderBy = "@timestamp"
		log.Debug("orderBy field is nil")
	} else {
		orderBy = query.OrderBy
	}

	// get groupBy
	var groupBy []string
	if len(query.GroupBy) == 0 {
		groupBy = append(groupBy, "timestamp(60s)")
		log.Debug("groupBy field is nil")
	} else {
		groupBy = query.GroupBy
	}

	// get order
	var order string
	if query.Order == "" {
		order = "asc"
		log.Debug("order field is nil")
	} else {
		order = query.Order
	}

	// get limit
	var limit int32
	if query.Limit == 0 {
		limit = defaultLimit
		log.Debug("limit field is nil")
	} else {
		limit = query.Limit
	}

	metricResult := &types.MetricMergedResult{}

	if len(query.Fields) == 0 {
		log.Debug("es search via simpleQuery")

		result, err := s.availableClient.SimpleQuery(table)
		if err != nil {
			log.Errorf("Failed to query es database", log.Err(err))
			return nil, errors.NewBadRequest(err.Error())
		}

		log.Debugf("%v", result)
		metricResult.Data = result
		return &types.MetricMergedResult{
			Columns: []string{},
			Data:    result,
		}, nil
	}

	log.Debug("es search via Query with fields")
	var columns [][]string
	var data [][]interface{}
	for _, fld := range query.Fields {
		_, fieldName, err := client.CheckAggsField(fld, true)
		if err != nil {
			return nil, errors.NewBadRequest(err.Error())
		}
		fieldFilter := make([]interface{}, 3)
		fieldFilter[0] = "name"
		fieldFilter[1] = "="
		fieldFilter[2] = fieldName

		col, result, err := s.availableClient.QueryMonitor(table, *startT, *endT, orderBy, order,
			limit, query.Offset, append(conditions, fieldFilter), groupBy, []string{fld})
		if err != nil {
			log.Errorf("Failed to query database", log.Err(err))
			return nil, errors.NewBadRequest(err.Error())
		}
		log.Debugf("columns %s", col)
		log.Debugf("data %s", result)
		columns = append(columns, col)
		data = append(data, result)
	}
	return MergeResultES(columns, data)
}

func MergeResultES(columns [][]string, data [][]interface{}) (*types.MetricMergedResult, error) {
	colSize := len(columns) - 1 + len(columns[0])
	cols := make([]string, colSize)
	maxLen := 0
	for _, v := range data {
		if len(v) > maxLen {
			maxLen = len(v)
		}
	}
	if maxLen == 0 {
		res := &types.MetricMergedResult{
			Columns: []string{},
			Data:    []interface{}{},
		}
		return res, nil
	}
	d := make([]interface{}, maxLen)
	for i, v := range columns[0] {
		if i == 0 {
			cols[0] = v
		} else if i != len(columns[0])-1 {
			cols[len(columns)+i] = v
		}
	}
	for i := range columns {
		cols[i+1] = columns[i][len(columns[i])-1]
	}

	tagIndex := make(map[string]int)
	for i := range data {
		for j := range data[i] {
			v, ok := data[i][j].([]interface{})
			if !ok {
				return nil, errors.NewInternalError(fmt.Errorf("failed to merge data"))
			}
			timestamp, ok := v[0].(json.Number)
			if !ok {
				return nil, errors.NewInternalError(fmt.Errorf("failed to merge data"))
			}
			tagList := []string{
				timestamp.String(),
			}
			for index, tag := range v {
				if index != 0 && index != len(v)-1 {
					tagList = append(tagList, tag.(string))
				}
			}
			tagValue := strings.Join(tagList, ",")
			_, ok = tagIndex[tagValue]
			if !ok {
				tagIndex[tagValue] = len(tagIndex)
			}
			ti := tagIndex[tagValue]
			if len(d) <= ti {
				tmp := make([]interface{}, ti+1)
				copy(tmp, d)
				d = tmp
			}
			if d[ti] == nil {
				d[ti] = make([]interface{}, colSize)
				val := d[ti].([]interface{})
				for k := range v {
					if k == 0 {
						val[0] = v[k]
					} else if k != len(v)-1 {
						val[len(data)+k] = v[k]
					}
				}
			}
			val := d[ti].([]interface{})
			val[i+1] = v[len(v)-1]
		}
	}

	return &types.MetricMergedResult{
		Columns: cols,
		Data:    d,
	}, nil
}
