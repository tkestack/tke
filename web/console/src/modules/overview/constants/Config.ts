/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

import { t } from '@tencent/tea-app/lib/i18n';

export const canNotOperateCluster = {
  clusterStatus: ['Initializing', 'Running', 'Failed', 'Upgrading', 'Terminating']
};

/** 集群状态 */
export const clusterStatus = {
  Initializing: {
    text: t('创建中'),
    classname: 'text-restart'
  },
  Running: {
    text: t('运行中'),
    classname: 'text-success'
  },
  Terminating: {
    text: t('删除中'),
    classname: 'text-restart'
  },
  Scaling: {
    text: t('规模调整中'),
    classname: 'text-restart'
  },
  Upgrading: {
    text: t('升级中'),
    classname: 'text-restart'
  },
  Failed: {
    text: t('异常'),
    classname: 'text-danger'
  },
  '-': {
    text: '-',
    classname: 'text-restart'
  }
};
