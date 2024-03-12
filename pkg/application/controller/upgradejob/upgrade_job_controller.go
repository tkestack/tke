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

package upgradejob

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"git.woa.com/kmetis/healthcheckpro/pb"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	applicationv1 "tkestack.io/tke/api/application/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	applicationv1informer "tkestack.io/tke/api/client/informers/externalversions/application/v1"
	applicationv1lister "tkestack.io/tke/api/client/listers/application/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/util/addon"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	controllerName = "upgradejob-controller"
)

type clsuterInfo struct {
	cluster *platformv1.Cluster
	client  *kubernetes.Clientset
}

// Controller is responsible for performing actions dependent upon an app phase.
type Controller struct {
	client            clientset.Interface
	platformClient    platformversionedclient.PlatformV1Interface
	queue             workqueue.RateLimitingInterface
	lister            applicationv1lister.UpgradeJobLister
	listerSynced      cache.InformerSynced
	healthcheckClient pb.HealthCheckerClient
	clusters          sync.Map
	region            string
	stopCh            <-chan struct{}
}

// NewController creates a new Controller object.
func NewController(
	client clientset.Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	upgradeJobInformer applicationv1informer.UpgradeJobInformer,
	resyncPeriod time.Duration, region string,
) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:         client,
		region:         region,
		platformClient: platformClient,
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	if client != nil &&
		client.ApplicationV1().RESTClient() != nil &&
		!reflect.ValueOf(client.ApplicationV1().RESTClient()).IsNil() &&
		client.ApplicationV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("tke_upgradejob_controller", client.ApplicationV1().RESTClient().GetRateLimiter())
	}

	upgradeJobInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.FilteringResourceEventHandler{
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc: controller.enqueue,
				UpdateFunc: func(oldObj, newObj interface{}) {
					old, ok1 := oldObj.(*applicationv1.UpgradeJob)
					cur, ok2 := newObj.(*applicationv1.UpgradeJob)
					if ok1 && ok2 && controller.needsUpdate(old, cur) {
						controller.enqueue(newObj)
					}
				},
				//DeleteFunc: controller.enqueue,
			},
			FilterFunc: func(obj interface{}) bool {
				up, ok := obj.(*applicationv1.UpgradeJob)
				if !ok {
					return false
				}
				if up.Status.BatchCompleteNum > up.Status.BatchOrder {
					return false
				}
				return true
			},
		},
		resyncPeriod,
	)
	controller.lister = upgradeJobInformer.Lister()
	controller.listerSynced = upgradeJobInformer.Informer().HasSynced
	controller.healthcheckClient = initHealthCheckerClient()

	return controller
}

// obj could be an *applicationv1.App, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("starting upgrade job controller")
	defer log.Info("shutting down upgrade job controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		log.Error("failed to wait for upgrade job caches to sync")
		return
	}

	c.stopCh = stopCh
	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	go c.garbageCollection()

	<-stopCh
}

// worker processes the queue of app objects.
// Each app can be in the queue at most once.
// The system ensures that no two workers can process
// the same app at the same time.
func (c *Controller) worker() {
	workFunc := func() bool {
		key, quit := c.queue.Get()
		if quit {
			return true
		}
		defer c.queue.Done(key)

		err := c.syncItem(key.(string))
		if err == nil {
			// no error, forget this entry and return
			c.queue.Forget(key)
			return false
		}

		// rather than wait for a full resync, re-add the app to the queue to be processed
		c.queue.AddRateLimited(key)
		runtime.HandleError(err)
		return false
	}

	for {
		quit := workFunc()

		if quit {
			return
		}
	}
}

// syncItem will sync the App with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// applications created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()
	defer func() {
		if time.Since(startTime) > 3*time.Second {
			log.Info("finished syncing upgradejob", log.String("upgradejob", key), log.Duration("processTime", time.Since(startTime)))
		}
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	up, err := c.lister.UpgradeJobs(namespace).Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Info("upgradejob has been deleted. Attempting to cleanup resources",
			log.String("namespace", namespace),
			log.String("name", name))
		_ = c.processDeletion(key)
		return nil
	case err != nil:
		log.Warn("unable to retrieve upgradejob from store",
			log.String("namespace", namespace),
			log.String("name", name), log.Err(err))
		return err
	default:
		return c.reconcileUpgradeJob(up)
	}
}

func (c *Controller) processDeletion(key string) error {
	return nil
}

func (c *Controller) needsUpdate(old *applicationv1.UpgradeJob, new *applicationv1.UpgradeJob) bool {
	// TODO
	return true
}

func (c *Controller) reconcileUpgradeJob(up *applicationv1.UpgradeJob) error {
	if up.Spec.Pause {
		log.Infof("upgradejob %s/%s has been paused", up.Namespace, up.Name)
		// retry on next sync
		return nil
	}

	if up.Spec.AppRefer != "" {
		app, err := c.client.ApplicationV1().Apps(up.Namespace).Get(context.TODO(), up.Spec.AppRefer, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return c.setUpgradeJobToCompleteStatus(up, "app not found")
			}
			// log error and retry on next sync
			log.Errorf("reconcileUpgradeJob get app for upgradejob %s/%s error: %s", up.Namespace, up.Name, err.Error())
			return nil
		}
		// TODO: app如果是安装在meta集群，先跳过处理，后续再看是否放开。同时需要注意，meta的client获取接口不同
		if app.Labels != nil && app.Labels["application.tkestack.io/meta"] == "true" {
			return c.setUpgradeJobToCompleteStatus(up, "app is install in meta cluster")
		}
	}

	ns, daemonsetName, err := cache.SplitMetaNamespaceKey(up.Spec.Target)
	if err != nil {
		log.Errorf("reconcileUpgradeJob parse target for upgradejob %s/%s failed: %s", up.Namespace, up.Name, err.Error())
		return c.setUpgradeJobToCompleteStatus(up, "parse target failed")
	}
	if ns == "" {
		ns = "default"
	}

	// 1.获取daemonset
	cluster, err := c.getCluster(up.Namespace)
	if err != nil {
		log.Errorf("reconcileUpgradeJob get cluster for upgradeJob %s/%s failed: %s", up.Namespace, up.Name, err.Error())
		return err
	}
	daemonSet, err := cluster.client.AppsV1().DaemonSets(ns).Get(context.TODO(), daemonsetName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return c.setUpgradeJobToCompleteStatus(up, "target daemonset not found")
		}
		log.Errorf("find daemonset %s failed: %v", daemonsetName, err)
		return err
	} else {
		if daemonSet.Spec.UpdateStrategy.Type != appsv1.OnDeleteDaemonSetStrategyType {
			log.Errorf("DaemonSet %s/%s updateStrategy uncorrect", daemonSet.Namespace, daemonSet.Name)
			return c.setUpgradeJobToCompleteStatus(up, "target daemonset updateStrategy is uncorrect")
		}
	}

	pods, err := getPodsForDaemonSet(cluster.client, daemonSet)
	if err != nil {
		log.Errorf("getPodsForDaemonSet %s failed: %v", daemonsetName, err)
		return err
	}

	// 2.获取ds的rv信息
	hash, err := getCRHashForDaemonSet(cluster.client, daemonSet)
	if err != nil {
		log.Errorf("getCRHashForDaemonSet %s failed: %v", daemonsetName, err)
		return err
	}

	log.Infof("Get pods count %d, targetVersion %s", len(pods), hash)

	// TODO: 对于在删除中的pod，如果一直terminating，是否尝试强制删除？
	if fc := getUnHealthPodsNum(pods); fc > int(*up.Spec.MaxFailed) {
		log.Warnf("%s/%s unhealth pods count %d exceed %d for ds %s", up.Namespace, up.Name, fc, *up.Spec.MaxFailed, daemonsetName)
		// 等下一个周期看是否ok
		return nil
	}

	// 3.判断是否已经升级完成，更新到complete状态
	if checkAllPodsUpdated(hash, pods) &&
		(up.Status.BatchCompleteNum >= *up.Spec.BatchNum || up.Status.BatchOrder == 0) { // 已经处理完了或者还未开始处理但其实已经升级完的场景，up资源在app升级之后创建
		log.Infof("Batch upgrade %s/%s has Completed for %s: clear status %v.", up.Namespace, up.Name, daemonsetName, up.Status)
		return c.setUpgradeJobToCompleteStatus(up, "")
	}

	// 4.新一批的处理：刚开始滚动或者上一批已完成
	if up.Status.BatchOrder == up.Status.BatchCompleteNum {
		// 对于非第一批的批次，等BatchInterval再进入下一批
		if up.Status.BatchOrder != 0 {
			intervalMinutes := time.Minute
			if up.Spec.BatchIntervalSeconds != nil {
				intervalMinutes = time.Duration(*up.Spec.BatchIntervalSeconds) * time.Second
			}
			if !up.Status.BatchCompleteTime.Add(intervalMinutes).Before(time.Now()) {
				log.Infof("Batch upgrade %s/%s/%d wait to handle next batch: %v.", up.Namespace, up.Name, up.Status.BatchOrder, up.Status)
				return nil
			}
		}

		batchOrder := up.Status.BatchOrder + 1
		updatePods := getUpdatePod(calculateBatchSize(len(pods), int(*up.Spec.BatchNum), int(batchOrder)), hash, pods)
		publisherCopy := up.DeepCopy()
		publisherCopy.Status.BatchOrder = batchOrder
		publisherCopy.Status.BatchStartTime = metav1.Time{Time: time.Now()}
		publisherCopy.Status.BatchCompleteTime = metav1.Time{}
		publisherCopy.Status.BatchUpdatedNode = getNodesFromPods(updatePods)

		if _, err := c.client.ApplicationV1().UpgradeJobs(up.Namespace).Update(context.TODO(), publisherCopy, metav1.UpdateOptions{}); err != nil {
			log.Errorf("Batch upgrade %s/%s/%d begin for next batch failed: %v %v", up.Namespace, up.Name, up.Status.BatchOrder, publisherCopy.Status, err)
			return err
		}
		return nil
	}

	// 5.当前batch的处理
	if checkPodsOnNodesUpdate(hash, pods, up.Status.BatchUpdatedNode) {
		publisherCopy := up.DeepCopy()
		publisherCopy.Status.BatchCompleteTime = metav1.Time{Time: time.Now()}
		publisherCopy.Status.BatchCompleteNum = publisherCopy.Status.BatchCompleteNum + 1
		publisherCopy.Status.BatchCompleteNodes = publisherCopy.Status.BatchCompleteNodes + int32(len(up.Status.BatchUpdatedNode))

		if _, err = c.client.ApplicationV1().UpgradeJobs(up.Namespace).Update(context.TODO(), publisherCopy, metav1.UpdateOptions{}); err != nil {
			log.Errorf("Batch upgrade %s/%s/%d handle failed: %v", up.Namespace, up.Name, up.Status.BatchOrder, err)
			return err
		}
		log.Infof("Batch upgrade %s/%s/%d handle completed for node %v", up.Namespace, up.Name, up.Status.BatchOrder, publisherCopy.Status.BatchUpdatedNode)
	} else {
		// healthcheck before action
		if c.healthcheckClient != nil {
			if err := checkHealth(c.healthcheckClient, c.region, up.Namespace, up.Spec.AppRefer, up.Status.BatchUpdatedNode); err != nil {
				log.Warnf("%s/%s healthcheck failed: %v", up.Namespace, up.Name, err)
				// 等下一个周期看是否ok
				return nil
			} else {
				log.Infof("%s/%s healthcheck success", up.Namespace, up.Name)
			}
		}

		order := up.Status.BatchOrder
		log.Infof("Batch upgrade %s/%s/%d begin to handle: %v", up.Namespace, up.Name, order, up.Status)
		ups := getPodsOnNodesToUpdate(hash, pods, up.Status.BatchUpdatedNode)
		for i, pod := range ups {
			if err := cluster.client.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{}); err != nil {
				log.Errorf("Batch upgrade %s/%s/%d: delete pod %s/%s on node %s faild: %v", up.Namespace, up.Name, order, pod.Namespace, pod.Name, pod.Spec.NodeName, err)
				return err
			}
			log.Infof("Batch upgrade %s/%s/%d: delete pod %s/%s on node %s", up.Namespace, up.Name, order, pod.Namespace, pod.Name, pod.Spec.NodeName)
			if int32(i+1) >= *up.Spec.MaxSurge {
				break
			}
		}
	}

	return nil
}

func (c *Controller) setUpgradeJobToCompleteStatus(up *applicationv1.UpgradeJob, failed string) error {
	// clear cluster cache
	c.deleteCluster(up.Namespace)

	upCopy := up.DeepCopy()
	upCopy.Status.BatchOrder = 0
	if upCopy.Status.BatchCompleteNum == 0 {
		upCopy.Status.BatchCompleteNum = 1
	}
	upCopy.Status.BatchUpdatedNode = nil
	upCopy.Status.BatchStartTime = metav1.Time{}
	upCopy.Status.BatchCompleteTime = metav1.Time{}
	if failed != "" {
		upCopy.Status.Reason = &failed
	}

	if _, err := c.client.ApplicationV1().UpgradeJobs(up.Namespace).Update(context.TODO(), upCopy, metav1.UpdateOptions{}); err != nil {
		log.Errorf("Batch upgrade %s/%s/%d handle failed: %v", up.Namespace, up.Name, up.Status.BatchOrder, err)
		return err
	}

	return c.notidyApp(up)
}

func (c *Controller) notidyApp(up *applicationv1.UpgradeJob) error {
	if up.Spec.AppRefer == "" {
		log.Infof("notidyApp for up %s/%s skiped", up.Namespace, up.Name)
		return nil
	}

	annotations := map[string]string{
		"upgradejob.application.tkestack.io/jobname": up.Name,
	}
	jsonAnnotations, _ := json.Marshal(annotations)
	patch := fmt.Sprintf(`{"metadata":{"annotations":%s}}`, jsonAnnotations)
	if _, err := c.client.ApplicationV1().Apps(up.Namespace).Patch(context.TODO(), up.Spec.AppRefer, types.MergePatchType, []byte(patch), metav1.PatchOptions{}); err != nil {
		log.Errorf("notidyApp %s/%s handle failed: %v", up.Namespace, up.Name, err)
		return err
	}
	return nil
}

func (c *Controller) getCluster(cls string) (*clsuterInfo, error) {
	if info, ok := c.clusters.Load(cls); ok {
		return info.(*clsuterInfo), nil
	}

	cluster, err := c.platformClient.Clusters().Get(context.TODO(), cls, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	client, err := addon.BuildExternalClientSet(context.TODO(), cluster, c.platformClient)
	if err != nil {
		return nil, err
	}

	info := &clsuterInfo{
		cluster: cluster,
		client:  client,
	}
	c.clusters.Store(cls, info)
	return info, nil
}

func (c *Controller) deleteCluster(cls string) {
	c.clusters.Delete(cls)
}

func (c *Controller) garbageCollection() {
	ticker := time.Tick(2 * time.Hour)
	for range ticker {
		func() {
			log.Infof("upgradejob garbageCollection begin")
			defer log.Infof("upgradejob garbageCollection end")

			ujs, err := c.lister.List(labels.Everything())
			if err != nil {
				log.Errorf("upgradejob garbageCollection list ujs failed: %v", err)
				return
			}

			for _, uj := range ujs {
				if uj.Status.BatchCompleteNum > uj.Status.BatchOrder &&
					uj.CreationTimestamp.Add(24*time.Hour).Before(time.Now()) {
					if err := c.client.ApplicationV1().UpgradeJobs(uj.Namespace).Delete(context.TODO(), uj.Name, metav1.DeleteOptions{}); err != nil {
						log.Errorf("upgradejob garbageCollection %s/%s %s failed: %v", uj.Namespace, uj.Name, uj.Spec.Target, err)
					} else {
						log.Infof("upgradejob garbageCollection %s/%s %s", uj.Namespace, uj.Name, uj.Spec.Target)
					}
				}
			}
		}()
	}
}

func getPodsForDaemonSet(clientset kubernetes.Interface, daemonSet *appsv1.DaemonSet) ([]corev1.Pod, error) {
	selector, err := metav1.LabelSelectorAsSelector(daemonSet.Spec.Selector)
	if err != nil {
		return nil, err
	}

	pods, err := clientset.CoreV1().Pods(daemonSet.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector.String(), ResourceVersion: "0",
	})
	if err != nil {
		return nil, err
	}

	var filteredPods []corev1.Pod
	for _, pod := range pods.Items {
		for _, ownerReference := range pod.OwnerReferences {
			if ownerReference.UID == daemonSet.UID {
				filteredPods = append(filteredPods, pod)
			}
		}
	}

	return filteredPods, nil
}

func getCRHashForDaemonSet(clientset kubernetes.Interface, daemonSet *appsv1.DaemonSet) (string, error) {
	selector, err := metav1.LabelSelectorAsSelector(daemonSet.Spec.Selector)
	if err != nil {
		return "", err
	}

	controllerRevisions, err := clientset.AppsV1().ControllerRevisions(daemonSet.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector.String(), ResourceVersion: "0",
	})
	if err != nil {
		return "", err
	}

	var max int
	for index, cr := range controllerRevisions.Items {
		if cr.Revision > controllerRevisions.Items[max].Revision {
			max = index
		}
	}

	return controllerRevisions.Items[max].Labels[appsv1.DefaultDaemonSetUniqueLabelKey], nil
}

// TODO: 升级过程中node发生变化，目前会重新计算，不会有影响
func calculateBatchSize(total int, batchNumber int, batchOrder int) int {
	batchSize := total / batchNumber

	remainder := total % batchNumber

	if batchOrder <= remainder {
		return batchSize + 1
	}

	return batchSize
}

func checkPodsOnNodesUpdate(targetVersion string, pods []corev1.Pod, nodes []string) bool {
	check := make(map[string]bool, len(pods))
	for _, pod := range pods {
		if pod.DeletionTimestamp != nil {
			// skip terminating pod
			continue
		}

		if pod.Labels[appsv1.DefaultDaemonSetUniqueLabelKey] == targetVersion {
			check[pod.Spec.NodeName] = true
		}
	}

	for _, node := range nodes {
		if !check[node] {
			return false
		}
	}
	return true
}

func getPodsOnNodesToUpdate(targetVersion string, pods []corev1.Pod, nodes []string) []corev1.Pod {
	var ups []corev1.Pod
	for _, pod := range pods {
		if pod.Labels[appsv1.DefaultDaemonSetUniqueLabelKey] != targetVersion && pod.DeletionTimestamp == nil {
			for _, node := range nodes {
				if pod.Spec.NodeName == node {
					ups = append(ups, pod)
					break
				}
			}
		}
	}
	return ups
}

func getUnHealthPodsNum(pods []corev1.Pod) int {
	count := 0
	for _, pod := range pods {
		if pod.Status.Phase != corev1.PodRunning {
			count++
		}
	}
	return count
}

func checkAllPodsUpdated(targetVersion string, pods []corev1.Pod) bool {
	for _, pod := range pods {
		if pod.DeletionTimestamp != nil {
			// skip terminating pod
			continue
		}

		if pod.Labels[appsv1.DefaultDaemonSetUniqueLabelKey] != targetVersion {
			return false
		}
	}
	return true
}

func getNodesFromPods(pods []corev1.Pod) []string {
	nodes := make([]string, len(pods))
	for i, pod := range pods {
		nodes[i] = pod.Spec.NodeName
	}
	return nodes
}

func getUpdatePod(num int, targetVersion string, pods []corev1.Pod) []corev1.Pod {
	var p []corev1.Pod

	for _, pod := range pods {
		if pod.DeletionTimestamp != nil {
			// skip terminating pod
			continue
		}
		if len(p) < num && pod.Labels[appsv1.DefaultDaemonSetUniqueLabelKey] != targetVersion {
			p = append(p, pod)
		}
	}

	return p
}
