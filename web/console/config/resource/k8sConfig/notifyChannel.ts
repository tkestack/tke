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

import { generateResourceInfo } from '../common';

/** notify的相关配置 */
export const notifyChannel = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'channel',
    requestType: {
      list: 'channels'
    }
  });
};

export const notifyTemplate = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'template',
    isRelevantToNamespace: true,
    requestType: {
      list: 'templates'
    }
  });
};

export const notifyMessage = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'message',
    requestType: {
      list: 'messages'
    }
  });
};
export const notifyReceiver = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'receiver',
    requestType: {
      list: 'receivers'
    }
  });
};
export const notifyReceiverGroup = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'receiverGroup',
    requestType: {
      list: 'receivergroups'
    }
  });
};
