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
	"encoding/json"
	"fmt"
	htmlTemplate "html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"tkestack.io/tke/pkg/util/log"
)

var (
	//SelfdefineURL is used to post message to a user-defined URL if it's not empty
	SelfdefineURL string
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

// SelfdefineBodyInfo represents the body info to request a user-defined URL
type SelfdefineBodyInfo struct {
	Type   string `json:"type"`
	Header string `json:"header,omitempty"`
	Body   string `json:"body"`
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
	c := http.Client{}
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

// RequestToSelfdefine is used to do a post request to selfdefine URL
func RequestToSelfdefine(reqBody interface{}) error {
	reqURL, err := url.Parse(SelfdefineURL)
	if err != nil {
		return err
	}
	option := Option{
		Protocol: reqURL.Scheme,
		Host:     reqURL.Host,
		Path:     reqURL.Path,
		Method:   http.MethodPost,
		Body:     reqBody,
		Headers:  map[string]string{"Content-Type": "application/json"},
	}
	_, err = Request(option)
	return err
}

// ParseTemplate is used to get body according to template
func ParseTemplate(name string, template string, variables map[string]string) (string, error) {
	var buffer bytes.Buffer
	tmplBody := htmlTemplate.Must(htmlTemplate.New(name).Parse(template))
	err := tmplBody.Execute(&buffer, variables)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// GetCurrentTime returns current timestamp
func GetCurrentTime() int64 {
	return time.Now().Unix()
}
