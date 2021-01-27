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

package kubernetes

import (
	"fmt"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"tkestack.io/tke/pkg/util/log"
)

// KubeConfigClusterProvider get kubeconfig files from env var
// Only for local test
type KubeConfigClusterProvider struct {
	KubeConfigFiles string
}

// kubeconfigGetter get kube configs by clusterName
//  Example:
//      export LOCAL_KUBECONFIG=~/.kube/config1.yaml:~/.kube/config2.yaml
func (k KubeConfigClusterProvider) kubeconfigGetter(clusterName string) (func() (*clientcmdapi.Config, error), error) {
	envVarFiles := k.KubeConfigFiles
	filenameMap := make(map[string]struct{})
	if len(envVarFiles) != 0 {
		fileList := filepath.SplitList(envVarFiles)
		// prevent the same path load multiple times
		for _, f := range fileList {
			if _, exists := filenameMap[f]; !exists {
				filenameMap[f] = struct{}{}
			}
		}
	}

	clusters := make(map[string]*clientcmdapi.Config)
	for filename := range filenameMap {
		clusterConfig, err := clientcmd.LoadFromFile(filename)
		if err != nil {
			log.Errorf("loading cluster kubeconfig file[%s] failed: %v", filename, err)
			continue
		}
		for name := range clusterConfig.Clusters {
			clusters[name] = clusterConfig
		}
	}

	return func() (*clientcmdapi.Config, error) {
		config, exist := clusters[clusterName]
		if !exist {
			return nil, fmt.Errorf("cluster[%s] kubeconfig not found", clusterName)
		}
		return config, nil
	}, nil
}

// Only for local test
func NewKubeConfigProvider(kubeconfigFiles string) *KubeConfigClusterProvider {
	return &KubeConfigClusterProvider{
		KubeConfigFiles: kubeconfigFiles,
	}
}

func (k *KubeConfigClusterProvider) RestConfig(clusterName string) (*rest.Config, error) {
	// 2020-11-05 switching kube configs context by clusterName
	getter, err := k.kubeconfigGetter(clusterName)
	if err != nil {
		return nil, err
	}
	return clientcmd.BuildConfigFromKubeconfigGetter("", getter)
}

func (k *KubeConfigClusterProvider) Client(clusterName string, scheme *runtime.Scheme) (ctrlclient.Client, error) {
	config, err := k.RestConfig(clusterName)
	if err != nil {
		return nil, err
	}

	return ctrlclient.New(config, ctrlclient.Options{Scheme: scheme})
}
