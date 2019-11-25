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

package namespace

import (
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"reflect"
	"sync"
	"time"
	"tkestack.io/tke/api/business/v1"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
)

type namespaceHealth struct {
	mu         sync.Mutex
	namespaces sets.String
}

func (s *namespaceHealth) Exist(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.namespaces.Has(key)
}

func (s *namespaceHealth) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.namespaces.Delete(key)
}

func (s *namespaceHealth) Set(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.namespaces.Insert(key)
}

func (c *Controller) startNamespaceHealthCheck(key string) {
	if !c.health.Exist(key) {
		c.health.Set(key)
		go func() {
			if err := wait.PollImmediateUntil(1*time.Minute, c.watchNamespaceHealth(key), c.stopCh); err != nil {
				log.Error("Failed to wait poll immediate until", log.Err(err))
			}
		}()
		log.Info("Namespace phase start new health check", log.String("namespace", key))
	} else {
		log.Info("Namespace phase health check exit", log.String("namespace", key))
	}
}

// for PollImmediateUntil, when return true ,an err while exit
func (c *Controller) watchNamespaceHealth(key string) func() (bool, error) {
	return func() (bool, error) {
		log.Debug("Check namespace health", log.String("namespace", key))

		if !c.health.Exist(key) {
			return true, nil
		}

		projectName, namespaceName, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			log.Error("Failed to split meta namespace key", log.String("key", key))
			c.health.Del(key)
			return true, nil
		}
		namespace, err := c.client.BusinessV1().Namespaces(projectName).Get(namespaceName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			log.Error("Namespace not found, to exit the health check loop", log.String("projectName", projectName), log.String("namespaceName", namespaceName))
			c.health.Del(key)
			return true, nil
		}
		if err != nil {
			log.Error("Check namespace health, namespace get failed", log.String("projectName", projectName), log.String("namespaceName", namespaceName), log.Err(err))
			return false, nil
		}
		// if status is terminated,to exit the  health check loop
		if namespace.Status.Phase == v1.NamespaceTerminating || namespace.Status.Phase == v1.NamespacePending {
			log.Warn("Namespace status is terminated, to exit the health check loop", log.String("projectName", projectName), log.String("namespaceName", namespaceName))
			c.health.Del(key)
			return true, nil
		}

		if err := c.checkNamespaceHealth(namespace); err != nil {
			log.Error("Failed to check namespace health", log.String("projectName", projectName), log.String("namespaceName", namespaceName), log.Err(err))
		}
		return false, nil
	}
}

func (c *Controller) checkNamespaceHealth(namespace *v1.Namespace) error {
	// build client
	kubeClient, err := util.BuildExternalClientSetWithName(c.platformClient, namespace.Spec.ClusterName)
	if err != nil {
		return err
	}
	message, reason := checkNamespaceOnCluster(kubeClient, namespace)
	if message != "" {
		namespace.Status.Phase = v1.NamespaceFailed
		namespace.Status.LastTransitionTime = metav1.Now()
		namespace.Status.Message = message
		namespace.Status.Reason = reason
		return c.persistUpdate(namespace)
	}
	message, reason, used := calculateNamespaceUsed(kubeClient, namespace)
	if message != "" {
		namespace.Status.Phase = v1.NamespaceFailed
		namespace.Status.LastTransitionTime = metav1.Now()
		namespace.Status.Message = message
		namespace.Status.Reason = reason
		return c.persistUpdate(namespace)
	}
	if namespace.Status.Phase != v1.NamespaceAvailable || !reflect.DeepEqual(namespace.Status.Used, used) {
		if namespace.Status.Phase != v1.NamespaceAvailable {
			namespace.Status.LastTransitionTime = metav1.Now()
		}
		namespace.Status.Phase = v1.NamespaceAvailable
		namespace.Status.Used = used
		namespace.Status.Message = ""
		namespace.Status.Reason = ""
		return c.persistUpdate(namespace)
	}
	return nil
}
