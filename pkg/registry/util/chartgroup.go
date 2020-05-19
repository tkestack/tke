package util

import (
	"context"
	"fmt"

	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	authv1 "tkestack.io/tke/api/auth/v1"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	authversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	businessversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	registryapi "tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/apiserver/filter"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	"tkestack.io/tke/pkg/util"
)

//TODO: fix here, use util-function to judge isPlatformAdmin
func isPlatformAdmin(ctx context.Context, businessClient businessversionedclient.BusinessV1Interface) (bool, error) {
	isPlatformAdmin := false
	username, tenantID := authentication.UsernameAndTenantID(ctx)
	platformList, err := businessClient.Platforms().List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s", tenantID),
	})
	if err != nil {
		return isPlatformAdmin, err
	}
	for _, platform := range platformList.Items {
		if util.InStringSlice(platform.Spec.Administrators, username) {
			isPlatformAdmin = true
			break
		}
	}
	return isPlatformAdmin, nil
}

// ListPersonalChartGroups list all chartgroups that belongs to personal
func prepareListPersonalChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	privilegedUsername string) (opt *metainternal.ListOptions, err error) {
	admin := authentication.IsAdministrator(ctx, privilegedUsername)
	fieldSelector := fields.OneTermEqualSelector("spec.type", string(registryapi.RepoTypePersonal))
	if admin {
		// do nothing, super admin has no tenantID
	} else {
		username, tenantID := authentication.UsernameAndTenantID(ctx)
		fieldSelector = fields.AndSelectors(fieldSelector, fields.OneTermEqualSelector("spec.tenantID", tenantID))

		platformAdmin, err := isPlatformAdmin(ctx, businessClient)
		if err != nil {
			return options, err
		}
		if !platformAdmin {
			fieldSelector = fields.AndSelectors(fieldSelector, fields.OneTermEqualSelector("spec.name", username))
		}
	}
	options = apiserverutil.FullListOptionsFieldSelector(options, fieldSelector)
	return options, nil
}

// ListPersonalChartGroupsFromStore list all chartgroups that belongs to personal
func ListPersonalChartGroupsFromStore(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	options, err := prepareListPersonalChartGroups(ctx, options, businessClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	return store.List(ctx, options)
}

// ListPersonalChartGroups list all chartgroups that belongs to personal
func ListPersonalChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string) (runtime.Object, error) {
	options, err := prepareListPersonalChartGroups(ctx, options, businessClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	return registryClient.ChartGroups().List(ctx, metav1.ListOptions{
		FieldSelector: options.FieldSelector.String(),
	})
}

// prepareListProjectChartGroups list all charts that belongs to project
func prepareListProjectChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	privilegedUsername string) (opt *metainternal.ListOptions, admin bool, platformAdmin bool, err error) {
	admin = authentication.IsAdministrator(ctx, privilegedUsername)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	platformAdmin, err = isPlatformAdmin(ctx, businessClient)
	if err != nil {
		return nil, admin, platformAdmin, err
	}

	fieldSelector := fields.OneTermEqualSelector("spec.type", string(registryapi.RepoTypeProject))
	if admin {
		// do nothing, admin has no tenantID
	} else {
		fieldSelector = fields.AndSelectors(fieldSelector, fields.OneTermEqualSelector("spec.tenantID", tenantID))
	}

	options = apiserverutil.FullListOptionsFieldSelector(options, fieldSelector)
	return options, admin, platformAdmin, nil
}

// filterProjectChartGroups filter all charts that belongs to project
func filterProjectChartGroups(ctx context.Context,
	authClient authversionedclient.AuthV1Interface,
	chartGroupList *registryapi.ChartGroupList,
	filterFromProjectBelongs bool) (runtime.Object, error) {

	targetPrj := func(prjs []string, prj string) bool {
		return prj == "" || util.InStringSlice(prjs, prj)
	}
	targetProjectID := filter.ProjectIDFrom(ctx)

	if filterFromProjectBelongs {
		uid := authentication.GetUID(ctx)
		_, tenantID := authentication.UsernameAndTenantID(ctx)
		belongs := &authv1.ProjectBelongs{}
		err := authClient.RESTClient().Get().
			Resource("users").
			Name(uid).
			SubResource("projects").
			SetHeader(filter.HeaderTenantID, tenantID).
			Do(ctx).Into(belongs)
		if err != nil {
			return nil, err
		}
		pidMap := make(map[string]bool)
		for pid := range belongs.MemberdProjects {
			pidMap[pid] = true
		}
		for pid := range belongs.ManagedProjects {
			pidMap[pid] = true
		}
		allList := []registryapi.ChartGroup{}
		for _, cg := range chartGroupList.Items {
			if targetPrj(cg.Spec.Projects, targetProjectID) {
				for _, pid := range cg.Spec.Projects {
					if _, ok := pidMap[pid]; ok {
						allList = append(allList, cg)
						break
					}
				}
			}
		}
		chartGroupList.Items = allList
	} else {
		allList := []registryapi.ChartGroup{}
		for _, cg := range chartGroupList.Items {
			if targetPrj(cg.Spec.Projects, targetProjectID) {
				allList = append(allList, cg)
			}
		}
		chartGroupList.Items = allList
	}

	return chartGroupList, nil
}

// ListProjectChartGroupsFromStore list all charts that belongs to project
func ListProjectChartGroupsFromStore(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	authClient authversionedclient.AuthV1Interface,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	options, admin, platformAdmin, err := prepareListProjectChartGroups(ctx, options, businessClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	obj, err := store.List(ctx, options)
	if err != nil {
		return nil, err
	}
	chartGroupList := obj.(*registryapi.ChartGroupList)
	if len(chartGroupList.Items) == 0 {
		return &registryapi.ChartGroupList{}, nil
	}
	return filterProjectChartGroups(ctx, authClient, chartGroupList, (!admin && !platformAdmin))
}

// ListProjectChartGroups list all charts that belongs to project
func ListProjectChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	authClient authversionedclient.AuthV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string) (runtime.Object, error) {
	options, admin, platformAdmin, err := prepareListProjectChartGroups(ctx, options, businessClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	chartGroupList, err := registryClient.ChartGroups().List(ctx, metav1.ListOptions{
		FieldSelector: options.FieldSelector.String(),
	})
	if err != nil {
		return nil, err
	}
	if len(chartGroupList.Items) == 0 {
		return &registryapi.ChartGroupList{}, nil
	}
	return filterProjectChartGroups(ctx, authClient, chartGroupList, (!admin && !platformAdmin))
}

// prepareListSystemChartGroups list all charts that belongs to system
func prepareListSystemChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	privilegedUsername string) (opt *metainternal.ListOptions, err error) {
	admin := authentication.IsAdministrator(ctx, privilegedUsername)

	fieldSelector := fields.OneTermEqualSelector("spec.type", string(registryapi.RepoTypeSystem))
	if admin {
		// do nothing, admin has no tenantID
	} else {
		_, tenantID := authentication.UsernameAndTenantID(ctx)
		fieldSelector = fields.AndSelectors(fieldSelector, fields.OneTermEqualSelector("spec.tenantID", tenantID))
	}
	options = apiserverutil.FullListOptionsFieldSelector(options, fieldSelector)
	return options, nil
}

// ListSystemChartGroupsFromStore list all charts that belongs to system
func ListSystemChartGroupsFromStore(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	options, err := prepareListSystemChartGroups(ctx, options, businessClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	return store.List(ctx, options)
}

// ListSystemChartGroups list all charts that belongs to system
func ListSystemChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string) (runtime.Object, error) {
	options, err := prepareListSystemChartGroups(ctx, options, businessClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	return registryClient.ChartGroups().List(ctx, metav1.ListOptions{
		FieldSelector: options.FieldSelector.String(),
	})
}

// ListPublicChartGroups list all charts that belongs to public
func prepareListPublicChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	privilegedUsername string) (opt *metainternal.ListOptions, err error) {
	admin := authentication.IsAdministrator(ctx, privilegedUsername)

	fieldSelector := fields.OneTermEqualSelector("spec.visibility", string(registryapi.VisibilityPublic))
	if admin {
		// do nothing, admin has no tenantID
	} else {
		_, tenantID := authentication.UsernameAndTenantID(ctx)
		fieldSelector = fields.AndSelectors(fieldSelector, fields.OneTermEqualSelector("spec.tenantID", tenantID))
	}
	options = apiserverutil.FullListOptionsFieldSelector(options, fieldSelector)
	return options, nil
}

// ListPublicChartGroupsFromStore list all charts that belongs to public
func ListPublicChartGroupsFromStore(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	options, err := prepareListPublicChartGroups(ctx, options, businessClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	return store.List(ctx, options)
}

// ListPublicChartGroups list all charts that belongs to public
func ListPublicChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string) (runtime.Object, error) {
	options, err := prepareListPublicChartGroups(ctx, options, businessClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	return registryClient.ChartGroups().List(ctx, metav1.ListOptions{
		FieldSelector: options.FieldSelector.String(),
	})
}

func mergeChartgroups(cgs ...runtime.Object) *registryapi.ChartGroupList {
	list := &registryapi.ChartGroupList{}
	exist := make(map[types.UID]bool)
	for _, v := range cgs {
		if v != nil {
			vList := v.(*registryapi.ChartGroupList)
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

// ListAllChartGroupsFromStore list all chartgroups
func ListAllChartGroupsFromStore(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	authClient authversionedclient.AuthV1Interface,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	personal, err := ListPersonalChartGroupsFromStore(ctx, options.DeepCopy(), businessClient, privilegedUsername, store)
	if err != nil {
		return nil, err
	}
	project, err := ListProjectChartGroupsFromStore(ctx, options.DeepCopy(), businessClient, authClient, privilegedUsername, store)
	if err != nil {
		return nil, err
	}
	public, err := ListPublicChartGroupsFromStore(ctx, options.DeepCopy(), businessClient, privilegedUsername, store)
	if err != nil {
		return nil, err
	}
	list := mergeChartgroups(personal, project, public)
	return list, nil
}

// ListAllChartGroups list all chartgroups
func ListAllChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	authClient authversionedclient.AuthV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string) (runtime.Object, error) {
	personal, err := ListPersonalChartGroups(ctx, options.DeepCopy(), businessClient, registryClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	project, err := ListProjectChartGroups(ctx, options.DeepCopy(), businessClient, authClient, registryClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	public, err := ListPublicChartGroups(ctx, options.DeepCopy(), businessClient, registryClient, privilegedUsername)
	if err != nil {
		return nil, err
	}

	list := mergeChartgroups(personal, project, public)
	return list, nil
}
