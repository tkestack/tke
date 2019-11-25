/*
 * Copyright 2019 THL A29 Limited, a Tencent company.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package namespace

import (
	"fmt"
	"strings"
	"testing"

	apimachineryresource "k8s.io/apimachinery/pkg/api/resource"
	apimachinerymetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/api/business"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/util/validation"
	"tkestack.io/tke/pkg/util/resource"
)

const (
	ClusterName   = "cluster-name"
	NamespaceName = "namespace-name"
	ProjectName   = "project-name"
	Phase         = "Active"
	TenantID      = "tenant-id"

	Namespace = "test-namespace"
)

var (
	_testCluster = platformv1.Cluster{
		ObjectMeta: apimachinerymetav1.ObjectMeta{
			Name: ClusterName,
		},
		Spec: platformv1.ClusterSpec{
			TenantID: TenantID,
		},
		Status: platformv1.ClusterStatus{},
	}
	_clusterMap = map[string]*platformv1.Cluster{
		ClusterName: &_testCluster,
	}

	_testProject = business.Project{
		ObjectMeta: apimachinerymetav1.ObjectMeta{
			Name: ProjectName,
		},
		Spec: business.ProjectSpec{
			TenantID: TenantID,
		},
		Status: business.ProjectStatus{},
	}
	_projectMap = map[string]*business.Project{
		ProjectName: &_testProject,
	}

	_testNamespace = business.Namespace{
		ObjectMeta: apimachinerymetav1.ObjectMeta{
			ClusterName:     ClusterName,
			Name:            NamespaceName,
			Namespace:       ProjectName,
			ResourceVersion: "v2",
		},
		Spec: business.NamespaceSpec{
			ClusterName: ClusterName,
			Namespace:   Namespace,
			TenantID:    TenantID,
		},
		Status: business.NamespaceStatus{
			Phase: Phase,
			Used:  business.ResourceList{},
		},
	}
	_oldNamespace = business.Namespace{
		ObjectMeta: apimachinerymetav1.ObjectMeta{
			ClusterName:     ClusterName,
			Name:            NamespaceName,
			Namespace:       ProjectName,
			ResourceVersion: "v1",
		},
		Spec: business.NamespaceSpec{
			ClusterName: ClusterName,
			Namespace:   Namespace,
			TenantID:    TenantID,
		},
		Status: business.NamespaceStatus{
			Phase: Phase,
			Used:  business.ResourceList{},
		},
	}
	_namespaceMap = map[string]*business.Namespace{
		NamespaceName: &_testNamespace,
	}
)

func TestClusterLimitation(t *testing.T) {
	_testNamespace.Spec.ClusterName = ClusterName

	hasError := false
	errors := ValidateNamespaceUpdate(&_testNamespace, &_oldNamespace, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		if strings.Contains(err.Error(), _clusterLimitErrorInfo) {
			hasError = true
		} else {
			t.Errorf("Unexpected: %s", err.Error())
		}
	}
	if !hasError {
		t.Errorf("Expect: %s", _clusterLimitErrorInfo)
	}

	_testProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {},
	}
	errors = ValidateNamespaceUpdate(&_testNamespace, &_oldNamespace, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		t.Errorf("Unexpected: %s", err.Error())
	}
}

func TestQuotaLimitation(t *testing.T) {
	quota, _ := apimachineryresource.ParseQuantity("1")
	_testProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}

	_testNamespace.Spec.Hard = business.ResourceList{}

	hasError := false
	errors := ValidateNamespaceUpdate(&_testNamespace, &_oldNamespace, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		if strings.Contains(err.Error(), resource.QuotaLimitErrorInfo) {
			hasError = true
		} else {
			t.Errorf("Unexpected: %s", err.Error())
		}
	}
	if !hasError {
		t.Errorf("Expect: %s", resource.QuotaLimitErrorInfo)
	}
}

func TestCreateAllocatable(t *testing.T) {
	quota, _ := apimachineryresource.ParseQuantity("1")
	_testProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}

	quota, _ = apimachineryresource.ParseQuantity("10")
	_testNamespace.Spec.Hard = business.ResourceList{
		"requests.cpu": quota,
	}

	hasError := false
	errors := ValidateNamespaceUpdate(&_testNamespace, &_oldNamespace, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		if strings.Contains(err.Error(), resource.AllocatableErrorInfo) {
			hasError = true
		} else {
			t.Errorf("Unexpected: %s", err.Error())
		}
	}
	if !hasError {
		t.Errorf("Expect: %s", resource.AllocatableErrorInfo)
	}

	quota, _ = apimachineryresource.ParseQuantity("1")
	_testNamespace.Spec.Hard["requests.cpu"] = quota
	errors = ValidateNamespaceUpdate(&_testNamespace, &_oldNamespace, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		t.Errorf("Unexpected: %s", err.Error())
	}
}

func TestUpdateAllocatable(t *testing.T) {
	/* remaining quota: 10 - 7 */
	quota, _ := apimachineryresource.ParseQuantity("10")
	_testProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}
	quota, _ = apimachineryresource.ParseQuantity("7")
	_testProject.Status.Clusters = business.ClusterUsed{
		ClusterName: {
			Used: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}

	/* request more quota: 7 - 3*/
	quota, _ = apimachineryresource.ParseQuantity("7")
	_testNamespace.Spec.Hard = business.ResourceList{
		"requests.cpu": quota,
	}
	_testProject.Status.CalculatedNamespaces = []string{_testNamespace.Name}
	quota, _ = apimachineryresource.ParseQuantity("3")
	_oldNamespace.Spec.Hard = business.ResourceList{
		"requests.cpu": quota,
	}

	hasError := false
	errors := ValidateNamespaceUpdate(&_testNamespace, &_oldNamespace, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		if strings.Contains(err.Error(), resource.AllocatableErrorInfo) {
			hasError = true
		} else {
			t.Errorf("Unexpected: %s", err.Error())
		}
	}
	if !hasError {
		t.Errorf("Expect: %s", resource.AllocatableErrorInfo)
	}

	quota, _ = apimachineryresource.ParseQuantity("6")
	_testNamespace.Spec.Hard["requests.cpu"] = quota
	errors = ValidateNamespaceUpdate(&_testNamespace, &_oldNamespace, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		t.Errorf("Unexpected: %s", err.Error())
	}
}

func newObjectGetter() validation.BusinessObjectGetter {
	return testObjectGetter{}
}

func newClusterGetter() validation.ClusterGetter {
	return testClusterGetter{}
}

type testObjectGetter struct{}

func (getter testObjectGetter) Project(name string, options apimachinerymetav1.GetOptions) (*business.Project, error) {
	project, has := _projectMap[name]
	if has {
		return project, nil
	}
	return nil, fmt.Errorf("failed to get project by name '%s'", name)
}

func (getter testObjectGetter) Namespace(project, name string, options apimachinerymetav1.GetOptions) (*business.Namespace, error) {
	namespace, has := _namespaceMap[name]
	if has {
		return namespace, nil
	}
	return nil, fmt.Errorf("failed to get namespace by project '%s' and name '%s'", project, name)
}

type testClusterGetter struct{}

func (getter testClusterGetter) Cluster(name string, options apimachinerymetav1.GetOptions) (*platformv1.Cluster, error) {
	cluster, has := _clusterMap[name]
	if has {
		return cluster, nil
	}
	return nil, fmt.Errorf("failed to get cluster by name '%s'", name)
}
