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

package persistentevent

import (
	"bytes"
	normalerrors "errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	"tkestack.io/tke/pkg/platform/controller/addon/persistentevent/images"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	v1 "tkestack.io/tke/api/platform/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	persistentEventClientRetryCount    = 5
	persistentEventClientRetryInterval = 5 * time.Second

	persistentEventMaxRetryCount = 5
	persistentEventTimeOut       = 5 * time.Minute
)

var (
	regionToDomain = map[string]string{
		"bj":   "ap-beijing.cls.myqcloud.com",
		"sh":   "ap-shanghai.cls.myqcloud.com",
		"gz":   "ap-guangzhou.cls.myqcloud.com",
		"cd":   "ap-chengdu.cls.myqcloud.com",
		"szjr": "ap-shenzhen-fsi.cls.myqcloud.com",
		"shjr": "ap-shanghai-fsi.cls.myqcloud.com",
		"na":   "na-toronto.cls.myqcloud.com",
	}
)

// Controller is used to synchronize the installation, upgrade and
// uninstallation of cluster event persistence components.
type Controller struct {
	client       clientset.Interface
	cache        *persistentEventCache
	health       *persistentEventHealth
	checking     *persistentEventChecking
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.PersistentEventLister
	listerSynced cache.InformerSynced
	stopCh       <-chan struct{}
}

// NewController creates a new Controller object.
func NewController(client clientset.Interface, persistentEventInformer platformv1informer.PersistentEventInformer, resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{

		client:   client,
		cache:    &persistentEventCache{persistentEventMap: make(map[string]*cachedPersistentEvent)},
		health:   &persistentEventHealth{persistentEventMap: make(map[string]*v1.PersistentEvent)},
		checking: &persistentEventChecking{persistentEventMap: make(map[string]*v1.PersistentEvent)},
		queue:    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "persistentevent"),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("persistentevent_controller", client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the persistent event informer event handlers
	persistentEventInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueuePersistentEvent,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldCluster, ok1 := oldObj.(*v1.PersistentEvent)
				curCluster, ok2 := newObj.(*v1.PersistentEvent)
				if ok1 && ok2 && controller.needsUpdate(oldCluster, curCluster) {
					controller.enqueuePersistentEvent(newObj)
				}
			},
			DeleteFunc: controller.enqueuePersistentEvent,
		},
		resyncPeriod,
	)
	controller.lister = persistentEventInformer.Lister()
	controller.listerSynced = persistentEventInformer.Informer().HasSynced

	return controller
}

func (c *Controller) enqueuePersistentEvent(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(oldPersistentEvent *v1.PersistentEvent, newPersistentEvent *v1.PersistentEvent) bool {
	return !reflect.DeepEqual(oldPersistentEvent, newPersistentEvent)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting persistent event controller")
	defer log.Info("Shutting down persistent event controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		return fmt.Errorf("failed to wait for persistent event caches to sync")
	}

	c.stopCh = stopCh

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
	return nil
}

// worker processes the queue of persistent event objects.
// Each cluster can be in the queue at most once.
// The system ensures that no two workers can process
// the same namespace at the same time.
func (c *Controller) worker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncPersistentEvent(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing persistent event %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncPersistentEvent will sync the PersistentEvent with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// persistent event created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncPersistentEvent(key string) error {
	startTime := time.Now()
	var cachedPersistentEvent *cachedPersistentEvent
	defer func() {
		log.Info("Finished syncing persistent event", log.String("clusterName", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// persistentEvent holds the latest persistentEvent info from apiserver
	persistentEvent, err := c.lister.Get(name)
	switch {
	case errors.IsNotFound(err):
		// There is no persistent event named by key in etcd, it is a deletion task.
		log.Info("Persistent Event has been deleted. Attempting to cleanup persistent event resources in backend persistentEvent", log.String("clusterName", key))
		err = c.processPersistentEventDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve persistent event from store", log.String("clusterName", key), log.Err(err))
	default:
		// otherwise, it is a addition task or updating task
		cachedPersistentEvent = c.cache.getOrCreate(key)
		err = c.processPersistentEventUpdate(cachedPersistentEvent, persistentEvent, key)
	}

	return err
}

func (c *Controller) processPersistentEventDeletion(key string) error {
	cachedPersistentEvent, ok := c.cache.get(key)
	if !ok {
		log.Error("PersistentEvent not in cache even though the watcher thought it was. Ignoring the deletion", log.String("namespaceSetName", key))
		return nil
	}
	return c.processPersistentEventDelete(cachedPersistentEvent, key)
}

func (c *Controller) processPersistentEventDelete(cachedPersistentEvent *cachedPersistentEvent, key string) error {
	log.Info("persistent event will be dropped", log.String("clusterName", key))

	if c.cache.Exist(key) {
		log.Info("delete the persistent event cache", log.String("clusterName", key))
		c.cache.delete(key)
	}

	if c.health.Exist(key) {
		log.Info("delete the persistent event health cache", log.String("clusterName", key))
		c.health.Del(key)
	}
	persistentEvent := cachedPersistentEvent.state
	return c.uninstallPersistentEventComponent(persistentEvent)
}

func (c *Controller) processPersistentEventUpdate(cachedPersistentEvent *cachedPersistentEvent, persistentEvent *v1.PersistentEvent, key string) error {
	if cachedPersistentEvent.state != nil {
		if cachedPersistentEvent.state.UID != persistentEvent.UID {
			// TODO check logic
			if err := c.processPersistentEventDelete(cachedPersistentEvent, key); err != nil {
				return err
			}
		}
	}

	err := c.createPersistentEventIfNeeded(key, cachedPersistentEvent, persistentEvent)
	if err != nil {
		return err
	}

	cachedPersistentEvent.state = persistentEvent
	// Always update the cache upon success.
	c.cache.set(key, cachedPersistentEvent)

	return nil
}

func (c *Controller) persistentEventReinitialize(key string, cachedPersistentEvent *cachedPersistentEvent, persistentEvent *v1.PersistentEvent) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		err := c.installPersistentEventComponent(persistentEvent)
		if err == nil {
			persistentEvent = persistentEvent.DeepCopy()
			persistentEvent.Status.Phase = v1.AddonPhaseChecking
			persistentEvent.Status.Reason = ""
			persistentEvent.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
			err = c.persistUpdate(persistentEvent)
			if err != nil {
				return true, err
			}
			return true, nil
		}
		// First, rollback the persistentEvent
		if err := c.uninstallPersistentEventComponent(persistentEvent); err != nil {
			log.Error("Uninstall persistent event component error.")
			return true, err
		}
		if persistentEvent.Status.RetryCount == persistentEventMaxRetryCount {
			persistentEvent = persistentEvent.DeepCopy()
			persistentEvent.Status.Phase = v1.AddonPhaseFailed
			persistentEvent.Status.Reason = fmt.Sprintf("Install error and retried max(%d) times already.", persistentEventMaxRetryCount)
			err := c.persistUpdate(persistentEvent)
			if err != nil {
				log.Error("Update persistent event error.")
				return true, err
			}
			return true, nil
		}
		// Add the retry count will trigger reinitialize function from the persistent controller again.
		persistentEvent = persistentEvent.DeepCopy()
		persistentEvent.Status.Phase = v1.AddonPhaseReinitializing
		persistentEvent.Status.Reason = err.Error()
		persistentEvent.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
		persistentEvent.Status.RetryCount++
		err = c.persistUpdate(persistentEvent)
		if err != nil {
			return true, err
		}
		return true, nil
	}
}

func (c *Controller) createPersistentEventIfNeeded(key string, cachedPersistentEvent *cachedPersistentEvent, persistentEvent *v1.PersistentEvent) error {
	if persistentEvent.Status.Phase == v1.AddonPhaseRunning &&
		cachedPersistentEvent != nil &&
		cachedPersistentEvent.state != nil &&
		!reflect.DeepEqual(cachedPersistentEvent.state.Spec.PersistentBackEnd, persistentEvent.Spec.PersistentBackEnd) {
		// delete from health check map
		if c.health.Exist(key) {
			c.health.Del(key)
		}
		if err := c.uninstallPersistentEventComponent(cachedPersistentEvent.state); err != nil {
			persistentEvent = persistentEvent.DeepCopy()
			persistentEvent.Status.Phase = v1.AddonPhaseFailed
			persistentEvent.Status.Reason = "Failed to delete the old storage backend"
			persistentEvent.Status.RetryCount = 0
			return c.persistUpdate(persistentEvent)
		}
		persistentEvent = persistentEvent.DeepCopy()
		persistentEvent.Status.Phase = v1.AddonPhaseInitializing
		persistentEvent.Status.Reason = ""
		persistentEvent.Status.RetryCount = 0
		return c.persistUpdate(persistentEvent)
	}

	switch persistentEvent.Status.Phase {
	case v1.AddonPhaseInitializing:
		log.Info("PersistentEvent will be created", log.String("persistentEvent", key))
		if err := c.installPersistentEventComponent(persistentEvent); err != nil {
			persistentEvent = persistentEvent.DeepCopy()
			persistentEvent.Status.Phase = v1.AddonPhaseReinitializing
			persistentEvent.Status.Reason = err.Error()
			persistentEvent.Status.RetryCount = 1
			persistentEvent.Status.LastReInitializingTimestamp = metav1.Now()
			return c.persistUpdate(persistentEvent)
		}
		persistentEvent = persistentEvent.DeepCopy()
		persistentEvent.Status.Phase = v1.AddonPhaseChecking
		persistentEvent.Status.Reason = ""
		persistentEvent.Status.RetryCount = 0
		return c.persistUpdate(persistentEvent)
	case v1.AddonPhaseReinitializing:
		var interval = time.Since(persistentEvent.Status.LastReInitializingTimestamp.Time)
		var waitTime time.Duration
		if interval >= persistentEventTimeOut {
			waitTime = time.Duration(1)
		} else {
			waitTime = persistentEventTimeOut - interval
		}
		go wait.Poll(waitTime, persistentEventTimeOut, c.persistentEventReinitialize(key, cachedPersistentEvent, persistentEvent))
	case v1.AddonPhaseChecking:
		if !c.checking.Exist(key) {
			c.checking.Set(persistentEvent)
			go wait.PollImmediate(5*time.Second, 5*time.Minute, c.checkDeploymentStatus(persistentEvent, key))
		}
	case v1.AddonPhaseRunning:
		if !c.health.Exist(key) {
			c.health.Set(persistentEvent)
			go wait.PollImmediateUntil(5*time.Minute, c.watchPersistentEventHealth(key), c.stopCh)
		}
	case v1.AddonPhaseFailed:
		log.Info("PersistentEvent is error", log.String("persistentEvent", key))
		if c.health.Exist(key) {
			c.health.Del(key)
		}
	}
	return nil
}

func (c *Controller) installPersistentEventComponent(persistentEvent *v1.PersistentEvent) error {
	log.Info("start to create Persistent event for the persistentEvent" + persistentEvent.Spec.ClusterName)

	cluster, err := c.client.PlatformV1().Clusters().Get(persistentEvent.Spec.ClusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return err
	}

	serviceAccount := c.makeServiceAccount()
	clusterRole := c.makeClusterRole()
	clusterRoleBinding := c.makeClusterRoleBinding()
	deployment := c.makeDeployment(persistentEvent.Spec.Version)
	config, err := c.makeConfigMap(&persistentEvent.Spec.PersistentBackEnd)
	if err != nil {
		return err
	}

	_, err = kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Get("tke-event-watcher", metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(serviceAccount); err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = kubeClient.RbacV1().ClusterRoles().Get("tke-event-watcher", metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		if _, err := kubeClient.RbacV1().ClusterRoles().Create(clusterRole); err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = kubeClient.RbacV1().ClusterRoleBindings().Get("tke-event-watcher-role-binding", metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(clusterRoleBinding); err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get("fluentd-config", metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		if _, err := kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Create(config); err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Get("tke-persistent-event", metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Create(deployment); err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func (c *Controller) uninstallPersistentEventComponent(persistentEvent *v1.PersistentEvent) error {
	cluster, err := c.client.PlatformV1().Clusters().Get(persistentEvent.Spec.ClusterName, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return err
	}

	var failed = false
	deployErr := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete("tke-persistent-event", &metav1.DeleteOptions{})
	if deployErr != nil && !errors.IsNotFound(deployErr) {
		failed = true
		log.Error("Failed to delete deployment", log.Err(deployErr))
	}
	configMapErr := kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Delete("fluentd-config", &metav1.DeleteOptions{})
	if configMapErr != nil && !errors.IsNotFound(configMapErr) {
		failed = true
		log.Error("Failed to delete configmap", log.Err(configMapErr))
	}
	serviceAccountErr := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete("tke-event-watcher", &metav1.DeleteOptions{})
	if serviceAccountErr != nil && !errors.IsNotFound(serviceAccountErr) {
		failed = true
		log.Error("Failed to delete service account", log.Err(deployErr))
	}
	clusterRoleErr := kubeClient.RbacV1().ClusterRoles().Delete("tke-event-watcher", &metav1.DeleteOptions{})
	if clusterRoleErr != nil && !errors.IsNotFound(clusterRoleErr) {
		failed = true
		log.Error("Failed to delete cluster role", log.Err(deployErr))
	}
	clusterRoleBindingErr := kubeClient.RbacV1().ClusterRoleBindings().Delete("tke-event-watcher-role-binding", &metav1.DeleteOptions{})
	if clusterRoleBindingErr != nil && !errors.IsNotFound(clusterRoleBindingErr) {
		failed = true
		log.Error("Failed to delete cluster role binding", log.Err(clusterRoleBindingErr))
	}

	if failed {
		return normalerrors.New("delete persistent event error")
	}

	return nil
}

func (c *Controller) watchPersistentEventHealth(key string) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start check persistent event in cluster health", log.String("cluster", key))
		persistentEvent, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		cluster, err := c.client.PlatformV1().Clusters().Get(persistentEvent.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if !c.health.Exist(cluster.Name) {
			log.Info("health check over.")
			return true, nil
		}

		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		_, err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Get("tke-persistent-event", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			persistentEvent = persistentEvent.DeepCopy()
			persistentEvent.Status.Phase = v1.AddonPhaseFailed
			persistentEvent.Status.Reason = "Persistent event deployment do not exist."
			if err = c.persistUpdate(persistentEvent); err != nil {
				return false, err
			}
			return true, nil
		}
		if err != nil {
			return false, err
		}
		return false, nil
	}
}

func (c *Controller) checkDeploymentStatus(persistentEvent *v1.PersistentEvent, key string) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to check the persistent event deployment health", log.String("clusterName", persistentEvent.Spec.ClusterName))
		cluster, err := c.client.PlatformV1().Clusters().Get(persistentEvent.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if !c.checking.Exist(key) {
			log.Info("checking over.")
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		persistentEvent, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}
		dep, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Get("tke-persistent-event", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			persistentEvent = persistentEvent.DeepCopy()
			persistentEvent.Status.Phase = v1.AddonPhaseFailed
			persistentEvent.Status.Reason = "Persistent event deployment do not exist."
			if err = c.persistUpdate(persistentEvent); err != nil {
				return false, err
			}
			c.checking.Del(key)
			return true, nil
		}
		if err != nil {
			return false, err
		}

		ok := true
		reason := ""
		for _, cond := range dep.Status.Conditions {
			if cond.Status != corev1.ConditionTrue {
				ok = false
				reason = cond.Message
				break
			}
		}
		if !ok && time.Since(dep.CreationTimestamp.Time) > 2*time.Minute {
			persistentEvent = persistentEvent.DeepCopy()
			persistentEvent.Status.Phase = v1.AddonPhaseFailed
			persistentEvent.Status.Reason = reason
			if err = c.persistUpdate(persistentEvent); err != nil {
				return false, err
			}
			c.checking.Del(key)
			return true, nil
		}
		if !ok {
			return false, nil
		}

		persistentEvent = persistentEvent.DeepCopy()
		persistentEvent.Status.Phase = v1.AddonPhaseRunning
		persistentEvent.Status.Reason = ""
		if err = c.persistUpdate(persistentEvent); err != nil {
			return false, err
		}
		c.checking.Del(key)
		return true, nil
	}
}

func (c *Controller) persistUpdate(persistentEvent *v1.PersistentEvent) error {
	var err error
	for i := 0; i < persistentEventClientRetryCount; i++ {
		_, err = c.client.PlatformV1().PersistentEvents().UpdateStatus(persistentEvent)
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to persistentEvent persistent event that no longer exists", log.String("clusterName", persistentEvent.Spec.ClusterName), log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to persistentEvent '%s' that has been changed since we received it: %v", persistentEvent.Spec.ClusterName, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of persistentEvent '%s/%s'", persistentEvent.Spec.ClusterName, persistentEvent.Status.Phase), log.String("clusterName", persistentEvent.Spec.ClusterName), log.Err(err))
		time.Sleep(persistentEventClientRetryInterval)
	}

	return err
}

func (c *Controller) makeServiceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tke-event-watcher",
			Namespace: metav1.NamespaceSystem,
		},
	}
}

func (c *Controller) makeClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tke-event-watcher",
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}
}

func (c *Controller) makeClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tke-event-watcher-role-binding",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "tke-event-watcher",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "tke-event-watcher",
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func (c *Controller) makeDeployment(version string) *appsv1.Deployment {
	var replicas int32 = 1
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tke-persistent-event",
			Labels: map[string]string{
				"qcloud-app": "tke-persistent-event",
				"k8s-app":    "tke-persistent-event",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"qcloud-app": "tke-persistent-event",
					"k8s-app":    "tke-persistent-event",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"qcloud-app": "tke-persistent-event",
						"k8s-app":    "tke-persistent-event",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "tke-persistent-event-watcher",
							Command: []string{"./tke-event-watcher"},
							Image:   images.Get(version).Watcher.FullName(),
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "event-data",
									MountPath: "/data/log",
								},
							},
						},
						{
							Name:  "tke-persistent-event-fluentd",
							Image: images.Get(version).Collector.FullName(),
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "fluentd-config",
									MountPath: "/root",
								},
								{
									Name:      "event-data",
									MountPath: "/data/log",
									ReadOnly:  true,
								},
							},
						},
					},
					ImagePullSecrets:   []corev1.LocalObjectReference{{Name: "qcloudregistrykey"}},
					RestartPolicy:      corev1.RestartPolicyAlways,
					ServiceAccountName: "tke-event-watcher",
					Volumes: []corev1.Volume{
						{
							Name: "fluentd-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "fluentd-config",
									},
								},
							},
						},
						{
							Name: "event-data",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
}

func (c *Controller) makeConfigMap(backend *v1.PersistentBackEnd) (*corev1.ConfigMap, error) {
	domain := ""
	if backend.CLS != nil {
		var err error
		domain, err = c.getCLSDomain()
		if err != nil {
			return nil, err
		}
	}

	config := fmt.Sprintf(`<source>
  @type tail
  path /data/log/*
  pos_file /data/pos
  tag host.path.*
  format json
  read_from_head true
  path_key path
</source>
<match **>
  {{if .CLS}}
  @type cls_buffered
  host %s
  port 80
  topic_id {{.CLS.TopicID}}
  {{else if .ES}}
  @type elasticsearch
  host {{.ES.IP}}
  port {{.ES.Port}}
  scheme {{.ES.Scheme}}
  index_name {{.ES.IndexName}}
  type_name tke-k8s-event
  flush_interval 5s
  {{end}}
  <buffer>
    flush_mode interval
    retry_type exponential_backoff
    total_limit_size 32MB
    chunk_limit_size 1MB
    chunk_full_threshold 0.8
    @type file
    path /var/log/td-agent/buffer/ccs.cluster.log_collector.buffer.audit-event-collector.host-path
    overflow_action block
    flush_interval 1s
    flush_thread_burst_interval 0.01
    chunk_limit_records 8000
   </buffer>
</match>`, domain)

	var err error
	t := template.New("fluentd-config")
	t, err = t.Parse(config)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	if err = t.Execute(&b, backend); err != nil {
		return nil, err
	}
	config = b.String()
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fluentd-config",
		},
		Data: map[string]string{
			"fluentd.conf": config,
		},
	}, nil
}

func (c *Controller) getCLSDomain() (string, error) {
	f, err := ioutil.ReadFile("/etc/qcloudzone")
	if err != nil {
		fmt.Printf("%s\n", err)
		return "", fmt.Errorf("read qcloudzone error")
	}
	zone := string(f)
	zone = strings.Replace(zone, "\n", "", -1)
	domain, ok := regionToDomain[zone]
	if !ok {
		domain, err := c.getCLSDomainFromConfigMap(zone)
		if err != nil {
			return "", fmt.Errorf("unsupported region")
		}
		return domain, nil
	}
	return domain, nil
}

func (c *Controller) getCLSDomainFromConfigMap(zone string) (string, error) {
	configMap, err := c.client.PlatformV1().ConfigMaps().Get("tke-controller-persistentevent", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	domain, ok := configMap.Data[fmt.Sprintf("cls-zone-%s", zone)]
	if !ok {
		return "", fmt.Errorf("unsupported region")
	}
	return domain, nil
}
