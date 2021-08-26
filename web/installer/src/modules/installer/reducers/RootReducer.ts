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

import { generateFetcherReducer, generateWorkflowReducer, reduceToPayload, ReduxAction, uuid } from '@tencent/ff-redux';

import { Record } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { initEdit } from './initState';

export const RootReducer = combineReducers({
  step: reduceToPayload(ActionType.StepNext, 'step2'),

  cluster: generateFetcherReducer<Record<any>>({
    actionType: ActionType.FetchCluster,
    initialData: {
      record: {
        config: {},
        progress: {}
      },
      auth: {
        isAuthorized: true
      }
    }
  }),

  isVerified: reduceToPayload(ActionType.VerifyLicense, -1),

  licenseConfig: reduceToPayload(ActionType.GetLicenseConfig, {}),

  clusterProgress: generateFetcherReducer<Record<any>>({
    actionType: ActionType.FetchProgress,
    initialData: {
      record: {},
      auth: {
        isAuthorized: true
      }
    }
  }),

  editState: (state = Object.assign({}, initEdit, { id: uuid() }), action: any) => {
    if (action.type === ActionType.UpdateEdit) {
      return Object.assign({}, state, action.payload);
    } else {
      return state;
    }
  },

  createCluster: generateWorkflowReducer({
    actionType: ActionType.CreateCluster
  })
});
