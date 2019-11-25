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

# ==============================================================================
# Makefile helper functions for deploy to developer env
#

KUBECTL := kubectl
NAMESPACE ?= tke
CONTEXT ?= tkestack.dev

DEPLOYS=tke-auth tke-registry-api tke-platform-api tke-platform-controller tke-business-api tke-business-controller tke-notify-api tke-notify-controller tke-monitor-api tke-monitor-controller tke-gateway

.PHONY: deploy.run.all
deploy.run.all:
	@echo "===========> Deploying all"
	@$(MAKE) deploy.run

.PHONY: deploy.run
deploy.run: $(addprefix deploy.run., $(DEPLOYS))

.PHONY: deploy.run.%
deploy.run.%:
	@echo "===========> Deploying $* $(VERSION)"
	@$(KUBECTL) -n $(NAMESPACE) --context=$(CONTEXT) set image deployment/$* $*=$(REGISTRY_PREFIX)/$*:$(VERSION)

.PHONY: deploy.run.tke-gateway
deploy.run.tke-gateway:
	@echo "===========> Deploying tke-gateway $(VERSION)"
	@$(KUBECTL) -n $(NAMESPACE) --context=$(CONTEXT) set image daemonset/tke-gateway tke-gateway=$(REGISTRY_PREFIX)/tke-gateway:$(VERSION)
