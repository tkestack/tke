import { combineReducers } from 'redux';

import {
    createFFListReducer, generateWorkflowReducer, RecordSet, reduceToPayload
} from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import { Resource } from '../../common';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { router } from '../router';
import { PeEditReducer } from './PeEditReducer';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  region: createFFListReducer(FFReduxActionName.REGION),

  cluster: createFFListReducer(FFReduxActionName.CLUSTER),

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
