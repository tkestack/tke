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

FROM alpine:3.10

RUN echo "hosts: files dns" >> /etc/nsswitch.conf

WORKDIR /app

ENV PATH="/app/bin:$PATH"

RUN apk add --no-cache \
    bash \
    busybox-extras \
    curl \
    tcpdump \
    docker \
    ansible

RUN apk --update-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ --allow-untrusted add lrzsz

RUN wget https://get.helm.sh/helm-v3.2.4-linux-amd64.tar.gz && \
    tar xvzf helm-v3.2.4-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/local/bin &&  \
    rm -rf helm-v3.2.4-linux-amd64.tar.gz linux-amd64

ADD . /app

RUN ln -s /app/.docker /root/.docker

ENTRYPOINT ["/app/bin/tke-installer"]
