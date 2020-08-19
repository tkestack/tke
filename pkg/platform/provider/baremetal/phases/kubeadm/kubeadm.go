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

package kubeadm

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	platformv1client "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	kubeadmv1beta2 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeadm/v1beta2"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubelet"
	"tkestack.io/tke/pkg/platform/provider/baremetal/res"
	platformapiclient "tkestack.io/tke/pkg/platform/util/apiclient"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/pkg/util/template"
)

const (
	kubeadmKubeletConf = "/usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf"

	initCmd = `kubeadm init phase {{.Phase}} --config={{.Config}}`
	joinCmd = `kubeadm join phase {{.Phase}} --config={{.Config}}`
)

var (
	ignoreErrors = []string{
		"ImagePull",
		"Port-10250",
		"FileContent--proc-sys-net-bridge-bridge-nf-call-iptables",
		"DirAvailable--etc-kubernetes-manifests",
	}
)

func Install(s ssh.Interface, version string) error {
	dstFile, err := res.Kubeadm.CopyToNode(s, version)
	if err != nil {
		return err
	}

	cmd := "tar xvaf %s -C %s "
	_, stderr, exit, err := s.Execf(cmd, dstFile, constants.DstBinDir)
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	data, err := template.ParseFile(path.Join(constants.ConfDir, "kubeadm/10-kubeadm.conf"), nil)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(data), kubeadmKubeletConf)
	if err != nil {
		return errors.Wrapf(err, "write %s error", kubeadmKubeletConf)
	}

	return nil
}

func Init(s ssh.Interface, kubeadmConfig *InitConfig, phase string) error {
	configData, err := kubeadmConfig.Marshal()
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(configData), constants.KubeadmConfigFileName)
	if err != nil {
		return err
	}

	cmd, err := template.ParseString(initCmd, map[string]interface{}{
		"Phase":  phase,
		"Config": constants.KubeadmConfigFileName,
	})
	if err != nil {
		return errors.Wrap(err, "parse initCmd error")
	}
	out, err := s.CombinedOutput(string(cmd))
	if err != nil {
		return fmt.Errorf("kubeadm.Init error: %w", err)
	}
	log.Debug(string(out))

	return nil
}

func Join(s ssh.Interface, config *kubeadmv1beta2.JoinConfiguration, phase string) error {
	configData, err := MarshalToYAML(config)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(configData), constants.KubeadmConfigFileName)
	if err != nil {
		return err
	}
	if phase == "preflight" {
		phase = fmt.Sprintf("preflight --ignore-preflight-errors=%s", strings.Join(ignoreErrors, ","))
	}

	cmd, err := template.ParseString(joinCmd, map[string]interface{}{
		"Phase":  phase,
		"Config": constants.KubeadmConfigFileName,
	})
	if err != nil {
		return errors.Wrap(err, "parse joinCmd error")
	}
	out, err := s.CombinedOutput(string(cmd))
	if err != nil {
		return fmt.Errorf("kubeadm.Join error: %w", err)
	}
	log.Debug(string(out))

	return nil
}

func RenewCerts(s ssh.Interface) error {
	err := fixKubeadmBug1753(s)
	if err != nil {
		return fmt.Errorf("fixKubeadmBug1753(https://github.com/kubernetes/kubeadm/issues/1753) error: %w", err)
	}

	cmd := fmt.Sprintf("kubeadm alpha certs renew all --config=%s", constants.KubeadmConfigFileName)
	_, err = s.CombinedOutput(cmd)
	if err != nil {
		return err
	}

	err = RestartControlPlane(s)
	if err != nil {
		return err
	}

	return nil
}

// https://github.com/kubernetes/kubeadm/issues/1753
func fixKubeadmBug1753(s ssh.Interface) error {
	needUpdate := false

	data, err := s.ReadFile(constants.KubeletKubeConfigFileName)
	if err != nil {
		return err
	}
	kubeletKubeconfig, err := clientcmd.Load(data)
	if err != nil {
		return err
	}
	for _, info := range kubeletKubeconfig.AuthInfos {
		if info.ClientKeyData == nil && info.ClientCertificateData == nil {
			continue
		}

		info.ClientKeyData = []byte{}
		info.ClientCertificateData = []byte{}
		info.ClientKey = constants.KubeletClientCurrent
		info.ClientCertificate = constants.KubeletClientCurrent

		needUpdate = true
	}

	if needUpdate {
		data, err := runtime.Encode(clientcmdlatest.Codec, kubeletKubeconfig)
		if err != nil {
			return err
		}
		err = s.WriteFile(bytes.NewReader(data), constants.KubeletKubeConfigFileName)
		if err != nil {
			return err
		}
	}

	return nil
}

// fixKubeadmBug88811 fix after upgrade, coredns deployment volumes still point to backup!
// https://github.com/kubernetes/kubernetes/pull/88811
func fixKubeadmBug88811(client kubernetes.Interface) error {
	patch := []byte(`{"spec":{"template":{"spec":{"volumes":[{"name": "config-volume", "configMap":{"name": "coredns", "items":[{"key": "Corefile", "path": "Corefile"}]}}]}}}}`)
	_, err := client.AppsV1().Deployments(metav1.NamespaceSystem).Patch(context.TODO(), "coredns", types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

func RestartControlPlane(s ssh.Interface) error {
	targets := []string{"kube-apiserver", "kube-controller-manager", "kube-scheduler"}
	for _, one := range targets {
		err := RestartContainerByFilter(s, DockerFilterForControlPlane(one))
		if err != nil {
			return err
		}
	}

	return nil
}

func DockerFilterForControlPlane(name string) string {
	return fmt.Sprintf("label=io.kubernetes.container.name=%s", name)
}

func RestartContainerByFilter(s ssh.Interface, filter string) error {
	cmd := fmt.Sprintf("docker rm -f $(docker ps -q -f '%s')", filter)
	_, err := s.CombinedOutput(cmd)
	if err != nil {
		return err
	}

	err = wait.PollImmediate(5*time.Second, 5*time.Minute, func() (bool, error) {
		cmd = fmt.Sprintf("docker ps -q -f '%s'", filter)
		output, err := s.CombinedOutput(cmd)
		if err != nil {
			return false, nil
		}
		if len(output) == 0 {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("restart container(%s) error: %w", filter, err)
	}

	return nil
}

type NodeRole string

const (
	NodeRoleMaster = NodeRole("Master")
	NodeRoleWorker = NodeRole("Worker")
)

type UpgradeOption struct {
	MachineName   string
	BootstrapNode bool
	NodeName      string
	NodeRole      NodeRole
	Version       string
	MaxUnready    *intstr.IntOrString
}

// UpgradeNode upgrades node by kubeadm.
// Refer: https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-upgrade/
func UpgradeNode(s ssh.Interface, client kubernetes.Interface, platformClient platformv1client.PlatformV1Interface, option UpgradeOption) (upgraded bool, err error) {
	if option.NodeRole == NodeRoleWorker {
		ok, err := checkMasterNodesVersion(client, option.Version)
		if err != nil {
			return upgraded, err
		}
		if !ok {
			return upgraded, fmt.Errorf("must wait for all master nodes to be upgraded, then upgrading worker nodes")
		}
	}
	needUpgrade, err := needUpgradeNode(client, option.NodeName, option.Version)
	if err != nil {
		return upgraded, err
	}
	if !needUpgrade {
		return upgraded, nil
	}

	// Step 1: drain node
	err = drainNodeCarefully(s, client, option.NodeName, option.MaxUnready)
	if err != nil {
		return upgraded, err
	}
	// ensure uncordon node
	defer uncordonNode(s, option.NodeName)

	// Step 2: install kubeadm
	err = Install(s, option.Version)
	if err != nil {
		return upgraded, err
	}

	// Step 3: do upgrade
	if option.NodeRole == NodeRoleMaster {
		needUpgrade, err := needUpgradeControlPlane(client, option.NodeName, option.Version)
		if err != nil {
			return upgraded, err
		}
		if needUpgrade {
			err = upgradeBootstrapNode(s, client, option.Version)
			if err != nil {
				return upgraded, err
			}
		}
	} else {
		err = upgradeNode(s)
		if err != nil {
			return upgraded, err
		}
	}

	// Step 4: upgrade kubelet and kubectl
	err = kubelet.Install(s, option.Version)
	if err != nil {
		return upgraded, err
	}

	// Step 5: wait for node information to be updated
	err = wait.PollImmediate(10*time.Second, 5*time.Minute, func() (bool, error) {
		node, err := client.CoreV1().Nodes().Get(context.TODO(), option.NodeName, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		if strings.HasSuffix(node.Status.NodeInfo.KubeletVersion, option.Version) {
			return true, err
		}
		return false, nil
	})
	if err != nil {
		return upgraded, err
	}

	return true, nil
}

func upgradeBootstrapNode(s ssh.Interface, client kubernetes.Interface, version string) error {
	cmd := fmt.Sprintf("kubeadm upgrade plan %s --ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", version)
	out, err := s.CombinedOutput(cmd)
	if err != nil {
		return err
	}
	log.Debug(string(out))

	cmd = fmt.Sprintf("kubeadm upgrade apply -f %s --ignore-preflight-errors=CoreDNSUnsupportedPlugins,CoreDNSMigration", version)
	out, err = s.CombinedOutput(cmd)
	if err != nil {
		return err
	}
	log.Debug(string(out))

	if ok := apiclient.CheckVersionOrDie(version, "<1.19"); ok {
		err = fixKubeadmBug88811(client)
		if err != nil {
			return fmt.Errorf("fixKubeadmBug88811(https://github.com/kubernetes/kubernetes/pull/88811) error: %w", err)
		}
	}

	return nil
}

func upgradeNode(s ssh.Interface) error {
	out, err := s.CombinedOutput("kubeadm upgrade node")
	if err != nil {
		return err
	}
	log.Debug(string(out))

	return nil
}

func needUpgradeControlPlane(client kubernetes.Interface, nodeName string, version string) (bool, error) {
	name := fmt.Sprintf("kube-apiserver-%s", nodeName)
	pod, err := client.CoreV1().Pods(metav1.NamespaceSystem).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	for _, container := range pod.Spec.Containers {
		if !apiclient.IsPodReady(pod) {
			continue
		}
		if !strings.HasSuffix(container.Image, version) {
			return true, nil
		}
	}

	return false, nil
}

// needUpgradeNode used to determine whether the node can be upgraded.
func needUpgradeNode(client kubernetes.Interface, nodeName string, version string) (bool, error) {
	node, err := client.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	nodeVersion, err := semver.NewVersion(node.Status.NodeInfo.KubeletVersion)
	if err != nil {
		return false, err
	}
	needVersion, err := semver.NewVersion(version)
	if err != nil {
		return false, err
	}

	if nodeVersion.Major() == needVersion.Major() &&
		nodeVersion.Minor() == needVersion.Minor() &&
		nodeVersion.Patch() == needVersion.Patch() {
		return false, nil
	}

	return true, nil
}

// checkMasterNodesVersion check all master nodes version.
func checkMasterNodesVersion(client kubernetes.Interface, version string) (bool, error) {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: fields.OneTermEqualSelector(constants.LabelNodeRoleMaster, "").String(),
	})
	if err != nil {
		return false, err
	}
	for _, node := range nodes.Items {
		if !strings.HasSuffix(node.Status.NodeInfo.KubeletVersion, version) {
			return false, fmt.Errorf("master node(%s) current version is %s, required version is %s", node.Name, node.Status.NodeInfo.KubeletVersion, version)
		}
	}

	return true, nil
}

// drainNodeCarefully drains node and ensure evicted pods are running in other node.
func drainNodeCarefully(s ssh.Interface, client kubernetes.Interface, nodeName string, maxUnready *intstr.IntOrString) error {
	err := drainNode(s, nodeName)
	if err != nil {
		_ = uncordonNode(s, nodeName) // drain node may cause error but cordon the node!
		return err
	}

	var totalPods, unreadyPods int
	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, namespace := range namespaces.Items {
		pods, err := client.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		totalPods += len(pods.Items)
		for _, pod := range pods.Items {
			if !apiclient.IsPodReady(&pod) {
				unreadyPods++
			}
		}
	}
	maxUnreadyThreshold, err := intstr.GetValueFromIntOrPercent(maxUnready, totalPods, true)
	if err != nil {
		return err
	}
	if unreadyPods > maxUnreadyThreshold {
		return fmt.Errorf("unready pods(%d) >= max unready threshold(%d %v/%d)", unreadyPods, maxUnreadyThreshold, maxUnready, totalPods)
	}

	// coredns must be ready, otherwise kubectl upgrade whill hang in waiting!
	err = wait.PollImmediate(5*time.Second, 5*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(context.TODO(), client, metav1.NamespaceSystem, "coredns")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
	if err != nil {
		return fmt.Errorf("coredns is not ready: %w", err)
	}

	return nil
}

// drainNode drains node
func drainNode(s ssh.Interface, nodeName string) error {
	cmd := fmt.Sprintf("kubectl drain %s --ignore-daemonsets", nodeName)
	out, err := s.CombinedOutput(cmd)
	if err != nil {
		return err
	}
	log.Debug(string(out))

	return nil
}

// uncordonNode undordons node
func uncordonNode(s ssh.Interface, nodeName string) error {
	cmd := fmt.Sprintf("kubectl uncordon %s", nodeName)
	out, err := s.CombinedOutput(cmd)
	if err != nil {
		return err
	}
	log.Debug(string(out))

	return nil
}

// markNextUpgradeWorkerNode marks next wokrer node to be upgraded.
func MarkNextUpgradeWorkerNode(client kubernetes.Interface, platformClient platformv1client.PlatformV1Interface, version string) error {
	machines, err := platformClient.Machines().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	// No machines need to be upgraded.
	if len(machines.Items) == 0 {
		return nil
	}

	// Get next upgraded machine by lowest name.
	var nextMachineName string
	for _, machine := range machines.Items {
		node, err := client.CoreV1().Nodes().Get(context.TODO(), machine.Spec.IP, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// Skip upgraded node.
		if strings.HasSuffix(node.Status.NodeInfo.KubeletVersion, version) {
			continue
		}

		if nextMachineName == "" {
			nextMachineName = machine.Name
		} else if strings.Compare(machine.Name, nextMachineName) < 0 {
			nextMachineName = machine.Name
		}
	}
	if nextMachineName != "" {
		err = platformapiclient.PatchMachine(context.TODO(), platformClient, nextMachineName, func(machine *platformv1.Machine) {
			if machine.Labels == nil {
				machine.Labels = make(map[string]string)
			}
			machine.Labels[constants.LabelNodeNeedUpgrade] = ""
			machine.Status.Phase = platformv1.MachineUpgrading
		})
		if err != nil {
			return err
		}
	}

	return nil
}
