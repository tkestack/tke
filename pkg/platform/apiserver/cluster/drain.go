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
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"math"
	"time"
	"tkestack.io/tke/pkg/platform/apiserver/cluster/drain"
	"tkestack.io/tke/pkg/util/log"
)

type drainCmdOptions struct {
	Namespace string

	drainer *drain.Helper
	node    *corev1.Node
}

// DrainNode work as kubectl drain node
func DrainNode(client kubernetes.Interface, node *corev1.Node) error {
	return newDrainCmdOptions(client, node).RunDrain()
}

func newDrainCmdOptions(client kubernetes.Interface, node *corev1.Node) *drainCmdOptions {
	return &drainCmdOptions{
		drainer: &drain.Helper{
			Client:              client,
			Force:               true,
			GracePeriodSeconds:  5,
			IgnoreAllDaemonSets: true,
			Timeout:             1 * time.Minute,
			DeleteLocalData:     true,
		},
		node: node,
	}
}

// RunDrain runs the 'drain' command
func (o *drainCmdOptions) RunDrain() error {
	if err := o.RunCordonOrUncordon(true); err != nil {
		return err
	}
	err := o.deleteOrEvictPodsSimple(o.node)
	if err != nil {
		log.Errorf("unable to drain node %q", o.node.Name)
	}

	return err
}

func (o *drainCmdOptions) deleteOrEvictPodsSimple(nodeInfo *corev1.Node) error {
	list, errs := o.drainer.GetPodsForDeletion(nodeInfo.Name)
	if errs != nil {
		return utilerrors.NewAggregate(errs)
	}
	if warnings := list.Warnings(); warnings != "" {
		log.Warn(warnings)
	}

	if err := o.deleteOrEvictPods(list.Pods()); err != nil {
		pendingList, newErrs := o.drainer.GetPodsForDeletion(nodeInfo.Name)

		log.Errorf("There are pending pods in node %q when an error occurred: %v\n", nodeInfo.Name, err)
		for _, pendingPod := range pendingList.Pods() {
			log.Errorf("%s/%s\n", "pod", pendingPod.Name)
		}
		if newErrs != nil {
			log.Errorf("following errors also occurred:\n%s", utilerrors.NewAggregate(newErrs))
		}
		return err
	}
	return nil
}

// deleteOrEvictPods deletes or evicts the pods on the api server
func (o *drainCmdOptions) deleteOrEvictPods(pods []corev1.Pod) error {
	if len(pods) == 0 {
		return nil
	}

	policyGroupVersion, err := drain.CheckEvictionSupport(o.drainer.Client)
	if err != nil {
		return err
	}

	getPodFn := func(namespace, name string) (*corev1.Pod, error) {
		return o.drainer.Client.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	}

	if len(policyGroupVersion) > 0 {
		return o.evictPods(pods, policyGroupVersion, getPodFn)
	}
	return o.deletePods(pods, getPodFn)
}

func (o *drainCmdOptions) evictPods(pods []corev1.Pod, policyGroupVersion string, getPodFn func(namespace, name string) (*corev1.Pod, error)) error {
	returnCh := make(chan error, 1)

	for _, pod := range pods {
		go func(pod corev1.Pod, returnCh chan error) {
			for {
				log.Infof("evicting pod %q\n", pod.Name)
				err := o.drainer.EvictPod(pod, policyGroupVersion)
				if err == nil {
					break
				} else if apierrors.IsNotFound(err) {
					returnCh <- nil
					return
				} else if apierrors.IsTooManyRequests(err) {
					log.Errorf("error when evicting pod %q (will retry after 5s): %v\n", pod.Name, err)
					time.Sleep(5 * time.Second)
				} else {
					returnCh <- fmt.Errorf("error when evicting pod %q: %v", pod.Name, err)
					return
				}
			}
			_, err := o.waitForDelete([]corev1.Pod{pod}, 1*time.Second, time.Duration(math.MaxInt64), true, getPodFn)
			if err == nil {
				returnCh <- nil
			} else {
				returnCh <- fmt.Errorf("error when waiting for pod %q terminating: %v", pod.Name, err)
			}
		}(pod, returnCh)
	}

	doneCount := 0
	var errors []error

	// 0 timeout means infinite, we use MaxInt64 to represent it.
	var globalTimeout time.Duration
	if o.drainer.Timeout == 0 {
		globalTimeout = time.Duration(math.MaxInt64)
	} else {
		globalTimeout = o.drainer.Timeout
	}
	globalTimeoutCh := time.After(globalTimeout)
	numPods := len(pods)
	for doneCount < numPods {
		select {
		case err := <-returnCh:
			doneCount++
			if err != nil {
				errors = append(errors, err)
			}
		case <-globalTimeoutCh:
			return fmt.Errorf("drain did not complete within %v", globalTimeout)
		}
	}
	return utilerrors.NewAggregate(errors)
}

func (o *drainCmdOptions) deletePods(pods []corev1.Pod, getPodFn func(namespace, name string) (*corev1.Pod, error)) error {
	// 0 timeout means infinite, we use MaxInt64 to represent it.
	var globalTimeout time.Duration
	if o.drainer.Timeout == 0 {
		globalTimeout = time.Duration(math.MaxInt64)
	} else {
		globalTimeout = o.drainer.Timeout
	}
	for _, pod := range pods {
		err := o.drainer.DeletePod(pod)
		if err != nil && !apierrors.IsNotFound(err) {
			return err
		}
	}
	_, err := o.waitForDelete(pods, 1*time.Second, globalTimeout, false, getPodFn)
	return err
}

func (o *drainCmdOptions) waitForDelete(pods []corev1.Pod, interval, timeout time.Duration, usingEviction bool, getPodFn func(string, string) (*corev1.Pod, error)) ([]corev1.Pod, error) {
	var verbStr string
	if usingEviction {
		verbStr = "evicted"
	} else {
		verbStr = "deleted"
	}

	err := wait.PollImmediate(interval, timeout, func() (bool, error) {
		pendingPods := []corev1.Pod{}
		for i, pod := range pods {
			p, err := getPodFn(pod.Namespace, pod.Name)
			if apierrors.IsNotFound(err) || (p != nil && p.ObjectMeta.UID != pod.ObjectMeta.UID) {
				log.Infof("%s %s/%s", verbStr, pod.Namespace, pod.Name)
				continue
			} else if err != nil {
				return false, err
			} else {
				pendingPods = append(pendingPods, pods[i])
			}
		}
		pods = pendingPods
		if len(pendingPods) > 0 {
			return false, nil
		}
		return true, nil
	})
	return pods, err
}

// RunCordonOrUncordon runs either Cordon or Uncordon.  The desired value for
// "Unschedulable" is passed as the first arg.
func (o *drainCmdOptions) RunCordonOrUncordon(desired bool) error {
	cordonOrUncordon := "cordon"
	if !desired {
		cordonOrUncordon = "un" + cordonOrUncordon
	}

	cordonHelper := drain.NewCordonHelper(o.node)
	if updateRequired := cordonHelper.UpdateIfRequired(desired); !updateRequired {
		return nil
	}
	err, patchErr := cordonHelper.PatchOrReplace(o.drainer.Client)
	if patchErr != nil {
		return fmt.Errorf("error: unable to %s node %q: %v", cordonOrUncordon, o.node.Name, patchErr)
	}
	if err != nil {
		return fmt.Errorf("error: unable to %s node %q: %v", cordonOrUncordon, o.node.Name, err)
	}

	return nil
}

// already() and changed() return suitable strings for {un,}cordoning

func already(desired bool) string {
	if desired {
		return "already cordoned"
	}
	return "already uncordoned"
}

func changed(desired bool) string {
	if desired {
		return "cordoned"
	}
	return "uncordoned"
}
