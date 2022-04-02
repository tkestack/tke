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
import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import { initValidator } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { ExternalTrafficPolicy, SessionAffinity } from '../constants/Config';
import { initPortsMap } from '../constants/initState';
import { Resource } from '../models';

const TempReducer = combineReducers({
  serviceName: reduceToPayload(ActionType.S_ServiceName, ''),

  v_serviceName: reduceToPayload(ActionType.SV_ServiceName, initValidator),

  description: reduceToPayload(ActionType.S_Description, ''),

  v_description: reduceToPayload(ActionType.SV_Description, initValidator),

  namespace: reduceToPayload(ActionType.S_Namespace, ''),

  v_namespace: reduceToPayload(ActionType.SV_Namespace, initValidator),

  communicationType: reduceToPayload(ActionType.S_CommunicationType, 'ClusterIP'),

  portsMap: reduceToPayload(ActionType.S_UpdatePortsMap, [initPortsMap]),

  isOpenHeadless: reduceToPayload(ActionType.S_IsOpenHeadless, false),

  selector: reduceToPayload(ActionType.S_Selector, []),

  isShowWorkloadDialog: reduceToPayload(ActionType.S_IsShowWorkloadDialog, false),

  workloadType: reduceToPayload(ActionType.S_WorkloadType, 'deployment'),

  workloadQuery: generateQueryReducer({
    actionType: ActionType.S_QueryWorkloadList
  }),

  workloadList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.S_FetchWorkloadList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  workloadSelection: reduceToPayload(ActionType.S_WorkloadSelection, []),

  externalTrafficPolicy: reduceToPayload(ActionType.S_ChooseExternalTrafficPolicy, ExternalTrafficPolicy.Cluster),

  sessionAffinity: reduceToPayload(ActionType.S_ChoosesessionAffinity, SessionAffinity.None),

  sessionAffinityTimeout: reduceToPayload(ActionType.S_InputsessionAffinityTimeout, 30),

  v_sessionAffinityTimeout: reduceToPayload(ActionType.SV_sessionAffinityTimeout, initValidator),

  vmiIsEnable: reduceToPayload(ActionType.S_VMI_IsEnable, false)
});

export const ServiceEditReducer = (state, action) => {
  let newState = state;
  // 销毁创建服务页面
  if (action.type === ActionType.ClearServiceEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
