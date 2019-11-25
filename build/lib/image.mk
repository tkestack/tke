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
# Makefile helper functions for docker image
#

DOCKER := docker
DOCKER_SUPPORTED_VERSIONS ?= 17|18|19

REGISTRY_PREFIX ?= tkestack
BASE_IMAGE = alpine:3.10

EXTRA_ARGS ?=
_DOCKER_BUILD_EXTRA_ARGS :=

ifdef HTTP_PROXY
_DOCKER_BUILD_EXTRA_ARGS += --build-arg HTTP_PROXY=${HTTP_PROXY}
endif

ifneq ($(EXTRA_ARGS), )
_DOCKER_BUILD_EXTRA_ARGS += $(EXTRA_ARGS)
endif

# Determine image files by looking into build/docker/*/Dockerfile
IMAGES_DIR ?= $(wildcard ${ROOT_DIR}/build/docker/*)
# Determine images names by stripping out the dir names
IMAGES ?= $(filter-out tools,$(foreach image,${IMAGES_DIR},$(notdir ${image})))

ifeq (${IMAGES},)
  $(error Could not determine IMAGES, set ROOT_DIR or run in source dir)
endif

.PHONY: image.verify
image.verify:
ifneq ($(shell $(DOCKER) -v | grep -q -E '\bversion ($(DOCKER_SUPPORTED_VERSIONS))\b' && echo 0 || echo 1), 0)
	$(error unsupported docker version. Please make install one of the following supported version: '$(DOCKER_SUPPORTED_VERSIONS)')
endif
	@echo "===========> Docker version verification passed"

.PHONY: image.build
image.build: image.verify go.build.verify $(addprefix image.build., $(IMAGES))

.PHONY: image.push
image.push: image.verify go.build.verify $(addprefix image.push., $(IMAGES))

.PHONY: image.build.%
image.build.%: go.build.linux_amd64.%
	@echo "===========> Building $* $(VERSION) docker image"
	@mkdir -p $(TMP_DIR)/$*
	@cat $(ROOT_DIR)/build/docker/$*/Dockerfile\
		| sed "s#BASE_IMAGE#$(BASE_IMAGE)#g" >$(TMP_DIR)/$*/Dockerfile
	@cp ${OUTPUT_DIR}/linux/amd64/$* $(TMP_DIR)/$*/
	@DST_DIR=$(TMP_DIR)/$* $(ROOT_DIR)/build/docker/$*/build.sh 2>/dev/null || true
	@$(DOCKER) build $(_DOCKER_BUILD_EXTRA_ARGS) --pull -t $(REGISTRY_PREFIX)/$*:$(VERSION) $(TMP_DIR)/$*
	@rm -rf $(TMP_DIR)/$*

.PHONY: image.push.%
image.push.%: image.build.%
	@echo "===========> Pushing $* $(VERSION) image to $(REGISTRY_PREFIX)"
	@$(DOCKER) push $(REGISTRY_PREFIX)/$*:$(VERSION)
