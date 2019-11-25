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
	"errors"
	"github.com/bitly/go-simplejson"
	"tkestack.io/tke/pkg/util/log"
)

// 普通查询结果
func commonResult(respJSON *simplejson.Json, fields []string) []interface{} {
	var result []interface{}
	data := respJSON.Get("hits").Get("hits").MustArray()
	for _, v := range data {
		source := v.(map[string]interface{})
		sourceMap := source["_source"]
		columns := sourceMap.(map[string]interface{})
		var record []interface{}
		for _, key := range fields {
			record = append(record, columns[key])
		}
		result = append(result, record)
	}
	return result
}

// 聚合查询结果
func aggsResult(respJSON *simplejson.Json, metrics []string, offset int32, limit int32) []interface{} {
	bucketsData := respJSON.Get("aggregations").Get("agg").Get("buckets").MustArray()
	var result []interface{}
	recursion(bucketsData, []interface{}{}, metrics, &result)

	// 聚合查询的limit
	maxIndex := int32(len(result) - 1)
	if offset > maxIndex+1 {
		offset = maxIndex + 1
	}
	end := offset + limit
	if end > maxIndex+1 {
		end = maxIndex + 1
	}
	return result[offset:end]
}

// 递归将ES的嵌套json拼装为行数据
func recursion(buckets []interface{}, keys []interface{}, metrics []string, result *[]interface{}) {
	for _, bucket := range buckets {
		bucketMap, _ := bucket.(map[string]interface{})
		key := bucketMap["key"]
		aggData, ok := bucketMap["agg"]
		if ok {
			aggMap, ok := aggData.(map[string]interface{})
			if !ok {
				log.Warnf("aggData not map[string]interface{}, assert failed")
			}
			bucketsData := aggMap["buckets"]
			bucketsArray, ok := bucketsData.([]interface{})
			if !ok {
				log.Warnf("bucketsData not []interface{}, assert failed")
			}
			recursion(bucketsArray, append(keys, key), metrics, result)
		} else {
			record := append(keys, key)
			for _, metricName := range metrics {
				metric, ok := bucketMap[metricName].(map[string]interface{})
				if !ok {
					log.Warnf("metricData not []interface{}, assert failed")
				}
				value, ok := metric["value"]
				if !ok {
					values, ok := metric["values"]
					if ok {
						for _, v := range values.(map[string]interface{}) {
							record = append(record, v)
						}
					}
				} else {
					record = append(record, value)
				}
			}
			*result = append(*result, record)
		}
	}
}

func mappingResultMonitor(respJSON *simplejson.Json) (map[string]string, error) {
	tableFields := map[string]string{}
	for _, v := range respJSON.MustMap() {
		m1, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.New("table field not exist")
		}
		for _, v1 := range m1 {
			m2, ok := v1.(map[string]interface{})
			if !ok {
				return nil, errors.New("table field not exist")
			}
			for _, v2 := range m2 {
				m3, ok := v2.(map[string]interface{})
				if !ok {
					return nil, errors.New("table field not exist")
				}
				properties := m3["properties"]
				propertiesMap, ok := properties.(map[string]interface{})
				if !ok {
					return nil, errors.New("table field not exist")
				}
				valueMap := propertiesMap["value"].(map[string]interface{})
				tableFields["value"] = valueMap["type"].(string)

				timestampMap := propertiesMap["@timestamp"].(map[string]interface{})
				tableFields["@timestamp"] = timestampMap["type"].(string)

				labels := propertiesMap["labels"].(map[string]interface{})
				innerProperties := labels["properties"]
				innerPropertiesMap, ok := innerProperties.(map[string]interface{})

				if !ok {
					return nil, errors.New("table field innerPropertiesMap not exist")
				}
				for k3, v3 := range innerPropertiesMap {
					m4, ok := v3.(map[string]interface{})
					if !ok {
						return nil, errors.New("table field not exist")
					}
					fieldType, ok := m4["type"].(string)
					if !ok {
						return nil, errors.New("table field not exist")
					}
					tableFields[k3] = fieldType
				}
			}
		}
	}
	return tableFields, nil
}

func mappingResult(respJSON *simplejson.Json) (map[string]string, error) {
	tableFields := make(map[string]string, 4)
	for _, v := range respJSON.MustMap() {
		m1, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.New("table field not exist")
		}
		for _, v1 := range m1 {
			m2, ok := v1.(map[string]interface{})
			if !ok {
				return nil, errors.New("table field not exist")
			}
			for _, v2 := range m2 {
				m3, ok := v2.(map[string]interface{})
				if !ok {
					return nil, errors.New("table field not exist")
				}
				properties := m3["properties"]
				propertiesMap, ok := properties.(map[string]interface{})
				if !ok {
					return nil, errors.New("table field not exist")
				}
				for k3, v3 := range propertiesMap {
					m4, ok := v3.(map[string]interface{})
					if !ok {
						return nil, errors.New("table field not exist")
					}
					fieldType, ok := m4["type"].(string)
					if !ok {
						return nil, errors.New("table field not exist")
					}
					tableFields[k3] = fieldType
				}
			}
		}
	}
	return tableFields, nil
}
