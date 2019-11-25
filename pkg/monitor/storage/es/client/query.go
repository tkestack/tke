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

package client

import (
	"encoding/json"
	"errors"
	"strings"
	"tkestack.io/tke/pkg/util/log"
)

// MaxLimit is the size used in aggs
const MaxLimit = 10000

func commonDSL(request map[string]interface{}, orderBy string, order string, fields []string, tableFields map[string]string) (string, error) {
	for _, fieldName := range fields {
		_, ok := tableFields[fieldName]
		if !ok {
			return "", errors.New("field [" + fieldName + "] not exist")
		}
	}

	if orderBy != "" {
		_, ok := tableFields[orderBy]
		if !ok {
			return "", errors.New("orderBy field [" + orderBy + "] not exist")
		}
		request["sort"] = []map[string]interface{}{
			{orderBy: order},
		}
	}
	request["_source"] = map[string]interface{}{
		"includes": fields,
	}
	jsonData, err := json.Marshal(&request)
	if err != nil {
		return "", err
	}
	jsonStr := string(jsonData)

	return jsonStr, nil
}

func aggsDSLMonitor(request map[string]interface{}, groupBy []string, orderBy string, order string, fields []string, tableFields map[string]string) (string, []string, []string, error) {
	request["size"] = 0 // 默认会同时返回明细数据

	// 对groupBy中的某个字段排序，需要把这个字段移到slice的最右边，最后出现在最外层的aggs
	groupBy = moveToRightSpecial(groupBy, orderBy)

	// 从内往外构造嵌套json
	var columns []string
	var metrics []string
	isOrderByMetric := false
	aggs := make(map[string]interface{}, 10)
	log.Infof("fields %v", fields)
	for _, v := range fields {
		if inSliceSpecial(groupBy, v) {
			continue
		}
		aggsType, fieldName, err := CheckAggsField(v, true)
		if err != nil {
			return "", nil, nil, err
		}

		var fieldParam string
		if aggsType == "percentiles" {
			aggsType, fieldName, fieldParam, err = checkAggsFieldSpecial(v, true)
			if err != nil {
				return "", nil, nil, err
			}
		}

		/*		tableFieldType, ok := tableFields[fieldName]
				if !ok {
					return "", nil, nil, errors.New("field [" + fieldName + "] not exist")
				}
				if tableFieldType != "float" && tableFieldType != "double" && aggsType != "value_count" {
					return "", nil, nil, errors.New("field [" + fieldName + "] not support aggregation")
				}*/
		metricName := fieldName + "_" + aggsType
		metrics = append(metrics, metricName)
		columns = append(columns, metricName)
		if aggsType == "percentiles" {
			aggs[metricName] = map[string]interface{}{
				aggsType: map[string]interface{}{
					"field":    fieldName,
					"percents": fieldParam,
				},
			}
		} else {
			aggs[metricName] = map[string]interface{}{
				aggsType: map[string]interface{}{
					"field": "value",
				},
			}
		}
		if orderBy == v {
			orderBy = metricName
			isOrderByMetric = true
		} else {
			_, ok := tableFields[orderBy]
			if !ok {
				return "", nil, nil, errors.New("orderBy field [" + orderBy + "] not exist")
			}
		}
	}
	columns = reverse(columns)

	groupBy = moveToRightSpecial(groupBy, "timestamp")

	for k, fieldName := range groupBy {
		if !strings.HasPrefix(fieldName, "timestamp") {
			_, ok := tableFields[fieldName]
			if !ok {
				return "", nil, nil, errors.New("groupBy field [" + fieldName + "] not exist")
			}
		}

		columns = append(columns, fieldName)

		if fieldName == "timestamp" {
			fieldName = "timestamp(1m)" // 默认按分钟粒度聚合
		}
		if strings.HasPrefix(fieldName, "timestamp(") {
			s := strings.Split(fieldName, "(")
			ss := strings.Split(s[1], ")")
			interval := ss[0]
			if interval == "" || len(ss) != 2 {
				return "", nil, nil, errors.New("query param error: group by timestamp must specify interval, such as timestamp(1m)")
			}

			terms := map[string]interface{}{
				// todo
				"field":    "@timestamp",
				"interval": interval,
			}
			if orderBy == "timestamp" {
				terms["order"] = map[string]string{
					"_term": order,
				}
			}
			if isOrderByMetric && k == 0 {
				terms["order"] = map[string]string{
					orderBy: order,
				}
			}
			aggs = map[string]interface{}{
				"agg": map[string]interface{}{
					"date_histogram": terms,
					"aggs":           aggs,
				},
			}
		} else {
			terms := map[string]interface{}{
				// todo
				"field": "labels." + fieldName,
				// todo
				// "size":  MaxLimit,
			}
			if orderBy == fieldName {
				terms["order"] = map[string]string{
					"_term": order,
				}
			}
			if isOrderByMetric && k == 0 {
				terms["order"] = map[string]string{
					orderBy: order,
				}
			}
			aggs = map[string]interface{}{
				"agg": map[string]interface{}{
					"terms": terms,
					"aggs":  aggs,
				},
			}
		}
	}
	request["aggs"] = aggs

	columns = reverse(columns)
	jsonData, err := json.Marshal(&request)
	if err != nil {
		return "", nil, nil, err
	}
	jsonStr := string(jsonData)

	return jsonStr, columns, metrics, nil
}

func aggsDSL(request map[string]interface{}, groupBy []string, orderBy string, order string, fields []string, tableFields map[string]string) (string, []string, []string, error) {
	request["size"] = 0 // 默认会同时返回明细数据

	// 对groupBy中的某个字段排序，需要把这个字段移到slice的最右边，最后出现在最外层的aggs
	groupBy = moveToRightSpecial(groupBy, orderBy)

	// 从内往外构造嵌套json
	var columns []string
	var metrics []string
	isOrderByMetric := false
	aggs := make(map[string]interface{}, 10)
	for _, v := range fields {
		if inSliceSpecial(groupBy, v) {
			continue
		}
		aggsType, fieldName, err := CheckAggsField(v, true)
		if err != nil {
			return "", nil, nil, err
		}

		var fieldParam string
		if aggsType == "percentiles" {
			aggsType, fieldName, fieldParam, err = checkAggsFieldSpecial(v, true)
			if err != nil {
				return "", nil, nil, err
			}
		}

		tableFieldType, ok := tableFields[fieldName]
		if !ok {
			return "", nil, nil, errors.New("field [" + fieldName + "] not exist")
		}
		if tableFieldType != "float" && aggsType != "value_count" {
			return "", nil, nil, errors.New("field [" + fieldName + "] not support aggregation")
		}
		metricName := fieldName + "_" + aggsType
		metrics = append(metrics, metricName)
		columns = append(columns, metricName)
		if aggsType == "percentiles" {
			aggs[metricName] = map[string]interface{}{
				aggsType: map[string]interface{}{
					"field":    fieldName,
					"percents": fieldParam,
				},
			}
		} else {
			aggs[metricName] = map[string]interface{}{
				aggsType: map[string]interface{}{
					"field": fieldName,
				},
			}
		}
		if orderBy == v {
			orderBy = metricName
			isOrderByMetric = true
		} else {
			_, ok := tableFields[orderBy]
			if !ok {
				return "", nil, nil, errors.New("orderBy field [" + orderBy + "] not exist")
			}
		}
	}
	columns = reverse(columns)

	for k, fieldName := range groupBy {
		if !strings.HasPrefix(fieldName, "timestamp") {
			_, ok := tableFields[fieldName]
			if !ok {
				return "", nil, nil, errors.New("groupBy field [" + fieldName + "] not exist")
			}
		}

		columns = append(columns, fieldName)

		if fieldName == "timestamp" {
			fieldName = "timestamp(1m)" // 默认按分钟粒度聚合
		}
		if strings.HasPrefix(fieldName, "timestamp(") {
			s := strings.Split(fieldName, "(")
			ss := strings.Split(s[1], ")")
			interval := ss[0]
			if interval == "" || len(ss) != 2 {
				return "", nil, nil, errors.New("query param error: group by timestamp must specify interval, such as timestamp(1m)")
			}

			terms := map[string]interface{}{
				"field":    "timestamp",
				"interval": interval,
			}
			if orderBy == "timestamp" {
				terms["order"] = map[string]string{
					"_term": order,
				}
			}
			if isOrderByMetric && k == 0 {
				terms["order"] = map[string]string{
					orderBy: order,
				}
			}
			aggs = map[string]interface{}{
				"agg": map[string]interface{}{
					"date_histogram": terms,
					"aggs":           aggs,
				},
			}
		} else {
			terms := map[string]interface{}{
				"field": fieldName,
				"size":  MaxLimit,
			}
			if orderBy == fieldName {
				terms["order"] = map[string]string{
					"_term": order,
				}
			}
			if isOrderByMetric && k == 0 {
				terms["order"] = map[string]string{
					orderBy: order,
				}
			}
			aggs = map[string]interface{}{
				"agg": map[string]interface{}{
					"terms": terms,
					"aggs":  aggs,
				},
			}
		}
	}
	request["aggs"] = aggs

	columns = reverse(columns)
	jsonData, err := json.Marshal(&request)
	if err != nil {
		return "", nil, nil, err
	}
	jsonStr := string(jsonData)
	log.Info(jsonStr)

	return jsonStr, columns, metrics, nil
}

func boolDSLMonitor(startTime int64, endTime int64, conditions []interface{}, tableFields map[string]string) (map[string]interface{}, error) {
	var mustNot []map[string]interface{}

	filters := []map[string]interface{}{
		{"range": map[string]interface{}{
			// todo
			"@timestamp": map[string]int64{ // timestamp字段名固定，ES中必须为毫秒时间戳
				"gte": startTime,
				"lt":  endTime,
			},
		}},
	}

	for _, cons := range conditions {
		condition, ok := cons.([]interface{})
		if !ok {
			return nil, errors.New("param conditions error")
		}
		if len(condition) != 3 {
			return nil, errors.New("param conditions must 3")
		}
		k, ok := condition[0].(string)
		if !ok {
			return nil, errors.New("param conditions field must string")
		}
		_, ok = tableFields[k]
		if !ok {
			return nil, errors.New("conditions field [" + k + "] not exist")
		}

		operator, ok := condition[1].(string)
		if !ok {
			return nil, errors.New("param conditions operator must string")
		}
		v := condition[2]
		filter := make(map[string]interface{}, 1)

		if operator == "!=" {
			filter["term"] = map[string]interface{}{
				k: v,
			}
			mustNot = append(mustNot, filter)
			continue
		}

		if strings.EqualFold(operator, "in") {
			vArray, ok := v.([]interface{})
			if !ok {
				return nil, errors.New("param conditions \"in\" must Array")
			}
			filter["terms"] = map[string]interface{}{
				"labels." + k: vArray,
			}
		} else if operator == "=" {
			filter["term"] = map[string]interface{}{
				// todo
				"labels." + k: v,
			}
		} else {
			if !checkNum(v) {
				return nil, errors.New("param conditions value must number")
			}
			if operator == ">" {
				filter["range"] = map[string]interface{}{
					k: map[string]interface{}{
						"gt": v,
					},
				}
			} else if operator == ">=" {
				filter["range"] = map[string]interface{}{
					k: map[string]interface{}{
						"gte": v,
					},
				}
			} else if operator == "<" {
				filter["range"] = map[string]interface{}{
					k: map[string]interface{}{
						"lt": v,
					},
				}
			} else if operator == "<=" {
				filter["range"] = map[string]interface{}{
					k: map[string]interface{}{
						"lte": v,
					},
				}
			} else {
				return nil, errors.New("param conditions error: " + operator + " is invalid")
			}
		}
		filters = append(filters, filter)
	}

	boolDsl := map[string]interface{}{
		"filter": filters,
	}
	if len(mustNot) > 0 {
		boolDsl["must_not"] = mustNot
	}
	return boolDsl, nil
}

func boolDSL(startTime int64, endTime int64, conditions []interface{}, tableFields map[string]string) (map[string]interface{}, error) {
	var mustNot []map[string]interface{}

	filters := []map[string]interface{}{
		{"range": map[string]interface{}{
			"timestamp": map[string]int64{ // timestamp字段名固定，ES中必须为毫秒时间戳
				"gte": startTime,
				"lt":  endTime,
			},
		}},
	}

	for _, cons := range conditions {
		condition, ok := cons.([]interface{})
		if !ok {
			return nil, errors.New("param conditions error")
		}
		if len(condition) != 3 {
			return nil, errors.New("param conditions must 3")
		}
		k, ok := condition[0].(string)
		if !ok {
			return nil, errors.New("param conditions field must string")
		}
		_, ok = tableFields[k]
		if !ok {
			return nil, errors.New("conditions field [" + k + "] not exist")
		}

		operator, ok := condition[1].(string)
		if !ok {
			return nil, errors.New("param conditions operator must string")
		}
		v := condition[2]
		filter := make(map[string]interface{}, 1)

		if operator == "!=" {
			filter["term"] = map[string]interface{}{
				k: v,
			}
			mustNot = append(mustNot, filter)
			continue
		}

		if strings.EqualFold(operator, "in") {
			vArray, ok := v.([]interface{})
			if !ok {
				return nil, errors.New("param conditions \"in\" must Array")
			}
			filter["terms"] = map[string]interface{}{
				k: vArray,
			}
		} else if operator == "=" {
			filter["term"] = map[string]interface{}{
				k: v,
			}
		} else {
			if !checkNum(v) {
				return nil, errors.New("param conditions value must number")
			}
			if operator == ">" {
				filter["range"] = map[string]interface{}{
					k: map[string]interface{}{
						"gt": v,
					},
				}
			} else if operator == ">=" {
				filter["range"] = map[string]interface{}{
					k: map[string]interface{}{
						"gte": v,
					},
				}
			} else if operator == "<" {
				filter["range"] = map[string]interface{}{
					k: map[string]interface{}{
						"lt": v,
					},
				}
			} else if operator == "<=" {
				filter["range"] = map[string]interface{}{
					k: map[string]interface{}{
						"lte": v,
					},
				}
			} else {
				return nil, errors.New("param conditions error: " + operator + " is invalid")
			}
		}
		filters = append(filters, filter)
	}

	boolDsl := map[string]interface{}{
		"filter": filters,
	}
	if len(mustNot) > 0 {
		boolDsl["must_not"] = mustNot
	}
	return boolDsl, nil
}
