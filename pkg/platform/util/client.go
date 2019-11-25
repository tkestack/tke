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

package util

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	monitoringclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	pkgerrors "github.com/pkg/errors"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
	"tkestack.io/tke/pkg/util/log"
)

const (
	contextName = "tke"
)

// ClientSetByCluster returns the backend kubernetes clientSet by given cluster object
func ClientSetByCluster(ctx context.Context, cluster *platform.Cluster, platformClient platforminternalclient.PlatformInterface) (*kubernetes.Clientset, error) {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) > 0 && cluster.Spec.TenantID != tenantID {
		return nil, errors.NewNotFound(platform.Resource("clusters"), cluster.ObjectMeta.Name)
	}
	credential, err := ClusterCredential(platformClient, cluster.Name)
	if err != nil {
		return nil, err
	}

	return BuildClientSet(cluster, credential)
}

// ClientSet returns the backend kubernetes clientSet
func ClientSet(ctx context.Context, platformClient platforminternalclient.PlatformInterface) (*kubernetes.Clientset, error) {
	clusterName := filter.ClusterFrom(ctx)
	if clusterName == "" {
		return nil, errors.NewBadRequest("clusterName is required")
	}

	cluster, err := platformClient.Clusters().Get(clusterName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return ClientSetByCluster(ctx, cluster, platformClient)
}

// RESTClient returns the versioned rest client of clientSet.
func RESTClient(ctx context.Context, platformClient platforminternalclient.PlatformInterface) (restclient.Interface, *request.RequestInfo, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, nil, errors.NewBadRequest("unable to get request info from context")
	}
	clientSet, err := ClientSet(ctx, platformClient)
	if err != nil {
		return nil, nil, err
	}
	client := RESTClientFor(clientSet, requestInfo.APIGroup, requestInfo.APIVersion)
	return client, requestInfo, nil
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

// ResourceFromKind returns the resource name by kind.
func ResourceFromKind(kind string) string {
	kindLower := strings.ToLower(kind)
	switch kindLower {
	case "policy":
		return "policies"
	case "storageclass":
		return "storageclasses"
	case "ingress":
		return "ingresses"
	case "networkpolicy":
		return "networkpolicies"
	case "podsecuritypolicy":
		return "podsecuritypolicies"
	case "priorityclass":
		return "priorityclasses"
	default:
		return kindLower + "s"
	}
}

// BuildTransport create the http transport for communicate to backend
// kubernetes api server.
func BuildTransport(credential *platform.ClusterCredential) (http.RoundTripper, error) {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	if len(credential.CACert) > 0 {
		transport.TLSClientConfig = &tls.Config{
			RootCAs: rootCertPool(credential.CACert),
		}
	} else {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	if credential.ClientKey != nil && credential.ClientCert != nil {
		cert, err := tls.X509KeyPair(credential.ClientCert, credential.ClientKey)
		if err != nil {
			return nil, err
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
	}

	return transport, nil
}

// GetRestConfig returns rest config according to cluster
func GetRestConfig(cluster *platformv1.Cluster, credential *platformv1.ClusterCredential) (*restclient.Config, error) {
	address, err := ClusterV1Address(cluster)
	if err != nil {
		return nil, err
	}
	config := api.NewConfig()
	config.CurrentContext = contextName

	if credential.CACert == nil {
		config.Clusters[contextName] = &api.Cluster{
			Server:                "https://" + address,
			InsecureSkipTLSVerify: true,
		}
	} else {
		config.Clusters[contextName] = &api.Cluster{
			Server:                   "https://" + address,
			CertificateAuthorityData: credential.CACert,
		}
	}

	if credential.Token != nil {
		config.AuthInfos[contextName] = &api.AuthInfo{
			Token: *credential.Token,
		}
	} else if credential.ClientCert != nil && credential.ClientKey != nil {
		config.AuthInfos[contextName] = &api.AuthInfo{
			ClientCertificateData: credential.ClientCert,
			ClientKeyData:         credential.ClientKey,
		}
	} else {
		return nil, fmt.Errorf("no credential for the cluster")
	}

	config.Contexts[contextName] = &api.Context{
		Cluster:  contextName,
		AuthInfo: contextName,
	}
	clientConfig := clientcmd.NewNonInteractiveClientConfig(*config, contextName, &clientcmd.ConfigOverrides{Timeout: "5s"}, nil)
	return clientConfig.ClientConfig()
}

// BuildExternalClientSetNoStatus creates the clientset of kubernetes by given
// cluster object and returns it.
func BuildExternalClientSetNoStatus(cluster *platformv1.Cluster, credential *platformv1.ClusterCredential) (*kubernetes.Clientset, error) {
	restConfig, err := GetRestConfig(cluster, credential)
	if err != nil {
		log.Error("Build cluster config error", log.String("clusterName", cluster.ObjectMeta.Name), log.Err(err))
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

// BuildExternalClientSet creates the clientset of kubernetes by given cluster
// object and returns it.
func BuildExternalClientSet(cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*kubernetes.Clientset, error) {
	credential, err := ClusterCredentialV1(client, cluster.Name)
	if err != nil {
		return nil, err
	}

	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	return BuildExternalClientSetNoStatus(cluster, credential)
}

// BuildExternalClientSetWithName creates the clientset of kubernetes by given cluster
// name and returns it.
func BuildExternalClientSetWithName(platformClient platformversionedclient.PlatformV1Interface, name string) (*kubernetes.Clientset, error) {
	cluster, err := platformClient.Clusters().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	clientset, err := BuildExternalClientSet(cluster, platformClient)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// BuildExternalExtensionClientSetNoStatus creates the api extension clientset of kubernetes by given
// cluster object and returns it.
func BuildExternalExtensionClientSetNoStatus(cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*apiextensionsclient.Clientset, error) {
	credential, err := ClusterCredentialV1(client, cluster.Name)
	if err != nil {
		return nil, err
	}
	restConfig, err := GetRestConfig(cluster, credential)
	if err != nil {
		log.Error("Build cluster config error", log.String("clusterName", cluster.ObjectMeta.Name), log.Err(err))
		return nil, err
	}
	return apiextensionsclient.NewForConfig(restConfig)
}

// BuildExternalExtensionClientSet creates the api extension clientset of kubernetes by given cluster
// object and returns it.
func BuildExternalExtensionClientSet(cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*apiextensionsclient.Clientset, error) {
	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	return BuildExternalExtensionClientSetNoStatus(cluster, client)
}

// BuildExternalMonitoringClientSetNoStatus creates the monitoring clientset of prometheus operator by given
// cluster object and returns it.
func BuildExternalMonitoringClientSetNoStatus(cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (monitoringclient.Interface, error) {
	credential, err := ClusterCredentialV1(client, cluster.Name)
	if err != nil {
		return nil, err
	}
	restConfig, err := GetRestConfig(cluster, credential)
	if err != nil {
		log.Error("Build cluster config error", log.String("clusterName", cluster.ObjectMeta.Name), log.Err(err))
		return nil, err
	}
	return monitoringclient.NewForConfig(restConfig)
}

// BuildExternalMonitoringClientSet creates the monitoring clientset of  prometheus operator by given cluster
// object and returns it.
func BuildExternalMonitoringClientSet(cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (monitoringclient.Interface, error) {
	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	return BuildExternalMonitoringClientSetNoStatus(cluster, client)
}

// BuildExternalMonitoringClientSetWithName creates the clientset of prometheus operator by given cluster
// name and returns it.
func BuildExternalMonitoringClientSetWithName(platformClient platformversionedclient.PlatformV1Interface, name string) (monitoringclient.Interface, error) {
	cluster, err := platformClient.Clusters().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	clientset, err := BuildExternalMonitoringClientSet(cluster, platformClient)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// BuildExternalDynamicClientSetNoStatus creates the dynamic clientset of kubernetes by given
// cluster object and returns it.
func BuildExternalDynamicClientSetNoStatus(cluster *platformv1.Cluster, credential *platformv1.ClusterCredential) (dynamic.Interface, error) {
	restConfig, err := GetRestConfig(cluster, credential)
	if err != nil {
		log.Error("Build cluster config error", log.String("clusterName", cluster.ObjectMeta.Name), log.Err(err))
		return nil, err
	}
	return dynamic.NewForConfig(restConfig)
}

// BuildExternalDynamicClientSet creates the dynamic clientset of kubernetes by given cluster
// object and returns it.
func BuildExternalDynamicClientSet(cluster *platformv1.Cluster, credential *platformv1.ClusterCredential) (dynamic.Interface, error) {
	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	return BuildExternalDynamicClientSetNoStatus(cluster, credential)
}

// BuildClientSet creates client based on cluster information and returns it.
func BuildClientSet(cluster *platform.Cluster, credential *platform.ClusterCredential) (*kubernetes.Clientset, error) {
	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}
	address, err := ClusterAddress(cluster)
	if err != nil {
		return nil, err
	}
	config := api.NewConfig()
	config.CurrentContext = contextName

	if credential.CACert == nil {
		config.Clusters[contextName] = &api.Cluster{
			Server:                "https://" + address,
			InsecureSkipTLSVerify: true,
		}
	} else {
		config.Clusters[contextName] = &api.Cluster{
			Server:                   "https://" + address,
			CertificateAuthorityData: credential.CACert,
		}
	}

	if credential.Token != nil {
		config.AuthInfos[contextName] = &api.AuthInfo{
			Token: *credential.Token,
		}
	} else if credential.ClientCert != nil && credential.ClientKey != nil {
		config.AuthInfos[contextName] = &api.AuthInfo{
			ClientCertificateData: credential.ClientCert,
			ClientKeyData:         credential.ClientKey,
		}
	} else {
		return nil, fmt.Errorf("no credential for the cluster")
	}

	config.Contexts[contextName] = &api.Context{
		Cluster:  contextName,
		AuthInfo: contextName,
	}
	clientConfig := clientcmd.NewNonInteractiveClientConfig(*config, contextName, &clientcmd.ConfigOverrides{Timeout: "5s"}, nil)
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		log.Error("Build cluster config error", log.String("clusterName", cluster.ObjectMeta.Name), log.Err(err))
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

// ClusterAddress returns the cluster address.
func ClusterAddress(cluster *platform.Cluster) (string, error) {
	addrs := make(map[platform.AddressType][]platform.ClusterAddress)
	for _, one := range cluster.Status.Addresses {
		addrs[one.Type] = append(addrs[one.Type], one)
	}

	var address *platform.ClusterAddress
	if cluster.Spec.Type == platform.ClusterEKSHosting {
		if len(addrs[platform.AddressInternal]) != 0 {
			address = &addrs[platform.AddressInternal][rand.Intn(len(addrs[platform.AddressInternal]))]
		}
	} else {
		if len(addrs[platform.AddressAdvertise]) != 0 {
			address = &addrs[platform.AddressAdvertise][rand.Intn(len(addrs[platform.AddressAdvertise]))]
		} else {
			if len(addrs[platform.AddressReal]) != 0 {
				address = &addrs[platform.AddressReal][rand.Intn(len(addrs[platform.AddressReal]))]
			}
		}
	}

	if address == nil {
		return "", pkgerrors.New("no valid address for the cluster")
	}

	return fmt.Sprintf("%s:%d", address.Host, address.Port), nil
}

// ClusterV1Address returns the cluster address.
func ClusterV1Address(c *platformv1.Cluster) (string, error) {
	var cluster platform.Cluster
	err := platformv1.Convert_v1_Cluster_To_platform_Cluster(c, &cluster, nil)
	if err != nil {
		return "", pkgerrors.Wrap(err, "Convert_v1_Cluster_To_platform_Cluster errror")
	}
	return ClusterAddress(&cluster)
}

// rootCertPool returns nil if caData is empty.  When passed along, this will mean "use system CAs".
// When caData is not empty, it will be the ONLY information used in the CertPool.
func rootCertPool(caData []byte) *x509.CertPool {
	// What we really want is a copy of x509.systemRootsPool, but that isn't exposed.  It's difficult to build (see the go
	// code for a look at the platform specific insanity), so we'll use the fact that RootCAs == nil gives us the system values
	// It doesn't allow trusting either/or, but hopefully that won't be an issue
	if len(caData) == 0 {
		return nil
	}

	// if we have caData, use it
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caData)
	return certPool
}

// CheckClusterHealthzWithTimeout check cluster status within timeout
func CheckClusterHealthzWithTimeout(platformClient platformversionedclient.PlatformV1Interface, name string, timeout time.Duration) error {
	err := wait.PollImmediate(1*time.Second, timeout, func() (bool, error) {
		clientset, err := BuildExternalClientSetWithName(platformClient, name)
		if err != nil {
			return false, nil
		}
		healthStatus := 0
		clientset.Discovery().RESTClient().Get().AbsPath("/healthz").Do().StatusCode(&healthStatus)
		if healthStatus != http.StatusOK {
			return false, nil
		}
		return true, nil
	})

	return err
}
