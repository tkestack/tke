package cluster

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/superedge/edgeadm/pkg/edgeadm/cmd"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kuberuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog"
	kubeadmscheme "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
	"os"
	"strings"
	platformv1 "tkestack.io/tke/api/platform/v1"
	v1 "tkestack.io/tke/pkg/platform/types/v1"

	superedgecommon "github.com/superedge/edgeadm/pkg/edgeadm/common"
	"github.com/superedge/edgeadm/pkg/edgeadm/constant"
	"github.com/superedge/edgeadm/pkg/edgeadm/steps"
	"github.com/superedge/edgeadm/pkg/util"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/util/log"
)

const (
	SuperEdgeRepo = "superedge.tencentcloudcr.com/superedge"
	EgressYaml    = `
apiVersion: apiserver.k8s.io/v1beta1
kind: EgressSelectorConfiguration
egressSelections:
- name: cluster
  connection:
    proxyProtocol: HTTPConnect
    transport:
      tcp:
        url: https://tunnel-cloud.edge-system.svc.cluster.local:8000
        tlsConfig:
          caBundle: /etc/kubernetes/pki/ca.crt
          clientCert: /etc/kubernetes/pki/tunnel-anp-client.crt
          clientKey: /etc/kubernetes/pki/tunnel-anp-client.key
`
)
const (
	EdgeImageRepository = "superedge.io/edgeImageResository"
	EdgeVersion         = "superedge.io/edge-version"
	EdgeVirtualAddr     = "superedge.io/edge-virtual-addr"
)

func (p *Provider) EnsureEdgeFlannel(ctx context.Context, c *v1.Cluster) error {
	edgeConf := &cmd.EdgeadmConfig{}
	edgeConf.ManifestsDir = ""

	cfg := &kubeadmapi.InitConfiguration{}
	cfg.ImageRepository = c.Annotations[EdgeImageRepository]
	cfg.Networking.PodSubnet = c.Spec.ClusterCIDR

	clientSet, err := c.Clientset()
	if err != nil {
		return err
	}
	// Deploy edge flannel
	return steps.EnsureFlannelAddon(cfg, edgeConf, clientSet)
}

func (p *Provider) EnsurePrepareEgdeCluster(ctx context.Context, c *v1.Cluster) error {
	// prepare node delay domain
	nodeDelayDomain := ""
	nodeDelayDomains := []string{constants.APIServerHostName}
	for _, domain := range nodeDelayDomains {
		nodeDelayDomain += fmt.Sprintf("%s\n", domain)
	}

	// prepare node hosts config
	nodeDomains := []string{
		p.bconfig.Registry.Domain,
		c.Cluster.Spec.TenantID + "." + p.bconfig.Registry.Domain,
	}
	hostsConfig := ""
	for _, one := range nodeDomains {
		hostsConfig += fmt.Sprintf("%s %s\n", p.bconfig.Registry.IP, one)
	}

	// prepare insecure registry config
	insecureRegistries := ""
	for _, registrie := range nodeDomains {
		insecureRegistries += fmt.Sprintf("%s\n", registrie)
	}

	// create edge-info configMap
	edgeInfoCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: constant.EdgeCertCM,
		},
		Data: map[string]string{
			constant.EdgeNodeHostConfig:  hostsConfig,
			constant.EdgeNodeDelayDomain: nodeDelayDomain,
			constant.InsecureRegistries:  insecureRegistries,
		},
	}
	clientSet, err := c.Clientset()
	if err != nil {
		return err
	}
	if err := superedgecommon.EnsureEdgeSystemNamespace(clientSet); err != nil {
		return err
	}
	alreayEdgeInfoCM, err := clientSet.CoreV1().ConfigMaps(constant.NamespaceEdgeSystem).
		Get(context.TODO(), constant.EdgeCertCM, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			cm, err := clientSet.CoreV1().ConfigMaps(
				constant.NamespaceEdgeSystem).Create(context.TODO(), edgeInfoCM, metav1.CreateOptions{})
			if err != nil {
				return err
			}
			log.Infof("Create success configMap: %v", constant.EdgeNodeHostConfig, util.ToJson(cm))
			return nil
		}
		return err
	}

	alreayEdgeInfoCM.Data[constant.EdgeNodeHostConfig] = hostsConfig
	alreayEdgeInfoCM.Data[constant.EdgeNodeHostConfig] = hostsConfig
	alreayEdgeInfoCM.Data[constant.EdgeNodeHostConfig] = hostsConfig
	if _, err := clientSet.CoreV1().ConfigMaps(constant.NamespaceEdgeSystem).
		Update(context.TODO(), alreayEdgeInfoCM, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureApplyEdgeApps(ctx context.Context, c *v1.Cluster) error {
	// get kube-apiserver ip
	apiserverIP := c.Spec.Machines[0].IP
	if len(c.Spec.PublicAlternativeNames) > 0 {
		apiserverIP = c.Spec.PublicAlternativeNames[0]
	}
	if c.Spec.Features.HA != nil {
		if c.Spec.Features.HA.TKEHA != nil {
			apiserverIP = c.Spec.Features.HA.TKEHA.VIP
		}
		if c.Spec.Features.HA.ThirdPartyHA != nil {
			apiserverIP = c.Spec.Features.HA.ThirdPartyHA.VIP
		}
	}

	os.MkdirAll(fmt.Sprintf("/tmp/%s", c.Name), os.ModePerm)
	// create edge cluster car key cart
	caKeyFile := fmt.Sprintf("/tmp/%s/ca.key", c.Name)
	err := ioutil.WriteFile(caKeyFile, c.ClusterCredential.CAKey, 0644)
	if err != nil {
		return err
	}

	caCertFile := fmt.Sprintf("/tmp/%s/ca.crt", c.Name)
	err = ioutil.WriteFile(caCertFile, c.ClusterCredential.CACert, 0644)
	if err != nil {
		return err
	}

	certSANs := []string{apiserverIP}
	for _, machine := range c.Spec.Machines {
		certSANs = append(certSANs, machine.IP)
	}
	if len(c.Spec.PublicAlternativeNames) > 0 {
		certSANs = append(certSANs, c.Spec.PublicAlternativeNames...)
	}

	// deploy superedge edge cluster apps
	clientset, err := c.Clientset()
	if err != nil {
		return err
	}
	edgeConf := &cmd.EdgeadmConfig{}
	edgeConf.ManifestsDir = ""
	edgeConf.TunnelCloudToken = util.GetRandToken(32)
	edgeConf.Version = c.Annotations[EdgeVersion]
	edgeConf.EdgeImageRepository = c.Annotations[EdgeImageRepository]

	virtualAddr, ok := c.Annotations[EdgeVirtualAddr]
	if ok {
		edgeConf.EdgeVirtualAddr = virtualAddr
	} else {
		edgeConf.EdgeVirtualAddr = constant.DefaultEdgeVirtualAddr
	}

	cfg := &kubeadmapi.InitConfiguration{}
	cfg.APIServer.CertSANs = certSANs
	cfg.CertificatesDir = fmt.Sprintf("/tmp/%s/", c.Name)
	cfg.ControlPlaneEndpoint = apiserverIP
	cfg.NodeRegistration.Name = c.Spec.Machines[0].IP
	version, err := kubeadmutil.KubernetesReleaseVersion(strings.Split(c.Spec.Version, "-")[0])
	if err != nil {
		klog.Errorf("Failed to get k8s version, cluster: %s, error: %v", c.Name, err)
		return err
	}
	cfg.KubernetesVersion = version

	err = steps.EnsureServiceGroupAddon(cfg, edgeConf, clientset)
	if err != nil {
		klog.Errorf("Failed to install ServiceGroup, cluster: %s, error: %v", c.Name, err)
		return err
	}

	steps.EnsureTunnelAddon(cfg, edgeConf, clientset)
	if err != nil {
		klog.Errorf("Failed to install ServiceGroup, cluster: %s, error: %v", c.Name, err)
		return err
	}

	steps.EnsureEdgeHealthAddon(cfg, edgeConf, clientset)
	if err != nil {
		klog.Errorf("Failed to install EdgeHealth, cluster: %s, error: %v", c.Name, err)
		return err
	}

	steps.EnsureEdgeCorednsAddon(cfg, edgeConf, clientset)
	if err != nil {
		klog.Errorf("Failed to install EdgeCoredns, cluster: %s, error: %v", c.Name, err)
		return err
	}

	steps.EnsureNodePrepare(cfg, edgeConf, clientset)
	if err != nil {
		klog.Errorf("Failed to install NodePrepar, cluster: %s, error: %v", c.Name, err)
		return err
	}

	steps.EnsureEdgeKubeConfig(cfg, edgeConf, clientset)
	if err != nil {
		klog.Errorf("Failed to install EdgeKubeConfig, cluster: %s, error: %v", c.Name, err)
		return err
	}
	return nil
}

func (p *Provider) EnsureKubeadmConfig(ctx context.Context, c *v1.Cluster) error {
	client, err := c.Clientset()
	if err != nil {
		klog.Errorf("Failed to get clientSet, cluster: %s, error: %v", c.Name, err)
		return err
	}
	cm, err := client.CoreV1().ConfigMaps(constant.NamespaceKubeSystem).Get(ctx, kubeadmconstants.KubeadmConfigConfigMap, metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		klog.Errorf("Failed to get configMap: %s, cluster: %s, error: %v", kubeadmconstants.KubeadmConfigConfigMap, c.Name, err)
		return err
	}
	clusterConfig, ok := cm.Data[kubeadmconstants.ClusterConfigurationConfigMapKey]
	if !ok {
		return fmt.Errorf("Fialed to get %s, cluster: %s ", kubeadmconstants.ClusterConfigurationConfigMapKey, c.Name)
	}

	f := bytes.NewBuffer([]byte(clusterConfig))
	d := yamlutil.NewYAMLOrJSONDecoder(f, 4096)
	ext := kuberuntime.RawExtension{}
	if err := d.Decode(&ext); err != nil {
		if err != io.EOF {
			return err
		}
	}
	obj, _, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, nil)
	if err != nil {
		return err
	}
	unstructuredMap, err := kuberuntime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return err
	}
	unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
	edgeRepo, ok := c.Annotations[EdgeImageRepository]
	if !ok {
		edgeRepo = EdgeImageRepository
	}
	unstructured.SetNestedField(unstructuredObj.Object, edgeRepo, "edgeImageRepository")

	clusterConfigByte, err := kubeadmutil.MarshalToYamlForCodecs(obj, schema.GroupVersion{
		Group:   strings.Split(unstructuredObj.GetAPIVersion(), "/")[0],
		Version: strings.Split(unstructuredObj.GetAPIVersion(), "/")[1],
	}, kubeadmscheme.Codecs)
	if err != nil {
		return err
	}
	cm.Data[kubeadmconstants.ClusterConfigurationConfigMapKey] = string(clusterConfigByte)
	cm.ResourceVersion = ""

	err = apiclient.CreateOrMutateConfigMap(client, cm, func(cm *corev1.ConfigMap) error {
		// Upgrade will call to UploadConfiguration with a modified KubernetesVersion reflecting the new
		// Kubernetes version. In that case, the mutation path will take place.
		cm.Data[kubeadmconstants.ClusterConfigurationConfigMapKey] = string(clusterConfigByte)
		return nil
	})
	if err != nil {
		return err
	}
	return err
}

func (p *Provider) EnsureEgressSelector(ctx context.Context, c *v1.Cluster) error {
	crt, err := tls.X509KeyPair(c.ClusterCredential.CACert, c.ClusterCredential.CAKey)
	if err != nil {
		return err
	}
	certs, err := x509.ParseCertificates(crt.Certificate[0])
	if err != nil {
		return err
	}
	crt.Leaf = certs[0]

	clientCrt, clientKey, err := util.GenerateClientCertAndKey(crt.Leaf, crt.PrivateKey.(*rsa.PrivateKey), "tunnel-anp")
	if err != nil {
		return err
	}
	keydata := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY", //tobe replaced
		Bytes: x509.MarshalPKCS1PrivateKey(clientKey),
	})
	crtdata := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE", //tobe replaced
		Bytes: clientCrt.Raw,
	})
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}
		err = machineSSH.WriteFile(bytes.NewReader(crtdata), "/etc/kubernetes/pki/tunnel-anp-client.crt")
		if err != nil {
			return err
		}

		err = machineSSH.WriteFile(bytes.NewReader(keydata), "/etc/kubernetes/pki/tunnel-anp-client.key")
		if err != nil {
			return err
		}

		err = machineSSH.WriteFile(bytes.NewReader([]byte(EgressYaml)), "/etc/kubernetes/egress-selector-configuration.yaml")
		if err != nil {
			return err
		}
		pod, err := machineSSH.ReadFile("/etc/kubernetes/manifests/kube-apiserver.yaml")
		if err != nil {
			return err
		}
		apiserver := &corev1.Pod{}
		_, _, err = clientsetscheme.Codecs.UniversalDeserializer().Decode(pod, nil, apiserver)
		if err != nil {
			return err
		}
		for k, v := range apiserver.Spec.Containers {
			if v.Name == "kube-apiserver" {
				v.Command = append(v.Command, "--egress-selector-config-file=/etc/kubernetes/egress-selector-configuration.yaml")
				apiserver.Spec.Containers[k] = v
			}
		}
		apiserver.Spec.DNSPolicy = corev1.DNSClusterFirstWithHostNet
		serialized, err := kubeadmutil.MarshalToYaml(apiserver, corev1.SchemeGroupVersion)
		if err != nil {
			return err
		}

		err = machineSSH.WriteFile(bytes.NewReader(serialized), "/etc/kubernetes/manifests/kube-apiserver.yaml")
		if err != nil {
			return err
		}

	}
	return nil
}
