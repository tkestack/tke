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
import { FetcherState, QueryState, RecordSet, WorkflowState } from '@tencent/ff-redux';

import { RouteState } from '../../../../helpers';
import { Resource, ResourceFilter } from '../../common/models/Resource';
import { DetailState, HelmCreation, ListState } from './';
import { Namespace } from './Namespace';

export interface RootState {
  /** 路由 */
  route?: RouteState;

  listState?: ListState;

  detailState?: DetailState;

  /**创建Helm所选参数*/
  helmCreation?: HelmCreation;

  /**是否显示提示 */
  isShowTips?: boolean;

  /**业务侧逻辑*/

  /** namespace列表 */
  namespaceList?: FetcherState<RecordSet<Namespace>>;

  /** namespace查询条件 */
  namespaceQuery?: QueryState<ResourceFilter>;

  /** namespace selection */
  namespaceSelection?: string;

  /** namespacesetQuery */
  projectNamespaceQuery?: QueryState<ResourceFilter>;

  /** namespaceset */
  projectNamespaceList?: FetcherState<RecordSet<Resource>>;

  /** projectList */
  projectList?: any[];

  /** projectSelection */
  projectSelection?: string;
  /**end */
}
