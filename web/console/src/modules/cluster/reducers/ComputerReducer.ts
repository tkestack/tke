import { generateWorkflowReducer, createFFListReducer } from '@tencent/ff-redux';
import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet, uuid } from '@tencent/qcloud-lib';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';
import * as ActionType from '../constants/ActionType';
import { Computer, ComputerFilter, Resource } from '../models';
import { FFReduxActionName } from '../constants/Config';

const TempReducer = combineReducers({
  computer: createFFListReducer<Computer, ComputerFilter>(FFReduxActionName.COMPUTER),

  createComputerWorkflow: generateWorkflowReducer({
    actionType: ActionType.CreateComputer
  }),

  deleteComputer: generateWorkflowReducer({ actionType: ActionType.DeleteComputer }),

  batchUnScheduleComputer: generateWorkflowReducer({
    actionType: ActionType.BatchUnScheduleComputer
  }),

  batchTurnOnSchedulingComputer: generateWorkflowReducer({
    actionType: ActionType.BatchTurnOnSchedulingComputer
  }),

  drainComputer: generateWorkflowReducer({
    actionType: ActionType.DrainComputer
  }),

  updateNodeLabel: generateWorkflowReducer({
    actionType: ActionType.UpdateNodeLabel
  }),

  labelEdition: reduceToPayload(ActionType.UpdateLabelEdition, {
    id: uuid(),
    labels: [],
    originLabel: {},
    computerName: ''
  }),

  updateNodeTaint: generateWorkflowReducer({
    actionType: ActionType.UpdateNodeTaint
  }),

  taintEdition: reduceToPayload(ActionType.UpdateTaintEdition, {
    id: uuid(),
    taints: [],
    computerName: ''
  }),

  computerPodList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.FetchComputerPodList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  computerPodQuery: generateQueryReducer({
    actionType: ActionType.QueryComputerPodList
  })
});

export const ComputerReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearComputer) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
