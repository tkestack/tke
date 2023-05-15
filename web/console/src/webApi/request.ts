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
import { changeForbiddentConfig } from '@/index.tke';
import { createCSRFHeader } from '@helper';
import Axios from 'axios';
import { v4 as uuidv4 } from 'uuid';
import { message } from 'tea-component';

const instance = Axios.create({
  timeout: 10000
});

// strategic-

instance.interceptors.request.use(
  config => {
    Object.assign(config.headers, {
      'X-Remote-Extra-RequestID': uuidv4(),
      ...createCSRFHeader()
    });
    return config;
  },
  error => {
    console.log('request error:', error);
    return Promise.reject(error);
  }
);

instance.interceptors.response.use(
  ({ data }) => data,
  error => {
    console.error('response error:', error);

    const errorMessage =
      error?.response?.data?.message ??
      `系统内部服务错误（${error?.config?.heraders?.['X-Remote-Extra-RequestID'] ?? ''}）`;

    switch (error?.response?.status) {
      case 401:
        location.reload();
        break;
      case 403:
        changeForbiddentConfig({
          isShow: true,
          message: errorMessage
        });
        break;
      case 404:
        // 404不一定要展示错误
        break;
      default:
        message.error({ content: errorMessage });
    }

    return Promise.reject(error);
  }
);

export default instance;

export const Request = instance;

export const generateQueryString = (query: Record<string, any>, joinKey = '&') => {
  return Object.entries(query)
    .filter(([_, value]) => value !== undefined && value !== null && value !== '')
    .map(([key, value]) => `${key}=${value}`)
    .join(joinKey);
};
