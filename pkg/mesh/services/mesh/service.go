/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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
 *
 */

package mesh

import (
	"context"
	"fmt"
	"sync"

	istionetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	meshconfig "tkestack.io/tke/pkg/mesh/apis/config"
	clusterclient "tkestack.io/tke/pkg/mesh/external/kubernetes"
	"tkestack.io/tke/pkg/mesh/external/tcmesh"
	"tkestack.io/tke/pkg/mesh/services"
	"tkestack.io/tke/pkg/mesh/services/rest"
	"tkestack.io/tke/pkg/mesh/util/constants"
	"tkestack.io/tke/pkg/mesh/util/errors"
	"tkestack.io/tke/pkg/util/log"
)

type meshClusterService struct {
	config    meshconfig.MeshConfiguration
	tcmClient *tcmesh.Client
	clients   clusterclient.Client
}

var _ services.MeshClusterService = &meshClusterService{}

func New(
	config meshconfig.MeshConfiguration, tcmClient *tcmesh.Client, clients clusterclient.Client,
) *meshClusterService {

	return &meshClusterService{
		config:    config,
		tcmClient: tcmClient,
		clients:   clients,
	}
}

// CreateMeshResource
func (m *meshClusterService) CreateMeshResource(ctx context.Context, meshName string, obj *unstructured.Unstructured) error {
	clusters := m.tcmClient.Cache().Clusters(meshName)
	var (
		errs  = errors.NewMultiError()
		size  = len(clusters)
		wg    = &sync.WaitGroup{}
	)
	wg.Add(size)

	for _, clusterName := range clusters {
		go func(clusterName string) {
			defer wg.Done()

			istioClient, err := m.clients.Istio(clusterName)
			if err != nil {
				e := fmt.Errorf("get istio[%s] client failed: %v", clusterName, err)
				log.Errorf("%v", e)
				errs.Add(e)
				return
			}

			err = istioClient.Create(ctx, obj)

			if err != nil {
				log.Errorf("create mesh[%s] cluster[%s] namespace[%s] resource[%s] errors: %v",
					meshName, clusterName, obj.GetNamespace(), obj.GetObjectKind().GroupVersionKind().String(), err)
				errs.Add(err)
			}

		}(clusterName)
	}

	wg.Wait()

	return errs
}

// ListMicroServices list mesh all micro services, which exists the istio "app" label
func (m *meshClusterService) ListMicroServices(
	ctx context.Context, meshName, namespace, serviceName string, selector labels.Selector,
) ([]rest.MicroService, *errors.MultiError) {

	clusters := m.tcmClient.Cache().Clusters(meshName)
	mainClustersMap := m.tcmClient.Cache().MainClustersMap(meshName)
	var (
		errs  = errors.NewMultiError()
		size  = len(clusters)
		wg    = &sync.WaitGroup{}
		retCh = make(chan []rest.MicroService, size)
	)
	wg.Add(size)

	fieldSelector := constants.ExcludeNamespacesSelector
	if serviceName != "" {
		fieldSelector = fields.AndSelectors(fieldSelector, fields.OneTermEqualSelector("metadata.name", serviceName))
	}

	labelSelector := constants.IstioAppSelector
	if selector != nil {
		rs, _ := selector.Requirements()
		labelSelector.Add(rs...)
	}

	for _, clusterName := range clusters {

		go func(clusterName string) {
			defer wg.Done()

			_, isMainCluster := mainClustersMap[clusterName]

			clusterClient, err := m.clients.Cluster(clusterName)
			if err != nil {
				e := fmt.Errorf("get cluster[%s] client failed: %v", clusterName, err)
				log.Errorf("%v", e)
				errs.Add(e)
				return
			}
			istioClient, err := m.clients.Istio(clusterName)
			if err != nil {
				e := fmt.Errorf("get istio[%s] client failed: %v", clusterName, err)
				log.Errorf("%v", e)
				errs.Add(e)
				return
			}

			ret, merr := fetchMicroServices(
				ctx, clusterClient, istioClient,
				meshName, clusterName, namespace,
				serviceName, isMainCluster, nil,
			)

			if len(ret) > 0 {
				retCh <- ret
			}
			if merr != nil && !merr.Empty() {
				log.Errorf("fetched mesh[%s] cluster[%s] namespace[%s] service[%s] errors: %v",
					meshName, clusterName, namespace, serviceName, merr)
				errs.Add(merr)
			}

		}(clusterName)

	}

	wg.Wait()
	close(retCh)

	rets := make([]rest.MicroService, 0)
	for rt := range retCh {
		if len(rt) > 0 {
			rets = append(rets, rt...)
		}
	}
	return rets, errs
}

/*func (m *meshClusterService) GetService(
	ctx context.Context, meshName, namespace, serviceName string,
) (*corev1.Service, error) {

	clusterName := m.tcmClient.Cache().MainClusters(meshName)
	cluster, err := m.clients.Cluster(clusterName)
	if err != nil {
		return nil, err
	}
	objectKey := types.NamespacedName{
		Namespace: namespace,
		Name:      serviceName,
	}
	ret := &corev1.Service{}
	err = cluster.Get(ctx, objectKey, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}*/

// 2020-11-06 implements micro service list
func fetchMicroServices(ctx context.Context, clusterClient ctrlclient.Client, istioClient ctrlclient.Client,
	meshName, clusterName, namespace, serviceName string,
	isMainCluster bool, sel labels.Selector) ([]rest.MicroService, *errors.MultiError) {

	var (
		wg    = &sync.WaitGroup{}
		merrs = errors.NewMultiError()

		role = "remote"

		fieldSelector = constants.ExcludeNamespacesSelector
		labelSelector = constants.IstioAppSelector

		svcs             []corev1.Service
		eps              []corev1.Endpoints
		pods             []corev1.Pod
		deployments      []appsv1.Deployment
		statefulSets     []appsv1.StatefulSet
		virtualServices  []istionetworking.VirtualService
		destinationRules []istionetworking.DestinationRule
		hpas             map[string]autoscaling.HorizontalPodAutoscaler
	)

	if isMainCluster {
		role = "master"
	}

	if sel != nil {
		rs, _ := sel.Requirements()
		labelSelector.Add(rs...)
	}

	wg.Add(8)

	go func() {
		// 1
		defer wg.Done()
		if serviceName != "" {
			fieldSelector = fields.AndSelectors(
				fieldSelector, fields.OneTermEqualSelector("metadata.name", serviceName),
			)
		}
		var err error
		ret := &corev1.ServiceList{}
		err = clusterClient.List(ctx, ret, &ctrlclient.ListOptions{
			LabelSelector: labelSelector,
			Namespace:     namespace,
			FieldSelector: fieldSelector,
		})
		svcs = ret.Items
		if err != nil {
			merrs.Add(err)
		}
	}()
	go func() {
		// 2
		defer wg.Done()
		ret := &corev1.EndpointsList{}
		err := clusterClient.List(ctx, ret, &ctrlclient.ListOptions{
			Namespace:     namespace,
			FieldSelector: fieldSelector,
			LabelSelector: labelSelector,
		})
		eps = ret.Items
		if err != nil {
			merrs.Add(err)
		}
	}()
	go func() {
		// 3
		defer wg.Done()
		ret := &corev1.PodList{}
		err := clusterClient.List(ctx, ret, &ctrlclient.ListOptions{
			Namespace:     namespace,
			FieldSelector: fieldSelector,
			LabelSelector: labelSelector,
		})
		pods = ret.Items
		if err != nil {
			merrs.Add(err)
		}
	}()
	go func() {
		// 4
		defer wg.Done()
		ret := &appsv1.DeploymentList{}
		err := clusterClient.List(ctx, ret, &ctrlclient.ListOptions{
			Namespace:     namespace,
			FieldSelector: fieldSelector,
			LabelSelector: labelSelector,
		})
		deployments = ret.Items
		if err != nil {
			merrs.Add(err)
		}
	}()
	go func() {
		// 5
		defer wg.Done()
		ret := &appsv1.StatefulSetList{}
		err := clusterClient.List(ctx, ret, &ctrlclient.ListOptions{
			Namespace:     namespace,
			FieldSelector: fieldSelector,
			LabelSelector: labelSelector,
		})
		statefulSets = ret.Items
		if err != nil {
			merrs.Add(err)
		}
	}()
	go func() {
		// 6
		defer wg.Done()
		ret := &istionetworking.VirtualServiceList{}
		err := istioClient.List(ctx, ret, &ctrlclient.ListOptions{
			Namespace:     namespace,
			FieldSelector: fieldSelector,
			LabelSelector: labelSelector,
		})
		virtualServices = ret.Items
		if err != nil {
			merrs.Add(err)
		}
	}()
	go func() {
		// 7
		defer wg.Done()
		ret := &istionetworking.DestinationRuleList{}
		err := istioClient.List(ctx, ret, &ctrlclient.ListOptions{
			Namespace:     namespace,
			FieldSelector: fieldSelector,
			LabelSelector: labelSelector,
		})
		destinationRules = ret.Items
		if err != nil {
			merrs.Add(err)
		}
	}()
	go func() {
		// 8
		defer wg.Done()
		ret := &autoscaling.HorizontalPodAutoscalerList{}
		err := clusterClient.List(ctx, ret, &ctrlclient.ListOptions{
			Namespace:     namespace,
			FieldSelector: fieldSelector,
			LabelSelector: labelSelector,
		})
		items := ret.Items
		if err != nil {
			merrs.Add(err)
		}
		for _, it := range items {
			ns := it.Namespace
			n := it.Name
			hpas[hpaKey(ns, n)] = it
		}
	}()

	wg.Wait()

	mss := buildMicroServices(meshName, clusterName, role, svcs, eps,
		pods, deployments, statefulSets, virtualServices, destinationRules, hpas)

	return mss, merrs
}

func buildMicroServices(
	meshName string, clusterName string, role string,
	svcs []corev1.Service, eps []corev1.Endpoints, pods []corev1.Pod,
	deployments []appsv1.Deployment, statefulSets []appsv1.StatefulSet,
	virtualServices []istionetworking.VirtualService,
	destinationRules []istionetworking.DestinationRule,
	hpas map[string]autoscaling.HorizontalPodAutoscaler,
) []rest.MicroService {

	mss := make([]rest.MicroService, 0)
	for _, svc := range svcs {
		var (
			ms = &rest.MicroService{
				Cluster: rest.Cluster{
					MeshName:    meshName,
					ClusterName: clusterName,
					Role:        role,
					Region:      "", // TODO:2020-11-09 cluster region
				},
				Service:          corev1.Service{},
				Endpoints:        corev1.Endpoints{},
				Workloads:        make([]rest.Workload, 0),
				VirtualServices:  make([]istionetworking.VirtualService, 0),
				DestinationRules: make([]istionetworking.DestinationRule, 0),
				Pods:             make([]corev1.Pod, 0),
			}
			svcLabels      = labels.Set(svc.Labels).AsSelector()
			svcPodSelector = labels.Set(svc.Spec.Selector).AsSelector()
			appSelector    = labels.Set{constants.IstioAppLabelKey: svc.Labels[constants.IstioAppLabelKey]}.AsSelector()
		)

		ms.Service = svc
		for _, ep := range eps {
			if svcLabels.Matches(labels.Set(ep.Labels)) {
				ms.Endpoints = ep
				break
			}
		}
		for _, pod := range pods {
			if svcPodSelector.Matches(labels.Set(pod.Labels)) {
				ms.Pods = append(ms.Pods, pod)
			}
		}
		for _, deploy := range deployments {
			if svcPodSelector.Matches(labels.Set(deploy.Spec.Template.Labels)) {
				w := rest.Workload{
					Object: &deploy,
				}
				k := hpaKey(deploy.Namespace, deploy.Name) // hpa name must match Deployment name
				if hpa, ok := hpas[k]; ok {
					w.HPA = hpa
				}
				ms.Workloads = append(ms.Workloads, w)
			}
		}
		for _, sts := range statefulSets {
			if svcPodSelector.Matches(labels.Set(sts.Spec.Template.Labels)) {
				w := rest.Workload{
					Object: &sts,
				}
				k := hpaKey(sts.Namespace, sts.Name) // hpa name must match StatefulSet name
				if hpa, ok := hpas[k]; ok {
					w.HPA = hpa
				}
				ms.Workloads = append(ms.Workloads, w)
			}
		}
		for _, vs := range virtualServices {
			if appSelector.Matches(labels.Set(vs.Labels)) {
				ms.VirtualServices = append(ms.VirtualServices, vs)
			}
		}
		for _, dr := range destinationRules {
			if appSelector.Matches(labels.Set(dr.Labels)) {
				ms.DestinationRules = append(ms.DestinationRules, dr)
			}
		}

		mss = append(mss, *ms)
	}

	return mss
}

func hpaKey(namespace, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}
