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

package project

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
	InvalidCluster    = "invalid-cluster"
	ClusterName       = "cluster-name"
	ParentDisplayName = "parent-display-name"
	DisplayName       = "display-name"
	NamespaceName     = "namespace-name"
	ParentProjectName = "parent-project-name"
	ProjectName       = "project-name"
	TenantID          = "tenant-id"
	Phase             = "Active"

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

	_parentProject = business.Project{
		ObjectMeta: apimachinerymetav1.ObjectMeta{
			Name:            ParentProjectName,
			ResourceVersion: "v2",
		},
		Spec: business.ProjectSpec{
			DisplayName: ParentDisplayName,
			TenantID:    TenantID,
		},
		Status: business.ProjectStatus{},
	}
	_childProject = business.Project{
		ObjectMeta: apimachinerymetav1.ObjectMeta{
			Name:            ProjectName,
			ResourceVersion: "v2",
		},
		Spec: business.ProjectSpec{
			DisplayName:       DisplayName,
			ParentProjectName: ParentProjectName,
			TenantID:          TenantID,
		},
		Status: business.ProjectStatus{},
	}
	_oldChildProject = business.Project{
		ObjectMeta: apimachinerymetav1.ObjectMeta{
			Name:            ProjectName,
			ResourceVersion: "v1",
		},
		Spec: business.ProjectSpec{
			DisplayName:       DisplayName,
			ParentProjectName: ParentProjectName,
			TenantID:          TenantID,
		},
		Status: business.ProjectStatus{},
	}
	_projectMap = map[string]*business.Project{
		ProjectName:       &_childProject,
		ParentProjectName: &_parentProject,
	}

	_childNamespace = business.Namespace{
		ObjectMeta: apimachinerymetav1.ObjectMeta{
			ClusterName:     ClusterName,
			Name:            NamespaceName,
			Namespace:       ParentProjectName,
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
	_namespaceMap = map[string]*business.Namespace{
		NamespaceName: &_childNamespace,
	}
)

func TestClusterLimitation(t *testing.T) {
	_childProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {},
	}

	_parentProject.Spec.Clusters = business.ClusterHard{
		InvalidCluster: {},
	}

	hasError := false
	errors := ValidateProjectUpdate(&_childProject, &_oldChildProject, newObjectGetter(), newClusterGetter())
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

	_parentProject.Spec.Clusters = business.ClusterHard{
		InvalidCluster: {},
		ClusterName:    {},
	}
	errors = ValidateProjectUpdate(&_childProject, &_oldChildProject, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		t.Errorf("Unexpected: %s", err.Error())
	}
}

func TestQuotaLimitation(t *testing.T) {
	quota, _ := apimachineryresource.ParseQuantity("1")
	_parentProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}

	_childProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {Hard: business.ResourceList{}},
	}

	hasError := false
	errors := ValidateProjectUpdate(&_childProject, &_oldChildProject, newObjectGetter(), newClusterGetter())
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
	_parentProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}

	quota, _ = apimachineryresource.ParseQuantity("10")
	_childProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}

	hasError := false
	errors := ValidateProjectUpdate(&_childProject, &_oldChildProject, newObjectGetter(), newClusterGetter())
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
	_childProject.Spec.Clusters[ClusterName].Hard["requests.cpu"] = quota
	errors = ValidateProjectUpdate(&_childProject, &_oldChildProject, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		t.Errorf("Unexpected: %s", err.Error())
	}
}

func TestUpdateAllocatable(t *testing.T) {
	/* remaining quota: 10 - 7 */
	quota, _ := apimachineryresource.ParseQuantity("10")
	_parentProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}
	quota, _ = apimachineryresource.ParseQuantity("7")
	_parentProject.Status.Clusters = business.ClusterUsed{
		ClusterName: {
			Used: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}

	/* request more quota: 7 - 3*/
	quota, _ = apimachineryresource.ParseQuantity("7")
	_childProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}
	_parentProject.Status.CalculatedChildProjects = []string{_childProject.Name}
	quota, _ = apimachineryresource.ParseQuantity("3")
	_oldChildProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}

	hasError := false
	errors := ValidateProjectUpdate(&_childProject, &_oldChildProject, newObjectGetter(), newClusterGetter())
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
	_childProject.Spec.Clusters[ClusterName].Hard["requests.cpu"] = quota
	errors = ValidateProjectUpdate(&_childProject, &_oldChildProject, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		t.Errorf("Unexpected: %s", err.Error())
	}
}

func TestUpdateQuota(t *testing.T) {
	_childProject.Spec.ParentProjectName = ""
	_oldChildProject.Spec.ParentProjectName = ""

	quota, _ := apimachineryresource.ParseQuantity("3")
	_childProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}
	quota, _ = apimachineryresource.ParseQuantity("7")
	_childProject.Status.Clusters = business.ClusterUsed{
		ClusterName: {
			Used: business.ResourceList{
				"requests.cpu": quota,
			},
		},
	}

	hasError := false
	errors := ValidateProjectUpdate(&_childProject, &_oldChildProject, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		if strings.Contains(err.Error(), resource.UpdateQuotaErrorInfo) {
			hasError = true
		} else {
			t.Errorf("Unexpected: %s", err.Error())
		}
	}
	if !hasError {
		t.Errorf("Expect: %s", resource.UpdateQuotaErrorInfo)
	}

	quota, _ = apimachineryresource.ParseQuantity("7")
	_childProject.Spec.Clusters[ClusterName].Hard["requests.cpu"] = quota
	errors = ValidateProjectUpdate(&_childProject, &_oldChildProject, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		t.Errorf("Unexpected: %s", err.Error())
	}

	_childProject.Spec.ParentProjectName = ParentProjectName
	_oldChildProject.Spec.ParentProjectName = ParentProjectName
}

func TestNewQuota(t *testing.T) {
	quota, _ := apimachineryresource.ParseQuantity("10")
	_childProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"pods": quota,
			},
		},
	}

	_parentProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"pods":         quota,
				"requests.cpu": quota,
			},
		},
	}

	_parentProject.Status.CalculatedChildProjects = []string{_childProject.Name}
	hasError := false
	errors := ValidateProjectUpdate(&_parentProject, &_parentProject, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		if strings.Contains(err.Error(), _addNewQuotaErrorInfo) {
			hasError = true
		} else {
			t.Errorf("Unexpected: %s", err.Error())
		}
	}
	if !hasError {
		t.Errorf("Expect: %s", _addNewQuotaErrorInfo)
	}
	_childProject.Spec.Clusters = business.ClusterHard{
		ClusterName: {
			Hard: business.ResourceList{
				"pods":         quota,
				"requests.cpu": quota,
			},
		},
	}
	errors = ValidateProjectUpdate(&_parentProject, &_parentProject, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		t.Errorf("Unexpected: %s", err.Error())
	}

	_parentProject.Status.CalculatedNamespaces = []string{_childNamespace.Name}
	hasError = false
	errors = ValidateProjectUpdate(&_parentProject, &_parentProject, newObjectGetter(), newClusterGetter())
	for _, err := range errors {
		if strings.Contains(err.Error(), _addNewQuotaErrorInfo) {
			hasError = true
		} else {
			t.Errorf("Unexpected: %s", err.Error())
		}
	}
	if !hasError {
		t.Errorf("Expect: %s", _addNewQuotaErrorInfo)
	}
	_childNamespace.Spec.Hard = business.ResourceList{
		"pods":         quota,
		"requests.cpu": quota,
	}
	errors = ValidateProjectUpdate(&_parentProject, &_parentProject, newObjectGetter(), newClusterGetter())
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
