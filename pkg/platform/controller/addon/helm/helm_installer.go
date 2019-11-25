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

package helm

import (
	normalerrors "errors"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"tkestack.io/tke/pkg/platform/controller/addon/helm/images"
	"tkestack.io/tke/pkg/util/apiclient"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/api/platform"
	controllerutil "tkestack.io/tke/pkg/controller"
)

const (
	deployTillerName   = "tiller-deploy"
	deployHelmAPIName  = "swift"
	svcTillerName      = "tiller-deploy"
	svcHelmAPIName     = "helm-api"
	svcAccountHelmName = "helm"
	crbHelmName        = "helm"
)

var selectorForTiller = metav1.LabelSelector{
	MatchLabels: map[string]string{"app": "helm", "name": "tiller"},
}

var selectorForOfficialTiller = metav1.LabelSelector{
	MatchLabels: map[string]string{"app": "helm", "name": "tiller", "qcloud-app": "tiller", "k8s-app": "tiller"},
}

var selectorForHelm = metav1.LabelSelector{
	MatchLabels: map[string]string{"app": "swift", "qcloud-app": "swift", "k8s-app": "swift"},
}

type Provisioner interface {
	Install() error
	Uninstall() error
	GetStatus() error
}

type provisioner struct {
	kubeClient kubernetes.Interface
	option     Option
}

// Option is option for this component
type Option struct {
	version              string
	isExtensionsAPIGroup bool
}

func NewProvisioner(kubeClient kubernetes.Interface, option *Option) Provisioner {
	return &provisioner{
		kubeClient: kubeClient,
		option:     *option,
	}
}

func (p *provisioner) Install() error {
	// if unOfficial tiller in cluster
	err := p.isOfficialTiller()
	if err != nil {
		return err
	}

	// begin install
	kubeClient := p.kubeClient
	option := p.option

	// ServiceAccount Helm
	if err := apiclient.CreateOrUpdateServiceAccount(kubeClient, serviceAccountHelm()); err != nil {
		return err
	}
	// ClusterRoleBinding Helm
	if err := apiclient.CreateOrUpdateClusterRoleBinding(kubeClient, crbHelm()); err != nil {
		return err
	}
	// Deployment Tiller
	if option.isExtensionsAPIGroup {
		if err := apiclient.CreateOrUpdateDeploymentExtensionsV1beta1(kubeClient, deploymentTillerExtensions(option.version)); err != nil {
			return err
		}
	} else {
		if err := apiclient.CreateOrUpdateDeployment(kubeClient, deploymentTiller(option.version)); err != nil {
			return err
		}
	}
	// Service Tiller
	if err := apiclient.CreateOrUpdateService(kubeClient, serviceTiller()); err != nil {
		return err
	}
	// Deployment Helm-api
	if option.isExtensionsAPIGroup {
		if err := apiclient.CreateOrUpdateDeploymentExtensionsV1beta1(kubeClient, deploymentHelmAPIExtensions(option.version)); err != nil {
			return err
		}
	} else {
		if err := apiclient.CreateOrUpdateDeployment(kubeClient, deploymentHelmAPI(option.version)); err != nil {
			return err
		}
	}
	// Service Helm-api
	if err := apiclient.CreateOrUpdateService(kubeClient, serviceHelmAPI()); err != nil {
		return err
	}
	return nil
}

func (p *provisioner) Uninstall() error {
	kubeClient := p.kubeClient
	option := p.option

	// Service Helm-api
	svcHelmAPIErr := apiclient.DeleteService(kubeClient, metav1.NamespaceSystem, svcHelmAPIName)
	// Deployment Helm-api
	deployHelmAPIErr := apiclient.DeleteDeployment(kubeClient, metav1.NamespaceSystem, deployHelmAPIName, option.isExtensionsAPIGroup, selectorForHelm)
	// Service Tiller
	svcTillerErr := apiclient.DeleteService(kubeClient, metav1.NamespaceSystem, svcTillerName)
	// Deployment Tiller
	deployTillerErr := apiclient.DeleteDeployment(kubeClient, metav1.NamespaceSystem, deployTillerName, option.isExtensionsAPIGroup, selectorForOfficialTiller)
	// ClusterRoleBinding Helm
	crbHelmErr := apiclient.DeleteClusterRoleBinding(kubeClient, crbHelmName)
	// ServiceAccount Helm
	svcAccountHelmErr := apiclient.DeleteServiceAccounts(kubeClient, metav1.NamespaceSystem, svcAccountHelmName)

	if (svcHelmAPIErr != nil && !errors.IsNotFound(svcHelmAPIErr)) ||
		(deployHelmAPIErr != nil && !errors.IsNotFound(deployHelmAPIErr)) ||
		(svcTillerErr != nil && !errors.IsNotFound(svcTillerErr)) ||
		(deployTillerErr != nil && !errors.IsNotFound(deployTillerErr)) ||
		(crbHelmErr != nil && !errors.IsNotFound(crbHelmErr)) ||
		(svcAccountHelmErr != nil && !errors.IsNotFound(svcAccountHelmErr)) {
		return normalerrors.New("delete helm error")
	}
	return nil
}

func (p *provisioner) GetStatus() error {
	if _, err := p.kubeClient.CoreV1().Services(metav1.NamespaceSystem).ProxyGet("http", svcHelmAPIName, "http", `/tiller/v2/version/json`, nil).DoRaw(); err != nil {
		// get more detailed checkErr about resource
		if checkErr := p.checkRsc(p.kubeClient); checkErr != nil {
			return checkErr
		}
		return err
	}
	return nil
}

func serviceAccountHelm() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcAccountHelmName,
			Namespace: metav1.NamespaceSystem,
		},
	}
}

func crbHelm() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: crbHelmName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      svcAccountHelmName,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func deploymentTiller(version string) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployTillerName,
			Labels:    selectorForOfficialTiller.MatchLabels,
			Namespace: metav1.NamespaceSystem,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: controllerutil.Int32Ptr(1),
			Selector: &selectorForOfficialTiller,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: selectorForOfficialTiller.MatchLabels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: svcAccountHelmName,
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{Name: "TILLER_NAMESPACE", Value: metav1.NamespaceSystem},
								{Name: "TILLER_HISTORY_MAX", Value: "0"},
							},
							Name:  "tiller",
							Image: images.Get(version).Tiller.FullName(),
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/liveness",
										Port: intstr.FromInt(44135),
									},
								},
								InitialDelaySeconds: 1,
								TimeoutSeconds:      1,
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 44134, Name: "tiller"},
								{ContainerPort: 44135, Name: "http"},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/readiness",
										Port: intstr.FromInt(44135),
									},
								},
								InitialDelaySeconds: 1,
								TimeoutSeconds:      1,
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(150, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(80*1024*1024, resource.BinarySI),
								},
							},
						},
					},
				},
			},
		},
	}
}

func deploymentTillerExtensions(version string) *extensionsv1beta1.Deployment {
	return &extensionsv1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployTillerName,
			Labels:    selectorForOfficialTiller.MatchLabels,
			Namespace: metav1.NamespaceSystem,
		},
		Spec: extensionsv1beta1.DeploymentSpec{
			Replicas: controllerutil.Int32Ptr(1),
			Selector: &selectorForOfficialTiller,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: selectorForOfficialTiller.MatchLabels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: svcAccountHelmName,
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{Name: "TILLER_NAMESPACE", Value: metav1.NamespaceSystem},
								{Name: "TILLER_HISTORY_MAX", Value: "0"},
							},
							Name:  "tiller",
							Image: images.Get(version).Tiller.FullName(),
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/liveness",
										Port: intstr.FromInt(44135),
									},
								},
								InitialDelaySeconds: 1,
								TimeoutSeconds:      1,
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 44134, Name: "tiller"},
								{ContainerPort: 44135, Name: "http"},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/readiness",
										Port: intstr.FromInt(44135),
									},
								},
								InitialDelaySeconds: 1,
								TimeoutSeconds:      1,
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(150, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(80*1024*1024, resource.BinarySI),
								},
							},
						},
					},
				},
			},
		},
	}
}

func serviceTiller() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcTillerName,
			Namespace: metav1.NamespaceSystem,
			Labels:    selectorForOfficialTiller.MatchLabels,
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorForOfficialTiller.MatchLabels,
			Ports: []corev1.ServicePort{
				{Name: "tiller", Port: 44134, TargetPort: intstr.FromString("tiller")},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

func deploymentHelmAPI(version string) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployHelmAPIName,
			Labels:    selectorForHelm.MatchLabels,
			Namespace: metav1.NamespaceSystem,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: controllerutil.Int32Ptr(1),
			Selector: &selectorForHelm,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      selectorForHelm.MatchLabels,
					Annotations: map[string]string{"scheduler.alpha.kubernetes.io/critical-pod": ""},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: svcAccountHelmName,
					Containers: []corev1.Container{
						{
							Name:  "swift",
							Image: images.Get(version).Swift.FullName(),
							Args:  []string{"run", "--v=3", "--connector=incluster", "--tiller-insecure-skip-verify=true", "--enable-analytics=true"},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 9855},
								{ContainerPort: 50055},
								{ContainerPort: 56790},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/tmp",
									Name:      "chart-volume",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(80*1024*1024, resource.BinarySI),
								},
							},
						},
						{
							Name:  "swift-reverse-proxy",
							Image: images.Get(version).HelmAPI.FullName(),
							Ports: []corev1.ContainerPort{
								{ContainerPort: 8080},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(30, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(20*1024*1024, resource.BinarySI),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "chart-volume",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "CriticalAddonsOnly",
							Operator: "Exists",
						},
					},
				},
			},
		},
	}
}

func deploymentHelmAPIExtensions(version string) *extensionsv1beta1.Deployment {
	return &extensionsv1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployHelmAPIName,
			Labels:    selectorForHelm.MatchLabels,
			Namespace: metav1.NamespaceSystem,
		},
		Spec: extensionsv1beta1.DeploymentSpec{
			Replicas: controllerutil.Int32Ptr(1),
			Selector: &selectorForHelm,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      selectorForHelm.MatchLabels,
					Annotations: map[string]string{"scheduler.alpha.kubernetes.io/critical-pod": ""},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: svcAccountHelmName,
					Containers: []corev1.Container{
						{
							Name:  "swift",
							Image: images.Get(version).Swift.FullName(),
							Args:  []string{"run", "--v=3", "--connector=incluster", "--tiller-insecure-skip-verify=true", "--enable-analytics=true"},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 9855},
								{ContainerPort: 50055},
								{ContainerPort: 56790},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/tmp",
									Name:      "chart-volume",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(80*1024*1024, resource.BinarySI),
								},
							},
						},
						{
							Name:  "swift-reverse-proxy",
							Image: images.Get(version).HelmAPI.FullName(),
							Ports: []corev1.ContainerPort{
								{ContainerPort: 8080},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(30, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(20*1024*1024, resource.BinarySI),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "chart-volume",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "CriticalAddonsOnly",
							Operator: "Exists",
						},
					},
				},
			},
		},
	}
}

func serviceHelmAPI() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcHelmAPIName,
			Namespace: metav1.NamespaceSystem,
			Labels: map[string]string{
				"kubernetes.io/cluster-service": "true",
				"app":                           "swift",
				"qcloud-app":                    "swift",
				"k8s-app":                       "swift",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorForHelm.MatchLabels,
			Ports: []corev1.ServicePort{
				{Name: "http", Port: 80, TargetPort: intstr.FromInt(8080)},
				{Name: "https", Port: 443, TargetPort: intstr.FromInt(50055)},
				{Name: "ops", Port: 56790, TargetPort: intstr.FromInt(56790)},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

func (p *provisioner) isOfficialTiller() error {
	kubeClient := p.kubeClient
	isExtensionsAPIGroup := p.option.isExtensionsAPIGroup

	var count int
	var deployLabels map[string]string
	tillerLabelSelector, err := metav1.LabelSelectorAsSelector(&selectorForTiller)
	if err != nil {
		return err
	}

	if !isExtensionsAPIGroup {
		deployList, err := kubeClient.AppsV1().Deployments(corev1.NamespaceAll).List(metav1.ListOptions{
			LabelSelector: tillerLabelSelector.String(),
		})
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		count = len(deployList.Items)
		if (err != nil && errors.IsNotFound(err)) || count == 0 {
			return nil
		}
		deployLabels = deployList.Items[0].Labels
	} else {
		deployList, err := kubeClient.ExtensionsV1beta1().Deployments(corev1.NamespaceAll).List(metav1.ListOptions{
			LabelSelector: tillerLabelSelector.String(),
		})
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		count = len(deployList.Items)
		if (err != nil && errors.IsNotFound(err)) || count == 0 {
			return nil
		}
		deployLabels = deployList.Items[0].Labels
	}

	conflictErr := errors.NewConflict(platform.Resource("helm"), "helm", fmt.Errorf("UnOfficial helm"))
	if count > 1 {
		return conflictErr
	}
	if reflect.DeepEqual(deployLabels, selectorForOfficialTiller.MatchLabels) {
		return nil
	}
	return conflictErr
}

func (p *provisioner) checkRsc(kubeClient kubernetes.Interface) error {
	if _, err := apiclient.GetServiceAccount(p.kubeClient, metav1.NamespaceSystem, svcAccountHelmName); err != nil {
		return err
	}
	if _, err := apiclient.GetClusterRoleBinding(p.kubeClient, crbHelmName); err != nil {
		return err
	}
	if _, err := apiclient.GetService(p.kubeClient, metav1.NamespaceSystem, svcHelmAPIName); err != nil {
		return err
	}
	if _, err := apiclient.GetService(p.kubeClient, metav1.NamespaceSystem, svcTillerName); err != nil {
		return err
	}
	if _, err := apiclient.CheckDeployment(p.kubeClient, metav1.NamespaceSystem, deployTillerName); err != nil {
		return err
	}
	if _, err := apiclient.CheckDeployment(p.kubeClient, metav1.NamespaceSystem, deployHelmAPIName); err != nil {
		return err
	}
	return nil
}
