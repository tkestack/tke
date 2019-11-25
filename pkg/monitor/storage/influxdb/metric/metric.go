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
	influxclient "github.com/influxdata/influxdb1-client/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sort"
	"strings"
	"time"
	"tkestack.io/tke/api/monitor"
	"tkestack.io/tke/pkg/monitor/storage/types"
	"tkestack.io/tke/pkg/monitor/util"
	"tkestack.io/tke/pkg/util/log"
)

func (s *InfluxDB) Query(query *monitor.MetricQuery) (*types.MetricMergedResult, error) {
	metricsReq, err := s.buildMetricsRequest(query)
	if err != nil {
		return nil, err
	}

	var results []types.MetricResult
	for dbname, cmd := range metricsReq {
		log.Debugf("query db %s, cmd:%s", dbname, cmd)
		influxQuery := influxclient.NewQuery(cmd, dbname, "s")
		res, err := s.availableClient.Query(influxQuery)
		if err != nil || res.Error() != nil {
			if err != nil {
				log.Error("Failed to connect on previous influxDB client", log.Err(err))
			} else if res.Error() != nil {
				log.Error("Failed to query on previous influxDB client", log.Err(res.Error()))
			}
			// todo: need lock the available influxDB client?
			log.Info("Starting to find available influxDB client")
			// check available influxDB client. Shall we check the abstract error, such as "connection refused" or "timeout"
			for _, db := range s.clients {
				if db == s.availableClient {
					_ = s.availableClient.Close()
					continue
				}
				res, err = db.Query(influxQuery)
				if err == nil && res.Error() == nil {
					log.Warn("InfluxDB client changed")
					s.availableClient = db
					break
				}
				_ = db.Close()
			}
		}
		if err != nil {
			_ = s.availableClient.Close()
			return nil, errors.NewBadRequest(err.Error())
		} else if res.Error() != nil {
			_ = s.availableClient.Close()
			return nil, errors.NewBadRequest(res.Error().Error())
		} else {
			for _, r := range res.Results {
				if r.Series == nil || len(r.Series) == 0 {
					log.Errorf("Nil query result from db(%s): %s", dbname, cmd)
					continue
				}
				metricRes := types.MetricResult{}
				metricRes.Series = r.Series
				// TODO get max/avg value
				results = append(results, metricRes)
			}
		}
	}
	return MergeResult(results, query.GroupBy, query.Fields), nil
}

func (s *InfluxDB) buildMetricsRequest(query *monitor.MetricQuery) (map[string]string, error) {
	// get cluster name
	cluster := ""
	// get query condition
	conditionStr := ""
	if len(query.Conditions) > 0 {
		for _, v := range query.Conditions {
			if v.Key == "tke_cluster_instance_id" {
				cluster = v.Value
				continue
			}
			conditionStr += fmt.Sprintf("%s%s'%s' and ", v.Key, v.Expr, v.Value)
		}
	}
	cluster = util.RenameInfluxDB(cluster)

	// get start time
	timeRange := ""
	startT := query.StartTime
	if startT == nil {
		return nil, errors.NewInvalid(monitor.Kind("Metric"), "query", field.ErrorList{
			field.Required(field.NewPath("query", "startTime"), "must be specify"),
		})
	}
	timeRange += fmt.Sprintf("time >= '%s' ", timestampToUTCStr(startT))
	// get end time
	endT := query.EndTime
	if endT != nil {
		timeRange += fmt.Sprintf("and time < '%s' ", timestampToUTCStr(endT))
	}

	// parse group by
	byStr := strings.Join(query.GroupBy, ",")
	byStr = strings.ReplaceAll(byStr, "timestamp", "time")
	if len(byStr) == 0 {
		byStr = "time(60s)"
	}

	// map cluster => selects, batch select multi query from db just once
	metricsReq := make(map[string]string)
	for _, fieldStr := range query.Fields {
		index := strings.Index(fieldStr, "(")
		if index == -1 {
			log.Warnf("Wrong field(%s), should be expr(resource_name)", fieldStr)
			continue
		}
		expr := fieldStr[:index]
		resource := fieldStr[index+1 : len(fieldStr)-1]

		// project metrics should handle differently
		if strings.HasPrefix(resource, "project_") {
			query, err := parseProjectRequest(resource, conditionStr, timeRange, byStr)
			if err != nil {
				log.Errorf("Wrong project query: %v", err)
				continue
			}
			if _, ok := metricsReq[util.ProjectDatabaseName]; ok {
				metricsReq[util.ProjectDatabaseName] = metricsReq[util.ProjectDatabaseName] + ";" + query
			} else {
				metricsReq[util.ProjectDatabaseName] = query
			}
			continue
		}

		// struct the whole query format
		query := fmt.Sprintf("select %s(value) as %s from %s where %s %s group by %s", expr, resource, resource, conditionStr, timeRange, byStr)

		if _, ok := metricsReq[cluster]; ok {
			metricsReq[cluster] = metricsReq[cluster] + "; " + query
		} else {
			metricsReq[cluster] = query
		}
	}
	return metricsReq, nil
}

func MergeResult(results []types.MetricResult, groupByAndTimestamp []string, fields []string) *types.MetricMergedResult {
	log.Debug("Merge influxDB query result")

	if len(results) == 0 {
		return &types.MetricMergedResult{
			Columns: []string{},
			Data:    []interface{}{},
		}
	}

	maxSeriesLen := 0
	for _, result := range results {
		if maxSeriesLen < len(result.Series) {
			maxSeriesLen = len(result.Series)
		}
	}
	maxValueSize := 0
	for _, result := range results {
		for _, row := range result.Series {
			if maxValueSize < len(row.Values) {
				maxValueSize = len(row.Values)
			}
		}
	}
	var groupBy []string
	var timestamp string
	for i, group := range groupByAndTimestamp {
		if strings.Split(group, "(")[0] == "timestamp" {
			timestamp = group
			groupBy = append(groupByAndTimestamp[:i], groupByAndTimestamp[i+1:]...)
			break
		}
	}
	if len(groupByAndTimestamp) == 0 {
		timestamp = "timestamp(60s)"
	}

	renameFields := make(map[string]string)
	for _, fld := range fields {
		index := strings.Index(fld, "(")
		if index == -1 {
			continue
		}
		expr := fld[:index]
		resource := fld[index+1 : len(fld)-1]
		renameFields[resource] = resource + "_" + expr
	}

	columns := make([]string, 1+len(results)+len(groupBy))
	columns[0] = timestamp
	for i, result := range results {
		if len(result.Series) > 0 {
			columns[i+1] = renameFields[result.Series[0].Name]
		}
	}
	for i, group := range groupBy {
		columns[1+len(results)+i] = group
	}

	values := make([]interface{}, maxValueSize*maxSeriesLen)
	for i := 0; i < maxValueSize*maxSeriesLen; i++ {
		values[i] = make([]interface{}, 1+len(results)+len(groupBy))
	}

	// tagIndex shows the index of one row in series according to its tag
	tagIndex := make(map[string]int)
	for fieldIndex, result := range results {
		for _, row := range result.Series {
			var sortedKey []string
			var sortedKV []string
			for k := range row.Tags {
				sortedKey = append(sortedKey, k)
			}
			sort.Strings(sortedKey)
			for _, k := range sortedKey {
				sortedKV = append(sortedKV, fmt.Sprintf("%v:%v", k, row.Tags[k]))
			}
			tagValue := strings.Join(sortedKV, ",")
			_, ok := tagIndex[tagValue]
			if !ok {
				tagIndex[tagValue] = len(tagIndex)
			}
			ti := tagIndex[tagValue]
			for i, v := range row.Values {
				val := values[i*maxSeriesLen+ti].([]interface{})
				if val[0] == nil {
					sec := v[0].(json.Number)
					secNum, _ := sec.Int64()
					val[0] = secNum * 1000
					for groupIndex, group := range groupBy {
						val[1+len(results)+groupIndex] = row.Tags[group]
					}
				}
				val[1+fieldIndex] = v[1]
			}
		}
	}
	return &types.MetricMergedResult{
		Columns: columns,
		Data:    values,
	}
}

// parseProjectRequest parse metrics requests for project, should first sum, then select
func parseProjectRequest(resource, condition, timeRange, byStr string) (string, error) {
	byStrNew := byStr
	if !strings.Contains(byStr, "time(60s)") {
		indexHalf := strings.Index(byStr, "time(")
		indexLast := strings.Index(byStr, "s)")
		byStrNew = byStr[:indexHalf] + "time(60s)" + byStr[indexLast+2:]
	}

	// first sum, then select
	query := fmt.Sprintf("select mean(sum) as %s from (select sum(value) from %s where %s %s group by %s ) where %s group by %s",
		resource, resource, condition, timeRange, byStrNew, timeRange, byStr)

	return query, nil
}

func timestampToUTCStr(msec *int64) string {
	tt := time.Unix(*msec/1000, 0)
	return tt.UTC().String()[:19]
}
