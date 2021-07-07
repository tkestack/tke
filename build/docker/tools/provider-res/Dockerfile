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

WORKDIR /data

## Make sure all files are end with .tar.gz

COPY linux-amd64/cni-*.tar.gz          res/linux-amd64/
COPY linux-arm64/cni-*.tar.gz          res/linux-arm64/

COPY linux-amd64/conntrack-tools-*.tar.gz          res/linux-amd64/
COPY linux-arm64/conntrack-tools-*.tar.gz          res/linux-arm64/

COPY linux-amd64/nerdctl-*.tar.gz res/linux-amd64/
COPY linux-arm64/nerdctl-*.tar.gz res/linux-arm64/

COPY linux-amd64/docker-*.tar.gz       res/linux-amd64/
COPY linux-arm64/docker-*.tar.gz       res/linux-arm64/

COPY linux-amd64/containerd-*.tar.gz       res/linux-amd64/
COPY linux-arm64/containerd-*.tar.gz       res/linux-arm64/

COPY linux-amd64/critools-*.tar.gz       res/linux-amd64/
COPY linux-arm64/critools-*.tar.gz       res/linux-arm64/

COPY linux-amd64/kubernetes-*.tar.gz   res/linux-amd64/
COPY linux-arm64/kubernetes-*.tar.gz   res/linux-arm64/

COPY linux-amd64/NVIDIA-*.tar.gz       res/linux-amd64/

COPY linux-amd64/nvidia-container-*.tar.gz res/linux-amd64/

COPY containerd conf/containerd

ENTRYPOINT ["sh"]
