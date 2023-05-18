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
	"context"
	"errors"
	"fmt"
	"reflect"

	"helm.sh/helm/v3/pkg/storage/driver"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	applicationv1 "tkestack.io/tke/api/application/v1"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	appconfig "tkestack.io/tke/pkg/application/config"
	"tkestack.io/tke/pkg/application/controller/app/action"
	"tkestack.io/tke/pkg/util/log"
)

// AppResourcesDeleterInterface to delete a app with all resources in
// it.
type AppResourcesDeleterInterface interface {
	Delete(ctx context.Context, namespace, appName string) error
}

// NewAppResourcesDeleter to create the appResourcesDeleter that
// implement the AppResourcesDeleterInterface by given applicationClient,
// applicationClient and configure.
func NewAppResourcesDeleter(
	applicationClient v1clientset.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	repo appconfig.RepoConfiguration,
	finalizerToken applicationv1.FinalizerName,
	deleteAppWhenDone bool) AppResourcesDeleterInterface {
	d := &applicationResourcesDeleter{
		applicationClient: applicationClient,
		platformClient:    platformClient,
		finalizerToken:    finalizerToken,
		repo:              repo,
		deleteAppWhenDone: deleteAppWhenDone,
	}
	return d
}

var _ AppResourcesDeleterInterface = &applicationResourcesDeleter{}

// applicationResourcesDeleter is used to delete all resources in a given app.
type applicationResourcesDeleter struct {
	// Client to manipulate the application.
	applicationClient v1clientset.ApplicationV1Interface
	// Client to manipulate the platform.
	platformClient platformversionedclient.PlatformV1Interface
	// The finalizer token that should be removed from the app
	// when all resources in that app have been deleted.
	finalizerToken applicationv1.FinalizerName
	// Also delete the app when all resources in the app have been deleted.
	deleteAppWhenDone bool
	// RepoConfiguration contains options to connect to a chart repo.
	repo appconfig.RepoConfiguration
}

// Delete deletes all resources in the given app.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   app (does nothing if deletion timestamp is missing).
// * Verifies that the app is in the "terminating" phase
//   (updates the app phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given app.
// * Deletes the app if deleteAppWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *applicationResourcesDeleter) Delete(ctx context.Context, namespace, appName string) error {
	// Multiple controllers may edit a app during termination
	// first get the latest state of the app before proceeding
	// if the app was deleted already, don't do anything
	app, err := d.applicationClient.Apps(namespace).Get(ctx, appName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if app.DeletionTimestamp == nil {
		return nil
	}

	log.Infof("app controller - syncApplication - app: %s, finalizerToken: %s", app.Name, d.finalizerToken)

	// ensure that the status is up to date on the app
	// if we get a not found error, we assume the app is truly gone
	app, err = d.retryOnConflictError(ctx, app, d.updateApplicationStatusFunc)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the app asserts that app is no longer deleting.
	if app.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the app if it is already finalized.
	if d.deleteAppWhenDone && finalized(app) {
		return d.deleteApplication(ctx, app)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(ctx, app)
	if err != nil {
		return err
	}

	// we have removed all content, so mark it finalized by us.
	app, err = d.retryOnConflictError(ctx, app, d.finalizeApplication)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do app deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if k8serrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check whether we can delete it now.
	if d.deleteAppWhenDone && finalized(app) {
		return d.deleteApplication(ctx, app)
	}
	return nil
}

// Deletes the given app.
func (d *applicationResourcesDeleter) deleteApplication(ctx context.Context, app *applicationv1.App) error {
	var opts metav1.DeleteOptions
	uid := app.UID
	if len(uid) > 0 {
		opts = metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.applicationClient.Apps(app.Namespace).Delete(ctx, app.Name, opts)
	if err != nil && !k8serrors.IsNotFound(err) {
		return err
	}
	return nil
}

// updateApplicationFunc is a function that makes an update to a app
type updateApplicationFunc func(ctx context.Context, app *applicationv1.App) (*applicationv1.App, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in authClient code
func (d *applicationResourcesDeleter) retryOnConflictError(ctx context.Context, app *applicationv1.App, fn updateApplicationFunc) (result *applicationv1.App, err error) {
	latestApplication := app
	for {
		result, err = fn(ctx, latestApplication)
		if err == nil {
			return result, nil
		}
		if !k8serrors.IsConflict(err) {
			return nil, err
		}
		prevApplication := latestApplication
		latestApplication, err = d.applicationClient.Apps(latestApplication.Namespace).Get(ctx, latestApplication.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevApplication.UID != latestApplication.UID {
			return nil, fmt.Errorf("app uid has changed across retries")
		}
	}
}

// updateApplicationStatusFunc will verify that the status of the app is correct
func (d *applicationResourcesDeleter) updateApplicationStatusFunc(ctx context.Context, app *applicationv1.App) (*applicationv1.App, error) {
	if app.DeletionTimestamp.IsZero() || app.Status.Phase == applicationv1.AppPhaseTerminating {
		return app, nil
	}
	newApplication := applicationv1.App{}
	newApplication.ObjectMeta = app.ObjectMeta
	newApplication.Status = app.Status
	newApplication.Status.Phase = applicationv1.AppPhaseTerminating
	return d.applicationClient.Apps(newApplication.Namespace).UpdateStatus(ctx, &newApplication, metav1.UpdateOptions{})
}

// finalized returns true if the app.Spec.Finalizers is an empty list
func finalized(app *applicationv1.App) bool {
	return len(app.Spec.Finalizers) == 0
}

// finalizeApplication removes the specified finalizerToken and finalizes the app
func (d *applicationResourcesDeleter) finalizeApplication(ctx context.Context, app *applicationv1.App) (*applicationv1.App, error) {
	applicationFinalize := applicationv1.App{}
	applicationFinalize.ObjectMeta = app.ObjectMeta
	applicationFinalize.Spec = app.Spec
	finalizerSet := sets.NewString()
	for i := range app.Spec.Finalizers {
		if app.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(app.Spec.Finalizers[i]))
		}
	}
	applicationFinalize.Spec.Finalizers = make([]applicationv1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		applicationFinalize.Spec.Finalizers = append(applicationFinalize.Spec.Finalizers, applicationv1.FinalizerName(value))
	}
	app = &applicationv1.App{}
	var err error
	if d.applicationClient.RESTClient() != nil && !reflect.ValueOf(d.applicationClient.RESTClient()).IsNil() {
		err = d.applicationClient.RESTClient().Put().
			Resource("apps").
			Name(applicationFinalize.Name).
			Namespace(applicationFinalize.Namespace).
			SubResource("finalize").
			Body(&applicationFinalize).
			Do(ctx).
			Into(app)
	} else {
		app, err = d.applicationClient.Apps(applicationFinalize.Namespace).Update(ctx, &applicationFinalize, metav1.UpdateOptions{})
	}

	if err != nil {
		// it was removed already, so life is good
		if k8serrors.IsNotFound(err) {
			return app, nil
		}
	}
	return app, err
}

type deleteResourceFunc func(ctx context.Context,
	deleter *applicationResourcesDeleter,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteApplication,
}

// deleteAllContent will use the dynamic authClient to delete each resource identified in groupVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *applicationResourcesDeleter) deleteAllContent(ctx context.Context, app *applicationv1.App) error {
	log.Info("App controller - deleteAllContent", log.String("namespace", app.Namespace), log.String("appName", app.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(ctx, d, app, d.repo)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.Info("App controller - deletedAllContent", log.String("namespace", app.Namespace), log.String("appName", app.Name))
	return nil
}

func deleteApplication(ctx context.Context,
	deleter *applicationResourcesDeleter,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration) error {
	_, err := action.Uninstall(ctx, deleter.applicationClient, deleter.platformClient, app, repo)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) || k8serrors.IsNotFound(err) {
			log.Warn(err.Error())
			return nil
		}
	}
	return err
}
