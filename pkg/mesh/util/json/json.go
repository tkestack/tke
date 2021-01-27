/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
 *
 */

package json

import (
	"bytes"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var (
	json                = jsoniter.ConfigCompatibleWithStandardLibrary
	MarshalToString     = json.MarshalToString
	Marshal             = json.Marshal
	MarshalIndent       = json.MarshalIndent
	Unmarshal           = json.Unmarshal
	UnmarshalFromString = json.UnmarshalFromString
	NewDecoder          = json.NewDecoder
	NewEncoder          = json.NewEncoder
	Get                 = json.Get
)

func NewJSONRequest(req *http.Request) (*Request, error) {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read the request body")
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	return &Request{
		req,
		buf,
	}, nil
}

type Request struct {
	req *http.Request
	raw []byte
}

func (j *Request) FindObject(jsonpath ...interface{}) jsoniter.Any {
	return json.Get(j.raw, jsonpath...)
}
