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

NVM_VERSION = v0.38.0

.PHONY: web.build
web.build: web.verify web.build.console web.build.installer

.PHONY: web.verify
.ONESHELL:
web.verify:
	curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/$(NVM_VERSION)/install.sh | bash

.PHONY: web.build.console
.ONESHELL:
web.build.console: web.verify
	@echo "===========> Building the console web app"
	@cd $(ROOT_DIR)/web/console
	@source $(HOME)/.nvm/nvm.sh
	@nvm install
	@npm run build

.PHONY: web.build.installer
.ONESHELL:
web.build.installer: web.verify
	@echo "===========> Building the Installer web app"
	@cd $(ROOT_DIR)/web/installer
	@source $(HOME)/.nvm/nvm.sh
	@nvm install
	@npm run build
