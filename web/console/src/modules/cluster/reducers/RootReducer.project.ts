/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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
