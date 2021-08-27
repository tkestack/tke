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
    createFFListReducer, generateWorkflowReducer, RecordSet, reduceToPayload
} from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import { Resource } from '../../common';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { router } from '../router';
import { PeEditReducer } from './PeEditReducer';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  region: createFFListReducer(FFReduxActionName.REGION),

  cluster: createFFListReducer(FFReduxActionName.CLUSTER),

  peList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.FetchPeList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  peQuery: generateQueryReducer({
    actionType: ActionType.QueryPeList
  }),

  peSelection: reduceToPayload(ActionType.SelectPe, []),

  peEdit: PeEditReducer,

  resourceInfo: reduceToPayload(ActionType.InitResourceInfo, {}),

  modifyPeFlow: generateWorkflowReducer({
    actionType: ActionType.ModifyPeFlow
  }),

  deletePeFlow: generateWorkflowReducer({
    actionType: ActionType.DeletePeFlow
  })
});
