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

package identityprovider

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	authv1 "tkestack.io/tke/api/auth/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	authversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	authv1informer "tkestack.io/tke/api/client/informers/externalversions/auth/v1"
	authv1lister "tkestack.io/tke/api/client/listers/auth/v1"
	registryv1 "tkestack.io/tke/api/registry/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	registryconfigv1 "tkestack.io/tke/pkg/registry/apis/config/v1"
	registrycontrollerconfig "tkestack.io/tke/pkg/registry/controller/config"
	"tkestack.io/tke/pkg/registry/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	controllerName = "identityprovider-controller"
)

// Controller is responsible for performing actions dependent upon a message request controller phase.
type Controller struct {
	client     clientset.Interface
	authClient authversionedclient.AuthV1Interface
	queue      workqueue.RateLimitingInterface

	identityProviderLister       authv1lister.IdentityProviderLister
	identityProviderListerSynced cache.InformerSynced

	stopCh                       <-chan struct{}
	registryDefaultConfiguration registrycontrollerconfig.RegistryDefaultConfiguration
	registryConfig               *registryconfigv1.RegistryConfiguration
	corednsClient                *util.CoreDNS
}

func NewController(authClient authversionedclient.AuthV1Interface,
	client clientset.Interface,
	identityProviderInformer authv1informer.IdentityProviderInformer,
	resyncPeriod time.Duration,
	registryDefaultConfiguration registrycontrollerconfig.RegistryDefaultConfiguration,
	registryConfig *registryconfigv1.RegistryConfiguration,
) *Controller {
	corednsClient, err := util.NewCoreDNS()
	if err != nil {
		log.Error("create coredns client failed", log.Err(err))
	}
	controller := Controller{
		client:                       client,
		authClient:                   authClient,
		queue:                        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		registryDefaultConfiguration: registryDefaultConfiguration,
		registryConfig:               registryConfig,
		corednsClient:                corednsClient,
	}

	if authClient != nil && authClient.RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("identityProvider_controller", authClient.RESTClient().GetRateLimiter())
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
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting identityprovider controller")
	defer log.Info("Shutting down identityprovider controller")

	if ok := cache.WaitForCacheSync(stopCh, c.identityProviderListerSynced); !ok {
		log.Error("Failed to wait for identyProvider caches to sync")
	}

	c.stopCh = stopCh
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
		log.Info("Init default registry setting for identityProvider", log.Any("name", idp.Name))
		err = c.processUpdateItem(context.Background(), idp, key)
		if c.corednsClient != nil {
			item := fmt.Sprintf("%s.%s", idp.Name, c.registryConfig.DomainSuffix)
			c.corednsClient.ParseCoreFile(item)
		}
	}
	return err
}

func (c *Controller) processUpdateItem(ctx context.Context, idp *authv1.IdentityProvider, key string) error {
	var errs []error
	if len(c.registryDefaultConfiguration.DefaultSystemChartGroups) == 0 {
		log.Info("DefaultSystemChartGroups is empty for identityProvider", log.Any("name", idp.Name))
		return nil
	}
	if c.client == nil {
		return nil
	}
	for _, cg := range c.registryDefaultConfiguration.DefaultSystemChartGroups {
		cgList, err := c.client.RegistryV1().ChartGroups().List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", idp.Spec.Name, cg),
		})
		if err != nil {
			log.Warn("identityprovider controller - list chartGroup failed",
				log.String("chartGroupName", cg),
				log.String("tenant", idp.Spec.Name),
				log.Err(err))
			errs = append(errs, err)
			continue
		}
		if len(cgList.Items) == 0 {
			_, err = c.client.RegistryV1().ChartGroups().Create(ctx, &registryv1.ChartGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name: cg,
				},
				Spec: registryv1.ChartGroupSpec{
					Name:        cg,
					DisplayName: cg,
					TenantID:    idp.Spec.Name,
					Visibility:  registryv1.VisibilityPublic,
					Type:        registryv1.RepoTypeSystem,
				}}, metav1.CreateOptions{})
			if err != nil {
				log.Warn("identityprovider controller - addChartGroup failed",
					log.String("chartGroupName", cg),
					log.String("tenant", idp.Spec.Name),
					log.Err(err))
				errs = append(errs, err)
			} else {
				log.Info("identityprovider controller - addChartGroup",
					log.String("chartGroupName", cg),
					log.String("tenant", idp.Spec.Name))
			}
		}
	}
	return utilerrors.NewAggregate(errs)
}
