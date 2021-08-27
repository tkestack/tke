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
import {
    createFFListActions, extend, FetchOptions, generateFetcherActionCreator,
    generateWorkflowActionCreator, OperationTrigger
} from '@tencent/ff-redux';

import * as ActionTypes from '../constants/ActionTypes';
import { RootState, Audit, AuditFilter } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
    noCache: false
};

const FFModelAuditActions = createFFListActions<Audit, AuditFilter>({
    actionName: ActionTypes.FetchAuditList,
    fetcher: async (query, getState: GetState) => {
        let response = await WebAPI.fetchAuditList(query);
        return response;
    },
    getRecord: (getState: GetState) => {
        return getState().auditList;
    },
    onFinish: (record, dispatch: Redux.Dispatch) => {
        if (record.data.recordCount) {
            let isNotNeedPoll = record.data.records.filter(item => item.status['phase'] === 'Terminating').length === 0;

            if (isNotNeedPoll) {
                dispatch(FFModelAuditActions.clearPolling());
            }
        }
    }
});

/**
 * 获取策略
 */
const getAuditFilterCondition = generateFetcherActionCreator({
    actionType: ActionTypes.FetchAuditFilterCondition,
    fetcher: async (getState: GetState, options: FetchOptions, dispatch) => {
        let result = await WebAPI.fetchAuditFilterCondition();
        return result;
    }
});

const restActions = {
    poll: () => {
        return async (dispatch: Redux.Dispatch, getState: GetState) => {
            dispatch(
                auditActions.polling({
                    delayTime: 5000
                })
            );
        };
    },
    getAuditFilterCondition,
};

export const auditActions = extend({}, FFModelAuditActions, restActions);
