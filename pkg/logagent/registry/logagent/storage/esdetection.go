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

package storage

import (
	//"bytes"
	"context"
	"io/ioutil"
	"strings"
	"time"

	//"encoding/json"
	"fmt"
	"io"
	"net/http"

	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/logagent"
	"tkestack.io/tke/pkg/logagent/util"
	"tkestack.io/tke/pkg/util/log"
)

// TokenREST implements the REST endpoint.
type ESDetectionREST struct {
	apiKeyStore    *registry.Store
	PlatformClient platformversionedclient.PlatformV1Interface
}

var _ = rest.Creater(&ESDetectionREST{})

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *ESDetectionREST) New() runtime.Object {
	return &logagent.LogEsDetection{}
}

type ESDetectionProxy struct {
	Scheme   string
	IP       string
	Port     string
	User     string
	Password string
}

func (p *ESDetectionProxy) GetReaderCloser() (io.ReadCloser, error) {
	if p.Scheme == "" || p.IP == "" || p.Port == "" {
		return nil, fmt.Errorf("es detection: scheme, ip, port maybe null")
	}

	url := fmt.Sprintf("%s://%s:%s@%s:%s/_cat/health", p.Scheme, p.User, p.Password, p.IP, p.Port)
	httpReq, err := http.NewRequest("GET", url, nil)
	httpReq.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Errorf("es detection: unable to generate request %v", err)
		return nil, fmt.Errorf("es detection: unable to generate request")
	}

	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Errorf("es detection: unable to connect %v", err)
		return nil, fmt.Errorf("es detection: unable to connect")
	}

	if resp.StatusCode != 200 {
		return ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"status": "Failure", "code": "%d"}`, resp.StatusCode))), nil
	}
	//return resp.Body, nil
	return ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"status": "Success", "code": "%d"}`, resp.StatusCode))), nil
}

func (r *ESDetectionREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	esDetection := obj.(*logagent.LogEsDetection)
	return &util.LocationStreamer{
		Request:     &ESDetectionProxy{Scheme: esDetection.Scheme, IP: esDetection.IP, Port: esDetection.Port, User: esDetection.User, Password: esDetection.Password},
		Transport:   nil,
		ContentType: "application/json",
	}, nil
}
