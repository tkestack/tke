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

package tenant

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/howeyc/fsnotify"
	"github.com/pborman/uuid"

	"tkestack.io/tke/pkg/auth/handler/policy"
	"tkestack.io/tke/pkg/auth/registry"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/etcd"
	"tkestack.io/tke/pkg/util/log"
)

const (
	watchDebounceDelay = 100 * time.Millisecond
)

// Helper operates the resource in the tenants.
type Helper struct {
	regist        *registry.Registry
	policyService *policy.Service
	watcher       *fsnotify.Watcher
	policyFile    string
	catogoryFile  string
	adminName     string
	adminSecret   string
}

// NewHelper create a new tenant operate instance.
func NewHelper(regist *registry.Registry, policyService *policy.Service, policyFile string, categoryFile string,
	admin string, secret string) *Helper {

	helper := &Helper{
		regist:        regist,
		policyService: policyService,
		policyFile:    policyFile,
		catogoryFile:  categoryFile,
		adminName:     admin,
		adminSecret:   secret,
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("New file watcher failed", log.Err(err))
		return helper
	}

	// watch the parent directory of the target files so we can catch
	// symlink updates of k8s ConfigMaps volumes.
	for _, file := range []string{helper.policyFile, helper.catogoryFile} {
		if file != "" {
			watchDir, _ := filepath.Split(file)
			if err := watcher.Watch(watchDir); err != nil {
				log.Error("could not watch %v: %v", log.String("file", file), log.Err(err))
				return helper
			}
		}
	}

	helper.watcher = watcher
	go helper.loadConfig()
	return helper
}

// LoadResourceAllTenant reloads initial resource for all tenants.
func (h *Helper) LoadResourceAllTenant() {
	conns, err := h.regist.DexStorage().ListConnectors()
	if err != nil {
		log.Error("List all tenant failed", log.Err(err))
		return
	}

	for _, conn := range conns {
		_ = h.LoadResource(conn.ID)
	}
}

//LoadResource loads initial resource for then tenant.
func (h *Helper) LoadResource(tenantID string) error {
	if len(h.catogoryFile) > 0 {
		if err := h.loadCategory(tenantID); err != nil {
			return err
		}
	}

	if len(h.policyFile) > 0 {
		if err := h.loadPolicy(tenantID); err != nil {
			return err
		}
	}

	return h.CreateAdmin(tenantID)
}

//CreateAdmin creates a admin user for then tenant.
func (h *Helper) CreateAdmin(tenantID string) error {
	_, err := h.regist.LocalIdentityStorage().Get(tenantID, h.adminName)
	if err == nil || err != etcd.ErrNotFound {
		return err
	}

	log.Info("Create admin for tenant", log.String("tenant", tenantID), log.String("user", h.adminName))
	hashPassword := base64.StdEncoding.EncodeToString([]byte(h.adminSecret))
	bcryptedPasswd, err := util.BcryptPassword(hashPassword)
	if err != nil {
		log.Error("Bcrypt hash password failed", log.Err(err))
		return err
	}

	identity := &types.LocalIdentity{
		Name: h.adminName,
		UID:  uuid.New(),
		Spec: &types.LocalIdentitySpec{
			HashedPassword: bcryptedPasswd,
			TenantID:       tenantID,
			Extra: map[string]string{
				"displayName": "Administrator",
			},
		},
		Status:   &types.LocalIdentityStatus{Locked: false},
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	err = h.regist.LocalIdentityStorage().Create(identity)
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) loadCategory(tenantID string) error {
	log.Info("Load categories for the tenant", log.String("tenant", tenantID))
	var categoryList []*types.Category
	bytes, err := ioutil.ReadFile(h.catogoryFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &categoryList)
	if err != nil {
		return err
	}

	categoryStorage := h.regist.CategoryStorage()

	return categoryStorage.Load(tenantID, categoryList)
}

func (h *Helper) loadPolicy(tenantID string) error {
	log.Info("Load predefine policies for the tenant", log.String("tenant", tenantID))

	var policyList []*types.Policy
	bytes, err := ioutil.ReadFile(h.policyFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &policyList)
	if err != nil {
		return err
	}

	for _, pol := range policyList {
		pol.UserName = h.adminName
	}

	return h.policyService.LoadPredefinePolicies(tenantID, policyList)
}

// loadConfig watch auth config file from the file system and reload config.
func (h *Helper) loadConfig() {
	defer h.watcher.Close()
	var timerC <-chan time.Time

	for {
		select {
		case <-timerC:
			timerC = nil
			log.Info("Config directory changed, load it")
			conns, err := h.regist.DexStorage().ListConnectors()
			if err != nil {
				log.Error("List all tenant failed", log.Err(err))
				break
			}

			if len(h.catogoryFile) > 0 {
				for _, conn := range conns {
					_ = h.loadCategory(conn.ID)
				}
			}

			if len(h.policyFile) > 0 {
				for _, conn := range conns {
					_ = h.loadPolicy(conn.ID)
				}
			}
		case event := <-h.watcher.Event:
			// use a timer to debounce configuration updates
			if (event.IsModify() || event.IsCreate()) && timerC == nil {
				timerC = time.After(watchDebounceDelay)
			}
		case err := <-h.watcher.Error:
			log.Errorf("Watcher error: %v", err)
		}
	}
}
