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

package storage

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/util/printers"
)

// AddHandlers adds print handlers for default TKE types dealing with internal versions.
// Refer kubernetes/pkg/printers/internalversion/printers.go:78
func AddHandlers(h printers.PrintHandler) {
	machineColumnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "Name", Type: "string", Format: "name", Description: metav1.ObjectMeta{}.SwaggerDoc()["name"]},
		{Name: "Type", Type: "string", Description: platformv1.MachineSpec{}.SwaggerDoc()["type"]},
		{Name: "IP", Type: "string", Description: platformv1.MachineSpec{}.SwaggerDoc()["ip"]},
		{Name: "Status", Type: "string", Description: platformv1.MachineStatus{}.SwaggerDoc()["phase"]},
		{Name: "Age", Type: "string", Description: metav1.ObjectMeta{}.SwaggerDoc()["creationTimestamp"]},
	}
	h.TableHandler(machineColumnDefinitions, printMachineList)
	h.TableHandler(machineColumnDefinitions, printMachine)
}

func printMachineList(machineList *platform.MachineList, options printers.PrintOptions) ([]metav1.TableRow, error) {
	rows := make([]metav1.TableRow, 0, len(machineList.Items))
	for i := range machineList.Items {
		r, err := printMachine(&machineList.Items[i], options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

func printMachine(machine *platform.Machine, options printers.PrintOptions) ([]metav1.TableRow, error) {
	row := metav1.TableRow{
		Object: runtime.RawExtension{Object: machine},
	}
	row.Cells = append(row.Cells, machine.Name, machine.Spec.Type, machine.Spec.IP, machine.Status.Phase, printers.TranslateTimestampSince(machine.CreationTimestamp))
	return []metav1beta1.TableRow{row}, nil
}
