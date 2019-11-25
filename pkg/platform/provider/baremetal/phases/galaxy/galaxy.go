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

package galaxy

import (
	"io"
	"net"
	"os/exec"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/galaxy/images"
	"tkestack.io/tke/pkg/util/log"
)

const (
	daemonsetFlannelName = "flannel"
	cmFlannel            = "kube-flannel-cfg"
	svcAccounFlannelName = "flannel"
	crFlannelName        = "flannel"
	crbFlannelName       = "flannel"
	daemonsetGalaxyName  = "galaxy-daemonset"
	cmGalaxy             = "galaxy-etc"
	svcAccountName       = "galaxy"
	crbName              = "galaxy"
)

// Option for coredns
type Option struct {
	Version   string
	NodeCIDR  string
	NetDevice string
}

// Install to install the galaxy workload
func Install(clientset kubernetes.Interface, option *Option) error {
	// old flannel interface should be deleted
	if err := cleanFlannelInterfaces(); err != nil {
		return err
	}
	// in private cloud, flannel must be installed
	if _, err := clientset.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(serviceAccountFlannel()); err != nil {
		if !errors.IsAlreadyExists(err) {
			// flannel service account will create automatically
			return err
		}
	}
	if _, err := clientset.RbacV1().ClusterRoles().Create(crFlannel()); err != nil {
		return err
	}
	if _, err := clientset.RbacV1().ClusterRoleBindings().Create(crbFlannel()); err != nil {
		return err
	}
	if _, err := clientset.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(cmFlannel, metav1.GetOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		cm, err := configMapFlannel(option.NodeCIDR)
		if err != nil {
			return err
		}
		if _, err := clientset.CoreV1().ConfigMaps(metav1.NamespaceSystem).Create(cm); err != nil {
			return err
		}
	}
	// Daemonset Flannel
	flannelObj, err := daemonsetFlannel(option.Version)
	if err != nil {
		return err
	}
	if _, err := clientset.AppsV1().DaemonSets(metav1.NamespaceSystem).Create(flannelObj); err != nil {
		log.Errorf("create daemonset with err: %v", err)
		return err
	}
	// flannel installation finished, begin to install galaxy-daemon
	if _, err := clientset.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(serviceAccountGalaxy()); err != nil {
		return err
	}
	// ClusterRoleBinding Galaxy
	if _, err := clientset.RbacV1().ClusterRoleBindings().Create(crbGalaxy()); err != nil {
		return err
	}
	// init galaxy configMap
	if _, err := clientset.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(cmGalaxy, metav1.GetOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		cms, err := configMapGalaxy(option.NetDevice)
		if err != nil {
			return err
		}
		for _, cm := range cms {
			if _, err := clientset.CoreV1().ConfigMaps(metav1.NamespaceSystem).Create(cm); err != nil {
				return err
			}
		}
	}
	// Daemonset Galaxy
	galaxyObj, err := daemonsetGalaxy(option.Version)
	if err != nil {
		return err
	}
	if _, err := clientset.AppsV1().DaemonSets(metav1.NamespaceSystem).Create(galaxyObj); err != nil {
		log.Errorf("create daemonset with err: %v", err)
		return err
	}

	return nil
}

func cleanFlannelInterfaces() error {
	var err error
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, iface := range ifaces {
		if strings.Contains(iface.Name, "flannel") {
			cmd := exec.Command("ip", "link", "delete", iface.Name)
			if err := cmd.Run(); err != nil {
				log.Errorf("fail to delete link %s : %v", iface.Name, err)
			}
		}
	}
	return err
}

func configMapFlannel(clusterCIDR string) (*corev1.ConfigMap, error) {
	reader := strings.NewReader(strings.Replace(FlannelCM, "{{ .Network }}", clusterCIDR, 1))
	decoder := yaml.NewYAMLOrJSONDecoder(reader, 4096)
	payload := &corev1.ConfigMap{}
	err := decoder.Decode(payload)
	if err != nil {
		return nil, err
	}
	payload.Name = cmFlannel
	return payload, nil
}

func configMapGalaxy(netDevice string) ([]*corev1.ConfigMap, error) {
	reader := strings.NewReader(strings.Replace(GalaxyCM, "{{ .DeviceName }}", netDevice, -1))
	var payloads []*corev1.ConfigMap
	decoder := yaml.NewYAMLOrJSONDecoder(reader, 4096)
	for {
		payload := &corev1.ConfigMap{}
		err := decoder.Decode(&payload)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		payloads = append(payloads, payload)
	}

	return payloads, nil
}

func daemonsetFlannel(version string) (*appsv1.DaemonSet, error) {
	imageName := images.Get(version).Flannel.FullName()
	reader := strings.NewReader(strings.Replace(FlannelDaemonset, "{{ .Image }}", imageName, -1))
	payload := &appsv1.DaemonSet{}
	err := yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(payload)
	if err != nil {
		return nil, err
	}
	payload.Name = daemonsetFlannelName
	return payload, nil
}

func daemonsetGalaxy(version string) (*appsv1.DaemonSet, error) {
	reader := strings.NewReader(GalaxyDaemonsetTemplate)
	payload := &appsv1.DaemonSet{}
	err := yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(payload)
	if err != nil {
		return nil, err
	}
	payload.Name = daemonsetGalaxyName
	payload.Spec.Template.Spec.Containers[0].Image = images.Get(version).GalaxyDaemon.FullName()
	return payload, nil
}

func serviceAccountFlannel() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcAccounFlannelName,
			Namespace: metav1.NamespaceSystem,
		},
	}
}

func serviceAccountGalaxy() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcAccountName,
			Namespace: metav1.NamespaceSystem,
		},
	}
}

func crFlannel() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: crFlannelName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"nodes/status"},
				Verbs:     []string{"patch"},
			},
		},
	}
}

func crbFlannel() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: crbFlannelName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      svcAccounFlannelName,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func crbGalaxy() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: crbName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      svcAccountName,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}
