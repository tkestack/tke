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
  createFFListReducer,
  generateWorkflowReducer,
  reduceToPayload,
  generateFetcherReducer,
  createFFObjectReducer
} from '@tencent/ff-redux';

import { Cluster, Namespace, ProjectNamespace, Project } from '../models';
import * as ActionType from '../constants/ActionType';
import { InitApiKey, InitRepo, InitChart, InitImage, Default_D_URL } from '../constants/Config';
import { router } from '../router';
import { createValidatorReducer } from '@tencent/ff-validator';
import { ChartGroupValidateSchema } from '../constants/ChartGroupValidateConfig';
import { ChartValidateSchema } from '../constants/ChartValidateConfig';
import { AppValidateSchema } from '../constants/AppValidateConfig';
import {
  initChartGroupCreationState,
  initChartGroupEditorState,
  initUserInfoState,
  initChartEditorState,
  initRemovedChartVersionsState,
  initAppCreationState,
  initCommonUserAssociationState
} from '../constants/initState';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  /** 访问凭证相关 */
  apiKey: createFFListReducer('apiKey'),

  createApiKey: generateWorkflowReducer({
    actionType: ActionType.CreateApiKey
  }),

  apiKeyCreation: reduceToPayload(ActionType.UpdateApiKeyCreation, InitApiKey),

  deleteApiKey: generateWorkflowReducer({
    actionType: ActionType.DeleteApiKey
  }),

  toggleKeyStatus: generateWorkflowReducer({
    actionType: ActionType.ToggleKeyStatus
  }),

  /** 镜像仓库相关 */
  repo: createFFListReducer('repo'),

  createRepo: generateWorkflowReducer({
    actionType: ActionType.CreateRepo
  }),

  repoCreation: reduceToPayload(ActionType.UpdateRepoCreation, InitRepo),

  deleteRepo: generateWorkflowReducer({
    actionType: ActionType.DeleteRepo
  }),

  /** 镜像相关 */
  image: createFFListReducer('image'),

  createImage: generateWorkflowReducer({
    actionType: ActionType.CreateImage
  }),

  imageCreation: reduceToPayload(ActionType.UpdateImageCreation, InitImage),

  deleteImage: generateWorkflowReducer({
    actionType: ActionType.DeleteImage
  }),

  dockerRegistryUrl: generateFetcherReducer({
    actionType: ActionType.FetchDockerRegUrl,
    initialData: Default_D_URL
  }),

  chart: createFFListReducer('chart'),

  chartIns: createFFListReducer('chartIns'),

  createChart: generateWorkflowReducer({
    actionType: ActionType.CreateChart
  }),

  chartCreation: reduceToPayload(ActionType.UpdateChartCreation, InitChart),

  deleteChart: generateWorkflowReducer({
    actionType: ActionType.DeleteChart
  }),

  /** chartGroup */
  chartGroupList: createFFListReducer(ActionType.ChartGroupList),
  chartGroupCreation: reduceToPayload(ActionType.UpdateChartGroupCreationState, initChartGroupCreationState),
  chartGroupEditor: reduceToPayload(ActionType.UpdateChartGroupEditorState, initChartGroupEditorState),
  chartGroupValidator: createValidatorReducer(ChartGroupValidateSchema),
  chartGroupAddWorkflow: generateWorkflowReducer({
    actionType: ActionType.AddChartGroup
  }),
  chartGroupUpdateWorkflow: generateWorkflowReducer({
    actionType: ActionType.UpdateChartGroup
  }),
  chartGroupRemoveWorkflow: generateWorkflowReducer({
    actionType: ActionType.RemoveChartGroup
  }),
  chartGroupRepoUpdateWorkflow: generateWorkflowReducer({
    actionType: ActionType.RepoUpdateChartGroup
  }),
  projectList: createFFListReducer(
    ActionType.ProjectList,
    '',
    (x: Project) => x.spec.displayName,
    (x: Project) => x.metadata.name
  ),
  userInfo: reduceToPayload(ActionType.UpdateUserInfo, initUserInfoState),

  /** chart */
  chartList: createFFListReducer(ActionType.ChartList, null, null, null, {
    query: {
      paging: {
        pageSize: 15
      }
    }
  }),
  // chartCreation: reduceToPayload(ActionType.UpdateChartCreationState, initChartCreationState),
  chartEditor: reduceToPayload(ActionType.UpdateChartEditorState, initChartEditorState),
  chartValidator: createValidatorReducer(ChartValidateSchema),
  chartAddWorkflow: generateWorkflowReducer({
    actionType: ActionType.AddChart
  }),
  chartUpdateWorkflow: generateWorkflowReducer({
    actionType: ActionType.UpdateChart
  }),
  chartRemoveWorkflow: generateWorkflowReducer({
    actionType: ActionType.RemoveChart
  }),
  chartVersionRemoveWorkflow: generateWorkflowReducer({
    actionType: ActionType.RemoveChartVersion
  }),
  removedChartVersions: reduceToPayload(ActionType.RemovedChartVersions, initRemovedChartVersionsState),
  chartDetail: createFFObjectReducer(ActionType.Chart),
  chartInfo: createFFObjectReducer(ActionType.ChartInfo),
  chartVersionFile: createFFObjectReducer(ActionType.ChartVersionFile),
  appCreation: reduceToPayload(ActionType.UpdateAppCreationState, initAppCreationState),
  appValidator: createValidatorReducer(AppValidateSchema),
  appAddWorkflow: generateWorkflowReducer({
    actionType: ActionType.AddApp
  }),
  appDryRun: reduceToPayload(ActionType.UpdateAppDryRunState, {}),

  /** 集群 */
  /** listActions.selectByValue依赖于valueField */
  clusterList: createFFListReducer(
    ActionType.ClusterList,
    '',
    (x: Cluster) => x.spec.displayName,
    (x: Cluster) => x.metadata.name
  ),
  /** 命名空间 */
  namespaceList: createFFListReducer(
    ActionType.NamespaceList,
    '',
    (x: Namespace) => x.metadata.name,
    (x: Namespace) => x.metadata.name
  ),
  projectNamespaceList: createFFListReducer(
    ActionType.ProjectNamespaceList,
    '',
    (x: ProjectNamespace) => x.metadata.name,
    (x: ProjectNamespace) => x.spec.clusterName + '/' + x.spec.namespace
  ),

  /** 关联用户相关 */
  userPlainList: createFFListReducer(ActionType.UserPlainList),
  commonUserAssociation: reduceToPayload(ActionType.UpdateCommonUserAssociation, initCommonUserAssociationState)
});
