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

.PHONY: all
all: lint test build

# ==============================================================================
# Build options

ROOT_PACKAGE=tkestack.io/tke
VERSION_PACKAGE=tkestack.io/tke/pkg/app/version

# ==============================================================================
# Includes

include build/lib/common.mk
include build/lib/golang.mk
include build/lib/image.mk
include build/lib/deploy.mk
include build/lib/asset.mk
include build/lib/web.mk
include build/lib/gen.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  DEBUG        Whether to generate debug symbols. Default is 0.
  BINS         The binaries to build. Default is all of cmd.
               This option is available when using: make build/build.multiarch
               Example: make build BINS="tke-platform-api tke-platform-controller"
  IMAGES       Backend images to make. Default is all of cmd starting with tke-.
               This option is available when using: make image/image.multiarch/push/push.multiarch
               Example: make image.multiarch IMAGES="tke-platform-api tke-platform-controller"
  PLATFORMS    The multiple platforms to build. Default is linux_amd64 and linux_arm64.
               This option is available when using: make build.multiarch/image.multiarch/push.multiarch
               Example: make image.multiarch IMAGES="tke-platform-api tke-platform-controller" PLATFORMS="linux_amd64 linux_arm64"
  VERSION      The version information compiled into binaries.
               The default is obtained from git.
  V            Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

# ==============================================================================
# Targets

## gen: Generate codes for API definitions.
.PHONY: gen
gen:
	@$(MAKE) gen.run

## asset: Embed front-end static files and documentation in the app.
.PHONY: asset
asset:
	@$(MAKE) asset.build

## web: Builds the web console app for production.
.PHONY: web
web:
	@$(MAKE) web.build

## build: Build source code for host platform.
.PHONY: build
build:
	@$(MAKE) go.build

## build.multiarch: Build source code for multiple platforms. See option PLATFORMS.
.PHONY: build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch

## image: Build docker images for host arch.
.PHONY: image
image:
	@$(MAKE) image.build

## image.multiarch: Build docker images for multiple platforms. See option PLATFORMS.
.PHONY: image.multiarch
image.multiarch:
	@$(MAKE) image.build.multiarch

## push: Build docker images for host arch and push images to registry.
.PHONY: push
push:
	@$(MAKE) image.push

## push.multiarch: Build docker images for multiple platforms and push images to registry.
.PHONY: push.multiarch
push.multiarch:
	@$(MAKE) image.push.multiarch

## manifest: Build docker images for host arch and push manifest list to registry.
.PHONY: manifest
manifest:
	@$(MAKE) image.manifest.push

## manifest.multiarch: Build docker images for multiple platforms and push manifest lists to registry.
.PHONY: manifest.multiarch
manifest.multiarch:
	@$(MAKE) image.manifest.push.multiarch

## deploy: Deploy updated components to development env.
.PHONY: deploy
deploy:
	@$(MAKE) deploy.run

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@$(MAKE) go.clean

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint:
	@$(MAKE) go.lint

## test: Run unit test.
.PHONY: test
test:
	@$(MAKE) go.test

.PHONY: release.provider
release.provider:
	cd build/docker/tools/provider-res && make all

.PHONY: release.build
release.build: release.provider
	make push.multiarch

## release: Release tke
.PHONY: release
release:
	build/docker/tools/tke-installer/release.sh

## release-test: test release
.PHONY: release-test
release-test:
	go test -v -timeout=1200m tkestack.io/tke/test/e2e_installer

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
