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

package cluster

import (
	"sync"
	"time"

	"github.com/blang/semver"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
	resourceutil "tkestack.io/tke/pkg/util/resource"

	coreV1 "k8s.io/api/core/v1"
	v1 "tkestack.io/tke/api/platform/v1"
)

const conditionTypeHealthCheck = "HealthCheck"
const conditionTypeSyncVersion = "SyncVersion"
const reasonHealthCheckFail = "HealthCheckFail"

type clusterHealth struct {
	mu         sync.Mutex
	clusterMap map[string]*v1.Cluster
}

func (s *clusterHealth) Exist(clusterName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.clusterMap[clusterName]
	return ok
}

func (s *clusterHealth) Set(cluster *v1.Cluster) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clusterMap[cluster.Name] = cluster
}

func (s *clusterHealth) Del(clusterName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clusterMap, clusterName)
}

func (c *Controller) ensureHealthCheck(key string, cluster *v1.Cluster) {
	if c.health.Exist(key) {
		return
	}

	log.Info("start health check for cluster", log.String("clusterName", key), log.String("phase", string(cluster.Status.Phase)))
	c.health.Set(cluster)
	go wait.PollImmediateUntil(5*time.Minute, c.watchClusterHealth(cluster.Name), c.stopCh)
}

func (c *Controller) checkClusterHealth(cluster *v1.Cluster) error {
	// wait for create clustercredential, optimize first health check for user experience
	if cluster.Status.Phase == v1.ClusterInitializing {
		err := wait.PollImmediate(time.Second, time.Minute, func() (bool, error) {
			_, err := util.ClusterCredentialV1(c.client.PlatformV1(), cluster.Name)
			if err != nil {
				return false, nil
			}
			return true, nil
		})
		if err != nil { // not return! execute next steps to show reason for user
			log.Warn("wait for create clustercredential error", log.String("clusterName", cluster.Name))
		}
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		cluster.Status.Phase = v1.ClusterFailed
		cluster.Status.Message = err.Error()
		cluster.Status.Reason = reasonHealthCheckFail
		now := metav1.Now()
		c.addOrUpdateCondition(cluster, v1.ClusterCondition{
			Type:               conditionTypeHealthCheck,
			Status:             v1.ConditionFalse,
			Message:            err.Error(),
			Reason:             reasonHealthCheckFail,
			LastTransitionTime: now,
			LastProbeTime:      now,
		})
		if err1 := c.persistUpdate(cluster); err1 != nil {
			log.Warn("Update cluster status failed", log.String("clusterName", cluster.Name), log.Err(err1))
			return err1
		}
		log.Warn("Failed to build the cluster client", log.String("clusterName", cluster.Name), log.Err(err))
		return err
	}

	res, err := c.caclClusterResource(kubeClient)
	if err != nil {
		cluster.Status.Phase = v1.ClusterFailed
		cluster.Status.Message = err.Error()
		cluster.Status.Reason = reasonHealthCheckFail
		now := metav1.Now()
		c.addOrUpdateCondition(cluster, v1.ClusterCondition{
			Type:               conditionTypeHealthCheck,
			Status:             v1.ConditionFalse,
			Message:            err.Error(),
			Reason:             reasonHealthCheckFail,
			LastTransitionTime: now,
			LastProbeTime:      now,
		})
		if err1 := c.persistUpdate(cluster); err1 != nil {
			log.Warn("Update cluster status failed", log.String("clusterName", cluster.Name), log.Err(err1))
			return err1
		}
		log.Warn("Failed to build the cluster client", log.String("clusterName", cluster.Name), log.Err(err))
		return err
	}
	cluster.Status.Resource = *res

	_, err = kubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		cluster.Status.Phase = v1.ClusterFailed
		cluster.Status.Message = err.Error()
		cluster.Status.Reason = reasonHealthCheckFail
		c.addOrUpdateCondition(cluster, v1.ClusterCondition{
			Type:          conditionTypeHealthCheck,
			Status:        v1.ConditionFalse,
			Message:       err.Error(),
			Reason:        reasonHealthCheckFail,
			LastProbeTime: metav1.Now(),
		})
	} else {
		cluster.Status.Phase = v1.ClusterRunning
		cluster.Status.Message = ""
		cluster.Status.Reason = ""
		c.addOrUpdateCondition(cluster, v1.ClusterCondition{
			Type:          conditionTypeHealthCheck,
			Status:        v1.ConditionTrue,
			Message:       "",
			Reason:        "",
			LastProbeTime: metav1.Now(),
		})

		// update version info
		if cluster.Status.Version == "" {
			log.Debug("Update version info", log.String("clusterName", cluster.Name))
			if version, err := kubeClient.ServerVersion(); err == nil {
				entireVersion, err := semver.ParseTolerant(version.GitVersion)
				if err != nil {
					return err
				}
				pureVersion := semver.Version{Major: entireVersion.Major, Minor: entireVersion.Minor, Patch: entireVersion.Patch}
				log.Info("Set cluster version", log.String("clusterName", cluster.Name), log.String("version", pureVersion.String()), log.String("entireVersion", entireVersion.String()))
				cluster.Status.Version = pureVersion.String()
				now := metav1.Now()
				c.addOrUpdateCondition(cluster, v1.ClusterCondition{
					Type:               conditionTypeSyncVersion,
					Status:             v1.ConditionTrue,
					Message:            "",
					Reason:             "",
					LastProbeTime:      now,
					LastTransitionTime: now,
				})
			}
		}
	}

	if err := c.persistUpdate(cluster); err != nil {
		log.Error("Update cluster status failed", log.String("clusterName", cluster.Name), log.Err(err))
		return err
	}
	return err
}

// cal the cluster's capacity , allocatable and allocated resource
func (c *Controller) caclClusterResource(kubeClient *kubernetes.Clientset) (*v1.ClusterResource, error) {
	// cal the node's capacity and allocatable
	var cpuCapacity, memoryCapcity, cpuAllocatable, memoryAllocatable, cpuAllocated, memoryAllocated resource.Quantity

	for {
		nodeList, err := kubeClient.CoreV1().Nodes().List(metav1.ListOptions{Limit: int64(300)})
		if err != nil {
			return &v1.ClusterResource{}, err
		}

		for _, node := range nodeList.Items {
			for resourceName, capacity := range node.Status.Capacity {
				if resourceName.String() == string(resourceutil.CPU) {
					cpuCapacity.Add(capacity)
				}
				if resourceName.String() == string(resourceutil.Memory) {
					memoryCapcity.Add(capacity)
				}
			}

			for resourceName, allocatable := range node.Status.Allocatable {
				if resourceName.String() == string(resourceutil.CPU) {
					cpuAllocatable.Add(allocatable)
				}
				if resourceName.String() == string(resourceutil.Memory) {
					memoryAllocatable.Add(allocatable)
				}
			}
		}

		if nodeList.Continue == "" {
			break
		}
	}

	// cal the pods's request resource as allocated resource
	for {
		podsList, err := kubeClient.CoreV1().Pods("").List(metav1.ListOptions{Limit: int64(500)})
		if err != nil {
			return &v1.ClusterResource{}, err
		}

		for _, pod := range podsList.Items {
			// same with kubectl skip those pods in failed or succeeded status
			if pod.Status.Phase == coreV1.PodFailed || pod.Status.Phase == coreV1.PodSucceeded {
				continue
			}
			for _, container := range pod.Spec.Containers {
				for resourceName, allocated := range container.Resources.Requests {
					if resourceName.String() == string(resourceutil.CPU) {
						cpuAllocated.Add(allocated)
					}
					if resourceName.String() == string(resourceutil.Memory) {
						memoryAllocated.Add(allocated)
					}
				}
			}
		}

		if podsList.Continue == "" {
			break
		}
	}
	result := &v1.ClusterResource{
		Capacity: v1.ResourceList{
			string(resourceutil.CPU):    cpuCapacity,
			string(resourceutil.Memory): memoryCapcity,
		},
		Allocatable: v1.ResourceList{
			string(resourceutil.CPU):    cpuAllocatable,
			string(resourceutil.Memory): memoryAllocatable,
		},
		Allocated: v1.ResourceList{
			string(resourceutil.CPU):    cpuAllocated,
			string(resourceutil.Memory): memoryAllocated,
		},
	}
	return result, nil
}

// for PollImmediateUntil, when return true ,an err while exit
func (c *Controller) watchClusterHealth(clusterName string) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Check cluster health", log.String("clusterName", clusterName))

		cluster, err := c.client.PlatformV1().Clusters().Get(clusterName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.Warn("Cluster not found, to exit the health check loop", log.String("clusterName", clusterName))
				return true, nil
			}
			log.Error("Check cluster health, cluster get failed", log.String("clusterName", clusterName), log.Err(err))
			return false, nil
		}

		if cluster.Status.Phase == v1.ClusterTerminating {
			log.Warn("Cluster status is Terminating, to exit the health check loop", log.String("clusterName", cluster.Name))
			return true, nil
		}

		_ = c.checkClusterHealth(cluster)
		return false, nil
	}
}
