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

# set the kubernetes apimachinery package dir
K8S_APIMACHINERY_DIR = $(shell go list -f '{{ .Dir }}' -m k8s.io/apimachinery)
# set the kubernetes api package dir
K8S_API_DIR = $(shell go list -f '{{ .Dir }}' -m k8s.io/api)
# set the gogo protobuf package dir
GOGO_PROTOBUF_DIR = $(shell go list -f '{{ .Dir }}' -m github.com/gogo/protobuf)

.PHONY: gen.run
gen.run: gen.clean gen.generator gen.api gen.openapi gen.gateway gen.registry gen.monitor

# ==============================================================================
# Generator

.PHONY: gen.generator
gen.generator:
	@echo "===========> Building code generator $(VERSION) docker image"
	@$(DOCKER) build --pull -t $(REGISTRY_PREFIX)/code-generator:$(VERSION) -f $(ROOT_DIR)/build/docker/tools/code-generator/Dockerfile $(ROOT_DIR)/build/docker/tools/code-generator

.PHONY: gen.api
gen.api:
	$(eval CODE_GENERATOR_VERSION := $(shell $(DOCKER) images --filter 'reference=$(REGISTRY_PREFIX)/code-generator' |sed 's/[ ][ ]*/,/g' |cut -d ',' -f 2 |sed -n '2p'))
	@$(DOCKER) run --rm \
		-v $(ROOT_DIR):/go/src/$(ROOT_PACKAGE) \
	 	$(REGISTRY_PREFIX)/code-generator:$(CODE_GENERATOR_VERSION) \
	 	/root/code.sh \
	 	all \
	 	$(ROOT_PACKAGE)/api/client \
	 	$(ROOT_PACKAGE)/api \
	 	$(ROOT_PACKAGE)/api \
	 	"platform:v1 business:v1 notify:v1 registry:v1 monitor:v1"

.PHONY: gen.gateway
gen.gateway:
	$(eval CODE_GENERATOR_VERSION := $(shell $(DOCKER) images --filter 'reference=$(REGISTRY_PREFIX)/code-generator' |sed 's/[ ][ ]*/,/g' |cut -d ',' -f 2 |sed -n '2p'))
	@$(DOCKER) run --rm \
		-v $(ROOT_DIR):/go/src/$(ROOT_PACKAGE) \
	 	$(REGISTRY_PREFIX)/code-generator:$(CODE_GENERATOR_VERSION) \
	 	/root/code.sh \
	 	deepcopy-internal,deepcopy-external,defaulter-external,conversion-external \
	 	$(ROOT_PACKAGE)/pkg/gateway/apis \
	 	$(ROOT_PACKAGE)/pkg/gateway/apis \
	 	$(ROOT_PACKAGE)/pkg/gateway/apis \
	 	"config:v1"

.PHONY: gen.registry
gen.registry:
	$(eval CODE_GENERATOR_VERSION := $(shell $(DOCKER) images --filter 'reference=$(REGISTRY_PREFIX)/code-generator' |sed 's/[ ][ ]*/,/g' |cut -d ',' -f 2 |sed -n '2p'))
	@$(DOCKER) run --rm \
		-v $(ROOT_DIR):/go/src/$(ROOT_PACKAGE) \
	 	$(REGISTRY_PREFIX)/code-generator:$(CODE_GENERATOR_VERSION) \
	 	/root/code.sh \
	 	deepcopy-internal,deepcopy-external,defaulter-external,conversion-external \
	 	$(ROOT_PACKAGE)/pkg/registry/apis \
	 	$(ROOT_PACKAGE)/pkg/registry/apis \
	 	$(ROOT_PACKAGE)/pkg/registry/apis \
	 	"config:v1"

.PHONY: gen.monitor
gen.monitor:
	$(eval CODE_GENERATOR_VERSION := $(shell $(DOCKER) images --filter 'reference=$(REGISTRY_PREFIX)/code-generator' |sed 's/[ ][ ]*/,/g' |cut -d ',' -f 2 |sed -n '2p'))
	@$(DOCKER) run --rm \
		-v $(ROOT_DIR):/go/src/$(ROOT_PACKAGE) \
	 	$(REGISTRY_PREFIX)/code-generator:$(CODE_GENERATOR_VERSION) \
	 	/root/code.sh \
	 	deepcopy-internal,deepcopy-external,defaulter-external,conversion-external \
	 	$(ROOT_PACKAGE)/pkg/monitor/apis \
	 	$(ROOT_PACKAGE)/pkg/monitor/apis \
	 	$(ROOT_PACKAGE)/pkg/monitor/apis \
	 	"config:v1"

.PHONY: gen.openapi
gen.openapi:
	$(eval CODE_GENERATOR_VERSION := $(shell $(DOCKER) images --filter 'reference=$(REGISTRY_PREFIX)/code-generator' |sed 's/[ ][ ]*/,/g' |cut -d ',' -f 2 |sed -n '2p'))
	@$(DOCKER) run --rm \
    	-v $(ROOT_DIR):/go/src/$(ROOT_PACKAGE) \
		-v $(K8S_APIMACHINERY_DIR):/go/src/k8s.io/apimachinery \
		-v $(K8S_API_DIR):/go/src/k8s.io/api \
	 	$(REGISTRY_PREFIX)/code-generator:$(CODE_GENERATOR_VERSION) \
	 	/root/openapi.sh \
	 	$(ROOT_PACKAGE)/api/openapi \
	 	$(shell ${ROOT_DIR}/build/script/openapi.sh)

.PHONY: gen.clean
gen.clean:
	@rm -rf ./api/client/{clientset,informers,listers}
	@find . -type f -name 'generated.*' -delete
	@find . -type f -name 'zz_generated*.go' -delete
	@find . -type f -name 'types_swagger_doc_generated.go' -delete

