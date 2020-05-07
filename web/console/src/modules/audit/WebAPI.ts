import { OperationResult, QueryState, RecordSet, uuid } from '@tencent/ff-redux';

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
        resourceInfo: apiKeyResourceInfo, specificName: 'list/'
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
    try {
        const resourceInfo: ResourceInfo = resourceConfig()['audit'];
        const url = reduceK8sRestfulPath({ resourceInfo, specificName: 'listFieldValues/' });
        // const url = '/apis/audit.tkestack.io/v1/events/listFieldValues/';
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
