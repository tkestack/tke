import { QueryState, RecordSet, uuid } from '@tencent/ff-redux';

import { resourceConfig } from '../../../../config/resourceConfig';
import { reduceK8sRestfulPath } from '../../../../helpers';
import { Method, reduceNetworkRequest } from '../../../../helpers/reduceNetwork';
import { RequestParams, Resource } from '../../../modules/common';
import { ResourceInfo } from '../../common/models/ResourceInfo';
import { Namespace, ResourceFilter } from '../models';

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
