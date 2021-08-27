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

import { createFFListReducer, generateWorkflowReducer, reduceToPayload } from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { Resource } from '../models/Resource';
import { router } from '../router';

interface ResourceFilter {}
export const RootReducer = combineReducers({
  route: router.getReducer(),
  channel: createFFListReducer<Resource, ResourceFilter>('channel'),
  template: createFFListReducer<Resource, ResourceFilter>('template'),
  receiver: createFFListReducer<Resource, ResourceFilter>('receiver'),
  receiverGroup: createFFListReducer<Resource, ResourceFilter>('receiverGroup'),
  resourceDeleteWorkflow: generateWorkflowReducer({ actionType: ActionType.DeleteResource }),
  modifyResourceFlow: generateWorkflowReducer({ actionType: ActionType.ModifyResource }),

  /**
   * 判断是否为国际版
   */
  isI18n: reduceToPayload(ActionType.isI18n, false)
});
