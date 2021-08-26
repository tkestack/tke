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

import { createFFListReducer, generateWorkflowReducer, RecordSet, reduceToPayload, uuid } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { Computer, ComputerFilter, Resource } from '../models';

const TempReducer = combineReducers({
  computer: createFFListReducer<Computer, ComputerFilter>(FFReduxActionName.COMPUTER),

  machine: createFFListReducer<Computer, ComputerFilter>(FFReduxActionName.MACHINE),

  createComputerWorkflow: generateWorkflowReducer({
    actionType: ActionType.CreateComputer,
  }),

  deleteComputer: generateWorkflowReducer({ actionType: ActionType.DeleteComputer }),

  batchUnScheduleComputer: generateWorkflowReducer({
    actionType: ActionType.BatchUnScheduleComputer,
  }),

  batchTurnOnSchedulingComputer: generateWorkflowReducer({
    actionType: ActionType.BatchTurnOnSchedulingComputer,
  }),

  drainComputer: generateWorkflowReducer({
    actionType: ActionType.DrainComputer,
  }),

  updateNodeLabel: generateWorkflowReducer({
    actionType: ActionType.UpdateNodeLabel,
  }),

  labelEdition: reduceToPayload(ActionType.UpdateLabelEdition, {
    id: uuid(),
    labels: [],
    originLabel: {},
    computerName: '',
  }),

  updateNodeTaint: generateWorkflowReducer({
    actionType: ActionType.UpdateNodeTaint,
  }),

  taintEdition: reduceToPayload(ActionType.UpdateTaintEdition, {
    id: uuid(),
    taints: [],
    computerName: '',
  }),

  computerPodList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.FetchComputerPodList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[],
    },
  }),

  computerPodQuery: generateQueryReducer({
    actionType: ActionType.QueryComputerPodList,
  }),

  isShowMachine: reduceToPayload(ActionType.IsShowMachine, false),

  deleteMachineResouceIns: reduceToPayload(ActionType.FetchDeleteMachineResouceIns, null),
});

export const ComputerReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearComputer) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
