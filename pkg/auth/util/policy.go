/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package util

import "fmt"

func ProjectOwnerPolicyID(tenantID string) string {
	return fmt.Sprintf("pol-%s-project-owner", tenantID)
}

func ProjectMemberPolicyID(tenantID string) string {
	return fmt.Sprintf("pol-%s-project-member", tenantID)
}

func ProjectViewerPolicyID(tenantID string) string {
	return fmt.Sprintf("pol-%s-project-viewer", tenantID)
}

func ChartGroupPullPolicyID(tenantID string) string {
	return fmt.Sprintf("pol-%s-chartgroup-pull-fake", tenantID)
}

func ChartGroupFullPolicyID(tenantID string) string {
	return fmt.Sprintf("pol-%s-chartgroup-full-fake", tenantID)
}

func ChartPullPolicyID(tenantID string) string {
	return fmt.Sprintf("pol-%s-chart-pull-fake", tenantID)
}

func ChartPushPolicyID(tenantID string) string {
	return fmt.Sprintf("pol-%s-chart-push-fake", tenantID)
}

func ChartDeletePolicyID(tenantID string) string {
	return fmt.Sprintf("pol-%s-chart-delete-fake", tenantID)
}

func ChartFullPolicyID(tenantID string) string {
	return fmt.Sprintf("pol-%s-chart-full-fake", tenantID)
}

func ChartGroupPolicyResources(cg string) []string {
	return []string{fmt.Sprintf("chartgroup:%s", cg)}
}

func ChartPolicyResources(registryNamespace string) []string {
	return []string{fmt.Sprintf("registrynamespace:%s/*", registryNamespace)}
}
