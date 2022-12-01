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

package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html"
	htmlTemplate "html/template"

	"io/ioutil"
	"net/http"
	"time"

	"tkestack.io/tke/pkg/util/log"
)

// Option is used to post request
type Option struct {
	Protocol string            `json:"protocol"`
	Host     string            `json:"host"`
	Path     string            `json:"path"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	Body     interface{}       `json:"body"`
}

// Request is used to do a post request
func Request(options Option) ([]byte, error) {
	var err error
	var rawBody []byte
	rawBody, err = json.Marshal(options.Body)
	if err != nil {
		return nil, err
	}
	log.Debugf("rawBody: %v", string(rawBody))
	body := bytes.NewReader(rawBody)
	URL := options.Protocol + `://` + options.Host + options.Path

	var req *http.Request
	req, err = http.NewRequest(options.Method, URL, body)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	for k, v := range options.Headers {
		req.Header.Add(k, v)
	}
	var resp *http.Response
	c := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err = c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : url=%v , statusCode=%v", URL, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

// ParseString is used to get body according to template
func ParseTemplate(name string, template string, variables map[string]string) (string, error) {
	var buffer bytes.Buffer
	tmplBody := htmlTemplate.Must(htmlTemplate.New(name).Parse(template))
	err := tmplBody.Execute(&buffer, variables)
	if err != nil {
		return "", err
	}
	return html.UnescapeString(buffer.String()), nil
}

// GetCurrentTime returns current timestamp
func GetCurrentTime() int64 {
	return time.Now().Unix()
}

//render the alert status before sending alert message
func GetAlertStatus(status string) string {
	var alertStatus string
	if status == string("firing") {
		alertStatus = "未恢复"
	} else {
		alertStatus = "已恢复"
	}
	alertStatus = fmt.Sprintf("告警状态： %s", alertStatus)
	return alertStatus
}
