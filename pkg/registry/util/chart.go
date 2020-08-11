package util

import (
	"context"

	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	authversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	businessversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	registryapi "tkestack.io/tke/api/registry"
)

func commonListAndMerge(ctx context.Context, options *metainternal.ListOptions, store *registry.Store, chartGroupList *registryapi.ChartGroupList) (runtime.Object, error) {
	allList := &registryapi.ChartList{}
	for k, v := range chartGroupList.Items {
		// ctx with namespace, Look through
		// https://github.com/kubernetes/apiserver/blob/release-1.17/pkg/endpoints/installer.go#L1083
		// https://github.com/kubernetes/apiserver/blob/release-1.17/pkg/endpoints/handlers/get.go#L186
		ctx = request.WithNamespace(ctx, v.Name)
		obj, err := store.List(ctx, options)
		if err != nil {
			return obj, err
		}
		chartList := obj.(*registryapi.ChartList)
		if k == 0 {
			allList = chartList
		} else {
			allList.Items = append(allList.Items, chartList.Items...)
		}
	}
	return allList, nil
}

// ListPersonalChartsFromStore list all charts that belongs to personal chartgroup
func ListPersonalChartsFromStore(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	obj, err := ListPersonalChartGroups(ctx, options.DeepCopy(), businessClient, registryClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	chartGroupList := obj.(*registryapi.ChartGroupList)
	if len(chartGroupList.Items) == 0 {
		return &registryapi.ChartList{}, nil
	}
	return commonListAndMerge(ctx, options, store, chartGroupList)
}

// ListProjectChartsFromStore list all charts that belongs to project chartgroup
func ListProjectChartsFromStore(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	authClient authversionedclient.AuthV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	obj, err := ListProjectChartGroups(ctx, options.DeepCopy(), businessClient, authClient, registryClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	chartGroupList := obj.(*registryapi.ChartGroupList)
	if len(chartGroupList.Items) == 0 {
		return &registryapi.ChartList{}, nil
	}

	return commonListAndMerge(ctx, options, store, chartGroupList)
}

// ListSystemChartsFromStore list all charts that belongs to system chartgroup
func ListSystemChartsFromStore(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	obj, err := ListSystemChartGroups(ctx, options.DeepCopy(), businessClient, registryClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	chartGroupList := obj.(*registryapi.ChartGroupList)
	if len(chartGroupList.Items) == 0 {
		return &registryapi.ChartList{}, nil
	}

	return commonListAndMerge(ctx, options, store, chartGroupList)
}

// ListPublicChartsFromStore list all charts that belongs to public chartgroup
func ListPublicChartsFromStore(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	obj, err := ListPublicChartGroups(ctx, options.DeepCopy(), businessClient, registryClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	chartGroupList := obj.(*registryapi.ChartGroupList)
	if len(chartGroupList.Items) == 0 {
		return &registryapi.ChartList{}, nil
	}
	return commonListAndMerge(ctx, options, store, chartGroupList)
}

func mergeCharts(cgs ...runtime.Object) *registryapi.ChartList {
	list := &registryapi.ChartList{}
	exist := make(map[types.UID]bool)
	for _, v := range cgs {
		if v != nil {
			vList := v.(*registryapi.ChartList)
			for _, v := range vList.Items {
				if _, ok := exist[v.UID]; !ok {
					list.Items = append(list.Items, v)
					exist[v.UID] = true
				}
			}
		}
	}
	return list
}

// ListAllChartsFromStore list all charts
func ListAllChartsFromStore(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	authClient authversionedclient.AuthV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	personal, err := ListPersonalChartsFromStore(ctx, options.DeepCopy(), businessClient, registryClient, privilegedUsername, store)
	if err != nil {
		return nil, err
	}
	project, err := ListProjectChartsFromStore(ctx, options.DeepCopy(), businessClient, authClient, registryClient, privilegedUsername, store)
	if err != nil {
		return nil, err
	}
	public, err := ListPublicChartsFromStore(ctx, options.DeepCopy(), businessClient, registryClient, privilegedUsername, store)
	if err != nil {
		return nil, err
	}
	list := mergeCharts(personal, project, public)
	return list, nil
}
