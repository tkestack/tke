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
import { WorkflowState, FetcherState, FFListModel, FFObjectModel } from '@tencent/ff-redux';
import { ApiKey, ApiKeyFilter, ApiKeyCreation } from './ApiKey';
import { Repo, RepoFilter, RepoCreation } from './Repo';
import { Image, ImageFilter, ImageCreation } from './Image';
import {
  Chart,
  ChartFilter,
  ChartIns,
  ChartInsFilter,
  ChartCreation,
  ChartEditor,
  ChartVersion,
  ChartDetailFilter,
  RemovedChartVersions,
  ChartInfo,
  ChartInfoFilter,
  ChartVersionFilter
} from './Chart';
import { ChartGroup, ChartGroupFilter, ChartGroupEditor, ChartGroupCreation } from './ChartGroup';
import { UserInfo } from './UserInfo';
import { Project, ProjectFilter } from './Project';
import { Cluster, ClusterFilter } from './Cluster';
import { Namespace, NamespaceFilter, ProjectNamespace, ProjectNamespaceFilter } from './Namespace';
import { UserPlain, CommonUserAssociation } from './CommonUser';
import { AppCreation, App, AppFilter } from './App';
import { RouteState } from '../../../../helpers';
import { Validation, ValidatorModel } from '@tencent/ff-validator';

type ApiKeyWorkflow = WorkflowState<ApiKey, void>;
type ApiKeyCreateWorkflow = WorkflowState<ApiKeyCreation, void>;
type RepoWorkflow = WorkflowState<Repo, void>;
type RepoCreateWorkflow = WorkflowState<RepoCreation, void>;
type ImageWorkflow = WorkflowState<Image, void>;
type ImageCreateWorkflow = WorkflowState<ImageCreation, void>;
type ChartWorkflow = WorkflowState<Chart, void>;
type ChartCreateWorkflow = WorkflowState<ChartCreation, void>;

export interface RootState {
  /** 路由 */
  route?: RouteState;

  /** -------- 访问凭证相关 --------- */

  apiKey?: FFListModel<ApiKey, ApiKeyFilter>;

  /** ApiKey 创建编辑参数 */
  apiKeyCreation?: ApiKeyCreation;

  /** 创建业务工作流 */
  createApiKey?: ApiKeyCreateWorkflow;

  /** 删除工作流 */
  deleteApiKey?: ApiKeyWorkflow;

  /** 启用禁用工作流 */
  toggleKeyStatus?: ApiKeyWorkflow;

  /** --------- 仓库管理相关 --------- */

  repo?: FFListModel<Repo, RepoFilter>;

  /** 创建仓库表单参数 */
  repoCreation?: RepoCreation;

  /** 创建仓库工作流 */
  createRepo?: RepoCreateWorkflow;

  /** 删除仓库工作流 */
  deleteRepo?: RepoWorkflow;

  /** --------- 镜像相关 --------- */

  image?: FFListModel<Image, ImageFilter>;

  /** 创建仓库表单参数 */
  imageCreation?: ImageCreation;

  /** 创建仓库工作流 */
  createImage?: ImageCreateWorkflow;

  /** 删除仓库工作流 */
  deleteImage?: ImageWorkflow;

  /** docker registry */
  dockerRegistryUrl?: FetcherState<string>;

  /** -------- chart group ----- */

  chart?: FFListModel<Chart, ChartFilter>;

  chartIns?: FFListModel<ChartIns, ChartInsFilter>;

  /** 创建仓库表单参数 */
  chartCreation?: ChartCreation;

  /** 创建仓库工作流 */
  createChart?: ChartCreateWorkflow;

  /** 删除仓库工作流 */
  deleteChart?: ChartWorkflow;

  /** 模板仓库 */
  chartGroupList?: FFListModel<ChartGroup, ChartGroupFilter>;
  chartGroupCreation?: ChartGroupCreation;
  chartGroupEditor?: ChartGroupEditor;
  chartGroupValidator?: ValidatorModel;
  chartGroupAddWorkflow?: WorkflowState<ChartGroup, any>;
  chartGroupUpdateWorkflow?: WorkflowState<ChartGroup, any>;
  chartGroupRemoveWorkflow?: WorkflowState<ChartGroup, any>;
  chartGroupRepoUpdateWorkflow?: WorkflowState<ChartGroup, any>;
  projectList?: FFListModel<Project, ProjectFilter>;
  userInfo?: UserInfo;

  /** 模板 */
  chartList?: FFListModel<Chart, ChartFilter>;
  // chartCreation?: ChartCreation;
  chartEditor?: ChartEditor;
  chartValidator?: ValidatorModel;
  chartAddWorkflow?: WorkflowState<Chart, any>;
  chartUpdateWorkflow?: WorkflowState<Chart, ChartDetailFilter>;
  chartRemoveWorkflow?: WorkflowState<Chart, ChartFilter>;
  chartVersionRemoveWorkflow?: WorkflowState<ChartVersion, ChartVersionFilter>;
  removedChartVersions?: RemovedChartVersions;
  chartDetail?: FFObjectModel<Chart, ChartDetailFilter>;
  chartInfo?: FFObjectModel<ChartInfo, ChartInfoFilter>;
  chartVersionFile?: FFObjectModel<any, ChartVersionFilter>;
  appCreation?: AppCreation;
  appValidator?: ValidatorModel;
  appAddWorkflow?: WorkflowState<App, AppFilter>;
  appDryRun?: App;

  /** 集群 */
  clusterList?: FFListModel<Cluster, ClusterFilter>;
  /** 命名空间 */
  namespaceList?: FFListModel<Namespace, NamespaceFilter>;
  projectNamespaceList?: FFListModel<ProjectNamespace, ProjectNamespaceFilter>;

  /** 关联用户相关，单独设置，不赋予任何场景相关的命名 */
  userPlainList?: FFListModel<UserPlain, void>;
  commonUserAssociation?: CommonUserAssociation;
}
