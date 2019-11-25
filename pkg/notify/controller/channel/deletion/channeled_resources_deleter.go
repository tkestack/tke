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

package deletion

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/notify/v1"
	"tkestack.io/tke/api/notify/v1"
	"tkestack.io/tke/pkg/util/log"
)

// ChanneledResourcesDeleterInterface to delete a channel with all resources in
// it.
type ChanneledResourcesDeleterInterface interface {
	Delete(channelName string) error
}

// NewChanneledResourcesDeleter to create the channeledResourcesDeleter that
// implement the ChanneledResourcesDeleterInterface by given client and
// configure.
func NewChanneledResourcesDeleter(channelClient v1clientset.ChannelInterface,
	notifyClient v1clientset.NotifyV1Interface,
	finalizerToken v1.FinalizerName,
	deleteChannelWhenDone bool) ChanneledResourcesDeleterInterface {
	d := &channeledResourcesDeleter{
		channelClient:         channelClient,
		notifyClient:          notifyClient,
		finalizerToken:        finalizerToken,
		deleteChannelWhenDone: deleteChannelWhenDone,
	}
	return d
}

var _ ChanneledResourcesDeleterInterface = &channeledResourcesDeleter{}

// channeledResourcesDeleter is used to delete all resources in a given channel.
type channeledResourcesDeleter struct {
	// Client to manipulate the channel.
	channelClient v1clientset.ChannelInterface
	notifyClient  v1clientset.NotifyV1Interface
	// The finalizer token that should be removed from the channel
	// when all resources in that channel have been deleted.
	finalizerToken v1.FinalizerName
	// Also delete the channel when all resources in the channel have been deleted.
	deleteChannelWhenDone bool
}

// Delete deletes all resources in the given channel.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   channel (does nothing if deletion timestamp is missing).
// * Verifies that the channel is in the "terminating" phase
//   (updates the channel phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given channel.
// * Deletes the channel if deleteChannelWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *channeledResourcesDeleter) Delete(channelName string) error {
	// Multiple controllers may edit a channel during termination
	// first get the latest state of the channel before proceeding
	// if the channel was deleted already, don't do anything
	channel, err := d.channelClient.Get(channelName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if channel.DeletionTimestamp == nil {
		return nil
	}

	log.Infof("channel controller - syncChannel - channel: %s, finalizerToken: %s", channel.Name, d.finalizerToken)

	// ensure that the status is up to date on the channel
	// if we get a not found error, we assume the channel is truly gone
	channel, err = d.retryOnConflictError(channel, d.updateChannelStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the channel asserts that channel is no longer deleting..
	if channel.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the channel if it is already finalized.
	if d.deleteChannelWhenDone && finalized(channel) {
		return d.deleteChannel(channel)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(channel)
	if err != nil {
		return err
	}

	// we have removed content, so mark it finalized by us
	channel, err = d.retryOnConflictError(channel, d.finalizeChannel)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do channel deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check if we can delete now.
	if d.deleteChannelWhenDone && finalized(channel) {
		return d.deleteChannel(channel)
	}
	return nil
}

// Deletes the given channel.
func (d *channeledResourcesDeleter) deleteChannel(channel *v1.Channel) error {
	var opts *metav1.DeleteOptions
	uid := channel.UID
	if len(uid) > 0 {
		opts = &metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.channelClient.Delete(channel.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// updateChannelFunc is a function that makes an update to a channel
type updateChannelFunc func(channel *v1.Channel) (*v1.Channel, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in client code
func (d *channeledResourcesDeleter) retryOnConflictError(channel *v1.Channel, fn updateChannelFunc) (result *v1.Channel, err error) {
	latestChannel := channel
	for {
		result, err = fn(latestChannel)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevChannel := latestChannel
		latestChannel, err = d.channelClient.Get(latestChannel.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevChannel.UID != latestChannel.UID {
			return nil, fmt.Errorf("channel uid has changed across retries")
		}
	}
}

// updateChannelStatusFunc will verify that the status of the channel is correct
func (d *channeledResourcesDeleter) updateChannelStatusFunc(channel *v1.Channel) (*v1.Channel, error) {
	if channel.DeletionTimestamp.IsZero() || channel.Status.Phase == v1.ChannelTerminating {
		return channel, nil
	}
	newChannel := v1.Channel{}
	newChannel.ObjectMeta = channel.ObjectMeta
	newChannel.Status = channel.Status
	newChannel.Status.Phase = v1.ChannelTerminating
	return d.channelClient.UpdateStatus(&newChannel)
}

// finalized returns true if the channel.Spec.Finalizers is an empty list
func finalized(channel *v1.Channel) bool {
	return len(channel.Spec.Finalizers) == 0
}

// finalizeChannel removes the specified finalizerToken and finalizes the channel
func (d *channeledResourcesDeleter) finalizeChannel(channel *v1.Channel) (*v1.Channel, error) {
	channelFinalize := v1.Channel{}
	channelFinalize.ObjectMeta = channel.ObjectMeta
	channelFinalize.Spec = channel.Spec
	finalizerSet := sets.NewString()
	for i := range channel.Spec.Finalizers {
		if channel.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(channel.Spec.Finalizers[i]))
		}
	}
	channelFinalize.Spec.Finalizers = make([]v1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		channelFinalize.Spec.Finalizers = append(channelFinalize.Spec.Finalizers, v1.FinalizerName(value))
	}

	channel = &v1.Channel{}
	err := d.notifyClient.RESTClient().Put().
		Resource("channels").
		Name(channelFinalize.Name).
		SubResource("finalize").
		Body(&channelFinalize).
		Do().
		Into(channel)

	if err != nil {
		// it was removed already, so life is good
		if errors.IsNotFound(err) {
			return channel, nil
		}
	}
	return channel, err
}

type deleteResourceFunc func(deleter *channeledResourcesDeleter, channel *v1.Channel) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteTemplate,
	deleteMessageRequest,
}

// deleteAllContent will use the dynamic client to delete each resource identified in groupVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *channeledResourcesDeleter) deleteAllContent(channel *v1.Channel) error {
	log.Debug("Channel controller - deleteAllContent", log.String("channelName", channel.ObjectMeta.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(d, channel)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.Debug("Channel controller - deletedAllContent", log.String("channelName", channel.ObjectMeta.Name))
	return nil
}

func deleteTemplate(deleter *channeledResourcesDeleter, channel *v1.Channel) error {
	log.Debug("Channel controller - deleteTemplate", log.String("channelName", channel.ObjectMeta.Name))

	background := metav1.DeletePropagationBackground
	deleteOpt := &metav1.DeleteOptions{PropagationPolicy: &background}
	if err := deleter.notifyClient.Templates(channel.ObjectMeta.Name).DeleteCollection(deleteOpt, metav1.ListOptions{}); err != nil {
		log.Error("Channel controller - failed to delete template collections", log.String("channelName", channel.ObjectMeta.Name), log.Err(err))
		return err
	}

	return nil
}

func deleteMessageRequest(deleter *channeledResourcesDeleter, channel *v1.Channel) error {
	log.Debug("Channel controller - deleteMessageRequest", log.String("channelName", channel.ObjectMeta.Name))

	background := metav1.DeletePropagationBackground
	deleteOpt := &metav1.DeleteOptions{PropagationPolicy: &background}
	if err := deleter.notifyClient.Templates(channel.ObjectMeta.Name).DeleteCollection(deleteOpt, metav1.ListOptions{}); err != nil {
		log.Error("Channel controller - failed to delete message request collections", log.String("channelName", channel.ObjectMeta.Name), log.Err(err))
		return err
	}

	return nil
}
