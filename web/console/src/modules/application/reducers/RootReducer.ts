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
import { combineReducers } from 'redux';

import {
  createFFListReducer,
  generateWorkflowReducer,
  reduceToPayload,
  generateFetcherReducer,
  createFFObjectReducer
} from '@tencent/ff-redux';

import { Cluster, Namespace, ProjectNamespace, Project } from '../models';
import * as ActionTypes from '../constants/ActionTypes';
import { router } from '../router';
import { createValidatorReducer } from '@tencent/ff-validator';
import { AppValidateSchema } from '../constants/AppValidateConfig';
import { initAppCreationState, initAppEditorState, initResourceList, initHistoryList } from '../constants/initState';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  /** 集群 */
  /** listActions.selectByValue依赖于valueField */
  clusterList: createFFListReducer(
    ActionTypes.ClusterList,
    '',
    (x: Cluster) => x.spec.displayName,
    (x: Cluster) => x.metadata.name
  ),
  /** 命名空间 */
  namespaceList: createFFListReducer(
    ActionTypes.NamespaceList,
    '',
    (x: Namespace) => x.metadata.name,
    (x: Namespace) => x.metadata.name
  ),
  projectNamespaceList: createFFListReducer(
    ActionTypes.ProjectNamespaceList,
    '',
    (x: ProjectNamespace) => x.metadata.name,
    (x: ProjectNamespace) => x.spec.clusterName + '/' + x.spec.namespace
  ),

  /** 模板 */
  chartList: createFFListReducer(ActionTypes.ChartList, null, null, null, {
    query: {
      paging: {
        pageSize: 15
      }
    }
  }),
  chartInfo: createFFObjectReducer({ actionName: ActionTypes.ChartInfo }),
  chartGroupList: createFFListReducer(ActionTypes.ChartGroupList),

  /** 业务 */
  projectList: createFFListReducer(
    ActionTypes.ProjectList,
    '',
    (x: Project) => x.metadata.name,
    (x: Project) => x.metadata.name
  ),

  /** 应用 */
  appList: createFFListReducer(ActionTypes.AppList),
  appCreation: reduceToPayload(ActionTypes.UpdateAppCreationState, initAppCreationState),
  appEditor: reduceToPayload(ActionTypes.UpdateAppEditorState, initAppEditorState),
  appDryRun: reduceToPayload(ActionTypes.UpdateAppDryRunState, {}),
  appValidator: createValidatorReducer(AppValidateSchema),
  appAddWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.AddApp
  }),
  appUpdateWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.UpdateApp
  }),
  appRemoveWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.RemoveApp
  }),
  appResource: createFFObjectReducer({ actionName: ActionTypes.AppResource }),
  resourceList: reduceToPayload(ActionTypes.ResourceList, initResourceList),
  appHistory: createFFObjectReducer({ actionName: ActionTypes.AppHistory }),
  historyList: reduceToPayload(ActionTypes.HistoryList, initHistoryList),
  appRollbackWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.RollbackApp
  })
});
