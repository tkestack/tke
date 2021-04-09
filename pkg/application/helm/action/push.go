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
 */

package action

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	cm "github.com/chartmuseum/helm-push/pkg/chartmuseum"
	"github.com/chartmuseum/helm-push/pkg/helm"
	"tkestack.io/tke/pkg/util/log"
)

// PushOptions is the options for pushing a chart.
type PushOptions struct {
	ChartPathOptions

	ChartFile   string
	ForceUpload bool
}

// Push is the action for pushing a chart.
func (c *Client) Push(options *PushOptions) (existed bool, err error) {
	client, err := cm.NewClient(
		cm.URL(options.RepoURL),
		cm.Username(options.Username),
		cm.Password(options.Password),
		cm.CAFile(options.CaFile),
		cm.CertFile(options.CertFile),
		cm.KeyFile(options.KeyFile),
		cm.InsecureSkipVerify(true),
	)
	if err != nil {
		return false, err
	}

	repo, err := helm.TempRepoFromURL(options.RepoURL)
	if err != nil {
		return false, err
	}
	index, err := helm.GetIndexByRepo(repo, getIndexDownloader(client))
	if err != nil {
		return false, err
	}
	client.Option(cm.ContextPath(index.ServerInfo.ContextPath))

	log.Infof("Pushing %s to %s", filepath.Base(options.ChartFile), options.RepoURL)
	resp, err := client.UploadChartPackage(options.ChartFile, options.ForceUpload)
	if err != nil {
		return false, err
	}
	return handlePushResponse(resp)
}

func handlePushResponse(resp *http.Response) (existed bool, err error) {
	if resp.StatusCode != 201 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		if resp.StatusCode == 409 {
			return true, getChartmuseumError(b, resp.StatusCode)
		}
		return false, getChartmuseumError(b, resp.StatusCode)
	}
	return false, nil
}

func getIndexDownloader(client *cm.Client) helm.IndexDownloader {
	return func() ([]byte, error) {
		resp, err := client.DownloadFile("index.yaml")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, getChartmuseumError(b, resp.StatusCode)
		}
		return b, nil
	}
}

func getChartmuseumError(b []byte, code int) error {
	var er struct {
		Error string `json:"error"`
	}
	err := json.Unmarshal(b, &er)
	if err != nil || er.Error == "" {
		return fmt.Errorf("%d: could not properly parse response JSON: %s", code, string(b))
	}
	return fmt.Errorf("%d: %s", code, er.Error)
}
