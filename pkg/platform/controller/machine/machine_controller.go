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

package machine

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"golang.org/x/time/rate"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	machineconfig "tkestack.io/tke/pkg/platform/controller/machine/config"
	"tkestack.io/tke/pkg/platform/controller/machine/deletion"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	machineprovider "tkestack.io/tke/pkg/platform/provider/machine"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	conditionTypeHealthCheck = "HealthCheck"
	failedHealthCheckReason  = "FailedHealthCheck"

	resyncInternal = 1 * time.Minute
)

// Controller is responsible for performing actions dependent upon a machine phase.
type Controller struct {
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.MachineLister
	listerSynced cache.InformerSynced

	log            log.Logger
	platformClient platformversionedclient.PlatformV1Interface
	deleter        deletion.MachineDeleterInterface
}

// NewController creates a new Controller object.
func NewController(
	platformclient platformversionedclient.PlatformV1Interface,
	machineInformer platformv1informer.MachineInformer,
	configuration machineconfig.MachineControllerConfiguration,
	finalizerToken platformv1.FinalizerName) *Controller {
	rateLimit := workqueue.NewMaxOfRateLimiter(
		workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(configuration.BucketRateLimiterLimit), configuration.BucketRateLimiterBurst)},
	)
	c := &Controller{
		queue: workqueue.NewNamedRateLimitingQueue(rateLimit, "machine"),

		log:            log.WithName("MachineController"),
		platformClient: platformclient,
		deleter:        deletion.NewMachineDeleter(platformclient.Machines(), platformclient, finalizerToken, true),
	}

	if platformclient != nil && platformclient.RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("machine_controller", platformclient.RESTClient().GetRateLimiter())
	}

	machineInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.addMachine,
			UpdateFunc: c.updateMachine,
		},
		configuration.MachineSyncPeriod,
	)
	c.lister = machineInformer.Lister()
	c.listerSynced = machineInformer.Informer().HasSynced

	return c
}

func (c *Controller) addMachine(obj interface{}) {
	machine := obj.(*platformv1.Machine)
	c.log.Info("Adding machine", "machine", machine.Name)
	c.enqueue(machine)
}

func (c *Controller) updateMachine(old, obj interface{}) {
	oldMachine := old.(*platformv1.Machine)
	machine := obj.(*platformv1.Machine)

	controllerNeedUpddateResult := c.needsUpdate(oldMachine, machine)
	var providerNeedUpddateResult bool
	provider, _ := machineprovider.GetProvider(machine.Spec.Type)
	if provider != nil {
		providerNeedUpddateResult = provider.NeedUpdate(oldMachine, machine)
	}
	if !(controllerNeedUpddateResult || providerNeedUpddateResult) {
		return
	}
	c.log.Info("Updating machine", "machine", machine.Name)
	c.enqueue(machine)
}

func (c *Controller) needsUpdate(old *platformv1.Machine, new *platformv1.Machine) bool {
	healthCondition := new.GetCondition(conditionTypeHealthCheck)
	if !reflect.DeepEqual(old.Spec, new.Spec) {
		return true

	}
	if !reflect.DeepEqual(old.ObjectMeta.Labels, new.ObjectMeta.Labels) {
		return true
	}
	if !reflect.DeepEqual(old.ObjectMeta.Annotations, new.ObjectMeta.Annotations) {
		return true
	}
	if old.Status.Phase != new.Status.Phase {
		return true

	}
	if new.Status.Phase == platformv1.MachineInitializing {
		// if ResourceVersion is equal, it's an resync envent, should return true.
		if old.ResourceVersion == new.ResourceVersion {
			return true
		}
		if len(new.Status.Conditions) == 0 {
			return true
		}
		if new.Status.Conditions[len(new.Status.Conditions)-1].Status == platformv1.ConditionUnknown {
			return true
		}
		// if user set last condition false block procesee until resync envent
		if new.Status.Conditions[len(new.Status.Conditions)-1].Status == platformv1.ConditionFalse {
			return false
		}
	}
	// if last health check is not long enough， return false
	if healthCondition != nil &&
		time.Since(healthCondition.LastProbeTime.Time) < resyncInternal {
		return false
	}
	return true
}

func (c *Controller) enqueue(obj *platformv1.Machine) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.Add(key)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting machine controller")
	defer log.Info("Shutting down machine controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		return fmt.Errorf("failed to wait for machine caches to sync")
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
	return nil
}

// worker processes the queue of persistent event objects.
// Each machine can be in the queue at most once.
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

	err := c.syncMachine(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing machine %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncMachine will sync the Machine with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncMachine(key string) error {
	ctx := c.log.WithValues("machine", key).WithContext(context.TODO())

	startTime := time.Now()
	defer func() {
		log.FromContext(ctx).Info("Finished syncing machine", "processTime", time.Since(startTime).String())
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	machine, err := c.lister.Get(name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.FromContext(ctx).Info("Machine has been deleted")
			return nil
		}

		utilruntime.HandleError(fmt.Errorf("unable to retrieve machine %v from store: %v", key, err))
		return err
	}

	ctx = log.FromContext(ctx).WithValues("cluster", machine.Spec.ClusterName).WithContext(ctx)

	return c.reconcile(ctx, key, machine)
}

func (c *Controller) reconcile(ctx context.Context, key string, machine *platformv1.Machine) error {
	var err error
	switch machine.Status.Phase {
	case platformv1.MachineInitializing:
		err = c.onCreate(ctx, machine)
	case platformv1.MachineRunning, platformv1.MachineFailed, platformv1.MachineUpgrading:
		err = c.onUpdate(ctx, machine)
	case platformv1.MachineTerminating:
		log.FromContext(ctx).Info("Machine has been terminated. Attempting to cleanup resources")
		err = c.deleter.Delete(ctx, key)
		if err == nil {
			log.FromContext(ctx).Info("Machine has been successfully deleted")
		}
	default:
		log.FromContext(ctx).Info("unknown machine phase", "status.phase", machine.Status.Phase)
	}

	return err
}

func (c *Controller) onCreate(ctx context.Context, machine *platformv1.Machine) error {
	provider, err := machineprovider.GetProvider(machine.Spec.Type)
	if err != nil {
		return err
	}
	cluster, err := clusterprovider.GetV1ClusterByName(ctx, c.platformClient, machine.Spec.ClusterName, clusterprovider.AdminUsername)
	if err != nil {
		return err
	}

	for machine.Status.Phase == platformv1.MachineInitializing {
		err = provider.OnCreate(ctx, machine, cluster)
		if err != nil {
			// Update status, ignore failure
			_, _ = c.platformClient.Machines().Update(ctx, machine, metav1.UpdateOptions{})
			return err
		}
		machine, err = c.platformClient.Machines().Update(ctx, machine, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return err
}

func (c *Controller) onUpdate(ctx context.Context, machine *platformv1.Machine) error {
	provider, err := machineprovider.GetProvider(machine.Spec.Type)
	if err != nil {
		return err
	}

	cluster, err := clusterprovider.GetV1ClusterByName(ctx, c.platformClient, machine.Spec.ClusterName, clusterprovider.AdminUsername)
	if err != nil {
		return err
	}

	err = provider.OnUpdate(ctx, machine, cluster)
	machine = c.checkHealth(ctx, machine)
	if err != nil {
		// Update status, ignore failure
		_, _ = c.platformClient.Machines().UpdateStatus(ctx, machine, metav1.UpdateOptions{})
		return err
	}
	_, err = c.platformClient.Machines().Update(ctx, machine, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) checkHealth(ctx context.Context, machine *platformv1.Machine) *platformv1.Machine {
	if !(machine.Status.Phase == platformv1.MachineRunning ||
		machine.Status.Phase == platformv1.MachineFailed) {
		return machine
	}

	healthCheckCondition := platformv1.MachineCondition{
		Type:   conditionTypeHealthCheck,
		Status: platformv1.ConditionFalse,
	}

	clientset, err := util.BuildExternalClientSetWithName(ctx, c.platformClient, machine.Spec.ClusterName)
	if err != nil {
		machine.Status.Phase = platformv1.MachineFailed

		healthCheckCondition.Reason = failedHealthCheckReason
		healthCheckCondition.Message = err.Error()
	} else {
		_, err = apiclient.GetNodeByMachineIP(ctx, clientset, machine.Spec.IP)
		if err != nil {
			machine.Status.Phase = platformv1.MachineFailed

			healthCheckCondition.Reason = failedHealthCheckReason
			healthCheckCondition.Message = err.Error()
		} else {
			machine.Status.Phase = platformv1.MachineRunning

			healthCheckCondition.Status = platformv1.ConditionTrue
		}
	}

	machine.SetCondition(healthCheckCondition)

	log.FromContext(ctx).Info("Update machine health status", "phase", machine.Status.Phase)

	return machine
}
