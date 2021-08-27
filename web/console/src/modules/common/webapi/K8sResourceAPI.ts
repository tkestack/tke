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
import { QueryState, RecordSet, uuid } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { apiServerVersion } from '../../../../config';
import {
    Method, operationResult, reduceK8sQueryString, reduceK8sRestfulPath, reduceNetworkRequest,
    reduceNetworkWorkflow, requestMethodForAction
} from '../../../../helpers';
import {
    CreateResource, RequestParams, Resource, ResourceFilter, ResourceInfo, UserDefinedHeader
} from '../models';

const tips = seajs.require('tips');

/**
 * Resource列表的查询
 * @param query:    Resource 的查询过滤条件
 * @param resourceInfo:ResourceInfo 资源的相关配置
 * @param isClearData:  是否清空数据
 * @param k8sQueryObj: any  是否有queryString
 * @param isNeedDes: boolean    是否需要降序展示
 */
export async function fetchResourceList<T = Resource>(options: {
  query: QueryState<ResourceFilter>;
  resourceInfo: ResourceInfo;
  isClearData?: boolean;
  k8sQueryObj?: any;
}): Promise<RecordSet<T>> {
  let { query, resourceInfo, k8sQueryObj = {}, isClearData = false } = options;

  let { filter, search } = query,
    { namespace, regionId, clusterId, specificName, logAgentName = '' } = filter;

  let resourceList = [];

  if (!isClearData) {
    let k8sUrl = '';
    // 如果有搜索字段的话
    if (search) {
      k8sUrl = reduceK8sRestfulPath({ resourceInfo, namespace, specificName: search, clusterId, logAgentName });
    } else {
      k8sUrl = reduceK8sRestfulPath({ resourceInfo, namespace, specificName, clusterId, logAgentName });
    }

    // 这里去拼接，是否需要在k8sUrl后面拼接一些queryString
    let queryString = reduceK8sQueryString({ k8sQueryObj });
    let url = k8sUrl + queryString;

    let params: RequestParams = {
      method: Method.get,
      url
    };

    try {
      let response = await reduceNetworkRequest(params, clusterId);

      if (response.code === 0) {
        let listItems = response.data;
        if (listItems.items) {
          resourceList = listItems.items.map(item => {
            return Object.assign({}, item, { id: uuid() });
          });
        } else {
          // 这里是拉取某个具体的resource的时候，没有items属性
          resourceList.push({
            metadata: listItems.metadata,
            spec: listItems.spec,
            status: listItems.status
          });
        }
      }
    } catch (error) {
      // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
      if (+error.response.status !== 404) {
        throw error;
      }
    }
  }

  const result: RecordSet<T> = {
    recordCount: resourceList.length,
    records: resourceList
  };
  return result;
}

/**
 * 拉取某个具体资源下的额外资源，如 event、pod等
 * @param query: ResourceFilter 的查询过滤条件
 * @param resourceInfo: ResourceInfo 资源的具体配置
 * @param isClearData: boolean 是否清除数据
 * @param extraResource: string 额外的资源，如event
 * @param k8sQueryObj: any  是否有queryString
 * @param isNeedDes: boolean    是否需要降序展示
 */
export async function fetchExtraResourceList<T = Resource>(options: {
  query: QueryState<ResourceFilter>;
  resourceInfo: ResourceInfo;
  isClearData?: boolean;
  extraResource?: string;
  k8sQueryObj?: any;
  isNeedDes?: boolean;
}): Promise<RecordSet<T>> {
  let { query, resourceInfo, isClearData = false, extraResource = '', k8sQueryObj = {}, isNeedDes = false } = options;
  let { filter } = query,
    { namespace, clusterId, regionId, specificName } = filter;

  let extraResourceList = [];

  if (!isClearData) {
    let k8sUrl = reduceK8sRestfulPath({ resourceInfo, namespace, specificName, extraResource, clusterId });
    // 这里是去拼接，是否需要在k8s url后面拼接一些queryString
    let queryString = reduceK8sQueryString({ k8sQueryObj });
    // 这里是拼接查询的 queryString
    let url = k8sUrl + queryString;

    // 构建参数
    let params: RequestParams = {
      method: Method.get,
      url
    };

    let response = await reduceNetworkRequest(params, clusterId);

    if (response.code === 0) {
      let listItems = response.data;
      if (listItems.items) {
        extraResourceList = listItems.items.map(item => {
          return Object.assign({}, item, { id: uuid() });
        });
      }
    }
  }

  const result: RecordSet<T> = {
    recordCount: extraResourceList.length,
    records: isNeedDes && extraResourceList.length ? extraResourceList.reverse() : extraResourceList
  };

  return result;
}

/**
 * 创建ResourceIns
 * @param resource: CreateResource   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function modifyResourceIns(resource: CreateResource[], regionId: number) {
  try {
    let { mode, resourceIns, clusterId, yamlData, resourceInfo, namespace, jsonData } = resource[0];

    let url = reduceK8sRestfulPath({ resourceInfo, namespace, specificName: resourceIns, clusterId });
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

/**
 * 同时创建多种资源
 * @param resource: CreateResource 创建resourceIns的相关信息
 * @param regionId: number 地域的ID
 */
export async function applyResourceIns(resource: CreateResource[], regionId: number) {
  try {
    let { clusterId, yamlData, jsonData } = resource[0];

    let url = `/${apiServerVersion.basicUrl}/${apiServerVersion.group}/${apiServerVersion.version}/clusters/${clusterId}/apply`;

    // 这里是独立部署版 和 控制台共用的参数，只有是yamlData的时候才需要userdefinedHeader，如果是jaonData的话，就不需要了
    let userDefinedHeader: UserDefinedHeader = yamlData
      ? {
          Accept: 'application/json',
          'Content-Type': 'application/yaml'
        }
      : {};

    // 构建参数
    let params: RequestParams = {
      method: Method.post,
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

/**
 * 删除ResourceIns
 * @param resource: CreateResource  创建resourceIns的相关信息
 * @param regionId: number  地域的id
 */
export async function deleteResourceIns(resource: CreateResource[], regionId: number) {
  console.log('delete---1', resource);
  try {
    let { resourceIns, clusterId, resourceInfo, namespace, mode, logAgentName = '' } = resource[0];

    let k8sUrl = reduceK8sRestfulPath({ resourceInfo, namespace, specificName: resourceIns, clusterId, logAgentName });
    let url = k8sUrl;

    console.log('delete---', k8sUrl);

    // 是用于后台去异步的删除resource当中的pod
    let extraParamsForDelete = {
      propagationPolicy: 'Background'
    };
    if (resourceInfo.headTitle === 'Namespace') {
      extraParamsForDelete['gracePeriodSeconds'] = 0;
    }

    // 构建参数 requestBody 当中
    let params: RequestParams = {
      method: Method.delete,
      url
    };

    let response = await reduceNetworkRequest(params, clusterId);

    if (response.code === 0) {
      // 更新页面删除东西，不要告诉别人删除了东西，会造成恐慌
      mode !== 'update' && tips.success(t('删除成功'), 2000);
      return operationResult(resource);
    } else {
      return operationResult(resource, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}

/**
 * 更新某个具体的资源
 * @param resource: CreateResource   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function updateResourceIns(resource: CreateResource[], regionId: number) {
  try {
    let { resourceIns, clusterId, resourceInfo, namespace, jsonData, isStrategic = true } = resource[0];

    let url = reduceK8sRestfulPath({ resourceInfo, namespace, specificName: resourceIns, clusterId });
    let params: RequestParams = {
      method: Method.patch,
      url,
      userDefinedHeader: {
        'Content-Type': isStrategic ? 'application/strategic-merge-patch+json' : 'application/merge-patch+json'
      },
      data: jsonData,
      apiParams: {
        module: 'tke',
        interfaceName: 'ForwardRequest',
        regionId,
        restParams: {
          Method: Method.patch,
          Path: url,
          Version: '2018-05-25',
          RequestBody: jsonData
        }
      }
    };

    let response = await reduceNetworkRequest(params, clusterId);

    if (response.code === 0) {
      tips.success(t('更新成功'), 2000);
      return operationResult(resource);
    } else {
      return operationResult(resource, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}
