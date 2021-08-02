package machine

import (
	"context"

	platformv1 "tkestack.io/tke/api/platform/v1"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
)

func (p *Provider) EnsureRemoveNode(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	clientset, err := cluster.Clientset()
	if err != nil {
		return err
	}

	node, err := apiclient.GetNodeByMachineIP(ctx, clientset, machine.Spec.IP)
	if err != nil {
		return err
	}
	err = apiclient.MarkNode(ctx, clientset, node.Name, machine.Spec.Labels, machine.Spec.Taints)
	if err != nil {
		return err
	}
	return nil
}
