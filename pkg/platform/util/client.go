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
	"math/rand"
	"net"
	"net/http"
	"path"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/fields"

	monitoringclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	pkgerrors "github.com/pkg/errors"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
)

const (
	contextName = "tke"
	clientQPS   = 100
	clientBurst = 200
)

func DynamicClientByCluster(ctx context.Context, cluster *platform.Cluster, platformClient platforminternalclient.PlatformInterface) (dynamic.Interface, error) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) > 0 && cluster.Spec.TenantID != tenantID {
		return nil, errors.NewNotFound(platform.Resource("clusters"), cluster.ObjectMeta.Name)
	}

	credential, err := GetClusterCredential(ctx, platformClient, cluster)
	if err != nil {
		return nil, err
	}

	return BuildInternalDynamicClientSet(cluster, credential)
}

// ClientSetByCluster returns the backend kubernetes clientSet by given cluster object
func ClientSetByCluster(ctx context.Context, cluster *platform.Cluster, platformClient platforminternalclient.PlatformInterface) (*kubernetes.Clientset, error) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) > 0 && cluster.Spec.TenantID != tenantID {
		return nil, errors.NewNotFound(platform.Resource("clusters"), cluster.ObjectMeta.Name)
	}
	credential, err := GetClusterCredential(ctx, platformClient, cluster)
	if err != nil {
		return nil, err
	}

	return BuildClientSet(ctx, cluster, credential)
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

// GetExternalRestConfig returns rest config according to cluster
func GetExternalRestConfig(cluster *platformv1.Cluster, credential *platformv1.ClusterCredential) (*restclient.Config, error) {
	host, err := ClusterV1Host(cluster)
	if err != nil {
		return nil, err
	}
	config := api.NewConfig()
	config.CurrentContext = contextName

	if credential.CACert == nil {
		config.Clusters[contextName] = &api.Cluster{
			Server:                fmt.Sprintf("https://%s", host),
			InsecureSkipTLSVerify: true,
		}
	} else {
		config.Clusters[contextName] = &api.Cluster{
			Server:                   fmt.Sprintf("https://%s", host),
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

// GetInternalRestConfig returns rest config according to cluster
func GetInternalRestConfig(cluster *platform.Cluster, credential *platform.ClusterCredential) (*restclient.Config, error) {
	host, err := ClusterHost(cluster)
	if err != nil {
		return nil, err
	}
	config := api.NewConfig()
	config.CurrentContext = contextName

	if credential.CACert == nil {
		config.Clusters[contextName] = &api.Cluster{
			Server:                fmt.Sprintf("https://%s", host),
			InsecureSkipTLSVerify: true,
		}
	} else {
		config.Clusters[contextName] = &api.Cluster{
			Server:                   fmt.Sprintf("https://%s", host),
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

// BuildInternalDynamicClientSet creates the dynamic clientset of kubernetes by given cluster
// object and returns it.
func BuildInternalDynamicClientSet(cluster *platform.Cluster, credential *platform.ClusterCredential) (dynamic.Interface, error) {
	if cluster.Status.Phase != platform.ClusterRunning {
		return nil, fmt.Errorf("cluster %s status is abnormal", cluster.ObjectMeta.Name)
	}

	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	return BuildInternalDynamicClientSetNoStatus(cluster, credential)
}

// BuildInternalDynamicClientSetNoStatus creates the dynamic clientset of kubernetes by given
// cluster object and returns it.
func BuildInternalDynamicClientSetNoStatus(cluster *platform.Cluster, credential *platform.ClusterCredential) (dynamic.Interface, error) {
	restConfig, err := GetInternalRestConfig(cluster, credential)
	if err != nil {
		log.Error("Build cluster config error", log.String("clusterName", cluster.ObjectMeta.Name), log.Err(err))
		return nil, err
	}
	return dynamic.NewForConfig(restConfig)
}

// BuildVersionedClientSet creates the clientset of kubernetes by given
// cluster object and returns it.
func BuildVersionedClientSet(cluster *platformv1.Cluster, credential *platformv1.ClusterCredential) (*kubernetes.Clientset, error) {
	restConfig, err := GetExternalRestConfig(cluster, credential)
	if err != nil {
		log.Error("Build cluster config error", log.String("clusterName", cluster.ObjectMeta.Name), log.Err(err))
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

// BuildExternalClientSet creates the clientset of kubernetes by given cluster
// object and returns it.
func BuildExternalClientSet(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*kubernetes.Clientset, error) {
	credential, err := GetClusterCredentialV1(ctx, client, cluster)
	if err != nil {
		return nil, err
	}

	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	return BuildVersionedClientSet(cluster, credential)
}

// BuildExternalClientSetWithName creates the clientset of kubernetes by given cluster
// name and returns it.
func BuildExternalClientSetWithName(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, name string) (*kubernetes.Clientset, error) {
	cluster, err := platformClient.Clusters().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	clientset, err := BuildExternalClientSet(ctx, cluster, platformClient)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// BuildExternalExtensionClientSetNoStatus creates the api extension clientset of kubernetes by given
// cluster object and returns it.
func BuildExternalExtensionClientSetNoStatus(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*apiextensionsclient.Clientset, error) {
	credential, err := GetClusterCredentialV1(ctx, client, cluster)
	if err != nil {
		return nil, err
	}
	restConfig, err := GetExternalRestConfig(cluster, credential)
	if err != nil {
		log.Error("Build cluster config error", log.String("clusterName", cluster.ObjectMeta.Name), log.Err(err))
		return nil, err
	}
	return apiextensionsclient.NewForConfig(restConfig)
}

// BuildExternalExtensionClientSet creates the api extension clientset of kubernetes by given cluster
// object and returns it.
func BuildExternalExtensionClientSet(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*apiextensionsclient.Clientset, error) {
	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	return BuildExternalExtensionClientSetNoStatus(ctx, cluster, client)
}

// BuildExternalMonitoringClientSetNoStatus creates the monitoring clientset of prometheus operator by given
// cluster object and returns it.
func BuildExternalMonitoringClientSetNoStatus(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (monitoringclient.Interface, error) {
	credential, err := GetClusterCredentialV1(ctx, client, cluster)
	if err != nil {
		return nil, err
	}
	restConfig, err := GetExternalRestConfig(cluster, credential)
	if err != nil {
		log.Error("Build cluster config error", log.String("clusterName", cluster.ObjectMeta.Name), log.Err(err))
		return nil, err
	}
	return monitoringclient.NewForConfig(restConfig)
}

// BuildExternalMonitoringClientSet creates the monitoring clientset of  prometheus operator by given cluster
// object and returns it.
func BuildExternalMonitoringClientSet(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (monitoringclient.Interface, error) {
	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	return BuildExternalMonitoringClientSetNoStatus(ctx, cluster, client)
}

// BuildExternalMonitoringClientSetWithName creates the clientset of prometheus operator by given cluster
// name and returns it.
func BuildExternalMonitoringClientSetWithName(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, name string) (monitoringclient.Interface, error) {
	cluster, err := platformClient.Clusters().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	clientset, err := BuildExternalMonitoringClientSet(ctx, cluster, platformClient)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// BuildExternalDynamicClientSetNoStatus creates the dynamic clientset of kubernetes by given
// cluster object and returns it.
func BuildExternalDynamicClientSetNoStatus(cluster *platformv1.Cluster, credential *platformv1.ClusterCredential) (dynamic.Interface, error) {
	restConfig, err := GetExternalRestConfig(cluster, credential)
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
func BuildClientSet(ctx context.Context, cluster *platform.Cluster, credential *platform.ClusterCredential) (*kubernetes.Clientset, error) {
	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}
	host, err := ClusterHost(cluster)
	if err != nil {
		return nil, err
	}
	config := api.NewConfig()
	config.CurrentContext = contextName

	if credential.CACert == nil {
		config.Clusters[contextName] = &api.Cluster{
			Server:                fmt.Sprintf("https://%s", host),
			InsecureSkipTLSVerify: true,
		}
	} else {
		config.Clusters[contextName] = &api.Cluster{
			Server:                   fmt.Sprintf("https://%s", host),
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
	clientConfig := clientcmd.NewNonInteractiveClientConfig(*config, contextName, &clientcmd.ConfigOverrides{Timeout: "30s"}, nil)
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		log.Error("Build cluster config error", log.String("clusterName", cluster.ObjectMeta.Name), log.Err(err))
		return nil, err
	}
	restConfig.QPS = clientQPS
	restConfig.Burst = clientBurst
	return kubernetes.NewForConfig(restConfig)
}

// ClusterHost returns host and port for kube-apiserver of cluster.
func ClusterHost(cluster *platform.Cluster) (string, error) {
	address, err := ClusterAddress(cluster)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("%s:%d", address.Host, address.Port)
	if address.Path != "" {
		result = path.Join(result, address.Path)
	}

	return result, nil
}

func ClusterAddress(cluster *platform.Cluster) (*platform.ClusterAddress, error) {
	addrs := make(map[platform.AddressType][]platform.ClusterAddress)
	for _, one := range cluster.Status.Addresses {
		addrs[one.Type] = append(addrs[one.Type], one)
	}

	var address *platform.ClusterAddress
	if len(addrs[platform.AddressInternal]) != 0 {
		address = &addrs[platform.AddressInternal][rand.Intn(len(addrs[platform.AddressInternal]))]
	} else if len(addrs[platform.AddressAdvertise]) != 0 {
		address = &addrs[platform.AddressAdvertise][rand.Intn(len(addrs[platform.AddressAdvertise]))]
	} else {
		if len(addrs[platform.AddressReal]) != 0 {
			address = &addrs[platform.AddressReal][rand.Intn(len(addrs[platform.AddressReal]))]
		}
	}
	if address == nil {
		return nil, pkgerrors.New("no valid address for the cluster")
	}

	return address, nil
}

// ClusterV1Host returns host and port for kube-apiserver of versioned cluster resource.
func ClusterV1Host(c *platformv1.Cluster) (string, error) {
	var cluster platform.Cluster
	err := platformv1.Convert_v1_Cluster_To_platform_Cluster(c, &cluster, nil)
	if err != nil {
		return "", pkgerrors.Wrap(err, "Convert_v1_Cluster_To_platform_Cluster errror")
	}
	return ClusterHost(&cluster)
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
func CheckClusterHealthzWithTimeout(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, name string, timeout time.Duration) error {
	err := wait.PollImmediate(1*time.Second, timeout, func() (bool, error) {
		clientset, err := BuildExternalClientSetWithName(ctx, platformClient, name)
		if err != nil {
			return false, nil
		}
		healthStatus := 0
		clientset.Discovery().RESTClient().Get().AbsPath("/healthz").Do(ctx).StatusCode(&healthStatus)
		if healthStatus != http.StatusOK {
			return false, nil
		}
		return true, nil
	})

	return err
}

// GetClusterCredential returns the cluster's credential
func GetClusterCredential(ctx context.Context, client platforminternalclient.PlatformInterface, cluster *platform.Cluster) (*platform.ClusterCredential, error) {
	var (
		credential *platform.ClusterCredential
		err        error
	)

	if cluster.Spec.ClusterCredentialRef != nil {
		credential, err = client.ClusterCredentials().Get(ctx, cluster.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
	} else {
		clusterName := cluster.Name
		fieldSelector := fields.OneTermEqualSelector("clusterName", clusterName).String()
		clusterCredentials, err := client.ClusterCredentials().List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
		if err != nil {
			return nil, err
		}
		if len(clusterCredentials.Items) == 0 {
			return nil, errors.NewNotFound(platform.Resource("ClusterCredential"), clusterName)
		}
		credential = &clusterCredentials.Items[0]
	}

	return credential, nil
}

// GetClusterCredentialV1 returns the versioned cluster's credential
func GetClusterCredentialV1(ctx context.Context, client platformversionedclient.PlatformV1Interface, cluster *platformv1.Cluster) (*platformv1.ClusterCredential, error) {
	var (
		credential *platformv1.ClusterCredential
		err        error
	)

	if cluster.Spec.ClusterCredentialRef != nil {
		credential, err = client.ClusterCredentials().Get(ctx, cluster.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
	} else {
		clusterName := cluster.Name
		fieldSelector := fields.OneTermEqualSelector("clusterName", clusterName).String()
		clusterCredentials, err := client.ClusterCredentials().List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
		if err != nil {
			return nil, err
		}
		if len(clusterCredentials.Items) == 0 {
			return nil, errors.NewNotFound(platform.Resource("ClusterCredential"), clusterName)
		}
		credential = &clusterCredentials.Items[0]
	}

	return credential, nil
}
