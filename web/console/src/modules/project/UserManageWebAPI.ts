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
import { QueryState, RecordSet } from '@tencent/ff-redux';
import {
  Method,
  operationResult,
  reduceK8sQueryString,
  reduceK8sRestfulPath,
  reduceNetworkRequest
} from '../../../helpers';
import { ResourceInfo, RequestParams } from '../common/models';
import { resourceConfig } from '../../../config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { METHODS } from 'http';
import { User, UserFilter, PolicyFilter, PolicyPlain } from './models';

const tips = seajs.require('tips');

class RequestResult {
  data: any;
  error: any;
}
const SEND = async (url: string, method: string, bodyData: any, tipErr = true) => {
  // 构建参数
  const params: RequestParams = {
    method: method,
    url,
    data: bodyData
  };
  const resp = new RequestResult();
  try {
    const response = await reduceNetworkRequest(params);
    if (response.code !== 0) {
      if (tipErr === true) {
        tips.error(response.message, 2000);
      }
      resp.data = bodyData;
      resp.error = response.message;
    } else {
      if (method !== Method.get) {
        tips.success('操作成功', 2000);
      }
      resp.data = response.data;
      resp.error = null;
    }
    return resp;
  } catch (error) {
    if (tipErr === true) {
      tips.error(error.response.data.message, 2000);
    }
    resp.data = bodyData;
    resp.error = error.response.data.message;
    return resp;
  }
};

const GET = async (url: string, tipErr = true) => {
  const response = await SEND(url, Method.get, null, tipErr);
  return response;
};
const DELETE = async (url: string, tipErr = true) => {
  const response = await SEND(url, Method.delete, null, tipErr);
  return response;
};
const POST = async (url: string, bodyData: any, tipErr = true) => {
  const response = await SEND(url, Method.post, JSON.stringify(bodyData), tipErr);
  return response;
};

const PUT = async (url: string, bodyData: any, tipErr = true) => {
  const response = await SEND(url, Method.put, JSON.stringify(bodyData), tipErr);
  return response;
};

const PATCH = async (url: string, bodyData: any, tipErr = true) => {
  const response = await SEND(url, Method.patch, JSON.stringify(bodyData), tipErr);
  return response;
};

/**
 * 用户列表的查询
 * @param query 列表查询条件参数
 */
export async function fetchUserList(query: QueryState<UserFilter>) {
  let users: User[] = [];
  const { search, filter } = query;
  const { projectId } = filter;
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['members'];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: projectId, extraResource: 'users' });
    const response = await reduceNetworkRequest(
      {
        method: Method.get,
        url
      },
      '',
      projectId,
      search
    );

    if (response.data.items) {
      users = response.data.items;
    } else {
      users = [];
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
  }
  const result: RecordSet<User> = {
    recordCount: users.length,
    records: users
  };
  return result;
}

/**
 * 增加用户
 * @param [userInfo] 用户数据, 这里和actions.user.addUser.start([userInfo]);的对应
 */
export async function addUser([userInfo]) {
  const { projectId, users, policies } = userInfo;
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['members'];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: projectId, extraResource: 'users' });
    const response = await reduceNetworkRequest({
      method: Method.post,
      url,
      data: {
        users,
        policies
      }
    });
    if (response.code === 0) {
      tips.success(t('操作成功'), 2000);
      return operationResult(userInfo);
    } else {
      // 是否给tip得看具体返回的数据
      return operationResult(userInfo, response);
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
    // 返回相关数据存储在redux中, 这里的error应该不用reduceNetworkWorkflow作数据处理
    return operationResult(userInfo, error.response);
  }
}

/**
 * 策略列表的查询，不参杂其他场景参数
 * @param query 列表查询条件参数
 */
export async function fetchPolicyPlainList(query: QueryState<PolicyFilter>) {
  const { search, filter, keyword } = query;
  let queryString = '';
  if (filter.resource === 'platform') {
    queryString = '?fieldSelector=spec.scope!=project';
  } else if (filter.resource === 'project') {
    queryString = '?fieldSelector=spec.scope=project';
  }

  const resourceInfo: ResourceInfo = resourceConfig()['policy'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  // const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  console.log('fetchPolicyPlainList url + queryString', url, queryString, 111, query);
  const rr: RequestResult = await GET(url + queryString);
  const items: PolicyPlain[] =
    !rr.error && rr.data.items
      ? rr.data.items.map(i => {
          return {
            id: i.metadata && i.metadata.name,
            name: i.metadata && i.metadata.name,
            displayName: i.spec && i.spec.displayName,
            category: i.spec && i.spec.category,
            description: i.spec && i.spec.description,
            tenantID: i.sepc && i.spec.tenantID
          };
        })
      : [];
  console.log('fetchPolicyPlainList items is:', items);
  const result: RecordSet<PolicyPlain> = {
    recordCount: items.length,
    records: items
  };
  return result;
}
