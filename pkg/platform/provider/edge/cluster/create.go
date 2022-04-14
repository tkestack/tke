package cluster

import (
	"context"
	"fmt"
	superedgecommon "github.com/superedge/superedge/pkg/edgeadm/common"
	"io/ioutil"
	platformv1 "tkestack.io/tke/api/platform/v1"
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
	// create edge cluster kubeconfig
	kubeadmConfig := p.bCluster.GetKubeadmInitConfig(c)
	configData, err := kubeadmConfig.Marshal()
	if err != nil {
		return err
	}
	kubeconfigFile := fmt.Sprintf("/tmp/%s-kubeconfig", c.Name)
	err = ioutil.WriteFile(kubeconfigFile, []byte(configData), 0644)
	if err != nil {
		return err
	}

	// create edge cluster car key cart
	caKeyFile := fmt.Sprintf("/tmp/%s.key", c.Name)
	err = ioutil.WriteFile(caKeyFile, c.ClusterCredential.CAKey, 0644)
	if err != nil {
		return err
	}

	caCertFile := fmt.Sprintf("/tmp/%s.crt", c.Name)
	err = ioutil.WriteFile(caCertFile, c.ClusterCredential.CACert, 0644)
	if err != nil {
		return err
	}

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
