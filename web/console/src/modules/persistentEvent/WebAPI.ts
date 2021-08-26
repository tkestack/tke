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

import { RequestParams, UserDefinedHeader } from '../common/models';
import {
  reduceNetworkRequest,
  reduceNetworkWorkflow,
  requestMethodForAction,
  reduceK8sRestfulPath,
  operationResult
} from '../../../helpers';
import { CreateResource } from './models';

/**
 * 设置集群持久化事件
 * @param resource: CreateResource 创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function modifyPeConfig(resource: CreateResource[], regionId: number) {
  try {
    let { mode, resourceInfo, clusterId, jsonData, namespace, resourceIns } = resource[0];
    let url = reduceK8sRestfulPath({ resourceInfo, namespace });
    let userDefinedHeader: UserDefinedHeader = {};

    if (mode === 'update') {
      url += `/${resourceIns}`;
      userDefinedHeader = {
        'Content-Type': 'application/strategic-merge-patch+json'
      };
    }

    let method = requestMethodForAction(mode);

    // 构建参数
    let params: RequestParams = {
      method,
      url,
      userDefinedHeader,
      data: jsonData
    };

    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      return operationResult(resource);
    } else {
      return operationResult(resource, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}
