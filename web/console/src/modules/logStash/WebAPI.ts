import { QueryState, RecordSet, uuid } from '@tencent/ff-redux';
import { tip } from '@tencent/tea-app/lib/bridge';
import { t } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../config';
import {
    Method, operationResult, reduceK8sQueryString, reduceK8sRestfulPath, reduceNetworkRequest,
    reduceNetworkWorkflow, requestMethodForAction
} from '../../../helpers';
import { CreateResource } from '../cluster/models';
import { Namespace, NamespaceFilter, RequestParams, ResourceInfo, Cluster } from '../common/models';
import { Resource, ResourceFilter } from './models';

/**
 * 校验当前的日志采集器的名称是否正确
 */
export async function checkStashNameIsExist(
  clusterVersion: string,
  logStashName: string,
  clusterId: string,
  regionId,
  namespace: string
) {
  let resourceInfo = resourceConfig(clusterVersion)['logcs'];
  let url = reduceK8sRestfulPath({
    isSpecialNamespace: true,
    resourceInfo,
    clusterId,
    namespace: namespace.replace(new RegExp(`^${clusterId}-`), ''),
    specificName: logStashName
  });
  let params: RequestParams = {
    method: Method.get,
    url
  };
  try {
    let response = await reduceNetworkRequest(params, clusterId);
    return response.code === 0;
  } catch (error) {
    return false;
  }
}

/**
 * namesapce列表的查询
 * @param query: QueryState<NamespaceFilter> namespace列表的查询
 * @param resourceInfo: ResourceInfo 当前namespace查询api的配置
 */
export async function fetchNamespaceList(
  query: QueryState<NamespaceFilter>,
  resourceInfo: ResourceInfo,
  isClearData: boolean
) {
  let { filter, search } = query;
  let { clusterId, projectName } = filter;
  let namespaceList = [];
  if (!isClearData) {
    // 获取k8s的url
    let url = projectName ? reduceK8sRestfulPath({ resourceInfo, specificName: projectName, extraResource: 'namespaces' }) : reduceK8sRestfulPath({ resourceInfo });
    if (search) {
      url = url + '/' + search;
    }
    /** 构建参数 */
    let params: RequestParams = {
      method: Method.get,
      url
    };
    // 如果是业务侧的话，
    let getCluster = item => {
      if (!projectName) {
        return {};
      }
      let clusterInfo: Cluster = {
        id: uuid(),
        metadata: {
          name: item.spec.clusterName,
        },
        spec: {
          displayName: item.spec.clusterDisplayName,
        },
        status: {
          version: item.spec.clusterVersion,
          phase: 'Running', // TODO: 让namespace接口返回集群的phase
        }
      };
      return clusterInfo;
    };
    try {
      let response = await reduceNetworkRequest(params, clusterId);
      if (response.code === 0) {
        let list = response.data;
        if (list.items) {
          namespaceList = list.items.map(item => {
            return {
              clusterId,
              cluster: getCluster(item),
              id: uuid(),
              namespace: item.metadata.name
            };
          });
        } else {
          namespaceList.push({
            clusterId,
            cluster: getCluster(list),
            id: uuid(),
            name: list.metadata.name
          });
        }
      }
    } catch (error) {
      // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
      if (error.code !== 'ResourceNotFound') {
        throw error;
      }
    }
  }

  const result: RecordSet<Namespace> = {
    recordCount: namespaceList.length,
    records: namespaceList
  };
  return result;
}

/**
 * Resource列表的查询
 * @param query:    Resource 的查询过滤条件
 * @param resourceInfo:ResourceInfo 资源的相关配置
 * @param isClearData:  是否清空数据
 * @param k8sQueryObj: any  是否有queryString
 * @param isNeedDes: boolean    是否需要降序展示
 */
export async function fetchResourceList(
  query: QueryState<ResourceFilter>,
  resourceInfo: ResourceInfo,
  isClearData: boolean = false,
  k8sQueryObj: any = {},
  isNeedDes: boolean = false
) {
  let { filter, search } = query,
    { namespace, clusterId, regionId, isCanFetchResourceList } = filter;

  let resourceList = [];
  if (!isClearData && isCanFetchResourceList) {
    let k8sUrl = reduceK8sRestfulPath({ resourceInfo, namespace });
    // 如果是搜索过滤的话，需要进行一下拼接
    if (search) {
      k8sUrl = k8sUrl + '/' + search;
    }

    // 这里是去拼接，是否需要在k8s url后面拼接一些queryString
    let queryString = reduceK8sQueryString(k8sQueryObj);
    let url = k8sUrl + queryString;

    // 构建参数
    let params: RequestParams = {
      method: Method.get,
      url,
      apiParams: {
        module: 'tke',
        interfaceName: 'ForwardRequest',
        regionId: regionId || 1,
        restParams: {
          Method: Method.get,
          Path: url,
          Version: '2018-05-25'
        },
        opts: {
          tipErr: false
        }
      }
    };

    try {
      let response = await reduceNetworkRequest(params, clusterId);

      if (response.code === 0) {
        let listItems = JSON.parse(response.data.ResponseBody);
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
      } else {
        // tips.error(zh2enMessage(response, _isI18n), 1000);
      }
    } catch (error) {
      // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
      if (error.code !== 'ResourceNotFound') {
        throw error;
      }
    }
  }

  const result: RecordSet<Resource> = {
    recordCount: resourceList.length,
    records: isNeedDes && resourceList.length > 1 ? resourceList.reverse() : resourceList
  };

  return result;
}

/**
 * 创建、更新日志采集规则
 */
export async function modifyLogStash(resources: CreateResource[], regionId: number) {
  try {
    let { clusterId, logAgentName, namespace, mode, jsonData, resourceInfo, isStrategic = true, resourceIns } = resources[0];
    let url = reduceK8sRestfulPath({ resourceInfo, namespace, specificName: resourceIns, clusterId, logAgentName, isSpecialNamespace: window.location.href.includes('tkestack-project') });
    // 构建参数
    let method = requestMethodForAction(mode);
    let params: RequestParams = {
      method: method,
      url,
      data: jsonData,
      apiParams: {
        module: 'tke',
        interfaceName: 'ForwardRequest',
        regionId: regionId || 1,
        restParams: {
          Method: method,
          Path: url,
          Version: '2018-05-25',
          RequestBody: jsonData
        }
      }
    };
    if (mode === 'update') {
      params.userDefinedHeader = {
        'Content-Type': isStrategic ? 'application/strategic-merge-patch+json' : 'application/merge-patch+json'
      };
    }

    let response = await reduceNetworkRequest(params, clusterId);
    let operateTip = mode === 'create' ? '创建成功' : '更新成功';
    if (response.code === 0) {
      tip.success(t(operateTip), 2000);
      return operationResult(resources);
    } else {
      return operationResult(resources, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resources, reduceNetworkWorkflow(error));
  }
}

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
