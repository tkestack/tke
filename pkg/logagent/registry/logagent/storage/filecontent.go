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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"
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
type FileContentREST struct {
	//apiKeyStore *registry.Store
	//rest.Storage
	apiKeyStore    *registry.Store
	PlatformClient platformversionedclient.PlatformV1Interface
	//*registry.Store
}

var _ = rest.Creater(&FileContentREST{})

func (r *FileContentREST) New() runtime.Object {
	return &logagent.LogFileContent{}
}

type FileContentProxy struct {
	Req  logagent.LogFileContentSpec
	IP   string
	Port string
}

func (p *FileContentProxy) GetReaderCloser() (io.ReadCloser, error) {
	jsonStr, err := json.Marshal(p.Req)
	if err != nil {
		log.Errorf("unable to marshal request to json %v", err)
		return nil, fmt.Errorf("unable to marshal request")
	}
	url := "http://" + p.IP + ":" + p.Port + "/v1/logfile/content"
	log.Infof("url is %v", url)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	httpReq.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Errorf("unable to generate request %v", err)
		return nil, fmt.Errorf("unable to generate request")
	}
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Errorf("unable to connect to log-agent %v", err)
		return nil, fmt.Errorf("unable to connect log-agent")
	}
	return resp.Body, nil
}

func (r *FileContentREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	//TODO: get cluster id from parent resource
	//userName, tenantID := authentication.GetUsernameAndTenantID(ctx)
	fileContent := obj.(*logagent.LogFileContent)
	//log.Infof("get userNmae %v tenantId %v and fileNode spec=%+v", userName, tenantID, fileContent.Spec)
	hostIP, err := util.GetClusterPodIP(ctx, fileContent.Spec.ClusterId, fileContent.Spec.Namespace, fileContent.Spec.Pod, r.PlatformClient)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("unable to get pod ip %v", err))
	}
	return &util.LocationStreamer{
		Request:     &FileContentProxy{Req: fileContent.Spec, IP: hostIP, Port: util.LogagentPort},
		Transport:   nil,
		ContentType: "application/json",
		IP:          hostIP,
	}, nil
}
