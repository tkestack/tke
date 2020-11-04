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

NODE_SUPPORTED_VERSION = 12
NPM = npm

.PHONY: web.build
web.build: web.verify web.build.console web.build.installer

.PHONY: web.verify
web.verify:
	@echo "===========> Check Node.js version"
	$(eval NODE_CURRENT_VERSION := $(shell node -v | cut -f1 -d. | cut -f2 -dv ))
	@if [ $(NODE_CURRENT_VERSION) -lt $(NODE_SUPPORTED_VERSION) ]; then \
		echo "you need upgrade Node.js version to $(NODE_SUPPORTED_VERSION) or newer"; \
		exit 1; \
	fi




.PHONY: web.build.console
web.build.console: web.verify
	@echo "===========> Building the console web app"
	@mkdir -p $(ROOT_DIR)/web/console/build
	@cd $(ROOT_DIR)/web/console && $(NPM) run build

.PHONY: web.build.installer
web.build.installer: web.verify
	@echo "===========> Building the installer web app"
	@mkdir -p $(ROOT_DIR)/web/installer/build
	@cd $(ROOT_DIR)/web/installer && $(NPM) run build
