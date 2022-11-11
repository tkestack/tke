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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"tkestack.io/tke/api/notify"
	notifyv1 "tkestack.io/tke/api/notify/v1"
	"tkestack.io/tke/pkg/util/printers"
)

// AddHandlers adds print handlers for default TKE types dealing with internal versions.
func AddHandlers(h printers.PrintHandler) {
	messageColumnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "NAME", Type: "string", Format: "name", Description: metav1.ObjectMeta{}.SwaggerDoc()["name"]},
		{Name: "ALARM POLICY", Type: "string", Description: notifyv1.MessageSpec{}.SwaggerDoc()["alarmPolicyType"]},
		{Name: "CLUSTER ID", Type: "string", Description: notifyv1.MessageSpec{}.SwaggerDoc()["clusterID"]},
		{Name: "STATUS", Type: "string", Description: notifyv1.MessageStatus{}.SwaggerDoc()["alertStatus"]},
		{Name: "CREATED AT", Type: "date", Description: metav1.ObjectMeta{}.SwaggerDoc()["creationTimestamp"]},
	}
	h.TableHandler(messageColumnDefinitions, printMSGList)
	h.TableHandler(messageColumnDefinitions, printMSG)
}

func printMSGList(msgList *notify.MessageList, options printers.PrintOptions) ([]metav1.TableRow, error) {
	rows := make([]metav1.TableRow, 0, len(msgList.Items))
	for i := range msgList.Items {
		r, err := printMSG(&msgList.Items[i], options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

func printMSG(msg *notify.Message, options printers.PrintOptions) ([]metav1.TableRow, error) {
	row := metav1.TableRow{
		Object: runtime.RawExtension{Object: msg},
	}
	row.Cells = append(row.Cells, msg.Name, msg.Spec.AlarmPolicyType, msg.Spec.ClusterID, msg.Status.AlertStatus, msg.CreationTimestamp)
	return []metav1beta1.TableRow{row}, nil
}
