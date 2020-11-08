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

func isAdmin(ctx context.Context, businessClient businessversionedclient.BusinessV1Interface, privilegedUsername string) (bool, error) {
	admin := authentication.IsAdministrator(ctx, privilegedUsername)
	if admin {
		return true, nil
	}
	return isPlatformAdmin(ctx, businessClient)
}

//TODO: fix here, use util-function to judge isPlatformAdmin
func isPlatformAdmin(ctx context.Context, businessClient businessversionedclient.BusinessV1Interface) (bool, error) {
	// if we don't use business component, we are all admin
	if businessClient == nil {
		return true, nil
	}
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

// filterUserChartGroups filter all chartgroups that belongs to user
func filterUserChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	username string,
	cgList *registryapi.ChartGroupList,
	isAdmin bool) (runtime.Object, error) {
	if !isAdmin {
		allList := []registryapi.ChartGroup{}
		for _, cg := range cgList.Items {
			if util.InStringSlice(cg.Spec.Users, username) {
				allList = append(allList, cg)
				break
			}
		}
		cgList.Items = allList
	}

	return cgList, nil
}

// ListUserChartGroups list all chartgroups that belongs to personal
func prepareListUserChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	privilegedUsername string) (opt *metainternal.ListOptions, isAdmin bool, err error) {
	admin := authentication.IsAdministrator(ctx, privilegedUsername)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	platformAdmin, err := isPlatformAdmin(ctx, businessClient)
	if err != nil {
		return nil, true, err
	}

	fieldSelector := fields.OneTermEqualSelector("spec.visibility", string(registryapi.VisibilityUser))
	if admin {
		// do nothing, super admin has no tenantID
	} else {
		fieldSelector = fields.AndSelectors(fieldSelector, fields.OneTermEqualSelector("spec.tenantID", tenantID))
	}
	options = apiserverutil.FullListOptionsFieldSelector(options, fieldSelector)
	return options, admin || platformAdmin, nil
}

// ListUserChartGroupsFromStore list all chartgroups that belongs to users
func ListUserChartGroupsFromStore(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	options, isAdmin, err := prepareListUserChartGroups(ctx, options, businessClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	obj, err := store.List(ctx, options)
	if err != nil {
		return nil, err
	}
	cgList := obj.(*registryapi.ChartGroupList)
	if len(cgList.Items) == 0 {
		return &registryapi.ChartGroupList{}, nil
	}
	username, _ := authentication.UsernameAndTenantID(ctx)
	return filterUserChartGroups(ctx, options, username, cgList, isAdmin)
}

// ListUserChartGroups list all chartgroups that belongs to users
func ListUserChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string) (runtime.Object, error) {
	options, isAdmin, err := prepareListUserChartGroups(ctx, options, businessClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	cgList, err := registryClient.ChartGroups().List(ctx, metav1.ListOptions{
		FieldSelector: options.FieldSelector.String(),
	})
	if err != nil {
		return nil, err
	}
	if len(cgList.Items) == 0 {
		return &registryapi.ChartGroupList{}, nil
	}
	username, _ := authentication.UsernameAndTenantID(ctx)
	return filterUserChartGroups(ctx, options, username, cgList, isAdmin)
}

// filterProjectChartGroups filter all charts that belongs to project
func filterProjectChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	targetProjectID string,
	authClient authversionedclient.AuthV1Interface,
	cgList *registryapi.ChartGroupList,
	isAdmin bool) (runtime.Object, error) {
	targetPrj := func(prjs []string, prj string) bool {
		return prj == "" || util.InStringSlice(prjs, prj)
	}

	if !isAdmin && authClient != nil {
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
		for _, cg := range cgList.Items {
			if targetPrj(cg.Spec.Projects, targetProjectID) {
				for _, pid := range cg.Spec.Projects {
					if _, ok := pidMap[pid]; ok {
						allList = append(allList, cg)
						break
					}
				}
			}
		}
		cgList.Items = allList
	} else {
		allList := []registryapi.ChartGroup{}
		for _, cg := range cgList.Items {
			if targetPrj(cg.Spec.Projects, targetProjectID) {
				allList = append(allList, cg)
			}
		}
		cgList.Items = allList
	}

	return cgList, nil
}

// prepareListProjectChartGroups list all charts that belongs to project
func prepareListProjectChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	businessClient businessversionedclient.BusinessV1Interface,
	privilegedUsername string) (opt *metainternal.ListOptions, isAdmin bool, err error) {
	admin := authentication.IsAdministrator(ctx, privilegedUsername)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	platformAdmin, err := isPlatformAdmin(ctx, businessClient)
	if err != nil {
		return nil, true, err
	}

	fieldSelector := fields.OneTermEqualSelector("spec.visibility", string(registryapi.VisibilityProject))
	if admin {
		// do nothing, admin has no tenantID
	} else {
		fieldSelector = fields.AndSelectors(fieldSelector, fields.OneTermEqualSelector("spec.tenantID", tenantID))
	}

	options = apiserverutil.FullListOptionsFieldSelector(options, fieldSelector)
	return options, admin || platformAdmin, nil
}

// ListProjectChartGroupsFromStore list all charts that belongs to project
func ListProjectChartGroupsFromStore(ctx context.Context,
	options *metainternal.ListOptions,
	targetProjectID string,
	businessClient businessversionedclient.BusinessV1Interface,
	authClient authversionedclient.AuthV1Interface,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	options, isAdmin, err := prepareListProjectChartGroups(ctx, options, businessClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	obj, err := store.List(ctx, options)
	if err != nil {
		return nil, err
	}
	cgList := obj.(*registryapi.ChartGroupList)
	if len(cgList.Items) == 0 {
		return &registryapi.ChartGroupList{}, nil
	}
	return filterProjectChartGroups(ctx, options, targetProjectID, authClient, cgList, isAdmin)
}

// ListProjectChartGroups list all charts that belongs to project
func ListProjectChartGroups(ctx context.Context,
	options *metainternal.ListOptions,
	targetProjectID string,
	businessClient businessversionedclient.BusinessV1Interface,
	authClient authversionedclient.AuthV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string) (runtime.Object, error) {
	options, isAdmin, err := prepareListProjectChartGroups(ctx, options, businessClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	cgList, err := registryClient.ChartGroups().List(ctx, metav1.ListOptions{
		FieldSelector: options.FieldSelector.String(),
	})
	if err != nil {
		return nil, err
	}
	if len(cgList.Items) == 0 {
		return &registryapi.ChartGroupList{}, nil
	}
	return filterProjectChartGroups(ctx, options, targetProjectID, authClient, cgList, isAdmin)
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
	targetProjectID string,
	businessClient businessversionedclient.BusinessV1Interface,
	authClient authversionedclient.AuthV1Interface,
	privilegedUsername string,
	store *registry.Store) (runtime.Object, error) {
	personal, err := ListUserChartGroupsFromStore(ctx, options.DeepCopy(), businessClient, privilegedUsername, store)
	if err != nil {
		return nil, err
	}
	project, err := ListProjectChartGroupsFromStore(ctx, options.DeepCopy(), targetProjectID, businessClient, authClient, privilegedUsername, store)
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
	targetProjectID string,
	businessClient businessversionedclient.BusinessV1Interface,
	authClient authversionedclient.AuthV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	privilegedUsername string) (runtime.Object, error) {
	personal, err := ListUserChartGroups(ctx, options.DeepCopy(), businessClient, registryClient, privilegedUsername)
	if err != nil {
		return nil, err
	}
	project, err := ListProjectChartGroups(ctx, options.DeepCopy(), targetProjectID, businessClient, authClient, registryClient, privilegedUsername)
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
