import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import { initValidator } from '../../common/models/Validation';
import * as ActionType from '../constants/ActionType';

const TempReducer = combineReducers({
  isOpen: reduceToPayload(ActionType.IsOpenPE, true),

  esAddress: reduceToPayload(ActionType.EsAddress, ''),

  v_esAddress: reduceToPayload(ActionType.V_EsAddress, initValidator),

  indexName: reduceToPayload(ActionType.IndexName, ''),

  v_indexName: reduceToPayload(ActionType.V_IndexName, initValidator),

  esUsername: reduceToPayload(ActionType.EsUsername, ''),

  esPassword: reduceToPayload(ActionType.EsPassword, '')
});

export const PeEditReducer = (state, action) => {
  let newState = state;
  // 销毁设置事件持久化的界面
  if (action.type === ActionType.ClearPeEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
