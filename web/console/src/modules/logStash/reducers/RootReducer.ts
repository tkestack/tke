import { Namespace } from 'react-i18next';
import { combineReducers } from 'redux';

import {
    createFFListReducer, generateWorkflowReducer, RecordSet, reduceToPayload
} from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import { Cluster, NamespaceFilter, Region } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { initLogDaemonsetStatus, initRegionInfo } from '../constants/initState';
import { Log } from '../models';
import { LogDaemonset } from '../models/LogDaemonset';
import { router } from '../router';
import { LogStashEditReducer } from './LogStashEditReducer';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  regionQuery: generateQueryReducer({
    actionType: ActionType.QueryRegion
  }),

  regionList: generateFetcherReducer<RecordSet<Region>>({
    actionType: ActionType.FetchRegion,
    initialData: {
      recordCount: 0,
      records: [] as Region[]
    }
  }),

  regionSelection: reduceToPayload(ActionType.SelectRegion, initRegionInfo),

  clusterQuery: generateQueryReducer({
    actionType: ActionType.QueryClusterList
  }),

  clusterList: generateFetcherReducer<RecordSet<Cluster>>({
    actionType: ActionType.FetchClusterList,
    initialData: {
      recordCount: 0,
      records: [] as Cluster[]
    }
  }),

  clusterSelection: reduceToPayload(ActionType.SelectCluster, []),

  clusterVersion: reduceToPayload(ActionType.ClusterVersion, '1.16'),

  namespaceList: generateFetcherReducer<RecordSet<Namespace>>({
    actionType: ActionType.FetchNamespaceList,
    initialData: {
      recordCount: 0,
      records: [] as Namespace[]
    }
  }),

  namespaceQuery: generateQueryReducer<NamespaceFilter>({
    actionType: ActionType.QueryNamespaceList
  }),

  namespaceSelection: reduceToPayload(ActionType.NamespaceSelection, ''),

  logQuery: generateQueryReducer({
    actionType: ActionType.QueryLogList
  }),

  logList: generateFetcherReducer<RecordSet<Log>>({
    actionType: ActionType.FetchLogList,
    initialData: {
      recordCount: 0,
      records: [] as Log[]
    }
  }),

  logSelection: reduceToPayload(ActionType.SelectLog, []),

  logDaemonsetQuery: generateQueryReducer({ actionType: ActionType.QueryLogDaemonset }),

  logDaemonset: generateFetcherReducer<RecordSet<LogDaemonset>>({
    actionType: ActionType.FetchLogDaemonset,
    initialData: {
      recordCount: 0,
      records: [] as LogDaemonset[]
    }
  }),

  isOpenLogStash: reduceToPayload(ActionType.IsOpenLogStash, false),

  isDaemonsetNormal: reduceToPayload(ActionType.IsDaemonsetNormal, initLogDaemonsetStatus),

  authorizeOpenLogFlow: generateWorkflowReducer({
    actionType: ActionType.AuthorizeOpenLog
  }),

  modifyLogStashFlow: generateWorkflowReducer({
    actionType: ActionType.ModifyLogStashFlow
  }),

  inlineDeleteLog: generateWorkflowReducer({
    actionType: ActionType.InlineDeleteLog
  }),

  isFetchDoneSpecificLog: reduceToPayload(ActionType.IsFetchDoneSpecificLog, false),

  logStashEdit: LogStashEditReducer,

  openAddon: createFFListReducer(FFReduxActionName.OPENADDON)
});
