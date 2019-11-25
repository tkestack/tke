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

package template

import (
	"bytes"
	"reflect"
	"text/template"

	"github.com/pkg/errors"
	kuberuntime "k8s.io/apimachinery/pkg/runtime"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

// ParseFile parse template file with obj
func ParseFile(filename string, obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		return nil, errors.Wrap(err, "error when parse file")
	}
	err = tmpl.Execute(&buf, obj)
	if err != nil {
		return nil, errors.Wrap(err, "error when executing template")
	}

	return buf.Bytes(), nil
}

// ParseFileToObject parse template file and decode to k8s object
func ParseFileToObject(filename string, obj interface{}, dst interface{}) error {
	var buf bytes.Buffer
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		return errors.Wrap(err, "error when parse file")
	}
	err = tmpl.Execute(&buf, obj)
	if err != nil {
		return errors.Wrap(err, "error when executing template")
	}

	if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), buf.Bytes(), dst.(kuberuntime.Object)); err != nil {
		return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(dst).String())
	}

	return nil
}
