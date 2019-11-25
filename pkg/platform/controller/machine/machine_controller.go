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
	"fmt"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	"sync"
	"time"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	"tkestack.io/tke/api/platform/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/controller/machine/deletion"
	"tkestack.io/tke/pkg/platform/provider"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	machineClientRetryCount    = 5
	machineClientRetryInterval = 5 * time.Second

	reasonFailedInit = "FailedInit"
)

// Controller is responsible for performing actions dependent upon a machine phase.
type Controller struct {
	client           clientset.Interface
	cache            *machineCache
	health           *machineHealth
	queue            workqueue.RateLimitingInterface
	lister           platformv1lister.MachineLister
	listerSynced     cache.InformerSynced
	stopCh           <-chan struct{}
	machineProviders *sync.Map

	deleter deletion.MachineDeleterInterface
}

// obj could be an *v1.Machine, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(oldMachine *v1.Machine, newMachine *v1.Machine) bool {
	if oldMachine.UID != newMachine.UID {
		return true
	}

	if !reflect.DeepEqual(oldMachine.Spec, newMachine.Spec) {
		return true
	}

	if !reflect.DeepEqual(oldMachine.Status, newMachine.Status) {
		return true
	}

	return false
}

// NewController creates a new Controller object.
func NewController(
	client clientset.Interface,
	machineInformer platformv1informer.MachineInformer,
	resyncPeriod time.Duration,
	finalizerToken v1.FinalizerName,
	machineProviders *sync.Map) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:           client,
		cache:            &machineCache{m: make(map[string]*cachedMachine)},
		health:           &machineHealth{m: make(map[string]*v1.Machine)},
		queue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "machine"),
		machineProviders: machineProviders,
		deleter: deletion.NewMachineDeleter(client.PlatformV1().Machines(),
			client.PlatformV1(),
			finalizerToken,
			true),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("machine_controller", client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the namespace informer event handlers
	machineInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldMachine, ok1 := oldObj.(*v1.Machine)
				curMachine, ok2 := newObj.(*v1.Machine)
				if ok1 && ok2 && controller.needsUpdate(oldMachine, curMachine) {
					controller.enqueue(newObj)
				} else {
					log.Debug("Update new machine not to add", log.String("machineName", curMachine.Name), log.String("resourceversion", curMachine.ResourceVersion), log.String("old-resourceversion", oldMachine.ResourceVersion), log.String("cur-resourceversion", curMachine.ResourceVersion))
				}
			},
		},
		resyncPeriod,
	)
	controller.lister = machineInformer.Lister()
	controller.listerSynced = machineInformer.Informer().HasSynced

	return controller
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

	c.stopCh = stopCh

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
	startTime := time.Now()
	var cachedMachine *cachedMachine
	defer func() {
		log.Info("Finished syncing machine", log.String("machineName", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// machine holds the latest machine info from apiserver
	machine, err := c.lister.Get(name)

	switch {
	case apierrors.IsNotFound(err):
		log.Info("Machine has been deleted. Attempting to cleanup resources", log.String("machineName", key))
		err = c.processMachineDeletion(key, machine)
	case err != nil:
		log.Warn("Unable to retrieve machine from store", log.String("machineName", key), log.Err(err))
	default:
		if (machine.Status.Phase == v1.MachineRunning) || (machine.Status.Phase == v1.MachineFailed) || (machine.Status.Phase == v1.MachineInitializing) {
			cachedMachine = c.cache.getOrCreate(key)
			err = c.processMachineUpdate(cachedMachine, machine, key)
		} else if machine.Status.Phase == v1.MachineTerminating {
			log.Info("Machine has been terminated. Attempting to cleanup resources", log.String("machineName", key))
			_ = c.processMachineDeletion(key, machine)
			err = c.deleter.Delete(key)
		} else {
			log.Debug(fmt.Sprintf("Machine %s status is %s, not to process", key, machine.Status.Phase), log.String("machineName", key))
		}
	}
	return err
}

func (c *Controller) processMachineUpdate(cachedMachine *cachedMachine, machine *v1.Machine, key string) error {
	if cachedMachine.state != nil {
		if cachedMachine.state.UID != machine.UID {
			err := c.processMachineDelete(key, machine)
			if err != nil {
				return err
			}
		}
	}

	// start update machine if needed
	err := c.handlePhase(key, cachedMachine, machine)
	if err != nil {
		return err
	}

	cachedMachine.state = machine
	// Always update the cache upon success.
	c.cache.set(key, cachedMachine)

	return nil
}

func (c *Controller) processMachineDeletion(key string, machine *v1.Machine) error {
	return c.processMachineDelete(key, machine)
}

func (c *Controller) processMachineDelete(key string, machine *v1.Machine) error {
	log.Info("Machine will be dropped", log.String("machineName", key))

	if c.cache.Exist(key) {
		log.Info("Delete the machine cache", log.String("machineName", key))
		c.cache.delete(key)
	}

	if c.health.Exist(key) {
		log.Info("Delete the machine health cache", log.String("machineName", key))
		c.health.Del(key)
	}

	return nil
}

func (c *Controller) handlePhase(key string, cachedMachine *cachedMachine, machine *v1.Machine) error {
	var err error
	switch machine.Status.Phase {
	case v1.MachineInitializing:
		err = c.doInitializing(machine)
	case v1.MachineRunning, v1.MachineFailed:
		err = c.startMachineHealthCheck(key, machine)
	default:
		err = errors.Errorf("no handler for %q", machine.Status.Phase)
	}

	return err
}

func (c *Controller) addOrUpdateCondition(machine *v1.Machine, newCondition v1.MachineCondition) {
	var conditions []v1.MachineCondition
	exist := false
	for _, condition := range machine.Status.Conditions {
		if condition.Type == newCondition.Type {
			exist = true
			if newCondition.Status != condition.Status {
				condition.Status = newCondition.Status
			}
			if newCondition.Message != condition.Message {
				condition.Message = newCondition.Message
			}
			if newCondition.Reason != condition.Reason {
				condition.Reason = newCondition.Reason
			}
			if !newCondition.LastProbeTime.IsZero() && newCondition.LastProbeTime != condition.LastProbeTime {
				condition.LastProbeTime = newCondition.LastProbeTime
			}
			if !newCondition.LastTransitionTime.IsZero() && newCondition.LastTransitionTime != condition.LastTransitionTime {
				condition.LastTransitionTime = newCondition.LastTransitionTime
			}
		}
		conditions = append(conditions, condition)
	}
	if !exist {
		if newCondition.LastProbeTime.IsZero() {
			newCondition.LastProbeTime = metav1.Now()
		}
		if newCondition.LastTransitionTime.IsZero() {
			newCondition.LastTransitionTime = metav1.Now()
		}
		conditions = append(conditions, newCondition)
	}
	machine.Status.Conditions = conditions
}

func (c *Controller) persistUpdate(machine *v1.Machine) error {
	var err error
	for i := 0; i < machineClientRetryCount; i++ {
		_, err = c.client.PlatformV1().Machines().UpdateStatus(machine)
		if err == nil {
			return nil
		}
		// if the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already
		if apierrors.IsNotFound(err) {
			log.Info("Not persisting update to machine set that no longer exists", log.String("machineName", machine.Name), log.Err(err))
			return nil
		}
		if apierrors.IsConflict(err) {
			return fmt.Errorf("not persisting update to machine '%s' that has been changed since we received it: %v", machine.Name, err)
		}
		log.Warn("Failed to persist updated status of machine", log.String("machineName", machine.Name), log.Err(err))
		time.Sleep(machineClientRetryInterval)
	}

	return err
}

func (c *Controller) doInitializing(machine *v1.Machine) error {
	newMachine, err := c.onInitialize(machine)
	if err != nil {
		machine.Status.Message = err.Error()
		machine.Status.Reason = reasonFailedInit
		_, _ = c.client.PlatformV1().Machines().Update(machine)
		return err
	}
	condition := newMachine.Status.Conditions[len(newMachine.Status.Conditions)-1]
	if condition.Status == v1.ConditionFalse { // means current condition run into error
		newMachine.Status.Message = condition.Message
		newMachine.Status.Reason = condition.Reason
		_, _ = c.client.PlatformV1().Machines().Update(newMachine)
		return fmt.Errorf("OnInitialize.%s [Failed] reason: %s message: %s",
			condition.Type, condition.Reason, condition.Message)
	}
	newMachine.Status.Message = ""
	newMachine.Status.Reason = ""
	_, err = c.client.PlatformV1().Machines().Update(newMachine)
	return err
}

func (c *Controller) onInitialize(machine *v1.Machine) (*v1.Machine, error) {
	machineProvider, err := provider.LoadMachineProvider(c.machineProviders, string(machine.Spec.Type))
	if err != nil {
		return nil, err
	}

	cluster, err := c.client.PlatformV1().Clusters().Get(machine.Spec.ClusterName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "get cluster error")
	}
	clusterCredential, err := util.ClusterCredentialV1(c.client.PlatformV1(), cluster.Name)
	if err != nil {
		return nil, errors.Wrap(err, "get cluster credential error")
	}

	resp, err := machineProvider.OnInitialize(*machine, *cluster, *clusterCredential)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
