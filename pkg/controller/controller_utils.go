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

package controller

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"tkestack.io/tke/pkg/util/log"
	// register auth group api scheme for api server.
	_ "tkestack.io/tke/api/platform/install"
)

var (
	// KeyFunc checks for DeletedFinalStateUnknown objects before calling MetaNamespaceKeyFunc.
	KeyFunc = cache.DeletionHandlingMetaNamespaceKeyFunc
)

// WaitForCacheSync is a wrapper around cache.WaitForCacheSync that generates log messages
// indicating that the controller identified by controllerName is waiting for syncs, followed by
// either a successful or failed sync.
func WaitForCacheSync(controllerName string, stopCh <-chan struct{}, cacheSyncs ...cache.InformerSynced) bool {
	log.Info("Waiting for caches to sync for controller", log.String("controllerName", controllerName))

	if !cache.WaitForCacheSync(stopCh, cacheSyncs...) {
		runtime.HandleError(fmt.Errorf("unable to sync caches for %s controller", controllerName))
		return false
	}

	log.Info("Caches are synced for controller", log.String("controllerName", controllerName))
	return true
}

// IsClusterVersionBefore1_9 to check if cluster version before v1.9.x
// DO NOT delete it because TKE support cluster v1.8.x
func IsClusterVersionBefore1_9(kubeClient *kubernetes.Clientset) bool {
	return isClusterVersionBefore(kubeClient, 1, 9)
}

func isClusterVersionBefore(kubeClient *kubernetes.Clientset, majorV int, minorV int) bool {
	version, err := kubeClient.Discovery().ServerVersion()
	if err != nil {
		log.Error("error in isClusterVersionBefore")
		return false
	}
	valid := regexp.MustCompile("[0-9]")
	version.Minor = strings.Join(valid.FindAllString(version.Minor, -1), "")

	versionMajor, err := strconv.Atoi(version.Major)
	if err != nil {
		log.Error("error in isClusterVersionBefore")
		return false
	}

	versionMinor, err := strconv.Atoi(version.Minor)
	if err != nil {
		log.Error("error in isClusterVersionBefore")
		return false
	}
	if versionMajor < majorV || (versionMajor == majorV && versionMinor < minorV) {
		return true
	}
	return false
}

// Int32Ptr translate int32 to pointer
func Int32Ptr(i int32) *int32 {
	o := i
	return &o
}

func BoolPtr(i bool) *bool {
	o := i
	return &o
}

func Int64Ptr(i int64) *int64 {
	o := i
	return &o
}

// DeleteReplicaSetApp delete the replicaset and pod additionally for deployment app with extension group
func DeleteReplicaSetApp(kubeClient *kubernetes.Clientset, options metav1.ListOptions) error {
	rsList, err := kubeClient.ExtensionsV1beta1().ReplicaSets(metav1.NamespaceSystem).List(options)
	if err != nil {
		return err
	}

	var errs []error
	for i := range rsList.Items {
		rs := &rsList.Items[i]
		// update replicas to zero
		rs.Spec.Replicas = Int32Ptr(0)
		_, err = kubeClient.ExtensionsV1beta1().ReplicaSets(metav1.NamespaceSystem).Update(rs)
		if err != nil {
			errs = append(errs, err)
		} else {
			// delete replicaset
			err = kubeClient.ExtensionsV1beta1().ReplicaSets(metav1.NamespaceSystem).Delete(rs.Name, &metav1.DeleteOptions{})
			if err != nil && !errors.IsNotFound(err) {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		errMsg := ""
		for _, e := range errs {
			errMsg += e.Error() + ";"
		}
		return fmt.Errorf("delete prometheus fail:%s", errMsg)
	}

	return nil
}

// CatchPanic handles any panics that might occur during the handlePhase
func CatchPanic(funcName string, addon string) {
	if err := recover(); err != nil {
		runtime.HandleError(fmt.Errorf("recover from %s.%s(), err is %v", addon, funcName, err))
	}
}
