import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';
import { generateWorkflowReducer } from '@tencent/ff-redux';
import * as ActionType from '../constants/ActionType';
import * as initState from './initState';
import { HelmHistory } from '../models';

const TempReducer = combineReducers({
  helm: reduceToPayload(ActionType.FetchHelm, null),
  isRefresh: reduceToPayload(ActionType.IsRefresh, false),
  historyQuery: generateQueryReducer({
    actionType: ActionType.QueryHistory
  }),
  histories: generateFetcherReducer<RecordSet<HelmHistory>>({
    actionType: ActionType.FetchHistory,
    initialData: {
      recordCount: 0,
      records: [] as HelmHistory[]
    }
  })
});

export const DetailReducer = (inputState, action) => {
  let state = inputState;
  // 销毁详情页面
  if (action.type === ActionType.ClearDetail) {
    state = undefined;
  }
  return TempReducer(state, action);
};
