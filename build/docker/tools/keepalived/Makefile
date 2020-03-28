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

IMAGE = keepalived
VERSION = 2.0.16-r0
REGISTRY_PREFIX ?= tkestack

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../../../.. && pwd -P))
DOCKERFILE_DIR := $(ROOT_DIR)/build/docker/tools/keepalived

.PHONY: all
all: image

# ==============================================================================
# Includes
include $(ROOT_DIR)/build/lib/common.mk
include $(ROOT_DIR)/build/lib/docker-buildx.mk

# ==============================================================================

## image: Build docker images for host arch.
.PHONY: image
image:
	@VERSION=$(VERSION) WORK_DIR=$(DOCKERFILE_DIR) $(MAKE) docker.buildx.$(IMAGE_PLAT).$(IMAGE)

## manifest: Build docker images for host arch and push manifest list to registry.
.PHONY: manifest
manifest:
	@VERSION=$(VERSION) WORK_DIR=$(DOCKERFILE_DIR) $(MAKE) docker.manifest.$(IMAGE_PLAT).$(IMAGE)

## manifest.multiarch: Build docker images for multiple platforms and push manifest lists to registry.
.PHONY: manifest.multiarch
manifest.multiarch:
	@set -o errexit; \
	for platform in $(PLATFORMS); do \
		VERSION=$(VERSION) WORK_DIR=$(DOCKERFILE_DIR) $(MAKE) docker.push.$${platform}.$(IMAGE); \
	done
	@echo "===========> Pushing manifest $(REGISTRY_PREFIX)/$(IMAGE):$(VERSION) and then remove the local manifest list"
	REGISTRY_PREFIX=$(REGISTRY_PREFIX) PLATFROMS="$(PLATFORMS)" IMAGE=$(IMAGE) VERSION=$(VERSION) \
	  DOCKER_CLI_EXPERIMENTAL=enabled $(ROOT_DIR)/build/lib/create-manifest.sh 

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
