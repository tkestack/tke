import { router } from './../router';
import { generateWorkflowReducer } from '@tencent/qcloud-redux-workflow';
import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';
import * as ActionType from '../constants/ActionType';
import { Namespace, Resource } from '../models';
import { SubReducer } from './SubReducer';
import { FFReduxActionName } from '../constants/Config';
import { createListReducer } from '@tencent/redux-list';

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

  cluster: createListReducer(FFReduxActionName.CLUSTER),

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

  region: createListReducer(FFReduxActionName.REGION),

  subRoot: SubReducer,

  mode: reduceToPayload(ActionType.ChangeMode, 'create'),

  isShowTips: reduceToPayload(ActionType.IsShowTips, false),

  isI18n: reduceToPayload(ActionType.isI18n, false)
});
