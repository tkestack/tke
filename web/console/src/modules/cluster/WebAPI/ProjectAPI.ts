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

import { resourceConfig } from '../../../../config/resourceConfig';
import { reduceK8sRestfulPath } from '../../../../helpers';
import { Method, reduceNetworkRequest } from '../../../../helpers/reduceNetwork';
import { RequestParams, Resource } from '../../../modules/common';
import { ResourceInfo } from '../../common/models/ResourceInfo';
import { ResourceFilter } from '../models';

//业务控制台api

/**
 * 获取资源的具体的 yaml文件
 * @param resourceIns: Resource[]   当前需要请求的具体资源数据
 * @param resourceInfo: ResouceInfo 当前请求数据url的基本配置
 */
export async function fetchUserPortal(resourceInfo: ResourceInfo) {
  let url = reduceK8sRestfulPath({ resourceInfo });

  // 构建参数
  let params: RequestParams = {
    method: Method.get,
    url
  };

  let response = await reduceNetworkRequest(params);
  return response.data;
}

/**
 * Namespace查询
 * @param query Namespace查询的一些过滤条件
 */
export async function fetchProjectNamespaceList(query: QueryState<ResourceFilter>) {
  let { filter } = query;
  let NamespaceResourceInfo: ResourceInfo = resourceConfig().namespaces;
  let url = reduceK8sRestfulPath({
    resourceInfo: NamespaceResourceInfo,
    specificName: filter.specificName,
    extraResource: 'namespaces'
  });
  /** 构建参数 */
  let method = 'GET';
  let params: RequestParams = {
    method,
    url
  };

  let response = await reduceNetworkRequest(params);
  let namespaceList = [],
    total = 0;
  if (response.code === 0) {
    let list = response.data;
    total = list.items.length;
    namespaceList = list.items.map(item => {
      return Object.assign({}, item, { id: uuid(), name: item.metadata.name });
    });
  }

  const result: RecordSet<Resource> = {
    recordCount: total,
    records: namespaceList
  };

  return result;
}
