import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { Resource } from '../models';

const TempReducer = combineReducers({
  detailResourceName: reduceToPayload(ActionType.InitDetailResourceName, ''),

  detailResourceInfo: reduceToPayload(ActionType.InitDetailResourceInfo, {}),

  detailResourceSelection: reduceToPayload(ActionType.SelectDetailResourceSelection, ''),

  detailResourceList: reduceToPayload(ActionType.InitDetailResourceList, []),

  detailDeleteResourceSelection: reduceToPayload(ActionType.SelectDetailDeleteResourceSelection, '')
});

export const DetailResourceReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearResource) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
