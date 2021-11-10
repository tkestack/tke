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
	"strings"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/registry/rest"
	authv1 "tkestack.io/tke/api/auth/v1"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	authversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	"tkestack.io/tke/api/registry"
	registryapi "tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/apiserver/authentication"
	authutil "tkestack.io/tke/pkg/auth/util"
	harbor "tkestack.io/tke/pkg/registry/harbor/client"
	harborHandler "tkestack.io/tke/pkg/registry/harbor/handler"
	helm "tkestack.io/tke/pkg/registry/harbor/helmClient"
	chartgroupstrategy "tkestack.io/tke/pkg/registry/registry/chartgroup"
	genericutil "tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
)

// REST adapts a service registry into apiserver's RESTStorage model.
type REST struct {
	chartgroup     ChartGroupStorage
	registryClient *registryinternalclient.RegistryClient
	authClient     authversionedclient.AuthV1Interface
	harborClient   *harbor.APIClient
	helmClient     *helm.APIClient
}

type ChartGroupStorage interface {
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
	chartgroup ChartGroupStorage,
	registryClient *registryinternalclient.RegistryClient,
	authClient authversionedclient.AuthV1Interface,
	harborClient *harbor.APIClient,
	helmClient *helm.APIClient,
) *REST {
	rest := &REST{
		chartgroup:     chartgroup,
		registryClient: registryClient,
		authClient:     authClient,
		harborClient:   harborClient,
		helmClient:     helmClient,
	}
	return rest
}

var (
	_ ChartGroupStorage           = &REST{}
	_ rest.ShortNamesProvider     = &REST{}
	_ rest.StorageVersionProvider = &REST{}
)

func (rs *REST) StorageVersion() runtime.GroupVersioner {
	return rs.chartgroup.StorageVersion()
}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (rs *REST) ShortNames() []string {
	return []string{"rcg"}
}

func (rs *REST) NamespaceScoped() bool {
	return rs.chartgroup.NamespaceScoped()
}

func (rs *REST) New() runtime.Object {
	return rs.chartgroup.New()
}

func (rs *REST) NewList() runtime.Object {
	return rs.chartgroup.NewList()
}

func (rs *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return rs.chartgroup.Get(ctx, name, options)
}

func (rs *REST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	return rs.chartgroup.List(ctx, options)
}

func (rs *REST) Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	return rs.chartgroup.Watch(ctx, options)
}

func (rs *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	if rs.harborClient != nil {
		o := obj.(*registryapi.ChartGroup)
		_, tenantID := authentication.UsernameAndTenantID(ctx)

		err := harborHandler.CreateProject(
			ctx,
			rs.harborClient,
			fmt.Sprintf("%s-chart-%s", tenantID, o.Spec.Name),
			o.Spec.Visibility == registryapi.VisibilityPublic,
		)
		if err != nil {
			return nil, err
		}
	}
	obj, err := rs.chartgroup.Create(ctx, obj, createValidation, options)
	if err != nil {
		if rs.harborClient != nil {
			o := obj.(*registryapi.ChartGroup)
			// cleanup harbor project
			harborHandler.DeleteProject(ctx, rs.harborClient, nil, fmt.Sprintf("%s-chart-%s", o.Spec.TenantID, o.Spec.Name))
		}
		return nil, err
	}
	cg := obj.(*registry.ChartGroup)
	// update policy binding
	err = rs.createOrUpdatePolicyBinding(ctx, cg)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (rs *REST) Delete(ctx context.Context, id string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	// TODO: handle graceful
	obj, _, err := rs.chartgroup.Delete(ctx, id, deleteValidation, options)
	if err != nil {
		return nil, false, err
	}
	if rs.harborClient != nil {
		o := obj.(*registryapi.ChartGroup)
		err := harborHandler.DeleteProject(ctx, rs.harborClient, rs.helmClient, fmt.Sprintf("%s-chart-%s", o.Spec.TenantID, o.Spec.Name))
		if err != nil {
			return nil, false, err
		}
	}
	// delete policy binding
	cg := obj.(*registry.ChartGroup)
	err = rs.deletePolicyBinding(ctx, cg)
	if err != nil {
		return nil, false, err
	}

	return obj, true, nil
}

func (rs *REST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	oldObj, err := rs.chartgroup.Get(ctx, name, &metav1.GetOptions{})
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
	oldChartGroup := oldObj.(*registry.ChartGroup)
	obj, err := objInfo.UpdatedObject(ctx, oldChartGroup)
	if err != nil {
		return nil, false, err
	}
	cg := obj.(*registry.ChartGroup)

	// update policy binding
	err = rs.createOrUpdatePolicyBinding(ctx, cg)
	if err != nil {
		return nil, false, err
	}

	// Copy over non-user fields
	strategy := chartgroupstrategy.NewStrategy(rs.registryClient)
	if err := rest.BeforeUpdate(strategy, ctx, cg, oldChartGroup); err != nil {
		return nil, false, err
	}

	return rs.chartgroup.Update(ctx, name, rest.DefaultUpdatedObjectInfo(cg), createValidation, updateValidation, forceAllowCreate, options)
}

func (rs *REST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return rs.chartgroup.ConvertToTable(ctx, object, tableOptions)
}

// createOrUpdatePolicyBinding add policy binding
func (rs *REST) createOrUpdatePolicyBinding(ctx context.Context, cg *registry.ChartGroup) error {
	username := cg.Spec.Creator
	isValidUsername := validUsername(username)
	pb := make([]*authv1.CustomPolicyBinding, 0)
	defaultAll := []authv1.Subject{{ID: authutil.DefaultAll, Name: authutil.DefaultAll}}

	switch cg.Spec.Visibility {
	case registry.VisibilityUser:
		{
			domain := authutil.DefaultDomain

			usernames := cg.Spec.Users[:]
			if isValidUsername && !genericutil.InStringSlice(usernames, username) {
				usernames = append(usernames, username)
			}
			users := make([]authv1.Subject, len(usernames))
			for k, v := range usernames {
				users[k] = authv1.Subject{ID: v, Name: v}
			}

			// owner policy
			policyID := authutil.ChartGroupFullPolicyID(cg.Spec.TenantID)
			pb = append(pb, buildCustomBinding(cg, policyID, policyID, domain, authutil.ChartGroupPolicyResources(cg.Name), users))

			policyID = authutil.ChartFullPolicyID(cg.Spec.TenantID)
			pb = append(pb, buildCustomBinding(cg, policyID, policyID, domain, authutil.ChartPolicyResources(cg.Name), users))
			break
		}
	case registry.VisibilityProject:
		{
			break
		}
	case registry.VisibilityPublic:
		{
			domain := authutil.DefaultDomain

			if isValidUsername {
				users := []authv1.Subject{{ID: username, Name: username}}

				// owner policy
				policyID := authutil.ChartGroupFullPolicyID(cg.Spec.TenantID)
				pb = append(pb, buildCustomBinding(cg, policyID, policyID, domain, authutil.ChartGroupPolicyResources(cg.Name), users))

				policyID = authutil.ChartFullPolicyID(cg.Spec.TenantID)
				pb = append(pb, buildCustomBinding(cg, policyID, policyID, domain, authutil.ChartPolicyResources(cg.Name), users))
			}

			// others policy
			policyID := authutil.ChartGroupPullPolicyID(cg.Spec.TenantID)
			pb = append(pb, buildCustomBinding(cg, policyID, policyID, domain, authutil.ChartGroupPolicyResources(cg.Name), defaultAll))

			policyID = authutil.ChartPullPolicyID(cg.Spec.TenantID)
			pb = append(pb, buildCustomBinding(cg, policyID, policyID, domain, authutil.ChartPolicyResources(cg.Name), defaultAll))
			break
		}
	}

	list, err := rs.authClient.CustomPolicyBindings(cg.Name).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, p := range pb {
		exist := filterCustomPolicyBinding(p.Name, list.Items)
		if exist != nil {
			log.Debugf("update custombinding %v/%v", cg.Name, p.Name)
			exist.Spec.LastDomain = exist.Spec.Domain
			exist.Spec.Domain = p.Spec.Domain
			exist.Spec.PolicyID = p.Spec.PolicyID
			exist.Spec.Users = append(exist.Spec.Users, p.Spec.Users...)
			exist.Spec.Groups = append(exist.Spec.Groups, p.Spec.Groups...)
			_, err = rs.authClient.CustomPolicyBindings(cg.Name).Update(ctx, exist, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		} else {
			log.Debugf("add custombinding %v/%v", cg.Name, p.Name)
			_, err := rs.authClient.CustomPolicyBindings(cg.Name).Create(ctx, p, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		}
	}
	oldIds := make([]string, len(list.Items))
	newIds := make([]string, len(pb))
	for k, v := range list.Items {
		oldIds[k] = v.Name
	}
	for k, v := range pb {
		newIds[k] = v.Name
	}
	_, removed := genericutil.DiffStringSlice(oldIds, newIds)
	for _, name := range removed {
		exist := filterCustomPolicyBinding(name, list.Items)
		if exist != nil {
			log.Debugf("delete custombinding %v/%v", cg.Name, name)
			err = rs.authClient.CustomPolicyBindings(cg.Name).Delete(ctx, name, metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func buildCustomBinding(cg *registry.ChartGroup, name, policyID, domain string, resources []string, users []authv1.Subject) *authv1.CustomPolicyBinding {
	return &authv1.CustomPolicyBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cg.Name,
			Name:      policyID,
		},
		Spec: authv1.CustomPolicyBindingSpec{
			TenantID:   cg.Spec.TenantID,
			Domain:     domain,
			LastDomain: "",
			PolicyID:   policyID,
			Resources:  resources,
			RulePrefix: cg.Name,
			Users:      users,
		},
	}
}

func filterCustomPolicyBinding(name string, list []authv1.CustomPolicyBinding) *authv1.CustomPolicyBinding {
	if len(list) == 0 {
		return nil
	}
	for _, v := range list {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func validUsername(name string) bool {
	return strings.TrimSpace(name) != ""
}

// deletePolicyBinding remove policy binding
func (rs *REST) deletePolicyBinding(ctx context.Context, cg *registry.ChartGroup) error {
	list, err := rs.authClient.CustomPolicyBindings(cg.Name).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, p := range list.Items {
		err = rs.authClient.CustomPolicyBindings(cg.Name).Delete(ctx, p.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
