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

package imagenamespace

import (
	"fmt"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	v1 "tkestack.io/tke/api/business/v1"
	"tkestack.io/tke/pkg/util/log"
)

type imageNamespaceHealth struct {
	mu              sync.Mutex
	imageNamespaces sets.String
}

func (s *imageNamespaceHealth) Exist(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.imageNamespaces.Has(key)
}

func (s *imageNamespaceHealth) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.imageNamespaces.Delete(key)
}

func (s *imageNamespaceHealth) Set(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.imageNamespaces.Insert(key)
}

func (c *Controller) startNamespaceHealthCheck(key string) {
	if !c.health.Exist(key) {
		c.health.Set(key)
		go func() {
			if err := wait.PollImmediateUntil(1*time.Minute, c.watchNamespaceHealth(key), c.stopCh); err != nil {
				log.Error("Failed to wait poll immediate until", log.Err(err))
			}
		}()
		log.Info("ImageNamespace phase start new health check", log.String("imageNamespace", key))
	} else {
		log.Info("ImageNamespace phase health check exit", log.String("imageNamespace", key))
	}
}

// for PollImmediateUntil, when return true ,an err while exit
func (c *Controller) watchNamespaceHealth(key string) func() (bool, error) {
	return func() (bool, error) {
		log.Debug("Check imageNamespace health", log.String("key", key))

		if !c.health.Exist(key) {
			return true, nil
		}

		projectName, imageNamespaceName, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			log.Error("Failed to split meta imagenamespace key", log.String("key", key))
			c.health.Del(key)
			return true, nil
		}

		imageNamespace, err := c.client.BusinessV1().ImageNamespaces(projectName).Get(imageNamespaceName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			log.Error("ImageNamespace not found, to exit the health check loop",
				log.String("projectName", projectName), log.String("imageNamespaceName", imageNamespaceName))
			c.health.Del(imageNamespaceName)
			return true, nil
		}
		if err != nil {
			log.Error("Check imageNamespace health, imageNamespace get failed",
				log.String("projectName", projectName), log.String("imageNamespaceName", imageNamespaceName), log.Err(err))
			return false, nil
		}
		// if status is terminated,to exit the  health check loop
		if imageNamespace.Status.Phase == v1.ImageNamespaceTerminating || imageNamespace.Status.Phase == v1.ImageNamespacePending {
			log.Warn("ImageNamespace status is terminated, to exit the health check loop",
				log.String("projectName", projectName), log.String("imageNamespaceName", imageNamespaceName))
			c.health.Del(imageNamespaceName)
			return true, nil
		}

		if err := c.checkNamespaceHealth(imageNamespace); err != nil {
			log.Error("Failed to check imageNamespace health",
				log.String("projectName", projectName), log.String("imageNamespaceName", imageNamespaceName), log.Err(err))
		}
		return false, nil
	}
}

func (c *Controller) checkNamespaceHealth(imageNamespace *v1.ImageNamespace) error {
	namespaceList, err := c.registryClient.Namespaces().List(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", imageNamespace.Spec.TenantID, imageNamespace.Spec.Name),
	})
	if err != nil {
		return err
	}
	if len(namespaceList.Items) == 0 {
		return fmt.Errorf("imagenamespace %s in tenant %s not exist", imageNamespace.Spec.Name, imageNamespace.Spec.TenantID)
	}
	return nil
}
