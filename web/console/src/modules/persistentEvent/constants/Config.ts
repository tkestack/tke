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

import { t, Trans } from '@tencent/tea-app/lib/i18n';

/** ========================= start FFRedux的相关配置 ======================== */
export const FFReduxActionName = {
  REGION: 'region',
  CLUSTER: 'cluster'
};

/** 事件轮询的列表 */
export const PollEventName = {
  peList: 'pollPeList'
};

/** 是否需要轮询事件持久列表 */
export const isNeedPollPE = ['initializing', 'reinitializing', 'checking', 'failed'];

/** 事件持久化状态 */
export const peStatus = {
  initializing: {
    text: t('初始化'),
    classname: 'text-restart'
  },
  reinitializing: {
    text: t('重新初始化'),
    classname: 'text-restart'
  },
  running: {
    text: t('已开启'),
    classname: 'text-success'
  },
  failed: {
    text: t('失败'),
    classname: 'text-danger'
  },
  checking: {
    text: t('检查中'),
    classname: 'text-restart'
  },
  '-': {
    text: '-',
    classname: 'text-restart'
  }
};
