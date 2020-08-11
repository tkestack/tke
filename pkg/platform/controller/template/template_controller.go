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

package template

import (
	"context"
	normalerrors "errors"
	"fmt"
	"reflect"
	"time"

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
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	templateClientRetryCount    = 5
	templateClientRetryInterval = 5 * time.Second

	templateMaxRetryCount = 5
	templateTimeOut       = 5 * time.Minute
)

// Controller is used to synchronize the installation, upgrade and
// uninstallation of cluster event persistence components.
type Controller struct {
	client       clientset.Interface
	cache        *templateCache
	health       *templateHealth
	checking     *templateChecking
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.TemplateLister
	listerSynced cache.InformerSynced
	stopCh       <-chan struct{}
}

// NewController creates a new Controller object.
func NewController(client clientset.Interface, templateInformer platformv1informer.TemplateInformer, resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:   client,
		cache:    &templateCache{m: make(map[string]*cachedTemplate)},
		health:   &templateHealth{m: make(map[string]*v1.Template)},
		checking: &templateChecking{templateMap: make(map[string]*v1.Template)},
		queue:    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "template"),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("template_controller", client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the template informer template handlers
	templateInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueTemplate,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldTemplate, ok1 := oldObj.(*v1.Template)
				curTemplate, ok2 := newObj.(*v1.Template)
				if ok1 && ok2 && controller.needsUpdate(oldTemplate, curTemplate) {
					controller.enqueueTemplate(newObj)
				}
			},
			DeleteFunc: controller.enqueueTemplate,
		},
		resyncPeriod,
	)
	controller.lister = templateInformer.Lister()
	controller.listerSynced = templateInformer.Informer().HasSynced

	return controller
}

func (c *Controller) enqueueTemplate(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(oldTemplate *v1.Template, newTemplate *v1.Template) bool {
	return !reflect.DeepEqual(oldTemplate, newTemplate)
}

// Run will set up the template handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting template controller")
	defer log.Info("Shutting down template controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		return fmt.Errorf("failed to wait for template caches to sync")
	}

	c.stopCh = stopCh

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
	return nil
}

// worker processes the queue of template objects.
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

	err := c.syncTemplate(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing template %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncTemplate will sync the Template with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// template created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncTemplate(key string) error {
	startTime := time.Now()
	var cachedTemplate *cachedTemplate
	defer func() {
		log.Info("Finished syncing template", log.String("templateName", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// template holds the latest template info from apiserver
	template, err := c.lister.Get(name)
	switch {
	case errors.IsNotFound(err):
		// There is no template named by key in etcd, it is a deletion task.
		log.Info("Template has been deleted. Attempting to cleanup template resources in backend template", log.String("templateName", key))
		err = c.processTemplateDeletion(context.Background(), key)
	case err != nil:
		log.Warn("Unable to retrieve template from store", log.String("templateName", key), log.Err(err))
	default:
		// otherwise, it is a addition task or updating task
		cachedTemplate = c.cache.getOrCreate(key)
		err = c.processTemplateUpdate(context.Background(), cachedTemplate, template, key)
	}

	return err
}

func (c *Controller) processTemplateDeletion(ctx context.Context, key string) error {
	cachedTemplate, ok := c.cache.get(key)
	if !ok {
		log.Error("Template not in cache even though the watcher thought it was. Ignoring the deletion", log.String("namespaceSetName", key))
		return nil
	}
	return c.processTemplateDelete(ctx, cachedTemplate, key)
}

func (c *Controller) processTemplateDelete(ctx context.Context, cachedTemplate *cachedTemplate, key string) error {
	log.Info("template will be dropped", log.String("templateName", key))

	if c.cache.Exist(key) {
		log.Info("delete the template cache", log.String("templateName", key))
		c.cache.delete(key)
	}

	if c.health.Exist(key) {
		log.Info("delete the template health cache", log.String("templateName", key))
		c.health.Del(key)
	}
	template := cachedTemplate.state
	return c.uninstallTemplateComponent(ctx, template)
}

func (c *Controller) processTemplateUpdate(ctx context.Context, cachedTemplate *cachedTemplate, template *v1.Template, key string) error {
	if cachedTemplate.state != nil {
		if cachedTemplate.state.UID != template.UID {
			// TODO check logic
			if err := c.processTemplateDelete(ctx, cachedTemplate, key); err != nil {
				return err
			}
		}
	}

	err := c.createTemplateIfNeeded(ctx, key, cachedTemplate, template)
	if err != nil {
		return err
	}

	cachedTemplate.state = template
	// Always update the cache upon success.
	c.cache.set(key, cachedTemplate)

	return nil
}

func (c *Controller) templateReinitialize(ctx context.Context, key string, cachedTemplate *cachedTemplate, template *v1.Template) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		err := c.installTemplateComponent(ctx, template)
		if err == nil {
			template = template.DeepCopy()
			template.Status.Phase = v1.AddonPhaseChecking
			template.Status.Reason = ""
			template.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
			err = c.persistUpdate(ctx, template)
			if err != nil {
				return true, err
			}
			return true, nil
		}
		// First, rollback the template
		if err := c.uninstallTemplateComponent(ctx, template); err != nil {
			log.Error("Uninstall template component error.")
			return true, err
		}
		if template.Status.RetryCount == templateMaxRetryCount {
			template = template.DeepCopy()
			template.Status.Phase = v1.AddonPhaseFailed
			template.Status.Reason = fmt.Sprintf("Install error and retried max(%d) times already.", templateMaxRetryCount)
			err := c.persistUpdate(ctx, template)
			if err != nil {
				log.Error("Update template error.")
				return true, err
			}
			return true, nil
		}
		// Add the retry count will trigger reinitialize function from the persistent controller again.
		template = template.DeepCopy()
		template.Status.Phase = v1.AddonPhaseReinitializing
		template.Status.Reason = err.Error()
		template.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
		template.Status.RetryCount++
		err = c.persistUpdate(ctx, template)
		if err != nil {
			return true, err
		}
		return true, nil
	}
}

func (c *Controller) createTemplateIfNeeded(ctx context.Context, key string, cachedTemplate *cachedTemplate, template *v1.Template) error {
	switch template.Status.Phase {
	case v1.AddonPhaseInitializing:
		log.Info("Template will be created", log.String("template", key))
		if err := c.installTemplateComponent(ctx, template); err != nil {
			template = template.DeepCopy()
			template.Status.Phase = v1.AddonPhaseReinitializing
			template.Status.Reason = err.Error()
			template.Status.RetryCount = 1
			template.Status.LastReInitializingTimestamp = metav1.Now()
			return c.persistUpdate(ctx, template)
		}
		template = template.DeepCopy()
		template.Status.Phase = v1.AddonPhaseChecking
		template.Status.Reason = ""
		template.Status.RetryCount = 0
		return c.persistUpdate(ctx, template)
	case v1.AddonPhaseReinitializing:
		var interval = time.Since(template.Status.LastReInitializingTimestamp.Time)
		var waitTime time.Duration
		if interval >= templateTimeOut {
			waitTime = time.Duration(1)
		} else {
			waitTime = templateTimeOut - interval
		}
		go wait.Poll(waitTime, templateTimeOut, c.templateReinitialize(ctx, key, cachedTemplate, template))
	case v1.AddonPhaseChecking:
		if !c.checking.Exist(key) {
			c.checking.Set(key, template)
			go wait.PollImmediate(5*time.Second, 5*time.Minute, c.checkDeploymentStatus(ctx, template, key))
		}
	case v1.AddonPhaseRunning:
		if !c.health.Exist(key) {
			c.health.Set(key, template)
			go wait.PollImmediateUntil(5*time.Minute, c.watchTemplateHealth(ctx, key), c.stopCh)
		}
	case v1.AddonPhaseFailed:
		log.Info("Template is error", log.String("template", key))
		if c.health.Exist(key) {
			c.health.Del(key)
		}
	}
	return nil
}

func (c *Controller) installTemplateComponent(ctx context.Context, template *v1.Template) error {
	log.Info("start to create the template" + template.ObjectMeta.Name)
	_, err := c.client.PlatformV1().Templates().Get(ctx, template.ObjectMeta.Name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		if _, err := c.client.PlatformV1().Templates().Create(ctx, template, metav1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (c *Controller) uninstallTemplateComponent(ctx context.Context, template *v1.Template) error {

	var failed = false
	deleteErr := c.client.PlatformV1().Templates().Delete(ctx, template.ObjectMeta.Name, metav1.DeleteOptions{})
	if deleteErr != nil && !errors.IsNotFound(deleteErr) {
		failed = true
		log.Error("Failed to delete template", log.Err(deleteErr))
	}
	if failed {
		return normalerrors.New("delete template error")
	}

	return nil
}

func (c *Controller) watchTemplateHealth(ctx context.Context, key string) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start check template health", log.String("templateName", key))
		template, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}
		if !c.health.Exist(key) {
			log.Info("health check over.")
			return true, nil
		}
		_, err = c.client.PlatformV1().Templates().Get(ctx, template.ObjectMeta.Name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			template = template.DeepCopy()
			template.Status.Phase = v1.AddonPhaseFailed
			template.Status.Reason = "Template do not exist."
			if err = c.persistUpdate(ctx, template); err != nil {
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

func (c *Controller) checkDeploymentStatus(ctx context.Context, template *v1.Template, key string) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to check the template health", log.String("templateName", template.ObjectMeta.Name))
		if !c.checking.Exist(key) {
			log.Info("checking over.")
			return true, nil
		}
		template, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}
		dep, err := c.client.PlatformV1().Templates().Get(ctx, template.ObjectMeta.Name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			template = template.DeepCopy()
			template.Status.Phase = v1.AddonPhaseFailed
			template.Status.Reason = "Template do not exist."
			if err = c.persistUpdate(ctx, template); err != nil {
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
		if !ok && time.Since(dep.CreationTimestamp.Time) > 2*time.Minute {
			template = template.DeepCopy()
			template.Status.Phase = v1.AddonPhaseFailed
			template.Status.Reason = reason
			if err = c.persistUpdate(ctx, template); err != nil {
				return false, err
			}
			c.checking.Del(key)
			return true, nil
		}
		if !ok {
			return false, nil
		}

		template = template.DeepCopy()
		template.Status.Phase = v1.AddonPhaseRunning
		template.Status.Reason = ""
		if err = c.persistUpdate(ctx, template); err != nil {
			return false, err
		}
		c.checking.Del(key)
		return true, nil
	}
}

func (c *Controller) persistUpdate(ctx context.Context, template *v1.Template) error {
	var err error
	for i := 0; i < templateClientRetryCount; i++ {
		_, err = c.client.PlatformV1().Templates().UpdateStatus(ctx, template, metav1.UpdateOptions{})
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to template template that no longer exists", log.String("templateName", template.ObjectMeta.Name), log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to template '%s' that has been changed since we received it: %v", template.ObjectMeta.Name, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of template '%s/%s'", template.ObjectMeta.Name, template.Status.Phase), log.String("templateName", template.ObjectMeta.Name), log.Err(err))
		time.Sleep(templateClientRetryInterval)
	}

	return err
}
