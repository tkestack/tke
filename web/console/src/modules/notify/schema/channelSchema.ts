/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
import { resourceConfig } from '@config';
import validatorjs from 'validator';

export const channelSchema = {
  properties: {
    apiVersion: {
      value: `${resourceConfig()['channel'].group}/${resourceConfig()['channel'].version}`
    },
    kind: {
      value: 'Channel'
    },
    metadata: {
      properties: {
        name: TYPES.string,
        namespace: TYPES.string
      }
    },
    spec: {
      type: 'pickOne',
      pick: 'smtp',
      properties: {
        displayName: { ...TYPES.string, required: true },
        smtp: {
          properties: {
            email: { ...TYPES.string, required: true },
            password: { ...TYPES.string, required: true },
            smtpHost: { ...TYPES.string, required: true },
            smtpPort: { ...TYPES.number, required: true },
            tls: TYPES.boolean
          }
        },
        tencentCloudSMS: {
          properties: {
            appKey: { ...TYPES.string, required: true },
            sdkAppID: { ...TYPES.string, required: true },
            extend: TYPES.string
          }
        },
        wechat: {
          properties: {
            appID: { ...TYPES.string, required: true },
            appSecret: { ...TYPES.string, required: true }
          }
        },
        webhook: {
          properties: {
            url: {
              ...TYPES.string,
              required: true,
              validator(value) {
                if (value === undefined) return {};

                if (value === '') return { status: 'error', message: 'url为必填' };

                if (!validatorjs.isURL(value)) return { status: 'error', message: 'url格式有误' };

                return {
                  status: 'success',
                  message: ''
                };
              }
            },
            headers: {
              ...TYPES.string,
              placeholder: '自定义Header，仅支持Key:Value格式，中间用;号分割。eg param1:1;param2:2'
            }
          }
        }
      }
    }
  }
};
