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

package hooks

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/fields"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/util/log"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	dexstorage "github.com/dexidp/dex/storage"
	"github.com/howeyc/fsnotify"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

const (
	watchDebounceDelay = 100 * time.Millisecond
)

type policyHookHandler struct {
	authClient   authinternalclient.AuthInterface
	dexStorage   dexstorage.Storage
	watcher      *fsnotify.Watcher
	policyFile   string
	categoryFile string
}

// NewPolicyHookHandler creates a new policyHookHandler object.
func NewPolicyHookHandler(authClient authinternalclient.AuthInterface, dexstorage dexstorage.Storage, policyFile string, categoryFile string) genericapiserver.PostStartHookProvider {
	return &policyHookHandler{
		authClient:   authClient,
		dexStorage:   dexstorage,
		policyFile:   policyFile,
		categoryFile: categoryFile,
	}
}

func (d *policyHookHandler) PostStartHook() (string, genericapiserver.PostStartHookFunc, error) {
	return "sync-default-policy", func(context genericapiserver.PostStartHookContext) error {
		if err := d.loadConfig(); err != nil {
			return err
		}

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Error("New file watcher failed", log.Err(err))
			return err
		}

		// watch the parent directory of the target files so we can catch
		// symlink updates of k8s ConfigMaps volumes.
		for _, file := range []string{d.policyFile, d.categoryFile} {
			if file != "" {
				watchDir, _ := filepath.Split(file)
				if err := watcher.Watch(watchDir); err != nil {
					log.Error("could not watch %v: %v", log.String("file", file), log.Err(err))
					return err
				}
			}
		}

		d.watcher = watcher

		go d.pollReload(context.StopCh)
		return nil
	}, nil
}

func (d *policyHookHandler) loadConfig() error {
	if d.categoryFile != "" {
		err := d.loadCategory()
		if err != nil {
			return err
		}
	}

	if d.policyFile != "" {
		conns, err := d.dexStorage.ListConnectors()
		if err != nil {
			log.Error("List all tenant failed", log.Err(err))
			return err
		}

		for _, conn := range conns {
			err := d.loadPolicy(conn.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// loadConfig watch auth config file from the file system and reload config.
func (d *policyHookHandler) pollReload(stopCh <-chan struct{}) {
	defer d.watcher.Close()
	var timerC <-chan time.Time

	for {
		select {
		case <-timerC:
			timerC = nil
			log.Info("Policy config directory changed, loadConfig it")

			if err := d.loadConfig(); err != nil {
				log.Errorf("Load config failed after changed", log.Err(err))
			}
		case event := <-d.watcher.Event:
			// use a timer to debounce configuration updates
			if (event.IsModify() || event.IsCreate()) && timerC == nil {
				timerC = time.After(watchDebounceDelay)
			}
		case err := <-d.watcher.Error:
			log.Errorf("Watcher error: %v", err)
		case <-stopCh:
			return
		}
	}
}

func (d *policyHookHandler) loadCategory() error {
	var categoryList []*auth.Category
	bytes, err := ioutil.ReadFile(d.categoryFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &categoryList)
	if err != nil {
		return err
	}

	var errs []error

	for _, cat := range categoryList {
		categorySelector := fields.OneTermEqualSelector("spec.categoryName", cat.Spec.CategoryName)

		result, err := d.authClient.Categories().List(metav1.ListOptions{FieldSelector: categorySelector.String()})
		if err != nil {
			return err
		}

		if len(result.Items) > 0 {
			exists := result.Items[0]
			if !reflect.DeepEqual(exists.Spec, cat.Spec) {
				exists.Spec = cat.Spec
				_, err = d.authClient.Categories().Update(&exists)

				if err != nil {
					errs = append(errs, err)
				}
			}
		} else {
			_, err = d.authClient.Categories().Create(cat)

			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func (d *policyHookHandler) loadPolicy(tenantID string) error {
	var policyList []*auth.Policy
	bytes, err := ioutil.ReadFile(d.policyFile)
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

		result, err := d.authClient.Policies().List(metav1.ListOptions{FieldSelector: policySelector.String()})
		if err != nil {
			return err
		}
		log.Info("selector result", log.String("selector", policySelector.String()), log.Any("result", result))
		pol.Spec.Type = auth.PolicyDefault
		pol.Spec.TenantID = tenantID
		pol.Spec.Username = "admin"
		pol.Spec.Finalizers = []auth.FinalizerName{
			auth.PolicyFinalize,
		}
		if len(result.Items) > 0 {
			exists := result.Items[0]
			if !reflect.DeepEqual(exists.Spec, pol.Spec) {
				exists.Spec = pol.Spec
				_, err = d.authClient.Policies().Update(&exists)

				if err != nil {
					errs = append(errs, err)
				}
			}
		} else {
			_, err = d.authClient.Policies().Create(pol)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}
