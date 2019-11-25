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
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"net/http"
	"strings"
	"time"
	"tkestack.io/tke/pkg/util/log"
)

// Client is a ElasticSearch Client
type Client struct {
	URL      string
	Username string
	Password string
}

func (c *Client) QueryMonitor(table string, startTime int64, endTime int64, orderBy string, order string, limit int32,
	offset int32, conditions []interface{}, groupBy []string, fields []string) ([]string, []interface{}, error) {
	mapping, err := c.mapping(table)
	if err != nil {
		return nil, nil, err
	}
	tableFields, err := mappingResultMonitor(mapping)
	if err != nil {
		return nil, nil, err
	}

	request := map[string]interface{}{
		"from": offset,
		"size": limit,
	}
	booDsl, err := boolDSLMonitor(startTime, endTime, conditions, tableFields)
	if err != nil {
		return nil, nil, err
	}
	request["query"] = map[string]interface{}{
		"bool": booDsl,
	}

	if len(groupBy) > 0 {
		dsl, columns, metrics, err := aggsDSLMonitor(request, groupBy, orderBy, order, fields, tableFields)
		log.Debugf("dsl %s \n columns %s \n metrics %s", dsl, columns, metrics)
		if err != nil {
			return nil, nil, err
		}
		respJSON, err := c.search(table, dsl)
		if err != nil {
			return nil, nil, err
		}
		result := aggsResult(respJSON, metrics, offset, limit)
		columns = replaceAggsField(fields, columns)

		return columns, result, nil
	}
	dsl, err := commonDSL(request, orderBy, order, fields, tableFields)
	if err != nil {
		return nil, nil, err
	}
	log.Debugf("commonDSL: %s", dsl)
	respJSON, err := c.search(table, dsl)
	if err != nil {
		return nil, nil, err
	}
	result := commonResult(respJSON, fields)
	columns := fields
	return columns, result, nil
}

func (c *Client) Query(table string, startTime int64, endTime int64, orderBy string, order string, limit int32,
	offset int32, conditions []interface{}, groupBy []string, fields []string) ([]string, []interface{}, error) {
	mapping, err := c.mapping(table)
	if err != nil {
		return nil, nil, err
	}
	tableFields, err := mappingResult(mapping)
	if err != nil {
		return nil, nil, err
	}

	request := map[string]interface{}{
		"from": offset,
		"size": limit,
	}
	booDsl, err := boolDSL(startTime, endTime, conditions, tableFields)
	if err != nil {
		return nil, nil, err
	}
	request["query"] = map[string]interface{}{
		"bool": booDsl,
	}

	if len(groupBy) > 0 {
		dsl, columns, metrics, err := aggsDSL(request, groupBy, orderBy, order, fields, tableFields)
		if err != nil {
			return nil, nil, err
		}
		respJSON, err := c.search(table, dsl)
		if err != nil {
			return nil, nil, err
		}
		result := aggsResult(respJSON, metrics, offset, limit)
		columns = replaceAggsField(fields, columns)
		return columns, result, nil
	}
	dsl, err := commonDSL(request, orderBy, order, fields, tableFields)
	if err != nil {
		return nil, nil, err
	}
	respJSON, err := c.search(table, dsl)
	if err != nil {
		return nil, nil, err
	}
	result := commonResult(respJSON, fields)
	columns := fields
	return columns, result, nil
}

func (c *Client) SimpleQuery(table string) ([]interface{}, error) {
	respJSON, err := c.search(table, `{"query":{"match_all":{}}}`)
	if err != nil {
		return nil, err
	}
	result := respJSON.Get("hits").Get("hits").MustArray()
	return result, nil
}

func (c *Client) search(table string, postData string) (*simplejson.Json, error) {
	url := fmt.Sprintf("%s/%s/_search", c.URL, table)
	return c.httpRequest(url, "POST", "json", postData)
}

func (c *Client) mapping(table string) (*simplejson.Json, error) {
	url := fmt.Sprintf("%s/%s/doc/_mapping", c.URL, table)
	return c.httpRequest(url, "GET", "", "")
}

func (c *Client) Indices() ([]string, error) {
	url := fmt.Sprintf("%s/_all/_settings", c.URL)
	respJSON, err := c.httpRequest(url, "GET", "", "")
	if err != nil {
		return nil, err
	}
	var names []string
	for name := range respJSON.MustMap() {
		names = append(names, name)
	}
	return names, nil
}

func (c *Client) httpRequest(requestURL string, method string, contentType string, body string) (*simplejson.Json, error) {
	var bodyData io.Reader
	if body != "" {
		bodyData = strings.NewReader(body)
	}
	request, _ := http.NewRequest(method, requestURL, bodyData)
	if contentType == "json" {
		request.Header.Add("Content-Type", "application/json")
	}
	request.SetBasicAuth(c.Username, c.Password)

	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respJSON, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	errReason1 := respJSON.Get("error").Get("reason").MustString("")
	if errReason1 != "" {
		return nil, errors.New(errReason1)
	}
	errReason2 := respJSON.Get("error").Get("caused_by").Get("reason").MustString("")
	if errReason2 != "" {
		return nil, errors.New(errReason2)
	}

	if respJSON.Get("_shards").Get("failed").MustInt() > 0 {
		return nil, errors.New("ES query shard error")
	}

	return respJSON, nil
}
