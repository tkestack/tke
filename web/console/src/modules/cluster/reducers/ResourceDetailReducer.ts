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

import { createFFListReducer, generateWorkflowReducer, RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { Event, Pod, Replicaset, Resource, ResourceFilter } from '../models';

/** ==== start 日志的相关处理 ============ */
const logOptionReducer = combineReducers({
  podName: reduceToPayload(ActionType.PodName, ''),

  containerName: reduceToPayload(ActionType.ContainerName, ''),

  logFile: reduceToPayload(ActionType.LogFile, 'stdout'),

  tailLines: reduceToPayload(ActionType.TailLines, '100'),

  isAutoRenew: reduceToPayload(ActionType.IsAutoRenewPodLog, false)
});
/** ==== start 日志的相关处理 ============ */

const TempReducer = combineReducers({
  resourceDetailInfo: createFFListReducer(FFReduxActionName.Resource_Detail_Info),

  yamlList: generateFetcherReducer<RecordSet<string>>({
    actionType: ActionType.FetchYaml,
    initialData: {
      recordCount: 0,
      records: [] as string[]
    }
  }),

  event: createFFListReducer<Event, ResourceFilter>(FFReduxActionName.DETAILEVENT),

  rsQuery: generateQueryReducer({
    actionType: ActionType.QueryRsList
  }),

  rsList: generateFetcherReducer<RecordSet<Replicaset>>({
    actionType: ActionType.FetchRsList,
    initialData: {
      recordCount: 0,
      records: [] as Replicaset[]
    }
  }),

  rollbackResourceFlow: generateWorkflowReducer({
    actionType: ActionType.RollBackResource
  }),

  removeTappPodFlow: generateWorkflowReducer({
    actionType: ActionType.RemoveTappPod
  }),

  rsSelection: reduceToPayload(ActionType.RsSelection, []),

  podQuery: generateQueryReducer({
    actionType: ActionType.QueryPodList
  }),

  podList: generateFetcherReducer<RecordSet<Pod>>({
    actionType: ActionType.FetchPodList,
    initialData: {
      recordCount: 0,
      records: [] as Pod[]
    }
  }),

  podFilterInNode: reduceToPayload(ActionType.PodFilterInNode, {}),

  containerList: reduceToPayload(ActionType.FetchContainerList, []),

  podSelection: reduceToPayload(ActionType.PodSelection, []),

  deletePodFlow: generateWorkflowReducer({
    actionType: ActionType.DeletePod
  }),

  updateGrayTappFlow: generateWorkflowReducer({
    actionType: ActionType.UpdateGrayTapp
  }),

  editTappGrayUpdate: reduceToPayload(ActionType.W_TappGrayUpdate, { containers: [] }),

  isShowLoginDialog: reduceToPayload(ActionType.IsShowLoginDialog, false),

  logQuery: generateQueryReducer({
    actionType: ActionType.QueryLogList
  }),

  logList: generateFetcherReducer<RecordSet<string>>({
    actionType: ActionType.FetchLogList,
    initialData: {
      recordCount: 0,
      records: [] as string[]
    }
  }),

  logAgent: reduceToPayload(ActionType.PodLogAgent, {}),

  logHierarchy: reduceToPayload(ActionType.PodLogHierarchy, []),

  logContent: reduceToPayload(ActionType.PodLogContent, ''),

  logOption: logOptionReducer,

  secretQuery: generateQueryReducer({
    actionType: ActionType.QuerySecretList
  }),

  secretList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.FetchSecretList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  secretSelection: reduceToPayload(ActionType.SecretSelection, []),

  modifyNamespaceSecretFlow: generateWorkflowReducer({
    actionType: ActionType.ModifyNamespaceSecret
  })
});

export const ResourceDetailReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearResourceDetail) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
