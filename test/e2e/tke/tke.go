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

package tke

import (
	"context"
	"fmt"
	"k8s.io/klog"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/onsi/gomega"

	"tkestack.io/tke/cmd/tke-installer/app/installer/types"

	"tkestack.io/tke/test/e2e"

	"github.com/onsi/ginkgo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/test/e2e/certs"
	testclient "tkestack.io/tke/test/util/client"
	"tkestack.io/tke/test/util/env"
)

func tkeHostName() string {
	restconf := testclient.GetRESTConfig()
	host := restconf.Host
	u, _ := url.Parse(host)
	return u.Hostname()
}

type TKE struct {
	Namespace string
	certs.TkeCert
	hostName string
	client   *kubernetes.Clientset
	steps    []types.Handler
}

func (t *TKE) Create() {
	t.initSteps()
	ctx := context.Background()
	for _, v := range t.steps {
		ginkgo.By(v.Name)
		fmt.Println(v.Name)
		gomega.Expect(v.Func(ctx)).To(gomega.BeNil())
	}
}

func (t *TKE) Delete() {
	t.ClearTmpDir()
	klog.Info("Delete namespace: ", t.Namespace)
	gracePeriodSeconds := int64(0)
	deleteOptions := metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
	}
	client := testclient.GetClientSet()

	// Workaround for below issue:
	// Operation cannot be fulfilled on namespaces \"platform\": The system is ensuring all content is removed from
	// this namespace.  Upon completion, this namespace will automatically be purged by the system.
	ns, _ := client.CoreV1().Namespaces().Get(context.Background(), t.Namespace, metav1.GetOptions{})
	ns.Spec.Finalizers = []corev1.FinalizerName{} // remove all finalizers
	_, err := client.CoreV1().Namespaces().Update(context.Background(), ns, metav1.UpdateOptions{})
	if err != nil {
		klog.Error(err)
	}
	time.Sleep(5 * time.Second)

	err = client.CoreV1().Namespaces().Delete(context.Background(), t.Namespace, deleteOptions)
	gomega.Expect(err).To(gomega.BeNil())
}

func (t *TKE) initialize(ctx context.Context) error {
	t.InitTmpDir(t.Namespace)
	str, _ := os.Getwd()
	fmt.Println("Before suite:", str, " version:", env.ImageVersion())
	t.client = testclient.GetClientSet()
	t.hostName = tkeHostName()
	fmt.Println("tkeHostName:", t.hostName)
	return nil
}

func (t *TKE) createNamespace(ctx context.Context) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: t.Namespace,
		},
	}

	return apiclient.CreateOrUpdateNamespace(context.Background(), t.client, ns)
}

func (t *TKE) createCertMap(ctx context.Context) error {
	ips := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP(t.hostName)}
	dns := []string{"etcd-client"}
	return t.CreateCertMap(context.Background(), t.client, dns, ips, t.Namespace)
}

func (t *TKE) createProviderConfigs(ctx context.Context) error {
	configMaps := []struct {
		Name string
		File string
	}{
		{
			Name: "provider-config",
			File: e2e.ConfDir + "config.yaml",
		},
		{
			Name: "docker",
			File: e2e.ConfDir + "docker/*",
		},
		{
			Name: "kubelet",
			File: e2e.ConfDir + "kubelet/*",
		},
		{
			Name: "kubeadm",
			File: e2e.ConfDir + "kubeadm/*",
		},
	}

	for _, one := range configMaps {
		err := apiclient.CreateOrUpdateConfigMapFromFile(ctx, t.client,
			&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      one.Name,
					Namespace: t.Namespace,
				},
			},
			one.File)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TKE) createEtcd(ctx context.Context) error {
	options := map[string]interface{}{
		"Servers":   []string{t.hostName},
		"Namespace": t.Namespace,
	}
	err := apiclient.CreateResourceWithDir(context.Background(), t.client, e2e.EtcdYamlFile, options)
	if err != nil {
		return err
	}
	return wait.PollImmediate(5*time.Second, 2*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckPodReadyWithLabel(context.Background(), t.client, t.Namespace, "app=etcd")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) createPlatformAPI(ctx context.Context) error {
	options := map[string]interface{}{
		"Image":       fmt.Sprintf("tkestack/tke-platform-api-amd64:%s", env.ImageVersion()),
		"EnableAuth":  false,
		"EnableAudit": false,
		"Namespace":   t.Namespace,
	}

	err := apiclient.CreateResourceWithDir(context.Background(), t.client, e2e.TKEPlatformAPIYAMLFile, options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 2*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckPodReadyWithLabel(context.Background(), t.client, t.Namespace, "app=tke-platform-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) createPlatformController(ctx context.Context) error {
	options := map[string]interface{}{
		"Image":             fmt.Sprintf("tkestack/tke-platform-controller-amd64:%s", env.ImageVersion()),
		"ProviderResImage":  fmt.Sprintf("tkestack/provider-res-amd64:%s", env.ProviderResImageVersion()),
		"RegistryDomain":    "docker.io",
		"RegistryNamespace": "tkestack",
		"Namespace":         t.Namespace,
	}

	err := apiclient.CreateResourceWithDir(context.Background(), t.client, e2e.TKEPlatformControllerYAMLFile, options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 2*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckPodReadyWithLabel(context.Background(), t.client, t.Namespace,
			"app=tke-platform-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) createKubeConfig(ctx context.Context) error {
	svc, err := t.client.CoreV1().Services(t.Namespace).Get(context.Background(), "tke-platform-api",
		metav1.GetOptions{})
	if err != nil {
		return err
	}

	nodePort := int(svc.Spec.Ports[0].NodePort)
	return t.WriteKubeConfig(t.hostName, nodePort, t.Namespace)
}

func (t *TKE) initSteps() {
	t.steps = append(t.steps, []types.Handler{
		{
			Name: "initialize",
			Func: t.initialize,
		},
	}...)
	t.steps = append(t.steps, []types.Handler{
		{
			Name: "create namespace",
			Func: t.createNamespace,
		},
	}...)
	t.steps = append(t.steps, []types.Handler{
		{
			Name: "create cert map",
			Func: t.createCertMap,
		},
	}...)
	t.steps = append(t.steps, []types.Handler{
		{
			Name: "create provider configs",
			Func: t.createProviderConfigs,
		},
	}...)
	t.steps = append(t.steps, []types.Handler{
		{
			Name: "create etcd",
			Func: t.createEtcd,
		},
	}...)
	t.steps = append(t.steps, []types.Handler{
		{
			Name: "create tke-platform-api",
			Func: t.createPlatformAPI,
		},
	}...)
	t.steps = append(t.steps, []types.Handler{
		{
			Name: "create tke-platform-controller",
			Func: t.createPlatformController,
		},
	}...)
	t.steps = append(t.steps, []types.Handler{
		{
			Name: "create kubeconfig",
			Func: t.createKubeConfig,
		},
	}...)
}
