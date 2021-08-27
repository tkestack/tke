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
import { FFListModel, WorkflowState } from '@tencent/ff-redux';

import { RouteState } from '../../../../helpers/Router';
import { Resource, ResourceFilter } from './Resource';

type ResourceOpWorkflow = WorkflowState<Resource, {}>;

export interface RootState {
  /**
   * 路由
   */
  route?: RouteState;

  channel?: FFListModel<Resource, ResourceFilter>;
  template?: FFListModel<Resource, ResourceFilter>;
  receiver?: FFListModel<Resource, ResourceFilter>;
  receiverGroup?: FFListModel<Resource, ResourceFilter>;
  resourceDeleteWorkflow?: ResourceOpWorkflow;
  modifyResourceFlow?: ResourceOpWorkflow;
  /** 是否为国际版 */
  isI18n?: boolean;
}
