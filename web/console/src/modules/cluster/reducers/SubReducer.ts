/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import { generateWorkflowReducer, RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { initAllcationRatioEdition } from '../constants/initState';
import { SubRouter } from '../models';
import { ComputerReducer } from './ComputerReducer';
import { ConfigMapEditReducer } from './ConfigMapEditReducer';
import { DetailResourceReducer } from './DetailResourceReducer';
import { LbcfEditReducer } from './LbcfEditReducer';
import { NamespaceEditReducer } from './NamespaceEditReducer';
import { ResourceDetailReducer } from './ResourceDetailReducer';
import { ResourceEventReducer } from './ResourceEventReducer';
import { ResourceLogReducer } from './ResourceLogReducer';
import { ResourceReducer } from './ResourceReducer';
import { SecretEditReducer } from './SecretEditReducer';
import { ServiceEditReducer } from './ServiceEditReducer';
import { WorkloadEditReducer } from './WorkloadEditReducer';

const TempReducer = combineReducers({
  computerState: ComputerReducer,

  clusterAllocationRatioEdition: reduceToPayload(
    ActionType.UpdateClusterAllocationRatioEdition,
    initAllcationRatioEdition
  ),

  updateClusterAllocationRatio: generateWorkflowReducer({
    actionType: ActionType.UpdateClusterAllocationRatio
  }),

  subRouterQuery: generateQueryReducer({
    actionType: ActionType.QuerySubRouterList
  }),

  subRouterList: generateFetcherReducer<RecordSet<SubRouter>>({
    actionType: ActionType.FetchSubRouterList,
    initialData: {
      recordCount: 0,
      records: [] as SubRouter[]
    }
  }),

  mode: reduceToPayload(ActionType.SelectMode, 'list'),

  applyResourceFlow: generateWorkflowReducer({
    actionType: ActionType.ApplyResource
  }),

  applyDifferentInterfaceResourceFlow: generateWorkflowReducer({
    actionType: ActionType.ApplyDifferentInterfaceResource
  }),

  modifyResourceFlow: generateWorkflowReducer({
    actionType: ActionType.ModifyResource
  }),

  modifyMultiResourceWorkflow: generateWorkflowReducer({
    actionType: ActionType.ModifyMultiResource
  }),

  deleteResourceFlow: generateWorkflowReducer({
    actionType: ActionType.DeleteResource
  }),

  updateResourcePart: generateWorkflowReducer({
    actionType: ActionType.UpdateResourcePart
  }),

  updateMultiResource: generateWorkflowReducer({
    actionType: ActionType.UpdateMultiResource
  }),

  resourceName: reduceToPayload(ActionType.InitResourceName, ''),

  resourceInfo: reduceToPayload(ActionType.InitResourceInfo, {}),

  detailResourceOption: DetailResourceReducer,

  resourceOption: ResourceReducer,

  resourceDetailState: ResourceDetailReducer,

  serviceEdit: ServiceEditReducer,

  namespaceEdit: NamespaceEditReducer,

  workloadEdit: WorkloadEditReducer,

  secretEdit: SecretEditReducer,

  cmEdit: ConfigMapEditReducer,

  lbcfEdit: LbcfEditReducer,

  resourceLogOption: ResourceLogReducer,

  resourceEventOption: ResourceEventReducer,

  isNeedFetchNamespace: reduceToPayload(ActionType.IsNeedFetchNamespace, true),

  isNeedExistedLb: reduceToPayload(ActionType.IsNeedExistedLb, false),
  addons: reduceToPayload(ActionType.FetchClusterAddons, {})
});

export const SubReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearSubRoot) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
