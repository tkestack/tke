import { RecordSet, uuid } from '@tencent/qcloud-lib';
import { QueryState } from '@tencent/qcloud-redux-query';
import { ResourceFilter, Resource } from './models';
import { RequestParams, NamespaceFilter, ResourceInfo, Namespace } from '../common/models';
import {
  reduceNetworkRequest,
  operationResult,
  reduceK8sRestfulPath,
  Method,
  reduceK8sQueryString,
  requestMethodForAction,
  reduceNetworkWorkflow
} from '../../../helpers';
import { tip } from '@tencent/tea-app/lib/bridge';
import { CreateResource } from '../cluster/models';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { CommonAPI } from '../common';
import { resourceConfig } from '../../../config';

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
    resourceInfo,
    clusterId,
    namespace,
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
  let { clusterId } = filter;
  let namespaceList = [];
  if (!isClearData) {
    // 获取k8s的url
    let url = reduceK8sRestfulPath({ resourceInfo });
    if (search) {
      url = url + '/' + search;
    }
    /** 构建参数 */
    let params: RequestParams = {
      method: Method.get,
      url
    };
    try {
      let response = await reduceNetworkRequest(params, clusterId);
      if (response.code === 0) {
        let list = response.data;
        if (list.items) {
          namespaceList = list.items.map(item => {
            return {
              id: uuid(),
              namespace: item.metadata.name
            };
          });
        } else {
          namespaceList.push({
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
    let { clusterId, namespace, mode, jsonData, resourceInfo, isStrategic = true, resourceIns } = resources[0];
    let url = reduceK8sRestfulPath({ resourceInfo, namespace, specificName: resourceIns, clusterId });
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
