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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	admissionv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	apps "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	rbac "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	aaclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	kubeaggregatorclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	utilsnet "k8s.io/utils/net"
	controllerutils "tkestack.io/tke/pkg/controller"
)

// PlatformLabel represents the type of platform.tkestack.io related label.
type PlatformLabel string

const (
	// APICallRetryInterval defines how long should wait before retrying a failed API operation
	APICallRetryInterval = 500 * time.Millisecond
	// PatchNodeTimeout specifies how long should wait for applying the label and taint on the master before timing out
	PatchNodeTimeout = 2 * time.Minute
	// UpdateNodeTimeout specifies how long should wait for updating node with the initial remote configuration of kubelet before timing out
	UpdateNodeTimeout = 2 * time.Minute
	// LabelHostname specifies the label in node.
	LabelHostname = "kubernetes.io/hostname"
	// LabelTopologyZone represents a logical failure domain. It is common for Kubernetes clusters to span multiple zones for increased availability.
	LabelTopologyZone = "topology.kubernetes.io/zone"
	// LabelMachineIPV4 specifies the label in node.
	LabelMachineIPV4 PlatformLabel = "platform.tkestack.io/machine-ip"
	// LabelMachineIPV6Head specifies the label in node.
	LabelMachineIPV6Head PlatformLabel = "platform.tkestack.io/machine-ipv6-head"
	// LabelMachineIPV6Tail specifies the label in node.
	LabelMachineIPV6Tail PlatformLabel = "platform.tkestack.io/machine-ipv6-tail"
	// LabelASNForCilium specifies the label in node when enable Cilium.
	LabelASNCilium PlatformLabel = "infra.tce.io/as"
	// LabelSwitchIPForCilium specifies the label in node when enable Cilium.
	LabelSwitchIPCilium PlatformLabel = "infra.tce.io/switch-ip"
)

// CreateOrUpdateConfigMap creates a ConfigMap if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateConfigMap(ctx context.Context, client clientset.Interface, cm *corev1.ConfigMap) error {
	if _, err := client.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).Create(ctx, cm, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create configmap")
		}

		if _, err := client.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).Update(ctx, cm, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update configmap")
		}
	}
	return nil
}

// CreateOrRetainConfigMap creates a ConfigMap if the target resource doesn't exist. If the resource exists already, this function will retain the resource instead.
func CreateOrRetainConfigMap(ctx context.Context, client clientset.Interface, cm *corev1.ConfigMap, configMapName string) error {
	if _, err := client.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).Get(ctx, configMapName, metav1.GetOptions{}); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		if _, err := client.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).Create(ctx, cm, metav1.CreateOptions{}); err != nil {
			if !apierrors.IsAlreadyExists(err) {
				return errors.Wrap(err, "unable to create configmap")
			}
		}
	}
	return nil
}

// CreateOrUpdateSecret creates a Secret if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateSecret(ctx context.Context, client clientset.Interface, secret *corev1.Secret) error {
	if _, err := client.CoreV1().Secrets(secret.ObjectMeta.Namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create secret")
		}

		if _, err := client.CoreV1().Secrets(secret.ObjectMeta.Namespace).Update(ctx, secret, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update secret")
		}
	}
	return nil
}

// CreateOrUpdateServiceAccount creates a ServiceAccount if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateServiceAccount(ctx context.Context, client clientset.Interface, sa *corev1.ServiceAccount) error {
	if _, err := client.CoreV1().ServiceAccounts(sa.ObjectMeta.Namespace).Create(ctx, sa, metav1.CreateOptions{}); err != nil {
		// Note: We don't run .Update here afterwards as that's probably not required
		// Only thing that could be updated is annotations/labels in .metadata, but we don't use that currently
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create serviceaccount")
		}
	}
	return nil
}

// CreateOrUpdateService creates a service if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdatePod(ctx context.Context, client clientset.Interface, pod *corev1.Pod) error {
	_, err := client.CoreV1().Pods(pod.ObjectMeta.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
	if err == nil {
		gracePeriodSeconds := int64(0)
		deleteOptions := metav1.DeleteOptions{
			GracePeriodSeconds: &gracePeriodSeconds,
		}

		err := client.CoreV1().Pods(pod.ObjectMeta.Namespace).Delete(ctx, pod.Name, deleteOptions)
		if err != nil {
			return err
		}
	}
	if _, err := client.CoreV1().Pods(pod.ObjectMeta.Namespace).Create(ctx, pod, metav1.CreateOptions{}); err != nil {
		return errors.Wrap(err, "unable to create pod")
	}

	return nil
}

// CreateOrUpdateDeployment creates a Deployment if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateDeployment(ctx context.Context, client clientset.Interface, deploy *apps.Deployment) error {
	if _, err := client.AppsV1().Deployments(deploy.ObjectMeta.Namespace).Create(ctx, deploy, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create deployment")
		}

		if _, err := client.AppsV1().Deployments(deploy.ObjectMeta.Namespace).Update(ctx, deploy, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update deployment")
		}
	}
	return nil
}

// CreateOrUpdateDaemonSet creates a DaemonSet if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateDaemonSet(ctx context.Context, client clientset.Interface, ds *apps.DaemonSet) error {
	if _, err := client.AppsV1().DaemonSets(ds.ObjectMeta.Namespace).Create(ctx, ds, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create daemonset")
		}

		if _, err := client.AppsV1().DaemonSets(ds.ObjectMeta.Namespace).Update(ctx, ds, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update daemonset")
		}
	}
	return nil
}

// DeleteDaemonSetForeground deletes the specified DaemonSet in foreground mode; i.e. it blocks until/makes sure all the managed Pods are deleted
func DeleteDaemonSetForeground(ctx context.Context, client clientset.Interface, namespace, name string) error {
	foregroundDelete := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &foregroundDelete,
	}
	return client.AppsV1().DaemonSets(namespace).Delete(ctx, name, deleteOptions)
}

// DeleteDeploymentForeground deletes the specified Deployment in foreground mode; i.e. it blocks until/makes sure all the managed Pods are deleted
func DeleteDeploymentForeground(ctx context.Context, client clientset.Interface, namespace, name string) error {
	foregroundDelete := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &foregroundDelete,
	}
	return client.AppsV1().Deployments(namespace).Delete(ctx, name, deleteOptions)
}

// CreateOrUpdateRole creates a Role if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateRole(ctx context.Context, client clientset.Interface, role *rbac.Role) error {
	if _, err := client.RbacV1().Roles(role.ObjectMeta.Namespace).Create(ctx, role, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create RBAC role")
		}

		if _, err := client.RbacV1().Roles(role.ObjectMeta.Namespace).Update(ctx, role, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update RBAC role")
		}
	}
	return nil
}

// CreateOrUpdateRoleBinding creates a RoleBinding if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateRoleBinding(ctx context.Context, client clientset.Interface, roleBinding *rbac.RoleBinding) error {
	if _, err := client.RbacV1().RoleBindings(roleBinding.ObjectMeta.Namespace).Create(ctx, roleBinding, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create RBAC rolebinding")
		}

		if _, err := client.RbacV1().RoleBindings(roleBinding.ObjectMeta.Namespace).Update(ctx, roleBinding, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update RBAC rolebinding")
		}
	}
	return nil
}

// CreateOrUpdateClusterRole creates a ClusterRole if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateClusterRole(ctx context.Context, client clientset.Interface, clusterRole *rbac.ClusterRole) error {
	if _, err := client.RbacV1().ClusterRoles().Create(ctx, clusterRole, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create RBAC clusterrole")
		}

		if _, err := client.RbacV1().ClusterRoles().Update(ctx, clusterRole, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update RBAC clusterrole")
		}
	}
	return nil
}

// CreateOrUpdateClusterRoleBinding creates a ClusterRoleBinding if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateClusterRoleBinding(ctx context.Context, client clientset.Interface, clusterRoleBinding *rbac.ClusterRoleBinding) error {
	if _, err := client.RbacV1().ClusterRoleBindings().Create(ctx, clusterRoleBinding, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create RBAC clusterrolebinding")
		}

		if _, err := client.RbacV1().ClusterRoleBindings().Update(ctx, clusterRoleBinding, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update RBAC clusterrolebinding")
		}
	}
	return nil
}

// PatchNodeOnce executes patchFn on the node object found by the node name.
// This is a condition function meant to be used with wait.Poll. false, nil
// implies it is safe to try again, an error indicates no more tries should be
// made and true indicates success.
func PatchNodeOnce(ctx context.Context, client clientset.Interface, nodeName string, patchFn func(*corev1.Node)) func() (bool, error) {
	return func() (bool, error) {
		// First get the node object
		n, err := client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
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

		patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, corev1.Node{})
		if err != nil {
			return false, errors.Wrap(err, "failed to create two way merge patch")
		}

		if _, err := client.CoreV1().Nodes().Patch(ctx, n.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{}); err != nil {
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
func PatchNode(ctx context.Context, client clientset.Interface, nodeName string, patchFn func(*corev1.Node)) error {
	// wait.Poll will rerun the condition function every interval function if
	// the function returns false. If the condition function returns an error
	// then the retries end and the error is returned.
	return wait.Poll(APICallRetryInterval, PatchNodeTimeout, PatchNodeOnce(ctx, client, nodeName, patchFn))
}

// CreateOrUpdateService creates a service if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateService(ctx context.Context, client clientset.Interface, svc *corev1.Service) error {
	_, err := client.CoreV1().Services(svc.ObjectMeta.Namespace).Get(ctx, svc.Name, metav1.GetOptions{})
	if err == nil {
		err := client.CoreV1().Services(svc.ObjectMeta.Namespace).Delete(ctx, svc.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	if _, err := client.CoreV1().Services(svc.ObjectMeta.Namespace).Create(ctx, svc, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create service")
		}

		if _, err := client.CoreV1().Services(svc.ObjectMeta.Namespace).Update(ctx, svc, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update service")
		}
	}
	return nil
}

// CreateOrUpdateStatefulSet creates a statefulSet if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateStatefulSet(ctx context.Context, client clientset.Interface, sts *apps.StatefulSet) error {
	_, err := client.AppsV1().StatefulSets(sts.ObjectMeta.Namespace).Get(ctx, sts.Name, metav1.GetOptions{})
	if err == nil {
		err := client.AppsV1().StatefulSets(sts.ObjectMeta.Namespace).Delete(ctx, sts.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	if _, err := client.AppsV1().StatefulSets(sts.ObjectMeta.Namespace).Create(ctx, sts, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create statefulSet")
		}

		if _, err := client.AppsV1().StatefulSets(sts.ObjectMeta.Namespace).Update(ctx, sts, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update statefulSet")
		}
	}
	return nil
}

// CreateOrUpdateNamespace creates a namespace if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateNamespace(ctx context.Context, client clientset.Interface, ns *corev1.Namespace) error {
	if _, err := client.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create namespace")
		}
		if _, err := client.CoreV1().Namespaces().Update(ctx, ns, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update namespace")
		}
	}
	return nil
}

// CreateOrUpdateEndpoints creates a Endpoints if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateEndpoints(ctx context.Context, client clientset.Interface, ep *corev1.Endpoints) error {
	if _, err := client.CoreV1().Endpoints(ep.ObjectMeta.Namespace).Create(ctx, ep, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create endpoints")
		}

		if _, err := client.CoreV1().Endpoints(ep.ObjectMeta.Namespace).Update(ctx, ep, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update endpoints")
		}
	}
	return nil
}

// CreateOrUpdateIngress creates a Ingress if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateIngress(ctx context.Context, client clientset.Interface, ing *extensionsv1beta1.Ingress) error {
	if _, err := client.ExtensionsV1beta1().Ingresses(ing.ObjectMeta.Namespace).Create(ctx, ing, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create ingress")
		}

		if _, err := client.ExtensionsV1beta1().Ingresses(ing.ObjectMeta.Namespace).Update(ctx, ing, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update ingress")
		}
	}
	return nil
}

// CreateOrUpdateJob creates a Job if the target resource doesn't exist. If the resource exists already, this function will update
func CreateOrUpdateJob(ctx context.Context, client clientset.Interface, job *batchv1.Job) error {
	if _, err := client.BatchV1().Jobs(job.ObjectMeta.Namespace).Create(ctx, job, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create job")
		}

		if _, err := client.BatchV1().Jobs(job.ObjectMeta.Namespace).Update(ctx, job, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update job")
		}
	}
	return nil
}

// CreateOrUpdateCronJob creates a Job if the target resource doesn't exist. If the resource exists already, this function will update
func CreateOrUpdateCronJob(ctx context.Context, client clientset.Interface, cronjob *batchv1beta1.CronJob) error {
	if _, err := client.BatchV1beta1().CronJobs(cronjob.ObjectMeta.Namespace).Create(ctx, cronjob, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create cronjob")
		}

		if _, err := client.BatchV1beta1().CronJobs(cronjob.ObjectMeta.Namespace).Update(ctx, cronjob, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update cronjob")
		}
	}
	return nil
}

// CreateOrUpdateAPIService creates a APIService if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateAPIService(ctx context.Context, client kubeaggregatorclientset.Interface, as *apiregistrationv1.APIService) error {
	if _, err := client.ApiregistrationV1().APIServices().Create(ctx, as, metav1.CreateOptions{}); err != nil {
		// Note: We don't run .Update here afterwards as that's probably not required
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create apiservice")
		}
	}
	return nil
}

// CreateOrUpdateCustomResourceDefinition creates a CustomResourceDefinition if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateCustomResourceDefinition(ctx context.Context, client aaclientset.Interface, crd *apiextensionsv1.CustomResourceDefinition) error {
	if _, err := client.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, crd, metav1.CreateOptions{}); err != nil {
		// Note: We don't run .Update here afterwards as that's probably not required
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create apiservice")
		}
	}
	return nil
}

// CreateOrUpdateConfigMapFromFile like kubectl apply configmap --from-file
func CreateOrUpdateConfigMapFromFile(ctx context.Context, client clientset.Interface, cm *corev1.ConfigMap, pattern string) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return errors.New("no matches found")
	}

	existCM, err := client.CoreV1().ConfigMaps(cm.Namespace).Get(ctx, cm.Name, metav1.GetOptions{})
	if err == nil {
		cm.Data = existCM.Data
	}
	if err != nil && !apierrors.IsNotFound(err) {
		return err
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

	return CreateOrUpdateConfigMap(ctx, client, cm)
}

// DeleteReplicaSetApp delete the replicaset and pod additionally for deployment app with extension group
func DeleteReplicaSetApp(ctx context.Context, client clientset.Interface, options metav1.ListOptions) error {
	rsList, err := client.ExtensionsV1beta1().ReplicaSets(metav1.NamespaceSystem).List(ctx, options)
	if err != nil {
		return err
	}

	var errs []error
	for i := range rsList.Items {
		rs := &rsList.Items[i]
		// update replicas to zero
		rs.Spec.Replicas = controllerutils.Int32Ptr(0)
		_, err = client.ExtensionsV1beta1().ReplicaSets(metav1.NamespaceSystem).Update(ctx, rs, metav1.UpdateOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			// delete replicaset
			err = client.ExtensionsV1beta1().ReplicaSets(metav1.NamespaceSystem).Delete(ctx, rs.Name, metav1.DeleteOptions{})
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
func CreateOrUpdateDeploymentExtensionsV1beta1(ctx context.Context, client clientset.Interface, deploy *extensionsv1beta1.Deployment) error {
	if _, err := client.ExtensionsV1beta1().Deployments(deploy.ObjectMeta.Namespace).Create(ctx, deploy, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create deployment")
		}

		if _, err := client.ExtensionsV1beta1().Deployments(deploy.ObjectMeta.Namespace).Update(ctx, deploy, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update deployment")
		}
	}
	return nil
}

func DeleteDeployment(ctx context.Context, client clientset.Interface, namespace string, deployName string, isExtensions bool, labelSelector metav1.LabelSelector) error {
	if isExtensions {
		return deleteDeploymentExtensionsV1beta1(ctx, client, deployName, labelSelector)
	} else {
		return deleteDeploymentAppsV1(ctx, client, namespace, deployName)
	}
}

// DeleteExtensionsV1beta1Deployment delete a deployment
func deleteDeploymentExtensionsV1beta1(ctx context.Context, client clientset.Interface, deployName string, labelSelector metav1.LabelSelector) error {
	retErr := client.ExtensionsV1beta1().Deployments(metav1.NamespaceSystem).Delete(ctx, deployName, metav1.DeleteOptions{})

	// Delete replicaset for extensions groups
	if retErr != nil {
		selector, err := metav1.LabelSelectorAsSelector(&labelSelector)
		if err != nil {
			retErr = err
		} else {
			options := metav1.ListOptions{
				LabelSelector: selector.String(),
			}
			err = DeleteReplicaSetApp(ctx, client, options)
			if err != nil {
				retErr = err
			}
		}
	}
	return retErr
}

// DeleteDeployment delete a deployment
func deleteDeploymentAppsV1(ctx context.Context, client clientset.Interface, namespace string, deployName string) error {
	return client.AppsV1().Deployments(namespace).Delete(ctx, deployName, metav1.DeleteOptions{})
}

// DeleteClusterRoleBinding delete a clusterrolebinding
func DeleteClusterRoleBinding(ctx context.Context, client clientset.Interface, crbName string) error {
	return client.RbacV1().ClusterRoleBindings().Delete(ctx, crbName, metav1.DeleteOptions{})
}

// DeleteServiceAccounts delete a serviceAccount
func DeleteServiceAccounts(ctx context.Context, client clientset.Interface, namespace string, svcAccountName string) error {
	return client.CoreV1().ServiceAccounts(namespace).Delete(ctx, svcAccountName, metav1.DeleteOptions{})
}

// DeleteService delete a service
func DeleteService(ctx context.Context, client clientset.Interface, namespace string, svcName string) error {
	return client.CoreV1().Services(namespace).Delete(ctx, svcName, metav1.DeleteOptions{})
}

// GetService get a service.
func GetService(ctx context.Context, client clientset.Interface, namespace string, name string) (*corev1.Service, error) {
	return client.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
}

// GetServiceAccount get a service.
func GetServiceAccount(ctx context.Context, client clientset.Interface, namespace string, name string) (*corev1.ServiceAccount, error) {
	return client.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})
}

// GetClusterRoleBinding get a cluster role binding.
func GetClusterRoleBinding(ctx context.Context, client clientset.Interface, name string) (*rbac.ClusterRoleBinding, error) {
	return client.RbacV1().ClusterRoleBindings().Get(ctx, name, metav1.GetOptions{})
}

// CreateOrUpdateValidatingWebhookConfiguration creates a ValidatingWebhookConfigurations if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateValidatingWebhookConfiguration(ctx context.Context, client clientset.Interface, obj *admissionv1beta1.ValidatingWebhookConfiguration) error {
	if _, err := client.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Create(ctx, obj, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create ValidatingWebhookConfiguration")
		}

		if _, err := client.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Update(ctx, obj, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update ValidatingWebhookConfiguration")
		}
	}

	return nil
}

// CreateOrUpdateMutatingWebhookConfiguration creates a MutatingWebhookConfigurations if the target resource doesn't exist. If the resource exists already, this function will update the resource instead.
func CreateOrUpdateMutatingWebhookConfiguration(ctx context.Context, client clientset.Interface, obj *admissionv1beta1.MutatingWebhookConfiguration) error {
	if _, err := client.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Create(ctx, obj, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "unable to create MutatingWebhookConfiguration")
		}

		if _, err := client.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Update(ctx, obj, metav1.UpdateOptions{}); err != nil {
			return errors.Wrap(err, "unable to update MutatingWebhookConfiguration")
		}
	}

	return nil
}

// MarkNode mark node by adding labels and taints
func MarkNode(ctx context.Context, client clientset.Interface, nodeName string, labels map[string]string, taints []corev1.Taint) error {
	return PatchNode(ctx, client, nodeName, func(n *corev1.Node) {
		for k, v := range labels {
			n.Labels[k] = v
		}

		for _, oldTaint := range n.Spec.Taints {
			existed := false
			for _, newTaint := range taints {
				if newTaint.MatchTaint(&oldTaint) {
					existed = true
					break
				}
			}
			if !existed {
				taints = append(taints, oldTaint)
			}
		}
		n.Spec.Taints = taints
	})
}

// RemoveNodeTaints remove taints from existed node taints
func RemoveNodeTaints(ctx context.Context, client clientset.Interface, nodeName string, taints []corev1.Taint) error {
	return PatchNode(ctx, client, nodeName, func(n *corev1.Node) {
		var newTaints []corev1.Taint
		for _, oldTaint := range n.Spec.Taints {
			existed := false
			for _, taint := range taints {
				if taint.MatchTaint(&oldTaint) {
					existed = true
					break
				}
			}
			if !existed {
				newTaints = append(newTaints, oldTaint)
			}
		}
		n.Spec.Taints = newTaints
	})
}

// GetNodeByMachineIP get node by machine ip.
func GetNodeByMachineIP(ctx context.Context, client clientset.Interface, ip string) (*corev1.Node, error) {
	// try to get node by name = machine ip
	node, err := client.CoreV1().Nodes().Get(ctx, ip, metav1.GetOptions{})
	if !apierrors.IsNotFound(err) {
		return node, err
	}
	labelSelector := fields.OneTermEqualSelector(string(LabelMachineIPV4), ip).String()
	if utilsnet.IsIPv6String(ip) {
		labelSelector = GetNodeIPV6Label(ip)
	}
	nodes, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return &corev1.Node{}, err
	}
	if len(nodes.Items) == 0 {
		return &corev1.Node{}, apierrors.NewNotFound(corev1.Resource("Node"), labelSelector)
	}
	return &nodes.Items[0], nil
}

// GetNodeIPV6Label split ip v6 address to head and tail ensure lable value
// less than 63 character, since k8s lable doesn't support ":" so that replace
// to "a", then return the consolidated label string
// Todo: add more check and corner case handle here later
func GetNodeIPV6Label(ip string) string {
	midLength := len(ip) / 2
	splitLength := int(math.Ceil(float64(midLength)))
	lableipv6Head := fmt.Sprintf("%s=%s", LabelMachineIPV6Head, strings.Replace(ip[0:splitLength], ":", "a", -1))
	lableipv6Tail := fmt.Sprintf("%s=%s", LabelMachineIPV6Tail, strings.Replace(ip[splitLength:], ":", "a", -1))
	return lableipv6Head + "," + lableipv6Tail
}
