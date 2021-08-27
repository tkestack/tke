/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
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

  projectList: reduceToPayload(ActionType.InitProjectList, []),

  projectSelection: reduceToPayload(ActionType.ProjectSelection, ''),

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
