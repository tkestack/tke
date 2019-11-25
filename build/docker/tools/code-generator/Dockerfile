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

FROM golang:1.12.12-alpine3.10

RUN apk add --update --no-cache autoconf automake libtool zlib-dev build-base curl git

# Install the protobuf
ENV PROTOBUF_VERSION 3.10.0
ENV PROTOBUF_URL=https://github.com/protocolbuffers/protobuf/archive/v${PROTOBUF_VERSION}.tar.gz
RUN curl -L ${PROTOBUF_URL} -o - | tar -xz -C /tmp
RUN cd /tmp/protobuf-* && \
    ./autogen.sh && \
    ./configure --disable-shared --enable-static && \
    make --silent -j `nproc` install-strip
RUN cd /tmp && rm -rf protobuf-$PROTOBUF_VERSION \
    && protoc --version

ENV GOPATH=/go

# Clone the kubernetes
ENV KUBERNETES_VERSION 1.16.3
RUN mkdir -p $GOPATH/src/k8s.io \
  && git clone -b v$KUBERNETES_VERSION --single-branch https://github.com/kubernetes/kubernetes.git $GOPATH/src/k8s.io/kubernetes
RUN apk add --update --no-cache rsync file patch
RUN cd $GOPATH/src/k8s.io/kubernetes \
  && make WHAT="cmd/genswaggertypedocs vendor/k8s.io/code-generator/cmd/defaulter-gen vendor/k8s.io/code-generator/cmd/client-gen vendor/k8s.io/code-generator/cmd/lister-gen vendor/k8s.io/code-generator/cmd/informer-gen vendor/k8s.io/code-generator/cmd/deepcopy-gen vendor/k8s.io/code-generator/cmd/conversion-gen vendor/k8s.io/code-generator/cmd/go-to-protobuf vendor/k8s.io/code-generator/cmd/go-to-protobuf/protoc-gen-gogo"
RUN ln -s $GOPATH/src/k8s.io/kubernetes/vendor/k8s.io/apimachinery $GOPATH/src/k8s.io/apimachinery \
  && ln -s $GOPATH/src/k8s.io/kubernetes/vendor/k8s.io/api $GOPATH/src/k8s.io/api

# work around 64MB tmpfs size in Docker 1.6
ENV TMPDIR /tmp.k8s
RUN mkdir $TMPDIR \
  && chmod a+rwx $TMPDIR \
  && chmod o+t $TMPDIR

# Get the code coverage tool and goimports
RUN go get golang.org/x/tools/cmd/cover \
           golang.org/x/tools/cmd/goimports \
    && go clean -cache

WORKDIR "$GOPATH/src/k8s.io/kubernetes"
ADD code.sh /root/
ADD openapi.sh /root/
ADD boilerplate.go.txt /root/
RUN chmod +x /root/code.sh
RUN chmod +x /root/openapi.sh

