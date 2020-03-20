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

ASSETS_GENERATOR_VERSION := v1.0.0

.PHONY: asset.build
asset.build: asset.build.addon asset.build.web.console asset.build.web.installer

.PHONY: asset.build.addon
asset.build.addon:
	@echo "===========> Bundling addon readme assets to application"
	@docker run --rm \
		-v $(ROOT_DIR)/hack/addon/readme:/src \
		-v $(ROOT_DIR)/pkg/platform/registry/clusteraddontype/assets:/assets \
		$(REGISTRY_PREFIX)/assets-generator:$(ASSETS_GENERATOR_VERSION) -o /assets/assets.go /src

.PHONY: asset.build.web.console
asset.build.web.console:
	@echo "===========> Bundling console web assets to application"
	@docker run --rm \
		-v $(ROOT_DIR)/web/console/build:/src \
		-v $(ROOT_DIR)/pkg/gateway/assets:/assets \
		$(REGISTRY_PREFIX)/assets-generator:$(ASSETS_GENERATOR_VERSION) -o /assets/assets.go /src

.PHONY: asset.build.web.installer
asset.build.web.installer:
	@echo "===========> Bundling installer web assets to application"
	@docker run --rm \
		-v $(ROOT_DIR)/web/installer/build:/src \
		-v $(ROOT_DIR)/cmd/tke-installer/assets:/assets \
		$(REGISTRY_PREFIX)/assets-generator:$(ASSETS_GENERATOR_VERSION) -o /assets/assets.go /src
