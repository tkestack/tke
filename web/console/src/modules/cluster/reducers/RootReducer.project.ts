import { combineReducers } from 'redux';

import {
    createFFListReducer, generateWorkflowReducer, RecordSet, reduceToPayload
} from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { Namespace, Resource } from '../models';
import { router } from '../router';
import { SubReducer } from './SubReducer';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  projectNamespaceQuery: generateQueryReducer({
    actionType: ActionType.QueryProjectNamespace
  }),

  projectNamespaceList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.FetchProjectNamespace,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  cluster: createFFListReducer(FFReduxActionName.CLUSTER),

  clusterVersion: reduceToPayload(ActionType.ClusterVersion, '1.16'),

  namespaceList: generateFetcherReducer<RecordSet<Namespace>>({
    actionType: ActionType.FetchNamespaceList,
    initialData: {
      recordCount: 0,
      records: [] as Namespace[]
    }
  }),

  namespaceQuery: generateQueryReducer({
    actionType: ActionType.QueryNamespaceList
  }),

  namespaceSelection: reduceToPayload(ActionType.SelectNamespace, ''),

  projectList: reduceToPayload(ActionType.InitProjectList, []),

  projectSelection: reduceToPayload(ActionType.ProjectSelection, ''),

  region: createFFListReducer(FFReduxActionName.REGION),

  subRoot: SubReducer,

  mode: reduceToPayload(ActionType.ChangeMode, 'create'),

  isShowTips: reduceToPayload(ActionType.IsShowTips, false),

  isI18n: reduceToPayload(ActionType.isI18n, false)
});
