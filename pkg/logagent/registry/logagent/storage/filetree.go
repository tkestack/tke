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
type FileNodeREST struct {
	//apiKeyStore *registry.Store
	//rest.Storage
	apiKeyStore    *registry.Store
	PlatformClient platformversionedclient.PlatformV1Interface
	//*registry.Store
}

var _ = rest.Creater(&FileNodeREST{})

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *FileNodeREST) New() runtime.Object {
	return &logagent.LogFileTree{}
}

type FileNodeRequest struct {
	PodName   string `json:"podName"`
	Namespace string `json:"namespace"`
	Container string `json:"container"`
}

type FileNodeProxy struct {
	Req  logagent.LogFileTreeSpec
	Ip   string
	Port string
}

func (p *FileNodeProxy) GetReaderCloser() (io.ReadCloser, error) {
	jsonStr, err := json.Marshal(p.Req)
	if err != nil {
		log.Errorf("unable to marshal request to json %v", err)
		return nil, fmt.Errorf("unable to marshal request")
	}
	url := "http://" + p.Ip + ":" + p.Port + "/v1/logfile/directory"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	httpReq.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Errorf("unable to generate request %v", err)
		return nil, fmt.Errorf("unable to generate request")
	}
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Errorf("unable to connect log-agent %v", err)
		return nil, fmt.Errorf("unable to connect log-agent")
	}
	return resp.Body, nil
}

func (r *FileNodeREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	//TODO: get cluster id from parent resource
	//userName, tenantID := authentication.GetUsernameAndTenantID(ctx)
	fileNode := obj.(*logagent.LogFileTree)
	//log.Infof("get userNmae %v tenantId %v and fileNode spec=%+v", userName, tenantID, fileNode.Spec)
	hostIp, err := util.GetClusterPodIp(ctx, fileNode.Spec.ClusterId, fileNode.Spec.Namespace, fileNode.Spec.Pod, r.PlatformClient)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("unable to get host ip"))
	}
	return &util.LocationStreamer{
		Request:     &FileNodeProxy{Req: fileNode.Spec, Ip: hostIp, Port: util.LogagentPort},
		Transport:   nil,
		ContentType: "application/json",
		Ip:          hostIp,
	}, nil
}
