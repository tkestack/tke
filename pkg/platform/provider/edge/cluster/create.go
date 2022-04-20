package cluster

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	superedgecommon "github.com/superedge/superedge/pkg/edgeadm/common"
	"github.com/superedge/superedge/pkg/edgeadm/constant"
	"github.com/superedge/superedge/pkg/util"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"

	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeconfig"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/log"
)

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
	if _, err := clientSet.CoreV1().ConfigMaps(constant.NamespaceEdgeSystem).
		Get(context.TODO(), constant.EdgeCertCM, metav1.GetOptions{}); err != nil {
		if apierrors.IsNotFound(err) {
			cm, err := clientSet.CoreV1().ConfigMaps(
				constant.NamespaceEdgeSystem).Create(context.TODO(), edgeInfoCM, metav1.CreateOptions{})
			if err != nil {
				return err
			}
			log.Infof("Create success configMap: %v", constant.EdgeNodeHostConfig, util.ToJson(cm))
			return nil
		} else {
			return err
		}
	}

	if _, err := clientSet.CoreV1().ConfigMaps(constant.NamespaceEdgeSystem).
		Update(context.TODO(), edgeInfoCM, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureApplyEdgeApps(ctx context.Context, c *v1.Cluster) error {
	// get kube-apiserver ip
	apiserverIP := c.Spec.Machines[0].IP
	if c.Spec.Features.HA != nil {
		if c.Spec.Features.HA.TKEHA != nil {
			apiserverIP = c.Spec.Features.HA.TKEHA.VIP
		}
		if c.Spec.Features.HA.ThirdPartyHA != nil {
			apiserverIP = c.Spec.Features.HA.ThirdPartyHA.VIP
		}
	}

	// create edge cluster kubeconfig
	kubeAPIAddr := fmt.Sprintf("https://%s:6443", apiserverIP)
	config := kubeconfig.CreateWithToken(kubeAPIAddr, c.Name,
		"kubernetes-admin", c.ClusterCredential.CACert, *c.ClusterCredential.Token)
	configData, err := runtime.Encode(clientcmdlatest.Codec, config)
	if err != nil {
		return err
	}
	os.MkdirAll(fmt.Sprintf("/tmp/%s", c.Name), os.ModePerm)
	kubeconfigFile := fmt.Sprintf("/tmp/%s/%s-kubeconfig", c.Name, c.Name)
	err = ioutil.WriteFile(kubeconfigFile, configData, 0644)
	if err != nil {
		return err
	}

	// create edge cluster car key cart
	caKeyFile := fmt.Sprintf("/tmp/%s/%s.key", c.Name, c.Name)
	err = ioutil.WriteFile(caKeyFile, c.ClusterCredential.CAKey, 0644)
	if err != nil {
		return err
	}

	caCertFile := fmt.Sprintf("/tmp/%s/%s.crt", c.Name, c.Name)
	err = ioutil.WriteFile(caCertFile, c.ClusterCredential.CACert, 0644)
	if err != nil {
		return err
	}

	certSANs := []string{apiserverIP}
	for _, machine := range c.Spec.Machines {
		certSANs = append(certSANs, machine.IP)
	}

	// deploy superedge edge cluster apps
	clientset, err := c.Clientset()
	if err != nil {
		return err
	}
	err = superedgecommon.DeployEdgeAPPS(clientset, "", caCertFile, caKeyFile, apiserverIP, certSANs, kubeconfigFile)
	if err != nil {
		return err
	}
	return nil
}
