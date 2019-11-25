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
	"io/ioutil"
	"reflect"
	"text/template"

	"github.com/pkg/errors"
	kuberuntime "k8s.io/apimachinery/pkg/runtime"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

// ParseFileTemplate parse file template
func ParseFileTemplate(filename string, obj interface{}) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ParseTemplate(string(data), obj)
}

// ParseTemplate validates and parses passed as argument template
func ParseTemplate(strtmpl string, obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	tmpl, err := template.New("template").Parse(strtmpl)
	if err != nil {
		return nil, errors.Wrap(err, "error when parsing template")
	}
	err = tmpl.Execute(&buf, obj)
	if err != nil {
		return nil, errors.Wrap(err, "error when executing template")
	}
	return buf.Bytes(), nil
}

func ParseTemplateTo(strtmpl string, obj interface{}, dst interface{}) error {
	b, err := ParseTemplate(strtmpl, obj)
	if err != nil {
		return errors.Wrapf(err, "error when parsing template")
	}

	if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), b, dst.(kuberuntime.Object)); err != nil {
		return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(dst).String())
	}

	return nil
}
