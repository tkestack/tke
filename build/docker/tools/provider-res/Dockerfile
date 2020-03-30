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

COPY linux-amd64/cni-*          res/linux-amd64/
COPY linux-arm64/cni-*          res/linux-arm64/

COPY linux-amd64/docker-*       res/linux-amd64/
COPY linux-arm64/docker-*       res/linux-arm64/

COPY linux-amd64/kubernetes-*   res/linux-amd64/
COPY linux-arm64/kubernetes-*   res/linux-arm64/

COPY linux-amd64/kubeadm-*      res/linux-amd64/
COPY linux-arm64/kubeadm-*      res/linux-arm64/

COPY linux-amd64/NVIDIA-*       res/linux-amd64/

COPY linux-amd64/nvidia-container-* res/linux-amd64/

ENTRYPOINT ["sh"]
