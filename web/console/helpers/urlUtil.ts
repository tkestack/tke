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
import { remove, isEmpty } from '../src/modules/common/utils';
import { ResourceInfo } from '../src/modules/common';
import { resourceConfig } from '../config';
import { isObject } from 'lodash';

export function parseQueryString(str = '') {
  const result = {};
  str
    .replace(/^\?*/, '')
    .split('&')
    .forEach(item => {
      const keyVal = item.split('=');
      if (keyVal.length > 0) {
        const key = decodeURIComponent(keyVal[0]);
        result[key] = keyVal[1] ? decodeURIComponent(keyVal[1]) : '';
      }
    });
  return result;
}

export function buildQueryString(obj: any = {}) {
  const keys = remove(Object.keys(obj), value => value === '');
  const queryStr = keys.map(key => `${encodeURIComponent(key)}=${encodeURIComponent(obj[key])}`).join('&');

  if (queryStr) {
    return '?' + queryStr;
  } else {
    return '';
  }
}

/**
 * 用于获取queryString
 * @param k8sQueryObj
 *  eg: ?fieldSelector=involvedObject.name=*,involvedObject.kind=*&limit=1
 *  传入进来的结构
 *  {
 *      fieldSelector: {
 *          involvedObject.name: *,
 *          involvedObject.kind: *
 *      },
 *      limit: 1
 *  }
 * @param options: K8sRestfulPathOptions
 */
export const reduceK8sQueryString = ({
  k8sQueryObj = {},
  restfulPath = ''
}: {
  k8sQueryObj: Record<string, string | number | boolean | Record<string, string | number | boolean>>;
  restfulPath?: string;
}) => {
  const queryString = Object.entries(k8sQueryObj)
    .filter(([_, value]) => value !== undefined)
    .map(([key, value]) => {
      // 也许value是object，即labelSelector 或者 fieldSelector
      value = isObject(value)
        ? Object.entries(value)
            .filter(([_, value]) => value)
            .map(([key, value]) => `${key}=${value}`)
            .join(',')
        : value;

      return `${key}=${encodeURIComponent(`${value}`)}`;
    })
    .join('&');

  const preFix = restfulPath.includes('?') ? '&' : '?';

  return queryString ? `${preFix}${queryString}` : '';
};

interface K8sRestfulPathOptions {
  /** 资源的配置 */
  resourceInfo: ResourceInfo;

  /** 命名空间，具体的ns */
  namespace?: string;

  /** 业务视图是否切分namespace */
  isSpecialNamespace?: boolean;

  /** 不在路径最后的变量，比如projectId*/
  middleKey?: string;

  /** 具体的资源名称 */
  specificName?: string;

  /** 某个具体资源下的子资源，eg: deployment/---/pos */
  extraResource?: string;

  /** 集群id，适用于addon 请求平台转发的场景 */
  clusterId?: string;

  /** 集群logAgentName */
  logAgentName?: string;

  meshId?: string;
}

/**
 * 获取k8s 的restful 风格的path
 * @param resourceInfo: ResourceInfo  资源的配置
 * @param namespace: string 具体的命名空间
 * @param specificName: string  具体的资源名称
 * @param extraResource: string 某个具体资源下的子资源
 * @param clusterId: string 集群id，适用于addon 请求平台转发的场景
 */
export const reduceK8sRestfulPath = (options: K8sRestfulPathOptions) => {
  let {
    resourceInfo,
    namespace = '',
    isSpecialNamespace = false,
    specificName = '',
    extraResource = '',
    clusterId = '',
    meshId = '',
    logAgentName = ''
  } = options;

  namespace = namespace.replace(new RegExp(`^${clusterId}-`), '');
  let url = '';
  const isAddon = resourceInfo.requestType && resourceInfo.requestType.addon ? resourceInfo.requestType.addon : false;

  /**
   * addon 和 非 addon的资源，请求的url 不太一样
   * addon:
   *  1. 如果包含有extraResource（目前仅支持 events 和 pods）
   *     => /apis/platfor.tke/v1/clusters/${cluster}/${addonNames}${extraResource}?message={namespace}&reason={specificName}
   *  2. 如果不包含extraResource
   *     => /apis/platfor.tke/v1/clusters/${cluster}/${addonNames}?namespace={namespace}&name={specificName}
   *
   * 非addon（以deployment为例):  /apis/apps/v1beta2/namespaces/${namespace}/deployments/${deployment}/${extraResource}
   */
  if (isAddon) {
    // 兼容新旧日志组件
    const baseInfo: ResourceInfo = resourceConfig()[logAgentName ? 'logagent' : 'cluster'];
    const baseValue = logAgentName || clusterId;
    url = `/${baseInfo.basicEntry}/${baseInfo.group}/${baseInfo.version}/${baseInfo.requestType['list']}/${baseValue}/${resourceInfo.requestType['list']}`;

    if (extraResource || resourceInfo['namespaces'] || specificName) {
      const queryArr: string[] = [];
      resourceInfo.namespaces && namespace && queryArr.push(`namespace=${namespace}`);
      specificName && queryArr.push(`name=${specificName}`);
      extraResource && queryArr.push(`action=${extraResource}`);
      url += `?${queryArr.join('&')}`;
    }
  } else {
    url =
      `/${resourceInfo.basicEntry}/` +
      (resourceInfo.group ? `${resourceInfo.group}/` : '') +
      `${resourceInfo.version}/` +
      (resourceInfo.namespaces ? `${resourceInfo.namespaces}/${namespace}/` : '') +
      `${resourceInfo.requestType.list}` +
      (specificName ? `/${specificName}` : '') +
      (extraResource ? `/${extraResource}` : '');
  }
  return url;
};
export function cutNsStartClusterId({ namespace, clusterId }) {
  return namespace.replace(new RegExp(`^${clusterId}-`), '');
}
export function reduceNs(namesapce) {
  let newNs = namesapce;
  /// #if project
  //业务侧ns eg: cls-xxx-ns 需要去除前缀
  if (newNs) {
    newNs = newNs.startsWith('global') ? newNs.split('-').splice(1).join('-') : newNs.split('-').splice(2).join('-');
  }
  /// #endif
  return newNs;
}

export function reverseReduceNs(clusterId: string, namespace: string) {
  let newNs = namespace;
  /// #if project
  //业务侧ns eg: cls-xxx-ns 需要去除前缀
  if (newNs) {
    newNs = `${clusterId}-${newNs}`;
  }
  /// #endif
  return newNs;
}
