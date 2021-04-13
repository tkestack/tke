package tke

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types2 "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	tkeclientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1 "tkestack.io/tke/api/platform/v1"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/test/util"
	"tkestack.io/tke/test/util/cloudprovider"
	"tkestack.io/tke/test/util/env"
)

type TestTKE struct {
	TkeClient tkeclientset.Interface
	provider  cloudprovider.Provider
	Clusters  []*platformv1.Cluster
}

func Init(tkeClient tkeclientset.Interface, provider cloudprovider.Provider) *TestTKE {
	return &TestTKE{
		TkeClient: tkeClient,
		provider:  provider,
		Clusters:  []*platformv1.Cluster{},
	}
}

func (testTke *TestTKE) K8sClient(cls *platformv1.Cluster) kubernetes.Interface {
	clusterWrapper, err := typesv1.GetCluster(context.Background(), testTke.TkeClient.PlatformV1(), cls)
	if err != nil {
		panic(err)
	}
	client, err := clusterWrapper.Clientset()
	if err != nil {
		panic(err)
	}
	return client
}

func (testTke *TestTKE) CreateCluster() (cluster *platformv1.Cluster, err error) {
	cluster = testTke.ClusterTemplate()
	return testTke.CreateClusterInternal(cluster)
}

func (testTke *TestTKE) ClusterTemplate(nodes ...cloudprovider.Instance) *platformv1.Cluster {
	cluster := &platformv1.Cluster{
		Spec: platformv1.ClusterSpec{
			Type: "Baremetal",
			Features: platformv1.ClusterFeature{
				HA:                   &platformv1.HA{},
				EnableMasterSchedule: true,
			},
			Version:       env.K8sVersion(),
			ClusterCIDR:   "10.244.0.0/16",
			NetworkDevice: "eth0",
			Machines:      []platformv1.ClusterMachine{},
		}}
	if len(nodes) == 0 {
		var err error
		nodes, err = testTke.provider.CreateInstances(1)
		if err != nil {
			panic(fmt.Errorf("CreateInstance failed. %v", err))
		}
	}
	for _, one := range nodes {
		cluster.Spec.Machines = append(cluster.Spec.Machines, platformv1.ClusterMachine{
			IP:       one.InternalIP,
			Port:     one.Port,
			Username: one.Username,
			Password: []byte(one.Password),
		})
	}
	return cluster
}

func (testTke *TestTKE) CreateClusterInternal(cls *platformv1.Cluster) (cluster *platformv1.Cluster, err error) {
	klog.Info("Create cluster: ", cls.String())

	err = wait.PollImmediate(10*time.Second, time.Minute, func() (bool, error) {
		cluster, err = testTke.TkeClient.PlatformV1().Clusters().Create(context.Background(), cls, metav1.CreateOptions{})
		if err != nil {
			return false, err
		}
		return true, nil
	})

	klog.Info("Cluster name: ", cluster.Name)
	testTke.Clusters = append(testTke.Clusters, cluster)
	return testTke.WaitClusterToBeRunning(cluster.Name)
}

func (testTke *TestTKE) UpgradeCluster(clsName string, targetVersion string, mode platformv1.UpgradeMode, drain bool) (*platformv1.Cluster, error) {
	body := &platformv1.Cluster{
		Spec: platformv1.ClusterSpec{
			Type: "Baremetal",
			Features: platformv1.ClusterFeature{
				Upgrade: platformv1.Upgrade{
					Mode: mode,
					Strategy: platformv1.UpgradeStrategy{
						DrainNodeBeforeUpgrade: common.BoolPtr(drain),
					},
				},
			},
			Version: targetVersion,
		},
	}
	patch, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	cls, err := testTke.TkeClient.PlatformV1().Clusters().Patch(context.Background(), clsName, types2.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return cls, fmt.Errorf("upgrade cluster failed. %v", err)
	}
	return testTke.WaitClusterToBeRunning(clsName)
}

// ScaleUp note: nodes must share the same VIP with the existing nodes
func (testTke *TestTKE) ScaleUp(clsName string, nodes []cloudprovider.Instance) (cls *platformv1.Cluster, err error) {
	klog.Info("Scale up")
	cls, err = testTke.TkeClient.PlatformV1().Clusters().Get(context.Background(), clsName, metav1.GetOptions{})
	if err != nil {
		return
	}

	patch := &platformv1.Cluster{
		Spec: cls.Spec,
	}
	for _, node := range nodes {
		patch.Spec.Machines = append(patch.Spec.Machines, platformv1.ClusterMachine{
			IP:       node.InternalIP,
			Port:     node.Port,
			Username: node.Username,
			Password: []byte(node.Password),
		})
	}

	patchData, _ := json.Marshal(patch)
	cls, err = testTke.TkeClient.PlatformV1().Clusters().Patch(context.Background(), cls.Name, types2.StrategicMergePatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return
	}
	return testTke.WaitClusterToBeRunning(cls.Name)
}

func (testTke *TestTKE) ScaleDown(clsName string, ipsToBeRemoved []string) (cls *platformv1.Cluster, err error) {
	klog.Info("Scale down")
	cls, err = testTke.TkeClient.PlatformV1().Clusters().Get(context.Background(), clsName, metav1.GetOptions{})
	if err != nil {
		return
	}

	patch := &platformv1.Cluster{
		Spec: cls.Spec,
	}
	patch.Spec.Machines = []platformv1.ClusterMachine{}
	for _, node := range cls.Spec.Machines {
		if !util.Contains(ipsToBeRemoved, node.IP) {
			patch.Spec.Machines = append(patch.Spec.Machines, platformv1.ClusterMachine{
				IP:       node.IP,
				Port:     node.Port,
				Username: node.Username,
				Password: node.Password,
			})
		}
	}
	patchData, _ := json.Marshal(patch)
	cls, err = testTke.TkeClient.PlatformV1().Clusters().Patch(context.Background(), cls.Name, types2.StrategicMergePatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return
	}
	return testTke.WaitClusterToBeRunning(cls.Name)
}

func (testTke *TestTKE) WaitClusterToBeRunning(clusterName string) (cluster *platformv1.Cluster, err error) {
	klog.Info("Wait cluster status to be running")
	err = wait.Poll(5*time.Second, 10*time.Minute, func() (bool, error) {
		cluster, err = testTke.TkeClient.PlatformV1().Clusters().Get(context.Background(), clusterName, metav1.GetOptions{})
		if err != nil {
			klog.Error(err)
			return false, nil
		}
		if len(cluster.Status.Conditions) > 0 {
			lastCondition := cluster.Status.Conditions[len(cluster.Status.Conditions)-1]
			klog.Info("Phase: ", cluster.Status.Phase, ", Type: ", lastCondition.Type, ", message: ", lastCondition.Message)
		}
		return cluster.Status.Phase == platformv1.ClusterRunning, nil
	})
	if err != nil {
		return cluster, fmt.Errorf("wait cluster to be running failed. %v", err)
	}
	return
}

func (testTke *TestTKE) ImportCluster(host string, port int32, caCert []byte, token *string) (cluster *platformv1.Cluster, err error) {
	credential := &platformv1.ClusterCredential{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "clustercredential",
		},
		CACert: caCert,
		Token:  token,
	}
	credential, err = testTke.TkeClient.PlatformV1().ClusterCredentials().Create(context.Background(), credential, metav1.CreateOptions{})
	if err != nil {
		return
	}

	cluster = &platformv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "cls",
		},
		Spec: platformv1.ClusterSpec{
			//DisplayName: baremetalClusterName,
			Type: "Imported",
			ClusterCredentialRef: &corev1.LocalObjectReference{
				Name: credential.Name,
			},
		},
		Status: platformv1.ClusterStatus{
			Addresses: []platformv1.ClusterAddress{
				{
					Host: host,
					Path: "",
					Port: port,
					Type: platformv1.AddressAdvertise,
				},
			},
		},
	}
	cluster, err = testTke.TkeClient.PlatformV1().Clusters().Create(context.Background(), cluster, metav1.CreateOptions{})
	if err != nil {
		klog.Error(err)
		return
	}

	klog.Info("Cluster name: ", cluster.Name)
	return testTke.WaitClusterToBeRunning(cluster.Name)
}

func (testTke *TestTKE) DeleteCluster(clusterName string) (err error) {
	klog.Info("Delete cluster: ", clusterName)
	err = testTke.TkeClient.PlatformV1().Clusters().Delete(context.Background(), clusterName, metav1.DeleteOptions{})
	if k8serror.IsNotFound(err) {
		klog.Info("Cluster was not found")
		return nil
	}
	if err != nil {
		return fmt.Errorf("delete cluster failed. %v", err)
	}
	klog.Info("Wait cluster to be deleted")
	return wait.Poll(5*time.Second, 10*time.Minute, func() (bool, error) {
		clusters, err := testTke.TkeClient.PlatformV1().Clusters().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			klog.Error(err)
			return false, err
		}
		for _, cls := range clusters.Items {
			if cls.Name == clusterName {
				klog.Info(cls.Status.Phase)
				return false, nil
			}
		}
		klog.Info("Cluster was deleted")
		return true, nil
	})
}

func (testTke *TestTKE) AddNode(clusterName string, workerNode cloudprovider.Instance) (machine *platformv1.Machine, err error) {
	klog.Info("Add node. InstanceId: ", workerNode.InstanceID, ", InternalIP: ", workerNode.InternalIP)
	machine = &platformv1.Machine{
		Spec: platformv1.MachineSpec{
			ClusterName: clusterName,
			Type:        "Baremetal",
			IP:          workerNode.InternalIP,
			Port:        workerNode.Port,
			Username:    workerNode.Username,
			Password:    []byte(workerNode.Password),
		},
	}
	machine, err = testTke.TkeClient.PlatformV1().Machines().Create(context.Background(), machine, metav1.CreateOptions{})
	if err != nil {
		klog.Error(err)
		return
	}

	klog.Info("Wait node status to be running")
	err = wait.Poll(5*time.Second, 5*time.Minute, func() (bool, error) {
		machine, err = testTke.TkeClient.PlatformV1().Machines().Get(context.Background(), machine.Name, metav1.GetOptions{})
		if err != nil {
			klog.Error(err)
			return false, nil
		}
		if len(machine.Status.Conditions) > 0 {
			lastCondition := machine.Status.Conditions[len(machine.Status.Conditions)-1]
			klog.Info("Phase: ", machine.Status.Phase, ", Type: ", lastCondition.Type, ", message: ", lastCondition.Message)
		}
		return machine.Status.Phase == platformv1.MachineRunning, nil
	})
	return
}

func (testTke *TestTKE) DeleteNode(machineName string) (err error) {
	klog.Info("Delete node: ", machineName)
	err = testTke.TkeClient.PlatformV1().Machines().Delete(context.Background(), machineName, metav1.DeleteOptions{})
	if err != nil {
		return
	}

	klog.Info("Wait node to be deleted")
	return wait.Poll(5*time.Second, 10*time.Minute, func() (bool, error) {
		_, err = testTke.TkeClient.PlatformV1().Machines().Get(context.Background(), machineName, metav1.GetOptions{})
		if k8serror.IsNotFound(err) {
			klog.Info("Node was deleted")
			return true, nil
		}
		return false, nil
	})
}

func (testTke *TestTKE) UnscheduleNode(cls *platformv1.Cluster, nodeName string) (node *corev1.Node, err error) {
	return testTke.updateNode(cls, nodeName, true)
}

func (testTke *TestTKE) CancleUnschedulableNode(cls *platformv1.Cluster, nodeName string) (node *corev1.Node, err error) {
	return testTke.updateNode(cls, nodeName, false)
}

func (testTke *TestTKE) updateNode(cls *platformv1.Cluster, nodeName string, unschedulable bool) (node *corev1.Node, err error) {
	err = wait.Poll(time.Second, 5*time.Second, func() (bool, error) {
		node, err = testTke.K8sClient(cls).CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		node.Spec.Unschedulable = unschedulable
		node, err = testTke.K8sClient(cls).CoreV1().Nodes().Update(context.Background(), node, metav1.UpdateOptions{})
		return true, err
	})
	if err != nil {
		return
	}

	return testTke.K8sClient(cls).CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
}

func (testTke *TestTKE) CreateInstances(count int64) ([]cloudprovider.Instance, error) {
	return testTke.provider.CreateInstances(count)
}
