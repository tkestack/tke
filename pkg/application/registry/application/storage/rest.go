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

package storage

import (
	"context"
	"fmt"
	"reflect"

	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/api/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/application"
	applicationapi "tkestack.io/tke/api/application"
	v1 "tkestack.io/tke/api/application/v1"
	applicationinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/application/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	registryversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	registryv1 "tkestack.io/tke/api/registry/v1"
	appconfig "tkestack.io/tke/pkg/application/config"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	helmutil "tkestack.io/tke/pkg/application/helm/util"
	applicationstrategy "tkestack.io/tke/pkg/application/registry/application"
	"tkestack.io/tke/pkg/application/util"
	"tkestack.io/tke/pkg/application/util/chartpath"
	authorizationutil "tkestack.io/tke/pkg/registry/util/authorization"
)

// REST adapts a service registry into apiserver's RESTStorage model.
type REST struct {
	application       ApplicationStorage
	applicationClient *applicationinternalclient.ApplicationClient
	platformClient    platformversionedclient.PlatformV1Interface
	registryClient    registryversionedclient.RegistryV1Interface
	authorizer        authorizer.Authorizer
	repo              appconfig.RepoConfiguration
}

type ApplicationStorage interface {
	rest.Scoper
	rest.Getter
	rest.Lister
	rest.CreaterUpdater
	rest.GracefulDeleter
	rest.Watcher
	rest.StorageVersionProvider
}

// NewREST returns a wrapper around the underlying generic storage and performs
// allocations and deallocations of various helm releases related resources like ports.
// TODO: all transactional behavior should be supported from within generic storage
//   or the strategy.
func NewREST(
	application ApplicationStorage,
	applicationClient *applicationinternalclient.ApplicationClient,
	platformClient platformversionedclient.PlatformV1Interface,
	registryClient registryversionedclient.RegistryV1Interface,
	authorizer authorizer.Authorizer,
	repo appconfig.RepoConfiguration,
) *REST {
	rest := &REST{
		application:       application,
		applicationClient: applicationClient,
		platformClient:    platformClient,
		registryClient:    registryClient,
		authorizer:        authorizer,
		repo:              repo,
	}
	return rest
}

var (
	_ ApplicationStorage          = &REST{}
	_ rest.ShortNamesProvider     = &REST{}
	_ rest.StorageVersionProvider = &REST{}
)

func (rs *REST) StorageVersion() runtime.GroupVersioner {
	return rs.application.StorageVersion()
}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (rs *REST) ShortNames() []string {
	return []string{"app"}
}

func (rs *REST) NamespaceScoped() bool {
	return rs.application.NamespaceScoped()
}

func (rs *REST) New() runtime.Object {
	return rs.application.New()
}

func (rs *REST) NewList() runtime.Object {
	return rs.application.NewList()
}

func (rs *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return rs.application.Get(ctx, name, options)
}

func (rs *REST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	return rs.application.List(ctx, options)
}

func (rs *REST) Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	return rs.application.Watch(ctx, options)
}

func (rs *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	app := obj.(*application.App)

	err := rs.prepareForCheck(ctx, app)
	if err != nil {
		return nil, err
	}

	if app.Spec.DryRun {
		rel, err := rs.dryRun(ctx, app)
		if err != nil {
			return nil, err
		}
		ret := app.DeepCopy()
		ret.Status.Manifest = rel.Manifest
		return ret, nil
	}

	if rs.registryClient != nil {
		chartGroup, err := rs.getChartGroup(ctx, app)
		if err != nil {
			return nil, err
		}
		app.Spec.Chart, err = chartpath.FullfillChartInfo(app.Spec.Chart, chartGroup)
		if err != nil {
			return nil, errors.NewInternalError(err)
		}
	}
	return rs.application.Create(ctx, obj, createValidation, options)
}

func (rs *REST) Delete(ctx context.Context, id string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	// TODO: handle graceful
	obj, _, err := rs.application.Delete(ctx, id, deleteValidation, options)
	if err != nil {
		return nil, false, err
	}

	return obj, true, nil
}

func (rs *REST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	oldObj, err := rs.application.Get(ctx, name, &metav1.GetOptions{})
	if err != nil {
		// Support create on update, if forced to.
		if forceAllowCreate {
			obj, err := objInfo.UpdatedObject(ctx, nil)
			if err != nil {
				return nil, false, err
			}
			createdObj, err := rs.Create(ctx, obj, createValidation, &metav1.CreateOptions{DryRun: options.DryRun})
			if err != nil {
				return nil, false, err
			}
			return createdObj, true, nil
		}
		return nil, false, err
	}
	oldApp := oldObj.(*application.App)

	obj, err := objInfo.UpdatedObject(ctx, oldApp)
	if err != nil {
		return nil, false, err
	}
	app := obj.(*application.App)

	requestNamespace, _ := genericapirequest.NamespaceFrom(ctx)
	err = rest.EnsureObjectNamespaceMatchesRequestNamespace(requestNamespace, &app.ObjectMeta)
	if err != nil {
		return nil, false, errors.NewConflict(applicationapi.Resource("apps"), app.Namespace, fmt.Errorf("App.Namespace does not match the provided context"))
	}

	// check chart permission
	// check value format
	err = rs.prepareForCheck(ctx, app)
	if err != nil {
		return nil, false, err
	}

	if app.Spec.DryRun {
		rel, err := rs.dryRun(ctx, app)
		if err != nil {
			return nil, false, err
		}
		ret := app.DeepCopy()
		ret.Status.Manifest = rel.Manifest
		return ret, true, nil
	}

	// Copy over non-user fields
	strategy := applicationstrategy.NewStrategy(rs.applicationClient)
	if err := rest.BeforeUpdate(strategy, ctx, app, oldApp); err != nil {
		return nil, false, err
	}

	if rs.registryClient != nil {
		chartGroup, err := rs.getChartGroup(ctx, app)
		if err != nil {
			return nil, false, err
		}
		app.Spec.Chart, err = chartpath.FullfillChartInfo(app.Spec.Chart, chartGroup)
		if err != nil {
			return nil, false, errors.NewInternalError(err)
		}
	}

	if !reflect.DeepEqual(oldApp.Spec, app.Spec) && app.Status.Phase != applicationapi.AppPhaseRolledBack && oldApp.Status.Phase != applicationapi.AppPhaseTerminating {
		app.Status.Phase = applicationapi.AppPhaseUpgrading
	}

	return rs.application.Update(ctx, name, rest.DefaultUpdatedObjectInfo(app), createValidation, updateValidation, forceAllowCreate, options)
}

func (rs *REST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return rs.application.ConvertToTable(ctx, object, tableOptions)
}

func (rs *REST) prepareForCheck(ctx context.Context, app *application.App) error {
	// check chart permission
	err := rs.canVisitChart(ctx, app)
	if err != nil {
		return err
	}

	// check value format
	_, err = helmutil.MergeValues(app.Spec.Values.Values, app.Spec.Values.RawValues, string(app.Spec.Values.RawValuesType))
	if err != nil {
		return errors.NewBadRequest(err.Error())
	}
	return nil
}

func (rs *REST) canVisitChart(ctx context.Context, app *application.App) error {
	// TODO: allowAlways if registryClient is empty?
	if rs.registryClient == nil {
		return nil
	}

	chartGroup, err := rs.getChartGroup(ctx, app)
	if err != nil {
		return err
	}
	chart, err := rs.getChart(ctx, app, &chartGroup)
	if err != nil {
		return err
	}

	u, exist := genericapirequest.UserFrom(ctx)
	if !exist || u == nil {
		return errors.NewUnauthorized("empty user info, not authenticated")
	}
	authorized, err := authorizationutil.AuthorizeForChart(ctx, u, rs.authorizer, "get", chartGroup, chart.Name)
	if err != nil {
		return err
	}
	if !authorized {
		return errors.NewForbidden(registryv1.Resource("charts"), "not authenticated", fmt.Errorf("can not get chart: %s/%s/%s", app.Spec.Chart.TenantID, app.Spec.Chart.ChartGroupName, app.Spec.Chart.ChartName))
	}
	return nil
}

func (rs *REST) getChartGroup(ctx context.Context, app *application.App) (registryv1.ChartGroup, error) {
	if rs.registryClient == nil {
		return registryv1.ChartGroup{}, nil
	}

	chartGroupList, err := rs.registryClient.ChartGroups().List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", app.Spec.Chart.TenantID, app.Spec.Chart.ChartGroupName),
	})
	if err != nil {
		return registryv1.ChartGroup{}, errors.NewInternalError(err)
	}
	if len(chartGroupList.Items) == 0 {
		return registryv1.ChartGroup{}, errors.NewNotFound(registryv1.Resource("chartgroups"), fmt.Sprintf("%s/%s", app.Spec.Chart.TenantID, app.Spec.Chart.ChartGroupName))
	}
	chartGroup := chartGroupList.Items[0]
	return chartGroup, nil
}

func (rs *REST) getChart(ctx context.Context, app *application.App, cg *registryv1.ChartGroup) (registryv1.Chart, error) {
	if rs.registryClient == nil {
		return registryv1.Chart{}, nil
	}

	chartList, err := rs.registryClient.Charts(cg.ObjectMeta.Name).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", app.Spec.Chart.TenantID, app.Spec.Chart.ChartName),
	})
	if err != nil {
		return registryv1.Chart{}, errors.NewInternalError(err)
	}
	if len(chartList.Items) == 0 {
		return registryv1.Chart{}, errors.NewNotFound(registryv1.Resource("charts"), fmt.Sprintf("%s/%s/%s", app.Spec.Chart.TenantID, app.Spec.Chart.ChartGroupName, app.Spec.Chart.ChartName))
	}
	chart := chartList.Items[0]
	return chart, nil
}

func (rs *REST) dryRun(ctx context.Context, app *application.App) (*release.Release, error) {
	chartGroup, err := rs.getChartGroup(ctx, app)
	if err != nil {
		return nil, err
	}
	appChart, err := chartpath.FullfillChartInfo(app.Spec.Chart, chartGroup)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	chartPathBasicOptions, err := chartpath.BuildChartPathBasicOptions(rs.repo, appChart)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	appv1 := &v1.App{}
	if err := v1.Convert_application_App_To_v1_App(app, appv1, nil); err != nil {
		return nil, err
	}
	client, err := util.NewHelmClientWithProvider(ctx, rs.platformClient, appv1)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}
	destfile, err := client.Pull(&helmaction.PullOptions{
		ChartPathOptions: chartPathBasicOptions,
	})
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	// check value format
	values, err := helmutil.MergeValues(app.Spec.Values.Values, app.Spec.Values.RawValues, string(app.Spec.Values.RawValuesType))
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	chartPathBasicOptions.ExistedFile = destfile
	rel, err := client.Install(ctx, &helmaction.InstallOptions{
		Namespace:        app.Spec.TargetNamespace,
		ReleaseName:      app.Spec.Name,
		DependencyUpdate: true,
		DryRun:           true,
		Values:           values,
		ChartPathOptions: chartPathBasicOptions,
	})
	return rel, err
}
