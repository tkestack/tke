import { combineReducers } from 'redux';
import { reduceToPayload } from '@tencent/qcloud-lib';
import * as ActionType from '../constants/ActionType';
import { initValidator } from '../../common';

export const PeEditReducer = combineReducers({
  esAddress: reduceToPayload(ActionType.EsAddress, ''),

  v_esAddress: reduceToPayload(ActionType.V_EsAddress, initValidator),

  indexName: reduceToPayload(ActionType.IndexName, ''),

  v_indexName: reduceToPayload(ActionType.V_IndexName, initValidator)
});
