/*
 * Tencent is pleased to support the open source community by making TKEStack available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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
	"context"
	"encoding/json"
	"sync"

	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/log"

	"github.com/caddyserver/caddy/caddyfile"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	DNS           = ".:53"
	REWRITE       = "rewrite"
	NAME          = "name"
	CoreFileKey   = "Corefile"
	KubeSystem    = "kube-system"
	ConfigMapName = "coredns"
	TargetSVC     = "tke-registry-api.tke.svc.cluster.local"
)

type CoreFileBlock struct {
	Body [][]string `json:"body"`
	Keys []string   `json:"keys"`
}

type CoreFile []*CoreFileBlock
type CoreFileHosts [][]string

type CoreDNS struct {
	sync.Mutex
	configMap  *corev1.ConfigMap
	kubeClient *kubernetes.Clientset
}

func NewCoreDNS() (*CoreDNS, error) {
	kubeClient, err := apiclient.BuildKubeClient()
	if err != nil {
		return nil, err
	}
	return &CoreDNS{
		kubeClient: kubeClient,
	}, nil
}

func (c *CoreDNS) LoadCoreFile() []byte {
	coreDNSConfigMap, err := c.kubeClient.CoreV1().ConfigMaps(KubeSystem).Get(context.Background(),
		ConfigMapName, metav1.GetOptions{})
	if err != nil {
		log.Error("get coredns's configmap failed", log.Err(err))
		return nil
	}
	c.configMap = coreDNSConfigMap
	return []byte(coreDNSConfigMap.Data[CoreFileKey])
}

func (c *CoreDNS) StoreCoreFile(content []byte) {
	c.configMap.Data[CoreFileKey] = string(content)
	_, err := c.kubeClient.CoreV1().ConfigMaps(KubeSystem).Update(context.Background(), c.configMap, metav1.UpdateOptions{})
	if err != nil {
		log.Error("update coredns's configmap failed", log.Err(err))
	}
}

func (c *CoreDNS) ParseCoreFile(item string) {
	c.Lock()
	defer c.Unlock()
	content := c.LoadCoreFile()
	if content == nil {
		return
	}
	coreFileBody, _ := caddyfile.ToJSON(content)
	coreFileObj := new(CoreFile)
	_ = json.Unmarshal(coreFileBody, coreFileObj)
	for _, block := range *coreFileObj {
		if c.ParseBlockDNS(block, item) {
			break
		}
	}
	coreFileBody, _ = json.Marshal(coreFileObj)
	content, _ = caddyfile.FromJSON(coreFileBody)
	c.StoreCoreFile(content)
}

func (c *CoreDNS) ParseBlockDNS(block *CoreFileBlock, item string) bool {
	if block.Keys == nil || len(block.Keys) == 0 || block.Keys[0] != DNS {
		return false
	}
	findItem := false
	for _, s := range block.Body {
		if c.ParseSectionRewrite(s, item) {
			findItem = true
			break
		}
	}
	if !findItem {
		block.Body = append(block.Body, []string{REWRITE, NAME, item, TargetSVC})
	}
	return true
}

func (c *CoreDNS) ParseSectionRewrite(section []string, item string) bool {
	if section == nil || len(section) != 4 || section[0] != REWRITE {
		return false
	}
	if section[2] == item {
		section[3] = TargetSVC
		return true
	}
	return false
}
