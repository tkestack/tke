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

import { QueryState, RecordSet, uuid } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { apiServerVersion } from '../../../../config';
import {
  Method,
  operationResult,
  reduceK8sQueryString,
  reduceK8sRestfulPath,
  reduceNetworkRequest,
  reduceNetworkWorkflow,
  requestMethodForAction
} from '../../../../helpers';
import {
  Cluster,
  CreateResource,
  RequestParams,
  Resource,
  ResourceFilter,
  ResourceInfo,
  UserDefinedHeader
} from '../models';

export async function createPromethus(resource: CreateResource) {
  try {
    let { mode, resourceIns, clusterId, yamlData, resourceInfo, namespace, jsonData } = resource;

    let url = '/apis/monitor.tkestack.io/v1/prometheuses';
    // 获取具体的请求方法，create为POST，modify为PUT
    let method = requestMethodForAction(mode);
    // 这里是独立部署版 和 控制台共用的参数，只有是yamlData的时候才需要userdefinedHeader，如果是jaonData的话，就不需要了
    let userDefinedHeader: UserDefinedHeader = yamlData
      ? {
          Accept: 'application/json',
          'Content-Type': 'application/yaml'
        }
      : {};

    // 构建参数
    let params: RequestParams = {
      method,
      url,
      userDefinedHeader,
      data: yamlData ? yamlData : jsonData
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

export async function deletePromethus(cluster: Cluster) {
  try {
    const clusterId = cluster.metadata.name;
    const url = cluster.spec.promethus.metadata.selfLink;
    // 构建参数 requestBody 当中
    let params: RequestParams = {
      method: Method.delete,
      url
    };

    let response = await reduceNetworkRequest(params, clusterId);

    if (response.code === 0) {
      return Promise.resolve();
    } else {
      return Promise.reject(response.code);
    }
  } catch (error) {
    return Promise.reject(error);
  }
}
