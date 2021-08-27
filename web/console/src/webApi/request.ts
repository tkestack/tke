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
import { Method } from '@helper';
import Axios from 'axios';
import { v4 as uuidv4 } from 'uuid';

const instance = Axios.create({
  timeout: 10000
});

instance.interceptors.request.use(
  config => {
    Object.assign(config.headers, {
      'X-Remote-Extra-RequestID': uuidv4(),
      'Content-Type':
        config.method === 'patch' ? 'application/strategic-merge-patch+json' : config.headers['Content-Type']
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
    if (!error.response) {
      error.response = {
        data: {
          message: `系统内部服务错误（${error?.config?.heraders?.['X-Remote-Extra-RequestID'] || ''}）`
        }
      };
    }

    if (error.response.status === 401) {
      location.reload();
    }

    if (error.response.status === 403) {
      changeForbiddentConfig({
        isShow: true,
        message: error.response.data.message
      });
    }

    return Promise.reject(error);
  }
);

export default instance;
