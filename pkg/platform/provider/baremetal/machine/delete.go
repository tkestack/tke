package machine

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	platformv1 "tkestack.io/tke/api/platform/v1"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/log"
)

func (p *Provider) EnsureRemoveNode(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	log.FromContext(ctx).Info("deleteNode doing")

	if cluster.Status.Phase == platformv1.ClusterTerminating {
		return nil
	}

	clientset, err := cluster.Clientset()
	if err != nil {
		return err
	}

	node, err := apiclient.GetNodeByMachineIP(ctx, clientset, machine.Spec.IP)
	if err != nil {
		return err
	}
	err = clientset.CoreV1().Nodes().Delete(ctx, node.Name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = wait.PollImmediate(5*time.Second, 5*time.Minute, waitForNodeDelete(ctx, clientset, node.Name)); err != nil {
		return err
	}

	log.FromContext(ctx).Info("deleteNode done")
	return nil
}

func waitForNodeDelete(ctx context.Context, c kubernetes.Interface, nodeName string) wait.ConditionFunc {
	return func() (done bool, err error) {
		if _, err := c.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{}); err != nil {
			if errors.IsNotFound(err) {
				return true, nil
			}
		}

		return false, nil
	}
}
