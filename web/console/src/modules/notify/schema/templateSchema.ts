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

export const templateSchema = {
  properties: {
    apiVersion: {
      value: `${resourceConfig()['template'].group}/${resourceConfig()['template'].version}`
    },
    kind: {
      value: 'Template'
    },
    metadata: {
      properties: {
        name: TYPES.string,
        namespace: TYPES.string
      }
    },
    spec: {
      type: 'pickOne',
      pick: 'text',
      properties: {
        displayName: { ...TYPES.string, required: true },
        text: {
          properties: {
            body: { ...TYPES.string, required: true, placeholder: '请输入body', bodyTip: true },
            header: { ...TYPES.string, required: true, placeholder: '请输入消息头' }
          }
        },
        tencentCloudSMS: {
          properties: {
            body: { ...TYPES.string, required: true, placeholder: '请输入body', bodyTip: true },
            sign: { ...TYPES.string, required: true, placeholder: '请输入腾讯云短信服务签名ID', smsSignTip: true },
            templateID: {
              ...TYPES.string,
              required: true,
              placeholder: '请输入腾讯云短信服务消息模板ID',
              smsTemplateIDTip: true
            }
          }
        },
        wechat: {
          properties: {
            templateID: { ...TYPES.string, required: true, placeholder: '请输入微信公众号上创建的消息模板ID' },
            miniProgramAppID: { ...TYPES.string, placeholder: '请输入小程序AppID' },
            miniProgramPagePath: { ...TYPES.string, placeholder: '请输入小程序页面地址' },
            url: { ...TYPES.string, placeholder: '请输入消息中的跳转链接' }
          }
        },
        webhook: {
          properties: {
            body: {
              ...TYPES.string,
              required: true,
              placeholder: '请输入body,例如：{{.summary}}',
              bodyTip: true,
              value: ''
            },
            header: { ...TYPES.string, required: true, placeholder: '请输入消息头' }
          }
        }
      }
    }
  }
};
