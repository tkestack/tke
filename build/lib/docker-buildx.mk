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

################################################
## Makefile helper functions for docker buildx
## Please set VERSION and WORK_DIR before use
## Example: VERSION=$(VERSION) WORK_DIR=$(Dockerfile_DIR) make docker.buildx.linux_amd64.keepalived
################################################

DOCKER := docker
DOCKER_SUPPORTED_API_VERSION ?= 1.40
DOCKER_VERSION ?= 19.03

REGISTRY_PREFIX ?= tkestack

EXTRA_ARGS ?=
_DOCKER_BUILD_EXTRA_ARGS :=

ifdef HTTP_PROXY
_DOCKER_BUILD_EXTRA_ARGS += --build-arg HTTP_PROXY=${HTTP_PROXY}
endif

ifneq ($(EXTRA_ARGS), )
_DOCKER_BUILD_EXTRA_ARGS += $(EXTRA_ARGS)
endif

.PHONY: docker.verify
docker.verify:
	$(eval API_VERSION := $(shell $(DOCKER) version | grep -E 'API version: {6}[0-9]' | awk '{print $$3} END { if (NR==0) print 0}' ))
	$(eval PASS := $(shell echo "$(API_VERSION) >= $(DOCKER_SUPPORTED_API_VERSION)" | bc))
	@if [ $(PASS) -ne 1 ]; then \
		$(DOCKER) -v ;\
		echo "Unsupported docker version. Docker API version should be greater than $(DOCKER_SUPPORTED_API_VERSION) (Or docker version: $(DOCKER_VERSION))"; \
		exit 1; \
	fi

.PHONY: docker.buildx.%
docker.buildx.%: docker.verify
	$(eval IMAGE := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	$(eval IMAGE_NAME := $(REGISTRY_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION))
	@echo "===========> Building docker image $(IMAGE) $(VERSION) for $(IMAGE_PLAT)"
	DOCKER_CLI_EXPERIMENTAL=enabled $(DOCKER) buildx build --pull --platform $(IMAGE_PLAT) --load \
	  -t $(IMAGE_NAME) $(_DOCKER_BUILD_EXTRA_ARGS) $(WORK_DIR)

.PHONY: docker.push.%
docker.push.%: docker.buildx.%
	@echo "===========> Pushing image $(IMAGE_NAME)"
	$(DOCKER) push $(IMAGE_NAME)

.PHONY: docker.manifest.%
docker.manifest.%: export DOCKER_CLI_EXPERIMENTAL := enabled
docker.manifest.%: docker.push.% docker.manifest.remove.%
	$(eval MANIFEST_NAME := $(REGISTRY_PREFIX)/$(IMAGE):$(VERSION))
	@echo "===========> Pushing manifest $(MANIFEST_NAME) and then remove the local manifest list"
	@$(DOCKER) manifest create $(MANIFEST_NAME) \
	  $(IMAGE_NAME)
	@$(DOCKER) manifest annotate $(MANIFEST_NAME) \
	  $(IMAGE_NAME) \
	  --os $(OS) --arch ${ARCH}
	@$(DOCKER) manifest push --purge $(MANIFEST_NAME)

# Docker cli has a bug: https://github.com/docker/cli/issues/954
# If you find your manifests were not updated,
# Please manually delete them in $HOME/.docker/manifests/
# and re-run.
.PHONY: docker.manifest.remove.%
docker.manifest.remove.%:
	@rm -rf ${HOME}/.docker/manifests/docker.io_$(REGISTRY_PREFIX)_$(IMAGE)-$(VERSION)
