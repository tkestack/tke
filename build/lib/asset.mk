# Tencent is pleased to support the open source community by making TKEStack
# available.
#
# Copyright (C) 2012-2019 Tencent. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use
# this file except in compliance with the License. You may obtain a copy of the
# License at
#
# https://opensource.org/licenses/Apache-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OF ANY KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations under the License.

.PHONY: asset.build
asset.build: asset.build.addon asset.build.web.console asset.build.web.installer

.PHONY: asset.verify
asset.verify:
ifeq (,$(wildcard $(GOBIN)/staticfiles))
	@echo "===========> Installing bou.ke staticfiles"
	@GO111MODULE=off $(GO) get -u bou.ke/staticfiles
endif

.PHONY: asset.build.addon
asset.build.addon: asset.verify
	@echo "===========> Bundling addon readme assets to application"
	@mkdir -p $(ROOT_DIR)/hack/addon/readme
	@mkdir -p $(ROOT_DIR)/pkg/platform/registry/clusteraddontype/assets
	$(eval ASSETS_DIR := $(shell (cd $(ROOT_DIR); ls -d -1 ./hack/addon/readme 2>/dev/null || echo ../../../hack/addon/readme)))
	@$(GOBIN)/staticfiles \
		-o $(ROOT_DIR)/pkg/platform/registry/clusteraddontype/assets/assets.go \
		$(ASSETS_DIR)/

.PHONY: asset.build.web.console
asset.build.web.console: asset.verify
	@echo "===========> Bundling console web assets to application"
	@mkdir -p $(ROOT_DIR)/web/console/build
	@mkdir -p $(ROOT_DIR)/pkg/gateway/assets
	$(eval ASSETS_DIR := $(shell (cd $(ROOT_DIR); ls -d -1 ./web/console/build 2>/dev/null || echo ../../../web/console/build)))
	@$(GOBIN)/staticfiles \
		-o $(ROOT_DIR)/pkg/gateway/assets/assets.go \
		$(ASSETS_DIR)/

.PHONY: asset.build.web.installer
asset.build.web.installer: asset.verify
	@echo "===========> Bundling installer web assets to application"
	@mkdir -p $(ROOT_DIR)/web/installer/build
	@mkdir -p $(ROOT_DIR)/cmd/tke-installer/assets
	$(eval ASSETS_DIR := $(shell (cd $(ROOT_DIR); ls -d -1 ./web/installer/build 2>/dev/null || echo ../../../web/installer/build)))
	@$(GOBIN)/staticfiles \
		-o $(ROOT_DIR)/cmd/tke-installer/assets/assets.go \
		$(ASSETS_DIR)/
