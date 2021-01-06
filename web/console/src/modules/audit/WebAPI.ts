import { OperationResult, QueryState, RecordSet, uuid } from '@tencent/ff-redux';
import { Base64 } from 'js-base64';

import { resourceConfig } from '../../../config/resourceConfig';
import {
  reduceK8sQueryString,
  reduceK8sRestfulPath,
  reduceNetworkRequest,
  reduceNetworkWorkflow
} from '../../../helpers';
import { Method } from '../../../helpers/reduceNetwork';
import { RequestParams, Resource, ResourceFilter, ResourceInfo } from '../common/models';
import { Audit, AuditFilter, AuditFilterConditionValues } from './models';

const tips = seajs.require('tips');

const enum Action {
  ListFieldValues = 'listFieldValues/',
  ListBlockClusters = 'listBlockClusters',
  GetStoreConfig = 'getStoreConfig',
  ConfigTest = 'configTest',
  UpdateStoreConfig = 'updateStoreConfig',
}

// 返回标准操作结果
function operationResult<T>(target: T[] | T, error?: any): OperationResult<T>[] {
  if (target instanceof Array) {
    return target.map(x => ({ success: !error, target: x, error }));
  }
  return [{ success: !error, target: target as T, error }];
}

/** 访问凭证相关 */
export async function fetchAuditList(query: QueryState<AuditFilter>) {
  const { search, paging, filter } = query;
  const { pageIndex, pageSize } = paging;
  const queryObj = {
    pageIndex,
    pageSize,
    ...filter
  };
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  const apiKeyResourceInfo: ResourceInfo = resourceConfig()['audit'];
  const url = reduceK8sRestfulPath({
    resourceInfo: apiKeyResourceInfo,
    specificName: 'list/'
  });
  // const url = '/apis/audit.tkestack.io/v1/events/list/';
  const params: RequestParams = {
    method: Method.get,
    url: url + queryString
  };

  const response = await reduceNetworkRequest(params);

  let auditList: any[] = [];
  let totalCount: number = 0;
  try {
    console.log('fetchAuditList response is: ', response);
    if (response.code === 0) {
      auditList = response.data.items;
      totalCount = response.data.total;
    }
  } catch (error) {
    if (+error.response.status !== 404) {
      throw error;
    }
  }

  const result: RecordSet<Audit> = {
    recordCount: totalCount,
    records: auditList
  };

  return result;
}

/**
 * 获取查询条件数据
 */
export async function fetchAuditFilterCondition() {
  return fetchAuditEvents(Action.ListFieldValues);
}

/**
 * 获取audit相关信息
 * anAction可以是：listFieldValues(), listBlockClusters(列出被屏蔽集群列表)
 * @param anAction
 */
export async function fetchAuditEvents(anAction: Action) {
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['audit'];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: anAction });
    const response = await reduceNetworkRequest({
      method: 'GET',
      url
    });
    if (response.code === 0) {
      return operationResult(response.data);
    } else {
      return operationResult('', response);
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
    return operationResult('', error.response);
  }
}

/**
 * 获取审计配置等。Action可以是：getStoreConfig(获取当前es配置)
 * @param anAction
 */
export async function fetchAuditRecord(anAction: Action) {
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['audit'];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: anAction });
    const response = await reduceNetworkRequest({
      method: 'GET',
      url
    });
    if (response.code === 0) {
      return response.data;
    } else {
      return operationResult('', response);
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
    return operationResult('', error.response);
  }
}

/**
 * 获取审计ES配置
 */
export async function getStoreConfig() {
  return fetchAuditRecord(Action.GetStoreConfig);
}

/**
 * 连接检测
 */
export async function configTest(payload) {
  return postAuditEvents(Action.ConfigTest, { elasticSearch: payload })
}

export async function updateStoreConfig(payload) {
  return postAuditEvents(Action.UpdateStoreConfig, { elasticSearch: payload })
}

export async function postAuditEvents(anAction: Action, payload) {
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['audit'];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: anAction });
    const response = await reduceNetworkRequest({
      method: 'POST',
      url,
      data: payload
    });
    if (response.code === 0) {
      return response.data;
    } else {
      return operationResult('', response);
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
    return operationResult('', error.response);
  }
}
