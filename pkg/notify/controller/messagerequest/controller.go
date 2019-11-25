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

package messagerequest

import (
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	notifyv1informer "tkestack.io/tke/api/client/informers/externalversions/notify/v1"
	notifyv1lister "tkestack.io/tke/api/client/listers/notify/v1"
	v1 "tkestack.io/tke/api/notify/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/notify/controller/messagerequest/smtp"
	"tkestack.io/tke/pkg/notify/controller/messagerequest/tencentcloudsms"
	"tkestack.io/tke/pkg/notify/controller/messagerequest/util"
	"tkestack.io/tke/pkg/notify/controller/messagerequest/wechat"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second
)

const (
	controllerName = "message-request"
)

// Controller is responsible for performing actions dependent upon a message request controller phase.
type Controller struct {
	client       clientset.Interface
	cache        *messageRequestCache
	queue        workqueue.RateLimitingInterface
	lister       notifyv1lister.MessageRequestLister
	listerSynced cache.InformerSynced
	stopCh       <-chan struct{}
}

// NewController creates a new Controller object.
func NewController(client clientset.Interface, informer notifyv1informer.MessageRequestInformer, resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client: client,
		cache:  &messageRequestCache{messageRequestMap: make(map[string]*cachedMessageRequest)},
		queue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage(controllerName, client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the tapp controller informer event handlers
	informer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueMessageRequest,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldMessageRequest, ok1 := oldObj.(*v1.MessageRequest)
				curMessageRequest, ok2 := newObj.(*v1.MessageRequest)
				if ok1 && ok2 && controller.needsUpdate(oldMessageRequest, curMessageRequest) {
					controller.enqueueMessageRequest(newObj)
				}
			},
			DeleteFunc: controller.enqueueMessageRequest,
		},
		resyncPeriod,
	)
	controller.lister = informer.Lister()
	controller.listerSynced = informer.Informer().HasSynced

	return controller
}

// obj could be an *v1.MessageRequest, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueueMessageRequest(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(old *v1.MessageRequest, new *v1.MessageRequest) bool {
	return !reflect.DeepEqual(old, new)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting message request controller")
	defer log.Info("Shutting down message request controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		log.Error("Failed to wait for message request caches to sync")
		return
	}

	c.stopCh = stopCh

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of message request objects.
// Each message request can be in the queue at most once.
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

	err := c.syncMessageRequest(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing message request controller %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncMessageRequest will sync the message request with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncMessageRequest(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing message request controller", log.String("messageRequestName", key), log.Duration("processTime", time.Since(startTime)))
	}()

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// messageRequest holds the latest messageRequest info from apiserver
	messageRequest, err := c.lister.MessageRequests(ns).Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Info("Message request has been deleted. Attempting to cleanup resources", log.String("messageRequestName", key))
		err = c.processMessageRequestDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve message request from store", log.String("messageRequestName", key), log.Err(err))
	default:
		cachedMessageRequest := c.cache.getOrCreate(key)
		err = c.processMessageRequestUpdate(cachedMessageRequest, messageRequest, key)
	}
	return err
}

func (c *Controller) processMessageRequestDeletion(key string) error {
	cachedMessageRequest, ok := c.cache.get(key)
	if !ok {
		log.Error("Message request controller not in cache even though the watcher thought it was. Ignoring the deletion", log.String("tappControllerName", key))
		return nil
	}
	return c.processMessageRequestDelete(cachedMessageRequest, key)
}

func (c *Controller) processMessageRequestDelete(cachedMessageRequest *cachedMessageRequest, key string) error {
	log.Info("Message request controller will be dropped", log.String("messageRequestName", key))

	if c.cache.Exist(key) {
		log.Info("Delete the message request controller cache", log.String("messageRequestName", key))
		c.cache.delete(key)
	}

	return nil
}

func (c *Controller) processMessageRequestUpdate(cachedMessageRequest *cachedMessageRequest, messageRequest *v1.MessageRequest, key string) error {
	if cachedMessageRequest.state != nil {
		// exist and the cluster name changed
		if cachedMessageRequest.state.UID != messageRequest.UID {
			if err := c.processMessageRequestDelete(cachedMessageRequest, key); err != nil {
				return err
			}
		}
	}
	err := c.createMessageRequestIfNeeded(key, cachedMessageRequest, messageRequest)
	if err != nil {
		return err
	}

	cachedMessageRequest.state = messageRequest
	// Always update the cache upon success.
	c.cache.set(key, cachedMessageRequest)
	return nil
}

func (c *Controller) createMessageRequestIfNeeded(key string, cachedMessageRequest *cachedMessageRequest, messageRequest *v1.MessageRequest) error {
	switch messageRequest.Status.Phase {
	case v1.MessageRequestPending:
		messageRequest = messageRequest.DeepCopy()
		messageRequest.Status.Phase = v1.MessageRequestSending
		messageRequest.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdate(messageRequest)
	case v1.MessageRequestSending:
		if cachedMessageRequest.state != nil && cachedMessageRequest.state.Status.Phase == v1.MessageRequestPending {
			sentMessages, failedReceiverErrors := c.sendMessage(messageRequest)
			if len(sentMessages) > 0 {
				c.archiveMessage(messageRequest, sentMessages)
			}
			messageRequest = messageRequest.DeepCopy()
			messageRequest.Status.LastTransitionTime = metav1.Now()
			if len(failedReceiverErrors) == 0 {
				messageRequest.Status.Phase = v1.MessageRequestSent
			} else {
				if len(sentMessages) == 0 {
					messageRequest.Status.Phase = v1.MessageRequestFailed
				} else {
					messageRequest.Status.Phase = v1.MessageRequestPartialFailure
				}
				messageRequest.Status.Errors = failedReceiverErrors
			}
			return c.persistUpdate(messageRequest)
		}
	}
	return nil
}

func (c *Controller) persistUpdate(messageRequest *v1.MessageRequest) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = c.client.NotifyV1().MessageRequests(messageRequest.ObjectMeta.Namespace).UpdateStatus(messageRequest)
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to message request that no longer exists", log.String("messageRequestName", messageRequest.ObjectMeta.Name), log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to message request '%s' that has been changed since we received it: %v", messageRequest.ObjectMeta.Name, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of message request '%s/%s'", messageRequest.ObjectMeta.Name, messageRequest.Status.Phase), log.String("messageRequestName", messageRequest.ObjectMeta.Name), log.Err(err))
		time.Sleep(clientRetryInterval)
	}

	return err
}

type sentMessage struct {
	receiverName    string
	receiverChannel v1.ReceiverChannel
	identity        string
	username        string
	header          string
	body            string
	messageID       string
}

func (c *Controller) sendMessage(messageRequest *v1.MessageRequest) (sentMessages []sentMessage, failedReceiverErrors map[string]string) {
	failedReceiverErrors = make(map[string]string)
	receivers := sets.NewString()
	for _, receiverGroupName := range messageRequest.Spec.ReceiverGroups {
		if receiverGroupName == "" {
			continue
		}
		if receiverGroup, err := c.client.NotifyV1().ReceiverGroups().Get(receiverGroupName, metav1.GetOptions{}); err != nil {
			log.Error("Failed to retrieve the specify receiver group", log.String("receiverGroupName", receiverGroupName), log.Err(err))
		} else {
			if len(receiverGroup.Spec.Receivers) != 0 {
				receivers.Insert(receiverGroup.Spec.Receivers...)
			}
		}
	}
	for _, receiver := range messageRequest.Spec.Receivers {
		if receiver != "" {
			receivers.Insert(receiver)
		}
	}

	if receivers.Len() == 0 {
		return
	}
	selfdefineChannelList, err := c.client.NotifyV1().Channels().List(metav1.ListOptions{LabelSelector: "selfdefine=true"})
	if err != nil {
		log.Error("Failed to get selfdefineChannelList", log.Err(err))
		return
	}
	if len(selfdefineChannelList.Items) != 0 {
		selfdefineURL, ok := selfdefineChannelList.Items[0].ObjectMeta.Annotations["selfdefineURL"]
		if !ok {
			log.Error("The selfdefineURL does not exist", log.Err(err))
			return
		}
		util.SelfdefineURL = selfdefineURL
	} else {
		log.Warn("The selfdefineChannel does not exist")
		util.SelfdefineURL = ""
	}

	channel, err := c.client.NotifyV1().Channels().Get(messageRequest.ObjectMeta.Namespace, metav1.GetOptions{})
	if err != nil {
		log.Error("Failed to get channel object", log.String("channelName", messageRequest.ObjectMeta.Namespace), log.Err(err))
		for _, receiverName := range receivers.List() {
			failedReceiverErrors[receiverName] = "Failed to get notification channel information"
		}
		return
	}
	template, err := c.client.NotifyV1().Templates(channel.ObjectMeta.Name).Get(messageRequest.Spec.TemplateName, metav1.GetOptions{})
	if err != nil {
		log.Error("Failed to get template object", log.String("channelName", channel.ObjectMeta.Name), log.String("templateName", messageRequest.Spec.TemplateName), log.Err(err))
		for _, receiverName := range receivers.List() {
			failedReceiverErrors[receiverName] = "Failed to get notification template information"
		}
		return
	}
	for _, receiverName := range receivers.List() {
		receiver, err := c.client.NotifyV1().Receivers().Get(receiverName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				failedReceiverErrors[receiverName] = "The specified notification recipient does not exist"
			} else {
				failedReceiverErrors[receiverName] = "Failed to get notification recipient information"
			}
			continue
		}
		templateCount := 0
		if template.Spec.TencentCloudSMS != nil {
			templateCount++
			if channel.Spec.TencentCloudSMS == nil {
				failedReceiverErrors[receiverName] = "The notification sending template is not configured with the corresponding tencent cloud account"
				continue
			}
			mobile, ok := receiver.Spec.Identities[v1.ReceiverChannelMobile]
			if !ok {
				failedReceiverErrors[receiverName] = "The notification recipient did not configure the mobile"
				continue
			}
			messageID, body, err := tencentcloudsms.Send(channel.Spec.TencentCloudSMS, template.Spec.TencentCloudSMS, mobile, messageRequest.Spec.Variables)
			if err != nil {
				failedReceiverErrors[receiverName] = err.Error()
				continue
			}
			sentMessages = append(sentMessages, sentMessage{
				receiverName:    receiverName,
				receiverChannel: v1.ReceiverChannelMobile,
				identity:        mobile,
				username:        receiver.Spec.Username,
				body:            body,
				messageID:       messageID,
			})
		}
		if template.Spec.Wechat != nil {
			templateCount++
			if channel.Spec.Wechat == nil {
				failedReceiverErrors[receiverName] = "The notification sending template is not configured with the corresponding Wechat account"
				continue
			}
			openID, ok := receiver.Spec.Identities[v1.ReceiverChannelWechatOpenID]
			if !ok {
				failedReceiverErrors[receiverName] = "The notification recipient did not configure the Wechat openid"
				continue
			}
			messageID, body, err := wechat.Send(channel.Spec.Wechat, template.Spec.Wechat, openID, messageRequest.Spec.Variables)
			if err != nil {
				failedReceiverErrors[receiverName] = err.Error()
				continue
			}
			sentMessages = append(sentMessages, sentMessage{
				receiverName:    receiverName,
				username:        receiver.Spec.Username,
				receiverChannel: v1.ReceiverChannelWechatOpenID,
				identity:        openID,
				body:            body,
				messageID:       messageID,
			})
		}
		if template.Spec.Text != nil {
			templateCount++
			if channel.Spec.SMTP == nil {
				failedReceiverErrors[receiverName] = "The notification sending template is not configured with the corresponding smtp server"
				continue
			}
			email, ok := receiver.Spec.Identities[v1.ReceiverChannelEmail]
			if !ok {
				failedReceiverErrors[receiverName] = "The notification recipient did not configure the email"
				continue
			}
			header, body, err := smtp.Send(channel.Spec.SMTP, template.Spec.Text, email, messageRequest.Spec.Variables)
			if err != nil {
				failedReceiverErrors[receiverName] = err.Error()
				continue
			}
			sentMessages = append(sentMessages, sentMessage{
				receiverName:    receiverName,
				username:        receiver.Spec.Username,
				receiverChannel: v1.ReceiverChannelEmail,
				identity:        email,
				header:          header,
				body:            body,
			})
		}
		if templateCount == 0 {
			failedReceiverErrors[receiverName] = "Notification sending template is not configured"
		}
	}
	return
}

func (c *Controller) archiveMessage(messageRequest *v1.MessageRequest, sentMessages []sentMessage) {
	for _, sentMessage := range sentMessages {
		message := &v1.Message{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s", messageRequest.ObjectMeta.Name, sentMessage.receiverName),
			},
			Spec: v1.MessageSpec{
				TenantID:         messageRequest.Spec.TenantID,
				ReceiverName:     sentMessage.receiverName,
				ReceiverChannel:  sentMessage.receiverChannel,
				Identity:         sentMessage.identity,
				Username:         sentMessage.username,
				Header:           sentMessage.header,
				Body:             sentMessage.body,
				ChannelMessageID: sentMessage.messageID,
			},
			Status: v1.MessageStatus{
				Phase: v1.MessageUnread,
			},
		}
		if _, err := c.client.NotifyV1().Messages().Create(message); err != nil {
			log.Error("Failed to create message object", log.String("messageRequestName", messageRequest.ObjectMeta.Name), log.String("receiverName", sentMessage.receiverName), log.Err(err))
		}
	}
}
