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

import { TYPES } from './schemaUtil';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { resourceConfig } from '@config';

export const receiverSchema = {
  properties: {
    apiVersion: {
      value: `${resourceConfig()['receiver'].group}/${resourceConfig()['receiver'].version}`
    },
    kind: {
      value: 'Receiver'
    },
    metadata: {
      properties: {
        name: TYPES.string,
        namespace: TYPES.string
      }
    },
    spec: {
      properties: {
        displayName: { ...TYPES.string, required: true, name: t('显示名称') },
        username: { ...TYPES.string, required: true, name: t('用户名') },
        identities: {
          properties: {
            mobile: { ...TYPES.string, required: true, name: t('移动电话') },
            email: { ...TYPES.string, required: true, name: t('电子邮件') },
            wechat_openid: { ...TYPES.string, required: true, name: t('微信OpenID') }
          }
        }
      }
    }
  }
};
