import { combineReducers } from 'redux';

import { reduceToPayload } from '@tencent/ff-redux';

import { initValidator } from '../../common';
import * as ActionType from '../constants/ActionType';

export const PeEditReducer = combineReducers({
  esAddress: reduceToPayload(ActionType.EsAddress, ''),

  v_esAddress: reduceToPayload(ActionType.V_EsAddress, initValidator),

  indexName: reduceToPayload(ActionType.IndexName, ''),

  v_indexName: reduceToPayload(ActionType.V_IndexName, initValidator),

  esUsername: reduceToPayload(ActionType.EsUsername, ''),

  esPassword: reduceToPayload(ActionType.EsPassword, '')
});
