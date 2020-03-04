import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { Event, Resource } from '../models';

let defaultNamespace = reduceToPayload(ActionType.E_NamespaceSelection, 'default');
/// #if tke
defaultNamespace = reduceToPayload(ActionType.E_NamespaceSelection, 'default');
/// #endif
/// #if project
defaultNamespace = reduceToPayload(ActionType.E_NamespaceSelection, '');
/// #endif

const TempReducer = combineReducers({
  workloadType: reduceToPayload(ActionType.E_WorkloadType, ''),

  namespaceSelection: defaultNamespace,

  workloadQuery: generateQueryReducer({
    actionType: ActionType.E_QueryWorkloadList
  }),

  workloadList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.E_FetchWorkloadList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  workloadSelection: reduceToPayload(ActionType.E_WorkloadSelection, ''),

  eventQuery: generateQueryReducer({
    actionType: ActionType.E_QueryEventList
  }),

  eventList: generateFetcherReducer<RecordSet<Event>>({
    actionType: ActionType.E_FetchEventList,
    initialData: {
      recordCount: 0,
      records: [] as Event[]
    }
  }),

  isAutoRenew: reduceToPayload(ActionType.E_IsAutoRenew, true)
});

export const ResourceEventReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearResourceEvent) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
