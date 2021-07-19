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

FROM alpine:3.14.0

RUN echo "hosts: files dns" >> /etc/nsswitch.conf

WORKDIR /app

ENV PATH="/app/bin:$PATH"
ENV DOCKER_CLI_EXPERIMENTAL=enabled
ARG ENV_ARCH

RUN apk add --no-cache \
    bash \
    busybox-extras \
    curl \
    tcpdump \
    docker \
    ansible \
    containerd

RUN apk --update-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ --allow-untrusted add lrzsz
RUN wget -O nerdctl-0.10.0-linux.tar.gz https://github.com/containerd/nerdctl/releases/download/v0.10.0/nerdctl-0.10.0-linux-"$ENV_ARCH".tar.gz \
    && tar -zvxf nerdctl-0.10.0-linux.tar.gz -C /usr/local/bin/ \
    && rm -rf nerdctl-0.10.0-linux.tar.gz \
    && wget -O crictl-v1.20.0-linux.tar.gz https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.20.0/crictl-v1.20.0-linux-"$ENV_ARCH".tar.gz \
    && tar -zvxf crictl-v1.20.0-linux.tar.gz -C /usr/local/bin/ \
    && rm -rf crictl-v1.20.0-linux.tar.gz

ADD . /app

RUN ln -s /app/.docker /root/.docker

ENTRYPOINT ["/app/bin/tke-installer"]
