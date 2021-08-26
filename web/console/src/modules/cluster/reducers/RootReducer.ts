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

import { Cluster } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { initClusterCreationState, initDialogState } from '../constants/initState';
import { Namespace } from '../models';
import { router } from '../router';
import { CreateICReducer } from './CreateICReducer';
import { SubReducer } from './SubReducer';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  region: createFFListReducer(FFReduxActionName.REGION),

  cluster: createFFListReducer(FFReduxActionName.CLUSTER),

  clustercredential: reduceToPayload(ActionType.FetchClustercredential, {
    name: '',
    clusterName: '',
    caCert: '',
    token: ''
  }),

  updateClusterToken: generateWorkflowReducer({
    actionType: ActionType.UpdateClusterToken
  }),

  clusterVersion: reduceToPayload(ActionType.ClusterVersion, '1.16'),

  clusterInfoQuery: generateQueryReducer({
    actionType: ActionType.QueryClusterInfo
  }),

  clusterInfoList: generateFetcherReducer<RecordSet<Cluster>>({
    actionType: ActionType.FetchClusterInfo,
    initialData: {
      recordCount: 0,
      records: [] as Cluster[]
    }
  }),

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

  namespaceSelection: reduceToPayload(ActionType.SelectNamespace, 'default'),

  subRoot: SubReducer,

  mode: reduceToPayload(ActionType.ChangeMode, 'create'),

  dialogState: reduceToPayload(ActionType.UpdateDialogState, initDialogState),

  deleteClusterFlow: generateWorkflowReducer({
    actionType: ActionType.DeleteCluster
  }),

  clusterCreationState: reduceToPayload(ActionType.UpdateclusterCreationState, initClusterCreationState),

  createClusterFlow: generateWorkflowReducer({
    actionType: ActionType.CreateCluster
  }),

  createIC: CreateICReducer,
  createICWorkflow: generateWorkflowReducer({
    actionType: ActionType.CreateIC
  }),

  modifyClusterName: generateWorkflowReducer({
    actionType: ActionType.ModifyClusterNameWorkflow
  })
});
