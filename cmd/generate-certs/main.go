/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

package main

import (
	"fmt"
	"net"

	"k8s.io/client-go/kubernetes"

	"github.com/spf13/pflag"
	"k8s.io/client-go/tools/clientcmd"
	"tkestack.io/tke/pkg/util/certs"
)

type Option struct {
	Kubeconfig string

	CACert string
	CAKey  string

	Name      string
	Namespace string
	IPs       []net.IP
	DNSNames  []string
}

var (
	option Option
)

func init() {
	pflag.StringVar(&option.Kubeconfig, "kubeconfig", "", "Path to the kubeconfig file to use for create certs")

	pflag.StringVar(&option.CACert, "cacert", "", "CA certificate file")
	pflag.StringVar(&option.CAKey, "cakey", "", "CA key file")

	pflag.StringVar(&option.Name, "name", "tke-certs", "The name of tke certs secrets")
	pflag.StringVar(&option.Namespace, "namespace", "", "The namespace of tke certs secrets")
	pflag.IPSliceVar(&option.IPs, "ips", nil, "Extra ips")
	pflag.StringSliceVar(&option.DNSNames, "dnsnames", nil, "Extra dns names")
}

func main() {
	pflag.Parse()

	if option.Namespace == "" {
		pflag.Usage()
		fmt.Println("namespace required")
		return
	}

	err := run(&option)
	if err != nil {
		fmt.Println(err)
	}
}

func run(option *Option) error {
	if option.Kubeconfig == "" {
		return certs.GenerateInDir("", option.Namespace, option.DNSNames, option.IPs)
	}
	config, err := clientcmd.BuildConfigFromFlags("", option.Kubeconfig)
	if err != nil {
		return err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	return certs.GenerateInK8s(client, option.Name, option.Namespace, option.DNSNames, option.IPs)
}
