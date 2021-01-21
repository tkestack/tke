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

package istio

import (
	"context"
	"fmt"
	"strings"
	"sync"

	istionetworkingapi "istio.io/api/networking/v1alpha3"
	istionetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	meshconfig "tkestack.io/tke/pkg/mesh/apis/config"
	clusterclient "tkestack.io/tke/pkg/mesh/external/kubernetes"
	"tkestack.io/tke/pkg/mesh/models"
	"tkestack.io/tke/pkg/mesh/services"
	"tkestack.io/tke/pkg/mesh/util/constants"
	"tkestack.io/tke/pkg/mesh/util/errors"
	"tkestack.io/tke/pkg/util/log"
)

const (
	appRuntimeLabelKey = "app"
)

type istioService struct {
	config  meshconfig.MeshConfiguration
	clients clusterclient.Client
}

func New(config meshconfig.MeshConfiguration, clients clusterclient.Client) services.IstioService {
	return &istioService{
		config:  config,
		clients: clients,
	}
}

func (c *istioService) ListGateways(
	ctx context.Context, clusterName, namespace string,
) ([]istionetworking.Gateway, error) {

	ret := &istionetworking.GatewayList{}
	err := c.ListResources(ctx, clusterName, ret, &ctrlclient.ListOptions{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	return ret.Items, nil
}

func (c *istioService) ListVirtualServices(
	ctx context.Context, clusterName, namespace string,
) ([]istionetworking.VirtualService, error) {

	ret := &istionetworking.VirtualServiceList{}
	err := c.ListResources(ctx, clusterName, ret, &ctrlclient.ListOptions{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	return ret.Items, nil
}

func (c *istioService) ListDestinationRules(
	ctx context.Context, clusterName, namespace string,
) ([]istionetworking.DestinationRule, error) {

	ret := &istionetworking.DestinationRuleList{}
	err := c.ListResources(ctx, clusterName, ret, &ctrlclient.ListOptions{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	return ret.Items, nil
}

func (c *istioService) ListServiceEntries(
	ctx context.Context, clusterName, namespace string) ([]istionetworking.ServiceEntry, error) {
	ret := &istionetworking.ServiceEntryList{}
	err := c.ListResources(ctx, clusterName, ret, &ctrlclient.ListOptions{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	return ret.Items, nil
}

func (c *istioService) ListWorkloadEntries(
	ctx context.Context, clusterName, namespace string) ([]istionetworking.WorkloadEntry, error) {
	ret := &istionetworking.WorkloadEntryList{}
	err := c.ListResources(ctx, clusterName, ret, &ctrlclient.ListOptions{Namespace: namespace})
	if err != nil {
		return nil, err
	}
	return ret.Items, nil
}

// ListAllResources namespace or kind is empty, query all namespaces and kinds
func (c *istioService) ListAllResources(
	ctx context.Context, clusterName, namespace, kind string, selector labels.Selector,
) (kindMaplist map[string][]unstructured.Unstructured, errs *errors.MultiError) {

	kindMaplist = make(map[string][]unstructured.Unstructured)
	errs = errors.NewMultiError()
	client, err := c.clients.Istio(clusterName)
	if err != nil {
		errs.Add(err)
		return nil, errs
	}

	var (
		gvTypes = constants.IstioNetworkingListGVK
		wg      = &sync.WaitGroup{}
		cnt     = len(gvTypes)
		retCh   = make(chan map[string][]unstructured.Unstructured, cnt)
	)

	wg.Add(cnt)
	opts := []ctrlclient.ListOption{&ctrlclient.ListOptions{
		LabelSelector: selector,
		Namespace:     namespace,
	}}
	for _, gvk := range gvTypes {
		gvk := gvk
		k := gvk.Kind
		if kind != "" && !strings.HasPrefix(k, kind) {
			wg.Done()
			continue
		}
		// 2020-11-03 implements multiple goroutine
		go func() {
			defer wg.Done()
			ret := &unstructured.UnstructuredList{}
			ret.SetGroupVersionKind(gvk)
			e := client.List(ctx, ret, opts...)
			if e != nil {
				log.Errorf("list all Istio %s resources error: %v", k, e)
				errs.Add(e)
				return
			}
			if len(ret.Items) > 0 {
				key := strings.ToLower(strings.Trim(k, "List"))
				retCh <- map[string][]unstructured.Unstructured{key: ret.Items}
			}
		}()
	}
	wg.Wait()
	close(retCh)

	for retMap := range retCh {
		for key, rets := range retMap {
			kindMaplist[key] = rets
		}
	}

	return kindMaplist, errs
}

// ListResources results in obj pointer
func (c *istioService) ListResources(
	ctx context.Context, clusterName string, obj runtime.Object, opt ...ctrlclient.ListOption,
) error {

	client, err := c.clients.Istio(clusterName)
	if err != nil {
		return err
	}
	err = client.List(ctx, obj, opt...)
	if err != nil {
		return err
	}
	return nil
}

func (c *istioService) GetResource(
	ctx context.Context, clusterName string, obj runtime.Object,
) error {

	client, err := c.clients.Istio(clusterName)
	if err != nil {
		return err
	}

	key, err := ctrlclient.ObjectKeyFromObject(obj)
	if err != nil {
		return err
	}

	err = client.Get(ctx, key, obj)
	if err != nil {
		return err
	}
	return nil
}

func (c *istioService) CreateResource(
	ctx context.Context, clusterName string, obj runtime.Object,
) (bool, error) {

	client, err := c.clients.Istio(clusterName)
	if err != nil {
		return false, err
	}
	err = client.Create(ctx, obj)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *istioService) UpdateResource(
	ctx context.Context, clusterName string, obj runtime.Object,
) (bool, error) {

	client, err := c.clients.Istio(clusterName)
	if err != nil {
		return false, err
	}
	err = client.Update(ctx, obj)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *istioService) DeleteResource(
	ctx context.Context, clusterName string, obj runtime.Object,
) (bool, error) {

	client, err := c.clients.Istio(clusterName)
	if err != nil {
		return false, err
	}
	err = client.Delete(ctx, obj)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CreateNorthTrafficGateway
func (c *istioService) CreateNorthTrafficGateway(
	ctx context.Context, clusterName string, obj *models.IstioNetworkingConfig,
) (bool, error) {

	// setting gateway attributes
	gateway := obj.Gateway
	serviceRuntime := gateway.GetLabels()[appRuntimeLabelKey]
	// format resourceName
	if len(serviceRuntime) <= 0 {
		return false, fmt.Errorf("gateway must be attatched to a backend service: [%s]", gateway.GetName())
	}
	newGatewayName := strings.Join([]string{obj.App, serviceRuntime, "gw"}, "-")
	gateway.SetName(newGatewayName)

	// Create gateway resource
	created, err := c.CreateResource(ctx, clusterName, gateway)
	if err != nil {
		return created, fmt.Errorf("created north traffic gateway[%s] failed. Err: %v", gateway.GetName(), err)
	}
	obj.Gateway = gateway

	// Create virtualservice resource
	virtualService := obj.VirtualService
	if virtualService != nil {
		// set metadata
		vsName := strings.Join([]string{obj.App, serviceRuntime, "gw-vs"}, "-")
		labs := virtualService.GetLabels()
		labs[appRuntimeLabelKey] = serviceRuntime
		virtualService.SetName(vsName)
		virtualService.SetLabels(labs)

		// set gateways
		virtualService.Spec.Gateways = []string{gateway.GetName()}

		// set exportto
		virtualService.Spec.ExportTo = []string{"*"}

		// set hosts
		requestedHosts := make([]string, 0)
		servers := gateway.Spec.Servers
		if len(servers) > 0 {
			for _, svr := range servers {
				hosts := svr.Hosts
				if len(hosts) > 0 {
					requestedHosts = append(requestedHosts, hosts...)
				}
			}
		}
		// if virtualService.Spec.Hosts == nil {
		virtualService.Spec.Hosts = requestedHosts
		// }

		// set destination while it is empty
		// currently ,we just support http
		if virtualService.Spec.Http == nil {
			virtualService.Spec.Http = []*istionetworkingapi.HTTPRoute{
				{
					Route: []*istionetworkingapi.HTTPRouteDestination{
						{
							Destination: &istionetworkingapi.Destination{
								Host: serviceRuntime,
							},
						},
					},
				},
			}
		}

		created, err = c.CreateResource(ctx, clusterName, virtualService)
		if err != nil {
			return created, fmt.Errorf("created north traffic virtual service[%s] failed. Err: %v",
				virtualService.GetName(), err)
		}
		obj.VirtualService = virtualService
	}

	return true, nil
}

// UpdateNorthTrafficGateway
func (c *istioService) UpdateNorthTrafficGateway(
	ctx context.Context, clusterName string, obj *models.IstioNetworkingConfig,
) (bool, error) {

	if obj == nil || obj.Gateway == nil {
		return false, fmt.Errorf("gateway resource can not be empty while update north traffic")
	}

	// get gateway resource
	gateway := obj.Gateway
	if err := c.GetResource(ctx, clusterName, gateway); err != nil {
		return false, fmt.Errorf("get nortth traffic gateway[%s] before updating failed. Err: [%v]",
			gateway.GetName(), err)
	}

	// update gateway resource
	updated, err := c.UpdateResource(ctx, clusterName, gateway)
	if err != nil {
		return updated, fmt.Errorf("updated north traffic gateway[%s] failed. Err: %v", gateway.GetName(), err)
	}
	obj.Gateway = gateway

	// update virtualservice resource
	virtualService := obj.VirtualService
	if virtualService != nil {

		// set ExportTo
		virtualService.Spec.ExportTo = []string{"*"}

		// set hosts
		requestedHosts := make([]string, 0)
		servers := gateway.Spec.Servers
		if len(servers) > 0 {
			for _, svr := range servers {
				hosts := svr.Hosts
				if len(hosts) > 0 {
					requestedHosts = append(requestedHosts, hosts...)
				}
			}
		}
		// if virtualService.Spec.Hosts == nil {
		virtualService.Spec.Hosts = requestedHosts
		// }

		if updated, err = c.UpdateResource(ctx, clusterName, virtualService); err != nil {
			return updated, fmt.Errorf("updated north traffic virtual service[%s] failed. Err: %v",
				virtualService.GetName(), err)
		}
		obj.VirtualService = virtualService
	}

	return true, nil
}

// DeleteNorthTrafficGateway
func (c *istioService) DeleteNorthTrafficGateway(
	ctx context.Context, clusterName string, gateway *istionetworking.Gateway,
) (bool, error) {

	err := c.GetResource(ctx, clusterName, gateway)
	if err != nil {
		return false, fmt.Errorf("get north traffic gateway[%s] before deleting failed. Err: [%v]",
			gateway.GetName(), err)
	}

	vsName := strings.Join([]string{gateway.GetName(), "vs"}, "-")
	vs := &istionetworking.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vsName,
			Namespace: gateway.GetNamespace(),
		},
	}

	err = c.GetResource(ctx, clusterName, vs)
	if err != nil {
		return false, fmt.Errorf("get north traffic virtual service[%s] before deleting failed. Err: [%v]",
			vs.GetName(), err)
	}

	deleted, err := c.DeleteResource(ctx, clusterName, gateway)
	if err != nil {
		return deleted, fmt.Errorf("deleted north traffic gateway[%s] failed. Err: [%v]",
			gateway.GetName(), err)
	}
	deleted, err = c.DeleteResource(ctx, clusterName, vs)
	if err != nil {
		return deleted, fmt.Errorf("deleted north traffic virtual service[%s] failed. Err: [%v]",
			vs.GetName(), err)
	}

	return true, nil
}

func (c *istioService) GetNorthTrafficGateway(
	ctx context.Context, clusterName string, gateway *istionetworking.Gateway,
) (*models.IstioNetworkingConfig, error) {

	ret := &models.IstioNetworkingConfig{}
	err := c.GetResource(ctx, clusterName, gateway)
	if err != nil {
		return ret, fmt.Errorf("get north traffic gateway[%s] failed. Err: [%v]", gateway.GetName(), err)
	}

	vsName := strings.Join([]string{gateway.GetName(), "vs"}, "-")
	vs := &istionetworking.VirtualService{}
	vs.SetNamespace(gateway.GetNamespace())
	vs.SetName(vsName)
	err = c.GetResource(ctx, clusterName, vs)
	if err != nil {
		return ret, fmt.Errorf("get north traffic virtual service[%s] failed. Err: [%v]", vs.GetName(), err)
	}

	ret.Gateway = gateway
	ret.VirtualService = vs
	return ret, nil
}
