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

package markcontrolplane

import (
	corev1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/util/apiclient"
	"tkestack.io/tke/pkg/util/log"
)

type Option struct {
	NodeName string
	Taints   []corev1.Taint
}

// Install taints the control-plane and sets the control-plane label
func Install(client clientset.Interface, option *Option) error {
	log.Infof("Marking the node %s as control-plane by adding the label \"%s=''\"\n", option.NodeName, constants.LabelNodeRoleMaster)
	if option.Taints != nil && len(option.Taints) > 0 {
		taintStrs := []string{}
		for _, taint := range option.Taints {
			taintStrs = append(taintStrs, taint.ToString())
		}
		log.Infof("Marking the node %s as control-plane by adding the taints %v\n", option.NodeName, taintStrs)
	}

	return apiclient.PatchNode(client, option.NodeName, func(n *corev1.Node) {
		markMasterNode(n, option.Taints)
	})
}

func markMasterNode(n *corev1.Node, taints []corev1.Taint) {
	n.ObjectMeta.Labels[constants.LabelNodeRoleMaster] = ""

	for _, nt := range n.Spec.Taints {
		if !taintExists(nt, taints) {
			taints = append(taints, nt)
		}
	}

	n.Spec.Taints = taints
}

func taintExists(taint corev1.Taint, taints []corev1.Taint) bool {
	for _, t := range taints {
		if t == taint {
			return true
		}
	}

	return false
}
