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
	"fmt"
	"strconv"
	"strings"
	"time"
	"tkestack.io/tke/pkg/util/log"
)

var tableNameFormat = "2006.01.02"
var aggsTypes = []string{"sum", "avg", "mean", "max", "min", "count", "percentile"}
var aggsTypeMap = map[string]string{
	"mean":       "avg",
	"count":      "value_count",
	"percentile": "percentiles",
}

// 将字符串s移动到切片a的最右边
func moveToRightSpecial(a []string, s string) []string {
	n := -1
	for k, v := range a {
		if strings.HasPrefix(v, "timestamp(") {
			v = "timestamp"
		}
		if v == s {
			n = k
			break
		}
	}
	if n != -1 {
		tmp := a[n]
		a = append(a[:n], a[n+1:]...)
		a = append(a, tmp)
	}
	return a
}

// 判断一个元素是否在slice中
func inSlice(s []string, item string) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

// 判断一个元素是否在slice中
func inSliceSpecial(s []string, item string) bool {
	for _, v := range s {
		if strings.HasPrefix(v, "timestamp(") {
			v = "timestamp"
		}
		if v == item {
			return true
		}
	}
	return false
}

// 将一个slice反序
func reverse(s []string) []string {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// 检查一个interface{}类型是否为数字
func checkNum(n interface{}) bool {
	switch n := n.(type) {
	case int, int64, uint, uint64, float64, float32, json.Number:
		return true
	case string: // ES的DSL中支持传入带双引号的数字
		_, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return false
		}
		return true
	default:
		return false
	}
}

// CheckInterfaceString 检查一个[]interface{}中的元素是否都为string
func CheckInterfaceString(a []interface{}) ([]string, error) {
	var new []string
	for _, v := range a {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%+v not string", v)
		}
		new = append(new, s)
	}
	return new, nil
}

// 检查聚合字段格式是否正确，如sum(success)
func CheckAggsField(field string, needReplace bool) (string, string, error) {
	s := strings.Split(field, "(")
	if len(s) != 2 {
		return "", "", errors.New("param fields error: [" + field + "] is invalid, must be aggs(field) form")
	}
	aggsType := s[0]
	if !inSlice(aggsTypes, aggsType) {
		return "", "", errors.New("param fields error: [" + aggsType + "] is invalid, aggs type not support")
	}

	s = strings.Split(s[1], ")")
	if len(s) != 2 {
		return "", "", errors.New("param fields error: [" + field + "] is invalid")
	}
	fieldName := s[0]

	if needReplace {
		for k, v := range aggsTypeMap {
			if aggsType == k {
				aggsType = v
			}
		}
	}
	return aggsType, fieldName, nil
}

// 检查聚合字段格式是否正确，如percentile(success,50)
func checkAggsFieldSpecial(field string, needReplace bool) (string, string, string, error) {
	s := strings.Split(field, "(")
	if len(s) != 2 {
		return "", "", "", errors.New("param fields error: [" + field + "] is invalid")
	}
	aggsType := s[0]
	if !inSlice(aggsTypes, aggsType) {
		return "", "", "", errors.New("param fields error: [" + aggsType + "] is invalid")
	}

	s = strings.Split(s[1], ")")
	if len(s) != 2 {
		return "", "", "", errors.New("param fields error: [" + field + "] is invalid")
	}

	s = strings.Split(s[0], ",")
	fieldName := s[0]
	fieldParam := s[1]
	if needReplace {
		for k, v := range aggsTypeMap {
			if aggsType == k {
				aggsType = v
			}
		}
	}
	return aggsType, fieldName, fieldParam, nil
}

// 将结果字段名中的value_count部分转回count
func replaceAggsField(fields []string, columns []string) []string {
	for _, v := range fields {
		var aggsType, fieldName string
		if strings.HasPrefix(v, "percentile(") {
			aggsType, fieldName, _, _ = checkAggsFieldSpecial(v, false)
		} else {
			aggsType, fieldName, _ = CheckAggsField(v, false)
		}
		_, ok := aggsTypeMap[aggsType]
		if ok {
			newName := fieldName + "_" + aggsType
			var aggsType2, fieldName2 string
			if aggsType == "percentile" {
				aggsType2, fieldName2, _, _ = checkAggsFieldSpecial(v, true)
			} else {
				aggsType2, fieldName2, _ = CheckAggsField(v, true)
			}
			oldName := fieldName2 + "_" + aggsType2
			replaceItem(&columns, oldName, newName)
		}
	}
	return columns
}

// 将一个slice中值为old的元素替换为new
func replaceItem(a *[]string, old string, new string) {
	for k, v := range *a {
		if v == old {
			(*a)[k] = new
		}
	}
}

// GetTablesMonitor get some indices according timestamp
func GetTablesMonitor(indices []string, tableName string, startTime int64, endTime int64) string {
	var sb strings.Builder
	startTimeObj := time.Unix(startTime/1000, 0)
	startTime -= int64(startTimeObj.Second()+startTimeObj.Minute()*60+startTimeObj.Hour()*3600) * 1000
	for i := startTime; i <= endTime; i += 86400000 {
		timeObj := time.Unix(i/1000, 0)
		date := timeObj.Format(tableNameFormat)
		indexName := tableName + "-" + date
		log.Infof("indices:%v indexName:%s", indices, indexName)
		if !inSlice(indices, indexName) {
			continue
		}
		sb.WriteString(indexName)
		sb.WriteString(",")
	}
	tables := strings.TrimSuffix(sb.String(), ",")
	return tables
}

// GetTables get some indices according timestamp
func GetTables(indices []string, tableName string, startTime int64, endTime int64) string {
	var sb strings.Builder
	startTimeObj := time.Unix(startTime/1000, 0)
	startTime -= int64(startTimeObj.Second()+startTimeObj.Minute()*60+startTimeObj.Hour()*3600) * 1000
	for i := startTime; i <= endTime; i += 86400000 {
		timeObj := time.Unix(i/1000, 0)
		date := timeObj.Format("20060102")
		indexName := tableName + "_" + date
		if !inSlice(indices, indexName) {
			continue
		}
		sb.WriteString(indexName)
		sb.WriteString(",")
	}
	tables := strings.TrimSuffix(sb.String(), ",")
	return tables
}
