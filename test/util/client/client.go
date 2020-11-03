/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package client

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
	"tkestack.io/tke/test/util/env"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	tkeclientset "tkestack.io/tke/api/client/clientset/versioned"

	// env for load env
	_ "tkestack.io/tke/test/util/env"
)

func LoadOrSetupTKE() (tkeclientset.Interface, error) {
	data := env.Kubeconfig()
	if data == "" {
		t := time.Now()
		log.Printf("env %s is not set, using tke-up", env.KUBECONFIG)
		cmd := exec.Command("tke-up")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, errors.Wrap(err, "run tke-up error")
		}
		log.Printf("tke-up is run sucessully [%s]", time.Since(t).String())
		os.Setenv(env.KUBECONFIG, string(output))
	}

	return GetTKEClientFromEnv()
}

func GetTKEClientFromEnv() (tkeclientset.Interface, error) {
	data := env.Kubeconfig()
	if data == "" {
		return nil, errors.Errorf("%s not set", env.KUBECONFIG)
	}
	kubeconfig, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, errors.Wrap(err, "base64 decode error")
	}

	return GetTKEClient(kubeconfig)
}

func GetTKEClient(kubeconfig []byte) (tkeclientset.Interface, error) {
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	return tkeclientset.NewForConfig(restConfig)
}

func GetRESTConfig() *rest.Config {
	data := env.Kubeconfig()
	if data == "" {
		panic(fmt.Sprintf("%s not set", env.KUBECONFIG))
	}
	kubeconfig, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	}
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		panic(err)
	}
	return restConfig
}

func GetClientSet() *kubernetes.Clientset {
	return kubernetes.NewForConfigOrDie(GetRESTConfig())
}
