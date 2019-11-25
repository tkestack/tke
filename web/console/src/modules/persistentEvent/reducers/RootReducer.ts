import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';
import { router } from '../router';
import * as ActionType from '../constants/ActionType';
import { PeEditReducer } from './PeEditReducer';
import { generateWorkflowReducer } from '@tencent/qcloud-redux-workflow';
import { createListReducer } from '@tencent/redux-list';
import { FFReduxActionName } from '../constants/Config';
import { Resource } from '../../common';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  region: createListReducer(FFReduxActionName.REGION),

  cluster: createListReducer(FFReduxActionName.CLUSTER),

  peList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.FetchPeList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  peQuery: generateQueryReducer({
    actionType: ActionType.QueryPeList
  }),

  peSelection: reduceToPayload(ActionType.SelectPe, []),

  peEdit: PeEditReducer,

  resourceInfo: reduceToPayload(ActionType.InitResourceInfo, {}),

  modifyPeFlow: generateWorkflowReducer({
    actionType: ActionType.ModifyPeFlow
  }),

  deletePeFlow: generateWorkflowReducer({
    actionType: ActionType.DeletePeFlow
  })
});
