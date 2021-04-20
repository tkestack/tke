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

package action

import (
	"bytes"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	"tkestack.io/tke/pkg/util/log"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/kubectl/pkg/cmd/get"
)

// StatusOptions is the query options to a status call.
type StatusOptions struct {
	Namespace   string
	ReleaseName string
}

type ReleaseStatus struct {
	Release   *release.Release
	Resources map[string][]string
}

// KubeClient represents a client capable of communicating with the Kubernetes API.
type KubeClient struct {
	kube.Interface
}

// ResourceActorFunc performs an action on a single resource.
type ResourceActorFunc func(*resource.Info) error

// Status returning release status.
func (c *Client) Status(options *StatusOptions) (ReleaseStatus, error) {
	actionConfig, err := c.buildActionConfig(options.Namespace)
	if err != nil {
		return ReleaseStatus{}, err
	}
	client := action.NewStatus(actionConfig)
	rel, err := client.Run(options.ReleaseName)
	if err != nil {
		return ReleaseStatus{}, err
	}

	resources, err := actionConfig.KubeClient.Build(bytes.NewBufferString(rel.Manifest), true)
	if err != nil {
		return ReleaseStatus{}, err
	}

	kubeClient := &KubeClient{actionConfig.KubeClient}
	resp, err := kubeClient.GetResource(resources)
	if err != nil {
		return ReleaseStatus{}, err
	}
	status := ReleaseStatus{
		Release:   rel,
		Resources: resp,
	}

	return status, nil
}

// GetResource gets Kubernetes resources as pretty-printed string.
//
// Namespace will set the namespace.
func (c *KubeClient) GetResource(resources kube.ResourceList) (map[string][]string, error) {
	printFlags := get.NewGetPrintFlags()
	typePrinter, _ := printFlags.JSONYamlPrintFlags.ToPrinter("json")
	printer := &get.TablePrinter{Delegate: typePrinter}
	objs := make(map[string][]runtime.Object)
	objsJSON := make(map[string][]string)
	err := resources.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}
		gvk := info.ResourceMapping().GroupVersionKind
		vk := gvk.Group + "/" + gvk.Version + "/" + gvk.Kind
		obj, err := getResource(info)
		if err != nil {
			return err
		}
		objs[vk] = append(objs[vk], obj)
		return nil
	})
	if err != nil {
		return nil, err
	}
	for t, ot := range objs {
		for _, o := range ot {
			buf := new(bytes.Buffer)
			if err := printer.PrintObj(o, buf); err != nil {
				log.Warnf("failed to print object type %s, object: %q :\n %v", t, o, err)
				return nil, err
			}
			objsJSON[t] = append(objsJSON[t], buf.String())
		}
	}
	return objsJSON, nil
}

func getResource(info *resource.Info) (runtime.Object, error) {
	obj, err := resource.NewHelper(info.Client, info.Mapping).Get(info.Namespace, info.Name)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
