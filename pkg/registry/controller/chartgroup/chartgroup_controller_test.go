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

package chartgroup

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	core "k8s.io/client-go/testing"
	businessv1 "tkestack.io/tke/api/business/v1"
	"tkestack.io/tke/api/client/clientset/versioned/fake"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	registryv1 "tkestack.io/tke/api/registry/v1"
	v1 "tkestack.io/tke/api/registry/v1"
)

func newChartGroup(name, tenantID string, uid types.UID, visibility v1.Visibility) *v1.ChartGroup {
	return &v1.ChartGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:     name,
			UID:      uid,
			SelfLink: "/apis/registry.tkestack.io/v1/chartgroups/" + name,
		},
		Spec: v1.ChartGroupSpec{
			TenantID:    tenantID,
			Name:        name,
			DisplayName: name,
			Visibility:  visibility,
			Projects:    []string{},
		},
		Status: v1.ChartGroupStatus{
			Phase: registryv1.ChartGroupAvailable,
		},
	}
}

func defaultPublicChartGroup() *v1.ChartGroup {
	return newChartGroup("public", "default", types.UID("456"), v1.VisibilityPublic)
}

func newController() (*Controller, *fake.Clientset, *fake.Clientset) {
	client := fake.NewSimpleClientset()
	businessClient := fake.NewSimpleClientset()

	informerFactory := versionedinformers.NewSharedInformerFactory(client, 10*time.Minute)
	informer := informerFactory.Registry().V1().ChartGroups()

	controller := NewController(businessClient.BusinessV1(), client, informer, 30*time.Second, v1.ChartGroupFinalize)
	client.ClearActions()         // ignore any client calls made in init()
	businessClient.ClearActions() // ignore any client calls made in init()
	return controller, client, businessClient
}

func TestProcessUpdateChartGroup(t *testing.T) {
	var controller *Controller
	var client *fake.Clientset
	var businessClient *fake.Clientset
	oldProjects := []string{"old-project"}
	newProjects := []string{"new-project"}

	testCases := []struct {
		testName   string
		key        string
		updateFn   func(*v1.ChartGroup) (chartGroup *v1.ChartGroup) //Manipulate the structure
		cg         *v1.ChartGroup
		expectedFn func(*v1.ChartGroup, error) error //Error comparison function
	}{
		{
			testName: "If updating a valid chartGroup",
			key:      "valid-key",
			cg:       defaultPublicChartGroup(),
			updateFn: func(cg *v1.ChartGroup) (chartGroup *v1.ChartGroup) {
				controller.cache.getOrCreate("valid-key")
				return cg
			},
			expectedFn: func(cg *v1.ChartGroup, err error) error {
				return err
			},
		},
		{
			testName: "If sync projects",
			key:      "sync-projects",
			cg:       newChartGroup("sync-projects", "default", types.UID("sync-projects-uid"), v1.VisibilityPublic),
			updateFn: func(cg *v1.ChartGroup) (chartGroup *v1.ChartGroup) {
				keyExpected := cg.GetObjectMeta().GetName()
				controller.enqueue(cg)
				cachedChartGroup := controller.cache.getOrCreate(keyExpected)
				cachedChartGroup.state = cg
				controller.cache.set(keyExpected, cachedChartGroup)

				keyGot, quit := controller.queue.Get()
				if quit {
					t.Fatalf("get no queue element")
				}
				if keyExpected != keyGot.(string) {
					t.Fatalf("get chartGroup key error, expected: %s, got: %s", keyExpected, keyGot.(string))
				}

				newChartGroup := cg.DeepCopy()
				newChartGroup.Spec.Projects = newProjects
				return newChartGroup
			},
			expectedFn: func(cg *v1.ChartGroup, err error) error {
				if err != nil {
					return err
				}

				keyExpected := cg.GetObjectMeta().GetName()
				cachedChartGroupGot, exist := controller.cache.get(keyExpected)
				if !exist {
					return fmt.Errorf("update chartGroup error, cache should contain chartGroup: %s", keyExpected)
				}
				if !reflect.DeepEqual(cachedChartGroupGot.state.Spec.Projects, newProjects) {
					return fmt.Errorf("update Projects error, expected: %s, got: %s", newProjects, cachedChartGroupGot.state.Spec.Projects)
				}
				return nil
			},
		},
		{
			testName: "If sync business chartGroup",
			key:      "sync-business-chartgroup",
			cg:       newChartGroup("sync-business-chartgroup", "default", types.UID("sync-business-chartgroup-uid"), v1.VisibilityPublic),
			updateFn: func(cg *v1.ChartGroup) (chartGroup *v1.ChartGroup) {
				if _, err := businessClient.BusinessV1().ChartGroups(oldProjects[0]).Create(context.Background(),
					&businessv1.ChartGroup{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "sync-business-chartgroup",
							Namespace: oldProjects[0],
						},
						Spec: businessv1.ChartGroupSpec{
							Name:     "sync-business-chartgroup",
							TenantID: "default",
						},
					}, metav1.CreateOptions{},
				); err != nil {
					t.Fatalf("Failed to prepare business chartGroup %s for testing: %v", "sync-business-chartgroup", err)
				}

				cg.Spec.Projects = oldProjects

				keyExpected := cg.GetObjectMeta().GetName()
				cachedChartGroup := controller.cache.getOrCreate(keyExpected)

				err := controller.processUpdate(context.Background(), cachedChartGroup, cg, keyExpected)
				if err != nil {
					t.Errorf("update chartGroup error, %v", err)
				}
				newChartGroup := cg.DeepCopy()
				newChartGroup.Spec.Projects = newProjects
				return newChartGroup
			},
			expectedFn: func(cg *v1.ChartGroup, err error) error {
				if err != nil {
					return err
				}

				actions := businessClient.Actions()
				if len(actions) != 5 {
					t.Errorf("Expected 5 business client action, but got %v", len(actions))
				}
				if !actions[0].Matches("create", "chartgroups") || actions[0].GetSubresource() != "" {
					t.Errorf("Expected business-chartgroup action %v", actions[0])
				}
				if !actions[1].Matches("get", "chartgroups") || actions[1].GetSubresource() != "" {
					t.Errorf("Expected business-chartgroup action %v", actions[1])
				}
				if !actions[2].Matches("get", "chartgroups") || actions[2].GetSubresource() != "" {
					t.Errorf("Expected business-chartgroup action %v", actions[2])
				}
				if !actions[3].Matches("create", "chartgroups") || actions[3].GetSubresource() != "" {
					t.Errorf("Expected business-chartgroup action %v", actions[3])
				}
				if !actions[4].Matches("delete", "chartgroups") || actions[4].GetSubresource() != "" {
					t.Errorf("Expected business-chartgroup action %v", actions[4])
				}
				businessChartGroup := actions[3].(core.CreateAction).GetObject().(*businessv1.ChartGroup)
				if businessChartGroup == nil {
					t.Errorf("There should be a business chartGroup remaining")
				} else if businessChartGroup.Namespace != newProjects[0] || businessChartGroup.Name != cg.Spec.Name {
					return fmt.Errorf("create business chartGroup error, expected: %s/%s, got: %s/%s", newProjects[0], cg.Spec.Name, businessChartGroup.Namespace, businessChartGroup.Name)
				}
				return nil
			},
		},
	}

	for _, tc := range testCases {
		controller, client, businessClient = newController()
		if _, err := client.RegistryV1().ChartGroups().Create(context.Background(), tc.cg, metav1.CreateOptions{}); err != nil {
			t.Fatalf("Failed to prepare chartGroup %s for testing: %v", tc.key, err)
		}
		newCg := tc.updateFn(tc.cg)
		cachedChartGroup, exist := controller.cache.get(tc.key)
		if !exist {
			t.Fatalf("update chartGroup error, cache should contain chartGroup: %s", tc.key)
		}
		obtErr := controller.processUpdate(context.Background(), cachedChartGroup, newCg, tc.key)
		if err := tc.expectedFn(newCg, obtErr); err != nil {
			t.Errorf("%v processUpdate() %v", tc.testName, err)
		}
	}

}

// TestProcessCreateOrUpdateK8sError tests processUpdate
// with various kubernetes errors when patching status.
func TestProcessCreateOrUpdateK8sError(t *testing.T) {
	cgName := "cg-k8s-err"
	conflictErr := apierrors.NewConflict(schema.GroupResource{}, cgName, errors.New("object conflict"))
	notFoundErr := apierrors.NewNotFound(schema.GroupResource{}, cgName)

	testCases := []struct {
		desc      string
		k8sErr    error
		expectErr error
	}{
		{
			desc:      "conflict error",
			k8sErr:    conflictErr,
			expectErr: fmt.Errorf("not persisting update to chartGroup 'cg-k8s-err' that has been changed since we received it: %v", conflictErr),
		},
		{
			desc:      "not found error",
			k8sErr:    notFoundErr,
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cg := newChartGroup(cgName, "default", types.UID("123"), v1.VisibilityPublic)
			// make sure we will go into the update logic
			cg.Status.Phase = v1.ChartGroupPending
			// Preset finalizer so k8s error only happens when patching status.
			cg.Finalizers = []string{"chartgroup"}
			controller, client, _ := newController()
			client.PrependReactor("update", "chartgroups", func(action core.Action) (bool, runtime.Object, error) {
				return true, nil, tc.k8sErr
			})

			cachedChartGroup := controller.cache.getOrCreate(cg.Name)
			cachedChartGroup.state = cg
			if err := controller.processUpdate(context.Background(), cachedChartGroup, cg, cgName); !reflect.DeepEqual(err, tc.expectErr) {
				t.Fatalf("processUpdate() = %v, want %v", err, tc.expectErr)
			}
			if tc.expectErr == nil {
				return
			}
		})
	}

}

func TestSyncChartGroup(t *testing.T) {
	var controller *Controller

	testCases := []struct {
		testName   string
		key        string
		updateFn   func()            //Function to manipulate the controller element to simulate error
		expectedFn func(error) error //Expected function if returns nil then test passed, failed otherwise
	}{
		{
			testName: "if an invalid chartGroup name is synced",
			key:      "invalid/key/string",
			updateFn: func() {
				controller, _, _ = newController()
			},
			expectedFn: func(e error) error {
				//TODO: should find a way to test for dependent package errors in such a way that it won't break
				//TODO:	our tests, currently we only test if there is an error.
				//Error should be unexpected key format: "invalid/key/string"
				expectedError := fmt.Sprintf("unexpected key format: %q", "invalid/key/string")
				if e == nil || e.Error() != expectedError {
					return fmt.Errorf("Expected=unexpected key format: %q, Obtained=%v", "invalid/key/string", e)
				}
				return nil
			},
		},
		//TODO: see if we can add a test for valid but error throwing chartGroup, its difficult right now because TestSyncChartGroup() currently runtime.HandleError
		{
			testName: "if valid chartGroup",
			key:      "valid-chartgroup",
			updateFn: func() {
				testCg := defaultPublicChartGroup()
				controller, _, _ = newController()
				controller.enqueue(testCg)
				cg := controller.cache.getOrCreate("valid-chartgroup")
				cg.state = testCg
			},
			expectedFn: func(e error) error {
				//error should be nil
				if e != nil {
					return fmt.Errorf("Expected=nil, Obtained=%v", e)
				}
				return nil
			},
		},
	}

	for _, tc := range testCases {
		tc.updateFn()
		obtainedErr := controller.syncItem(tc.key)

		//expected matches obtained ??.
		if exp := tc.expectedFn(obtainedErr); exp != nil {
			t.Errorf("%v Error:%v", tc.testName, exp)
		}

		//Post processing, the element should not be in the sync queue.
		_, exist := controller.cache.get(tc.key)
		if exist {
			t.Fatalf("%v working Queue should be empty, but contains %s", tc.testName, tc.key)
		}
	}
}

func TestProcessChartGroupDeletion(t *testing.T) {
	var controller *Controller
	// Add a global cgKey name
	cgKey := "test-cg"

	testCases := []struct {
		testName   string
		updateFn   func(*Controller)       // Update function used to manipulate srv and controller values
		expectedFn func(cgErr error) error // Function to check if the returned value is expected
	}{
		{
			testName: "If a non-existent chartGroup is deleted",
			updateFn: func(controller *Controller) {
				// Does not do anything
			},
			expectedFn: func(cgErr error) error {
				return cgErr
			},
		},
		{
			testName: "If delete was successful",
			updateFn: func(controller *Controller) {
				testCg := defaultPublicChartGroup()
				controller.enqueue(testCg)
				cg := controller.cache.getOrCreate(cgKey)
				cg.state = testCg
				controller.cache.set(cgKey, cg)
			},
			expectedFn: func(cgErr error) error {
				if cgErr != nil {
					return fmt.Errorf("Expected=nil Obtained=%v", cgErr)
				}

				// It should no longer be in the workqueue.
				_, exist := controller.cache.get(cgKey)
				if exist {
					return fmt.Errorf("delete chartGroup error, queue should not contain chartGroup: %s any more", cgKey)
				}

				return nil
			},
		},
	}

	for _, tc := range testCases {
		//Create a new controller.
		controller, _, _ = newController()
		tc.updateFn(controller)
		obtainedErr := controller.processDeletion(cgKey)
		if err := tc.expectedFn(obtainedErr); err != nil {
			t.Errorf("%v processChartGroupDeletion() %v", tc.testName, err)
		}
	}

}

func TestNeedsUpdate(t *testing.T) {
	var oldCg, newCg *v1.ChartGroup

	testCases := []struct {
		testName            string //Name of the test case
		updateFn            func() //Function to update the chartGroup object
		expectedNeedsUpdate bool   //needsupdate always returns bool
	}{
		{
			testName: "If the chartGroup Visibility is changed from public to private",
			updateFn: func() {
				oldCg = defaultPublicChartGroup()
				newCg = defaultPublicChartGroup()
				newCg.Spec.Visibility = v1.VisibilityPrivate
			},
			expectedNeedsUpdate: true,
		},
		{
			testName: "If the Projects are different",
			updateFn: func() {
				oldCg = defaultPublicChartGroup()
				newCg = defaultPublicChartGroup()
				oldCg.Spec.Projects = []string{"p1"}
				newCg.Spec.Projects = []string{"p1", "p2"}
			},
			expectedNeedsUpdate: true,
		},
		{
			testName: "If DisplayName are different",
			updateFn: func() {
				oldCg = defaultPublicChartGroup()
				newCg = defaultPublicChartGroup()
				oldCg.Spec.DisplayName = "name1"
				newCg.Spec.DisplayName = "name2"
			},
			expectedNeedsUpdate: true,
		},
		{
			testName: "If TenantID are different",
			updateFn: func() {
				oldCg = defaultPublicChartGroup()
				newCg = defaultPublicChartGroup()
				oldCg.Spec.TenantID = "t1"
				newCg.Spec.TenantID = "t2"
			},
			expectedNeedsUpdate: true,
		},
		{
			testName: "If UID is different",
			updateFn: func() {
				oldCg = defaultPublicChartGroup()
				newCg = defaultPublicChartGroup()
				oldCg.UID = types.UID("UID old")
				newCg.UID = types.UID("UID new")
			},
			expectedNeedsUpdate: true,
		},
		{
			testName: "If Type is different",
			updateFn: func() {
				oldCg = defaultPublicChartGroup()
				newCg = defaultPublicChartGroup()
				oldCg.Spec.Type = v1.RepoTypeProject
				newCg.Spec.Type = v1.RepoTypePersonal
			},
			expectedNeedsUpdate: true,
		},
		{
			testName: "If ChartCount is different",
			updateFn: func() {
				oldCg = defaultPublicChartGroup()
				newCg = defaultPublicChartGroup()
				oldCg.Status.ChartCount = 1
				newCg.Status.ChartCount = 2
			},
			expectedNeedsUpdate: true,
		},
		{
			testName: "If Phase is different 1",
			updateFn: func() {
				oldCg = defaultPublicChartGroup()
				newCg = defaultPublicChartGroup()
				oldCg.Status.Phase = v1.ChartGroupPending
				newCg.Status.Phase = v1.ChartGroupTerminating
			},
			expectedNeedsUpdate: true,
		},
	}

	controller, _, _ := newController()
	for _, tc := range testCases {
		tc.updateFn()
		obtainedResult := controller.needsUpdate(oldCg, newCg)
		if obtainedResult != tc.expectedNeedsUpdate {
			t.Errorf("%v needsUpdate() should have returned %v but returned %v", tc.testName, tc.expectedNeedsUpdate, obtainedResult)
		}
	}
}

//All the test cases for ChartGroupCache uses a single cache, these below test cases should be run in order,
//as tc1 (addCache would add elements to the cache)
//and tc2 (delCache would remove element from the cache without it adding automatically)
//Please keep this in mind while adding new test cases.
func TestChartGroupCache(t *testing.T) {
	//ChartGroupCache a common chartGroup cache for all the test cases
	sc := &chartGroupCache{m: make(map[string]*cachedChartGroup)}

	testCases := []struct {
		testName     string
		setCacheFn   func()
		checkCacheFn func() error
	}{
		{
			testName: "Add",
			setCacheFn: func() {
				cS := sc.getOrCreate("addTest")
				cS.state = defaultPublicChartGroup()
			},
			checkCacheFn: func() error {
				//There must be exactly one element
				if len(sc.m) != 1 {
					return fmt.Errorf("Expected=1 Obtained=%d", len(sc.m))
				}
				return nil
			},
		},
		{
			testName: "Del",
			setCacheFn: func() {
				sc.delete("addTest")

			},
			checkCacheFn: func() error {
				//Now it should have no element
				if len(sc.m) != 0 {
					return fmt.Errorf("Expected=0 Obtained=%d", len(sc.m))
				}
				return nil
			},
		},
		{
			testName: "Set and Get",
			setCacheFn: func() {
				sc.set("addTest", &cachedChartGroup{state: defaultPublicChartGroup()})
			},
			checkCacheFn: func() error {
				//Now it should have one element
				Cs, bool := sc.get("addTest")
				if !bool {
					return fmt.Errorf("is Available Expected=true Obtained=%v", bool)
				}
				if Cs == nil {
					return fmt.Errorf("cachedChartGroup expected:non-nil Obtained=nil")
				}
				return nil
			},
		},
		{
			testName: "ListKeys",
			setCacheFn: func() {
				//Add one more entry here
				sc.set("addTest1", &cachedChartGroup{state: defaultPublicChartGroup()})
			},
			checkCacheFn: func() error {
				//It should have two elements
				keys := sc.listKeys()
				if len(keys) != 2 {
					return fmt.Errorf("elements Expected=2 Obtained=%v", len(keys))
				}
				return nil
			},
		},
		{
			testName:   "GetbyKeys",
			setCacheFn: nil, //Nothing to set
			checkCacheFn: func() error {
				//It should have two elements
				cg, exist := sc.get("addTest")
				if cg == nil || exist == false {
					return fmt.Errorf("Expected(non-nil, true) Obtained(%v,%v)", cg, exist)
				}
				return nil
			},
		},
		{
			testName:   "allChartGroups",
			setCacheFn: nil, //Nothing to set
			checkCacheFn: func() error {
				//It should return two elements
				svcArray := sc.allChartGroups()
				if len(svcArray) != 2 {
					return fmt.Errorf("Expected(2) Obtained(%v)", len(svcArray))
				}
				return nil
			},
		},
	}

	for _, tc := range testCases {
		if tc.setCacheFn != nil {
			tc.setCacheFn()
		}
		if err := tc.checkCacheFn(); err != nil {
			t.Errorf("%v returned %v", tc.testName, err)
		}
	}
}

// TODO(@MrHohn): Verify the end state when below issue is resolved:
// https://github.com/kubernetes/client-go/issues/607
func TestUpdateStatus(t *testing.T) {
	testCases := []struct {
		desc         string
		cg           *v1.ChartGroup
		newStatus    *v1.ChartGroupStatus
		expectUpdate bool
	}{
		{
			desc: "no-op add status",
			cg: &v1.ChartGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-patch-status",
				},
				Status: v1.ChartGroupStatus{
					ChartCount: 1,
					Phase:      v1.ChartGroupAvailable,
				},
			},
			newStatus: &v1.ChartGroupStatus{
				ChartCount: 1,
				Phase:      v1.ChartGroupAvailable,
			},
			expectUpdate: false,
		},
		{
			desc: "add status",
			cg: &v1.ChartGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-patch-status",
				},
				Status: v1.ChartGroupStatus{
					Phase: v1.ChartGroupAvailable,
				},
			},
			newStatus: &v1.ChartGroupStatus{
				ChartCount: 1,
				Phase:      v1.ChartGroupAvailable,
			},
			expectUpdate: true,
		},
		{
			desc: "no-op clear status",
			cg: &v1.ChartGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-patch-status",
				},
				Status: v1.ChartGroupStatus{},
			},
			newStatus:    &v1.ChartGroupStatus{},
			expectUpdate: false,
		},
		{
			desc: "clear status",
			cg: &v1.ChartGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-patch-status",
				},
				Status: v1.ChartGroupStatus{
					ChartCount: 1,
					Phase:      v1.ChartGroupAvailable,
				},
			},
			newStatus:    &v1.ChartGroupStatus{},
			expectUpdate: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			c := fake.NewSimpleClientset()
			s := &Controller{
				client: c,
			}
			if _, err := s.client.RegistryV1().ChartGroups().Create(context.Background(), tc.cg, metav1.CreateOptions{}); err != nil {
				t.Fatalf("Failed to prepare chartGroup for testing: %v", err)
			}
			if _, err := s.updateStatus(context.Background(), tc.cg, &tc.cg.Status, tc.newStatus); err != nil {
				t.Fatalf("updateStatus() = %v, want nil", err)
			}
			updateActionFound := false
			for _, action := range c.Actions() {
				if action.Matches("update", "chartgroups") {
					updateActionFound = true
				}
			}
			if updateActionFound != tc.expectUpdate {
				t.Errorf("Got updateActionFound = %t, want %t", updateActionFound, tc.expectUpdate)
			}
		})
	}
}
