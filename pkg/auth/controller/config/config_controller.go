/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"time"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider/local"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"tkestack.io/tke/pkg/util/metrics"

	v1 "tkestack.io/tke/api/auth/v1"
	controllerutil "tkestack.io/tke/pkg/controller"

	"github.com/howeyc/fsnotify"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"tkestack.io/tke/api/auth"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	authv1informer "tkestack.io/tke/api/client/informers/externalversions/auth/v1"
	authv1lister "tkestack.io/tke/api/client/listers/auth/v1"
	"tkestack.io/tke/pkg/util/log"
)

const (
	watchDebounceDelay = 100 * time.Millisecond

	controllerName = "config-controller"
)

// Controller is responsible for performing actions dependent upon a message request controller phase.
type Controller struct {
	client       clientset.Interface
	queue        workqueue.RateLimitingInterface
	watcher      *fsnotify.Watcher
	categoryFile string
	policyFile   string
	username     string
	password     string

	identityProviderLister       authv1lister.IdentityProviderLister
	identityProviderListerSynced cache.InformerSynced

	stopCh <-chan struct{}
}

func NewController(client clientset.Interface, identityProviderInformer authv1informer.IdentityProviderInformer, resyncPeriod time.Duration,
	policyFile string, categoryFile string, username string, password string) *Controller {
	controller := Controller{
		client:       client,
		queue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		categoryFile: categoryFile,
		policyFile:   policyFile,
		username:     username,
		password:     password,
	}

	if client != nil && client.AuthV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("config_controller", client.AuthV1().RESTClient().GetRateLimiter())
	}

	identityProviderInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    controller.enqueue,
			DeleteFunc: controller.enqueue,
		},
		resyncPeriod,
	)

	controller.identityProviderLister = identityProviderInformer.Lister()
	controller.identityProviderListerSynced = identityProviderInformer.Informer().HasSynced

	return &controller

}

// obj could be an *v1.identityprovider, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.AddRateLimited(key)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer runtime.HandleCrash()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting config controller")
	defer log.Info("Shutting down config controller")

	if ok := cache.WaitForCacheSync(stopCh, c.identityProviderListerSynced); !ok {
		log.Error("Failed to wait for identyProvider caches to sync")
	}

	c.stopCh = stopCh
	if err := c.loadConfig(); err != nil {
		log.Errorf("Preload config failed", log.Err(err))
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("New file watcher failed", log.Err(err))
		return
	}

	// watch the parent directory of the target files so we can catch
	// symlink updates of k8s ConfigMaps volumes.
	for _, file := range []string{c.policyFile, c.categoryFile} {
		if file != "" {
			watchDir, _ := filepath.Split(file)
			if err := watcher.Watch(watchDir); err != nil {
				log.Error("could not watch %v: %v", log.String("file", file), log.Err(err))
				return
			}
		}
	}

	c.watcher = watcher
	go c.pollReload(stopCh)

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of localIdentity objects.
// Each localIdentity can be in the queue at most once.
// The system ensures that no two workers can process
// the same localIdentity at the same time.
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

		// rather than wait for a full resync, re-add the localIdentity to the queue to be processed
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

// syncItem will sync the localIdentity with the given key if it has had
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()

	defer func() {
		log.Info("Finished syncing identityProvider", log.String("name", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	idp, err := c.identityProviderLister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Infof("IdentityProvider has been deleted %v", key)
		return nil
	case err != nil:
		log.Warn("Unable to retrieve identityProvider from store", log.String("name", key), log.Err(err))
	default:
		log.Info("Init config and tenant admin for identityProvider", log.Any("name", idp.Name))
		err = c.processConfig(idp, key)
	}
	return err
}

func (c *Controller) processConfig(idp *v1.IdentityProvider, key string) error {
	var errs []error

	if c.policyFile != "" {
		if err := c.loadPolicy(idp.Name); err != nil {
			errs = append(errs, err)
		}
	}

	if idp.Spec.Type == local.ConnectorType {
		if err := c.createAdmin(idp.Name); err != nil {
			errs = append(errs, err)
		}
	}

	return utilerrors.NewAggregate(errs)
}

// pollReload watch auth config file from the file system and reload config into storage.
func (c *Controller) pollReload(stopCh <-chan struct{}) {
	defer c.watcher.Close()
	var timerC <-chan time.Time

	for {
		select {
		case <-timerC:
			timerC = nil
			log.Info("Policy config directory changed, loadConfig it")

			if err := c.loadConfig(); err != nil {
				log.Error("Load config failed after changed", log.Err(err))
			}
		case event := <-c.watcher.Event:
			// use a timer to debounce configuration updates
			if (event.IsModify() || event.IsCreate()) && timerC == nil {
				timerC = time.After(watchDebounceDelay)
			}
		case err := <-c.watcher.Error:
			log.Error("Watcher error: %v", log.Err(err))
		case <-stopCh:
			return
		}
	}
}

func (c *Controller) loadConfig() error {
	if c.categoryFile != "" {
		err := c.loadCategory()
		if err != nil {
			return err
		}
	}

	if c.policyFile != "" {
		idps, err := c.identityProviderLister.List(labels.Everything())
		if err != nil {
			log.Error("List all tenant failed", log.Err(err))
			return err
		}

		for _, idp := range idps {
			err := c.loadPolicy(idp.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Controller) loadCategory() error {
	var categoryList []*v1.Category
	bytes, err := ioutil.ReadFile(c.categoryFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &categoryList)
	if err != nil {
		return err
	}

	var errs []error

	for _, cat := range categoryList {
		result, err := c.client.AuthV1().Categories().Get(cat.Name, metav1.GetOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			errs = append(errs, err)
			continue
		}

		if err != nil {
			log.Info("Create a new policy category", log.String("id", cat.Name), log.String("displayName", cat.Spec.DisplayName))
			_, err = c.client.AuthV1().Categories().Create(cat)

			if err != nil {
				errs = append(errs, err)
			}
		} else {
			if !reflect.DeepEqual(result.Spec, cat.Spec) {
				log.Info("Update policy category", log.String("id", cat.Name), log.String("displayName", cat.Spec.DisplayName))
				result.Spec = cat.Spec
				_, err = c.client.AuthV1().Categories().Update(result)

				if err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func (c *Controller) loadPolicy(tenantID string) error {
	log.Info("Handle default policy for tenant", log.String("tenantID", tenantID))
	var policyList []*v1.Policy

	bytes, err := ioutil.ReadFile(c.policyFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &policyList)
	if err != nil {
		return err
	}

	var errs []error

	for _, pol := range policyList {
		policySelector := fields.AndSelectors(
			fields.OneTermEqualSelector("spec.tenantID", tenantID),
			fields.OneTermEqualSelector("spec.displayName", pol.Spec.DisplayName),
			fields.OneTermEqualSelector("spec.type", string(auth.PolicyDefault)),
		)

		result, err := c.client.AuthV1().Policies().List(metav1.ListOptions{FieldSelector: policySelector.String()})
		if err != nil {
			return err
		}
		pol.Spec.Type = v1.PolicyDefault
		pol.Spec.TenantID = tenantID
		pol.Spec.Username = "admin"
		pol.Spec.Finalizers = []v1.FinalizerName{
			v1.PolicyFinalize,
		}
		if len(result.Items) > 0 {
			exists := result.Items[0]
			if !reflect.DeepEqual(exists.Spec, pol.Spec) {
				log.Info("Update default policy", log.String("displayName", pol.Spec.DisplayName))
				exists.Spec = pol.Spec
				_, err = c.client.AuthV1().Policies().Update(&exists)

				if err != nil {
					errs = append(errs, err)
				}
			}
		} else {
			log.Info("Create a new default policy", log.String("displayName", pol.Spec.DisplayName))
			_, err = c.client.AuthV1().Policies().Create(pol)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func (c *Controller) createAdmin(tenantID string) error {
	log.Info("Handle create admin for tenant", log.String("tenantID", tenantID))
	tenantUserSelector := fields.AndSelectors(
		fields.OneTermEqualSelector("spec.tenantID", tenantID),
		fields.OneTermEqualSelector("spec.username", c.username))

	localIdentityList, err := c.client.AuthV1().LocalIdentities().List(metav1.ListOptions{FieldSelector: tenantUserSelector.String()})
	if err != nil {
		return err
	}

	if len(localIdentityList.Items) != 0 {
		return nil
	}

	_, err = c.client.AuthV1().LocalIdentities().Create(&v1.LocalIdentity{
		Spec: v1.LocalIdentitySpec{
			HashedPassword: base64.StdEncoding.EncodeToString([]byte(c.password)),
			TenantID:       tenantID,
			Username:       c.username,
			DisplayName:    "Administrator",
			Extra: map[string]string{
				"platformadmin": "true",
			},
		},
	})
	if err != nil {
		log.Error("Failed to create the default admin identity for tenant", log.String("tenant", tenantID), log.Err(err))
		return err
	}
	return nil
}
