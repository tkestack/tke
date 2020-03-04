import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { Resource } from '../models';

const TempReducer = combineReducers({
  resourceQuery: generateQueryReducer({
    actionType: ActionType.QueryResourceList
  }),

  resourceList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.FetchResourceList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  resourceSelection: reduceToPayload(ActionType.SelectResource, []),

  resourceMultipleSelection: reduceToPayload(ActionType.SelectMultipleResource, []),

  resourceDeleteSelection: reduceToPayload(ActionType.SelectDeleteResource, [])
});

export const ResourceReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearResource) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
