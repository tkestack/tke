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

package deletion

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	core "k8s.io/client-go/testing"
	businessv1 "tkestack.io/tke/api/business/v1"
	"tkestack.io/tke/api/client/clientset/versioned/fake"
	v1 "tkestack.io/tke/api/registry/v1"
)

func TestFinalized(t *testing.T) {
	testChartGroup := &v1.ChartGroup{
		Spec: v1.ChartGroupSpec{
			Finalizers: []v1.FinalizerName{"a", "b"},
		},
	}
	if finalized(testChartGroup) {
		t.Errorf("Unexpected result, namespace is not finalized")
	}
	testChartGroup.Spec.Finalizers = []v1.FinalizerName{}
	if !finalized(testChartGroup) {
		t.Errorf("Expected object to be finalized")
	}
}

func TestFinalizeChartGroupFunc(t *testing.T) {
	registryClient := &fake.Clientset{}
	testChartGroup := &v1.ChartGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "test",
			ResourceVersion: "1",
		},
		Spec: v1.ChartGroupSpec{
			Finalizers: []v1.FinalizerName{"chartgroup", "other"},
		},
	}
	d := chartGroupResourcesDeleter{
		registryClient: registryClient.RegistryV1(),
		finalizerToken: v1.ChartGroupFinalize,
	}
	d.finalizeChartGroup(context.Background(), testChartGroup)
	actions := registryClient.Actions()
	if len(actions) != 1 {
		t.Errorf("Expected 1 mock client action, but got %v", len(actions))
	}
	if !actions[0].Matches("update", "chartgroups") || actions[0].GetSubresource() != "" {
		t.Errorf("Expected finalize-chartgroup action %v", actions[0])
	}
	finalizers := actions[0].(core.UpdateAction).GetObject().(*v1.ChartGroup).Spec.Finalizers
	if len(finalizers) != 1 {
		t.Errorf("There should be a single finalizer remaining")
	}
	if string(finalizers[0]) != "other" {
		t.Errorf("Unexpected finalizer value, %v", finalizers[0])
	}
}

func testSyncChartGroupThatIsTerminating(t *testing.T, versions *metav1.APIVersions) {
	now := metav1.Now()
	chartGroupName := "rcg-test"
	chartGroupSpecName := "test"
	projectID := "project"
	tenantID := "t"
	testChartGroupPendingFinalize := &v1.ChartGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:              chartGroupName,
			ResourceVersion:   "1",
			DeletionTimestamp: &now,
		},
		Spec: v1.ChartGroupSpec{
			TenantID:   tenantID,
			Name:       chartGroupSpecName,
			Projects:   []string{projectID},
			Finalizers: []v1.FinalizerName{"chartgroup"},
		},
		Status: v1.ChartGroupStatus{
			Phase: v1.ChartGroupTerminating,
		},
	}
	testChartGroupFinalizeComplete := &v1.ChartGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:              chartGroupName,
			ResourceVersion:   tenantID,
			DeletionTimestamp: &now,
		},
		Spec: v1.ChartGroupSpec{},
		Status: v1.ChartGroupStatus{
			Phase: v1.ChartGroupTerminating,
		},
	}
	testChart := &v1.Chart{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testc",
			Namespace: chartGroupName,
		},
		Spec: v1.ChartSpec{
			ChartGroupName: chartGroupSpecName,
			TenantID:       tenantID,
		},
	}
	testBusinessChartGroup := &businessv1.ChartGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      chartGroupSpecName,
			Namespace: projectID,
		},
		Spec: businessv1.ChartGroupSpec{
			Name:     chartGroupSpecName,
			TenantID: tenantID,
		},
	}

	scenarios := map[string]struct {
		testChartGroup          *v1.ChartGroup
		registryClientActionSet sets.String
		businessClientActionSet sets.String
	}{
		"pending-finalize": {
			testChartGroup: testChartGroupPendingFinalize,
			registryClientActionSet: sets.NewString(
				strings.Join([]string{"registry.tkestack.io", "v1", "get", "chartgroups", ""}, "-"),
				strings.Join([]string{"registry.tkestack.io", "v1", "delete", "chartgroups", ""}, "-"),
				strings.Join([]string{"registry.tkestack.io", "v1", "update", "chartgroups", ""}, "-"),
				strings.Join([]string{"registry.tkestack.io", "v1", "list", "charts", ""}, "-"),
				strings.Join([]string{"registry.tkestack.io", "v1", "delete", "charts", ""}, "-"),
			),
			businessClientActionSet: sets.NewString(
				strings.Join([]string{"business.tkestack.io", "v1", "get", "chartgroups", ""}, "-"),
				strings.Join([]string{"business.tkestack.io", "v1", "delete", "chartgroups", ""}, "-"),
			),
		},
		"complete-finalize": {
			testChartGroup: testChartGroupFinalizeComplete,
			registryClientActionSet: sets.NewString(
				strings.Join([]string{"registry.tkestack.io", "v1", "get", "chartgroups", ""}, "-"),
				strings.Join([]string{"registry.tkestack.io", "v1", "delete", "chartgroups", ""}, "-"),
			),
			businessClientActionSet: sets.NewString(),
		},
	}

	for scenario, testInput := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			registryClient := fake.NewSimpleClientset(testInput.testChartGroup, testChart)
			businessClient := fake.NewSimpleClientset(testBusinessChartGroup)

			d := NewChartGroupResourcesDeleter(businessClient.BusinessV1(), registryClient.RegistryV1(), v1.ChartGroupFinalize, true, nil)
			if err := d.Delete(context.Background(), testInput.testChartGroup.Name); err != nil {
				t.Errorf("when syncing chartGroup, got %q", err)
			}

			// validate traffic from kube client
			actionSet := sets.NewString()
			for _, action := range registryClient.Actions() {
				actionSet.Insert(strings.Join([]string{action.GetResource().Group, action.GetResource().Version, action.GetVerb(), action.GetResource().Resource, action.GetSubresource()}, "-"))
			}
			if !actionSet.Equal(testInput.registryClientActionSet) {
				t.Errorf("mock client expected actions:\n%v\n but got:\n%v\nDifference:\n%v",
					testInput.registryClientActionSet, actionSet, testInput.registryClientActionSet.Difference(actionSet))
			}

			// validate traffic from business client
			businessActionSet := sets.NewString()
			for _, action := range businessClient.Actions() {
				businessActionSet.Insert(strings.Join([]string{action.GetResource().Group, action.GetResource().Version, action.GetVerb(), action.GetResource().Resource, action.GetSubresource()}, "-"))
			}
			if !businessActionSet.Equal(testInput.businessClientActionSet) {
				t.Errorf("mock client expected actions:\n%v\n but got:\n%v\nDifference:\n%v",
					testInput.businessClientActionSet, actionSet, testInput.businessClientActionSet.Difference(businessActionSet))
			}
		})
	}
}

func TestRetryOnConflictError(t *testing.T) {
	registryClient := &fake.Clientset{}
	numTries := 0
	retryOnce := func(ctx context.Context, cg *v1.ChartGroup) (*v1.ChartGroup, error) {
		numTries++
		if numTries <= 1 {
			return cg, errors.NewConflict(v1.Resource("chartgroups"), cg.Name, fmt.Errorf("ERROR"))
		}
		return cg, nil
	}
	cg := &v1.ChartGroup{}
	d := chartGroupResourcesDeleter{
		registryClient: registryClient.RegistryV1(),
	}
	_, err := d.retryOnConflictError(context.Background(), cg, retryOnce)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if numTries != 2 {
		t.Errorf("Expected %v, but got %v", 2, numTries)
	}
}

func TestSyncChartGroupThatIsTerminatingNonExperimental(t *testing.T) {
	testSyncChartGroupThatIsTerminating(t, &metav1.APIVersions{})
}

func TestSyncNamespaceThatIsTerminatingV1(t *testing.T) {
	testSyncChartGroupThatIsTerminating(t, &metav1.APIVersions{Versions: []string{"registry.tkestack.io/v1"}})
}

func TestSyncChartGroupThatIsAvailable(t *testing.T) {
	registryClient := &fake.Clientset{}
	testChartGroup := &v1.ChartGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "test",
			ResourceVersion: "1",
		},
		Spec: v1.ChartGroupSpec{
			Finalizers: []v1.FinalizerName{"chartgroup"},
		},
		Status: v1.ChartGroupStatus{
			Phase: v1.ChartGroupAvailable,
		},
	}
	d := NewChartGroupResourcesDeleter(nil, registryClient.RegistryV1(), v1.ChartGroupFinalize, true, nil)
	err := d.Delete(context.Background(), testChartGroup.Name)
	if err != nil {
		t.Errorf("Unexpected error when synching namespace %v", err)
	}
	if len(registryClient.Actions()) != 1 {
		t.Errorf("Expected only one action from controller, but got: %d %v", len(registryClient.Actions()), registryClient.Actions())
	}
	action := registryClient.Actions()[0]
	if !action.Matches("get", "chartgroups") {
		t.Errorf("Expected get chartgroups, got: %v", action)
	}
}
