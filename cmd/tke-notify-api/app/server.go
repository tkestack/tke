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

package app

import (
	"fmt"

	v1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/cmd/tke-notify-api/app/config"
	"tkestack.io/tke/pkg/notify/apiserver"
	"tkestack.io/tke/pkg/util/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

const (
	NotifyApiConfigMapName = "tke-notify-api"
	NotifyAPIAddressKey    = "notifyAPIAddress"
)

// CreateServerChain creates the apiservers connected via delegation.
func CreateServerChain(cfg *config.Config) (*genericapiserver.GenericAPIServer, error) {
	apiServerConfig := createAPIServerConfig(cfg)
	apiServer, err := CreateAPIServer(apiServerConfig, genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	apiServer.GenericAPIServer.AddPostStartHookOrDie("start-notify-api-server-informers", func(context genericapiserver.PostStartHookContext) error {
		cfg.VersionedSharedInformerFactory.Start(context.StopCh)

		// Store notify api address in configmap named tke-notify-api; It is used by prometheus addon for alertmanager to send alarms
		notifyAPIAddress := fmt.Sprintf("https://%s:%d", cfg.ExternalHost, cfg.ExternalPort)
		notifyConfigMap := &v1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				APIVersion: v1.GroupName + "/" + v1.Version,
				Kind:       "ConfigMap",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:        NotifyApiConfigMapName,
				Annotations: map[string]string{NotifyAPIAddressKey: notifyAPIAddress},
			},
		}
		cm, err := cfg.PlatformClient.ConfigMaps().Get(NotifyApiConfigMapName, metav1.GetOptions{})
		if err == nil && cm != nil {
			v, ok := cm.Annotations[NotifyAPIAddressKey]
			if !ok || v != notifyConfigMap.Annotations[NotifyAPIAddressKey] {
				notifyConfigMap.ResourceVersion = cm.ResourceVersion
				_, err = cfg.PlatformClient.ConfigMaps().Update(notifyConfigMap)
				if err != nil {
					log.Warnf("failed to update configmap for tke-notify-api due to %v", err)
					return err
				}
			}
			return nil
		}
		cm, err = cfg.PlatformClient.ConfigMaps().Create(notifyConfigMap)
		if err != nil {
			log.Warnf("failed to create configmap for tke-notify-api due to %v", err)
			return err
		}
		return nil
	})

	return apiServer.GenericAPIServer, nil
}

// CreateAPIServer creates and wires a workable tke-business-api
func CreateAPIServer(apiServerConfig *apiserver.Config, delegateAPIServer genericapiserver.DelegationTarget) (*apiserver.APIServer, error) {
	return apiServerConfig.Complete().New(delegateAPIServer)
}

func createAPIServerConfig(cfg *config.Config) *apiserver.Config {
	return &apiserver.Config{
		GenericConfig: &genericapiserver.RecommendedConfig{
			Config: *cfg.GenericAPIServerConfig,
		},
		ExtraConfig: apiserver.ExtraConfig{
			ServerName:              cfg.ServerName,
			VersionedInformers:      cfg.VersionedSharedInformerFactory,
			StorageFactory:          cfg.StorageFactory,
			APIResourceConfigSource: cfg.StorageFactory.APIResourceConfigSource,
			PrivilegedUsername:      cfg.PrivilegedUsername,
		},
	}
}
