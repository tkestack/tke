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
import { AlarmRecord, AlarmRecordFilter } from './models';

const tips = seajs.require('tips');

// 返回标准操作结果
function operationResult<T>(target: T[] | T, error?: any): OperationResult<T>[] {
    if (target instanceof Array) {
        return target.map(x => ({ success: !error, target: x, error }));
    }
    return [{ success: !error, target: target as T, error }];
}

/** 告警记录 */
export async function fetchAlarmRecord(
        query: QueryState<AlarmRecordFilter>,
        options: {
            continueToken?: string;
        }
    ) {
    const { search, paging, filter } = query;
    let { pageIndex, pageSize: limit } = paging;
    let { continueToken = undefined } = options;
    let queryObj = {
        limit,
    };
    if (filter.clusterID || search) {
        queryObj['fieldSelector'] = {};
    }
    if (filter.clusterID) {
        queryObj['fieldSelector']['spec.clusterID'] = filter.clusterID;
    }
    if (search) {
        delete queryObj.limit;
        queryObj['fieldSelector']['spec.alarmPolicyName'] = search;
        continueToken = undefined;
    }
    if (continueToken) {
        queryObj = JSON.parse(
            JSON.stringify(
                Object.assign({}, queryObj, { continue: continueToken ? continueToken : undefined })
            )
        );
    }
    const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
    const apiKeyResourceInfo: ResourceInfo = resourceConfig()['alarmRecord'];
    const url = reduceK8sRestfulPath({
        resourceInfo: apiKeyResourceInfo
    });
    const params: RequestParams = {
        method: Method.get,
        url: url + queryString
    };
    let alarmRecord: any[] = [];
    // let totalCount: number = 0;
    let nextContinueToken: string;
    try {
        const response = await reduceNetworkRequest(params);
        if (response.code === 0) {
            const listItems = response.data;

            nextContinueToken = listItems.metadata && listItems.metadata.continue ? listItems.metadata.continue : undefined;

            if (listItems.items) {
                alarmRecord = listItems.items.map(item => {
                    return Object.assign({}, item, { id: uuid() });
                });
            }
        }
    } catch (error) {
        // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
        if (+error.response.status !== 404) {
            throw error;
        }
    }

    const result: RecordSet<AlarmRecord> = {
        recordCount: alarmRecord.length, // 无意义
        records: alarmRecord,
        continueToken: nextContinueToken
    };

    return result;
}
