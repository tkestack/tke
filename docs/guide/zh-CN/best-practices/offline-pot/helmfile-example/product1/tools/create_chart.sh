#!/bin/bash

# Tencent is pleased to support the open source community by making TKEStack
# available.
#
# Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

app=$1

if [ ! -n "$1" ]; then
    echo "Usage: ./tools/create_chart.sh product1-new"
    exit
fi

if [ -d "$1" ]; then
    echo "$1 directory exists"
    exit
fi

cp -R ./tools/product1-demo ./${app}

sed -i "s/product1-demo/${app}/g" ${app}/Chart.yaml ${app}/values/ty.yaml
