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

import { resourceConfig } from '@config';
import { RecordSet, uuid } from '@tencent/ff-redux';
import {
  Method,
  operationResult,
  reduceK8sRestfulPath,
  reduceNetworkRequest,
  reduceNetworkWorkflow
} from '../../../../helpers';
import { CreateResource, RequestParams } from '../../common/models';
import { isEmpty } from '../../common';
import { Resource } from '../models';

/**
 * 获取集群日志组件列表
 */
export async function fetchLogagents() {
  let resourceList = [];
  let resourceInfo = resourceConfig()['logagent'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  // 构建参数
  let params: RequestParams = {
    method: Method.get,
    url,
  };

  let response = await reduceNetworkRequest(params);

  if (response.code === 0) {
    const { items } = response.data;
    if (!isEmpty(items)) {
      resourceList = items.map(item => {
        return Object.assign({}, item, { id: uuid() });
      });
    }
  }
  const result: RecordSet<Resource> = {
    recordCount: resourceList.length,
    records: resourceList
  };

  return result;
}

/**
 * 创建集群新日志组件
 */
export async function createLogAgent(resources: CreateResource, tenantID: string = 'default') {
  try {
    let { clusterId, resourceInfo } = resources;
    let { group, version, headTitle } = resourceInfo;
    let url = reduceK8sRestfulPath({ resourceInfo, clusterId });
    let params: RequestParams = {
      method: Method.post,
      url,
      data: {
        kind: headTitle,
        apiVersion: `${group}/${version}`,
        metadata: {
          name: ''
        },
        spec: {
          tenantID,
          clusterName: clusterId,
          version: 'v1.0.0'
        }
      }
    };

    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      return operationResult(resources);
    } else {
      return operationResult(resources, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resources, reduceNetworkWorkflow(error));
  }
}

/**
 * 删除集群 LogAgent 组件
 */
export async function deleteLogAgent(resource: CreateResource, logAgentName: string) {
  try {
    let { resourceInfo } = resource;

    let url = reduceK8sRestfulPath({ resourceInfo, specificName: logAgentName });

    // 构建参数 requestBody 当中
    let params: RequestParams = {
      method: Method.delete,
      url
    };

    let response = await reduceNetworkRequest(params);

    if (response.code === 0) {
      return operationResult(resource);
    } else {
      return operationResult(resource, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}

