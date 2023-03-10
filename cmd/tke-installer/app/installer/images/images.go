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

package images

import (
	"reflect"
	"sort"

	"tkestack.io/tke/pkg/app/version"

	"tkestack.io/tke/pkg/util/containerregistry"
)

type BaseComponents struct {
	TKEAuthAPI            containerregistry.Image
	TKEAuthController     containerregistry.Image
	TKEPlatformAPI        containerregistry.Image
	TKEPlatformController containerregistry.Image
	TKERegistryAPI        containerregistry.Image
	TKERegistryController containerregistry.Image
	ProviderRes           containerregistry.Image
	TKEGateway            containerregistry.Image

	NginxIngress       containerregistry.Image
	KebeWebhookCertgen containerregistry.Image

	NFSProvisioner containerregistry.Image

	CsiNodeDriverRegistrar containerregistry.Image
	CsiProvisioner         containerregistry.Image
	CsiAttacher            containerregistry.Image
	CsiResizer             containerregistry.Image
	CsiSnapshotter         containerregistry.Image
	CephCsi                containerregistry.Image
}

type ExComponents struct {
	Registry containerregistry.Image
	Busybox  containerregistry.Image
	InfluxDB containerregistry.Image
	Thanos   containerregistry.Image

	TKEBusinessAPI           containerregistry.Image
	TKEBusinessController    containerregistry.Image
	TKEMonitorAPI            containerregistry.Image
	TKEMonitorController     containerregistry.Image
	TKENotifyAPI             containerregistry.Image
	TKENotifyController      containerregistry.Image
	TKELogagentAPI           containerregistry.Image
	TKELogagentController    containerregistry.Image
	TKEAudit                 containerregistry.Image
	TKEApplicationAPI        containerregistry.Image
	TKEApplicationController containerregistry.Image
	TKEMeshAPI               containerregistry.Image
	TKEMeshController        containerregistry.Image
}

type Components struct {
	BaseComponents
	ExComponents
}

func (c Components) Get(name string) *containerregistry.Image {
	v := reflect.ValueOf(c)
	for i := 0; i < v.NumField(); i++ {
		v, _ := v.Field(i).Interface().(containerregistry.Image)
		if v.Name == name {
			return &v
		}
	}
	return nil
}

var Version = version.Get().GitVersion

var exComponents = ExComponents{
	Registry: containerregistry.Image{Name: "registry", Tag: "2.7.1"},
	Busybox:  containerregistry.Image{Name: "busybox", Tag: "1.31.1"},
	InfluxDB: containerregistry.Image{Name: "influxdb", Tag: "1.8.10"},
	Thanos:   containerregistry.Image{Name: "thanos", Tag: "v0.15.0"},

	TKEBusinessAPI:           containerregistry.Image{Name: "tke-business-api", Tag: Version},
	TKEBusinessController:    containerregistry.Image{Name: "tke-business-controller", Tag: Version},
	TKEMonitorAPI:            containerregistry.Image{Name: "tke-monitor-api", Tag: Version},
	TKEMonitorController:     containerregistry.Image{Name: "tke-monitor-controller", Tag: Version},
	TKENotifyAPI:             containerregistry.Image{Name: "tke-notify-api", Tag: Version},
	TKENotifyController:      containerregistry.Image{Name: "tke-notify-controller", Tag: Version},
	TKELogagentAPI:           containerregistry.Image{Name: "tke-logagent-api", Tag: Version},
	TKELogagentController:    containerregistry.Image{Name: "tke-logagent-controller", Tag: Version},
	TKEAudit:                 containerregistry.Image{Name: "tke-audit-api", Tag: Version},
	TKEApplicationAPI:        containerregistry.Image{Name: "tke-application-api", Tag: Version},
	TKEApplicationController: containerregistry.Image{Name: "tke-application-controller", Tag: Version},
	TKEMeshAPI:               containerregistry.Image{Name: "tke-mesh-api", Tag: Version},
	TKEMeshController:        containerregistry.Image{Name: "tke-mesh-controller", Tag: Version},
}

var baseComponents = BaseComponents{
	TKEAuthAPI:            containerregistry.Image{Name: "tke-auth-api", Tag: Version},
	TKEAuthController:     containerregistry.Image{Name: "tke-auth-controller", Tag: Version},
	TKEPlatformAPI:        containerregistry.Image{Name: "tke-platform-api", Tag: Version},
	TKEPlatformController: containerregistry.Image{Name: "tke-platform-controller", Tag: Version},
	TKERegistryAPI:        containerregistry.Image{Name: "tke-registry-api", Tag: Version},
	TKERegistryController: containerregistry.Image{Name: "tke-registry-controller", Tag: Version},
	ProviderRes:           containerregistry.Image{Name: "provider-res", Tag: "v1.21.4-5"},
	TKEGateway:            containerregistry.Image{Name: "tke-gateway", Tag: Version},

	NginxIngress:       containerregistry.Image{Name: "ingress-nginx-controller", Tag: "v1.1.3"},
	KebeWebhookCertgen: containerregistry.Image{Name: "kube-webhook-certgen", Tag: "v1.1.1"},

	NFSProvisioner: containerregistry.Image{Name: "nfs-subdir-external-provisioner", Tag: "v4.0.2"},

	CsiNodeDriverRegistrar: containerregistry.Image{Name: "csi-node-driver-registrar", Tag: "v2.4.0"},
	CsiProvisioner:         containerregistry.Image{Name: "csi-provisioner", Tag: "v3.1.0"},
	CsiAttacher:            containerregistry.Image{Name: "csi-attacher", Tag: "v3.4.0"},
	CsiResizer:             containerregistry.Image{Name: "csi-resizer", Tag: "v1.4.0"},
	CsiSnapshotter:         containerregistry.Image{Name: "csi-snapshotter", Tag: "v4.2.0"},
	CephCsi:                containerregistry.Image{Name: "cephcsi", Tag: "v3.6.1-csp2.8.3.1216"},
}

var components = Components{baseComponents, exComponents}

func List() []string {
	var items []string
	v := reflect.ValueOf(components)
	for i := 0; i < v.NumField(); i++ {
		v, _ := v.Field(i).Interface().(containerregistry.Image)
		items = append(items, v.BaseName())
	}
	sort.Strings(items)

	return items
}

func Get() Components {
	return components
}

func ListBaseComponents() []string {
	var items []string
	v := reflect.ValueOf(baseComponents)
	for i := 0; i < v.NumField(); i++ {
		v, _ := v.Field(i).Interface().(containerregistry.Image)
		items = append(items, v.BaseName())
	}
	sort.Strings(items)

	return items
}

func ListExComponents() []string {
	var items []string
	v := reflect.ValueOf(exComponents)
	for i := 0; i < v.NumField(); i++ {
		v, _ := v.Field(i).Interface().(containerregistry.Image)
		items = append(items, v.BaseName())
	}
	sort.Strings(items)

	return items
}
