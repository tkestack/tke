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

package proxy

import (
	"context"
	"fmt"
	"strings"
	"sync"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
	"tkestack.io/tke/pkg/platform/types"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/pkiutil"
)

type clientX509Pool struct {
	sm sync.Map
}

var pool clientX509Pool

type clientX509Cache struct {
	clientCertData []byte
	clientKeyData  []byte
}

func makeClientKey(username string, groups []string) string {
	return fmt.Sprintf("%s###%v", username, groups)
}

func ClientSet(ctx context.Context, platformClient platforminternalclient.PlatformInterface) (*kubernetes.Clientset,
	error) {
	clusterName := filter.ClusterFrom(ctx)
	if clusterName == "" {
		return nil, errors.NewBadRequest("clusterName is required")
	}

	cluster, err := platformClient.Clusters().Get(ctx, clusterName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) > 0 && cluster.Spec.TenantID != tenantID {
		return nil, errors.NewNotFound(platform.Resource("clusters"), cluster.ObjectMeta.Name)
	}

	clusterWrapper, err := types.GetCluster(ctx, platformClient, cluster)
	if err != nil {
		return nil, err
	}

	config := &rest.Config{}
	//if cluster.AuthzWebhookEnabled() {
	//	clientCertData, clientKeyData, err := getOrCreateClientCert(ctx, clusterWrapper.ClusterCredential)
	//	if err != nil {
	//		return nil, err
	//	}
	//	config, err = clusterWrapper.RESTConfigForClientX509(config, clientCertData, clientKeyData)
	//	if err != nil {
	//		return nil, err
	//	}
	//} else {
	//	config, err = clusterWrapper.RESTConfig(config)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	//转发给api-server的请求，都需要使用当前用户的证书去访问，如果没有证书，则生成证书
	clientCertData, clientKeyData, err := getOrCreateClientCert(ctx, clusterWrapper)
	if err != nil {
		return nil, err
	}
	config, err = clusterWrapper.RESTConfigForClientX509(config, clientCertData, clientKeyData)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func getOrCreateClientCert(ctx context.Context, clusterWrapper *types.Cluster) ([]byte, []byte, error) {
	credential := clusterWrapper.ClusterCredential
	groups := authentication.Groups(ctx)
	username, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID != "" {
		groups = append(groups, fmt.Sprintf("tenant:%s", tenantID))
	}

	ns, ok := request.NamespaceFrom(ctx)
	if ok {
		groups = append(groups, fmt.Sprintf("namespace:%s", ns))
	}
	uid := authentication.GetUID(ctx)
	clusterName := filter.ClusterFrom(ctx)
	if clusterName == "" {
		return nil, nil, errors.NewBadRequest("clusterName is required")
	}
	var clientCertData, clientKeyData []byte
	client, _ := clusterWrapper.Clientset()
	cache, err := client.CoreV1().ConfigMaps("kube-system").Get(ctx, clusterName+"-"+uid, metav1.GetOptions{})
	if err != nil || cache == nil {
		clientCertData, clientKeyData, err := pkiutil.GenerateClientCertAndKey(username, groups, credential.CACert,
			credential.CAKey)
		if err != nil {
			return nil, nil, err
		}
		confMap := &v1.ConfigMap{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterName + "-" + uid,
			},
			Data: map[string]string{
				"CommonName": username,
			},
			BinaryData: map[string][]byte{
				"clientCertData": clientCertData,
				"clientKeyData":  clientKeyData,
			},
		}
		client.CoreV1().ConfigMaps("kube-system").Create(ctx, confMap, metav1.CreateOptions{})
	} else {
		clientCertData = cache.BinaryData["clientCertData"]
		clientKeyData = cache.BinaryData["clientKeyData"]
	}

	log.Debugf("generateClientCert success. username:%s groups:%v\n clientCertData:\n %s clientKeyData:\n %s",
		username, groups, clientCertData, clientKeyData)

	return clientCertData, clientKeyData, nil
}

// RESTClient returns the versioned rest client of clientSet.
func RESTClient(ctx context.Context, platformClient platforminternalclient.PlatformInterface) (restclient.Interface, *request.RequestInfo, error) {
	request, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, nil, errors.NewBadRequest("unable to get request info from context")
	}
	clientSet, err := ClientSet(ctx, platformClient)
	if err != nil {
		return nil, nil, err
	}
	client := RESTClientFor(clientSet, request.APIGroup, request.APIVersion)
	return client, request, nil
}

// RESTClientFor returns the versioned rest client of clientSet by given api
// version.
func RESTClientFor(clientSet *kubernetes.Clientset, apiGroup, apiVersion string) restclient.Interface {
	gv := fmt.Sprintf("%s/%s", strings.ToLower(apiGroup), strings.ToLower(apiVersion))
	switch gv {
	case "/v1":
		return clientSet.CoreV1().RESTClient()
	case "apps/v1":
		return clientSet.AppsV1().RESTClient()
	case "apps/v1beta1":
		return clientSet.AppsV1beta1().RESTClient()
	case "admissionregistration.k8s.io/v1beta1":
		return clientSet.AdmissionregistrationV1beta1().RESTClient()
	case "apps/v1beta2":
		return clientSet.AppsV1beta2().RESTClient()
	case "autoscaling/v1":
		return clientSet.AutoscalingV1().RESTClient()
	case "autoscaling/v2beta1":
		return clientSet.AutoscalingV2beta1().RESTClient()
	case "batch/v1":
		return clientSet.BatchV1().RESTClient()
	case "batch/v1beta1":
		return clientSet.BatchV1beta1().RESTClient()
	case "batch/v2alpha1":
		return clientSet.BatchV2alpha1().RESTClient()
	case "certificates.k8s.io/v1beta1":
		return clientSet.CertificatesV1beta1().RESTClient()
	case "events.k8s.io/v1beta1":
		return clientSet.EventsV1beta1().RESTClient()
	case "extensions/v1beta1":
		return clientSet.ExtensionsV1beta1().RESTClient()
	case "networking.k8s.io/v1":
		return clientSet.NetworkingV1().RESTClient()
	case "networking.k8s.io/v1beta1":
		return clientSet.NetworkingV1beta1().RESTClient()
	case "coordination.k8s.io/v1":
		return clientSet.CoordinationV1().RESTClient()
	case "coordination.k8s.io/v1beta1":
		return clientSet.CoordinationV1beta1().RESTClient()
	case "policy/v1beta1":
		return clientSet.PolicyV1beta1().RESTClient()
	case "rbac.authorization.k8s.io/v1alpha1":
		return clientSet.RbacV1alpha1().RESTClient()
	case "rbac.authorization.k8s.io/v1":
		return clientSet.RbacV1().RESTClient()
	case "rbac.authorization.k8s.io/v1beta1":
		return clientSet.RbacV1beta1().RESTClient()
	case "scheduling.k8s.io/v1alpha1":
		return clientSet.SchedulingV1alpha1().RESTClient()
	case "scheduling.k8s.io/v1beta1":
		return clientSet.SchedulingV1beta1().RESTClient()
	case "node.k8s.io/v1alpha1":
		return clientSet.NodeV1alpha1().RESTClient()
	case "node.k8s.io/v1beta1":
		return clientSet.NodeV1beta1().RESTClient()
	case "scheduling.k8s.io/v1":
		return clientSet.SchedulingV1().RESTClient()
	case "settings.k8s.io/v1alpha1":
		return clientSet.SettingsV1alpha1().RESTClient()
	case "storage.k8s.io/v1alpha1":
		return clientSet.StorageV1alpha1().RESTClient()
	case "storage.k8s.io/v1":
		return clientSet.StorageV1().RESTClient()
	case "storage.k8s.io/v1beta1":
		return clientSet.StorageV1beta1().RESTClient()
	default:
		return clientSet.RESTClient()
	}
}
