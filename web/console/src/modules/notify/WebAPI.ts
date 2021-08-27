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
import { oc } from 'ts-optchain';

import { OperationResult, QueryState, RecordSet } from '@tencent/ff-redux';
import {} from '@tencent/qcloud-redux-query';
import { t } from '@tencent/tea-app/lib/i18n';

import { reduceK8sRestfulPath } from '../../../helpers';
import { reduceNetworkRequest } from '../../../helpers/reduceNetwork';
import { RequestParams } from '../common/models';
import { Resource } from './models/Resource';

/** RESTFUL风格的请求方法 */
const Method = {
  get: 'GET',
  post: 'POST',
  patch: 'PATCH',
  delete: 'DELETE',
  put: 'PUT'
};
const tips = seajs.require('tips');

// 返回标准操作结果
function operationResult<T>(target: T[] | T, error?: any): OperationResult<T>[] {
  if (target instanceof Array) {
    return target.map(x => ({ success: !error, target: x, error }));
  }
  return [{ success: !error, target: target as T, error }];
}
/**fetchResourceList */
export async function fetchResourceList(query: QueryState<{}>, resourceInfo) {
  let { paging, search } = query;
  let url = reduceK8sRestfulPath({
    resourceInfo
  });
  if (search) {
    url += '/' + search;
  }
  let params: RequestParams = {
    method: Method.get,
    url
  };
  if (paging) {
    let { pageIndex, pageSize } = paging;
    params['page'] = pageIndex;
    params['page_size'] = pageSize;
  }

  let resourceList = [];
  try {
    let response = await reduceNetworkRequest(params);

    if (response.code === 0) {
      let listItems = response.data;
      if (listItems.items) {
        resourceList = listItems.items.map(item => {
          return Object.assign({}, item, { id: item.metadata.name });
        });
      } else {
        // 这里是拉取某个具体的resource的时候，没有items属性
        if (listItems.metadata) {
          resourceList.push({
            metadata: listItems.metadata,
            spec: listItems.spec,
            status: listItems.status
          });
        }
      }
    }
  } catch (error) {
    // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
    if (error.code !== 'ResourceNotFound') {
      throw error;
    }
  }

  const result: RecordSet<Resource> = {
    recordCount: resourceList.length,
    records: resourceList
  };

  return result;
}
/**
 * 更新Resource
 * @param strategy 新的策略
 */
export async function modifyResource([resource]: Resource[]) {
  let { mode, resourceIns, resourceInfo, jsonData, namespace, isSpecialNamespace } = resource;

  try {
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: resourceIns, namespace, isSpecialNamespace });
    const response = await reduceNetworkRequest({
      method: mode === 'create' ? 'POST' : 'PUT',
      url,
      data: jsonData
    });
    if (response.code === 0) {
      tips.success(t('修改成功'), 2000);
      return operationResult(resource);
    } else {
      return operationResult(resource, response);
    }
  } catch (e) {
    let error = oc(e).response.data.message() ? e.response.data : e;
    tips.error(error, 2000);
    return operationResult(resource, error);
  }
}

export async function deleteResource(resources: Resource[]) {
  for (let resource of resources) {
    let { resourceIns, resourceInfo, jsonData, namespace, isSpecialNamespace } = resource;

    try {
      const url = reduceK8sRestfulPath({ resourceInfo, specificName: resource.metadata.name, namespace: namespace || resource.metadata.namespace, isSpecialNamespace });
      const response = await reduceNetworkRequest({
        method: 'DELETE',
        url,
        data: jsonData
      });
      if (response.code === 0) {
        tips.success(t('删除成功'), 2000);
      } else {
        return operationResult(resource, response);
      }
    } catch (error) {
      tips.error(error, 2000);
      return operationResult(resource, error);
    }
  }
  return operationResult(resources);
}
