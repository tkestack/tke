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
EXT_PB_APIS = "k8s.io/api/core/v1 k8s.io/api/apps/v1"
# set the code generator image version
CODE_GENERATOR_VERSION := v1.19.7
FIND := find . ! -path './pkg/platform/provider/baremetal/apis/*'

.PHONY: gen.run
gen.run: gen.clean gen.api gen.openapi gen.gateway gen.registry gen.monitor gen.mesh gen.audit

# ==============================================================================
# Generator

.PHONY: gen.api
gen.api:
	@$(DOCKER) run --rm \
		-v $(ROOT_DIR):/go/src/$(ROOT_PACKAGE) \
		-e EXT_PB_APIS=$(EXT_PB_APIS)\
	 	$(REGISTRY_PREFIX)/code-generator:$(CODE_GENERATOR_VERSION) \
	 	/root/code.sh \
	 	all \
	 	$(ROOT_PACKAGE)/api/client \
	 	$(ROOT_PACKAGE)/api \
	 	$(ROOT_PACKAGE)/api \
		"platform:v1 business:v1 notify:v1 registry:v1 monitor:v1 auth:v1 logagent:v1 application:v1 mesh:v1"

.PHONY: gen.gateway
gen.gateway:
	@$(DOCKER) run --rm \
		-v $(ROOT_DIR):/go/src/$(ROOT_PACKAGE) \
		-e EXT_PB_APIS=$(EXT_PB_APIS)\
	 	$(REGISTRY_PREFIX)/code-generator:$(CODE_GENERATOR_VERSION) \
	 	/root/code.sh \
	 	deepcopy-internal,deepcopy-external,defaulter-external,conversion-external \
	 	$(ROOT_PACKAGE)/pkg/gateway/apis \
	 	$(ROOT_PACKAGE)/pkg/gateway/apis \
	 	$(ROOT_PACKAGE)/pkg/gateway/apis \
	 	"config:v1"

.PHONY: gen.audit
gen.audit:
	@$(DOCKER) run --rm \
		-v $(ROOT_DIR):/go/src/$(ROOT_PACKAGE) \
	 	$(REGISTRY_PREFIX)/code-generator:$(CODE_GENERATOR_VERSION) \
	 	/root/code.sh \
	 	deepcopy-internal,deepcopy-external,defaulter-external,conversion-external \
	 	$(ROOT_PACKAGE)/pkg/audit/apis \
	 	$(ROOT_PACKAGE)/pkg/audit/apis \
	 	$(ROOT_PACKAGE)/pkg/audit/apis \
	 	"config:v1"

.PHONY: gen.registry
gen.registry:
	@$(DOCKER) run --rm \
		-v $(ROOT_DIR):/go/src/$(ROOT_PACKAGE) \
		-e EXT_PB_APIS=$(EXT_PB_APIS)\
	 	$(REGISTRY_PREFIX)/code-generator:$(CODE_GENERATOR_VERSION) \
	 	/root/code.sh \
	 	deepcopy-internal,deepcopy-external,defaulter-external,conversion-external \
	 	$(ROOT_PACKAGE)/pkg/registry/apis \
	 	$(ROOT_PACKAGE)/pkg/registry/apis \
	 	$(ROOT_PACKAGE)/pkg/registry/apis \
	 	"config:v1"

.PHONY: gen.monitor
gen.monitor:
	@$(DOCKER) run --rm \
		-v $(ROOT_DIR):/go/src/$(ROOT_PACKAGE) \
		-e EXT_PB_APIS=$(EXT_PB_APIS)\
	 	$(REGISTRY_PREFIX)/code-generator:$(CODE_GENERATOR_VERSION) \
	 	/root/code.sh \
	 	deepcopy-internal,deepcopy-external,defaulter-external,conversion-external \
	 	$(ROOT_PACKAGE)/pkg/monitor/apis \
	 	$(ROOT_PACKAGE)/pkg/monitor/apis \
	 	$(ROOT_PACKAGE)/pkg/monitor/apis \
	 	"config:v1"

.PHONY: gen.mesh
gen.mesh:
	@$(DOCKER) run --rm \
		-v $(ROOT_DIR):/go/src/$(ROOT_PACKAGE) \
		-e EXT_PB_APIS=$(EXT_PB_APIS)\
	 	$(REGISTRY_PREFIX)/code-generator:$(CODE_GENERATOR_VERSION) \
	 	/root/code.sh \
	 	deepcopy-internal,deepcopy-external,defaulter-external,conversion-external \
	 	$(ROOT_PACKAGE)/pkg/mesh/apis \
	 	$(ROOT_PACKAGE)/pkg/mesh/apis \
	 	$(ROOT_PACKAGE)/pkg/mesh/apis \
	 	"config:v1"

.PHONY: gen.openapi
gen.openapi:
	@$(DOCKER) run --rm \
    	-v $(ROOT_DIR):/go/src/$(ROOT_PACKAGE) \
		-e EXT_PB_APIS=$(EXT_PB_APIS)\
		-v $(K8S_APIMACHINERY_DIR):/go/src/k8s.io/apimachinery \
		-v $(K8S_API_DIR):/go/src/k8s.io/api \
	 	$(REGISTRY_PREFIX)/code-generator:$(CODE_GENERATOR_VERSION) \
	 	/root/openapi.sh \
	 	$(ROOT_PACKAGE)/api/openapi \
	 	$(shell ${ROOT_DIR}/build/script/openapi.sh)

.PHONY: gen.clean
gen.clean:
	@rm -rf ./api/client/{clientset,informers,listers}
	@$(FIND) -type f -name 'generated.*' -delete
	@$(FIND) -type f -name 'zz_generated*.go' -delete
	@$(FIND) -type f -name 'types_swagger_doc_generated.go' -delete

