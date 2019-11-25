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

package apiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/api/extensions/v1beta1"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	apps "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	rbac "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	controllerutils "tkestack.io/tke/pkg/controller"
)

const (
	// APICallRetryInterval defines how long should wait before retrying a failed API operation
	APICallRetryInterval = 500 * time.Millisecond
	// PatchNodeTimeout specifies how long should wait for applying the label and taint on the master before timing out
	PatchNodeTimeout = 2 * time.Minute
	// UpdateNodeTimeout specifies how long should wait for updating node with the initial remote configuration of kubelet before timing out
	UpdateNodeTimeout = 2 * time.Minute
	// LabelHostname specifies the lable in node.
	LabelHostname = "kubernetes.io/hostname"
)

// CreateOrUpdateConfigMap creates a ConfigMap if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateConfigMap(client clientset.Interface, cm *v1.ConfigMap) error {
	if _, err := client.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).Create(cm); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create configmap")
		}

		if _, err := client.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).Update(cm); err != nil {
			return errors.Wrap(err, "unable to update configmap")
		}
	}
	return nil
}

// CreateOrRetainConfigMap creates a ConfigMap if the target resource doesn't exist. If the resource exists already, this function will retain the resource instead.
func CreateOrRetainConfigMap(client clientset.Interface, cm *v1.ConfigMap, configMapName string) error {
	if _, err := client.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).Get(configMapName, metav1.GetOptions{}); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		if _, err := client.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).Create(cm); err != nil {
			if !apierrors.IsAlreadyExists(err) {
				return errors.Wrap(err, "unable to create configmap")
			}
		}
	}
	return nil
}

// CreateOrUpdateSecret creates a Secret if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateSecret(client clientset.Interface, secret *v1.Secret) error {
	if _, err := client.CoreV1().Secrets(secret.ObjectMeta.Namespace).Create(secret); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create secret")
		}

		if _, err := client.CoreV1().Secrets(secret.ObjectMeta.Namespace).Update(secret); err != nil {
			return errors.Wrap(err, "unable to update secret")
		}
	}
	return nil
}

// CreateOrUpdateServiceAccount creates a ServiceAccount if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateServiceAccount(client clientset.Interface, sa *v1.ServiceAccount) error {
	if _, err := client.CoreV1().ServiceAccounts(sa.ObjectMeta.Namespace).Create(sa); err != nil {
		// Note: We don't run .Update here afterwards as that's probably not required
		// Only thing that could be updated is annotations/labels in .metadata, but we don't use that currently
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create serviceaccount")
		}
	}
	return nil
}

// CreateOrUpdateDeployment creates a Deployment if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateDeployment(client clientset.Interface, deploy *apps.Deployment) error {
	if _, err := client.AppsV1().Deployments(deploy.ObjectMeta.Namespace).Create(deploy); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create deployment")
		}

		if _, err := client.AppsV1().Deployments(deploy.ObjectMeta.Namespace).Update(deploy); err != nil {
			return errors.Wrap(err, "unable to update deployment")
		}
	}
	return nil
}

// CreateOrUpdateDaemonSet creates a DaemonSet if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateDaemonSet(client clientset.Interface, ds *apps.DaemonSet) error {
	if _, err := client.AppsV1().DaemonSets(ds.ObjectMeta.Namespace).Create(ds); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create daemonset")
		}

		if _, err := client.AppsV1().DaemonSets(ds.ObjectMeta.Namespace).Update(ds); err != nil {
			return errors.Wrap(err, "unable to update daemonset")
		}
	}
	return nil
}

// DeleteDaemonSetForeground deletes the specified DaemonSet in foreground mode; i.e. it blocks until/makes sure all the managed Pods are deleted
func DeleteDaemonSetForeground(client clientset.Interface, namespace, name string) error {
	foregroundDelete := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		PropagationPolicy: &foregroundDelete,
	}
	return client.AppsV1().DaemonSets(namespace).Delete(name, deleteOptions)
}

// DeleteDeploymentForeground deletes the specified Deployment in foreground mode; i.e. it blocks until/makes sure all the managed Pods are deleted
func DeleteDeploymentForeground(client clientset.Interface, namespace, name string) error {
	foregroundDelete := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		PropagationPolicy: &foregroundDelete,
	}
	return client.AppsV1().Deployments(namespace).Delete(name, deleteOptions)
}

// CreateOrUpdateRole creates a Role if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateRole(client clientset.Interface, role *rbac.Role) error {
	if _, err := client.RbacV1().Roles(role.ObjectMeta.Namespace).Create(role); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create RBAC role")
		}

		if _, err := client.RbacV1().Roles(role.ObjectMeta.Namespace).Update(role); err != nil {
			return errors.Wrap(err, "unable to update RBAC role")
		}
	}
	return nil
}

// CreateOrUpdateRoleBinding creates a RoleBinding if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateRoleBinding(client clientset.Interface, roleBinding *rbac.RoleBinding) error {
	if _, err := client.RbacV1().RoleBindings(roleBinding.ObjectMeta.Namespace).Create(roleBinding); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create RBAC rolebinding")
		}

		if _, err := client.RbacV1().RoleBindings(roleBinding.ObjectMeta.Namespace).Update(roleBinding); err != nil {
			return errors.Wrap(err, "unable to update RBAC rolebinding")
		}
	}
	return nil
}

// CreateOrUpdateClusterRole creates a ClusterRole if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateClusterRole(client clientset.Interface, clusterRole *rbac.ClusterRole) error {
	if _, err := client.RbacV1().ClusterRoles().Create(clusterRole); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create RBAC clusterrole")
		}

		if _, err := client.RbacV1().ClusterRoles().Update(clusterRole); err != nil {
			return errors.Wrap(err, "unable to update RBAC clusterrole")
		}
	}
	return nil
}

// CreateOrUpdateClusterRoleBinding creates a ClusterRoleBinding if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateClusterRoleBinding(client clientset.Interface, clusterRoleBinding *rbac.ClusterRoleBinding) error {
	if _, err := client.RbacV1().ClusterRoleBindings().Create(clusterRoleBinding); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create RBAC clusterrolebinding")
		}

		if _, err := client.RbacV1().ClusterRoleBindings().Update(clusterRoleBinding); err != nil {
			return errors.Wrap(err, "unable to update RBAC clusterrolebinding")
		}
	}
	return nil
}

// PatchNodeOnce executes patchFn on the node object found by the node name.
// This is a condition function meant to be used with wait.Poll. false, nil
// implies it is safe to try again, an error indicates no more tries should be
// made and true indicates success.
func PatchNodeOnce(client clientset.Interface, nodeName string, patchFn func(*v1.Node)) func() (bool, error) {
	return func() (bool, error) {
		// First get the node object
		n, err := client.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		// The node may appear to have no labels at first,
		// so we wait for it to get hostname label.
		if _, found := n.ObjectMeta.Labels[LabelHostname]; !found {
			return false, nil
		}

		oldData, err := json.Marshal(n)
		if err != nil {
			return false, errors.Wrapf(err, "failed to marshal unmodified node %q into JSON", n.Name)
		}

		// Execute the mutating function
		patchFn(n)

		newData, err := json.Marshal(n)
		if err != nil {
			return false, errors.Wrapf(err, "failed to marshal modified node %q into JSON", n.Name)
		}

		patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1.Node{})
		if err != nil {
			return false, errors.Wrap(err, "failed to create two way merge patch")
		}

		if _, err := client.CoreV1().Nodes().Patch(n.Name, types.StrategicMergePatchType, patchBytes); err != nil {
			if apierrors.IsConflict(err) {
				fmt.Println("[patchnode] Temporarily unable to update node metadata due to conflict (will retry)")
				return false, nil
			}
			return false, errors.Wrapf(err, "error patching node %q through apiserver", n.Name)
		}

		return true, nil
	}
}

// PatchNode tries to patch a node using patchFn for the actual mutating logic.
// Retries are provided by the wait package.
func PatchNode(client clientset.Interface, nodeName string, patchFn func(*v1.Node)) error {
	// wait.Poll will rerun the condition function every interval function if
	// the function returns false. If the condition function returns an error
	// then the retries end and the error is returned.
	return wait.Poll(APICallRetryInterval, PatchNodeTimeout, PatchNodeOnce(client, nodeName, patchFn))
}

// CreateOrUpdateService creates a service if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateService(client clientset.Interface, svc *v1.Service) error {
	_, err := client.CoreV1().Services(svc.ObjectMeta.Namespace).Get(svc.Name, metav1.GetOptions{})
	if err == nil {
		err := client.CoreV1().Services(svc.ObjectMeta.Namespace).Delete(svc.Name, &metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	if _, err := client.CoreV1().Services(svc.ObjectMeta.Namespace).Create(svc); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create service")
		}

		if _, err := client.CoreV1().Services(svc.ObjectMeta.Namespace).Update(svc); err != nil {
			return errors.Wrap(err, "unable to update service")
		}
	}
	return nil
}

// CreateOrUpdateStatefulSet creates a statefulSet if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateStatefulSet(client clientset.Interface, sts *apps.StatefulSet) error {
	_, err := client.AppsV1().StatefulSets(sts.ObjectMeta.Namespace).Get(sts.Name, metav1.GetOptions{})
	if err == nil {
		err := client.AppsV1().StatefulSets(sts.ObjectMeta.Namespace).Delete(sts.Name, &metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	if _, err := client.AppsV1().StatefulSets(sts.ObjectMeta.Namespace).Create(sts); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create statefulSet")
		}

		if _, err := client.AppsV1().StatefulSets(sts.ObjectMeta.Namespace).Update(sts); err != nil {
			return errors.Wrap(err, "unable to update statefulSet")
		}
	}
	return nil
}

// CreateOrUpdateNamespace creates a namespace if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateNamespace(client clientset.Interface, ns *v1.Namespace) error {
	if _, err := client.CoreV1().Namespaces().Create(ns); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create namespace")
		}

		if _, err := client.CoreV1().Namespaces().Update(ns); err != nil {
			return errors.Wrap(err, "unable to update namespace")
		}
	}
	return nil
}

// CreateOrUpdateEndpoints creates a Endpoints if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateEndpoints(client clientset.Interface, ep *v1.Endpoints) error {
	if _, err := client.CoreV1().Endpoints(ep.ObjectMeta.Namespace).Create(ep); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create endpoints")
		}

		if _, err := client.CoreV1().Endpoints(ep.ObjectMeta.Namespace).Update(ep); err != nil {
			return errors.Wrap(err, "unable to update endpoints")
		}
	}
	return nil
}

// CreateOrUpdateIngress creates a Ingress if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateIngress(client clientset.Interface, ing *v1beta1.Ingress) error {
	if _, err := client.ExtensionsV1beta1().Ingresses(ing.ObjectMeta.Namespace).Create(ing); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create ing")
		}

		if _, err := client.ExtensionsV1beta1().Ingresses(ing.ObjectMeta.Namespace).Update(ing); err != nil {
			return errors.Wrap(err, "unable to update ing")
		}
	}
	return nil
}

// CreateOrUpdateConfigMapFromFile like kubectl create configmap --from-file
func CreateOrUpdateConfigMapFromFile(client clientset.Interface, cm *v1.ConfigMap, pattern string) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return errors.New("no matches found")
	}

	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}
	for _, filename := range matches {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		cm.Data[filepath.Base(filename)] = string(data)
	}

	return CreateOrUpdateConfigMap(client, cm)
}

// DeleteReplicaSetApp delete the replicaset and pod additionally for deployment app with extension group
func DeleteReplicaSetApp(client clientset.Interface, options metav1.ListOptions) error {
	rsList, err := client.ExtensionsV1beta1().ReplicaSets(metav1.NamespaceSystem).List(options)
	if err != nil {
		return err
	}

	var errs []error
	for i := range rsList.Items {
		rs := &rsList.Items[i]
		// update replicas to zero
		rs.Spec.Replicas = controllerutils.Int32Ptr(0)
		_, err = client.ExtensionsV1beta1().ReplicaSets(metav1.NamespaceSystem).Update(rs)
		if err != nil {
			errs = append(errs, err)
		} else {
			// delete replicaset
			err = client.ExtensionsV1beta1().ReplicaSets(metav1.NamespaceSystem).Delete(rs.Name, &metav1.DeleteOptions{})
			if err != nil && !apierrors.IsNotFound(err) {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		errMsg := ""
		for _, e := range errs {
			errMsg += e.Error() + ";"
		}
		return fmt.Errorf("delete replicaSet fail:%s", errMsg)
	}

	return nil
}

// CreateOrUpdateDeploymentExtensionsV1beta1 creates a ExtensionsV1beta1 Deployment if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateDeploymentExtensionsV1beta1(client clientset.Interface, deploy *extensionsv1beta1.Deployment) error {
	if _, err := client.ExtensionsV1beta1().Deployments(deploy.ObjectMeta.Namespace).Create(deploy); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create deployment")
		}

		if _, err := client.ExtensionsV1beta1().Deployments(deploy.ObjectMeta.Namespace).Update(deploy); err != nil {
			return errors.Wrap(err, "unable to update deployment")
		}
	}
	return nil
}

func DeleteDeployment(client clientset.Interface, namespace string, deployName string, isExtensions bool, labelSelector metav1.LabelSelector) error {
	if isExtensions {
		return deleteDeploymentExtensionsV1beta1(client, deployName, labelSelector)
	} else {
		return deleteDeploymentAppsV1(client, namespace, deployName)
	}
}

// DeleteExtensionsV1beta1Deployment delete a deployment
func deleteDeploymentExtensionsV1beta1(client clientset.Interface, deployName string, labelSelector metav1.LabelSelector) error {
	retErr := client.ExtensionsV1beta1().Deployments(metav1.NamespaceSystem).Delete(deployName, &metav1.DeleteOptions{})

	// Delete replicaset for extensions groups
	if retErr != nil {
		selector, err := metav1.LabelSelectorAsSelector(&labelSelector)
		if err != nil {
			retErr = err
		} else {
			options := metav1.ListOptions{
				LabelSelector: selector.String(),
			}
			err = DeleteReplicaSetApp(client, options)
			if err != nil {
				retErr = err
			}
		}
	}
	return retErr
}

// DeleteDeployment delete a deployment
func deleteDeploymentAppsV1(client clientset.Interface, namespace string, deployName string) error {
	return client.AppsV1().Deployments(namespace).Delete(deployName, &metav1.DeleteOptions{})
}

// DeleteClusterRoleBinding delete a clusterrolebinding
func DeleteClusterRoleBinding(client clientset.Interface, crbName string) error {
	return client.RbacV1().ClusterRoleBindings().Delete(crbName, &metav1.DeleteOptions{})
}

// DeleteServiceAccounts delete a serviceAccount
func DeleteServiceAccounts(client clientset.Interface, namespace string, svcAccountName string) error {
	return client.CoreV1().ServiceAccounts(namespace).Delete(svcAccountName, &metav1.DeleteOptions{})
}

// DeleteService delete a service
func DeleteService(client clientset.Interface, namespace string, svcName string) error {
	return client.CoreV1().Services(namespace).Delete(svcName, &metav1.DeleteOptions{})
}

// GetService get a service.
func GetService(client clientset.Interface, namespace string, name string) (*v1.Service, error) {
	return client.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
}

// GetServiceAccount get a service.
func GetServiceAccount(client clientset.Interface, namespace string, name string) (*v1.ServiceAccount, error) {
	return client.CoreV1().ServiceAccounts(namespace).Get(name, metav1.GetOptions{})
}

// GetClusterRoleBinding get a cluster role binding.
func GetClusterRoleBinding(client clientset.Interface, name string) (*rbac.ClusterRoleBinding, error) {
	return client.RbacV1().ClusterRoleBindings().Get(name, metav1.GetOptions{})
}
