package cluster

import (
	"context"

	superedgecommon "github.com/superedge/superedge/pkg/edgeadm/common"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/util/mark"
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
	clientset, err := c.Clientset()
	if err != nil {
		return err
	}

	superedgecommon.DeployEdgeAPPS(clientset, "manifestDir", "caCertFile", "caKeyFile", "masterPublicAddr", "certSANs", "kubeConfig")
	return mark.Create(ctx, clientset)
}
