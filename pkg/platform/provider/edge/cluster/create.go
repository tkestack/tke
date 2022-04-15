package cluster

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	superedgecommon "github.com/superedge/superedge/pkg/edgeadm/common"
	"k8s.io/apimachinery/pkg/runtime"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"

	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeconfig"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
)

func (p *Provider) EnsurePrepareEgdeCluster(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]

	for _, machine := range machines {
		_, err := machine.SSH()
		if err != nil {
			return err
		}
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
	masterPublicAddr := apiserverIP

	// create edge cluster kubeconfig
	kubeAPIAddr := fmt.Sprintf("https://%s:6443", masterPublicAddr)
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

	certSANs := []string{masterPublicAddr}
	for _, machine := range c.Spec.Machines {
		certSANs = append(certSANs, machine.IP)
	}

	// deploy superedge edge cluster apps
	clientset, err := c.Clientset()
	if err != nil {
		return err
	}
	err = superedgecommon.DeployEdgeAPPS(clientset, "", caCertFile, caKeyFile, masterPublicAddr, certSANs, kubeconfigFile)
	if err != nil {
		return err
	}
	return nil
}
