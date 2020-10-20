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

package resourcelock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "tkestack.io/tke/api/application/v1"
	applicationv1client "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
)

// ApplicationConfigMapLock defines the structure of using configmap resources to implement
// distributed locks.
type ApplicationConfigMapLock struct {
	// ConfigMapMeta should contain a Name and a Namespace of a
	// ConfigMapMeta object that the LeaderElector will attempt to lead.
	ConfigMapMeta metav1.ObjectMeta
	Client        applicationv1client.ConfigMapsGetter
	LockConfig    Config
	cm            *v1.ConfigMap
}

// Get returns the election record from a ConfigMap Annotation
func (cml *ApplicationConfigMapLock) Get(ctx context.Context) (*LeaderElectionRecord, error) {
	var record LeaderElectionRecord
	var err error
	cml.cm, err = cml.Client.ConfigMaps().Get(ctx, cml.ConfigMapMeta.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if cml.cm.Annotations == nil {
		cml.cm.Annotations = make(map[string]string)
	}
	if recordBytes, found := cml.cm.Annotations[LeaderElectionRecordAnnotationKey]; found {
		if err := json.Unmarshal([]byte(recordBytes), &record); err != nil {
			return nil, err
		}
	}
	return &record, nil
}

// Create attempts to create a LeaderElectionRecord annotation
func (cml *ApplicationConfigMapLock) Create(ctx context.Context, ler LeaderElectionRecord) error {
	recordBytes, err := json.Marshal(ler)
	if err != nil {
		return err
	}
	cml.cm, err = cml.Client.ConfigMaps().Create(ctx, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cml.ConfigMapMeta.Name,
			Namespace: cml.ConfigMapMeta.Namespace,
			Annotations: map[string]string{
				LeaderElectionRecordAnnotationKey: string(recordBytes),
			},
		},
	}, metav1.CreateOptions{})
	return err
}

// Update will update an existing annotation on a given resource.
func (cml *ApplicationConfigMapLock) Update(ctx context.Context, ler LeaderElectionRecord) error {
	if cml.cm == nil {
		return errors.New("endpoint not initialized, call get or create first")
	}
	recordBytes, err := json.Marshal(ler)
	if err != nil {
		return err
	}
	cml.cm.Annotations[LeaderElectionRecordAnnotationKey] = string(recordBytes)
	cml.cm, err = cml.Client.ConfigMaps().Update(ctx, cml.cm, metav1.UpdateOptions{})
	return err
}

// Describe is used to convert details on current resource lock
// into a string
func (cml *ApplicationConfigMapLock) Describe() string {
	return fmt.Sprintf("%v/%v", cml.ConfigMapMeta.Namespace, cml.ConfigMapMeta.Name)
}

// Identity returns the Identity of the lock
func (cml *ApplicationConfigMapLock) Identity() string {
	return cml.LockConfig.Identity
}
