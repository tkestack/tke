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

NODE_SUPPORTED_VERSION = v12

.PHONY: web.build
web.build: web.verify web.build.console web.build.installer

.PHONY: web.verify
web.verify:
	@echo "===========> Check NOde.js version"
ifneq ($(shell node -v | cut -f1 -d.), $(NODE_SUPPORTED_VERSION))
	@echo "===========> Install Node.js v12"
	curl -sL https://deb.nodesource.com/setup_12.x | sudo -E bash -
	sudo apt-get install -y nodejs
endif

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
