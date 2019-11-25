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

NPM := npm
NPM_SUPPORTED_VERSIONS ?= 6

.PHONY: web.build
web.build: web.verify web.build.console web.build.installer

.PHONY: web.verify
web.verify:
ifneq ($(shell $(NPM) -v | grep -q -E '\b($(NPM_SUPPORTED_VERSIONS))\b' && echo 0 || echo 1), 0)
	$(error unsupported npm version. Please make install one of the following supported version: '$(NPM_SUPPORTED_VERSIONS)')
endif

.PHONY: web.build.console
web.build.console:
	@echo "===========> Building the console web app"
	@mkdir -p $(ROOT_DIR)/web/console/build
	@cd $(ROOT_DIR)/web/console && $(NPM) run build

.PHONY: web.build.installer
web.build.installer:
	@echo "===========> Building the installer web app"
	@mkdir -p $(ROOT_DIR)/web/installer/build
	@cd $(ROOT_DIR)/web/installer && $(NPM) run build
